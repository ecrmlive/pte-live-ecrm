<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchSeckillList, fetchSeckillTimes, type SeckillActive, type SeckillTimeSlot } from "@/api/seckill";

const router = useRouter();
const loading = ref(false);
const list = ref<SeckillActive[]>([]);
const slots = ref<SeckillTimeSlot[]>([]);
const selectedSlot = ref("");

const visibleList = computed(() => selectedSlot.value
  ? list.value.filter((item) => item.time_slots?.includes(selectedSlot.value))
  : list.value,
);
const activeSlot = computed(() => slots.value.find((slot) => slot.label === selectedSlot.value));

onMounted(async () => {
  loading.value = true;
  try {
    const [activities, timeSlots] = await Promise.all([fetchSeckillList(), fetchSeckillTimes()]);
    list.value = activities.list || [];
    slots.value = timeSlots.list || [];
    selectedSlot.value = slots.value.find((slot) => slot.active)?.label || slots.value[0]?.label || "";
  } finally {
    loading.value = false;
  }
});

function go(id: number) {
  void router.push(`/goods/${id}`);
}
</script>

<template>
  <div class="seckill-page">
    <section class="seckill-hero">
      <div class="pc-container seckill-hero__inner">
        <p>SECKILL</p>
        <h1>限时秒杀</h1>
        <span>精选好物，限时抢购</span>
      </div>
    </section>

    <div class="pc-container content">
      <section class="time-panel" aria-label="秒杀场次">
        <button
          v-for="slot in slots"
          :key="slot.label"
          type="button"
          :class="{ active: selectedSlot === slot.label }"
          @click="selectedSlot = slot.label"
        >
          <strong>{{ slot.label }}</strong>
          <span>{{ slot.active ? "抢购中" : "即将开始" }}</span>
        </button>
      </section>

      <section class="activity-head">
        <div><h2>{{ selectedSlot || "全部场次" }}</h2><p>{{ activeSlot?.active ? "本场正在进行，优惠价以结算页为准" : "活动场次即将开始" }}</p></div>
        <span>共 {{ visibleList.length }} 件秒杀商品</span>
      </section>

      <p v-if="loading" class="hint">秒杀活动加载中…</p>
      <p v-else-if="!visibleList.length" class="hint empty">当前场次暂无秒杀商品</p>
      <section v-else class="seckill-grid">
        <article v-for="item in visibleList" :key="item.seckill_active_id" class="seckill-card" @click="go(item.product_id)">
          <div class="cover"><img v-if="item.image" :src="item.image" :alt="item.store_name || item.name" /><span v-else>暂无商品图片</span><b>{{ item.in_window ? "抢购中" : "未开场" }}</b></div>
          <div class="meta">
            <p class="shop">{{ item.shop_name || item.mer_name || "七禧商城" }}</p>
            <h3>{{ item.store_name || item.name }}</h3>
            <div class="price"><strong>¥{{ Number(item.seckill_price).toFixed(2) }}</strong><del v-if="item.price">¥{{ Number(item.price).toFixed(2) }}</del></div>
            <div class="buy-line"><span>已售 {{ item.sales || 0 }}</span><button type="button" @click.stop="go(item.product_id)">查看商品</button></div>
          </div>
        </article>
      </section>
    </div>
  </div>
</template>

<style scoped>
.seckill-page { min-height: 680px; padding-bottom: 54px; background: #f5f5f5; }.seckill-hero { min-height: 214px; color: #fff; background: #e84239; }.seckill-hero__inner { display: flex; flex-direction: column; justify-content: center; align-items: center; min-height: 214px; }.seckill-hero p { margin: 0; letter-spacing: .36em; font-size: 12px; opacity: .78; }.seckill-hero h1 { margin: 8px 0 4px; font-size: 40px; letter-spacing: .12em; }.seckill-hero span { font-size: 14px; opacity: .86; }.content { margin-top: -32px; }.time-panel { display: grid; grid-template-columns: repeat(4, 1fr); overflow: hidden; border: 1px solid #e7e7e7; background: #fff; box-shadow: 0 6px 18px rgb(0 0 0 / 5%); }.time-panel button { display: grid; justify-items: center; gap: 4px; min-height: 82px; border: 0; border-right: 1px solid #eee; color: #555; background: #fff; }.time-panel button:last-child { border-right: 0; }.time-panel strong { font-size: 24px; }.time-panel span { color: #999; font-size: 12px; }.time-panel button.active { color: #fff; background: #ef3727; }.time-panel button.active span { color: #fff; opacity: .88; }.activity-head { display: flex; align-items: end; justify-content: space-between; margin: 24px 0 14px; }.activity-head h2 { margin: 0; color: #333; font-size: 20px; }.activity-head p, .activity-head > span { margin: 6px 0 0; color: #999; font-size: 13px; }.seckill-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 15px; }.seckill-card { overflow: hidden; background: #fff; cursor: pointer; }.cover { position: relative; aspect-ratio: 1; display: grid; place-items: center; color: #aaa; background: #f1f1f1; }.cover img { width: 100%; height: 100%; object-fit: cover; }.cover b { position: absolute; top: 10px; left: 10px; padding: 4px 7px; color: #fff; background: #ef3727; font-size: 12px; }.meta { padding: 13px 14px 15px; }.shop { overflow: hidden; margin: 0 0 5px; color: #999; font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }.meta h3 { display: -webkit-box; min-height: 42px; margin: 0; overflow: hidden; color: #333; font-size: 15px; line-height: 21px; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }.price { display: flex; align-items: baseline; gap: 8px; margin-top: 12px; }.price strong { color: #ef3727; font-size: 22px; }.price del { color: #aaa; font-size: 12px; }.buy-line { display: flex; align-items: center; justify-content: space-between; margin-top: 14px; color: #999; font-size: 12px; }.buy-line button { border: 0; padding: 6px 10px; color: #fff; background: #ef3727; font-size: 12px; }.hint { margin: 50px 0; color: #999; text-align: center; }.hint.empty { padding: 60px 0; background: #fff; }
@media (max-width: 800px) { .seckill-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.seckill-hero h1 { font-size: 30px; }.time-panel strong { font-size: 19px; } }
</style>
