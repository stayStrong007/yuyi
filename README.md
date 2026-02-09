# 羽译 (Yuyi)

<div align="center">

**轻量级桌面翻译工具**

一个基于 Wails v2 构建的快速、优雅的桌面翻译助手

[功能特性](#功能特性) • [安装使用](#安装使用) • [构建指南](#构建指南) • [配置说明](#配置说明) • [开发](#开发)

</div>

---

## 简介

羽译是一个轻量级的桌面翻译工具，通过全局热键快速唤起，支持 OpenAI 兼容的翻译 API。设计理念是"快速、简洁、不打扰"，让翻译像呼吸一样自然。

## 功能特性

- 🚀 **全局热键** - `Ctrl+Space` 快速唤起/隐藏窗口
- 🎯 **智能翻译** - 支持 OpenAI 兼容 API，自动识别语言方向
- ⌨️ **键盘优先** - 方向键导航，`Enter` 复制，`1-9` 快速选择，`Esc` 关闭
- 🎨 **简洁界面** - 无边框弹窗，始终置顶，自动防抖
- 💾 **配置持久化** - 设置自动保存到用户目录
- 🔒 **隐私安全** - API 密钥本地存储，不上传云端
- 📦 **单文件运行** - 无需安装，开箱即用

## 安装使用

### 下载安装

1. 前往 [Releases](https://github.com/stayStrong007/yuyi/releases) 页面
2. 下载最新版本的 `yuyi.exe`
3. 双击运行即可

### 首次配置

1. 点击系统托盘图标，选择「设置」
2. 填写以下信息：
   - **API Key**: 你的 OpenAI API 密钥
   - **API URL**: API 端点（默认：OpenAI 官方）
   - **Model**: 模型名称（默认：gpt-3.5-turbo）
   - **目标语言**: 翻译目标语言（默认：中文）
3. 点击「保存设置」

### 使用方法

1. 按下 `Ctrl+Space` 唤起翻译窗口
2. 输入要翻译的文本（支持自动翻译）
3. 使用键盘操作：
   - `↑/↓` - 切换选中项
   - `Enter` - 复制当前选中的翻译
   - `1-9` - 快速选择对应序号的翻译
   - `Esc` - 关闭窗口

## 配置说明

配置文件位置：
- **Windows**: `%APPDATA%\yuyi\config.json`
- **macOS**: `~/Library/Application Support/yuyi/config.json`
- **Linux**: `~/.config/yuyi/config.json`

配置文件示例：
```json
{
  "api_key": "sk-...",
  "api_url": "https://api.openai.com/v1/chat/completions",
  "model": "gpt-3.5-turbo",
  "target_lang": "ZH"
}
```

### 支持的 API 提供商

羽译支持所有 OpenAI 兼容的 API，包括但不限于：
- OpenAI 官方 API
- Azure OpenAI
- 国内各大模型厂商（通义千问、文心一言等）
- 自部署的兼容服务

只需修改 `api_url` 和 `api_key` 即可切换。

## 构建指南

### 环境要求

- Go 1.21+
- Node.js 18+
- Wails CLI v2

### 安装依赖

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装项目依赖
cd frontend
npm install
cd ..
```

### 开发模式

```bash
# 启动开发服务器（支持热重载）
wails dev
```

### 构建生产版本

```bash
# 构建可执行文件
wails build

# 构建产物位于 build/bin/ 目录
```

### 前端独立开发

```bash
cd frontend
npm run dev      # 启动 Vite 开发服务器
npm run build    # 构建前端
npm run check    # TypeScript 类型检查
```

## 技术栈

### 后端
- **Go** - 核心逻辑
- **Wails v2** - 桌面应用框架
- **systray** - 系统托盘集成
- **hotkey** - 全局热键监听
- **clipboard** - 剪贴板操作

### 前端
- **Svelte** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具

## 项目结构

```
yuyi/
├── main.go           # Wails 应用入口
├── app.go            # 核心应用逻辑（热键、托盘、窗口管理）
├── config.go         # 配置文件读写
├── translator.go     # 翻译 API 客户端
├── frontend/         # 前端代码
│   ├── src/
│   │   ├── App.svelte    # 主界面组件
│   │   ├── main.ts       # 前端入口
│   │   └── style.css     # 样式
│   └── wailsjs/      # Wails 自动生成的绑定代码
└── build/            # 构建资源（图标、配置）
```

## 开发路线图

查看 [TODO.md](./TODO.md) 了解详细的开发计划。

**已完成：**
- ✅ 全局热键监听
- ✅ 系统托盘集成
- ✅ OpenAI API 翻译
- ✅ 键盘导航
- ✅ 配置持久化

**计划中：**
- [ ] 窗口失焦自动隐藏
- [ ] 淡入淡出动画
- [ ] 翻译历史记录
- [ ] 多语言界面
- [ ] 自定义热键

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交改动 (`git commit -m '添加某个特性'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- [Wails](https://wails.io/) - 优秀的 Go 桌面应用框架
- [Svelte](https://svelte.dev/) - 简洁高效的前端框架
- [OpenAI](https://openai.com/) - 强大的 AI 能力

---

<div align="center">

**如果这个项目对你有帮助，请给个 ⭐️ Star 支持一下！**

</div>
