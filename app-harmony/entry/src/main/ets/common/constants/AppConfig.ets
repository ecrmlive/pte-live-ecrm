/**
 * 构建时由受控环境配置替换，禁止在源码中写入真实服务凭据。
 */
export const AppConfig = {
  apiBaseUrl: 'https://api.example.invalid/api/app/v1',
  loginChannel: 'harmony',
  authorizationHeader: 'Authori-zation',
  storeAppIdHeader: 'X-AppId',
  requestTimeoutMs: 10000
} as const;
