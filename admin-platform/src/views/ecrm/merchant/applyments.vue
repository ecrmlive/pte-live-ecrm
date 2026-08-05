<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { fetchProfitsharingApplications, reviewProfitsharingApplication, saveProfitsharingApplicationNote, type ProfitsharingApplication } from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';

const rows = ref<ProfitsharingApplication[]>([]);
const can = ref(false);
const loading = ref(false);
const validateNote = (value: string) => {
  const note = value.trim();
  return note && [...note].length <= 500 ? true : '审核说明不能为空，且不能超过 500 个字符。';
};
const isPromptDismissed = (error: unknown) => error === 'cancel' || error === 'close' || error === 'escape';

async function load() {
  loading.value = true;
  try {
    rows.value = (await fetchProfitsharingApplications()).list || [];
  } finally {
    loading.value = false;
  }
}

async function review(row: ProfitsharingApplication, approved: boolean) {
  try {
    const { value } = await ElMessageBox.prompt('填写审核说明。', approved ? '同意分账申请' : '拒绝分账申请', { inputValidator: validateNote });
    await reviewProfitsharingApplication(row.id, approved, value.trim());
    ElMessage.success('审核已保存');
    await load();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function note(row: ProfitsharingApplication) {
  try {
    const { value } = await ElMessageBox.prompt('填写内部审核备注。', '分账申请备注', { inputValue: row.review_note, inputValidator: validateNote });
    await saveProfitsharingApplicationNote(row.id, value.trim());
    ElMessage.success('备注已保存');
    await load();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

onMounted(async () => {
  const [codes] = await Promise.all([getAccessCodesApi(), load()]);
  can.value = codes.includes('merchant.profitsharing.review');
});
</script>

<template>
  <Page title="店铺分账申请" description="平台审核商户分账申请；页面不展示或保存渠道账户资料。">
    <el-table v-loading="loading" :data="rows">
      <el-table-column prop="application_no" label="申请编号" />
      <el-table-column prop="merchant_id" label="商户 ID" />
      <el-table-column prop="description" label="申请说明" min-width="260" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="review_note" label="审核备注" min-width="180" />
      <el-table-column v-if="can" label="操作" width="210">
        <template #default="{ row }">
          <template v-if="row.status === 'applied'">
            <el-button link type="success" @click="review(row, true)">同意</el-button>
            <el-button link type="danger" @click="review(row, false)">拒绝</el-button>
          </template>
          <el-button link type="primary" @click="note(row)">备注</el-button>
        </template>
      </el-table-column>
    </el-table>
  </Page>
</template>
