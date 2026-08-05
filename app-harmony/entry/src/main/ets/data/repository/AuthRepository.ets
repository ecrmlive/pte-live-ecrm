import { HttpClient } from '../../common/network/HttpClient';
import { AppConfig } from '../../common/constants/AppConfig';
import { ApiEnvelope, AuthSession } from '../../domain/model/Auth';

export class AuthRepository {
  private readonly client: HttpClient = new HttpClient();

  async login(account: string, password: string): Promise<AuthSession> {
    const response = await this.client.post<ApiEnvelope<AuthSession>>('/auth/login', {
      body: {
        account,
        password,
        channel: AppConfig.loginChannel
      }
    });
    return response.data;
  }

  async currentUser(accessToken: string): Promise<AuthSession> {
    const response = await this.client.get<ApiEnvelope<AuthSession>>('/auth/me', {
      token: accessToken
    });
    return response.data;
  }
}
