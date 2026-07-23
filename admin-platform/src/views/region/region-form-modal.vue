<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue';

import type { FormRules } from 'element-plus';
import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElRadio,
  ElSelect,
} from 'element-plus';

import RegionApi from '#/api/core/region';

import {
  regionCityOptions,
  type RegionAreaItem,
  type RegionFormModel,
  type RegionRow,
} from './types';

const props = defineProps<{
  mode: 'add' | 'edit';
  open: boolean;
  region?: RegionRow;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const loading = ref(false);
const formRef = ref<InstanceType<typeof ElForm>>();
const areaList = ref<Record<number | string, RegionAreaItem>>({});
const provinceOptions = ref<RegionAreaItem[]>([]);

function createEmptyForm(): RegionFormModel {
  return {
    level: 1,
    name: '',
    shortname: '',
    merger_name: '',
    pinyin: '',
    code: '',
    ad_code: '',
    zip_code: '',
    first: '',
    lng: '',
    lat: '',
    sort: '',
  };
}

const form = reactive<RegionFormModel>(createEmptyForm());

const dialogTitle = computed(() =>
  props.mode === 'edit' ? '编辑地区' : '添加地区',
);

const requiredBlur = {
  message: '必填',
  required: true,
  trigger: 'blur',
} as const;

const formRules = computed<FormRules>(() => {
  const rules: FormRules = {
    name: [requiredBlur],
    shortname: [requiredBlur],
    merger_name: [requiredBlur],
    pinyin: [requiredBlur],
    ad_code: [requiredBlur],
    zip_code: [requiredBlur],
    first: [requiredBlur],
    lng: [requiredBlur],
    lat: [requiredBlur],
    sort: [requiredBlur],
  };
  if (form.level > 1) {
    rules.province_id = [
      {
        message: '请选择省份',
        required: true,
        trigger: 'change',
        validator: (_rule, value, callback) => {
          if (value !== undefined && value !== null && Number(value) > 0) {
            callback();
          } else {
            callback(new Error('请选择省份'));
          }
        },
      },
    ];
  }
  if (form.level > 2) {
    rules.city_id = [
      {
        message: '请选择城市',
        required: true,
        trigger: 'change',
        validator: (_rule, value, callback) => {
          if (value !== undefined && value !== null && Number(value) > 0) {
            callback();
          } else {
            callback(new Error('请选择城市'));
          }
        },
      },
    ];
  }
  return rules;
});

function cityOptions() {
  return regionCityOptions(areaList.value, form.province_id ?? 0);
}

function parseSelectId(value: unknown): number | undefined {
  if (value === undefined || value === null || value === '') {
    return undefined;
  }
  const n = Number(value);
  return Number.isFinite(n) && n > 0 ? n : undefined;
}

function normalizeFormModel(model: Record<string, unknown>) {
  Object.assign(form, createEmptyForm(), model);
  form.province_id = parseSelectId(model.province_id);
  form.city_id = parseSelectId(model.city_id);
  if (model.id !== undefined && model.id !== null && model.id !== '') {
    form.id = Number(model.id);
  }
  if (model.sort !== undefined && model.sort !== null) {
    form.sort = model.sort as string | number;
  }
}

function initCity() {
  form.city_id = undefined;
}

async function loadFormData() {
  loading.value = true;
  try {
    if (props.mode === 'edit' && props.region) {
      const res = await RegionApi.regionDetail({ id: props.region.id }, true);
      const model =
        (res.data as { model?: Record<string, unknown> })?.model ?? {};
      normalizeFormModel(model);
      areaList.value =
        (res.data as { regionData?: Record<number, RegionAreaItem> })?.regionData ??
        {};
    } else {
      Object.assign(form, createEmptyForm());
      const res = await RegionApi.toAddRegion({}, true);
      areaList.value =
        (res.data as { regionData?: Record<number, RegionAreaItem> })?.regionData ??
        {};
    }
    provinceOptions.value = Object.values(areaList.value).filter(
      (item) => Number(item.id) > 0,
    );
    formRef.value?.clearValidate();
  } finally {
    loading.value = false;
  }
}

watch(
  () => form.level,
  (level) => {
    if (level <= 1) {
      form.province_id = undefined;
    }
    if (level <= 2) {
      form.city_id = undefined;
    }
  },
);

watch(
  () => [props.open, props.mode, props.region?.id] as const,
  ([open]) => {
    if (open) {
      loadFormData();
    }
  },
);

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    try {
      const payload = {
        ...form,
        province_id: form.province_id ?? '',
        city_id: form.city_id ?? '',
      };
      if (props.mode === 'edit') {
        const id = props.region?.id ?? form.id;
        if (!id) {
          ElMessage.error('缺少地区 ID');
          return;
        }
        const res = await RegionApi.editRegion({ ...payload, id }, true);
        if (res.code === 1) {
          ElMessage.success(res.msg || '修改成功');
          emit('update:open', false);
          emit('success');
        }
      } else {
        const res = await RegionApi.addRegion(payload, true);
        if (res.code === 1) {
          ElMessage.success(res.msg || '添加成功');
          emit('update:open', false);
          emit('success');
        }
      }
    } finally {
      loading.value = false;
    }
  });
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :model-value="open"
    :title="dialogTitle"
    width="560px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <ElForm
      ref="formRef"
      v-loading="loading"
      :model="form"
      :rules="formRules"
      label-width="100px"
    >
      <ElFormItem v-if="mode === 'edit'" label="ID">
        <ElInput :model-value="String(form.id ?? '')" disabled />
      </ElFormItem>
      <ElFormItem v-else label="ID">
        <ElInput disabled placeholder="新增后自动生成" />
      </ElFormItem>

      <ElFormItem label="地区类型">
        <ElRadio v-model="form.level" :label="1">省份</ElRadio>
        <ElRadio v-model="form.level" :label="2">城市</ElRadio>
        <ElRadio v-model="form.level" :label="3">地区</ElRadio>
      </ElFormItem>
      <ElFormItem
        v-if="form.level > 1"
        label="选择省份"
        prop="province_id"
      >
        <ElSelect
          :key="`province-${provinceOptions.length}-${form.province_id ?? ''}`"
          v-model="form.province_id"
          class="w-full"
          placeholder="请选择省份"
          @change="initCity"
        >
          <ElOption
            v-for="item in provinceOptions"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem
        v-if="form.level > 2 && form.province_id"
        label="选择城市"
        prop="city_id"
      >
        <ElSelect
          :key="`city-${form.province_id ?? ''}-${cityOptions().length}-${form.city_id ?? ''}`"
          v-model="form.city_id"
          class="w-full"
          placeholder="请选择城市"
        >
          <ElOption
            v-for="item in cityOptions()"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="地区名称" prop="name">
        <ElInput v-model="form.name" placeholder="请输入地区名称" />
      </ElFormItem>
      <ElFormItem label="简称" prop="shortname">
        <ElInput v-model="form.shortname" placeholder="简称" />
      </ElFormItem>
      <ElFormItem label="全称" prop="merger_name">
        <ElInput v-model="form.merger_name" placeholder="全称" />
      </ElFormItem>
      <ElFormItem label="拼音" prop="pinyin">
        <ElInput v-model="form.pinyin" placeholder="拼音" />
      </ElFormItem>
      <ElFormItem label="区划编号" prop="ad_code">
        <ElInput v-model="form.ad_code" placeholder="如 110101" />
      </ElFormItem>
      <ElFormItem label="长途区号" prop="code">
        <ElInput v-model="form.code" placeholder="如 010" />
      </ElFormItem>
      <ElFormItem label="邮编" prop="zip_code">
        <ElInput v-model="form.zip_code" placeholder="邮编" />
      </ElFormItem>
      <ElFormItem label="首字母" prop="first">
        <ElInput v-model="form.first" placeholder="首字母" />
      </ElFormItem>
      <ElFormItem label="经度" prop="lng">
        <ElInput v-model="form.lng" placeholder="经度" />
      </ElFormItem>
      <ElFormItem label="纬度" prop="lat">
        <ElInput v-model="form.lat" placeholder="纬度" />
      </ElFormItem>
      <ElFormItem label="排序" prop="sort">
        <ElInput v-model="form.sort" placeholder="100" type="number" />
        <div class="text-xs text-[#999]">数字越小越靠前</div>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">
        保存
      </ElButton>
    </template>
  </ElDialog>
</template>
