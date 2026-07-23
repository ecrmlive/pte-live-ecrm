<script setup lang="ts">
import type { LiveVodVideoItem } from '#/api/core/live-vod';
import type {
  LiveAnchorListItem,
  LiveRoomForm,
  LiveRoomListItem,
  LiveStreamInfo,
} from '#/api/core/live';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { reactive, ref, watch } from 'vue';

import {
  getLiveAnchorListApi,
  getLiveStreamApi,
  refreshLiveStreamUrlApi,
  updateLiveRoomApi,
} from '#/api/core/live';
import { buildLiveRoomPayload, defaultLiveRoomForm } from '#/utils/live-room-payload';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

import LiveRoomFormFields from '#/components/shop/live-room-form-fields.vue';
import VodVideoLibrary from '#/components/shop/vod-video-library.vue';

defineOptions({ name: 'LiveRoomEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveRoomListItem | null;
}>();

const emit = defineEmits<{ success: [] }>();

const formRef = ref<InstanceType<typeof LiveRoomFormFields>>();
const submitting = ref(false);
const streamLoading = ref(false);
const vodPickerOpen = ref(false);
const anchorOptions = ref<LiveAnchorListItem[]>([]);
const streamInfo = ref<LiveStreamInfo>({});

const form = reactive<LiveRoomForm & ReturnType<typeof defaultLiveRoomForm>>(defaultLiveRoomForm());

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

const [VodPickerModal, vodPickerModalApi] = useVbenModal({
  onOpenChange(isOpen) {
    vodPickerOpen.value = isOpen;
  },
});

watch(vodPickerOpen, (visible) => {
  if (visible) {
    vodPickerModalApi.open();
    return;
  }
  vodPickerModalApi.close();
});

function fillForm(row: LiveRoomListItem) {
  Object.assign(form, defaultLiveRoomForm(), {
    allow_chat: Number(row.allow_chat ?? 1),
    anchor_id: row.anchor_id || 0,
    anchor_name: row.anchor_name || '',
    anchor_wechat: row.anchor_wechat || '',
    background_img: resolveCosMediaUrl(String(row.background_img ?? '')),
    close_comment: Number(row.close_comment ?? 0),
    close_goods: Number(row.close_goods ?? 0),
    close_like: Number(row.close_like ?? 0),
    close_replay: Number(row.close_replay ?? 0),
    cover_img: resolveCosMediaUrl(row.cover_img || ''),
    enable_gift: Number(row.enable_gift ?? 1),
    enable_linkmic: Number(row.enable_linkmic ?? 0),
    enable_replay: Number(row.enable_replay ?? 1),
    end_time: row.end_time_text || '',
    feeds_img: resolveCosMediaUrl(row.feeds_img || ''),
    fire_value: Number(row.fire_value ?? 0),
    is_visible: Number(row.is_visible ?? 1),
    live_id: row.live_id,
    name: row.name || '',
    record_video_path: String(row.record_video_path ?? ''),
    record_vod_file_id: row.record_vod_file_id || '',
    record_vod_media_url: resolveCosMediaUrl(
      String(row.record_vod_media_url || row.record_video_path || ''),
    ),
    room_type: row.room_type || 1,
    share_img: resolveCosMediaUrl(row.share_img || ''),
    share_intro: row.share_intro || '',
    show_heat: Number(row.show_heat ?? 1),
    show_home: Number(row.show_home ?? 1),
    show_online_count: Number(row.show_online_count ?? 1),
    show_total_count: Number(row.show_total_count ?? 1),
    start_time: row.start_time_text || '',
    stream_orientation: Number(row.stream_orientation ?? 2),
    system_notice: row.system_notice || '',
    watch_password: row.watch_password || '',
  });
}

async function loadAnchors() {
  try {
    const res = await getLiveAnchorListApi({ list_rows: 200, page: 1, status: 1 });
    anchorOptions.value = res.list.data ?? [];
  } catch {
    anchorOptions.value = [];
  }
}

async function loadStream() {
  if (!form.live_id || form.room_type !== 1) return;
  try {
    streamInfo.value = await getLiveStreamApi(form.live_id);
  } catch {
    streamInfo.value = {};
  }
}

async function refreshStream() {
  if (!form.live_id) return;
  streamLoading.value = true;
  try {
    const res = await refreshLiveStreamUrlApi(form.live_id);
    streamInfo.value = res;
    ElMessage.success(res.msg || '地址已刷新');
  } finally {
    streamLoading.value = false;
  }
}

function onAnchorChange(anchorId?: number) {
  const hit = anchorOptions.value.find((item) => item.anchor_id === anchorId);
  if (hit) {
    form.anchor_name = hit.nick_name;
    form.anchor_wechat = hit.wechat || '';
  }
}

function onVodPicked(item: LiveVodVideoItem) {
  form.record_vod_file_id = item.file_id || '';
  form.record_vod_media_url = item.media_url || '';
  form.record_video_path = '';
  vodPickerOpen.value = false;
}

async function submit() {
  const valid = (await formRef.value?.validate()) ?? false;
  if (!valid) return;
  if (
    form.room_type === 2 &&
    !form.record_video_path?.trim() &&
    !form.record_vod_media_url?.trim() &&
    !form.record_vod_file_id?.trim()
  ) {
    ElMessage.error('录播房请选择录播视频');
    return;
  }
  submitting.value = true;
  try {
    const res = await updateLiveRoomApi(buildLiveRoomPayload(form));
    ElMessage.success(res.msg || '修改成功');
    open.value = false;
    emit('success');
  } finally {
    submitting.value = false;
  }
}

watch(
  () => [open.value, props.row] as const,
  ([visible, row]) => {
    if (visible) {
      modalApi.open();
      if (row) {
        fillForm(row);
        void loadAnchors();
        void loadStream();
      }
      return;
    }
    modalApi.close();
  },
);
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[900px]"
    title="编辑直播间"
  >
    <LiveRoomFormFields
      ref="formRef"
      :form="form"
      is-edit
      :anchor-options="anchorOptions"
      show-stream
      :stream-info="streamInfo"
      :stream-loading="streamLoading"
      @anchor-change="onAnchorChange"
      @pick-vod="vodPickerOpen = true"
      @refresh-stream="refreshStream"
    />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">提交</ElButton>
    </template>
  </Modal>

  <VodPickerModal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[960px]"
    title="选择录播视频"
  >
    <VodVideoLibrary picker @select="onVodPicked" />
  </VodPickerModal>
</template>
