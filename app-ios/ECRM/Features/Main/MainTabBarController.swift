import UIKit

final class MainTabBarController: UITabBarController {
    override func viewDidLoad() {
        super.viewDidLoad()
        viewControllers = [
            makeTab(title: "首页", symbol: "house", detail: "推荐商品、活动与店铺入口将在此接入。"),
            makeTab(title: "分类", symbol: "square.grid.2x2", detail: "商品分类和搜索将在此接入。"),
            makeTab(title: "购物车", symbol: "cart", detail: "多商户购物车将在此接入。"),
            makeTab(title: "订单", symbol: "list.bullet.rectangle", detail: "订单、支付与售后将在此接入。"),
            makeTab(title: "我的", symbol: "person", detail: "登录、会员与账户资产将在此接入。"),
        ]
        tabBar.tintColor = .ecrmPrimary
        tabBar.backgroundColor = .systemBackground
    }

    private func makeTab(title: String, symbol: String, detail: String) -> UIViewController {
        let screen = PlaceholderViewController(titleText: title, detail: detail)
        let navigationController = ECRMNavigationController(rootViewController: screen)
        navigationController.tabBarItem = UITabBarItem(title: title, image: UIImage(systemName: symbol), selectedImage: nil)
        return navigationController
    }
}
