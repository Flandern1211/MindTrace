// internal/service/external/url_extractor_test.go
package external

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLExtractor_ExtractFromReader(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
<html>
<head><title>测试页面</title></head>
<body>
<article>
<h1>标题一</h1>
<p>这是第一段正文内容。</p>
<h2>标题二</h2>
<p>这是第二段正文，包含<strong>加粗</strong>和<em>斜体</em>。</p>
<ul><li>列表项1</li><li>列表项2</li></ul>
</article>
</body>
</html>`

	extractor := NewURLExtractor()
	result, err := extractor.ExtractFromReader(
		&noopCloser{strings.NewReader(htmlContent)},
	)
	if err != nil {
		t.Fatalf("提取失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("提取结果不应为空")
	}

	if !strings.Contains(result, "标题一") {
		t.Error("结果应包含 '标题一'")
	}
	if !strings.Contains(result, "这是第一段正文内容") {
		t.Error("结果应包含正文内容")
	}
	if !strings.Contains(result, "列表项1") {
		t.Error("结果应包含列表项")
	}
	t.Logf("提取结果:\n%s", result)
}

func TestURLExtractor_ExtractFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Mock Page</title></head>
<body>
<article>
<h1>Mock 标题</h1>
<p>Mock 正文内容。</p>
</article>
</body>
</html>`))
	}))
	defer server.Close()

	extractor := NewURLExtractor()
	result, err := extractor.ExtractFromURL(server.URL)
	if err != nil {
		t.Fatalf("URL抓取失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("抓取结果不应为空")
	}

	if !strings.Contains(result, "Mock 标题") {
		t.Error("结果应包含 'Mock 标题'")
	}
	if !strings.Contains(result, "Mock 正文内容") {
		t.Error("结果应包含正文内容")
	}
}

func TestURLExtractor_CheckURL_HTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
	}))
	defer server.Close()

	extractor := NewURLExtractor()
	ok, ct, err := extractor.CheckURL(server.URL)
	if err != nil {
		t.Fatalf("CheckURL失败: %v", err)
	}
	if !ok {
		t.Error("HTML页面应该可抓取")
	}
	if !strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type 应包含 text/html, got %s", ct)
	}
}

func TestURLExtractor_CheckURL_UnsupportedType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(200)
	}))
	defer server.Close()

	extractor := NewURLExtractor()
	ok, _, err := extractor.CheckURL(server.URL)
	if ok {
		t.Error("视频文件不应可抓取")
	}
	if err == nil {
		t.Error("应返回错误")
	}
}

// noopCloser 包装 strings.Reader 为 io.ReadCloser
type noopCloser struct {
	*strings.Reader
}

func (n *noopCloser) Close() error {
	return nil
}
