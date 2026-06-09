package external

import (
	"testing"
)

func TestMarkitdownClientImplementsDocumentConverter(t *testing.T) {
	// 创建 markitdown 客户端
	client := NewMarkitdownClient("http://localhost:8080")

	// 验证是否实现了 DocumentConverter 接口
	_, ok := client.(DocumentConverter)
	if !ok {
		t.Error("markitdownClient 未实现 DocumentConverter 接口")
	}

	// 验证是否实现了 MarkitdownClient 接口
	_, ok = client.(MarkitdownClient)
	if !ok {
		t.Error("markitdownClient 未实现 MarkitdownClient 接口")
	}

	// 验证 SupportedFormats 方法
	formats := client.SupportedFormats()
	if len(formats) == 0 {
		t.Error("SupportedFormats 返回空列表")
	}

	t.Logf("markitdownClient 支持 %d 种格式", len(formats))
}