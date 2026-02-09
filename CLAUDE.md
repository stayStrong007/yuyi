# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

羽译 (Yuyi) is a lightweight desktop translation tool built with Wails v2 (Go backend + Svelte frontend). It provides a global hotkey-triggered popup window for quick translation via OpenAI-compatible APIs.

## Build Commands

```bash
# Development mode with hot reload
wails dev

# Build production executable
wails build

# Frontend only commands (run from frontend/ directory)
cd frontend && npm install    # Install dependencies
cd frontend && npm run dev    # Run Vite dev server
cd frontend && npm run build  # Build frontend
cd frontend && npm run check  # TypeScript/Svelte type checking
```

## Architecture

### Backend (Go)

- **main.go**: Wails app initialization with window options (frameless, always-on-top, start hidden)
- **app.go**: Core application logic
  - System tray integration (`systray`)
  - Global hotkey listener (`Ctrl+Space` via `golang.design/x/hotkey`)
  - Window state management (show/hide/toggle)
  - Methods exposed to frontend: `Translate()`, `CopyToClipboard()`, `GetConfig()`, `SaveSettings()`, `Hide()`
- **config.go**: JSON config persistence in user config directory (`%APPDATA%/yuyi/config.json`)
- **translator.go**: OpenAI-compatible API client for translation

### Frontend (Svelte + TypeScript)

- **frontend/src/App.svelte**: Single-file component containing:
  - Main translation UI with input and result list
  - Settings panel for API configuration
  - Keyboard navigation (arrows, Enter, 1-9 quick select, Esc to close)
  - Auto-translate with 600ms debounce
- **frontend/wailsjs/**: Auto-generated Wails bindings for Go methods

### Communication Flow

1. User presses `Ctrl+Space` → Go hotkey listener calls `ToggleWindow()`
2. Window shows → Go emits `window_show_event` → Frontend focuses input
3. User types → Frontend debounces → calls `Translate()` → Go translator calls API
4. User selects result → `CopyToClipboard()` → `Hide()`

### Key Dependencies

- Backend: `github.com/wailsapp/wails/v2`, `github.com/getlantern/systray`, `golang.design/x/clipboard`, `golang.design/x/hotkey`
- Frontend: Svelte 3, Vite, TypeScript

## Configuration

Config stored at `%APPDATA%/yuyi/config.json` with fields:
- `api_key`: OpenAI API key
- `api_url`: API endpoint (default: OpenAI)
- `model`: Model name (default: gpt-3.5-turbo)
- `target_lang`: Target language code
