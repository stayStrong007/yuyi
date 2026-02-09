# 羽译 (Yuyi) - 开发任务清单

> **目标**: 实现 MVP v1.0 核心功能  
> **技术栈**: Wails v2 + Go + Svelte + TypeScript  
> **最后更新**: 2026-01-20

---

## 📦 Phase 0: 环境准备

- [x] **安装核心依赖库**
  ```bash
  go get golang.design/x/hotkey
  go get github.com/getlantern/systray
  ```
- [ ] **了解 Wails 运行时 API**
  - 文档: https://wails.io/docs/reference/runtime/intro

---

## 🖥️ Phase 1: 系统托盘与窗口行为 ✅

### 1.1 配置隐藏启动与系统托盘 (System Tray)

- [x] 修改 `main.go` 配置项：
  - [x] 设置 `StartHidden: true` 启动时隐藏窗口
  - [x] 设置 `Frameless: true` 无边框窗口
  - [x] 调整窗口尺寸为 `Width: 500, Height: 300`
- [x] 实现系统托盘功能：
  - [x] 添加托盘图标 (`build/windows/icon.ico`)
  - [x] 实现托盘右键菜单：「显示/隐藏」「退出」
  - [x] 点击托盘图标显示/隐藏窗口
- [x] 窗口行为优化：
  - [ ] 窗口失去焦点时自动隐藏 (`OnBlur` 事件) - 待实现
  - [x] 窗口始终居中显示
  - [x] 窗口置顶显示 (`AlwaysOnTop: true`)

**涉及文件**: `main.go`, `app.go`

---

## ⌨️ Phase 2: 全局热键 (Global Hotkey) ✅

### 2.1 注册全局热键 (Ctrl+Space)

- [x] 在 `app.go` 中实现热键监听：
  - [x] 引入 `golang.design/x/hotkey`
  - [x] 在 `startup()` 中启动热键监听 goroutine
  - [x] 实现 `Ctrl+Space` 热键注册
- [x] 热键触发行为：
  - [x] 窗口隐藏时 → 显示窗口 + 发送事件通知前端
  - [x] 窗口显示时 → 隐藏窗口
- [x] 错误处理：
  - [x] 热键注册失败时的提示
  - [x] 应用退出时注销热键

**涉及文件**: `app.go`  
**参考库**: `golang.design/x/hotkey`

---

## 📋 Phase 3: 智能聚焦与前端联调 ✅

### 3.1 实现智能聚焦功能

- [x] 后端事件通知：
  - [x] 显示窗口时发送 `window_show_event` 事件
  - [x] 前端监听事件并聚焦输入框
- [x] 前端 UI 实现：
  - [x] 创建输入框 `<textarea id="input-box">`
  - [x] 监听 `window_show_event` 事件
  - [x] 收到事件后自动聚焦输入框
- [x] 用户交互：
  - [x] 用户可直接打字输入
  - [x] 用户可按 `Ctrl+V` 粘贴（依赖系统原生粘贴）

**涉及文件**: `app.go`, `frontend/src/App.svelte`

---

## 🌐 Phase 4: 翻译 API 接入 (进行中)

### 4.1 Mock 翻译模式 ✅

- [x] 实现 Mock 翻译方法：
  ```go
  func (a *App) Translate(text string) []string
  ```
- [x] Mock 逻辑：返回 3 个模拟候选结果
- [x] 日志输出：在终端打印翻译请求

### 4.2 接入真实翻译 API（待实现）

- [ ] 设计配置结构：
  ```go
  type Config struct {
      APIProvider string // "openai" | "deepl" | "youdao"
      APIKey      string
      APIEndpoint string
      CandidateCount int // 候选译文数量 1-5
  }
  ```
- [ ] 实现翻译接口抽象：
  - [ ] 定义 `Translator` 接口
  - [ ] 实现 OpenAI 翻译 Provider
  - [ ] 实现 DeepL 翻译 Provider（可选）
- [ ] 智能语言识别：
  - [ ] 自动判断输入是中文还是英文
  - [ ] 决定翻译方向

**涉及文件**: `app.go`, 新建 `translator.go`, `config.go`

---

## 🎨 Phase 5: 前端 UI 改造 (进行中)

### 5.1 重构 UI 界面

- [x] 替换默认 Hello World 界面
- [x] 实现主界面布局：
  - [x] 输入框区域（支持多行）
  - [ ] 语言方向指示器（中→英 / 英→中）
  - [x] 翻译结果列表（1-5 个候选）
  - [x] 底部快捷键提示栏
- [x] 实现交互逻辑：
  - [x] `↑/↓` 键切换选中项
  - [x] `Enter` 复制当前选中
  - [x] `1-3` 快速选择
  - [x] `Esc` 清空输入
  - [ ] `Tab` 回到输入框
- [x] 实现视觉效果：
  - [ ] 窗口淡入淡出动画
  - [ ] 复制成功 Toast 提示
  - [x] 加载状态指示器
  - [x] 错误状态展示

**涉及文件**: `frontend/src/App.svelte`, `frontend/src/style.css`

---

## ⚙️ Phase 6: 配置与设置

### 6.1 本地配置存储

- [ ] 实现配置文件读写：
  - [ ] 配置文件路径：`%APPDATA%/Yuyi/config.json`
  - [ ] 便携模式：exe 同目录下 `config.json`
- [ ] 配置项：
  - [ ] API Provider 选择
  - [ ] API Key
  - [ ] 热键自定义
  - [ ] 候选数量 (1-5)
  - [ ] 自动复制开关

### 6.2 设置界面

- [ ] 实现设置页面 UI
- [ ] 托盘菜单「设置」打开设置页

**涉及文件**: 新建 `config.go`, `frontend/src/Settings.svelte`

---

## 🧪 Phase 7: 测试与优化

- [ ] 功能测试：
  - [ ] 热键响应测试
  - [ ] 剪贴板读取测试
  - [ ] 翻译 API 调用测试
  - [ ] 窗口行为测试
- [ ] 性能优化：
  - [ ] 启动速度优化
  - [ ] 内存占用检查
- [ ] 打包发布：
  - [ ] 配置应用图标
  - [ ] `wails build` 生成 release 版本
  - [ ] 测试单文件运行

---

## 📚 参考资源

- [Wails 官方文档](https://wails.io/docs/introduction)
- [Wails 运行时 API](https://wails.io/docs/reference/runtime/intro)
- [golang.design/x/hotkey](https://github.com/nicholaskh/hotkey)
- [golang.design/x/clipboard](https://github.com/nicholaskh/clipboard)
- [Svelte 官方文档](https://svelte.dev/docs)

---

## ✅ 完成标准

MVP v1.0 发布条件：
1. ✅ 单文件 `.exe`，体积 < 10MB
2. ✅ 系统托盘常驻，双击 Alt 唤醒
3. ✅ 自动读取剪贴板并翻译
4. ✅ 显示 3 个候选译文
5. ✅ 快捷键操作流畅
6. ✅ 复制后自动隐藏
