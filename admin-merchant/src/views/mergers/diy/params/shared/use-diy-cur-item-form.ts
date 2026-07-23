import type { VbenFormSchema } from '#/adapter/form';

import type { ComputedRef, MaybeRefOrGetter } from 'vue';

import { computed, reactive, toValue, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';

import {
  extractFieldNames,
  getByPath,
  pickPathsNested,
  setByPath,
  type FormRecord,
} from './path-utils';

export interface UseDiyCurItemFormOptions {
  fieldPaths?: string[];
  labelWidth?: number;
  onInit?: (curItem: FormRecord) => void;
  onValuesChange?: (values: FormRecord, fieldsChanged: string[]) => void;
}

export function useDiyCurItemForm(
  curItem: MaybeRefOrGetter<FormRecord>,
  schema: ComputedRef<VbenFormSchema[]>,
  options: UseDiyCurItemFormOptions = {},
) {
  const fieldPaths = computed(() => {
    if (options.fieldPaths?.length) {
      return options.fieldPaths;
    }
    return extractFieldNames(schema.value);
  });

  const [Form, formApi] = useVbenForm(
    reactive({
      commonConfig: {
        componentProps: {
          size: 'small',
        },
        labelWidth: options.labelWidth ?? 100,
      },
      handleValuesChange(values: FormRecord, fieldsChanged: string[]) {
        const item = toValue(curItem);
        if (options.onValuesChange) {
          options.onValuesChange(values, fieldsChanged);
          return;
        }
        for (const path of fieldsChanged) {
          setByPath(item, path, getByPath(values, path));
        }
      },
      layout: 'horizontal',
      schema,
      showDefaultActions: false,
    }),
  );

  async function syncFromCurItem() {
    const item = toValue(curItem);
    await formApi.setValues(pickPathsNested(item, fieldPaths.value));
  }

  watch(
    () => toValue(curItem),
    (item) => {
      if (item) {
        options.onInit?.(item);
        void syncFromCurItem();
      }
    },
    { immediate: true },
  );

  return { Form, fieldPaths, formApi, syncFromCurItem };
}
