import { AppTheme } from '../../common/theme/AppTheme';
import { LoginPage } from '../auth/LoginPage';
import { PlaceholderPage } from '../common/PlaceholderPage';
import { HomePage } from '../home/HomePage';

@Component
export struct RootPage {
  @State private selectedTab: number = 0;
  @State private loginVisible: boolean = false;

  build() {
    Stack({ alignContent: Alignment.Bottom }) {
      Column() {
        this.navigationBar()
        if (this.selectedTab === 0) {
          HomePage({ selectedTab: $selectedTab, loginVisible: $loginVisible })
        } else if (this.selectedTab === 1) {
          PlaceholderPage({ title: '商品', description: '商品分类、搜索、店铺消费视图与商品详情将在此接入。' })
        } else if (this.selectedTab === 2) {
          PlaceholderPage({ title: '购物车', description: '购物车按商户拆分结算，提交订单前由服务端校验库存和营销规则。' })
        } else if (this.selectedTab === 3) {
          PlaceholderPage({ title: '订单', description: '订单列表、详情、物流、核销和售后入口将在此接入。' })
        } else {
          PlaceholderPage({ title: '我的', description: '会员资料、地址、优惠券、积分、分销和客服入口将在此接入。' })
        }
        this.tabBar()
      }
      .width('100%')
      .height('100%')
      .backgroundColor(AppTheme.color.canvas)

      if (this.loginVisible) {
        Column() {
          LoginPage({ visible: $loginVisible })
        }
        .width('100%')
        .backgroundColor(AppTheme.color.surface)
        .borderRadius({ topLeft: 20, topRight: 20 })
        .shadow({ radius: 20, color: '#33000000', offsetY: -4 })
      }
    }
    .width('100%')
    .height('100%')
  }

  @Builder
  private navigationBar() {
    Row() {
      Text('七禧商城')
        .fontSize(18)
        .fontWeight(FontWeight.Bold)
        .fontColor(AppTheme.color.textPrimary)
      Blank()
      Button('登录')
        .height(AppTheme.size.touchTarget)
        .fontColor(AppTheme.color.primary)
        .backgroundColor(Color.Transparent)
        .onClick(() => this.loginVisible = true)
    }
    .width('100%')
    .height(AppTheme.size.navigationHeight)
    .padding({ left: AppTheme.size.pagePadding, right: 8 })
    .backgroundColor(AppTheme.color.surface)
  }

  @Builder
  private tabBar() {
    Row() {
      this.tabItem('首页', 0)
      this.tabItem('商品', 1)
      this.tabItem('购物车', 2)
      this.tabItem('订单', 3)
      this.tabItem('我的', 4)
    }
    .width('100%')
    .height(56)
    .backgroundColor(AppTheme.color.surface)
    .border({ width: { top: 1 }, color: AppTheme.color.divider })
  }

  @Builder
  private tabItem(title: string, index: number) {
    Text(title)
      .layoutWeight(1)
      .height(56)
      .textAlign(TextAlign.Center)
      .fontSize(12)
      .fontColor(this.selectedTab === index ? AppTheme.color.primary : AppTheme.color.textSecondary)
      .onClick(() => this.selectedTab = index)
  }
}
