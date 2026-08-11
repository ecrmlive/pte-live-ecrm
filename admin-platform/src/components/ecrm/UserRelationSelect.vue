<script setup lang="ts">
import { computed, ref } from 'vue';

import { ElButton, ElTag } from 'element-plus';

import CustomerUserPicker, {
  type PickedCustomerUser,
} from '#/views/ecrm/customer-service/customer-user-picker.vue';

const modelValue = defineModel<number | number[] | undefined>({ default: undefined });

const props = withDefaults(
  defineProps<{
    multiple?: boolean;
    placeholder?: string;
  }>(),
  {
    multiple: false,
    placeholder: '请选择用户',
  },
);

const open = ref(false);
const selectedUsers = ref<PickedCustomerUser[]>([]);

const selectedCount = computed(() =>
  Array.isArray(modelValue.value)
    ? modelValue.value.length
    : modelValue.value
      ? 1
      : 0,
);

function selectUser(user: PickedCustomerUser) {
  selectedUsers.value = [user];
  modelValue.value = user.id;
}

function selectUsers(users: PickedCustomerUser[]) {
  selectedUsers.value = users;
  modelValue.value = users.map((user) => user.id);
}

function clear() {
  selectedUsers.value = [];
  modelValue.value = props.multiple ? [] : undefined;
}
</script>

<template>
  <div class="user-relation-select">
    <ElButton plain type="primary" @click="open = true">选择用户</ElButton>
    <span v-if="selectedCount" class="user-relation-select__summary">
      已关联 {{ selectedCount }} 位用户
    </span>
    <ElTag
      v-for="user in selectedUsers"
      :key="user.id"
      closable
      @close="clear"
    >
      {{ user.nickname || user.mobile || '已关联用户' }}
    </ElTag>
  </div>

  <CustomerUserPicker
    v-model:open="open"
    :multiple="props.multiple"
    @select="selectUser"
    @select-multiple="selectUsers"
  />
</template>

<style scoped>
.user-relation-select {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  width: 100%;
}

.user-relation-select__summary {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
