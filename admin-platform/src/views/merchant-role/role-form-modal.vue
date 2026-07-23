<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue';

import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElTree,
} from 'element-plus';

import MerchantRoleApi from '#/api/core/merchant-role';
import type { MerchantTemplateRoleForm, MerchantTemplateRoleRow } from '#/api/core/merchant-role';
import type { AccessNode } from '#/views/access/types';

const props = defineProps<{
  mode: 'add' | 'edit';
  open: boolean;
  role?: MerchantTemplateRoleRow;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const loading = ref(false);
const menuTree = ref<AccessNode[]>([]);
const checkedKeys = ref<number[]>([]);
const treeRef = ref<InstanceType<typeof ElTree>>();

const formData = reactive<MerchantTemplateRoleForm>({
  role_name: '',
  sort: 100,
  remark: '',
  access_id: [],
});

const formRules = {
  role_name: [{ message: '请输入角色名称', required: true, trigger: 'blur' }],
};

const treeProps = {
  children: 'children',
  label: 'name',
};

const dialogTitle = computed(() =>
  props.mode === 'edit' ? '编辑角色' : '添加角色',
);

async function loadFormData() {
  loading.value = true;
  try {
    if (props.mode === 'edit' && props.role) {
      const res = await MerchantRoleApi.roleEditInfo(props.role.role_id, true);
      const data = res.data as {
        menu?: AccessNode[];
        model?: MerchantTemplateRoleRow;
        select_menu?: number[];
      };
      menuTree.value = data.menu ?? [];
      checkedKeys.value = data.select_menu ?? [];
      Object.assign(formData, {
        role_id: props.role.role_id,
        role_name: data.model?.role_name ?? props.role.role_name,
        sort: data.model?.sort ?? 100,
        remark: data.model?.remark ?? '',
        access_id: [...checkedKeys.value],
      });
    } else {
      const res = await MerchantRoleApi.roleAddInfo(true);
      menuTree.value = (res.data as { menu?: AccessNode[] })?.menu ?? [];
      Object.assign(formData, {
        role_name: '',
        sort: 100,
        remark: '',
        access_id: [],
      });
      checkedKeys.value = [];
    }
    await nextTickSetChecked();
  } finally {
    loading.value = false;
  }
}

async function nextTickSetChecked() {
  await new Promise((r) => setTimeout(r, 0));
  treeRef.value?.setCheckedKeys(checkedKeys.value, false);
}

function handleCheck(_: unknown, state: { checkedKeys: number[]; halfCheckedKeys: number[] }) {
  formData.access_id = [...state.checkedKeys, ...state.halfCheckedKeys];
}

watch(
  () => [props.open, props.mode, props.role?.role_id] as const,
  ([open]) => {
    if (open) {
      loadFormData();
    }
  },
  { immediate: true },
);

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!formData.role_name.trim()) {
    ElMessage.warning('请输入角色名称');
    return;
  }
  if (!formData.access_id.length) {
    ElMessage.warning('请选择权限');
    return;
  }
  loading.value = true;
  try {
    const payload = JSON.stringify(formData);
    const res =
      props.mode === 'edit' && props.role
        ? await MerchantRoleApi.roleEdit(props.role.role_id, payload, true)
        : await MerchantRoleApi.roleAdd(payload, true);
    if (res.code === 1) {
      ElMessage.success(res.msg || '保存成功');
      emit('update:open', false);
      emit('success');
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :model-value="open"
    :title="dialogTitle"
    width="720px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <ElForm v-loading="loading" :model="formData" :rules="formRules" label-width="100px">
      <ElFormItem label="角色名称" prop="role_name">
        <ElInput v-model="formData.role_name" placeholder="请输入角色名称" />
      </ElFormItem>
      <ElFormItem label="权限列表">
        <ElTree
          ref="treeRef"
          :data="menuTree"
          :default-expand-all="true"
          :props="treeProps"
          node-key="access_id"
          show-checkbox
          @check="handleCheck"
        />
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInput v-model.number="formData.sort" type="number" />
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="formData.remark" type="textarea" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">保存</ElButton>
    </template>
  </ElDialog>
</template>
