import { AppTheme } from '../../common/theme/AppTheme';

@Component
export struct PlaceholderPage {
  @Prop title: string;
  @Prop description: string;

  build() {
    Column({ space: 12 }) {
      Text(this.title)
        .fontSize(22)
        .fontWeight(FontWeight.Bold)
        .fontColor(AppTheme.color.textPrimary)
      Text(this.description)
        .fontSize(14)
        .fontColor(AppTheme.color.textSecondary)
        .lineHeight(21)
        .textAlign(TextAlign.Center)
      Text('页面骨架已建立，待接入对应 Repository 与 ViewModel。')
        .fontSize(13)
        .fontColor(AppTheme.color.primary)
    }
    .width('100%')
    .layoutWeight(1)
    .justifyContent(FlexAlign.Center)
    .padding(AppTheme.size.pagePadding)
  }
}
