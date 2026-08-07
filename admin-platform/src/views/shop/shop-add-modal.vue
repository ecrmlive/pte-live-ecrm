<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElCheckbox,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import ShopApi from '#/api/core/shop';
import MerchantRoleApi from '#/api/core/merchant-role';
import type { MerchantTemplateRoleRow } from '#/api/core/merchant-role';
import LiveTrafficApi from '#/api/core/live-traffic';
import { hasSpace, isAllSpace, replaceSpace } from '#/utils/form-text';
import { validateBEndPassword } from '#/utils/b-end-password';

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);
const roleOptions = ref<MerchantTemplateRoleRow[]>([]);

const form = reactive({
  app_name: '',
  expire_time: '',
  no_expire: false,
  user_name: '',
  password: '',
  password_confirm: '',
  weixin_service: false,
  initial_traffic_gb: 0,
  initial_amount_yuan: 0,
  merchant_role_id: undefined as number | undefined,
});

const rules = {
  app_name: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!value) return callback(new Error('请输入商城名称'));
        if (isAllSpace(value)) return callback(new Error('不能全是空格'));
        callback();
      },
    },
  ],
  user_name: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!replaceSpace(value)) return callback(new Error('商家账户名'));
        if (hasSpace(value)) return callback(new Error('不能包含空格'));
        callback();
      },
    },
  ],
  password: [
    {
      required: true,
      trigger: 'change',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        const msg = validateBEndPassword(value);
        if (msg) return callback(new Error(msg));
        callback();
      },
    },
  ],
  password_confirm: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!value) return callback(new Error('请填写确认密码'));
        if (value !== form.password) return callback(new Error('确认密码不一致'));
        callback();
      },
    },
  ],
  merchant_role_id: [
    {
      required: true,
      trigger: 'change',
      validator: (_: unknown, value: number | undefined, callback: (err?: Error) => void) => {
        if (!value) return callback(new Error('请选择商城角色'));
        callback();
      },
    },
  ],
  expire_time: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (form.no_expire) return callback();
        if (!value) return callback(new Error('请输入过期时间'));
        callback();
      },
    },
  ],
};

async function loadRoleOptions() {
  try {
    const res = await MerchantRoleApi.roleList(true);
    roleOptions.value = res.data?.list ?? [];
    if (!form.merchant_role_id && roleOptions.value.length === 1) {
      form.merchant_role_id = roleOptions.value[0]?.role_id;
    }
  } catch {
    roleOptions.value = [];
  }
}

function resetForm() {
  Object.assign(form, {
    app_name: '',
    expire_time: '',
    no_expire: false,
    user_name: '',
    password: '',
    password_confirm: '',
    weixin_service: false,
    initial_traffic_gb: 0,
    initial_amount_yuan: 0,
    merchant_role_id: undefined,
  });
  void loadRoleOptions();
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      resetForm();
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
      drawerApi.setState({ title: '新增小程序商城' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    drawerApi.setState({ confirmLoading: true });
    try {
      const res = await ShopApi.addShop({ ...form }, true);
      if (res.code === 1) {
        const appId = res.data?.app_id ?? res.data?.appId;
        const initialGb = Number(form.initial_traffic_gb || 0);
        const initialAmount = Number(form.initial_amount_yuan || 0);
        const finish = () => {
          ElMessage.success('恭喜你，添加成功');
          open.value = false;
          emit('success');
        };
        if (appId && initialGb > 0) {
          try {
            await LiveTrafficApi.recharge({
              amount_yuan: initialAmount,
              app_id: appId,
              delta_gb: initialGb,
              recharge_type: 'initial',
              remark: '创建商城初始流量',
            });
          } catch {
            /* ignore recharge failure */
          }
        }
        finish();
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
    title="新增小程序商城"
  >
    <div style="height: 0; overflow: hidden">
      <input type="password" />
    </div>
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="132px">
      <ElFormItem label="商城角色" prop="merchant_role_id">
        <ElSelect v-model="form.merchant_role_id" placeholder="请选择商城角色" style="width: 100%">
          <ElOption
            v-for="item in roleOptions"
            :key="item.role_id"
            :label="item.role_name"
            :value="item.role_id"
          />
        </ElSelect>
        <p class="text-xs text-[#909399]">创建的管理员仅拥有该角色配置的商户功能</p>
      </ElFormItem>
      <ElFormItem label="商城名称" prop="app_name">
        <ElInput v-model="form.app_name" autocomplete="off" placeholder="请输入商城名称" />
      </ElFormItem>
      <ElFormItem label="过期时间" prop="expire_time">
        <ElDatePicker
          v-model="form.expire_time"
          :disabled="form.no_expire"
          placeholder="过期时间"
          type="date"
          value-format="YYYY-MM-DD"
        />
        <ElCheckbox v-model="form.no_expire" class="pl-4">永不过期</ElCheckbox>
      </ElFormItem>
      <ElFormItem label="商家账户名" prop="user_name">
        <ElInput v-model="form.user_name" autocomplete="off" placeholder="请输入商家账户名" />
        <p class="text-xs text-[#909399]">注：商家后台用户名</p>
      </ElFormItem>
      <ElFormItem label="商家账户密码" prop="password">
        <ElInput
          v-model="form.password"
          autocomplete="off"
          placeholder="请输入密码"
          type="password"
        />
        <p class="text-xs text-[#909399]">注：商家后台用户密码</p>
      </ElFormItem>
      <ElFormItem label="确认密码" prop="password_confirm">
        <ElInput
          v-model="form.password_confirm"
          autocomplete="off"
          placeholder="请输入确认密码"
          type="password"
        />
      </ElFormItem>
      <ElFormItem label="微信服务商支付">
        <ElCheckbox v-model="form.weixin_service">开启</ElCheckbox>
      </ElFormItem>
      <ElFormItem label="初始流量(GB)">
        <ElInputNumber v-model="form.initial_traffic_gb" :min="0" :step="100" />
      </ElFormItem>
      <ElFormItem label="初始充值金额(元)">
        <ElInputNumber v-model="form.initial_amount_yuan" :min="0" :precision="2" />
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
