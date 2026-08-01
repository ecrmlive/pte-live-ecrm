<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { listPlatformAgreementsApi, savePlatformAgreementApi, type PlatformAgreement } from '#/api/core/platform-content';
const loading = ref(false); const saving = ref(false); const rows = ref<PlatformAgreement[]>([]); const current = ref<PlatformAgreement>(); const content = ref('');
async function load() { loading.value = true; try { const data = await listPlatformAgreementsApi(); rows.value = data.list || []; if (!current.value && rows.value.length) select(rows.value[0]); } finally { loading.value = false; } }
function select(row: PlatformAgreement) { current.value = row; content.value = row.content; }
async function save() { if (!current.value) return; if (!content.value.trim()) { ElMessage.warning('协议正文不能为空'); return; } saving.value = true; try { const saved = await savePlatformAgreementApi(current.value.key, content.value.trim()); current.value = saved; const index = rows.value.findIndex((item) => item.key === saved.key); if (index >= 0) rows.value[index] = saved; ElMessage.success('协议已保存'); } finally { saving.value = false; } }
onMounted(() => void load());
</script>
<template>
  <Page title="协议规则" description="维护 C 端协议正文；保存操作受 agreement/update 按钮权限控制。"><div class="grid gap-4 lg:grid-cols-[280px_1fr]"><el-card v-loading="loading" shadow="never"><template #header>协议列表</template><el-menu :default-active="current?.key" @select="(key) => { const row = rows.find((item) => item.key === key); if (row) select(row); }"><el-menu-item v-for="row in rows" :key="row.key" :index="row.key">{{ row.label }}</el-menu-item></el-menu></el-card><el-card shadow="never"><template #header>{{ current?.label || '请选择协议' }}</template><template v-if="current"><el-input v-model="content" :rows="18" type="textarea" placeholder="请输入协议正文" /><div class="mt-4 flex justify-end"><el-button :loading="saving" type="primary" @click="save">保存协议</el-button></div></template></el-card></div></Page>
</template>
