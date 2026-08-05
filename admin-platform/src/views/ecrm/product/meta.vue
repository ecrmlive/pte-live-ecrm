<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import type { AttachmentItem } from '#/api/core/attachment';

import {
  deleteProductGuarantee, deleteProductLabel, deleteProductParameterTemplate,
  fetchProductGuarantees, fetchProductLabels, fetchProductParameterTemplates,
  saveProductGuarantee, saveProductLabel, saveProductParameterTemplate,
  type ProductGuaranteeRow, type ProductLabelRow, type ProductParameterTemplateRow,
} from '#/api/core/ecrm';

const route = useRoute();
type Kind = 'guarantee' | 'label' | 'parameter';
const kind = computed<Kind>(() => route.path.includes('guarantee') ? 'guarantee' : route.path.includes('spec') ? 'parameter' : 'label');
const title = computed(() => ({ label: '商品标签', guarantee: '保障服务', parameter: '平台商品参数' })[kind.value]);
const description = computed(() => ({ label: '维护商品运营标签；不包含店铺私有标签。', guarantee: '维护平台统一售后保障说明；实际退款仍由售后状态机处理。', parameter: '维护平台商品参数模板；参数值仅限展示与筛选。' })[kind.value]);
const loading = ref(false); const dialog = ref(false); const editing = ref<number>();
const labels = ref<ProductLabelRow[]>([]); const guarantees = ref<ProductGuaranteeRow[]>([]); const parameters = ref<ProductParameterTemplateRow[]>([]);
const form = reactive({ name: '', description: '', color: '', content: '', icon_url: '', values: '', sort: 0, status: 1 }); const iconPicker = ref(false); const icon = ref<AttachmentItem>();

function reset() { editing.value = undefined; icon.value=undefined; Object.assign(form, { name: '', description: '', color: '', content: '', icon_url: '', values: '', sort: 0, status: 1 }); }
function open(row?: ProductLabelRow | ProductGuaranteeRow | ProductParameterTemplateRow) {
  reset(); if (row) { editing.value = row.id; form.name = row.name; form.sort = row.sort; form.status = row.status;
    if ('description' in row) { form.description = row.description; form.color = row.color; }
    if ('content' in row) { form.content = row.content; form.icon_url = row.icon_url; }
    if ('values_json' in row) { try { form.values = (JSON.parse(row.values_json) as string[]).join('、'); } catch { form.values = ''; } }
  } dialog.value = true;
}
async function load() { loading.value = true; try {
  if (kind.value === 'label') labels.value = (await fetchProductLabels()).list || [];
  else if (kind.value === 'guarantee') guarantees.value = (await fetchProductGuarantees()).list || [];
  else parameters.value = (await fetchProductParameterTemplates()).list || [];
} finally { loading.value = false; } }
async function save() {
  if (!form.name.trim()) { ElMessage.warning('请填写名称'); return; }
  if (kind.value === 'label') await saveProductLabel(editing.value, { name: form.name.trim(), description: form.description.trim(), color: form.color.trim(), sort: form.sort, status: form.status });
  else if (kind.value === 'guarantee') await saveProductGuarantee(editing.value, { name: form.name.trim(), content: form.content.trim(), icon_url: form.icon_url.trim(), sort: form.sort, status: form.status });
  else { const values = form.values.split(/[、,，\n]/).map((value) => value.trim()).filter(Boolean); if (!values.length) { ElMessage.warning('请至少填写一个参数值'); return; } await saveProductParameterTemplate(editing.value, { name: form.name.trim(), values, sort: form.sort, status: form.status }); }
  dialog.value = false; ElMessage.success('已保存'); await load();
}
function parameterValues(raw: string) { try { const values = JSON.parse(raw); return Array.isArray(values) ? values.map((value) => String(value)) : []; } catch { return []; } }
async function remove(id: number, name: string) { try { await ElMessageBox.confirm(`确认删除“${name}”？已关联的商品仍保留其历史快照。`, '删除确认', { type: 'warning' }); if (kind.value === 'label') await deleteProductLabel(id); else if (kind.value === 'guarantee') await deleteProductGuarantee(id); else await deleteProductParameterTemplate(id); ElMessage.success('已删除'); await load(); } catch { /* 用户取消或请求失败由统一层反馈。 */ } }
onMounted(() => void load());
watch(kind, () => void load());
</script>

<template>
  <Page :title="title" :description="description">
    <template #extra><el-button type="primary" @click="open()">新增{{ title }}</el-button></template>
    <el-card shadow="never">
      <el-table v-if="kind === 'label'" v-loading="loading" :data="labels" row-key="id"><el-table-column label="ID" prop="id" width="90" /><el-table-column label="名称" prop="name" min-width="140" /><el-table-column label="说明" prop="description" min-width="240" /><el-table-column label="颜色" width="120"><template #default="{ row }"><span :style="{ color: row.color || undefined }">{{ row.color || '默认' }}</span></template></el-table-column><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="open(row)">编辑</el-button><el-button link type="danger" @click="remove(row.id, row.name)">删除</el-button></template></el-table-column></el-table>
      <el-table v-else-if="kind === 'guarantee'" v-loading="loading" :data="guarantees" row-key="id"><el-table-column label="ID" prop="id" width="90" /><el-table-column label="名称" prop="name" min-width="140" /><el-table-column label="保障说明" prop="content" min-width="320" show-overflow-tooltip /><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="open(row)">编辑</el-button><el-button link type="danger" @click="remove(row.id, row.name)">删除</el-button></template></el-table-column></el-table>
      <el-table v-else v-loading="loading" :data="parameters" row-key="id"><el-table-column label="ID" prop="id" width="90" /><el-table-column label="参数名称" prop="name" min-width="160" /><el-table-column label="候选值" min-width="280"><template #default="{ row }"><el-tag v-for="value in parameterValues(row.values_json)" :key="value" class="mr-1">{{ value }}</el-tag><span v-if="!parameterValues(row.values_json).length" class="text-gray-400">参数值异常</span></template></el-table-column><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="open(row)">编辑</el-button><el-button link type="danger" @click="remove(row.id, row.name)">删除</el-button></template></el-table-column></el-table>
    </el-card>
    <el-dialog v-model="dialog" :title="`${editing ? '编辑' : '新增'}${title}`" width="560px" destroy-on-close><el-form label-width="92px"><el-form-item label="名称" required><el-input v-model="form.name" maxlength="64" show-word-limit /></el-form-item><template v-if="kind === 'label'"><el-form-item label="说明"><el-input v-model="form.description" maxlength="255" /></el-form-item><el-form-item label="颜色"><el-input v-model="form.color" placeholder="#2563eb（可选）" maxlength="32" /></el-form-item></template><template v-else-if="kind === 'guarantee'"><el-form-item label="保障说明"><el-input v-model="form.content" type="textarea" :rows="4" maxlength="1000" show-word-limit /></el-form-item><el-form-item label="保障图标"><el-button @click="iconPicker=true">选择平台图片</el-button><el-image v-if="form.icon_url" :src="form.icon_url" class="ml-2 h-8 w-8" fit="cover" /><el-button v-if="form.icon_url" link type="danger" @click="form.icon_url=''">清除</el-button></el-form-item></template><el-form-item v-else label="参数值" required><el-input v-model="form.values" type="textarea" :rows="3" placeholder="用中文顿号、逗号或换行分隔，例如：小杯、中杯、大杯" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" class="w-full" /></el-form-item><el-form-item label="启用"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialog = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog><ImagePickerDialog v-model:open="iconPicker" default-library="system" kind="image" :limit="1" @select="items=>{icon=items[0];form.icon_url=items[0]?.attachment_src||''}" />
  </Page>
</template>
