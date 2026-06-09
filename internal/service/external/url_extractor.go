// internal/service/external/url_extractor.go
package external

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"

	"YoudaoNoteLm/pkg/logger"
	"go.uber.org/zap"
)

// URLExtractor 网页 URL 抓取器
// 使用 goquery 提取正文，转为 Markdown
type URLExtractor struct {
	client *http.Client
}

// NewURLExtractor 创建 URL 抓取器
func NewURLExtractor() *URLExtractor {
	return &URLExtractor{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// ExtractFromURL 抓取网页并转为 Markdown
func (e *URLExtractor) ExtractFromURL(url string) (string, error) {
	resp, err := e.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求URL失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP状态码 %d", resp.StatusCode)
	}

	return e.ExtractFromReader(resp.Body)
}

// ExtractFromReader 从 HTML Reader 提取正文并转为 Markdown
func (e *URLExtractor) ExtractFromReader(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	// 移除无关标签
	doc.Find("script, style, nav, footer, header, aside, .ad, .advertisement, .sidebar").Remove()

	// 尝试提取正文：优先 <article>，其次 <main>，最后 <body>
	var articleHTML string
	var htmlErr error
	if article := doc.Find("article"); article.Length() > 0 {
		articleHTML, htmlErr = article.Html()
	} else if main := doc.Find("main"); main.Length() > 0 {
		articleHTML, htmlErr = main.Html()
	} else if body := doc.Find("body"); body.Length() > 0 {
		articleHTML, htmlErr = body.Html()
	}
	if htmlErr != nil {
		return "", fmt.Errorf("提取HTML内容失败: %w", htmlErr)
	}

	if articleHTML == "" {
		return "", fmt.Errorf("无法提取网页正文")
	}

	// HTML → Markdown
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())

	markdown, err := converter.ConvertString(articleHTML)
	if err != nil {
		return "", fmt.Errorf("HTML转Markdown失败: %w", err)
	}

	// 清理多余空行
	markdown = cleanMarkdown(markdown)

	if strings.TrimSpace(markdown) == "" {
		return "", fmt.Errorf("未能从网页提取到有效内容")
	}

	logger.Info("URL抓取成功", zap.Int("content_length", len(markdown)))
	return markdown, nil
}

// CheckURL 预检 URL 是否可抓取
func (e *URLExtractor) CheckURL(url string) (bool, string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, "", fmt.Errorf("创建预检请求失败: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("预检请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, "", fmt.Errorf("HTTP状态码 %d", resp.StatusCode)
	}

	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	supportedTypes := []string{
		"text/html",
		"text/plain",
		"application/xhtml+xml",
		"application/xml",
		"text/xml",
	}

	for _, st := range supportedTypes {
		if strings.Contains(ct, st) {
			return true, ct, nil
		}
	}

	// 无 Content-Type 时默认尝试
	if ct == "" {
		return true, "unknown", nil
	}

	return false, ct, fmt.Errorf("不支持的内容类型: %s", ct)
}

// cleanMarkdown 清理 Markdown 中的多余空行
func cleanMarkdown(s string) string {
	for strings.Contains(s, "\n\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n\n", "\n\n\n")
	}
	return strings.TrimSpace(s)
}
