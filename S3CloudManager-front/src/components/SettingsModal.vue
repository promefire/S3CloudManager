<template>
  <div>
    <!-- 触发按钮 -->
    <a href="#" class="header-icon-btn" @click.prevent="openModal" title="配置">
      <i class="material-icons">settings</i>
      <span v-if="hasConfig" class="config-indicator"></span>
    </a>

    <!-- 配置模态框 -->
    <div v-if="isOpen" class="settings-overlay" @click.self="closeModal">
      <div class="settings-modal">
        <div class="settings-header">
          <h5>
            <i class="material-icons left" style="font-size: 22px; margin-right: 8px;">settings</i>
            S3 连接配置
          </h5>
          <button @click="closeModal" class="btn-floating btn-small btn-secondary waves-effect waves-light">
            <i class="material-icons" style="color: var(--color-text-secondary);">close</i>
          </button>
        </div>

        <div class="settings-body">
          <!-- 未配置提示 -->
          <div v-if="!hasConfig" class="settings-notice">
            <i class="material-icons">info</i>
            <span>配置您自己的 S3 连接信息，数据将加密存储在浏览器本地。</span>
          </div>

          <div v-else class="settings-notice success">
            <i class="material-icons">check_circle</i>
            <span>已配置自定义连接。数据加密存储在浏览器本地，清空缓存后失效。</span>
          </div>

          <form @submit.prevent="saveSettings">
            <!-- Endpoint -->
            <div class="settings-field">
              <label class="settings-label">S3 Endpoint <span class="required">*</span></label>
              <input
                type="text"
                v-model="form.endpoint"
                placeholder="例如: account-id.r2.cloudflarestorage.com"
                required
              />
              <span class="settings-hint">S3 兼容服务端点（不含 https://）</span>
            </div>

            <!-- Access Key ID -->
            <div class="settings-field">
              <label class="settings-label">Access Key ID <span class="required">*</span></label>
              <input
                type="text"
                v-model="form.accessKeyId"
                placeholder="输入 Access Key ID"
                required
              />
            </div>

            <!-- Secret Access Key -->
            <div class="settings-field">
              <label class="settings-label">Secret Access Key <span class="required">*</span></label>
              <div class="password-field">
                <input
                  :type="showPassword ? 'text' : 'password'"
                  v-model="form.secretAccessKey"
                  placeholder="输入 Secret Access Key"
                  required
                />
                <button type="button" class="password-toggle" @click="showPassword = !showPassword">
                  <i class="material-icons" style="font-size: 18px;">
                    {{ showPassword ? 'visibility_off' : 'visibility' }}
                  </i>
                </button>
              </div>
            </div>

            <!-- Region & SSL -->
            <div class="settings-row">
              <div class="settings-field flex-1">
                <label class="settings-label">Region</label>
                <input
                  type="text"
                  v-model="form.region"
                  placeholder="auto"
                />
              </div>
              <div class="settings-field" style="width: 120px;">
                <label class="settings-label">SSL</label>
                <div class="toggle-field">
                  <label class="toggle-switch">
                    <input type="checkbox" v-model="form.useSsl" />
                    <span class="toggle-slider"></span>
                  </label>
                  <span class="toggle-label">{{ form.useSsl ? '启用' : '禁用' }}</span>
                </div>
              </div>
            </div>

            <!-- Signature Type -->
            <div class="settings-field">
              <label class="settings-label">签名类型</label>
              <div class="radio-group">
                <label class="radio-label">
                  <input type="radio" v-model="form.signatureType" value="V4" />
                  <span>V4</span>
                </label>
                <label class="radio-label">
                  <input type="radio" v-model="form.signatureType" value="V2" />
                  <span>V2</span>
                </label>
              </div>
              <span class="settings-hint">Cloudflare R2 / AWS S3 通常使用 V4</span>
            </div>

            <!-- 操作按钮 -->
            <div class="settings-actions">
              <button type="submit" class="btn btn-primary waves-effect waves-light">
                <i class="material-icons left">save</i>保存配置
              </button>
              <button v-if="hasConfig" type="button" class="btn btn-danger waves-effect waves-light" @click="clearSettings">
                <i class="material-icons left">delete</i>清除配置
              </button>
              <button type="button" class="btn btn-secondary waves-effect waves-light" @click="testConnection" :disabled="isTesting">
                <i class="material-icons left">{{ isTesting ? 'hourglass_empty' : 'wifi_tethering' }}</i>
                {{ isTesting ? '测试中...' : '测试连接' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { saveConfig, loadConfig, clearConfig, hasConfig, DEFAULT_CONFIG } from '../config/storage.js';

export default {
  name: 'SettingsModal',
  data() {
    return {
      isOpen: false,
      showPassword: false,
      isTesting: false,
      hasConfig: false,
      form: { ...DEFAULT_CONFIG }
    };
  },
  mounted() {
    this.checkConfig();
  },
  methods: {
    openModal() {
      this.isOpen = true;
      this.checkConfig();
      if (this.hasConfig) {
        const saved = loadConfig();
        if (saved) {
          this.form = { ...DEFAULT_CONFIG, ...saved };
        }
      }
    },
    closeModal() {
      this.isOpen = false;
      this.showPassword = false;
    },
    checkConfig() {
      this.hasConfig = hasConfig();
    },
    saveSettings() {
      if (!this.form.endpoint || !this.form.accessKeyId || !this.form.secretAccessKey) {
        this.$root.notify('请填写所有必填项', 'error');
        return;
      }
      const success = saveConfig(this.form);
      if (success) {
        this.$root.notify('配置已保存（加密存储）', 'success');
        this.hasConfig = true;
        this.$emit('config-saved', this.form);
      } else {
        this.$root.notify('保存失败', 'error');
      }
    },
    clearSettings() {
      if (confirm('确定要清除配置吗？清除后将使用服务器默认配置。')) {
        clearConfig();
        this.hasConfig = false;
        this.form = { ...DEFAULT_CONFIG };
        this.$root.notify('配置已清除', 'success');
        this.$emit('config-cleared');
      }
    },
    async testConnection() {
      if (!this.form.endpoint || !this.form.accessKeyId || !this.form.secretAccessKey) {
        this.$root.notify('请先填写连接信息', 'error');
        return;
      }
      this.isTesting = true;
      try {
        const response = await fetch('/api/health', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-S3-Config': JSON.stringify(this.form)
          }
        });
        if (response.ok) {
          this.$root.notify('连接测试成功！', 'success');
        } else {
          const data = await response.json();
          this.$root.notify('连接失败: ' + (data.error || '未知错误'), 'error');
        }
      } catch (error) {
        this.$root.notify('连接失败: ' + error.message, 'error');
      } finally {
        this.isTesting = false;
      }
    }
  }
};
</script>

<style scoped>
.header-icon-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: background-color 0.15s ease;
}

.header-icon-btn:hover {
  background-color: var(--color-bg);
  color: var(--color-text);
  text-decoration: none;
}

.header-icon-btn i {
  font-size: 22px;
}

.config-indicator {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-success);
  border: 2px solid white;
}

/* 遮罩层 */
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.5);
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
}

/* 模态框 */
.settings-modal {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-lg) var(--spacing-2xl);
  border-bottom: 1px solid var(--color-border);
}

.settings-header h5 {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
  display: flex;
  align-items: center;
}

.settings-body {
  padding: var(--spacing-2xl);
}

/* 提示信息 */
.settings-notice {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  background: var(--color-primary-light);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-xl);
  font-size: var(--font-size-sm);
  color: var(--color-text);
  line-height: 1.5;
}

.settings-notice i {
  font-size: 20px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.settings-notice.success {
  background: #F0FDF4;
}

.settings-notice.success i {
  color: var(--color-success);
}

/* 表单字段 */
.settings-field {
  margin-bottom: var(--spacing-lg);
}

.settings-row {
  display: flex;
  gap: var(--spacing-lg);
}

.flex-1 {
  flex: 1;
}

.settings-label {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: var(--spacing-sm);
}

.settings-label .required {
  color: var(--color-danger);
}

.settings-field input[type="text"],
.settings-field input[type="password"] {
  width: 100%;
  height: 40px;
  padding: 0 var(--spacing-md);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-md);
  font-family: var(--font-family);
  color: var(--color-text);
  background: var(--color-surface);
  transition: border-color 0.15s ease;
  outline: none;
}

.settings-field input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.settings-field input::placeholder {
  color: var(--color-text-secondary);
}

.settings-hint {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
}

/* 密码字段 */
.password-field {
  position: relative;
}

.password-field input {
  padding-right: 40px;
}

.password-toggle {
  position: absolute;
  right: var(--spacing-sm);
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-secondary);
  padding: var(--spacing-xs);
}

/* Toggle 开关 */
.toggle-field {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  height: 40px;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 26px;
  vertical-align: middle;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
  position: absolute;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  border: none;
  outline: none;
  margin: 0;
  padding: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #CBD5E1;
  border-radius: 13px;
  transition: background-color 0.2s ease;
  box-sizing: border-box;
  display: block;
}

.toggle-slider::before {
  content: "";
  position: absolute;
  width: 20px;
  height: 20px;
  left: 3px;
  top: 3px;
  background-color: #FFFFFF;
  border-radius: 50%;
  transition: transform 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--color-primary);
}

.toggle-switch input:checked + .toggle-slider::before {
  transform: translateX(22px);
}

.toggle-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-left: var(--spacing-sm);
}

/* 单选按钮组 */
.radio-group {
  display: flex;
  gap: var(--spacing-lg);
}

.radio-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  cursor: pointer;
  font-size: var(--font-size-md);
  color: var(--color-text);
}

.radio-label input[type="radio"] {
  width: 18px;
  height: 18px;
  accent-color: var(--color-primary);
}

/* 操作按钮 */
.settings-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-2xl);
  padding-top: var(--spacing-xl);
  border-top: 1px solid var(--color-border);
  flex-wrap: wrap;
}

/* 按钮样式 */
.btn-primary {
  background-color: var(--color-primary) !important;
  color: white !important;
}

.btn-primary:hover {
  background-color: var(--color-primary-hover) !important;
}

.btn-primary i {
  color: white !important;
}

.btn-secondary {
  background-color: var(--color-surface) !important;
  color: var(--color-text) !important;
  border: 1px solid var(--color-border) !important;
}

.btn-secondary:hover {
  background-color: var(--color-bg) !important;
}

.btn-danger {
  background-color: var(--color-danger) !important;
  color: white !important;
}

.btn-danger:hover {
  background-color: var(--color-danger-hover) !important;
}

.btn-danger i {
  color: white !important;
}
</style>
