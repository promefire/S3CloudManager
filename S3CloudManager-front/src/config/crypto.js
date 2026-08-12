import CryptoJS from 'crypto-js';

// 加密密钥（内置固定密钥，用于浏览器端加密存储）
// 注意：这是对称加密，密钥嵌入在代码中。
// 目的是防止他人直接从 localStorage 明文查看配置信息，
// 而非用于传输安全。清空浏览器缓存后数据消失。
const SECRET_KEY = 'S3Storage-Config-2026';

/**
 * 加密数据
 * @param {Object} data - 要加密的数据对象
 * @returns {string} 加密后的字符串
 */
export function encryptConfig(data) {
  try {
    const jsonStr = JSON.stringify(data);
    const encrypted = CryptoJS.AES.encrypt(jsonStr, SECRET_KEY).toString();
    return encrypted;
  } catch (error) {
    console.error('加密失败:', error);
    return null;
  }
}

/**
 * 解密数据
 * @param {string} encryptedData - 加密后的字符串
 * @returns {Object|null} 解密后的数据对象
 */
export function decryptConfig(encryptedData) {
  try {
    const bytes = CryptoJS.AES.decrypt(encryptedData, SECRET_KEY);
    const decryptedStr = bytes.toString(CryptoJS.enc.Utf8);
    if (!decryptedStr) {
      return null;
    }
    return JSON.parse(decryptedStr);
  } catch (error) {
    console.error('解密失败:', error);
    return null;
  }
}

/**
 * 生成配置指纹（用于后端识别配置是否变化）
 * @param {Object} config - 配置对象
 * @returns {string} MD5 指纹
 */
export function getConfigFingerprint(config) {
  const str = `${config.endpoint}|${config.accessKeyId}|${config.region}|${config.useSsl}`;
  return CryptoJS.MD5(str).toString();
}
