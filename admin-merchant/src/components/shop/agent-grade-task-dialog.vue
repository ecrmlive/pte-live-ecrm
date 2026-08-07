<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { AgentTaskItem, AgentTaskTypeOption } from '#/api/core/plus-agent';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@element-plus/icons-vue';
import { ElButton, ElMessage, ElMessageBox } from 'element-plus';
import { computed, reactive, ref, shallowRef, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  addAgentTaskApi,
  deleteAgentTaskApi,
  editAgentTaskApi,
  getAgentTaskListApi,
  setAgentTaskStateApi,
} from '#/api/core/plus-agent';

defineOptions({ name: 'AgentGradeTaskDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  gradeId?: number;
}>();

const typeList = shallowRef<AgentTaskTypeOption[]>([]);
const formOpen = ref(false);
const formMode = ref<'add' | 'edit'>('add');
const submitting = ref(false);
const editingTaskId = ref(0);

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Select',
    componentProps: {
      class: 'w-full',
      options: typeList.value.map((item) => ({
        label: item.name,
        value: item.type_id,
      })),
      placeholder: '-请选择任务类型-',
    },
    fieldName: 'task_type',
    label: '任务类型',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    componentProps: { placeholder: '请输入任务名称' },
    fieldName: 'name',
    label: '任务名称',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { placeholder: '请输入任务数量', type: 'number' },
    fieldName: 'number',
    label: '任务数量',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { placeholder: '请输入排序', type: 'number' },
    fieldName: 'sort',
    label: '排序',
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '显示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
  },
]);

const [TaskForm, taskFormApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { size: 'small' },
      labelWidth: 120,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        const payload = {
          grade_id: props.gradeId ?? 0,
          name: String(values.name ?? ''),
          number: String(values.number ?? ''),
          sort: String(values.sort ?? ''),
          status: Number(values.status ?? 1),
          task_id: editingTaskId.value,
          task_type: values.task_type as number | string,
        };
        if (formMode.value === 'add') {
          await addAgentTaskApi(payload);
        } else {
          await editAgentTaskApi(payload);
        }
        ElMessage.success('保存成功');
        formOpen.value = false;
        await gridApi.reload();
      } finally {
        submitting.value = false;
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const gridOptions = reactive<VxeGridProps<AgentTaskItem>>({
  columns: [
    { field: 'task_id', title: '任务ID', width: 90 },
    { field: 'name', minWidth: 140, showOverflow: true, title: '任务名称' },
    { field: 'type_name', minWidth: 120, title: '任务类型' },
    { field: 'number', minWidth: 100, title: '任务数量' },
    {
      field: 'status',
      minWidth: 100,
      slots: { default: 'status' },
      title: '状态',
    },
    { field: 'sort', minWidth: 80, title: '排序' },
    {
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 120,
    },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.gradeId) {
          return { items: [], total: 0 };
        }
        const res = await getAgentTaskListApi({
          grade_id: props.gradeId,
          list_rows: page.pageSize,
          page: page.currentPage,
        });
        typeList.value = res.typeList ?? [];
        return {
          items: res.list.data,
          total: res.list.total,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'task_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

const [FormModal, formModalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    formOpen.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    void gridApi.reload();
    return;
  }
  modalApi.close();
});

watch(formOpen, (visible) => {
  if (visible) formModalApi.open();
  else formModalApi.close();
});

function openAddForm() {
  formMode.value = 'add';
  editingTaskId.value = 0;
  void taskFormApi.resetForm();
  void taskFormApi.setValues({
    name: '',
    number: '',
    sort: '',
    status: 1,
    task_type: '',
  });
  formOpen.value = true;
}

function openEditForm(row: AgentTaskItem) {
  formMode.value = 'edit';
  editingTaskId.value = row.task_id;
  void taskFormApi.setValues({ ...row });
  formOpen.value = true;
}

async function submitForm() {
  await taskFormApi.validateAndSubmitForm();
}

async function toggleStatus(row: AgentTaskItem, status: boolean) {
  await setAgentTaskStateApi({
    status: status ? 1 : 0,
    task_id: row.task_id,
  });
  ElMessage.success('操作成功');
  await gridApi.reload();
}

async function deleteTask(row: AgentTaskItem) {
  await ElMessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
    type: 'warning',
  });
  await deleteAgentTaskApi(row.task_id);
  ElMessage.success('删除成功');
  await gridApi.reload();
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[1000px]"
    title="等级任务"
  >
    <div class="mb-4">
      <el-button
        v-auth="'/plus/agent/task/add'"
        :icon="Plus"
        size="small"
        type="primary"
        @click="openAddForm"
      >
        添加任务
      </el-button>
    </div>
    <Grid>
      <template #status="{ row }">
        <el-switch
          :model-value="row.status === 1"
          @change="(val: boolean) => toggleStatus(row, val)"
        />
      </template>
      <template #action="{ row }">
        <ElButton
          v-auth="'/plus/agent/task/edit'"
          link
          size="small"
          type="primary"
          @click="openEditForm(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-auth="'/plus/agent/task/delete'"
          link
          size="small"
          type="danger"
          @click="deleteTask(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>
  </Modal>

  <FormModal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[600px]"
    :title="formMode === 'add' ? '添加任务' : '编辑任务'"
  >
    <TaskForm />
    <template #footer>
      <ElButton @click="formOpen = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submitForm">确定</ElButton>
    </template>
  </FormModal>
</template>
