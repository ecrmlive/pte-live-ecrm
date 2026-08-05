<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { createMemberLevel, deleteMemberLevel, listMemberLevels, updateMemberLevel, type MemberLevel, type MemberLevelInput } from '#/api/core/platform-member-level';

const canManage = ref(false);
const rows = ref<MemberLevel[]>([]);
const open = ref(false);
const saving = ref(false);
const editing = ref<MemberLevel>();
const form = reactive<MemberLevelInput>({ name: '', rank: 1, rules: '{\n  "description": "满足成长规则后自动升级"\n}', benefits: '[\n  "会员专享活动"\n]', status: 1 });
function reset() { Object.assign(form, { name: '', rank: 1, rules: '{\n  "description": "满足成长规则后自动升级"\n}', benefits: '[\n  "会员专享活动"\n]', status: 1, version: undefined }); }
function prettyJSON(value: string) { try { return JSON.stringify(JSON.parse(value), null, 2); } catch { return value; } }
function benefitText(value: string) { try { const list = JSON.parse(value); return Array.isArray(list) ? list.join('、') : '—'; } catch { return '—'; } }
async function load() { if (canManage.value) rows.value = (await listMemberLevels()).list || []; }
function edit(row?: MemberLevel) { editing.value = row; if (!row) reset(); else Object.assign(form, { name: row.name, rank: row.rank, rules: prettyJSON(row.rules), benefits: prettyJSON(row.benefits), status: row.status, version: row.version }); open.value = true; }
async function save() { if (!form.name.trim() || form.rank < 1) { ElMessage.warning('请填写等级名称与排序等级'); return; } try { JSON.parse(form.rules); const benefits = JSON.parse(form.benefits); if (!Array.isArray(benefits) || !benefits.length) throw new Error('benefits'); } catch { ElMessage.warning('等级规则必须是 JSON 对象，会员权益必须是非空 JSON 字符串数组'); return; } saving.value = true; try { const data = { ...form, name: form.name.trim(), rules: prettyJSON(form.rules), benefits: prettyJSON(form.benefits) }; if (editing.value) await updateMemberLevel(editing.value.id, data); else await createMemberLevel(data); ElMessage.success('会员等级已保存；不会修改现有用户等级或历史变更记录'); open.value = false; await load(); } finally { saving.value = false; } }
async function remove(row: MemberLevel) { try { await ElMessageBox.confirm(`删除“${row.name}”只会逻辑隐藏配置。若有用户正在使用或已有历史变更记录，服务端将拒绝删除。`, '删除会员等级', { type: 'warning' }); await deleteMemberLevel(row.id); ElMessage.success('会员等级已逻辑删除'); await load(); } catch {} }
onMounted(async () => { const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]); canManage.value = profile.roles.includes('platform') && codes.includes('user.member.level.manage'); await load(); });
</script>

<template>
  <Page title="会员等级管理" description="维护普通会员等级、成长规则和权益。配置变更不会直接升级、降级或删除任何用户等级；用户等级调整仍走独立的带审计命令。">
    <el-alert v-if="!canManage" title="当前账号没有会员等级管理权限" type="warning" :closable="false" />
	    <el-card v-else shadow="never"><template #header><div class="flex justify-between"><span>会员等级</span><el-button type="primary" @click="edit()">新增等级</el-button></div></template><el-table :data="rows"><el-table-column prop="name" label="等级名称" min-width="140"/><el-table-column prop="rank" label="等级排序" width="100"/><el-table-column label="权益" min-width="240"><template #default="{ row }">{{ benefitText(row.benefits) }}</template></el-table-column><el-table-column label="当前用户数" prop="assigned_count" width="120"/><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" :disabled="row.assigned_count > 0" @click="remove(row)">删除</el-button></template></el-table-column></el-table><el-alert class="mt-4" type="info" :closable="false" title="删除还会检查不可变会员变更日志；已使用过的等级不会被删除，以保护历史监管数据。"/></el-card>
    <el-dialog v-model="open" :title="editing ? '编辑会员等级' : '新增会员等级'" width="680px" destroy-on-close><el-form label-width="98px"><el-form-item label="等级名称" required><el-input v-model="form.name" maxlength="64"/></el-form-item><el-form-item label="等级排序" required><el-input-number v-model="form.rank" :min="1" :max="10000"/></el-form-item><el-form-item label="成长规则 JSON" required><el-input v-model="form.rules" type="textarea" :rows="6" class="font-mono"/></el-form-item><el-form-item label="会员权益 JSON" required><el-input v-model="form.benefits" type="textarea" :rows="6" class="font-mono" placeholder='["权益一", "权益二"]'/></el-form-item><el-form-item label="状态"><el-radio-group v-model="form.status"><el-radio :value="1">启用</el-radio><el-radio :value="0">停用</el-radio></el-radio-group></el-form-item></el-form><template #footer><el-button @click="open = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
