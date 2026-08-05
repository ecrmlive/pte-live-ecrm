<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import { clearMaintainCacheApi } from '#/api/core/platform-maintain';

const clearing = ref(false);
const canManage = ref(false);
const lastNote = ref('');

async function clearCache() {
  try {
    await ElMessageBox.confirm('确认提交缓存清理请求？此操作不会回显或保存任何密钥。', '清除缓存', { type: 'warning' });
    clearing.value = true;
    const result = await clearMaintainCacheApi();
    lastNote.value = result.note;
    ElMessage.success('缓存清理请求已提交');
  } catch {
    // 取消
  } finally {
    clearing.value = false;
  }
}

onMounted(async () => {
  canManage.value = (await getAccessCodesApi()).includes('maintain.cache.manage');
});
</script>

<template>
  <Page title="清除缓存" description="提交平台缓存清理 stub；真实 Redis/应用缓存清理策略待运维域接入。">
    <el-card shadow="never">
      <el-alert class="mb-4" title="当前为监管 stub：点击后将记录清理意图，不含密钥或敏感凭据。" type="warning" :closable="false" />
      <el-alert v-if="lastNote" class="mb-4" :title="lastNote" type="success" :closable="false" />
      <el-button v-if="canManage" :loading="clearing" type="danger" @click="clearCache">清除缓存</el-button>
      <el-alert v-else title="当前账号无缓存清理权限。" type="info" :closable="false" />
    </el-card>
  </Page>
</template>
