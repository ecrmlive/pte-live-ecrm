<script setup lang="ts">
/**
 * 左侧类型 Select + 右侧关键词 Input（Element Plus Input prepend）。
 * 用户搜索标准选项见 `#/components/ecrm/user-search-field`。
 */
import { ElInput, ElOption, ElSelect } from 'element-plus';

export type PrefixedKeywordValue = {
  keyword: string;
  type: string;
};

export type PrefixedKeywordOption = {
  label: string;
  value: string;
};

const props = withDefaults(
  defineProps<{
    modelValue?: PrefixedKeywordValue;
    options?: PrefixedKeywordOption[];
    placeholder?: string;
    typeWidth?: string;
  }>(),
  {
    modelValue: () => ({ type: 'nickname', keyword: '' }),
    options: () => [],
    placeholder: '请输入内容',
    typeWidth: '96px',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: PrefixedKeywordValue];
}>();

function currentType() {
  const type = props.modelValue?.type;
  if (type) return type;
  return props.options[0]?.value || '';
}

function currentKeyword() {
  return props.modelValue?.keyword ?? '';
}

function patch(partial: Partial<PrefixedKeywordValue>) {
  emit('update:modelValue', {
    type: currentType(),
    keyword: currentKeyword(),
    ...partial,
  });
}
</script>

<template>
  <ElInput
    class="prefixed-keyword-field"
    :model-value="currentKeyword()"
    clearable
    :placeholder="placeholder"
    @update:model-value="(v) => patch({ keyword: String(v ?? '') })"
  >
    <template #prepend>
      <ElSelect
        class="prefixed-keyword-field__type"
        :model-value="currentType()"
        :style="{ width: typeWidth }"
        @update:model-value="(v) => patch({ type: String(v ?? '') })"
      >
        <ElOption
          v-for="opt in options"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </ElSelect>
    </template>
  </ElInput>
</template>

<style scoped>
.prefixed-keyword-field {
  width: 100%;
}

.prefixed-keyword-field :deep(.el-input-group__prepend) {
  padding: 0;
  background: var(--el-fill-color-blank);
}

.prefixed-keyword-field__type {
  margin: -1px 0;
}

.prefixed-keyword-field__type :deep(.el-select__wrapper) {
  box-shadow: none !important;
  background: transparent;
}
</style>
