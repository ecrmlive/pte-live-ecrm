<script lang="ts">
import '#/assets/font/iconfont.css';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import {
  isCenterDiyMode,
  isHomeDiyMode,
  loadDiyEditorApi,
  loadDiySingletonApi,
  saveDiyEditorApi,
  saveDiySingletonApi,
  type DiyEditorMode,
  type DiySingletonScope,
} from '#/api/core/diy';
import { normalizeCenterDiyItems } from '#/utils/diy/center-diy-normalize';

import Model from './Model.vue';
import Params from './Params.vue';
import Type from './Type.vue';
import {
  createSingletonFallbackBootstrap,
  saveSingletonFallbackDraft,
} from './singleton-fallback';

export default {
  name: 'NativeDiyEditorPage',
  components: {
    ElButton,
    Model,
    Page,
    Params,
    Type,
  },
  data() {
    return {
      defaultData: {} as Record<string, unknown>,
      diyData: {
        items: [] as Array<Record<string, unknown>>,
        page: {} as Record<string, unknown>,
      },
      form: {
        curItem: {} as Record<string, unknown>,
        selectedIndex: -1,
        umeditor: {},
      },
      loading: true,
      mode: 'custom-add' as DiyEditorMode,
      opts: {},
      pageId: 0,
      activeConfigTab: 'content',
    };
  },
  computed: {
    selectedConfigTitle() {
      return this.form.curItem?.name || this.form.curItem?.title || '页面设置';
    },
    diyType() {
      return isCenterDiyMode(this.mode) ? 'center' : '';
    },
    isDiy() {
      return this.isSingletonEditor || isHomeDiyMode(this.mode);
    },
    singletonScope(): DiySingletonScope | null {
      const path = this.$route.path;
      if (path === '/setting/diy/cart') return 'cart';
      if (path === '/setting/diy/store') return 'store';
      if (path === '/setting/diy/list' || path === '/setting/diy/index') return 'home';
      return null;
    },
    isSingletonEditor() {
      return this.singletonScope !== null;
    },
    isCartEditor() {
      return this.singletonScope === 'cart';
    },
    pageTitle() {
      if (this.singletonScope === 'home') {
        return '首页装修';
      }
      if (this.singletonScope === 'cart') {
        return '购物车装修';
      }
      if (this.singletonScope === 'store') {
        return '店铺装修';
      }
      const titles: Record<DiyEditorMode, string> = {
        'center-add': '新增个人中心页',
        'center-edit': '个人中心',
        'custom-add': '返回 页面列表',
        'custom-edit': '编辑页面',
        'home-add': '新增页面',
        'home-edit': '编辑页面',
      };
      return titles[this.mode] ?? '页面装修';
    },
    isHeaderBackLink() {
      return !this.isSingletonEditor;
    },
    listPath() {
      const p = this.$route.path || '';
      if (p.includes('/devise/')) return '/devise/diy/list';
      return '/setting/diy/list';
    },
    headerBackRoute() {
      if (this.isSingletonEditor) return null;
      if (this.mode.startsWith('custom')) return '/setting/micro/list';
      if (this.mode.startsWith('home')) return this.listPath;
      return null;
    },
  },
  async created() {
    this.pageId = Number(this.$route.query.id ?? this.$route.query.page_id ?? 0);
    if (this.singletonScope) {
      this.mode = 'home-edit';
      await this.loadData();
      return;
    }
    const types = String(this.$route.query.types ?? '1');
    if (this.pageId > 0) {
      this.mode = types === '0' ? 'custom-edit' : 'home-edit';
    } else {
      this.mode = types === '0' ? 'custom-add' : 'home-add';
    }
    await this.loadData();
  },
  methods: {
    deepClone<T>(obj: T): T {
      return JSON.parse(JSON.stringify(obj)) as T;
    },
    async loadEditorBootstrap() {
      const scope = this.singletonScope;
      if (!scope) {
        return loadDiyEditorApi(this.mode, this.pageId || undefined);
      }

      try {
        return await loadDiySingletonApi(scope);
      } catch {
        try {
          // 旧服务尚未提供 singleton 接口时，仍使用通用编辑器配置渲染完整工作台。
          return await loadDiyEditorApi('home-edit');
        } catch {
          // 接口断网或 404 时，不让单例装修页退化为三个空白区域。
          return createSingletonFallbackBootstrap(scope);
        }
      }
    },
    async loadData() {
      this.loading = true;
      try {
        const scope = this.singletonScope;
        const res = await this.loadEditorBootstrap();
        const fallback = scope
          ? createSingletonFallbackBootstrap(scope)
          : null;
        const data = res as {
          defaultData?: Record<string, unknown>;
          jsonData?: { items?: Array<Record<string, unknown>>; page?: Record<string, unknown> };
          opts?: Record<string, unknown>;
          pageId?: number;
        };
        const pageId = Number(data.pageId ?? 0);
        if (Number.isFinite(pageId) && pageId > 0) {
          this.pageId = pageId;
        }
        this.defaultData = {
          ...(fallback?.defaultData ?? {}),
          ...(data.defaultData ?? {}),
        };
        const jsonData = this.deepClone(
          data.jsonData ?? fallback?.jsonData ?? { items: [], page: {} },
        ) as { items?: Array<Record<string, unknown>>; page?: Record<string, unknown> };
        if (!Array.isArray(jsonData.items)) {
          jsonData.items = [];
        }
        if (!jsonData.page || typeof jsonData.page !== 'object') {
          jsonData.page = {};
        }
        if (this.isSingletonEditor && jsonData.items.length === 0) {
          jsonData.items = this.createSingletonDefaultItems();
        }
        this.diyData = jsonData;
        this.ensureDiyItemKeys(this.diyData.items);
        if (isCenterDiyMode(this.mode)) {
          normalizeCenterDiyItems(this.diyData);
        }
        this.form.curItem = this.diyData.page ?? {};
        this.opts = {
          ...(fallback?.opts ?? {}),
          ...(data.opts ?? {}),
        };
      } finally {
        this.loading = false;
      }
    },
    stampDiyItem(item: Record<string, unknown>, key: string) {
      item._diyUid = `${key}-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
      return item;
    },
    ensureDiyItemKeys(items: Array<Record<string, unknown>>) {
      items.forEach((item, index) => {
        if (!item._diyUid) {
          item._diyUid = `loaded-${index}-${String(item.type ?? 'item')}`;
        }
      });
    },
  resolveDefaultItem(key: string) {
      const source = this.defaultData[key];
      if (source) {
        return this.deepClone(source);
      }
      for (const defaultKey in this.defaultData) {
        const candidate = this.defaultData[defaultKey] as { type?: string } | undefined;
        if (candidate?.type === key) {
          return this.deepClone(candidate);
        }
      }
    return null;
  },
  singletonDefaultItemKeys() {
    switch (this.singletonScope) {
      case 'home':
        return ['search', 'banner', 'navBar', 'product', 'bottomNav'];
      case 'cart':
        return ['search', 'title', 'product', 'bottomNav'];
      case 'store':
        return ['search', 'banner', 'store', 'product', 'bottomNav'];
      default:
        return [];
    }
  },
  createSingletonDefaultItems() {
    const items: Array<Record<string, unknown>> = [];
    this.singletonDefaultItemKeys().forEach((key) => {
      const item = this.resolveDefaultItem(key);
      if (item && typeof item === 'object') {
        items.push(this.stampDiyItem(item as Record<string, unknown>, key));
      }
    });
    return items;
  },
  onAddItem(key: string) {
    if (
      this.isSingletonEditor ||
      isHomeDiyMode(this.mode) ||
      this.mode === 'custom-add' ||
      this.mode === 'custom-edit'
    ) {
        this.onAddItemWithRules(key);
        return;
      }
      const resolved = this.resolveDefaultItem(key);
      if (!resolved) {
        ElMessage.error('组件模板不存在');
        return;
      }
      const item = this.stampDiyItem(resolved as Record<string, unknown>, key);
      let insertIndex = 0;
      if (this.form.selectedIndex < 0) {
        this.diyData.items.unshift(item);
      } else {
        insertIndex = this.form.selectedIndex + 1;
        this.diyData.items.splice(insertIndex, 0, item);
      }
      this.$refs.model?.onEditer(insertIndex);
    },
    onAddItemWithRules(key: string) {
		if (key === 'bottomNav' && this.diyData.items.some((i) => i.type === 'bottomNav')) {
			ElMessage.error('每个 DIY 页面只能配置一个底部导航');
			return;
		}
      if (key === 'option' || key === 'search') {
        if (this.diyData.items.some((i) => i.type === 'topMerge')) {
          ElMessage.error('轮播搜索不能与选项卡或者搜索框同时存在');
          return;
        }
        if (this.diyData.items.some((i) => i.type === key)) {
          ElMessage.error('该组件不可重复添加');
          return;
        }
      }
      if (key === 'topMerge') {
        if (this.diyData.items.some((i) => i.type === 'option' || i.type === 'search')) {
          ElMessage.error('轮播搜索不能与选项卡或者搜索框同时存在');
          return;
        }
        if (this.diyData.items.some((i) => i.type === 'topMerge')) {
          ElMessage.error('该组件不可重复添加');
          return;
        }
      }

      const resolved = this.resolveDefaultItem(key);
      if (!resolved) {
        ElMessage.error('组件模板不存在');
        return;
      }
      const item = this.stampDiyItem(resolved as Record<string, unknown>, key);
      let insertIndex = 0;

      if (key === 'bottomNav') {
        insertIndex = this.diyData.items.length;
        this.diyData.items.push(item);
      } else if (key === 'topMerge' || key === 'search' || key === 'option') {
        if (key === 'option' && this.diyData.items[0]?.type === 'search') {
          this.diyData.items.splice(1, 0, item);
          insertIndex = 1;
        } else {
          this.diyData.items.unshift(item);
          insertIndex = 0;
        }
      } else {
        insertIndex = this.form.selectedIndex < 0 ? 0 : this.form.selectedIndex + 1;
        const specialTypes = ['option', 'search', 'topMerge'];
        let lastSpecialIndex = -1;
        this.diyData.items.forEach((it, idx) => {
          if (specialTypes.includes(String(it.type))) {
            lastSpecialIndex = Math.max(lastSpecialIndex, idx);
          }
        });
        if (insertIndex <= lastSpecialIndex) {
          insertIndex = lastSpecialIndex + 1;
        }
        this.diyData.items.splice(insertIndex, 0, item);
      }

      this.$refs.model?.onEditer(insertIndex);
    },
    async submit() {
      if (!this.diyData.items?.length) {
        ElMessage.error('请至少添加一个组件');
        return;
      }
      this.loading = true;
      try {
        const scope = this.singletonScope;
        if (scope) {
          saveSingletonFallbackDraft(scope, this.diyData);
          try {
            const result = await saveDiySingletonApi(
              scope,
              JSON.stringify(this.diyData),
            );
            this.pageId = Number(result?.page_id ?? this.pageId);
            ElMessage.success('保存成功');
          } catch {
            ElMessage.warning('服务暂不可用，已保存到本机草稿');
          }
          this.form.selectedIndex = -1;
          return;
        }
        const res = await saveDiyEditorApi(
          this.mode,
          JSON.stringify(this.diyData),
          this.pageId || undefined,
        );
        ElMessage.success('保存成功');
        if (this.mode.includes('add')) {
          const pageId = Number(res?.page_id ?? 0);
          if (pageId > 0 && this.mode === 'custom-add') {
            await this.$router.replace({
              path: this.$route.path,
              query: { id: pageId, types: this.mode.startsWith('custom') ? '0' : '1' },
            });
            this.mode = 'custom-edit';
            this.pageId = pageId;
            await this.loadData();
            this.form.selectedIndex = -1;
            return;
          }
          this.$router.back();
        } else {
          await this.loadData();
          this.form.selectedIndex = -1;
        }
      } finally {
        this.loading = false;
      }
    },
    gotoBack() {
      this.$router.back();
    },
    onHeaderClick() {
      if (this.headerBackRoute) {
        this.$router.push(this.headerBackRoute);
      }
    },
  },
};
</script>

<template>
  <Page
    auto-content-height
    content-class="diy-editor-page flex min-h-0 flex-col overflow-hidden !p-0"
  >
    <div v-loading="loading" class="diy-editor-shell native-form-page native-form-shell">
      <div class="diy-editor-header flex-shrink-0">
        <div
          class="common-form"
          :class="{ 'diy-editor-header__back': isHeaderBackLink }"
          @click="onHeaderClick"
        >
          {{ pageTitle }}
        </div>
      </div>

      <div class="diy-editor-body">
        <div class="diy-container flex min-h-0 flex-1">
          <aside class="diy-menu diy-panel">
            <div class="diy-panel__scroll">
              <Type v-if="!loading" :default-data="defaultData" @add-item="onAddItem" />
            </div>
          </aside>

          <main class="diy-phone">
            <Model
              v-if="!loading"
              ref="model"
              :default-data="defaultData"
              :diy-data="diyData"
              :diy-type="diyType"
              :form="form"
              :is-diy="isDiy"
              @select-item="activeConfigTab = 'content'"
            />
          </main>

          <aside class="diy-info diy-panel">
            <div class="diy-config-head">
              <div class="diy-config-title">{{ selectedConfigTitle }}</div>
              <div class="diy-config-tabs" role="tablist" aria-label="组件配置">
                <button
                  type="button"
                  class="diy-config-tab"
                  :class="{ 'is-active': activeConfigTab === 'content' }"
                  @click="activeConfigTab = 'content'"
                >
                  内容
                </button>
                <button
                  type="button"
                  class="diy-config-tab"
                  :class="{ 'is-active': activeConfigTab === 'style' }"
                  @click="activeConfigTab = 'style'"
                >
                  样式
                </button>
              </div>
            </div>
            <div class="diy-panel__scroll">
              <Params
                v-if="!loading"
                :default-data="defaultData"
                :diy-data="diyData"
                :diy-type="diyType"
                :form="form"
                :is-diy="isDiy"
                :opts="opts"
                :config-tab="activeConfigTab"
              />
            </div>
          </aside>
        </div>
      </div>

      <div class="form-footer">
        <ElButton @click="gotoBack">取消</ElButton>
        <ElButton :loading="loading" type="primary" @click="submit">保存</ElButton>
      </div>
    </div>
  </Page>
</template>

<!-- Unscoped: palette (Type.vue) and preview utilities live in child components -->
<style lang="scss">
@use '../../../styles/diy.scss';
@use '../../../styles/diy-legacy-utils.scss';
</style>
