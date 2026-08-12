<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  clearMaintainCacheApi,
  type MaintainCacheScope,
} from '#/api/core/platform-maintain';

const canManage = ref(false);
const processing = ref<MaintainCacheScope | null>(null);
const domainDialogVisible = ref(false);
const domainForm = reactive({ oldDomain: '', newDomain: '' });

const cards = computed(() => [
  {
    action: 'replace_domain' as const,
    button: '立即更换',
    description: '替换所有本地上传的图片域名',
    title: '更换域名',
    type: 'primary' as const,
  },
  {
    action: 'all' as const,
    button: '立即清理',
    description: '清除所有商城运行缓存，谨慎操作',
    title: '清除全部缓存',
    type: 'default' as const,
  },
  {
    action: 'store' as const,
    button: '立即清理',
    description: '清除所有店铺运行缓存，谨慎操作',
    title: '清理店铺缓存',
    type: 'default' as const,
  },
  {
    action: 'config' as const,
    button: '立即清理',
    description: '清除配置运行缓存，谨慎操作',
    title: '清理配置缓存',
    type: 'default' as const,
  },
]);

function openDomainDialog() {
  domainForm.oldDomain = '';
  domainForm.newDomain = '';
  domainDialogVisible.value = true;
}

async function runCacheAction(scope: Exclude<MaintainCacheScope, 'replace_domain'>) {
  const title = cards.value.find((card) => card.action === scope)?.title ?? '清除缓存';
  try {
    await ElMessageBox.confirm(
      '将只清理本平台运行缓存，不会删除商城、店铺、商品、用户等业务数据。',
      title,
      { cancelButtonText: '取消', confirmButtonText: '立即清理', type: 'warning' },
    );
    processing.value = scope;
    const result = await clearMaintainCacheApi({ scope });
    ElMessage.success(`${result.note}（${result.deleted_keys} 项）`);
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 请求客户端已展示服务端错误。
    }
  } finally {
    processing.value = null;
  }
}

async function saveDomainReplacement() {
  try {
    processing.value = 'replace_domain';
    const result = await clearMaintainCacheApi({
      new_domain: domainForm.newDomain,
      old_domain: domainForm.oldDomain,
      scope: 'replace_domain',
    });
    domainDialogVisible.value = false;
    ElMessage.success(`${result.note}（已更新 ${result.updated_assets ?? 0} 个素材）`);
  } finally {
    processing.value = null;
  }
}

onMounted(async () => {
  canManage.value = (await getAccessCodesApi()).includes('maintain.cache.manage');
});
</script>

<template>
  <Page class="cache-page">
    <section class="cache-page__title">
      <h2>清除缓存</h2>
      <p><span class="cache-page__warning-icon">i</span> 清除数据请谨慎，清除就无法恢复哦！</p>
    </section>

    <section class="cache-page__content">
      <div class="cache-page__grid">
        <article v-for="card in cards" :key="card.action" class="cache-page__card">
          <h3>{{ card.title }}</h3>
          <p>{{ card.description }}</p>
          <el-button
            :disabled="!canManage"
            :loading="processing === card.action"
            :type="card.type"
            @click="card.action === 'replace_domain' ? openDomainDialog() : runCacheAction(card.action)"
          >
            {{ card.button }}
          </el-button>
        </article>
      </div>
      <p v-if="!canManage" class="cache-page__permission">当前账号无缓存管理权限。</p>
    </section>

    <el-dialog v-model="domainDialogVisible" title="更换资源域名" width="560px" destroy-on-close>
      <el-form label-width="108px" @submit.prevent="saveDomainReplacement">
        <el-form-item label="原资源域名" required>
          <el-input v-model="domainForm.oldDomain" placeholder="例如：https://old.example.com" />
        </el-form-item>
        <el-form-item label="新资源域名" required>
          <el-input v-model="domainForm.newDomain" placeholder="例如：https://cdn.example.com" />
        </el-form-item>
        <p class="cache-page__domain-help">仅替换以原资源域名开头的系统素材地址。</p>
      </el-form>
      <template #footer>
        <div class="cache-page__dialog-actions">
          <el-button @click="domainDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="processing === 'replace_domain'" @click="saveDomainReplacement">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </Page>
</template>

<style scoped>
.cache-page :deep(.vben-page__content) {
  padding: 0;
}

.cache-page__title,
.cache-page__content {
  border-radius: 8px;
  background: var(--el-bg-color);
}

.cache-page__title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 98px;
  padding: 0 28px;
}

.cache-page__title h2 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 22px;
  font-weight: 500;
}

.cache-page__title p {
  display: flex;
  align-items: center;
  margin: 0;
  color: var(--el-color-danger);
  font-size: 14px;
}

.cache-page__warning-icon {
  display: inline-flex;
  width: 16px;
  height: 16px;
  align-items: center;
  justify-content: center;
  margin-right: 6px;
  border-radius: 50%;
  background: var(--el-color-danger);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  font-style: italic;
}

.cache-page__content {
  min-height: 560px;
  margin-top: 16px;
  padding: 28px;
}

.cache-page__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(240px, 1fr));
  gap: 28px 34px;
}

.cache-page__card {
  display: flex;
  min-height: 232px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  text-align: center;
}

.cache-page__card h3 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 20px;
  font-weight: 500;
}

.cache-page__card p {
  margin: 20px 0;
  color: var(--el-text-color-secondary);
  font-size: 15px;
}

.cache-page__permission,
.cache-page__domain-help {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.cache-page__dialog-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

@media (max-width: 1100px) {
  .cache-page__grid {
    grid-template-columns: repeat(2, minmax(240px, 1fr));
  }
}
</style>
