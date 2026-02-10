package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置结构
type Config struct {
	APIKey     string `json:"api_key"`
	APIUrl     string `json:"api_url"`     // 例如 https://api.openai.com/v1/chat/completions
	Model      string `json:"model"`       // 例如 gpt-3.5-turbo
	TargetLang string `json:"target_lang"` // 默认 "ZH"
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		APIKey:     "",
		APIUrl:     "https://api.openai.com/v1/chat/completions",
		Model:      "gpt-3.5-turbo",
		TargetLang: "ZH",
	}
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取用户配置目录失败: %w", err)
	}

	// 创建 yuyi 文件夹
	yuyiDir := filepath.Join(configDir, "yuyi")
	if err := os.MkdirAll(yuyiDir, 0755); err != nil {
		return "", fmt.Errorf("创建配置目录失败: %w", err)
	}

	return filepath.Join(yuyiDir, "config.json"), nil
}

// LoadConfig 从文件加载配置
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// 如果文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 JSON
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return DefaultConfig(), fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 序列化为 JSON（格式化输出）
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
