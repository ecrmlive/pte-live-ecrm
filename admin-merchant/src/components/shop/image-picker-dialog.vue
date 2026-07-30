<script setup lang="ts">
import type { AttachmentCategory, AttachmentItem } from '#/api/core/attachment';

import { ElButton, ElDialog, ElEmpty, ElMessage, ElPagination, ElUpload } from 'element-plus';
import { onMounted, ref, watch } from 'vue';

import { listAttachmentCategoriesApi, listAttachmentsApi, uploadAttachmentApi } from '#/api/core/attachment';

type PickerItem = AttachmentItem & { file_id: number; file_path: string };

const open = defineModel<boolean>('open', { default: false });
const props = withDefaults(defineProps<{ defaultLibrary?: 'merchant' | 'system'; kind?: 'image' | 'video'; limit?: number }>(), { defaultLibrary: 'merchant', kind: 'image', limit: 1 });
const emit = defineEmits<{ select: [PickerItem[]] }>();
const categories = ref<AttachmentCategory[]>([]);
const categoryID = ref(0);
const files = ref<PickerItem[]>([]);
const selected = ref<number[]>([]);
const loading = ref(false);
const page = ref(1);
const total = ref(0);
const pageSize = 18;
async function loadCategories() { const result = await listAttachmentCategoriesApi(); categories.value = result.list ?? []; }
async function loadFiles() { loading.value = true; try { const result = await listAttachmentsApi({ category_id: categoryID.value || undefined, limit: pageSize, page: page.value, type: props.kind }); files.value = (result.list ?? []).map((item) => ({ ...item, file_id: item.attachment_id, file_path: item.attachment_src })); total.value = result.total ?? 0; } finally { loading.value = false; } }
function selectCategory(id: number) { categoryID.value = id; page.value = 1; selected.value = []; void loadFiles(); }
function toggle(row: PickerItem) { const index = selected.value.indexOf(row.attachment_id); if (index >= 0) { selected.value.splice(index, 1); return; } if (selected.value.length >= props.limit) { ElMessage.warning(`最多选择 ${props.limit} 张图片`); return; } selected.value.push(row.attachment_id); }
async function upload({ file }: { file: File }) { if (props.kind === 'image' ? !file.type.startsWith('image/') : !file.type.startsWith('video/')) { ElMessage.warning(props.kind === 'image' ? '请选择图片文件' : '请选择视频文件'); return; } await uploadAttachmentApi(file, categoryID.value); await loadFiles(); ElMessage.success(props.kind === 'image' ? '图片已上传' : '视频已上传'); }
function confirm() { const rows = files.value.filter((item) => selected.value.includes(item.attachment_id)); if (!rows.length) { ElMessage.warning('请选择图片'); return; } emit('select', rows); open.value = false; }
async function initialize() { selected.value = []; await Promise.all([loadCategories(), loadFiles()]); }
watch(open, (visible) => { if (visible) void initialize(); });
onMounted(() => { if (open.value) void initialize(); });
</script>

<template>
  <ElDialog v-model="open" :title="`选择${props.kind === 'image' ? '图片' : '视频'}素材`" width="900px" destroy-on-close append-to-body>
    <div class="picker"><aside class="picker__categories"><button :class="{ active: categoryID === 0 }" @click="selectCategory(0)">全部素材</button><button v-for="row in categories" :key="row.attachment_category_id" :class="{ active: categoryID === row.attachment_category_id }" @click="selectCategory(row.attachment_category_id)">{{ row.attachment_category_name }}</button></aside><section class="picker__content"><div class="picker__toolbar"><span>从本店素材库选择{{ props.kind === 'image' ? '图片' : '视频' }}</span><ElUpload :accept="props.kind === 'image' ? 'image/jpeg,image/png,image/webp,image/gif' : 'video/mp4,video/quicktime,video/webm'" :http-request="upload" :show-file-list="false"><ElButton type="primary">上传{{ props.kind === 'image' ? '图片' : '视频' }}</ElButton></ElUpload></div><div v-loading="loading" class="picker__grid"><ElEmpty v-if="!loading && files.length === 0" :description="`暂无${props.kind === 'image' ? '图片' : '视频'}素材`" /><button v-for="row in files" :key="row.attachment_id" class="picker__item" :class="{ selected: selected.includes(row.attachment_id) }" @click="toggle(row)"><img v-if="props.kind === 'image'" :src="row.attachment_src" :alt="row.attachment_name" /><video v-else :src="row.attachment_src" preload="metadata" /><span>{{ row.attachment_name }}</span></button></div><ElPagination v-if="total > pageSize" v-model:current-page="page" :page-size="pageSize" :total="total" small background layout="prev, pager, next" @current-change="loadFiles" /></section></div>
    <template #footer><ElButton @click="open = false">取消</ElButton><ElButton type="primary" @click="confirm">确定</ElButton></template>
  </ElDialog>
</template>

<style scoped lang="scss">
.picker { display: flex; min-height: 420px; border: 1px solid hsl(var(--border)); }.picker__categories { width: 170px; padding: 8px; border-right: 1px solid hsl(var(--border)); }.picker__categories button { display: block; width: 100%; padding: 8px; overflow: hidden; text-align: left; text-overflow: ellipsis; white-space: nowrap; border-radius: 4px; }.picker__categories button.active, .picker__categories button:hover { color: hsl(var(--primary)); background: hsl(var(--accent)); }.picker__content { display: flex; flex: 1; flex-direction: column; gap: 12px; min-width: 0; padding: 12px; }.picker__toolbar { display: flex; align-items: center; justify-content: space-between; }.picker__grid { display: grid; grid-template-columns: repeat(6, minmax(0, 1fr)); gap: 10px; min-height: 330px; align-content: start; }.picker__item { overflow: hidden; border: 2px solid transparent; border-radius: 4px; text-align: left; }.picker__item.selected { border-color: hsl(var(--primary)); }.picker__item img, .picker__item video { display: block; width: 100%; height: 92px; object-fit: cover; }.picker__item span { display: block; overflow: hidden; padding: 4px; font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }
</style>
