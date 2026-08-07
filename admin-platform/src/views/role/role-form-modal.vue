<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElTree,
} from 'element-plus';

import PlatformRoleApi from '#/api/core/platform-role';
import type { AccessNode } from '#/views/access/types';

import type { PlatformRoleForm, PlatformRoleRow } from '#/api/core/platform-role';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  mode: 'add' | 'edit';
  role?: PlatformRoleRow;
}>();

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const menuTree = ref<AccessNode[]>([]);
const checkedKeys = ref<number[]>([]);
const treeRef = ref<InstanceType<typeof ElTree>>();

const formData = reactive<PlatformRoleForm>({
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

async function nextTickSetChecked() {
  await new Promise((r) => setTimeout(r, 0));
  treeRef.value?.setCheckedKeys(checkedKeys.value, false);
}

async function loadFormData() {
  loading.value = true;
  try {
    if (props.mode === 'edit' && props.role) {
      const res = await PlatformRoleApi.roleEditInfo(props.role.role_id, true);
      const data = res.data as {
        menu?: AccessNode[];
        model?: PlatformRoleRow;
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
      const res = await PlatformRoleApi.roleAddInfo(true);
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

function handleCheck(_: unknown, state: { checkedKeys: number[]; halfCheckedKeys: number[] }) {
  formData.access_id = [...state.checkedKeys, ...state.halfCheckedKeys];
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void loadFormData();
    }
  },
  onConfirm: () => {
    void handleSubmit();
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      drawerApi.setState({ title: dialogTitle.value }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(dialogTitle, (value) => {
  if (open.value) {
    drawerApi.setState({ title: value });
  }
});

watch(
  () => [props.mode, props.role?.role_id] as const,
  () => {
    if (open.value) {
      void loadFormData();
    }
  },
);

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
  drawerApi.setState({ confirmLoading: true });
  try {
    const payload = JSON.stringify(formData);
    const res =
      props.mode === 'edit' && props.role
        ? await PlatformRoleApi.roleEdit(props.role.role_id, payload, true)
        : await PlatformRoleApi.roleAdd(payload, true);
    if (res.code === 1) {
      ElMessage.success(res.msg || '保存成功');
      open.value = false;
      emit('success');
    }
  } finally {
    loading.value = false;
    drawerApi.setState({ confirmLoading: false });
  }
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :confirm-loading="loading"
    :destroy-on-close="true"
    :title="dialogTitle"
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
  </Drawer>
</template>
