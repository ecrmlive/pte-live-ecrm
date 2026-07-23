<script setup lang="ts">
import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { ref, watch } from 'vue';

import { getAgentPosterEditMetaApi } from '#/api/core/plus-agent';

defineOptions({ name: 'AgentPosterPreviewDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  posterId?: number;
}>();

const loading = ref(false);
const imageUrl = ref('');

async function loadPreview() {
  if (!props.posterId) return;
  loading.value = true;
  try {
    const res = await getAgentPosterEditMetaApi(props.posterId);
    const raw = res as Record<string, unknown>;
    const nested = raw.data;
    const model =
      nested && typeof nested === 'object'
        ? (nested as Record<string, unknown>)
        : raw;
    imageUrl.value = String(model.poster_image ?? raw.model?.poster_image ?? '');
  } finally {
    loading.value = false;
  }
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      imageUrl.value = '';
      void loadPreview();
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal :destroy-on-close="true" class="w-[340px]" title="预览">
    <div v-loading="loading" class="flex justify-center py-2">
      <img v-if="imageUrl" :src="imageUrl" alt="poster" class="max-w-[300px]" />
      <span v-else class="text-sm text-muted-foreground">暂无预览图</span>
    </div>
    <template #footer>
      <ElButton type="primary" @click="open = false">关闭</ElButton>
    </template>
  </Modal>
</template>
