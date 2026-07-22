<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  assistCreate,
  fetchAssist,
  fetchSets,
  helpAssist,
  startAssist,
  type AssistSet,
  type ProductAssist,
} from "@/api/assist";
import { fetchAddresses } from "@/api/trade";

const route = useRoute();
const router = useRouter();
const active = ref<ProductAssist | null>(null);
const sets = ref<AssistSet[]>([]);
const msg = ref("");

async function reload() {
  const id = Number(route.params.id || 0);
  if (!id) return;
  active.value = await fetchAssist(id);
  const data = await fetchSets(id);
  sets.value = data.list || [];
}

onMounted(async () => {
  try {
    await reload();
  } catch (e) {
    msg.value = (e as Error).message || "加载失败";
  }
});

async function pickAddressID(): Promise<number> {
  const data = await fetchAddresses();
  const arr = data.list || [];
  if (!arr.length) {
    msg.value = "请先添加收货地址";
    router.push("/user");
    return 0;
  }
  const def = arr.find((a) => a.is_default === 1) || arr[0];
  return def.address_id;
}

async function start() {
  const id = Number(route.params.id || 0);
  try {
    const set = await startAssist(id);
    if (set.status === 10) {
      await order(set.product_assist_set_id);
      return;
    }
    msg.value = "已发起助力，邀请好友帮忙";
    await reload();
  } catch (e) {
    msg.value = (e as Error).message || "发起失败";
  }
}

async function help(setID: number) {
  try {
    await helpAssist(setID);
    msg.value = "助力成功";
    await reload();
  } catch (e) {
    msg.value = (e as Error).message || "助力失败";
  }
}

async function order(setID: number) {
  const addressID = await pickAddressID();
  if (!addressID) return;
  try {
    const g = await assistCreate({
      product_assist_set_id: setID,
      cart_num: 1,
      address_id: addressID,
    });
    router.push(`/pay/${g.group_order_id}`);
  } catch (e) {
    msg.value = (e as Error).message || "下单失败";
  }
}

function back() {
  router.push("/assist");
}
</script>

<template>
  <div class="page">
    <button class="back" type="button" @click="back">← 返回</button>
    <header v-if="active" class="head">
      <h1>{{ active.store_name }}</h1>
      <p>需{{ active.assist_count }}人 · 助力价 ¥{{ active.assist_price }}</p>
    </header>
    <p v-if="msg" class="hint">{{ msg }}</p>
    <button class="cta" type="button" @click="start">发起助力 / 下单</button>
    <h2>进行中</h2>
    <ul>
      <li v-for="s in sets" :key="s.product_assist_set_id">
        <span
          >{{ s.nickname || "用户" }} · {{ s.yet_assist_count }}/{{ s.assist_count }} · status={{
            s.status
          }}</span
        >
        <button v-if="s.status === 1" type="button" class="link" @click="help(s.product_assist_set_id)">
          帮TA
        </button>
        <button v-if="s.status === 10" type="button" class="link" @click="order(s.product_assist_set_id)">
          去下单
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.page {
  max-width: 720px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.back {
  border: 0;
  background: transparent;
  color: #666;
  cursor: pointer;
  margin-bottom: 16px;
}
.head h1 {
  margin: 0 0 8px;
}
.head p,
.hint {
  color: #888;
}
.cta {
  margin: 16px 0;
  background: #e23030;
  color: #fff;
  border: 0;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
}
ul {
  padding-left: 0;
  list-style: none;
}
li {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}
.link {
  border: 0;
  background: transparent;
  color: #e23030;
  cursor: pointer;
}
</style>
