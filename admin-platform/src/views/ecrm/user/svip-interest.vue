<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { createSvipInterest, deleteSvipInterest, listSvipInterests, updateSvipInterest, type SvipInterest, type SvipInterestInput } from '#/api/core/platform-svip-interest';

const canManage = ref(false);
const rows = ref<SvipInterest[]>([]);
const open = ref(false);
const saving = ref(false);
const editing = ref<SvipInterest>();
const form = reactive<SvipInterestInput>({ name: '', description: '', icon_url: '', status: 1, sort: 0 });
function reset() { Object.assign(form, { name: '', description: '', icon_url: '', status: 1, sort: 0, version: undefined }); }
async function load() { if (canManage.value) rows.value = (await listSvipInterests()).list || []; }
function edit(row?: SvipInterest) { editing.value = row; if (!row) reset(); else Object.assign(form, { name: row.name, description: row.description, icon_url: row.icon_url, status: row.status, sort: row.sort, version: row.version }); open.value = true; }
async function save() { const icon = form.icon_url.trim(); if (!form.name.trim() || (icon && !icon.startsWith('/demo/') && !icon.startsWith('https://'))) { ElMessage.warning('请填写权益名称；图标仅允许 /demo/ 或 https:// 地址'); return; } saving.value = true; try { const data = { ...form, name: form.name.trim(), description: form.description.trim(), icon_url: icon }; if (editing.value) await updateSvipInterest(editing.value.id, data); else await createSvipInterest(data); ElMessage.success('会员权益已保存；套餐只能选择启用权益'); open.value = false; await load(); } finally { saving.value = false; } }
async function remove(row: SvipInterest) { try { await ElMessageBox.confirm(`删除“${row.name}”前会检查所有启用会员类型；如仍被使用，服务端将拒绝删除。`, '删除会员权益', { type: 'warning' }); await deleteSvipInterest(row.id); ElMessage.success('会员权益已逻辑删除'); await load(); } catch {} }
onMounted(async () => { const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]); canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && codes.includes('user.svip.interest.manage'); await load(); });
</script>

<template>
  <Page title="会员权益" description="维护 C 端 SVIP 可展示权益。已启用会员类型只能选择启用权益；删除会检查现有可售类型，已付款订单的权益快照不受影响。">
    <el-alert v-if="!canManage" title="当前账号没有会员权益维护权限" type="warning" :closable="false" />
    <el-card v-else shadow="never"><template #header><div class="flex justify-between"><span>付费会员权益</span><el-button type="primary" @click="edit()">新增权益</el-button></div></template><el-table :data="rows"><el-table-column prop="name" label="权益名称" min-width="140"/><el-table-column prop="description" label="权益说明" min-width="260" show-overflow-tooltip/><el-table-column label="图标" width="100"><template #default="{ row }">{{ row.icon_url || '—' }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column prop="sort" label="排序" width="80"/><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table></el-card>
    <el-dialog v-model="open" :title="editing ? '编辑会员权益' : '新增会员权益'" width="560px" destroy-on-close><el-form label-width="88px"><el-form-item label="权益名称" required><el-input v-model="form.name" maxlength="64"/></el-form-item><el-form-item label="权益说明"><el-input v-model="form.description" type="textarea" :rows="4" maxlength="500" show-word-limit/></el-form-item><el-form-item label="图标地址"><el-input v-model="form.icon_url" maxlength="1024" placeholder="/demo/ 或 https://；可留空"/></el-form-item><el-form-item label="状态"><el-radio-group v-model="form.status"><el-radio :value="1">启用</el-radio><el-radio :value="0">停用</el-radio></el-radio-group></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0"/></el-form-item></el-form><template #footer><el-button @click="open = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
