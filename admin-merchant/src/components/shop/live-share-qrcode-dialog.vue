<script setup lang="ts">
import type {
  LiveRoomListItem,
  LiveRoomShareH5Item,
  LiveRoomShareQrcodeResult,
} from '#/api/core/live';

import { computed, ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { ElButton, ElMessage } from 'element-plus';

import { getLiveRoomShareQrcodeApi } from '#/api/core/live';

defineOptions({ name: 'LiveShareQrcodeDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveRoomListItem | Record<string, never>;
}>();

const loading = ref(false);
const info = ref<LiveRoomShareQrcodeResult>({});

const h5Items = computed<LiveRoomShareH5Item[]>(() => {
  if (info.value.h5_list?.length) {
    return info.value.h5_list;
  }
  if (info.value.h5_qrcode || info.value.h5_url) {
    return [
      {
        domain: '平台 H5',
        h5_qrcode: info.value.h5_qrcode,
        h5_url: info.value.h5_url,
      },
    ];
  }
  return [];
});

const roomTitle = computed(() => {
  const name = props.row?.name;
  if (name) return name;
  const liveId = props.row?.live_id;
  return liveId ? `直播间 #${liveId}` : '';
});

async function loadQrcodes() {
  const liveId = props.row?.live_id;
  const roomId = props.row?.roomid ?? (props.row as { room_id?: number }).room_id;
  if (!liveId) return;
  loading.value = true;
  info.value = {};
  try {
    info.value = await getLiveRoomShareQrcodeApi({
      live_id: liveId,
      room_id: roomId || liveId,
    });
  } finally {
    loading.value = false;
  }
}

async function copyH5Url(url?: string) {
  if (!url) return;
  try {
    await navigator.clipboard.writeText(url);
    ElMessage.success('链接已复制');
  } catch {
    ElMessage.error('复制失败，请手动复制');
  }
}

function buildQrcodeFilename(suffix: string) {
  const liveId = props.row?.live_id ?? 'live';
  const rawName = String(props.row?.name ?? '').trim();
  const safeName = rawName.replace(/[\\/:*?"<>|]/g, '_').slice(0, 24);
  return safeName ? `${safeName}-${liveId}-${suffix}.png` : `live-room-${liveId}-${suffix}.png`;
}

function h5FilenameSuffix(domain?: string) {
  const slug = String(domain ?? 'h5')
    .replace(/[\\/:*?"<>|]/g, '_')
    .replace(/\s+/g, '-')
    .slice(0, 40);
  return slug || 'h5';
}

function saveQrcodeImage(dataUri: string | undefined, filename: string) {
  if (!dataUri) {
    ElMessage.warning('暂无二维码');
    return;
  }
  const link = document.createElement('a');
  link.href = dataUri;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

function saveWxQrcode() {
  saveQrcodeImage(info.value.wx_qrcode, buildQrcodeFilename('wx'));
}

function saveH5QrcodeItem(item: LiveRoomShareH5Item) {
  saveQrcodeImage(item.h5_qrcode, buildQrcodeFilename(h5FilenameSuffix(item.domain)));
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    void loadQrcodes();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[min(760px,96vw)]"
    title="分享二维码"
  >
    <div v-loading="loading" class="min-h-[220px]">
      <p v-if="roomTitle" class="mb-3 text-center text-sm text-foreground">
        {{ roomTitle }}
        <span class="text-muted-foreground"> · 扫码直达直播间</span>
      </p>
      <el-alert
        v-if="info.wx_error"
        class="mb-3"
        :closable="false"
        show-icon
        :title="info.wx_error"
        type="warning"
      />

      <section class="share-qrcode-section">
        <p class="share-qrcode-section__title">小程序直播间</p>
        <div class="share-qrcode-section__body share-qrcode-section__body--center">
          <div class="share-qrcode-card share-qrcode-card--wx">
            <div class="share-qrcode-card__img-wrap">
              <img
                v-if="info.wx_qrcode"
                :src="info.wx_qrcode"
                alt="小程序码"
                class="share-qrcode-card__img share-qrcode-card__img--round"
              />
              <el-empty v-else :image-size="72" description="暂无小程序码" />
            </div>
            <div class="share-qrcode-actions">
              <ElButton
                v-if="info.wx_qrcode"
                link
                size="small"
                type="primary"
                @click="saveWxQrcode"
              >
                保存二维码
              </ElButton>
            </div>
          </div>
        </div>
      </section>

      <section class="share-qrcode-section">
        <p class="share-qrcode-section__title">
          H5 直播间
          <span v-if="h5Items.length > 1" class="share-qrcode-section__hint">
            （共 {{ h5Items.length }} 个入口，含平台默认与商户域名）
          </span>
        </p>
        <div v-if="h5Items.length" class="h5-share-list">
          <div
            v-for="(item, index) in h5Items"
            :key="`${item.domain}-${index}`"
            class="share-qrcode-card share-qrcode-card--h5"
          >
            <p class="share-qrcode-card__label" :title="item.domain">{{ item.domain }}</p>
            <div class="share-qrcode-card__img-wrap share-qrcode-card__img-wrap--h5">
              <img
                v-if="item.h5_qrcode"
                :src="item.h5_qrcode"
                alt="H5码"
                class="share-qrcode-card__img"
              />
              <el-empty v-else :image-size="56" description="暂无 H5 码" />
            </div>
            <div class="share-qrcode-actions share-qrcode-actions--stack">
              <ElButton
                v-if="item.h5_url"
                link
                size="small"
                type="primary"
                @click="copyH5Url(item.h5_url)"
              >
                复制链接
              </ElButton>
              <ElButton
                v-if="item.h5_qrcode"
                link
                size="small"
                type="primary"
                @click="saveH5QrcodeItem(item)"
              >
                保存二维码
              </ElButton>
            </div>
          </div>
        </div>
        <el-empty v-else :image-size="72" description="暂无 H5 分享码" />
      </section>
    </div>
  </Modal>
</template>

<style scoped>
.share-qrcode-section + .share-qrcode-section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid hsl(var(--border));
}

.share-qrcode-section__title {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: hsl(var(--foreground));
  text-align: center;
}

.share-qrcode-section__hint {
  font-size: 12px;
  font-weight: 400;
  color: hsl(var(--muted-foreground));
}

.share-qrcode-section__body--center {
  display: flex;
  justify-content: center;
}

.h5-share-list {
  display: flex;
  gap: 16px;
  padding-bottom: 4px;
  overflow-x: auto;
  justify-content: safe center;
}

.share-qrcode-card {
  display: flex;
  flex: 0 0 auto;
  flex-direction: column;
  align-items: center;
}

.share-qrcode-card--wx {
  width: 200px;
}

.share-qrcode-card--h5 {
  width: 168px;
}

.share-qrcode-card__label {
  width: 100%;
  margin: 0 0 8px;
  overflow: hidden;
  font-size: 12px;
  font-weight: 500;
  line-height: 1.4;
  color: hsl(var(--foreground));
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.share-qrcode-card__img-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 180px;
  height: 180px;
}

.share-qrcode-card__img-wrap--h5 {
  width: 148px;
  height: 148px;
}

.share-qrcode-card__img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.share-qrcode-card__img--round {
  border-radius: 9999px;
}

.share-qrcode-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 8px;
  align-items: center;
  justify-content: center;
  min-height: 28px;
  margin-top: 8px;
}

.share-qrcode-actions--stack {
  flex-direction: column;
  gap: 2px;
}
</style>
