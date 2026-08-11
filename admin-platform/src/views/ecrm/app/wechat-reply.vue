<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createWechatReplyApi,
  deleteWechatReplyApi,
  getWechatReplySpecialApi,
  listWechatRepliesApi,
  matchWechatReplyApi,
  saveWechatReplySpecialApi,
  setWechatReplyStatusApi,
  updateWechatReplyApi,
  type WechatReplyRow,
} from '#/api/core/platform-wechat-reply';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const canRead = ref(false);
const canManage = ref(false);
const activeTab = ref('subscribe');
const editingId = ref(0);

const subscribeForm = reactive({
  content: '',
  status: true,
});
const defaultForm = reactive({
  content: '',
  status: true,
});
const specialSaving = ref(false);
const specialLoading = ref(false);

const previewKey = ref('1');
const previewText = ref('');
const previewLoading = ref(false);

const keywordForm = reactive({
  key: '',
  content: '',
  status: true,
  sort: 0,
});

const gridOptions: VxeGridProps<WechatReplyRow> = {
  columns: [
    { field: 'wechat_reply_id', title: 'ID', width: 80 },
    {
      field: 'key',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '关键字',
    },
    {
      field: 'content',
      minWidth: 260,
      showOverflow: 'tooltip',
      title: '回复内容',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'sort',
      title: '排序',
      width: 80,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 160 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await listWechatRepliesApi({
          page: page.currentPage,
          limit: page.pageSize,
          kind: 'keyword',
        });
        return {
          items: data.list || [],
          total: Number(data.count || 0),
        };
      },
    },
  },
  toolbarConfig: { search: false },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => saveKeyword(),
});

async function loadSpecial(key: 'subscribe' | 'default') {
  specialLoading.value = true;
  try {
    const row = await getWechatReplySpecialApi(key);
    const target = key === 'subscribe' ? subscribeForm : defaultForm;
    target.content = row.content || '';
    target.status = Number(row.status) === 1;
  } finally {
    specialLoading.value = false;
  }
}

async function saveSpecial(key: 'subscribe' | 'default') {
  if (!canManage.value) return;
  const target = key === 'subscribe' ? subscribeForm : defaultForm;
  if (!target.content.trim()) {
    ElMessage.warning('请填写回复内容');
    return;
  }
  specialSaving.value = true;
  try {
    await saveWechatReplySpecialApi(key, {
      content: target.content.trim(),
      status: target.status ? 1 : 0,
      type: 'text',
    });
    ElMessage.success('已保存');
    await loadSpecial(key);
  } finally {
    specialSaving.value = false;
  }
}

function openCreate() {
  editingId.value = 0;
  Object.assign(keywordForm, {
    key: '',
    content: '',
    status: true,
    sort: 0,
  });
  formDrawerApi.setState({ title: '新增关键字回复' });
  formDrawerApi.open();
}

function openEdit(row: WechatReplyRow) {
  editingId.value = row.wechat_reply_id;
  Object.assign(keywordForm, {
    key: row.key || '',
    content: row.content || '',
    status: Number(row.status) === 1,
    sort: Number(row.sort || 0),
  });
  formDrawerApi.setState({ title: '编辑关键字回复' });
  formDrawerApi.open();
}

async function saveKeyword() {
  if (!canManage.value) return;
  const key = keywordForm.key.trim();
  const content = keywordForm.content.trim();
  if (!key) {
    ElMessage.warning('请填写关键字');
    return;
  }
  if (key === 'subscribe' || key === 'default') {
    ElMessage.warning('关键字不能使用 subscribe / default');
    return;
  }
  if (!content) {
    ElMessage.warning('请填写回复内容');
    return;
  }
  formDrawerApi.setState({ confirming: true });
  try {
    const payload = {
      key,
      content,
      type: 'text',
      status: keywordForm.status ? 1 : 0,
      sort: Number(keywordForm.sort || 0),
    };
    if (editingId.value) {
      await updateWechatReplyApi(editingId.value, payload);
      ElMessage.success('已更新');
    } else {
      await createWechatReplyApi(payload);
      ElMessage.success('已新增');
    }
    formDrawerApi.close();
    await gridApi.reload();
  } finally {
    formDrawerApi.setState({ confirming: false });
  }
}

async function onToggleStatus(row: WechatReplyRow, enabled: boolean) {
  if (!canManage.value) return;
  await setWechatReplyStatusApi(row.wechat_reply_id, enabled ? 1 : 0);
  ElMessage.success(enabled ? '已启用' : '已停用');
  await gridApi.reload();
}

async function onDelete(row: WechatReplyRow) {
  if (!canManage.value) return;
  try {
    await confirm({
      title: '删除确认',
      content: `确定删除关键字「${row.key}」吗？`,
    });
  } catch {
    return;
  }
  await deleteWechatReplyApi(row.wechat_reply_id);
  ElMessage.success('已删除');
  await gridApi.reload();
}

async function runPreview() {
  const key = previewKey.value.trim();
  if (!key) {
    ElMessage.warning('请输入要模拟的用户消息');
    return;
  }
  previewLoading.value = true;
  try {
    const data = await matchWechatReplyApi(key);
    if (!data.matched || !data.reply) {
      previewText.value = '（无匹配回复）';
      return;
    }
    previewText.value = data.reply.content || '';
  } finally {
    previewLoading.value = false;
  }
}

async function onTabChange(name: string | number) {
  const tab = String(name);
  if (tab === 'subscribe' || tab === 'default') {
    await loadSpecial(tab);
  } else if (tab === 'keyword') {
    await gridApi.reload();
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canRead.value =
    codes.includes('app.wechat_reply.read') ||
    codes.includes('app.wechat_reply.manage');
  canManage.value = codes.includes('app.wechat_reply.manage');
  if (canRead.value) {
    await loadSpecial('subscribe');
  }
});
</script>

<template>
  <Page auto-content-height>
    <div class="wechat-reply">
      <ElTabs
        v-model="activeTab"
        class="wechat-reply__tabs"
        @tab-change="onTabChange"
      >
        <ElTabPane label="关注回复" name="subscribe">
          <div v-loading="specialLoading" class="wechat-reply__panel">
            <div class="wechat-reply__editor-heading">
              <div>
                <h2>关注欢迎语</h2>
                <p>新用户关注公众号后，会立即收到这条消息。</p>
              </div>
              <div class="wechat-reply__switch-box">
                <div>
                  <span>自动发送</span>
                  <small>{{ subscribeForm.status ? '已启用' : '已停用' }}</small>
                </div>
                <ElSwitch
                  v-model="subscribeForm.status"
                  :disabled="!canManage"
                />
              </div>
            </div>
            <ElForm
              label-position="top"
              :disabled="!canManage"
              class="wechat-reply__editor-form"
            >
              <ElFormItem label="欢迎内容" required>
                <ElInput
                  v-model="subscribeForm.content"
                  type="textarea"
                  :rows="10"
                  maxlength="2000"
                  show-word-limit
                  placeholder="建议包含数字菜单，例如：&#10;欢迎关注！请回复：&#10;1. 商城首页&#10;2. 热门活动&#10;3. 售后须知&#10;4. 联系客服"
                />
              </ElFormItem>
              <div class="wechat-reply__editor-footer">
                <p>建议以简短问候开头，再使用换行和数字菜单引导用户操作。</p>
                <ElButton
                  v-if="canManage"
                  type="primary"
                  :loading="specialSaving"
                  @click="saveSpecial('subscribe')"
                >
                  保存关注回复
                </ElButton>
              </div>
            </ElForm>
          </div>
        </ElTabPane>

        <ElTabPane label="关键字回复" name="keyword">
          <div class="wechat-reply__keyword wechat-reply__content-card">
            <div class="wechat-reply__editor-heading wechat-reply__editor-heading--compact">
              <div>
                <h2>关键字规则</h2>
                <p>用户消息命中规则后，按排序优先级返回对应内容。</p>
              </div>
            </div>
            <Grid>
              <template #toolbar-actions>
                <ElButton v-if="canManage" type="primary" @click="openCreate">
                  新增关键字
                </ElButton>
              </template>
              <template #status="{ row }">
                <ElSwitch
                  :model-value="Number(row.status) === 1"
                  :disabled="!canManage"
                  @change="(val: boolean) => onToggleStatus(row, val)"
                />
              </template>
              <template #action="{ row }">
                <ElButton
                  v-if="canManage"
                  link
                  type="primary"
                  @click="openEdit(row)"
                >
                  编辑
                </ElButton>
                <ElButton
                  v-if="canManage"
                  link
                  type="danger"
                  @click="onDelete(row)"
                >
                  删除
                </ElButton>
              </template>
            </Grid>
          </div>
        </ElTabPane>

        <ElTabPane label="默认回复" name="default">
          <div v-loading="specialLoading" class="wechat-reply__panel">
            <div class="wechat-reply__editor-heading">
              <div>
                <h2>未命中兜底回复</h2>
                <p>用户消息没有匹配关键字规则时，会发送这条消息。</p>
              </div>
              <div class="wechat-reply__switch-box">
                <div>
                  <span>自动发送</span>
                  <small>{{ defaultForm.status ? '已启用' : '已停用' }}</small>
                </div>
                <ElSwitch
                  v-model="defaultForm.status"
                  :disabled="!canManage"
                />
              </div>
            </div>
            <ElForm
              label-position="top"
              :disabled="!canManage"
              class="wechat-reply__editor-form"
            >
              <ElFormItem label="默认内容" required>
                <ElInput
                  v-model="defaultForm.content"
                  type="textarea"
                  :rows="8"
                  maxlength="2000"
                  show-word-limit
                  placeholder="用户消息未命中任何关键字时的兜底回复"
                />
              </ElFormItem>
              <div class="wechat-reply__editor-footer">
                <p>可提示用户回复菜单数字，或引导进入人工客服。</p>
                <ElButton
                  v-if="canManage"
                  type="primary"
                  :loading="specialSaving"
                  @click="saveSpecial('default')"
                >
                  保存默认回复
                </ElButton>
              </div>
            </ElForm>
          </div>
        </ElTabPane>

        <ElTabPane label="本地预览" name="preview">
          <div class="wechat-reply__panel wechat-reply__preview-panel">
            <div class="wechat-reply__editor-heading wechat-reply__editor-heading--compact">
              <div>
                <h2>本地消息模拟</h2>
                <p>输入用户消息，确认关键字规则与默认回复是否符合预期。</p>
              </div>
            </div>
            <ElForm label-position="top" @submit.prevent>
              <ElFormItem label="用户消息">
                <div class="flex w-full gap-2">
                  <ElInput
                    v-model="previewKey"
                    clearable
                    placeholder="例如输入 1 / 2 / 你好"
                    @keyup.enter="runPreview"
                  />
                  <ElButton
                    type="primary"
                    :loading="previewLoading"
                    @click="runPreview"
                  >
                    模拟回复
                  </ElButton>
                </div>
              </ElFormItem>
              <ElFormItem label="系统回复">
                <div class="wechat-reply__preview">
                  <template v-if="previewText">
                    <ElTag size="small" class="mb-2" type="primary">text</ElTag>
                    <pre class="wechat-reply__preview-text">{{ previewText }}</pre>
                  </template>
                  <span v-else class="text-[var(--el-text-color-secondary)]">
                    输入用户消息后点击「模拟回复」
                  </span>
                </div>
              </ElFormItem>
            </ElForm>
          </div>
        </ElTabPane>
      </ElTabs>
    </div>

    <FormDrawer>
      <ElForm label-width="96px" class="p-4" :disabled="!canManage">
        <ElFormItem label="关键字" required>
          <ElInput
            v-model="keywordForm.key"
            maxlength="64"
            show-word-limit
            placeholder="如 1、2、商城、活动"
          />
        </ElFormItem>
        <ElFormItem label="回复内容" required>
          <ElInput
            v-model="keywordForm.content"
            type="textarea"
            :rows="8"
            maxlength="2000"
            show-word-limit
            placeholder="用户回复该关键字后推送的文本"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="keywordForm.sort" :min="0" :max="9999" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="keywordForm.status" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.wechat-reply {
  max-width: 1440px;
  margin: 0 auto;
}

.wechat-reply__tabs {
  padding: 0 24px 24px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
}

.wechat-reply__tabs :deep(.el-tabs__header) {
  margin: 0 -24px 24px;
  padding: 0 24px;
  background: var(--el-fill-color-lighter);
  border-bottom: 1px solid var(--el-border-color-lighter);
  border-radius: 12px 12px 0 0;
}

.wechat-reply__tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.wechat-reply__tabs :deep(.el-tabs__item) {
  height: 58px;
  padding: 0 22px;
  font-size: 15px;
  font-weight: 600;
}

.wechat-reply__tabs :deep(.el-tabs__active-bar) {
  height: 3px;
  border-radius: 3px 3px 0 0;
}

.wechat-reply__panel {
  max-width: 1040px;
  padding: 4px 0 8px;
}

.wechat-reply__content-card {
  min-height: 480px;
  padding: 4px 0 0;
}

.wechat-reply__editor-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 22px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.wechat-reply__editor-heading--compact {
  margin-bottom: 16px;
}

.wechat-reply__editor-heading h2 {
  margin: 0;
  font-size: 17px;
  line-height: 1.4;
}

.wechat-reply__editor-heading p {
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.6;
}

.wechat-reply__switch-box {
  display: flex;
  gap: 12px;
  align-items: center;
  min-width: 150px;
  padding: 10px 12px;
  background: var(--el-fill-color-lighter);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.wechat-reply__switch-box span,
.wechat-reply__switch-box small {
  display: block;
}

.wechat-reply__switch-box span {
  color: var(--el-text-color-primary);
  font-size: 13px;
  font-weight: 600;
}

.wechat-reply__switch-box small {
  margin-top: 2px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.wechat-reply__editor-form :deep(.el-form-item__label) {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.wechat-reply__editor-form :deep(.el-textarea__inner) {
  min-height: 260px !important;
  padding: 14px 16px;
  line-height: 1.75;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.wechat-reply__editor-form :deep(.el-textarea__inner:focus) {
  background: var(--el-bg-color);
}

.wechat-reply__editor-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding-top: 2px;
}

.wechat-reply__editor-footer p {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.wechat-reply__editor-footer :deep(.el-button) {
  min-width: 132px;
  height: 38px;
  border-radius: 7px;
  font-weight: 600;
}

.wechat-reply__preview-panel {
  max-width: 820px;
}

.wechat-reply__preview-panel :deep(.el-form-item__label) {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.wechat-reply__preview {
  width: 100%;
  min-height: 168px;
  padding: 16px;
  background:
    linear-gradient(135deg, rgb(64 158 255 / 7%), transparent 45%),
    var(--el-fill-color-lighter);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.wechat-reply__preview-text {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  color: #303133;
  font-family: inherit;
  font-size: 13px;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .wechat-reply__editor-heading,
  .wechat-reply__editor-footer {
    align-items: flex-start;
    flex-direction: column;
  }

  .wechat-reply__tabs {
    padding-right: 16px;
    padding-left: 16px;
  }

  .wechat-reply__tabs :deep(.el-tabs__header) {
    margin-right: -16px;
    margin-left: -16px;
    padding: 0 16px;
    overflow-x: auto;
  }

  .wechat-reply__tabs :deep(.el-tabs__item) {
    padding: 0 14px;
  }

  .wechat-reply__switch-box {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
