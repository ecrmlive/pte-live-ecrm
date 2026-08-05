import { AppTheme } from '../../common/theme/AppTheme';
import { LoginViewModel } from './LoginViewModel';

@Component
export struct LoginPage {
  @Link visible: boolean;
  @State private account: string = '';
  @State private password: string = '';
  @State private message: string = '';
  @State private submitting: boolean = false;
  private readonly viewModel: LoginViewModel = new LoginViewModel();

  private async submit(): Promise<void> {
    this.submitting = true;
    this.message = '';
    try {
      const session = await this.viewModel.submit(this.account, this.password);
      AppStorage.setOrCreate('ecrm.session.token', session.accessToken);
      AppStorage.setOrCreate('ecrm.session.nickname', session.user.nickname);
      this.visible = false;
    } catch (error) {
      this.message = error instanceof Error ? error.message : '登录暂时不可用';
    } finally {
      this.submitting = false;
    }
  }

  build() {
    Column({ space: 16 }) {
      Row() {
        Text('账号登录')
          .fontSize(20)
          .fontWeight(FontWeight.Bold)
          .fontColor(AppTheme.color.textPrimary)
        Blank()
        Button('关闭')
          .type(ButtonType.Normal)
          .fontColor(AppTheme.color.textSecondary)
          .backgroundColor(Color.Transparent)
          .height(AppTheme.size.touchTarget)
          .onClick(() => this.visible = false)
      }
      .width('100%')
      .height(AppTheme.size.navigationHeight)

      TextInput({ placeholder: '账号', text: this.account })
        .onChange((value: string) => this.account = value)
        .height(48)
        .backgroundColor(AppTheme.color.canvas)
        .borderRadius(8)

      TextInput({ placeholder: '密码（至少 6 位）', text: this.password })
        .type(InputType.Password)
        .onChange((value: string) => this.password = value)
        .height(48)
        .backgroundColor(AppTheme.color.canvas)
        .borderRadius(8)

      if (this.message.length > 0) {
        Text(this.message).fontSize(13).fontColor(AppTheme.color.primary)
      }

      Button(this.submitting ? '登录中…' : '登录')
        .width('100%')
        .height(AppTheme.size.touchTarget)
        .fontColor(Color.White)
        .backgroundColor(AppTheme.color.primary)
        .borderRadius(22)
        .enabled(!this.submitting)
        .onClick(() => this.submit())

      Text('登录来源将以 harmony 记录，并与其他 C 端共享同一用户主体。')
        .fontSize(12)
        .fontColor(AppTheme.color.textSecondary)
        .lineHeight(18)
    }
    .width('100%')
    .padding(AppTheme.size.pagePadding)
    .backgroundColor(AppTheme.color.surface)
  }
}
