package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

//go:embed build/windows/icon.ico
var iconData []byte

// App struct
type App struct {
	ctx       context.Context
	isVisible bool    // 窗口是否可见
	config    *Config // 应用配置
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		isVisible: false,
		config:    DefaultConfig(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("加载配置失败:", err)
	}
	a.config = config

	// 初始化剪贴板
	err = clipboard.Init()
	if err != nil {
		fmt.Println("剪贴板初始化失败:", err)
	}

	// 启动系统托盘（在单独的 goroutine 中运行）
	go systray.Run(a.onTrayReady, a.onTrayExit)

	// 启动全局热键监听
	go mainthread.Init(a.startHotkeyListener)
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// DOM 加载完成后的初始化逻辑
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// 退出托盘
	systray.Quit()
}

// onTrayReady 托盘初始化
func (a *App) onTrayReady() {
	// 设置托盘图标（使用嵌入的 icon.ico）
	systray.SetIcon(iconData)
	systray.SetTitle("羽译")
	systray.SetTooltip("羽译 - 轻量级翻译工具")

	// 添加菜单项
	mShow := systray.AddMenuItem("显示/隐藏", "显示或隐藏窗口")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出羽译")

	// 监听菜单点击事件
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				a.ToggleWindow()
			case <-mQuit.ClickedCh:
				runtime.Quit(a.ctx)
				return
			}
		}
	}()
}

// onTrayExit 托盘退出时调用
func (a *App) onTrayExit() {
	// 清理资源
}

// startHotkeyListener 启动全局热键监听
func (a *App) startHotkeyListener() {
	// 注册 Ctrl+Space 热键（翻译工具常用快捷键）
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl}, hotkey.KeySpace)
	err := hk.Register()
	if err != nil {
		fmt.Println("热键注册失败:", err)
		return
	}
	defer hk.Unregister()

	fmt.Println("全局热键已启动，按 Ctrl+Space 唤醒/隐藏窗口")

	for {
		<-hk.Keydown()
		// 按下热键时切换窗口状态
		a.ToggleWindow()
	}
}

// ToggleWindow 切换窗口显示/隐藏状态（可供前端调用）
func (a *App) ToggleWindow() {
	if a.isVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

// ShowWindow 显示窗口
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowCenter(a.ctx) // 居中显示
	a.isVisible = true
	// 通知前端窗口已显示，触发智能聚焦
	runtime.EventsEmit(a.ctx, "window_show_event", nil)
}

// HideWindow 隐藏窗口
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	a.isVisible = false
}

// Hide 隐藏窗口（供前端调用的简化方法）
func (a *App) Hide() {
	a.HideWindow()
}

// CopyToClipboard 将文本写入剪贴板
func (a *App) CopyToClipboard(text string) bool {
	clipboard.Write(clipboard.FmtText, []byte(text))
	fmt.Printf("[Clipboard] 已复制: %s\n", text)
	return true
}

// IsWindowVisible 返回窗口是否可见（供前端查询）
func (a *App) IsWindowVisible() bool {
	return a.isVisible
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Translate 翻译文本并返回多个候选结果
func (a *App) Translate(text string) []string {
	// 空文本返回空数组
	if text == "" {
		return []string{}
	}

	fmt.Printf("[Translate] 收到翻译请求: %s\n", text)

	// 检查 API Key 配置
	if a.config.APIKey == "" {
		return []string{"请先在设置中配置 API Key"}
	}

	// 检查 API URL 配置
	if a.config.APIUrl == "" {
		return []string{"请先在设置中配置 API URL"}
	}

	// 创建翻译器并执行翻译
	translator := NewOpenAITranslator(
		a.config.APIKey,
		a.config.APIUrl,
		a.config.Model,
	)

	results := translator.Translate(text)

	fmt.Printf("[Translate] 返回 %d 个候选结果\n", len(results))
	return results
}

// GetConfig 返回当前配置（供前端调用）
func (a *App) GetConfig() *Config {
	return a.config
}

// SaveSettings 保存设置（供前端调用）
func (a *App) SaveSettings(apiKey, apiUrl, model, targetLang string) bool {
	a.config.APIKey = apiKey
	a.config.APIUrl = apiUrl
	a.config.Model = model
	a.config.TargetLang = targetLang

	if err := SaveConfig(a.config); err != nil {
		fmt.Println("保存配置失败:", err)
		return false
	}

	fmt.Println("[Settings] 配置已更新")
	return true
}
