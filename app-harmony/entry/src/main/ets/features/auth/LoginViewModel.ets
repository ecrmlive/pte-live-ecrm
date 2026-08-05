import { AuthRepository } from '../../data/repository/AuthRepository';
import { AuthSession } from '../../domain/model/Auth';

/** 登录用例层；UI 不直接依赖 HTTP 细节。 */
export class LoginViewModel {
  private readonly authRepository: AuthRepository = new AuthRepository();

  async submit(account: string, password: string): Promise<AuthSession> {
    const normalizedAccount = account.trim();
    if (normalizedAccount.length === 0 || password.length < 6) {
      throw new Error('请输入账号和至少 6 位密码');
    }
    return this.authRepository.login(normalizedAccount, password);
  }
}
