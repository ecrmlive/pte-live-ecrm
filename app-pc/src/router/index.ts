import { createRouter, createWebHistory } from "vue-router";
import PcLayout from "@/layouts/PcLayout.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: PcLayout,
      children: [
        {
          path: "",
          name: "home",
          component: () => import("@/pages/home/HomePage.vue"),
        },
        {
          path: "category",
          name: "category",
          component: () => import("@/pages/category/CategoryPage.vue"),
        },
        {
          path: "goods",
          name: "goods-list",
          component: () => import("@/pages/goods/GoodsListPage.vue"),
        },
        {
          path: "goods/:id",
          name: "goods-detail",
          component: () => import("@/pages/goods/GoodsDetailPage.vue"),
        },
        {
          path: "store/:id",
          name: "store",
          component: () => import("@/pages/store/StorePage.vue"),
        },
        {
          path: "cart",
          name: "cart",
          component: () => import("@/pages/cart/CartPage.vue"),
        },
        {
          path: "checkout",
          name: "checkout",
          component: () => import("@/pages/checkout/CheckoutPage.vue"),
        },
        {
          path: "orders",
          name: "orders",
          component: () => import("@/pages/order/OrdersPage.vue"),
        },
        {
          path: "pay/:id",
          name: "pay-result",
          component: () => import("@/pages/order/PayResultPage.vue"),
        },
        {
          path: "user",
          name: "user",
          component: () => import("@/pages/user/UserPage.vue"),
        },
        {
          path: "coupons",
          name: "coupons",
          component: () => import("@/pages/coupon/CouponsPage.vue"),
        },
        {
          path: "seckill",
          name: "seckill",
          component: () => import("@/pages/seckill/SeckillPage.vue"),
        },
        {
          path: "combination",
          name: "combination",
          component: () => import("@/pages/combination/CombinationPage.vue"),
        },
        {
          path: "reservation",
          name: "reservation",
          component: () => import("@/pages/reservation/ReservationPage.vue"),
        },
        {
          path: "presell",
          name: "presell",
          component: () => import("@/pages/presell/PresellPage.vue"),
        },
        {
          path: "live",
          name: "live",
          component: () => import("@/pages/live/LivePage.vue"),
        },
        {
          path: "live/:id",
          name: "live-detail",
          component: () => import("@/pages/live/LiveDetailPage.vue"),
        },
        {
          path: "community",
          name: "community",
          component: () => import("@/pages/community/CommunityPage.vue"),
        },
        {
          path: "community/create",
          name: "community-create",
          component: () => import("@/pages/community/CommunityCreatePage.vue"),
        },
        {
          path: "community/:id",
          name: "community-detail",
          component: () => import("@/pages/community/CommunityDetailPage.vue"),
        },
        {
          path: "assist",
          name: "assist",
          component: () => import("@/pages/assist/AssistPage.vue"),
        },
        {
          path: "assist/:id",
          name: "assist-detail",
          component: () => import("@/pages/assist/AssistDetailPage.vue"),
        },
        {
          path: "points",
          name: "points",
          component: () => import("@/pages/points/PointsPage.vue"),
        },
        {
          path: "points/checkout",
          name: "points-checkout",
          component: () => import("@/pages/points/PointsCheckoutPage.vue"),
        },
        {
          path: "notices",
          name: "notices",
          component: () => import("@/pages/notice/NoticesPage.vue"),
        },
        {
          path: "notices/:id",
          name: "notice-detail",
          component: () => import("@/pages/notice/NoticeDetailPage.vue"),
        },
        {
          path: "agreements/:key",
          name: "agreement",
          component: () => import("@/pages/agreement/AgreementPage.vue"),
        },
        {
          path: "login",
          name: "login",
          component: () => import("@/pages/auth/LoginPage.vue"),
        },
      ],
    },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
