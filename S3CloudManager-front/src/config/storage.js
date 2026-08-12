import { encryptConfig, decryptConfig, getConfigFingerprint } from './crypto.js';

const STORAGE_KEY = 's3_storage_config';

/**
 * 保存配置到 localStorage（加密存储）
 * @param {Object} config - S3 配置对象
 */
export function saveConfig(config) {
  const encrypted = encryptConfig(config);
  if (encrypted) {
    localStorage.setItem(STORAGE_KEY, encrypted);
    return true;
  }
  return false;
}

/**
 * 从 localStorage 加载配置（解密）
 * @returns {Object|null} 解密后的配置对象，未找到或解密失败返回 null
 */
export function loadConfig() {
  const encrypted = localStorage.getItem(STORAGE_KEY);
  if (!encrypted) {
    return null;
  }
  return decryptConfig(encrypted);
}

/**
 * 清除存储的配置
 */
export function clearConfig() {
  localStorage.removeItem(STORAGE_KEY);
}

/**
 * 检查是否存在已保存的配置
 * @returns {boolean}
 */
export function hasConfig() {
  return localStorage.getItem(STORAGE_KEY) !== null;
}

/**
 * 获取用于发送给后端的配置数据（包含加密数据和指纹）
 * @returns {Object|null}
 */
export function getConfigForRequest() {
  const config = loadConfig();
  if (!config) {
    return null;
  }
  return {
    encrypted: encryptConfig(config),
    fingerprint: getConfigFingerprint(config),
    // 解密后的明文配置（后端需要明文来连接 S3）
    ...config
  };
}

/**
 * 默认配置模板
 */
export const DEFAULT_CONFIG = {
  endpoint: '',
  accessKeyId: '',
  secretAccessKey: '',
  region: 'auto',
  useSsl: true,
  signatureType: 'V4',
  port: 9000
};
