import { http } from '@kit.NetworkKit';
import { AppConfig } from '../constants/AppConfig';

export interface RequestOptions {
  token?: string;
  storeAppId?: string;
  body?: Object;
}

/**
 * C 端请求统一入口。店铺范围只通过 X-AppId 传递，不能由客户端传 merchantId 或 storeId。
 */
export class HttpClient {
  async get<T>(path: string, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(path, http.RequestMethod.GET, options);
  }

  async post<T>(path: string, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(path, http.RequestMethod.POST, options);
  }

  private async request<T>(
    path: string,
    method: http.RequestMethod,
    options: RequestOptions
  ): Promise<T> {
    const client = http.createHttp();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json'
    };
    if (options.token) {
      headers[AppConfig.authorizationHeader] = `Bearer ${options.token}`;
    }
    if (options.storeAppId) {
      headers[AppConfig.storeAppIdHeader] = options.storeAppId;
    }

    try {
      const response = await client.request(`${AppConfig.apiBaseUrl}${path}`, {
        method,
        header: headers,
        extraData: options.body ? JSON.stringify(options.body) : undefined,
        expectDataType: http.HttpDataType.STRING,
        connectTimeout: AppConfig.requestTimeoutMs,
        readTimeout: AppConfig.requestTimeoutMs
      });
      if (response.responseCode < 200 || response.responseCode >= 300) {
        throw new Error(`请求失败：HTTP ${response.responseCode}`);
      }
      return JSON.parse(response.result as string) as T;
    } finally {
      client.destroy();
    }
  }
}
