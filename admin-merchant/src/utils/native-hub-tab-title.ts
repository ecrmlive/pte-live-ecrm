import type { ComputedRef, Ref } from 'vue';

import { useTabs } from '@vben/hooks';
import { onMounted, watch } from 'vue';

/** 插件 hub 页：同步 Vben 顶栏标签标题（对齐 legacy Head menu_title + 子 Tab 文案） */
export function useNativeHubTabTitle(
  activeTab: Ref<string>,
  tabLabels: Readonly<Record<string, string>>,
  options?: { hubTitle?: string },
) {
  const { setTabTitle } = useTabs();

  async function syncTopTabTitle() {
    const label =
      tabLabels[activeTab.value]?.trim() ||
      options?.hubTitle?.trim() ||
      '';
    if (label) {
      await setTabTitle(label);
    }
  }

  onMounted(() => {
    void syncTopTabTitle();
  });

  watch(activeTab, () => {
    void syncTopTabTitle();
  });

  return { syncTopTabTitle };
}

/** 单页插件（无 hub 子 Tab） */
export function useNativePageTabTitle(title: ComputedRef<string> | string) {
  const { setTabTitle } = useTabs();

  onMounted(() => {
    void setTabTitle(title);
  });
}
