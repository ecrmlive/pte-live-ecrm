<template>
  <div>
    <a-card :bordered="false" class="card" title="待发货订单">
      <a-list :loading="loading" :data-source="list" item-layout="vertical">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta
              :title="item.order_sn"
              :description="`${item.real_name} ${item.user_phone} · ¥${Number(item.pay_price).toFixed(2)}`"
            />
            <div class="addr">{{ item.user_address }}</div>
            <template v-if="canDeliver" #actions>
              <a @click="openDeliver(item)">发货</a>
            </template>
          </a-list-item>
        </template>
      </a-list>
      <a-empty v-if="!loading && !list.length" description="暂无待发货" />
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      title="确认发货"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="快递公司" required>
          <a-input v-model:value="form.delivery_name" placeholder="如：顺丰速运" />
        </a-form-item>
        <a-form-item label="运单号" required>
          <a-input v-model:value="form.delivery_id" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { deliverOrder, fetchAwaitShip, type StoreOrder } from '@/api/order';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canDeliver = computed(() => auth.user?.is_goods === 1);

const loading = ref(false);
const saving = ref(false);
const list = ref<StoreOrder[]>([]);
const modalOpen = ref(false);
const currentId = ref(0);
const form = reactive({ delivery_name: '演示快递', delivery_id: '' });

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAwaitShip();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

function openDeliver(row: StoreOrder) {
  currentId.value = row.order_id;
  form.delivery_name = '演示快递';
  form.delivery_id = `SF${Date.now().toString().slice(-10)}`;
  modalOpen.value = true;
}

async function submit() {
  if (!form.delivery_name || !form.delivery_id) {
    message.warning('请填写物流信息');
    return;
  }
  saving.value = true;
  try {
    await deliverOrder(currentId.value, {
      delivery_name: form.delivery_name,
      delivery_id: form.delivery_id,
      delivery_type: '1',
    });
    message.success('已发货');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
.card {
  margin-bottom: 12px;
  border-radius: 12px;
}
.addr {
  color: #888;
  font-size: 13px;
  margin-top: 4px;
}
</style>
