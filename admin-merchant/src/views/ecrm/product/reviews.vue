<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { listMerchantProductCommentsApi, replyMerchantProductCommentApi, type MerchantProductComment } from '#/api/core/merchant-product-comment';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<MerchantProductComment[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20, product_id: undefined as number | undefined, status: undefined as string | undefined });
const replyOpen = ref(false);
const replySaving = ref(false);
const current = ref<MerchantProductComment>();
const replyContent = ref('');

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantProductCommentsApi({ ...query });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function openReply(row: MerchantProductComment) {
  current.value = row;
  replyContent.value = row.reply_content || '';
  replyOpen.value = true;
}

async function saveReply() {
  if (!current.value || !replyContent.value.trim()) {
    ElMessage.warning('请填写回复内容');
    return;
  }
  replySaving.value = true;
  try {
    await replyMerchantProductCommentApi(current.value.id, { reply_content: replyContent.value.trim() });
    replyOpen.value = false;
    ElMessage.success('商家回复已保存');
    await load();
  } finally {
    replySaving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="商品评论" description="查看本店商品评论并回复买家。">
    <EcrmListPage title="评论列表">
      <template #filters>
        <el-form inline @submit.prevent="query.page = 1; load()">
          <el-form-item label="商品 ID"><el-input-number v-model="query.product_id" :min="1" /></el-form-item>
          <el-form-item label="状态"><el-select v-model="query.status" clearable><el-option label="待审核" value="pending" /><el-option label="已展示" value="published" /><el-option label="已隐藏" value="hidden" /></el-select></el-form-item>
        </el-form>
      </template>
      <template #actions><el-button type="primary" @click="query.page = 1; load()">查询</el-button></template>
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="商品" min-width="180"><template #default="{ row }">{{ row.product_title || `#${row.product_id}` }}</template></el-table-column>
        <el-table-column prop="user_id" label="用户 ID" width="90" />
        <el-table-column label="评分" width="120"><template #default="{ row }"><el-rate :model-value="row.score" disabled /></template></el-table-column>
        <el-table-column prop="content" label="评论" min-width="220" show-overflow-tooltip />
        <el-table-column prop="reply_content" label="回复" min-width="160" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90" />
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column label="操作" width="90" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="openReply(row)">回复</el-button></template></el-table-column>
      </el-table>
      <template #pager><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></template>
    </EcrmListPage>
    <el-dialog v-model="replyOpen" title="商家回复" width="520px"><el-input v-model="replyContent" type="textarea" :rows="4" maxlength="500" show-word-limit /><template #footer><el-button @click="replyOpen = false">取消</el-button><el-button type="primary" :loading="replySaving" @click="saveReply">保存</el-button></template></el-dialog>
  </Page>
</template>
