package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAITranslator 使用 OpenAI 兼容接口的翻译器
type OpenAITranslator struct {
	APIKey string
	APIUrl string
	Model  string
}

// ChatMessage OpenAI 消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest OpenAI 请求结构
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatResponse OpenAI 响应结构
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// SystemPrompt 翻译系统提示词
const SystemPrompt = `You are a professional translator. User will provide text. You must detect the language. If it is English, translate to Chinese. If it is Chinese, translate to English. Output only the translation results. Provide 3 different versions/styles of translation separated by newline. Do not include any numbering, explanations, or extra formatting.`

// NewOpenAITranslator 创建新的翻译器实例
func NewOpenAITranslator(apiKey, apiUrl, model string) *OpenAITranslator {
	return &OpenAITranslator{
		APIKey: apiKey,
		APIUrl: apiUrl,
		Model:  model,
	}
}

// Translate 执行翻译
func (t *OpenAITranslator) Translate(text string) []string {
	// 构建请求体
	reqBody := ChatRequest{
		Model: t.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: SystemPrompt},
			{Role: "user", Content: text},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return []string{fmt.Sprintf("Error: 请求序列化失败 - %v", err)}
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", t.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return []string{fmt.Sprintf("Error: 创建请求失败 - %v", err)}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.APIKey)

	// 发送请求（设置超时）
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Printf("[Translator] 发送请求到: %s\n", t.APIUrl)
	resp, err := client.Do(req)
	if err != nil {
		return []string{fmt.Sprintf("Error: 网络请求失败 - %v", err)}
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{fmt.Sprintf("Error: 读取响应失败 - %v", err)}
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return []string{fmt.Sprintf("Error: API 返回错误 (HTTP %d) - %s", resp.StatusCode, string(body))}
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return []string{fmt.Sprintf("Error: 解析响应失败 - %v", err)}
	}

	// 检查 API 错误
	if chatResp.Error != nil {
		return []string{fmt.Sprintf("Error: API 错误 - %s", chatResp.Error.Message)}
	}

	// 提取翻译结果
	if len(chatResp.Choices) == 0 {
		return []string{"Error: API 未返回翻译结果"}
	}

	content := chatResp.Choices[0].Message.Content
	fmt.Printf("[Translator] 收到响应: %s\n", content)

	// 按换行符分割结果
	results := splitResults(content)

	// 确保至少返回一个结果
	if len(results) == 0 {
		return []string{content}
	}

	return results
}

// splitResults 分割翻译结果（处理多种换行格式）
func splitResults(content string) []string {
	// 统一换行符
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// 按换行符分割
	lines := strings.Split(content, "\n")

	var results []string
	for _, line := range lines {
		// 清理每行
		line = strings.TrimSpace(line)

		// 移除可能的序号前缀（如 "1.", "1)", "1:", "1、"）
		if len(line) > 2 {
			runes := []rune(line)
			if len(runes) > 1 && runes[0] >= '1' && runes[0] <= '9' {
				if runes[1] == '.' || runes[1] == ')' || runes[1] == ':' || runes[1] == '、' {
					line = strings.TrimSpace(string(runes[2:]))
				}
			}
		}

		// 跳过空行
		if line != "" {
			results = append(results, line)
		}
	}

	// 限制最多返回 5 个结果
	if len(results) > 5 {
		results = results[:5]
	}

	return results
}
