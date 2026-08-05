<script setup lang="ts">
import type { MerchantStaff } from '#/api/core/staff';
import { Page } from '@vben/common-ui';
import { onMounted, ref } from 'vue';
import { listMerchantStaffApi } from '#/api/core/staff';

const rows = ref<MerchantStaff[]>([]);
const loading = ref(false);

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantStaffApi({ page: 1, limit: 100 });
    rows.value = (data.list ?? []).filter((item) => item.role_code === 'delivery' || item.is_goods === 1);
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="配送员管理" description="具备发货/配送能力的店铺员工。">
    <el-card v-loading="loading" shadow="never">
      <el-table :data="rows">
        <el-table-column prop="service_id" label="ID" width="80" />
        <el-table-column prop="nickname" label="姓名" min-width="140" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="role_code" label="角色" width="100" />
        <el-table-column label="可发货" width="90"><template #default="{ row }">{{ row.is_goods ? '是' : '否' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
