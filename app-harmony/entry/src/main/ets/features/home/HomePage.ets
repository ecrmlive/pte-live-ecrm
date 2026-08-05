import { AppTheme } from '../../common/theme/AppTheme';

@Component
export struct HomePage {
  @Link selectedTab: number;
  @Link loginVisible: boolean;

  build() {
    Scroll() {
      Column({ space: 16 }) {
        Column({ space: 8 }) {
          Text('七禧商城').fontSize(24).fontWeight(FontWeight.Bold).fontColor(Color.White)
          Text('多商户好物，原生鸿蒙体验').fontSize(14).fontColor('#FFE9E9')
        }
        .width('100%')
        .padding(20)
        .borderRadius(12)
        .backgroundColor(AppTheme.color.primary)

        Text('快捷入口').fontSize(18).fontWeight(FontWeight.Medium).fontColor(AppTheme.color.textPrimary)
        Row({ space: 10 }) {
          this.actionCard('商品分类', '浏览全量商品', () => this.selectedTab = 1)
          this.actionCard('购物车', '查看待结算商品', () => this.selectedTab = 2)
        }
        .width('100%')

        Row({ space: 10 }) {
          this.actionCard('我的订单', '查看履约与售后', () => this.selectedTab = 3)
          this.actionCard('登录', '同步账号与权益', () => this.loginVisible = true)
        }
        .width('100%')

        Text('首期范围').fontSize(18).fontWeight(FontWeight.Medium).fontColor(AppTheme.color.textPrimary)
        Text('按 C 端功能清单逐步接入登录、商品、购物车、订单、售后、会员、营销和客服。支付、库存、优惠与资产由服务端完成幂等与事务控制。')
          .fontSize(14)
          .lineHeight(21)
          .fontColor(AppTheme.color.textSecondary)
          .padding(16)
          .backgroundColor(AppTheme.color.surface)
          .borderRadius(12)
      }
      .width('100%')
      .padding(AppTheme.size.pagePadding)
    }
    .width('100%')
    .layoutWeight(1)
  }

  @Builder
  private actionCard(title: string, description: string, action: () => void) {
    Column({ space: 8 }) {
      Text(title).fontSize(16).fontWeight(FontWeight.Medium).fontColor(AppTheme.color.textPrimary)
      Text(description).fontSize(12).fontColor(AppTheme.color.textSecondary)
    }
    .alignItems(HorizontalAlign.Start)
    .layoutWeight(1)
    .height(92)
    .padding(14)
    .backgroundColor(AppTheme.color.surface)
    .borderRadius(12)
    .onClick(action)
  }
}
