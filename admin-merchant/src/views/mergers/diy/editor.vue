<script lang="ts">
import '#/assets/font/iconfont.css';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import {
  isCenterDiyMode,
  isHomeDiyMode,
  loadDiyEditorApi,
  resolveDiyEditorMode,
  saveDiyEditorApi,
  type DiyEditorMode,
} from '#/api/core/diy';
import { normalizeCenterDiyItems } from '#/utils/diy/center-diy-normalize';

import Model from './Model.vue';
import Params from './Params.vue';
import Type from './Type.vue';

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
    };
  },
  computed: {
    diyType() {
      return isCenterDiyMode(this.mode) ? 'center' : '';
    },
    isDiy() {
      return isHomeDiyMode(this.mode);
    },
    pageTitle() {
      const titles: Record<DiyEditorMode, string> = {
        'center-add': '添加个人中心页',
        'center-edit': '个人中心',
        'custom-add': '返回 页面列表',
        'custom-edit': '编辑页面',
        'home-add': '添加首页',
        'home-edit': '首页装修',
      };
      return titles[this.mode] ?? '页面装修';
    },
    isHeaderBackLink() {
      return true;
    },
    listPath() {
      const p = this.$route.path || '';
      if (p.includes('/devise/')) return '/devise/diy/list';
      return '/setting/diy/list';
    },
    headerBackRoute() {
      if (this.mode.startsWith('custom')) return { path: this.listPath, query: { kind: 'micro' } };
      if (this.mode.startsWith('home')) return this.listPath;
      return null;
    },
  },
  async created() {
    this.pageId = Number(this.$route.query.id ?? this.$route.query.page_id ?? 0);
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
    async loadData() {
      this.loading = true;
      try {
        const res = await loadDiyEditorApi(this.mode, this.pageId || undefined);
        const data = res as {
          defaultData?: Record<string, unknown>;
          jsonData?: { items?: Array<Record<string, unknown>>; page?: Record<string, unknown> };
          opts?: Record<string, unknown>;
        };
        this.defaultData = data.defaultData ?? {};
        const jsonData = data.jsonData ?? { items: [], page: {} };
        if (!Array.isArray(jsonData.items)) {
          jsonData.items = [];
        }
        if (!jsonData.page || typeof jsonData.page !== 'object') {
          jsonData.page = {};
        }
        this.diyData = jsonData;
        this.ensureDiyItemKeys(this.diyData.items);
        if (isCenterDiyMode(this.mode)) {
          normalizeCenterDiyItems(this.diyData);
        }
        this.form.curItem = this.diyData.page ?? {};
        this.opts = data.opts ?? {};
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
    onAddItem(key: string) {
      if (isHomeDiyMode(this.mode) || this.mode === 'custom-add' || this.mode === 'custom-edit') {
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

      if (key === 'topMerge' || key === 'search' || key === 'option') {
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
    <div
      v-loading="loading"
      class="diy-editor-shell native-form-page native-form-shell"
    >
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
            <div class="diy-panel__head">组件库</div>
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
            />
          </main>

          <aside class="diy-info diy-panel">
            <div class="diy-panel__head">页面设置</div>
            <div class="diy-panel__scroll">
              <Params
                v-if="!loading"
                :default-data="defaultData"
                :diy-data="diyData"
                :diy-type="diyType"
                :form="form"
                :is-diy="isDiy"
                :opts="opts"
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
