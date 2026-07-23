<script setup lang="ts">
import type { LiveSessionChatItem } from '#/api/core/live';

import { ref, watch } from 'vue';

import { getLiveSessionChatListApi } from '#/api/core/live';
import { getDanmakuRoleColor } from '#/utils/live/liveDanmakuStyle.js';

defineOptions({ name: 'LiveSessionChatPanel' });

const props = defineProps<{
  embedded?: boolean;
  liveId?: number;
  sessionId?: string;
}>();

const loading = ref(false);
const list = ref<LiveSessionChatItem[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);

function roleColor(item: LiveSessionChatItem) {
  return getDanmakuRoleColor(item?.role, item?.source);
}

function roleStyle(item: LiveSessionChatItem) {
  const color = roleColor(item);
  return {
    background: `${color}2E`,
    border: `1px solid ${color}55`,
    color,
  };
}

async function loadList() {
  const liveId = Number(props.liveId ?? 0);
  const sessionId = String(props.sessionId ?? '').trim();
  if (!liveId || !sessionId) {
    list.value = [];
    total.value = 0;
    return;
  }
  loading.value = true;
  try {
    const res = await getLiveSessionChatListApi({
      live_id: liveId,
      list_rows: pageSize.value,
      page: page.value,
      session_id: sessionId,
    });
    list.value = res.list ?? [];
    total.value = Number(res.total ?? 0);
  } catch {
    list.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  page.value = 1;
  void loadList();
}

function onPageChange(val: number) {
  page.value = val;
  void loadList();
}

watch(
  () => [props.liveId, props.sessionId] as const,
  () => reload(),
  { immediate: true },
);
</script>

<template>
  <div class="session-chat-panel" :class="{ 'session-chat-panel--embedded': embedded }">
    <div v-if="!embedded" class="stat-panel__head">
      <span class="stat-panel__bar" />
      <span class="stat-panel__title">聊天消息</span>
      <span v-if="total > 0" class="session-chat-panel__count">共 {{ total }} 条</span>
    </div>
    <div v-else-if="total > 0" class="session-chat-panel__toolbar">
      <span class="session-chat-panel__count">共 {{ total }} 条</span>
    </div>

    <div v-loading="loading" class="session-chat-panel__body">
      <el-empty v-if="!loading && !list.length" :image-size="72" description="本场暂无聊天消息" />
      <div v-for="item in list" :key="item.message_id" class="session-chat-row">
        <span class="session-chat-row__role" :style="roleStyle(item)">{{ item.role_text }}</span>
        <span class="session-chat-row__uid" :style="{ color: roleColor(item) }">{{ item.user_id }}</span>
        <span class="session-chat-row__nick" :style="{ color: roleColor(item) }">
          {{ item.nick_name || '用户' }}
        </span>
        <span class="session-chat-row__text">{{ item.content }}</span>
        <span class="session-chat-row__time">{{ item.send_time_text }}</span>
      </div>
    </div>

    <div v-if="total > 0" class="session-chat-panel__pager">
      <el-pagination
        background
        :current-page="page"
        layout="total, prev, pager, next"
        :page-size="pageSize"
        small
        :total="total"
        @current-change="onPageChange"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.session-chat-panel {
  margin-top: 8px;
}

.session-chat-panel--embedded {
  margin-top: 0;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.session-chat-panel__toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.session-chat-panel--embedded .session-chat-panel__body {
  max-height: min(56vh, 520px);
  flex: 1;
  min-height: 200px;
}

.session-chat-panel__count {
  margin-left: auto;
  font-size: 12px;
  color: #909399;
}

.session-chat-panel__body {
  min-height: 120px;
  max-height: 360px;
  overflow-y: auto;
  padding: 8px 4px 4px;
}

.session-chat-row {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px 6px;
  padding: 8px 10px;
  margin-bottom: 6px;
  border-radius: 8px;
  background: #f7f8fa;
  font-size: 13px;
  line-height: 1.5;
}

.session-chat-row__role {
  display: inline-block;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  line-height: 18px;
}

.session-chat-row__uid,
.session-chat-row__nick {
  font-weight: 600;
}

.session-chat-row__text {
  flex: 1 1 100%;
  color: rgba(48, 49, 51, 0.92);
  word-break: break-word;
}

.session-chat-row__time {
  margin-left: auto;
  font-size: 11px;
  color: #909399;
}

.session-chat-panel__pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
}
</style>
