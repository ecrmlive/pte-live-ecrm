<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { PteChatListener, PteChatMessage, PteSendAck } from "@pte-live/im-web-sdk";
import {
  fetchCustomerServiceThreads,
  openCustomerServiceThread,
  type CustomerServiceThread,
} from "@/api/customer-service";
import { connectCustomerServiceIM, type CustomerServiceIMSession } from "@/services/customer-service-im";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const router = useRouter();
const user = useUserStore();

const threads = ref<CustomerServiceThread[]>([]);
const selectedThread = ref<CustomerServiceThread>();
const messages = ref<PteChatMessage[]>([]);
const draft = ref("");
const hint = ref("");
const loading = ref(false);
const sending = ref(false);
const connected = ref(false);
const messagePanel = ref<HTMLElement>();
let session: CustomerServiceIMSession | undefined;
let loadVersion = 0;

const requestedMerID = computed(() => {
  const raw = route.query.mer_id;
  const value = Number(Array.isArray(raw) ? raw[0] : raw);
  return Number.isSafeInteger(value) && value > 0 ? value : 0;
});

const activeConversationID = computed(() => session?.conversationID || 0);

function messageKey(message: PteChatMessage) {
  return message.serverMsgId || message.clientMsgId;
}

function mergeMessage(message: PteChatMessage) {
  if (String(message.conversationId) !== String(activeConversationID.value)) return;
  const index = messages.value.findIndex((item) => item.clientMsgId === message.clientMsgId || messageKey(item) === messageKey(message));
  if (index >= 0) messages.value.splice(index, 1, { ...messages.value[index], ...message });
  else messages.value.push(message);
  messages.value.sort((a, b) => (a.serverSeq || 0) - (b.serverSeq || 0) || a.createdAt - b.createdAt);
  void scrollToLatest();
}

function mergeSendAck(ack: PteSendAck) {
  const message = messages.value.find((item) => item.clientMsgId === ack.clientMsgId);
  if (message) {
    message.serverMsgId = ack.serverMsgId;
    message.serverSeq = ack.serverSeq;
    message.sendState = "sent";
  }
}

const imListener: PteChatListener = {
  onConnectionChanged: (value) => {
    connected.value = value;
    if (value) hint.value = "";
  },
  onMessage: mergeMessage,
  onSendAck: mergeSendAck,
  onError: (message) => {
    hint.value = message || "客服 IM 连接异常";
  },
};

async function scrollToLatest() {
  await nextTick();
  if (messagePanel.value) messagePanel.value.scrollTop = messagePanel.value.scrollHeight;
}

function closeSession() {
  session?.close();
  session = undefined;
  connected.value = false;
  messages.value = [];
}

async function selectThread(thread: CustomerServiceThread, version = loadVersion) {
  closeSession();
  selectedThread.value = thread;
  hint.value = "正在连接客服 IM…";
  try {
    const nextSession = await connectCustomerServiceIM(thread, imListener);
    if (version !== loadVersion) {
      nextSession.close();
      return;
    }
    session = nextSession;
    const history = await nextSession.history();
    if (version !== loadVersion || session !== nextSession) return;
    messages.value = history;
    hint.value = nextSession.client.isConnected() ? "" : "IM 已初始化，正在连接…";
    void nextSession.client.markConversationRead(nextSession.conversationID).catch(() => undefined);
    void scrollToLatest();
  } catch (error) {
    if (version === loadVersion) hint.value = error instanceof Error ? error.message : "客服 IM 初始化失败";
  }
}

async function load() {
  const version = ++loadVersion;
  closeSession();
  selectedThread.value = undefined;
  if (!user.isLogin) {
    await router.replace({ name: "login", query: { redirect: route.fullPath } });
    return;
  }
  loading.value = true;
  hint.value = "加载客服会话…";
  try {
    const page = await fetchCustomerServiceThreads();
    if (version !== loadVersion) return;
    threads.value = page.list || [];
    let target = threads.value.find((item) => item.mer_id === requestedMerID.value);
    if (!target && requestedMerID.value) {
      target = await openCustomerServiceThread(requestedMerID.value);
      if (version !== loadVersion) return;
      threads.value = [target, ...threads.value.filter((item) => item.thread_id !== target?.thread_id)];
    }
    target ||= threads.value[0];
    if (!target) {
      hint.value = "暂无客服会话。请从商品或店铺页面发起咨询。";
      return;
    }
    await selectThread(target, version);
  } catch (error) {
    if (version === loadVersion) hint.value = error instanceof Error ? error.message : "客服会话加载失败";
  } finally {
    if (version === loadVersion) loading.value = false;
  }
}

async function sendMessage() {
  const text = draft.value.trim();
  if (!text || !session || sending.value) return;
  if (!session.client.isConnected()) {
    hint.value = "客服 IM 正在连接，请稍后再试";
    return;
  }
  sending.value = true;
  try {
    const clientMsgId = await session.client.sendText(session.conversationID, text);
    mergeMessage({
      clientMsgId,
      conversationId: String(session.conversationID),
      senderId: session.client.currentUserId(),
      type: "text",
      content: { text },
      createdAt: Date.now(),
      sendState: "sending",
    });
    draft.value = "";
  } catch (error) {
    hint.value = error instanceof Error ? error.message : "消息发送失败";
  } finally {
    sending.value = false;
  }
}

function isMine(message: PteChatMessage) {
  return message.senderId === session?.client.currentUserId();
}

function messageText(message: PteChatMessage) {
  const text = message.content?.text;
  return typeof text === "string" && text.trim() ? text : "暂不支持展示此消息类型";
}

function threadTitle(thread: CustomerServiceThread) {
  return thread.mer_name || `店铺 #${thread.mer_id}`;
}

function formatTime(value: number) {
  if (!value) return "";
  return new Intl.DateTimeFormat("zh-CN", { hour: "2-digit", minute: "2-digit" }).format(new Date(value));
}

watch(() => route.fullPath, () => void load(), { immediate: true });
onBeforeUnmount(() => {
  loadVersion += 1;
  closeSession();
});
</script>

<template>
  <div class="pc-container service-page">
    <section class="service-card">
      <aside class="thread-list" aria-label="客服会话列表">
        <div class="thread-list__title"><strong>我的客服</strong><span>{{ threads.length }} 个会话</span></div>
        <button
          v-for="thread in threads"
          :key="thread.thread_id"
          class="thread-row"
          :class="{ active: selectedThread?.thread_id === thread.thread_id }"
          type="button"
          @click="void selectThread(thread)"
        >
          <span class="thread-row__avatar">客</span>
          <span class="thread-row__body">
            <b>{{ threadTitle(thread) }}</b>
            <em>{{ thread.last_msg || "点击开始咨询" }}</em>
          </span>
          <i v-if="thread.user_unread">{{ thread.user_unread }}</i>
        </button>
        <p v-if="!threads.length && !loading" class="thread-empty">从商品或店铺页点击“咨询客服”即可创建会话。</p>
      </aside>

      <main class="chat-panel">
        <header class="chat-panel__header">
          <div>
            <h1>{{ selectedThread ? threadTitle(selectedThread) : "在线客服" }}</h1>
            <p>消息由 PTE Live IM 加密传输</p>
          </div>
          <span class="connection" :class="{ online: connected }">{{ connected ? "已连接" : "未连接" }}</span>
        </header>
        <p v-if="hint" class="chat-hint">{{ hint }}</p>
        <div ref="messagePanel" class="message-list" aria-live="polite">
          <p v-if="!messages.length && selectedThread && !loading" class="message-empty">开始与店铺客服沟通吧</p>
          <article v-for="message in messages" :key="messageKey(message)" class="message" :class="{ mine: isMine(message) }">
            <span class="message__avatar">{{ isMine(message) ? "我" : "客" }}</span>
            <div>
              <p class="message__meta">{{ isMine(message) ? "我" : "客服" }} {{ formatTime(message.createdAt) }}</p>
              <p class="message__bubble">{{ messageText(message) }}</p>
              <small v-if="isMine(message) && message.sendState === 'sending'">发送中…</small>
            </div>
          </article>
        </div>
        <form class="composer" @submit.prevent="void sendMessage()">
          <textarea v-model="draft" :disabled="!selectedThread" maxlength="1000" placeholder="输入消息，按发送按钮提交" @keydown.enter.exact.prevent="void sendMessage()" />
          <div><span>文本消息通过 pte-live-im-sdk 直接发送</span><button class="pc-btn" type="submit" :disabled="!draft.trim() || !selectedThread || sending">{{ sending ? "发送中…" : "发送" }}</button></div>
        </form>
      </main>
    </section>
  </div>
</template>

<style scoped>
.service-page { padding-top: 1.7rem; padding-bottom: 2.8rem; }
.service-card { display: grid; grid-template-columns: 300px minmax(0, 1fr); min-height: 660px; border: 1px solid var(--pc-line); background: #fff; box-shadow: var(--pc-shadow); }
.thread-list { border-right: 1px solid var(--pc-line); background: #fafafa; overflow: auto; }.thread-list__title { display: flex; justify-content: space-between; align-items: center; padding: 1.2rem 1.1rem; border-bottom: 1px solid var(--pc-line); }.thread-list__title span { color: #999; font-size: .78rem; }
.thread-row { position: relative; display: flex; width: 100%; gap: .75rem; padding: 1rem 1.1rem; border: 0; border-bottom: 1px solid #eee; background: transparent; text-align: left; cursor: pointer; }.thread-row:hover, .thread-row.active { background: #fff0ee; }.thread-row__avatar, .message__avatar { display: grid; flex: 0 0 auto; place-items: center; width: 38px; height: 38px; border-radius: 50%; color: #fff; background: #ef3727; font-size: .85rem; }.thread-row__body { display: grid; min-width: 0; gap: .25rem; }.thread-row__body b { color: #444; font-size: .9rem; }.thread-row__body em { overflow: hidden; color: #999; font-size: .78rem; font-style: normal; text-overflow: ellipsis; white-space: nowrap; }.thread-row i { position: absolute; top: .8rem; right: .9rem; min-width: 18px; padding: 1px 5px; border-radius: 9px; color: #fff; background: #ef3727; font-size: .68rem; font-style: normal; text-align: center; }.thread-empty { margin: 1.2rem; color: #999; font-size: .84rem; line-height: 1.7; }
.chat-panel { display: grid; grid-template-rows: auto auto minmax(300px, 1fr) auto; min-width: 0; }.chat-panel__header { display: flex; align-items: center; justify-content: space-between; padding: 1.05rem 1.4rem; border-bottom: 1px solid var(--pc-line); }.chat-panel__header h1 { margin: 0; color: #333; font-size: 1.05rem; }.chat-panel__header p { margin: .35rem 0 0; color: #999; font-size: .75rem; }.connection { padding: .3rem .55rem; border-radius: 12px; color: #999; background: #f2f2f2; font-size: .75rem; }.connection.online { color: #16814c; background: #e8f8ef; }.chat-hint { margin: 0; padding: .7rem 1.4rem; color: #9b6416; background: #fff8e9; font-size: .82rem; }
.message-list { overflow: auto; padding: 1.3rem 1.4rem; background: #f7f7f7; }.message-empty { margin: 7rem 0; color: #aaa; text-align: center; }.message { display: flex; max-width: 76%; gap: .7rem; margin-bottom: 1.1rem; }.message.mine { flex-direction: row-reverse; margin-right: 0; margin-left: auto; }.message.mine > div { text-align: right; }.message.mine .message__avatar { background: #4f83cc; }.message__meta { margin: 0 0 .3rem; color: #aaa; font-size: .7rem; }.message__bubble { display: inline-block; max-width: 100%; margin: 0; padding: .65rem .85rem; border-radius: 3px; color: #444; background: #fff; text-align: left; line-height: 1.55; white-space: pre-wrap; word-break: break-word; }.message.mine .message__bubble { color: #fff; background: #ef3727; }.message small { display: block; margin-top: .2rem; color: #aaa; font-size: .67rem; }
.composer { padding: .9rem 1.1rem; border-top: 1px solid var(--pc-line); }.composer textarea { width: 100%; min-height: 86px; resize: vertical; border: 1px solid #ddd; padding: .65rem; color: #444; font: inherit; outline: none; }.composer textarea:focus { border-color: #ef3727; }.composer > div { display: flex; align-items: center; justify-content: space-between; margin-top: .55rem; }.composer span { color: #999; font-size: .74rem; }.composer button { min-width: 90px; border-radius: 0; }
@media (max-width: 760px) { .service-card { grid-template-columns: 1fr; }.thread-list { max-height: 210px; border-right: 0; border-bottom: 1px solid var(--pc-line); }.chat-panel { min-height: 520px; }.message { max-width: 90%; } }
</style>
