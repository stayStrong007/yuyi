<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime.js';
  import { Translate, CopyToClipboard, Hide, GetConfig, SaveSettings } from '../wailsjs/go/main/App.js';

  // 主界面状态
  let inputText: string = '';
  let inputBox: HTMLTextAreaElement;
  let results: string[] = [];
  let selectedIndex: number = 0;
  let isLoading: boolean = false;
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  let copyFeedback: string = '';

  // 设置界面状态
  let isSettingsView: boolean = false;
  let settingsForm = {
    apiKey: '',
    apiUrl: '',
    model: '',
    targetLang: ''
  };
  let isSaving: boolean = false;
  let showApiKey: boolean = false;

  // 切换 API Key 显示/隐藏
  function toggleApiKeyVisibility() {
    showApiKey = !showApiKey;
  }

  // 防抖函数：停止输入 600ms 后触发翻译
  function debounce(fn: () => void, delay: number) {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    debounceTimer = setTimeout(fn, delay);
  }

  // 翻译函数
  async function doTranslate() {
    const text = inputText.trim();
    if (!text) {
      results = [];
      return;
    }

    isLoading = true;
    try {
      results = await Translate(text);
      selectedIndex = 0; // 重置选中项
    } catch (err) {
      console.error('翻译失败:', err);
      results = ['翻译出错，请重试'];
    } finally {
      isLoading = false;
    }
  }

  // 监听输入变化，触发防抖翻译
  $: if (inputText !== undefined) {
    debounce(doTranslate, 600);
  }

  // 窗口显示事件处理：智能聚焦
  function onWindowShow() {
    if (inputBox) {
      inputBox.focus();
    }
  }

  // 处理键盘事件
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      // ESC: 清空并隐藏
      inputText = '';
      results = [];
      Hide();
    } else if (event.key === 'ArrowDown') {
      event.preventDefault();
      if (results.length > 0) {
        selectedIndex = (selectedIndex + 1) % results.length;
      }
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      if (results.length > 0) {
        selectedIndex = (selectedIndex - 1 + results.length) % results.length;
      }
    } else if (event.key >= '1' && event.key <= '9') {
      const idx = parseInt(event.key) - 1;
      if (idx < results.length) {
        event.preventDefault();
        copyAndHide(idx);
      }
    } else if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      if (results.length > 0) {
        copyAndHide(selectedIndex);
      }
    }
  }

  // 复制结果并隐藏窗口
  async function copyAndHide(index: number) {
    if (index >= 0 && index < results.length) {
      const text = results[index];
      try {
        await CopyToClipboard(text);
        // 显示复制成功反馈
        showCopyFeedback();
        // 延迟隐藏，让用户看到反馈
        setTimeout(() => {
          inputText = '';
          results = [];
          Hide();
        }, 300);
      } catch (err) {
        console.error('复制失败:', err);
      }
    }
  }

  // 显示复制成功反馈
  function showCopyFeedback() {
    copyFeedback = '已复制 ✓';
    setTimeout(() => {
      copyFeedback = '';
    }, 800);
  }

  // 复制结果到剪贴板（点击时使用）
  async function copyResult(index: number) {
    await copyAndHide(index);
  }

  // 选择结果项
  function selectResult(index: number) {
    selectedIndex = index;
  }

  // 切换设置视图
  function toggleSettings() {
    isSettingsView = !isSettingsView;
    if (isSettingsView) {
      loadSettings();
    }
  }

  // 加载设置
  async function loadSettings() {
    try {
      const config = await GetConfig();
      settingsForm = {
        apiKey: config.api_key || '',
        apiUrl: config.api_url || '',
        model: config.model || '',
        targetLang: config.target_lang || ''
      };
    } catch (err) {
      console.error('加载配置失败:', err);
    }
  }

  // 保存设置
  async function saveSettingsHandler() {
    isSaving = true;
    try {
      const success = await SaveSettings(
        settingsForm.apiKey,
        settingsForm.apiUrl,
        settingsForm.model,
        settingsForm.targetLang
      );
      if (success) {
        showCopyFeedback();
        copyFeedback = '保存成功 ✓';
        isSettingsView = false;
      } else {
        copyFeedback = '保存失败';
      }
    } catch (err) {
      console.error('保存配置失败:', err);
      copyFeedback = '保存失败';
    } finally {
      isSaving = false;
    }
  }

  // 取消设置编辑
  function cancelSettings() {
    isSettingsView = false;
  }

  onMount(() => {
    EventsOn('window_show_event', onWindowShow);
    loadSettings(); // 初始化时加载配置
    if (inputBox) {
      inputBox.focus();
    }
  });

  onDestroy(() => {
    EventsOff('window_show_event');
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
  });
</script>

<main>
  <div class="container">
    <!-- 复制成功反馈 Toast -->
    {#if copyFeedback}
      <div class="toast">{copyFeedback}</div>
    {/if}

    <!-- 顶部拖拽区域 + 设置按钮 -->
    <div class="header">
      <div class="drag-area" style="--wails-draggable: drag;"></div>
      <button class="settings-btn" on:click|stopPropagation={toggleSettings} title="设置">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </button>
    </div>

    {#if isSettingsView}
      <!-- 设置界面 -->
      <div class="settings-view">
        <h3 class="settings-title">API 设置</h3>
        
        <div class="form-group">
          <label for="api-url">API URL</label>
          <input
            id="api-url"
            type="text"
            bind:value={settingsForm.apiUrl}
            placeholder="https://api.openai.com/v1/chat/completions"
          />
        </div>

        <div class="form-group">
          <label for="api-key">API Key</label>
          <div class="input-with-toggle">
            {#if showApiKey}
              <input
                id="api-key"
                type="text"
                bind:value={settingsForm.apiKey}
                placeholder="sk-..."
              />
            {:else}
              <input
                id="api-key"
                type="password"
                bind:value={settingsForm.apiKey}
                placeholder="sk-..."
              />
            {/if}
            <button type="button" class="toggle-visibility" on:click={toggleApiKeyVisibility} title={showApiKey ? '隐藏' : '显示'}>
              {#if showApiKey}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              {:else}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              {/if}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label for="model">Model</label>
          <input
            id="model"
            type="text"
            bind:value={settingsForm.model}
            placeholder="gpt-3.5-turbo"
          />
        </div>

        <div class="form-group">
          <label for="target-lang">目标语言</label>
          <select id="target-lang" bind:value={settingsForm.targetLang}>
            <option value="ZH">中文 (ZH)</option>
            <option value="EN">英文 (EN)</option>
            <option value="JA">日文 (JA)</option>
            <option value="KO">韩文 (KO)</option>
          </select>
        </div>

        <div class="settings-actions">
          <button class="btn btn-secondary" on:click={cancelSettings}>取消</button>
          <button class="btn btn-primary" on:click={saveSettingsHandler} disabled={isSaving}>
            {isSaving ? '保存中...' : '保存'}
          </button>
        </div>
      </div>
    {:else}
      <!-- 主界面 -->
      <!-- 输入区域 -->
      <div class="input-area">
        <textarea
          bind:this={inputBox}
          bind:value={inputText}
          id="input-box"
          class="input-box"
          placeholder="输入或粘贴要翻译的文本..."
          on:keydown={handleKeydown}
          spellcheck="false"
        ></textarea>
      </div>

      <!-- 翻译结果区域 -->
      <div class="result-area">
        {#if isLoading}
          <div class="loading">翻译中...</div>
        {:else if results.length > 0}
          <div class="result-list">
            {#each results as result, index}
              <div
                class="result-item"
                class:selected={index === selectedIndex}
                on:click={() => copyResult(index)}
                on:mouseenter={() => selectResult(index)}
                on:keydown={(e) => e.key === 'Enter' && copyResult(index)}
                role="button"
                tabindex="0"
              >
                <span class="result-index">{index + 1}</span>
                <span class="result-text">{result}</span>
              </div>
            {/each}
          </div>
          <div class="result-hint">↑↓ 切换 | Enter 复制 | 1-3 快选 | Esc 关闭</div>
        {:else if inputText.length > 0}
          <div class="result-hint">输入完成后自动翻译...</div>
        {:else}
          <div class="result-hint">Ctrl+Space 唤醒 | 直接输入或 Ctrl+V 粘贴</div>
        {/if}
      </div>
    {/if}
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    overflow: hidden;
    background: transparent;
  }

  main {
    width: 100vw;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    background: rgba(27, 38, 54, 0.95);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
    position: relative;
  }

  .container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  /* Toast 复制反馈 */
  .toast {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(79, 140, 255, 0.95);
    color: white;
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 500;
    z-index: 1000;
    animation: toast-fade 0.3s ease;
  }

  @keyframes toast-fade {
    from {
      opacity: 0;
      transform: translate(-50%, -50%) scale(0.9);
    }
    to {
      opacity: 1;
      transform: translate(-50%, -50%) scale(1);
    }
  }

  /* 顶部区域：拖拽 + 设置按钮 */
  .header {
    display: flex;
    align-items: center;
    height: 28px;
    background: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .drag-area {
    flex: 1;
    height: 100%;
  }

  .settings-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: color 0.2s ease;
    -webkit-app-region: no-drag;
  }

  .settings-btn:hover {
    color: rgba(255, 255, 255, 0.9);
  }

  /* 设置界面样式 */
  .settings-view {
    padding: 16px;
    overflow-y: auto;
    flex: 1;
  }

  .settings-title {
    margin: 0 0 16px 0;
    color: #ffffff;
    font-size: 16px;
    font-weight: 500;
  }

  .form-group {
    margin-bottom: 14px;
  }

  .form-group label {
    display: block;
    color: rgba(255, 255, 255, 0.7);
    font-size: 12px;
    margin-bottom: 4px;
  }

  .form-group input,
  .form-group select {
    width: 100%;
    padding: 8px 10px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.08);
    color: #ffffff;
    font-size: 13px;
    outline: none;
    box-sizing: border-box;
    transition: border-color 0.2s ease;
  }

  .form-group input:focus,
  .form-group select:focus {
    border-color: #4f8cff;
  }

  .form-group input::placeholder {
    color: rgba(255, 255, 255, 0.3);
  }

  /* API Key 输入框带切换按钮 */
  .input-with-toggle {
    position: relative;
    display: flex;
  }

  .input-with-toggle input {
    flex: 1;
    padding-right: 36px;
  }

  .toggle-visibility {
    position: absolute;
    right: 4px;
    top: 50%;
    transform: translateY(-50%);
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    transition: color 0.2s ease;
  }

  .toggle-visibility:hover {
    color: rgba(255, 255, 255, 0.8);
  }

  .form-group select option {
    background: #1b2636;
    color: #ffffff;
  }

  .settings-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
  }

  .btn {
    padding: 8px 18px;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
  }

  .btn-primary {
    background: #4f8cff;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #3d7ae8;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
  }

  .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.15);
  }

  .input-area {
    padding: 12px 16px;
    display: flex;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .input-box {
    width: 100%;
    min-height: 60px;
    max-height: 100px;
    border: none;
    outline: none;
    background: transparent;
    color: #ffffff;
    font-size: 16px;
    line-height: 1.5;
    resize: none;
    font-family: inherit;
  }

  .input-box::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .result-area {
    flex: 1;
    padding: 8px 0;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
  }

  .result-list {
    flex: 1;
  }

  .result-item {
    display: flex;
    align-items: flex-start;
    padding: 10px 16px;
    cursor: pointer;
    transition: background-color 0.15s ease;
  }

  .result-item:hover {
    background: rgba(255, 255, 255, 0.08);
  }

  .result-item.selected {
    background: rgba(79, 140, 255, 0.2);
    border-left: 3px solid #4f8cff;
  }

  .result-index {
    color: rgba(255, 255, 255, 0.4);
    font-size: 12px;
    margin-right: 12px;
    min-width: 16px;
    padding-top: 2px;
  }

  .result-text {
    color: #ffffff;
    font-size: 14px;
    line-height: 1.5;
    flex: 1;
    word-break: break-word;
  }

  .result-hint {
    color: rgba(255, 255, 255, 0.4);
    font-size: 12px;
    text-align: center;
    padding: 8px 16px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    margin-top: auto;
  }

  .loading {
    color: rgba(255, 255, 255, 0.6);
    font-size: 14px;
    text-align: center;
    padding: 20px;
  }
</style>
