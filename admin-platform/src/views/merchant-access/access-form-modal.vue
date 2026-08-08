<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue';

import { IconPicker, useVbenDrawer } from '@vben/common-ui';
import {
  ElCascader,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
} from 'element-plus';

import MerchantAccessApi from '#/api/core/merchant-access';
import { PLATFORM_LUCIDE_ICONS } from '#/constants/platform-lucide-icons';
import { deepClone, formatModel } from '#/utils/base';

import type { AccessAddType, AccessFormModel, AccessNode } from '#/views/access/types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  addType?: AccessAddType;
  mode: 'add' | 'edit';
  rawData: AccessNode[];
  selectModel?: AccessNode;
}>();

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);
const parentsVal = ref<number[]>([]);

const defaultForm = (): AccessFormModel => ({
  name: '',
  path: '',
  component: '',
  permission_code: '',
  api_path: '',
  icon: '',
  is_menu: 0,
  is_route: 1,
  is_show: 1,
  sort: 1,
  parent_id: 0,
  redirect_name: '',
  remark: '',
});

const formData = reactive<AccessFormModel>(defaultForm());

const formRules = {
  name: [{ message: '请输入菜单名称', required: true, trigger: 'blur' }],
  path: [{ message: '请输入路径', required: true, trigger: 'blur' }],
};

const dialogTitle = computed(() =>
  props.mode === 'edit' ? '修改菜单&权限' : '新增菜单&权限',
);

const accessList = computed(() => {
  const list = deepClone(props.rawData) as AccessNode[];
  list.unshift({ access_id: 0, children: [], is_route: 1, is_show: 1, name: '顶级菜单', path: '', sort: 0 });
  return list;
});

const cascaderProps = {
  checkStrictly: true,
  label: 'name',
  value: 'access_id',
};

function findParentsID(list: AccessNode[], parentId: number): boolean {
  for (const item of list) {
    if (item.access_id === parentId) {
      parentsVal.value.unshift(item.access_id);
      return true;
    }
    if (item.children?.length && findParentsID(item.children, parentId)) {
      parentsVal.value.unshift(item.access_id);
      return true;
    }
  }
  return false;
}

function resetForm() {
  Object.assign(formData, defaultForm());
  parentsVal.value = [];
}

function initForm() {
  resetForm();
  if (props.mode === 'edit' && props.selectModel) {
    Object.assign(formData, deepClone(props.selectModel));
    findParentsID(accessList.value, formData.parent_id ?? 0);
    return;
  }
  if (props.addType === 'copy' && props.selectModel) {
    Object.assign(formData, formatModel(defaultForm(), props.selectModel as Record<string, unknown>));
    findParentsID(accessList.value, formData.parent_id ?? 0);
    return;
  }
  if (props.addType === 'child' && props.selectModel) {
    formData.parent_id = props.selectModel.access_id;
    findParentsID(accessList.value, formData.parent_id);
  }
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      initForm();
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
  () => [props.mode, props.addType, props.selectModel?.access_id] as const,
  () => {
    if (open.value) {
      initForm();
    }
  },
);

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    const params: Record<string, unknown> = { ...formData };
    if (parentsVal.value.length > 0) {
      params.parent_id = parentsVal.value[parentsVal.value.length - 1] ?? 0;
    }
    if (Number(params.is_route) === 0) {
      params.is_menu = 0;
    }
    loading.value = true;
    drawerApi.setState({ confirmLoading: true });
    try {
      const res =
        props.mode === 'edit'
          ? await MerchantAccessApi.editAccess(params, true)
          : await MerchantAccessApi.addAccess(params, true);
      if (res.code === 1) {
        ElMessage.success(res.msg || '操作成功');
        open.value = false;
        emit('success');
      }
    } finally {
      loading.value = false;
      drawerApi.setState({ confirmLoading: false });
    }
  });
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :confirm-loading="loading"
    :destroy-on-close="true"
    :title="dialogTitle"
  >
    <ElForm ref="formRef" :model="formData" :rules="formRules" label-width="120px">
      <ElFormItem label="菜单名称" prop="name">
        <ElInput v-model="formData.name" placeholder="请输入菜单名称" />
      </ElFormItem>
      <ElFormItem label="类型">
        <ElRadioGroup v-model="formData.is_route">
          <ElRadio :label="1">页面</ElRadio>
          <ElRadio :label="0">按钮</ElRadio>
          <ElRadio :label="2">独立单页面</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="上级菜单">
        <ElCascader
          v-model="parentsVal"
          :options="accessList"
          :props="cascaderProps"
          class="w-full"
          clearable
        />
      </ElFormItem>
      <ElFormItem label="路径" prop="path">
        <ElInput v-model="formData.path" placeholder="路由 path 或 #btn:module:action" />
      </ElFormItem>
      <ElFormItem v-if="formData.is_route === 1" label="组件路径">
        <ElInput
          v-model="formData.component"
          placeholder="views/home/index.vue"
        />
      </ElFormItem>
      <ElFormItem label="权限码">
        <ElInput
          v-model="formData.permission_code"
          placeholder="platform:shop:list"
        />
      </ElFormItem>
      <ElFormItem label="接口路径">
        <ElInput
          v-model="formData.api_path"
          placeholder="platform.access/index 或 live:..."
        />
      </ElFormItem>
      <ElFormItem label="图标">
        <IconPicker
          v-model="formData.icon"
          :auto-fetch-api="false"
          :icons="PLATFORM_LUCIDE_ICONS"
          :input-component="ElInput"
          class="w-full"
          icon-slot="append"
          model-value-prop="model-value"
          prefix="lucide"
          type="input"
        />
      </ElFormItem>
      <ElFormItem v-if="formData.is_route === 1" label="是否是菜单">
        <ElRadioGroup v-model="formData.is_menu">
          <ElRadio :label="1">是</ElRadio>
          <ElRadio :label="0">否</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="是否显示">
        <ElRadioGroup v-model="formData.is_show">
          <ElRadio :label="1">是</ElRadio>
          <ElRadio :label="0">否</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem v-if="formData.is_route === 1" label="重定向">
        <ElInput v-model="formData.redirect_name" placeholder="请输入重定向地址" />
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="formData.remark" placeholder="请输入备注" type="textarea" />
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInput v-model.number="formData.sort" placeholder="请输入排序" type="number" />
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
