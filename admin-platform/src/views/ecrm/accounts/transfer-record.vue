<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { fetchTransferRecords } from '#/api/core/platform-maintain';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<Array<{ id: string; name: string; enabled: boolean; remark: string }>>([]);

async function load() {
  loading.value = true;
  try {
    rows.value = (await fetchTransferRecords()).list || [];
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="转账记录" description="转账记录监管 stub；真实打款流水待 finance 域 API 接入。">
    <el-alert class="mb-4" title="当前为只读监管视图；写入与打款凭据不在后台保存或回显。" type="warning" :closable="false" />
    <EcrmListPage title="转账记录">
      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column label="标识" prop="id" min-width="120" />
        <el-table-column label="名称" prop="name" min-width="160" />
        <el-table-column label="备注" prop="remark" min-width="220" show-overflow-tooltip />
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '有效' : '停用' }}</el-tag></template></el-table-column>
      </el-table>
      <el-empty v-if="!loading && !rows.length" description="暂无转账记录 stub 数据" />
    </EcrmListPage>
  </Page>
</template>
