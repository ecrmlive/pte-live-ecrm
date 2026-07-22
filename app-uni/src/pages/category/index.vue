<template>
  <view class="page">
    <view class="aside">
      <view
        v-for="c in roots"
        :key="c.id"
        class="aside-item"
        :class="{ active: c.id === activeId }"
        @click="activeId = c.id"
      >
        <text>{{ c.name }}</text>
      </view>
    </view>
    <scroll-view scroll-y class="main">
      <view class="main-head">{{ activeName }}</view>
      <view class="chips">
        <view v-for="c in children" :key="c.id" class="chip" @click="goList(c.id)">
          {{ c.name }}
        </view>
        <view v-if="!children.length && activeId" class="chip" @click="goList(activeId)">全部商品</view>
        <view v-if="!roots.length" class="empty">暂无分类</view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchCategories, type CategoryItem } from "@/api/catalog";

const list = ref<CategoryItem[]>([]);
const activeId = ref(0);

const roots = computed(() => list.value.filter((c) => c.pid === 0));
const children = computed(() => list.value.filter((c) => c.pid === activeId.value));
const activeName = computed(
  () => roots.value.find((c) => c.id === activeId.value)?.name || "分类"
);

onShow(async () => {
  try {
    list.value = await fetchCategories();
    if (!activeId.value && roots.value.length) {
      activeId.value = roots.value[0].id;
    }
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goList(cateId: number) {
  uni.navigateTo({ url: `/pages/goods/list?cate_id=${cateId}` });
}
</script>

<style lang="scss" scoped>
.page {
  display: flex;
  height: 100vh;
  background: var(--qx-bg);
}
.aside {
  width: 200rpx;
  background: #f7f7f7;
}
.aside-item {
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  color: var(--qx-text-secondary);
}
.aside-item.active {
  background: #fff;
  color: var(--qx-brand);
  font-weight: 600;
}
.main {
  flex: 1;
  background: #fff;
  padding: 24rpx;
  box-sizing: border-box;
}
.main-head {
  font-size: 30rpx;
  font-weight: 700;
  margin-bottom: 20rpx;
}
.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}
.chip {
  padding: 16rpx 28rpx;
  background: #f5f5f5;
  border-radius: 12rpx;
  font-size: 26rpx;
}
.empty {
  color: var(--qx-text-secondary);
  font-size: 26rpx;
  padding: 40rpx 0;
}
</style>
