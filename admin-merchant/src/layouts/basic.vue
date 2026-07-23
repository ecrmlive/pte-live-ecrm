<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { AuthenticationLoginExpiredModal } from '@vben/common-ui';
import { LockKeyhole } from '@vben/icons';
import { useWatermark } from '@vben/hooks';
import { BasicLayout, LockScreen, UserDropdown } from '@vben/layouts';
import { preferences, usePreferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';
import { storeToRefs } from 'pinia';
import { ElNotification } from 'element-plus';

import { useAuthStore } from '#/store';
import { getLegacyUserInfo } from '#/utils/qixi-live-token';
import {
  applyTenantAppBranding,
  resolveShopDisplayName,
} from '#/utils/shop-display-name';
import LoginForm from '#/views/_core/authentication/login.vue';
import ShopUpdatePasswordDialog from '#/components/shop/shop-update-password-dialog.vue';

const userStore = useUserStore();
const authStore = useAuthStore();
const { isDefaultPassword } = storeToRefs(authStore);
const accessStore = useAccessStore();
const { destroyWatermark, updateWatermark } = useWatermark();
const { isDark } = usePreferences();
const showPasswordDialog = ref(false);

const legacyUser = computed(() => getLegacyUserInfo());

const shopName = computed(() => {
  const info = userStore.userInfo as Record<string, unknown> | null;
  return resolveShopDisplayName(
    (info?.shopName as string) || legacyUser.value?.shopName,
  );
});

const userLabel = computed(() => {
  return (
    legacyUser.value?.userName ||
    userStore.userInfo?.username ||
    userStore.userInfo?.realName ||
    shopName.value
  );
});

const userDescription = computed(() => {
  const version = legacyUser.value?.version;
  return version ? `当前版本：${version}` : shopName.value;
});

const userMenus = computed(() => [
  {
    handler: () => {
      showPasswordDialog.value = true;
    },
    icon: LockKeyhole,
    text: '修改密码',
  },
]);

const avatar = computed(() => {
  return userStore.userInfo?.avatar ?? preferences.app.defaultAvatar;
});

async function handleLogout() {
  await authStore.logout(false);
}

function promptDefaultPasswordChange() {
  ElNotification({
    title: '安全提示',
    message: '您当前使用的是默认密码，为了账户安全，请及时修改密码。',
    type: 'warning',
    duration: 10000,
  });
  showPasswordDialog.value = true;
}

watch(
  isDefaultPassword,
  (value) => {
    if (value) {
      promptDefaultPasswordChange();
    }
  },
  { immediate: true },
);

watch(
  () => userStore.userInfo,
  (info) => {
    const legacy = getLegacyUserInfo();
    const name =
      (info as Record<string, unknown> | undefined)?.shopName ||
      legacy?.shopName;
    const logoUrl =
      (info as Record<string, unknown> | undefined)?.logoUrl ||
      legacy?.logoUrl;

    if (name && typeof name === 'string') {
      applyTenantAppBranding(name, logoUrl as string | undefined);
    } else if (logoUrl && typeof logoUrl === 'string') {
      applyTenantAppBranding(undefined, logoUrl);
    }
  },
  { immediate: true, deep: true },
);

watch(
  () => ({
    enable: preferences.app.watermark,
    content: preferences.app.watermarkContent,
    isDark: isDark.value,
  }),
  async ({ enable, content, isDark: isDarkValue }) => {
    if (enable) {
      const watermarkColor = isDarkValue
        ? 'rgba(255, 255, 255, 0.12)'
        : 'rgba(0, 0, 0, 0.12)';

      await updateWatermark({
        advancedStyle: {
          colorStops: [
            { color: watermarkColor, offset: 0 },
            { color: watermarkColor, offset: 1 },
          ],
          type: 'linear',
        },
        content: content || `${shopName.value} · ${userStore.userInfo?.username || ''}`,
      });
    } else {
      destroyWatermark();
    }
  },
  { immediate: true },
);
</script>

<template>
  <BasicLayout @clear-preferences-and-logout="handleLogout">
    <template #logo-text>
      <span class="truncate">{{ shopName }}</span>
    </template>
    <template #user-dropdown>
      <UserDropdown
        :avatar
        :menus="userMenus"
        :text="userLabel"
        :description="userDescription"
        @logout="handleLogout"
        @clear-preferences-and-logout="handleLogout"
      />
    </template>
    <template #extra>
      <AuthenticationLoginExpiredModal
        v-model:open="accessStore.loginExpired"
        :avatar
      >
        <LoginForm />
      </AuthenticationLoginExpiredModal>
    </template>
    <template #lock-screen>
      <LockScreen :avatar @to-login="handleLogout" />
    </template>
  </BasicLayout>

  <ShopUpdatePasswordDialog
    v-model:open="showPasswordDialog"
    @success="authStore.clearDefaultPasswordFlag()"
  />
</template>
