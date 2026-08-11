<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getWechatMenusApi,
  saveWechatMenusApi,
  type WechatMenuButton,
} from '#/api/core/platform-wechat-menus';
import wechatMenuBottom from '#/assets/wechat-menu-bottom.png';
import wechatMenuHead from '#/assets/wechat-menu-head.png';

const RULE_OPTIONS = [
  { label: '关键字', value: 'click' },
  { label: '跳转网页', value: 'view' },
  { label: '小程序', value: 'miniprogram' },
] as const;

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const list = ref<WechatMenuButton[]>([]);

/** 当前选中：顶级 index；子菜单时 parentIndex + childIndex */
const parentIndex = ref<number | null>(null);
const childIndex = ref<number | null>(null);

const form = reactive({
  name: '',
  type: 'view' as string,
  key: '',
  url: '',
  appid: '',
  pagepath: '',
});

const hasSelection = computed(
  () => parentIndex.value !== null && parentIndex.value >= 0,
);

const selectedIsParentWithChildren = computed(() => {
  if (parentIndex.value === null || childIndex.value !== null) return false;
  const item = list.value[parentIndex.value];
  return !!item && (item.sub_button?.length || 0) > 0;
});

const selectedLabel = computed(() => {
  if (!hasSelection.value) return '';
  if (childIndex.value !== null && parentIndex.value !== null) {
    return (
      list.value[parentIndex.value]?.sub_button?.[childIndex.value]?.name ||
      '子菜单'
    );
  }
  return list.value[parentIndex.value!]?.name || '一级菜单';
});

function emptyTop(): WechatMenuButton {
  return { type: 'click', name: '一级菜单', key: '', sub_button: [] };
}

function emptyChild(): WechatMenuButton {
  return { type: 'click', name: '子菜单', key: '' };
}

function childCount(item?: WechatMenuButton) {
  return item?.sub_button?.length || 0;
}

function syncFormFromSelection() {
  if (parentIndex.value === null) {
    form.name = '';
    form.type = 'view';
    form.key = '';
    form.url = '';
    form.appid = '';
    form.pagepath = '';
    return;
  }
  const parent = list.value[parentIndex.value];
  if (!parent) return;
  const item =
    childIndex.value !== null
      ? parent.sub_button?.[childIndex.value]
      : parent;
  if (!item) return;
  form.name = item.name || '';
  form.type = item.type || 'click';
  form.key = item.key || '';
  form.url = item.url || '';
  form.appid = item.appid || '';
  form.pagepath = item.pagepath || '';
}

function applyFormToSelection() {
  if (parentIndex.value === null) return;
  const parent = list.value[parentIndex.value];
  if (!parent) return;

  const next: WechatMenuButton = {
    name: form.name.trim(),
  };

  if (childIndex.value === null && childCount(parent) > 0) {
    next.sub_button = parent.sub_button;
    list.value[parentIndex.value] = next;
    return;
  }

  next.type = form.type;
  if (form.type === 'click') {
    next.key = form.key.trim();
  } else if (form.type === 'view') {
    next.url = form.url.trim();
  } else if (form.type === 'miniprogram') {
    next.appid = form.appid.trim();
    next.pagepath = form.pagepath.trim();
    next.url = form.url.trim();
  }

  if (childIndex.value === null) {
    next.sub_button = parent.sub_button || [];
    list.value[parentIndex.value] = next;
  } else if (parent.sub_button) {
    parent.sub_button[childIndex.value] = next;
  }
}

function selectTop(index: number) {
  if (hasSelection.value) applyFormToSelection();
  parentIndex.value = index;
  childIndex.value = null;
  syncFormFromSelection();
}

function selectChild(pIndex: number, cIndex: number) {
  if (hasSelection.value) applyFormToSelection();
  parentIndex.value = pIndex;
  childIndex.value = cIndex;
  syncFormFromSelection();
}

function addTopMenu() {
  if (!canManage.value) return;
  if (list.value.length >= 3) {
    ElMessage.warning('最多添加 3 个一级菜单');
    return;
  }
  if (hasSelection.value) applyFormToSelection();
  list.value.push(emptyTop());
  parentIndex.value = list.value.length - 1;
  childIndex.value = null;
  syncFormFromSelection();
}

function addChildMenu(pIndex: number) {
  if (!canManage.value) return;
  const parent = list.value[pIndex];
  if (!parent) return;
  if (!parent.sub_button) parent.sub_button = [];
  if (parent.sub_button.length >= 5) {
    ElMessage.warning('每个一级菜单最多 5 个子菜单');
    return;
  }
  if (hasSelection.value) applyFormToSelection();
  // 保留父级 type/url 供表单回显（发布时 toPayload 会去掉）
  parent.sub_button.push(emptyChild());
  parentIndex.value = pIndex;
  childIndex.value = parent.sub_button.length - 1;
  syncFormFromSelection();
}

async function removeSelected() {
  if (!canManage.value || !hasSelection.value) {
    ElMessage.warning('请选择菜单');
    return;
  }
  try {
    await confirm({
      title: '删除确认',
      content: `确定删除「${selectedLabel.value}」吗？`,
    });
  } catch {
    return;
  }
  const p = parentIndex.value!;
  if (childIndex.value !== null) {
    list.value[p]?.sub_button?.splice(childIndex.value, 1);
  } else {
    list.value.splice(p, 1);
  }
  parentIndex.value = null;
  childIndex.value = null;
  syncFormFromSelection();
}

function validateBeforeSave(): boolean {
  if (hasSelection.value) applyFormToSelection();
  if (list.value.length === 0) {
    ElMessage.warning('请添加至少一个按钮');
    return false;
  }
  const urlRe =
    /^https?:\/\/[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/;

  for (const top of list.value) {
    if (!top.name?.trim()) {
      ElMessage.warning('请输入按钮名称');
      return false;
    }
    const children = top.sub_button || [];
    if (children.length > 0) {
      for (const child of children) {
        if (!child.name?.trim()) {
          ElMessage.warning('请输入子菜单名称');
          return false;
        }
        if (!validateLeaf(child, urlRe)) return false;
      }
      continue;
    }
    if (!validateLeaf(top, urlRe)) return false;
  }
  return true;
}

function validateLeaf(item: WechatMenuButton, urlRe: RegExp): boolean {
  const type = item.type || '';
  if (type === 'click') {
    if (!item.key?.trim()) {
      ElMessage.warning('请输入关键字');
      return false;
    }
  } else if (type === 'view') {
    if (!item.url?.trim() || !urlRe.test(item.url.trim())) {
      ElMessage.warning('请输入正确的跳转地址');
      return false;
    }
  } else if (type === 'miniprogram') {
    if (!item.appid?.trim() || !item.pagepath?.trim() || !item.url?.trim()) {
      ElMessage.warning('请填写完整小程序配置');
      return false;
    }
    if (!urlRe.test(item.url.trim())) {
      ElMessage.warning('请输入正确的跳转地址');
      return false;
    }
  } else {
    ElMessage.warning('请选择规则状态');
    return false;
  }
  return true;
}

function toPayload(): WechatMenuButton[] {
  return list.value.map((top) => {
    const children = (top.sub_button || []).filter((c) => c.name?.trim());
    if (children.length > 0) {
      return {
        name: top.name.trim(),
        sub_button: children.map((c) => leafPayload(c)),
      };
    }
    return leafPayload(top);
  });
}

function leafPayload(item: WechatMenuButton): WechatMenuButton {
  const type = item.type || 'click';
  const base: WechatMenuButton = { name: item.name.trim(), type };
  if (type === 'click') base.key = (item.key || '').trim();
  if (type === 'view') base.url = (item.url || '').trim();
  if (type === 'miniprogram') {
    base.appid = (item.appid || '').trim();
    base.pagepath = (item.pagepath || '').trim();
    base.url = (item.url || '').trim();
  }
  return base;
}

async function load() {
  loading.value = true;
  try {
    const data = await getWechatMenusApi();
    list.value = (data.wechat_menus || []).map((item) => ({
      ...item,
      sub_button: item.sub_button ? [...item.sub_button] : [],
    }));
    parentIndex.value = list.value.length ? 0 : null;
    childIndex.value = null;
    syncFormFromSelection();
  } catch {
    list.value = [];
    parentIndex.value = null;
    childIndex.value = null;
  } finally {
    loading.value = false;
  }
}

async function saveAndPublish() {
  if (!canManage.value) return;
  if (!validateBeforeSave()) return;
  saving.value = true;
  try {
    const data = await saveWechatMenusApi(toPayload());
    list.value = (data.wechat_menus || []).map((item) => ({
      ...item,
      sub_button: item.sub_button ? [...item.sub_button] : [],
    }));
    ElMessage.success(data.published ? '已保存并发布' : '已保存菜单配置');
    if (parentIndex.value !== null && parentIndex.value >= list.value.length) {
      parentIndex.value = list.value.length ? 0 : null;
      childIndex.value = null;
    }
    syncFormFromSelection();
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value = codes.includes('app.wechat_menus.manage');
  await load();
});
</script>

<template>
  <Page auto-content-height>
    <div v-loading="loading" class="wechat-menus">
      <div class="wechat-menus__body">
        <!-- 左侧：手机预览（对齐 CRMEB：头图 + 灰区 + 底栏图含键盘） -->
        <div class="wechat-phone">
          <div class="wechat-phone__frame">
            <img
              class="wechat-phone__head"
              :src="wechatMenuHead"
              alt=""
            />
            <div class="wechat-phone__body" />
            <img
              class="wechat-phone__foot"
              :src="wechatMenuBottom"
              alt=""
            />

            <div class="wechat-phone__textbot">
              <div
                v-for="(item, index) in list"
                :key="`top-${index}`"
                class="wechat-phone__li"
                :class="{
                  'is-active': parentIndex === index && childIndex === null,
                }"
              >
                <div>
                  <!-- 仅当前选中列显示：添加子菜单（图2） -->
                  <div
                    v-if="
                      canManage &&
                      parentIndex === index &&
                      childCount(item) < 5
                    "
                    class="wechat-phone__add"
                    title="新增子菜单"
                    @click.stop="addChildMenu(index)"
                  >
                    <span class="wechat-phone__plus">+</span>
                    <i class="wechat-phone__arrow" />
                  </div>

                  <!-- 子菜单列表叠在添加条上方 -->
                  <div
                    v-if="
                      parentIndex === index && childCount(item) > 0
                    "
                    class="wechat-phone__tianjia"
                  >
                    <div
                      v-for="(child, cIndex) in item.sub_button"
                      :key="`child-${index}-${cIndex}`"
                      class="wechat-phone__addadd"
                      :class="{
                        'is-active':
                          parentIndex === index && childIndex === cIndex,
                      }"
                      @click.stop="selectChild(index, cIndex)"
                    >
                      {{ child.name || '二级菜单' }}
                    </div>
                  </div>
                </div>

                <div
                  class="wechat-phone__text"
                  :title="item.name || '一级菜单'"
                  @click="selectTop(index)"
                >
                  {{ item.name || '一级菜单' }}
                </div>
              </div>

              <div
                v-if="canManage && list.length < 3"
                class="wechat-phone__li wechat-phone__li--add"
                title="新增一级菜单"
                @click="addTopMenu"
              >
                <div class="wechat-phone__text">
                  <span class="wechat-phone__plus wechat-phone__plus--muted">+</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：菜单信息 -->
        <div class="wechat-editor">
          <div class="wechat-editor__head">
            <div class="wechat-editor__tabs">
              <span class="wechat-editor__tab is-active">菜单信息</span>
            </div>
            <ElButton
              v-if="canManage && hasSelection"
              type="danger"
              size="small"
              @click="removeSelected"
            >
              删除
            </ElButton>
          </div>

          <div class="wechat-editor__body">
            <template v-if="hasSelection">
              <ElAlert
                v-if="selectedIsParentWithChildren"
                class="wechat-editor__alert"
                type="success"
                show-icon
                closable
                title="已添加子菜单，仅可设置菜单名称"
              />

              <ElForm
                class="wechat-editor__form"
                label-width="96px"
                :disabled="!canManage"
                @submit.prevent
              >
                <ElFormItem label="菜单名称" required>
                  <ElInput
                    v-model="form.name"
                    maxlength="12"
                    show-word-limit
                    placeholder="请填写菜单名称"
                    @input="applyFormToSelection"
                  />
                </ElFormItem>

                <ElFormItem label="规则状态" required>
                  <ElSelect
                    v-model="form.type"
                    class="w-full"
                    placeholder="请选择规则状态"
                    :disabled="!canManage || selectedIsParentWithChildren"
                    @change="applyFormToSelection"
                  >
                    <ElOption
                      v-for="opt in RULE_OPTIONS"
                      :key="opt.value"
                      :label="opt.label"
                      :value="opt.value"
                    />
                  </ElSelect>
                </ElFormItem>

                <ElFormItem
                  v-if="form.type === 'click'"
                  label="关键字"
                  required
                >
                  <ElInput
                    v-model="form.key"
                    maxlength="128"
                    placeholder="请输入关键字"
                    :disabled="!canManage || selectedIsParentWithChildren"
                    @input="applyFormToSelection"
                  />
                </ElFormItem>

                <ElFormItem
                  v-if="form.type === 'view'"
                  label="跳转地址"
                  required
                >
                  <ElInput
                    v-model="form.url"
                    placeholder="请填写跳转地址"
                    :disabled="!canManage || selectedIsParentWithChildren"
                    @input="applyFormToSelection"
                  />
                </ElFormItem>

                <template v-if="form.type === 'miniprogram'">
                  <ElFormItem label="appid" required>
                    <ElInput
                      v-model="form.appid"
                      placeholder="请填写 appid"
                      :disabled="!canManage || selectedIsParentWithChildren"
                      @input="applyFormToSelection"
                    />
                  </ElFormItem>
                  <ElFormItem label="跳转地址" required>
                    <ElInput
                      v-model="form.url"
                      placeholder="备用网页，不支持小程序时打开"
                      :disabled="!canManage || selectedIsParentWithChildren"
                      @input="applyFormToSelection"
                    />
                  </ElFormItem>
                  <ElFormItem label="小程序路径" required>
                    <ElInput
                      v-model="form.pagepath"
                      placeholder="请填写小程序路径"
                      :disabled="!canManage || selectedIsParentWithChildren"
                      @input="applyFormToSelection"
                    />
                  </ElFormItem>
                </template>
              </ElForm>
            </template>
            <div v-else class="wechat-editor__empty">
              请在左侧选择或添加菜单
            </div>
          </div>

          <div class="wechat-editor__footer">
            <ElButton
              v-if="canManage"
              type="primary"
              :loading="saving"
              @click="saveAndPublish"
            >
              保存并发布
            </ElButton>
          </div>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.wechat-menus {
  min-height: 100%;
}

.wechat-menus__body {
  display: grid;
  grid-template-columns: 390px minmax(0, 1fr);
  gap: 30px;
  align-items: start;
}

.wechat-phone {
  position: relative;
  min-width: 390px;
  min-height: 550px;
  padding-left: 40px;
}

.wechat-phone__frame {
  position: relative;
  width: 320px;
  min-height: 550px;
}

.wechat-phone__head {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 1;
  display: block;
  width: 320px;
  height: 64px;
  pointer-events: none;
  user-select: none;
}

.wechat-phone__body {
  position: absolute;
  top: 63px;
  left: 0;
  z-index: 0;
  width: 320px;
  height: 437px;
  background: #f4f5f9;
}

.wechat-phone__foot {
  position: absolute;
  bottom: 0;
  left: 0;
  z-index: 1;
  display: block;
  width: 320px;
  height: 50px;
  pointer-events: none;
  user-select: none;
}

/* 底图键盘区约到 43px，菜单从分隔线右侧起 */
.wechat-phone__textbot {
  position: absolute;
  bottom: 0;
  left: 44px;
  z-index: 3;
  width: 276px;
  font-size: 14px;
}

.wechat-phone__li {
  position: relative;
  float: left;
  box-sizing: border-box;
  width: 92px;
  color: #353535;
  line-height: 48px;
  text-align: center;
  background: #fafafa;
  border: 1px solid #e7e7eb;
  cursor: pointer;
}

.wechat-phone__li--add {
  color: #999;
}

.wechat-phone__li.is-active {
  color: #44b549 !important;
  border-color: #44b549 !important;
}

.wechat-phone__text {
  height: 48px;
  padding: 0 4px;
  overflow: hidden;
  color: inherit;
  font-size: 14px;
  line-height: 48px;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wechat-phone__text:hover {
  color: #000;
}

.wechat-phone__li.is-active .wechat-phone__text:hover {
  color: #44b549;
}

.wechat-phone__plus {
  display: inline-block;
  color: #44b549;
  font-size: 22px;
  font-weight: 300;
  line-height: 1;
  vertical-align: middle;
}

.wechat-phone__plus--muted {
  color: #999;
  font-size: 20px;
}

/* 添加子菜单条：同宽白底 + 绿色 + + 向下三角（图2） */
.wechat-phone__add {
  position: absolute;
  bottom: 58px;
  left: -1px;
  z-index: 4;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  width: calc(100% + 2px);
  height: 48px;
  color: #44b549;
  background: #fff;
  border: 1px solid #e7e7eb;
  cursor: pointer;
}

.wechat-phone__arrow {
  position: absolute;
  bottom: -12px;
  left: 50%;
  width: 0;
  height: 0;
  margin-left: -6px;
  border-style: solid;
  border-width: 6px 6px 0;
  border-color: #e7e7eb transparent transparent;
}

.wechat-phone__arrow::after {
  position: absolute;
  top: -7px;
  left: -5px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 5px 5px 0;
  border-color: #fff transparent transparent;
  content: '';
}

/* 子菜单列表 */
.wechat-phone__tianjia {
  position: absolute;
  bottom: 106px;
  left: -1px;
  z-index: 5;
  box-sizing: border-box;
  width: calc(100% + 2px);
  background: #fff;
}

.wechat-phone__addadd {
  box-sizing: border-box;
  width: 100%;
  height: 48px;
  overflow: hidden;
  color: #353535;
  font-size: 13px;
  line-height: 46px;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: #fff;
  border: 1px solid #e7e7eb;
  border-bottom: none;
  cursor: pointer;
}

.wechat-phone__addadd:last-child {
  border-bottom: 1px solid #e7e7eb;
}

.wechat-phone__addadd.is-active {
  color: #44b549 !important;
  border-color: #44b549 !important;
}

.wechat-editor {
  display: flex;
  flex-direction: column;
  min-height: 540px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.wechat-editor__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  border-bottom: 1px solid #ebeef5;
}

.wechat-editor__tabs {
  display: flex;
  align-items: center;
}

.wechat-editor__tab {
  position: relative;
  padding: 14px 4px;
  color: var(--el-text-color-regular);
  font-size: 14px;
  font-weight: 500;
}

.wechat-editor__tab.is-active {
  color: var(--el-color-primary);
}

.wechat-editor__tab.is-active::after {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  height: 2px;
  background: var(--el-color-primary);
  content: '';
}

.wechat-editor__body {
  flex: 1;
  padding: 20px 24px 8px;
}

.wechat-editor__alert {
  margin-bottom: 18px;
}

.wechat-editor__alert :deep(.el-alert__close-btn) {
  color: #13ce33 !important;
}

.wechat-editor__form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.wechat-editor__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 280px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.wechat-editor__footer {
  display: flex;
  justify-content: center;
  padding: 16px 0 24px;
  border-top: 1px solid #f0f2f5;
}

@media (max-width: 960px) {
  .wechat-menus__body {
    grid-template-columns: 1fr;
  }
}
</style>
