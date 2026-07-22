<script setup lang="ts">
import { onLaunch } from "@dcloudio/uni-app";
import { useUserStore } from "@/stores/user";

onLaunch(() => {
  const user = useUserStore();
  if (user.isLogin) {
    user.refreshMe();
  }
});
</script>

<style>
page {
  --qx-brand: #e23030;
  --qx-brand-dark: #c41f1f;
  --qx-text: #1f1f1f;
  --qx-text-secondary: #8a8a8a;
  --qx-bg: #f5f5f5;
  --qx-card: #ffffff;
  --qx-line: #eeeeee;
  --qx-price: #e23030;
  --qx-nav-h: 44px;
  background: var(--qx-bg);
  color: var(--qx-text);
  font-family: "PingFang SC", "Noto Sans SC", "Helvetica Neue", sans-serif;
  font-size: 28rpx;
}

.qx-btn {
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 12rpx;
  font-size: 30rpx;
  font-weight: 600;
  text-align: center;
}

.qx-btn-primary {
  background: var(--qx-brand);
  color: #fff;
}

.qx-btn-ghost {
  background: #fff;
  color: var(--qx-brand);
  border: 1px solid var(--qx-brand);
}

.qx-safe-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
