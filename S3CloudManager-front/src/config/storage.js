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

// 存储桶自定义域映射
const BUCKET_DOMAINS_KEY = 's3_bucket_domains';

/**
 * 获取所有存储桶的自定义域映射
 * @returns {Object} { bucketName: domain, ... }
 */
export function getBucketDomains() {
  try {
    const data = localStorage.getItem(BUCKET_DOMAINS_KEY);
    const domains = data ? JSON.parse(data) : {};
    // 清理可能存在的旧格式（带协议前缀）
    let cleaned = false;
    for (const key in domains) {
      if (domains[key] && typeof domains[key] === 'string' && domains[key].startsWith('http')) {
        domains[key] = domains[key].replace(/^https?:\/\//, '').replace(/\/+$/, '');
        cleaned = true;
      }
    }
    if (cleaned) {
      localStorage.setItem(BUCKET_DOMAINS_KEY, JSON.stringify(domains));
    }
    return domains;
  } catch {
    return {};
  }
}

/**
 * 获取指定存储桶的自定义域（纯域名，不含协议）
 * @param {string} bucketName - 存储桶名称
 * @returns {string|null}
 */
export function getBucketDomain(bucketName) {
  const domains = getBucketDomains();
  const domain = domains[bucketName] || null;
  // 清理可能存在的旧格式（带协议前缀）
  if (domain && domain.startsWith('http')) {
    const cleaned = domain.replace(/^https?:\/\//, '').replace(/\/+$/, '');
    domains[bucketName] = cleaned;
    localStorage.setItem(BUCKET_DOMAINS_KEY, JSON.stringify(domains));
    return cleaned;
  }
  return domain;
}

/**
 * 设置存储桶的自定义域（自动清理协议前缀，统一由 getImageDomain 添加 https://）
 * @param {string} bucketName - 存储桶名称
 * @param {string} domain - 自定义域名
 */
export function setBucketDomain(bucketName, domain) {
  const domains = getBucketDomains();
  if (domain) {
    // 清理协议前缀和多余空格，只保留纯域名
    domains[bucketName] = domain.replace(/^https?:\/\//, '').replace(/\/+$/, '').trim();
  } else {
    delete domains[bucketName];
  }
  localStorage.setItem(BUCKET_DOMAINS_KEY, JSON.stringify(domains));
}
