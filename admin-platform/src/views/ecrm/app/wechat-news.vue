<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElEmpty,
  ElImage,
  ElInput,
  ElMessage,
  ElPagination,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  deleteWechatNewsApi,
  listWechatNewsApi,
  type WechatNewsRow,
} from '#/api/core/platform-wechat-news';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

import WechatNewsSave from './wechat-news-save.vue';

const loading = ref(false);
const canManage = ref(false);
const list = ref<WechatNewsRow[]>([]);
const total = ref(0);
const keyword = ref('');
const pager = reactive({ page: 1, limit: 10 });
const editingId = ref(0);
const editorRef = ref<InstanceType<typeof WechatNewsSave>>();

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[1200px] max-w-[98vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => {
    const saved = await editorRef.value?.submit();
    if (!saved) return;
    formDrawerApi.close();
    await load();
  },
});

async function load() {
  loading.value = true;
  try {
    const data = await listWechatNewsApi({
      page: pager.page,
      limit: pager.limit,
      cate_name: keyword.value.trim() || undefined,
    });
    list.value = data.list || [];
    total.value = Number(data.count || 0);
  } catch {
    list.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function coverOf(row: WechatNewsRow) {
  return resolveCosMediaUrl(row.article?.[0]?.image || '');
}

function titleOf(row: WechatNewsRow) {
  return row.article?.[0]?.title || '未命名图文';
}

function goCreate() {
  editingId.value = 0;
  formDrawerApi.setState({ title: '新增图文消息' });
  formDrawerApi.open();
}

function goEdit(id: number) {
  editingId.value = id;
  formDrawerApi.setState({ title: '编辑图文消息' });
  formDrawerApi.open();
}

async function onDelete(row: WechatNewsRow) {
  if (!canManage.value) return;
  try {
    await confirm({
      title: '删除确认',
      content: `确定删除「${titleOf(row)}」吗？`,
    });
  } catch {
    return;
  }
  await deleteWechatNewsApi(row.wechat_news_id);
  ElMessage.success('已删除');
  await load();
}

function onSearch() {
  pager.page = 1;
  void load();
}

function onReset() {
  keyword.value = '';
  pager.page = 1;
  void load();
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value = codes.includes('app.wechat_news.manage');
  await load();
});
</script>

<template>
  <Page auto-content-height>
    <div v-loading="loading" class="wechat-news">
      <section class="wechat-news__filter-card">
        <div class="wechat-news__toolbar">
          <div class="wechat-news__filters">
            <ElInput
              v-model="keyword"
              clearable
              class="wechat-news__keyword"
              placeholder="请输入图文标题"
              @keyup.enter="onSearch"
            />
            <ElButton @click="onReset">重置</ElButton>
            <ElButton type="primary" @click="onSearch">搜索</ElButton>
          </div>
        </div>
      </section>

      <section class="wechat-news__list-card">
        <div class="wechat-news__list-toolbar">
          <ElButton v-if="canManage" type="primary" @click="goCreate">
            新增图文消息
          </ElButton>
        </div>
        <div v-if="list.length" class="wechat-news__grid">
          <div
            v-for="row in list"
            :key="row.wechat_news_id"
            class="wechat-news-card"
          >
            <div
              class="wechat-news-card__cover"
              @click="goEdit(row.wechat_news_id)"
            >
              <ElImage
                v-if="coverOf(row)"
                :src="coverOf(row)"
                fit="cover"
                class="wechat-news-card__img"
              />
              <div v-else class="wechat-news-card__placeholder">暂无封面</div>
              <div class="wechat-news-card__title">{{ titleOf(row) }}</div>
            </div>
            <div
              v-for="(item, idx) in (row.article || []).slice(1)"
              :key="`${row.wechat_news_id}-${idx}`"
              class="wechat-news-card__sub"
              @click="goEdit(row.wechat_news_id)"
            >
              <span class="wechat-news-card__sub-title">{{ item.title }}</span>
              <ElImage
                v-if="item.image"
                :src="resolveCosMediaUrl(item.image)"
                fit="cover"
                class="wechat-news-card__sub-img"
              />
            </div>
            <div class="wechat-news-card__foot">
              <span class="wechat-news-card__time">
                {{ formatShanghaiDateTime(row.create_time) || '—' }}
              </span>
              <div class="wechat-news-card__actions">
                <ElButton
                  link
                  type="primary"
                  @click="goEdit(row.wechat_news_id)"
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
              </div>
            </div>
          </div>
        </div>
        <ElEmpty v-else description="暂无图文消息" />

        <div v-if="total > 0" class="wechat-news__pager">
          <ElPagination
            v-model:current-page="pager.page"
            v-model:page-size="pager.limit"
            background
            layout="total, prev, pager, next"
            :total="total"
            @current-change="load"
          />
        </div>
      </section>
    </div>

    <FormDrawer>
      <WechatNewsSave
        ref="editorRef"
        :can-manage="canManage"
        :news-id="editingId"
      />
    </FormDrawer>
  </Page>
</template>

<style scoped>
.wechat-news {
  min-height: 100%;
}

.wechat-news__filter-card,
.wechat-news__list-card {
  overflow: hidden;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.wechat-news__filter-card {
  margin-bottom: 16px;
}

.wechat-news__toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  padding: 16px 20px;
}

.wechat-news__list-toolbar {
  display: flex;
  align-items: center;
  min-height: 56px;
  padding: 12px 20px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.wechat-news__filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.wechat-news__keyword {
  width: 220px;
}

.wechat-news__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  padding: 20px;
}

.wechat-news-card {
  overflow: hidden;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.wechat-news-card__cover {
  position: relative;
  height: 160px;
  cursor: pointer;
  background: #f5f7fa;
}

.wechat-news-card__img {
  width: 100%;
  height: 100%;
}

.wechat-news-card__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #c0c4cc;
  font-size: 13px;
}

.wechat-news-card__title {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  padding: 8px 12px;
  overflow: hidden;
  color: #fff;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: linear-gradient(transparent, rgb(0 0 0 / 55%));
}

.wechat-news-card__sub {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-top: 1px dashed #eee;
  cursor: pointer;
}

.wechat-news-card__sub-title {
  flex: 1;
  overflow: hidden;
  color: #303133;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wechat-news-card__sub-img {
  flex: none;
  width: 48px;
  height: 48px;
  border-radius: 4px;
}

.wechat-news-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-top: 1px solid #f0f2f5;
}

.wechat-news-card__time {
  color: #909399;
  font-size: 12px;
}

.wechat-news-card__actions {
  display: flex;
  gap: 4px;
}

.wechat-news__pager {
  display: flex;
  justify-content: flex-end;
  min-height: 60px;
  padding: 14px 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
