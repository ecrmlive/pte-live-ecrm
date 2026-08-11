<script setup lang="ts">
import CustomerUserPicker, {
  type PickedCustomerUser,
} from '#/views/ecrm/customer-service/customer-user-picker.vue';

/**
 * 兼容既有调用。用户关联一律复用 CustomerUserPicker，禁止新增页面级用户选择弹窗。
 */
export type PickedPlatformUser = PickedCustomerUser;

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  select: [user: PickedPlatformUser];
}>();

function onSelect(user: PickedPlatformUser) {
  emit('select', user);
}
</script>

<template>
  <CustomerUserPicker v-model:open="open" @select="onSelect" />
</template>
