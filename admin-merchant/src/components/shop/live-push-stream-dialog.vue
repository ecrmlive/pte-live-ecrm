<script setup lang="ts">
import type { LiveRoomListItem, LiveStreamInfo } from '#/api/core/live';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { getLiveStreamApi } from '#/api/core/live';

defineOptions({ name: 'LivePushStreamDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveRoomListItem | Record<string, never>;
}>();

const loading = ref(false);
const qrLoading = ref(false);
const activePanels = ref(['push']);
const streamInfo = reactive<LiveStreamInfo>({
  pull_url_flv: '',
  pull_url_hls: '',
  pull_url_rtmp: '',
  pull_url_webrtc: '',
  push_auth_query: '',
  push_stream_key: '',
  push_server_url: '',
  push_url: '',
  stream_status: 0,
});

const pushQrImage = computed(() => {
  const url = streamInfo.push_url;
  if (!url) return '';
  return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(url)}`;
});

const showQrSpinner = computed(() => loading.value || qrLoading.value);

// OBS 的“服务器密钥”由流名和可选鉴权参数组成，不能再次放入完整 rtmp 地址。
const obsStreamKey = computed(() => {
  const streamKey = streamInfo.push_stream_key?.trim();
  if (streamKey) return streamKey;

  const streamName = streamInfo.stream_name?.trim();
  const authQuery = streamInfo.push_auth_query?.trim();
  if (streamName) return authQuery ? `${streamName}?${authQuery}` : streamName;

  return '';
});

watch(pushQrImage, (url) => {
  if (url) {
    qrLoading.value = true;
  } else if (!loading.value) {
    qrLoading.value = false;
  }
});

function onQrImageSettled() {
  qrLoading.value = false;
}

function resetStreamInfo() {
  Object.assign(streamInfo, {
    pull_url_flv: '',
    pull_url_hls: '',
    pull_url_rtmp: '',
    pull_url_webrtc: '',
    push_auth_query: '',
    push_stream_key: '',
    push_server_url: '',
    push_url: '',
    stream_status: 0,
  });
}

const pullRows = computed(() => {
  const rows: Array<{ label: string; url: string }> = [];
  if (streamInfo.pull_url_rtmp) rows.push({ label: 'RTMP', url: streamInfo.pull_url_rtmp });
  if (streamInfo.pull_url_flv) rows.push({ label: 'FLV', url: streamInfo.pull_url_flv });
  if (streamInfo.pull_url_hls) rows.push({ label: 'HLS', url: streamInfo.pull_url_hls });
  if (streamInfo.pull_url_webrtc) rows.push({ label: 'WebRTC', url: streamInfo.pull_url_webrtc });
  return rows;
});

async function loadStream() {
  const liveId = props.row?.live_id;
  if (!liveId) return;
  loading.value = true;
  qrLoading.value = true;
  resetStreamInfo();
  try {
    const res = await getLiveStreamApi(liveId);
    Object.assign(streamInfo, res);
    if (!streamInfo.push_url && props.row?.push_url) {
      streamInfo.push_url = String(props.row.push_url);
    }
    if (!streamInfo.push_url) {
      qrLoading.value = false;
    }
  } catch {
    if (props.row?.push_url) {
      streamInfo.push_url = String(props.row.push_url);
    } else {
      qrLoading.value = false;
    }
  } finally {
    loading.value = false;
  }
}

async function copyText(text?: string) {
  if (!text) return;
  try {
    await navigator.clipboard.writeText(text);
    ElMessage.success('已复制');
  } catch {
    ElMessage.warning('复制失败，请手动复制');
  }
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    void loadStream();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :footer="false"
    class="w-[min(680px,96vw)]"
    title="推流地址"
  >
    <div class="push-stream-dialog min-h-[200px]">
      <div class="flex flex-col items-center pb-5">
        <div
          v-loading="showQrSpinner"
          class="push-qr-box"
          element-loading-text="二维码生成中..."
        >
          <img
            v-if="pushQrImage"
            :src="pushQrImage"
            alt="扫码推流"
            class="push-qr-box__img"
            @error="onQrImageSettled"
            @load="onQrImageSettled"
          />
          <el-empty
            v-else-if="!showQrSpinner"
            :image-size="80"
            description="暂无推流地址，请确认已配置腾讯云 LVB"
          />
        </div>
        <p v-if="streamInfo.push_url && !showQrSpinner" class="push-stream-dialog__scan-hint">
          扫码推流
        </p>
      </div>

      <div v-if="streamInfo.push_server_url" class="stream-url-row mb-2">
        <span class="stream-url-row__label">服务器地址</span>
        <el-input :model-value="streamInfo.push_server_url" readonly />
        <el-button link type="primary" @click="copyText(streamInfo.push_server_url)">复制</el-button>
      </div>
      <div v-if="obsStreamKey" class="stream-url-row mb-2">
        <span class="stream-url-row__label">服务器密钥</span>
        <el-input :model-value="obsStreamKey" readonly />
        <el-button link type="primary" @click="copyText(obsStreamKey)">复制</el-button>
      </div>
      <p v-if="streamInfo.push_server_url && obsStreamKey" class="push-stream-dialog__tip">
        OBS 等桌面推流：服务器填“服务器地址”，串流密钥填“服务器密钥”。
      </p>

      <el-collapse v-model="activePanels" class="mb-3">
        <el-collapse-item name="push" title="完整推流地址（备用）">
          <div v-if="streamInfo.push_url" class="flex items-center gap-2">
            <el-input :model-value="streamInfo.push_url" readonly />
            <el-button link type="primary" @click="copyText(streamInfo.push_url)">复制</el-button>
          </div>
          <p v-else class="text-xs text-gray-500">暂无推流地址</p>
        </el-collapse-item>
      </el-collapse>

      <section v-if="pullRows.length" class="pull-stream-card" aria-label="拉流地址">
        <header class="pull-stream-card__header">
          <span class="pull-stream-card__title">拉流地址</span>
          <span class="pull-stream-card__hint">用于播放器预览或第三方接入</span>
        </header>
        <div
          v-for="item in pullRows"
          :key="item.label"
          class="pull-stream-card__row"
        >
          <span class="pull-stream-card__protocol">{{ item.label }}</span>
          <code class="pull-stream-card__url" :title="item.url">{{ item.url }}</code>
          <el-button link type="primary" @click="copyText(item.url)">复制</el-button>
        </div>
      </section>
    </div>
  </Modal>
</template>

<style scoped>
.push-stream-dialog {
  color: var(--el-text-color-primary);
}

.push-stream-dialog__scan-hint,
.push-stream-dialog__tip,
.stream-url-row__label {
  color: var(--el-text-color-secondary);
}

.push-stream-dialog__scan-hint {
  margin-top: 8px;
  font-size: 14px;
}

.push-stream-dialog__tip {
  margin-bottom: 12px;
  font-size: 12px;
  line-height: 20px;
}

.stream-url-row {
  display: grid;
  grid-template-columns: 88px minmax(0, 1fr) auto;
  align-items: center;
  gap: 8px;
}

.stream-url-row__label {
  font-size: 14px;
}

.push-qr-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 216px;
  height: 216px;
  padding: 8px;
  background: #ffffff;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
}

.push-qr-box__img {
  display: block;
  width: 200px;
  height: 200px;
  object-fit: contain;
}

.pull-stream-card {
  padding: 14px 16px;
  color: #1e293b;
  background: #f8fafc;
  border: 1px solid #dbe4ef;
  border-radius: 12px;
}

.pull-stream-card__header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.pull-stream-card__title {
  font-size: 14px;
  font-weight: 600;
}

.pull-stream-card__hint {
  color: #64748b;
  font-size: 12px;
}

.pull-stream-card__row {
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr) auto;
  align-items: start;
  gap: 12px;
  padding: 11px 0;
  border-bottom: 1px solid #e2e8f0;
}

.pull-stream-card__row:last-child {
  padding-bottom: 0;
  border-bottom: 0;
}

.pull-stream-card__protocol {
  color: #475569;
  font-size: 12px;
  font-weight: 600;
  line-height: 20px;
}

.pull-stream-card__url {
  min-width: 0;
  color: #334155;
  font-family: var(--font-mono, ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace);
  font-size: 12px;
  font-style: normal;
  line-height: 20px;
  overflow-wrap: anywhere;
}

:global(.dark) .pull-stream-card {
  color: #e2e8f0;
  background: #151b26;
  border-color: #334155;
}

:global(.dark) .pull-stream-card__hint,
:global(.dark) .pull-stream-card__protocol {
  color: #94a3b8;
}

:global(.dark) .pull-stream-card__url {
  color: #e2e8f0;
}

:global(.dark) .pull-stream-card__row {
  border-color: #334155;
}

@media (max-width: 640px) {
  .stream-url-row {
    grid-template-columns: 1fr auto;
  }

  .stream-url-row__label {
    grid-column: 1 / -1;
  }

  .pull-stream-card__header {
    align-items: flex-start;
    flex-direction: column;
    gap: 2px;
  }
}
</style>
