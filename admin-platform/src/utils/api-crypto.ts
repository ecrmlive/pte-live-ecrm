/** mergers 不启用 shop 响应体加解密 */
export async function attachAPIEncryption(config: any) {
  return config;
}

export async function decryptAPIResponse(data: any) {
  return data;
}
