<template>
  <a-card :bordered="false" class="page-card" title="店铺资料">
    <a-form layout="vertical" style="max-width: 520px">
      <a-form-item label="店铺名称" required>
        <a-input v-model:value="form.mer_name" :disabled="!canUpdate" />
      </a-form-item>
      <a-form-item label="联系人">
        <a-input v-model:value="form.real_name" :disabled="!canUpdate" />
      </a-form-item>
      <a-form-item label="手机">
        <a-input v-model:value="form.mer_phone" :disabled="!canUpdate" />
      </a-form-item>
      <a-form-item label="地址">
        <a-input v-model:value="form.mer_address" :disabled="!canUpdate" />
      </a-form-item>
      <a-form-item label="简介">
        <a-textarea v-model:value="form.mer_info" :rows="3" :disabled="!canUpdate" />
      </a-form-item>
      <a-button v-if="canUpdate" type="primary" :loading="saving" @click="submit">保存</a-button>
      <span v-else class="hint">无「保存店铺资料」权限</span>
    </a-form>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchShop, updateShop } from '@/api/setting';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canUpdate = computed(() => auth.hasPerm('shop/update'));

const saving = ref(false);
const form = reactive({
  mer_name: '',
  real_name: '',
  mer_phone: '',
  mer_address: '',
  mer_info: '',
});

async function load() {
  const { data } = await fetchShop();
  form.mer_name = data.mer_name || '';
  form.real_name = data.real_name || '';
  form.mer_phone = data.mer_phone || '';
  form.mer_address = data.mer_address || '';
  form.mer_info = data.mer_info || '';
}

async function submit() {
  if (!canUpdate.value) {
    message.warning('无保存权限');
    return;
  }
  saving.value = true;
  try {
    await updateShop({ ...form });
    message.success('已保存');
    void load();
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<style scoped>
.hint {
  color: #888;
  font-size: 13px;
  margin-left: 8px;
}
</style>
