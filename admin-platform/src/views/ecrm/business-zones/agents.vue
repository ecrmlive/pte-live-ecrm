<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  createBusinessZoneAgent,
  fetchBusinessZoneAgents,
  updateBusinessZoneAgent,
  type BusinessZoneAgentRow,
} from '#/api/core/ecrm';

const rows = ref<BusinessZoneAgentRow[]>([]);
const total = ref(0);
const loading = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const query = reactive({ keyword: '', status: undefined as number | undefined, page: 1, limit: 20 });
const form = reactive({ uid: 0, name: '', phone: '', qualification: '', remark: '', payment_method: 0, payment_name: '', payment_account: '', payment_bank: '', payment_qr_img: '', type: 0, business_name: '', business_store_category: 0, business_store_type: 0 });

const statusText = (value: number) => ({ '-2': '已撤销', '-1': '已驳回', '0': '待审核', '1': '已通过' }[String(value)] || '未知');
const statusTag = (value: number) => ({ '-2': 'info', '-1': 'danger', '0': 'warning', '1': 'success' }[String(value)] || 'info');
const payText = (value: number) => ['银行卡', '微信', '支付宝'][value] || '银行卡';
const dialogTitle = computed(() => editingID.value ? '编辑代理申请' : '新增代理申请');

async function load() { loading.value = true; try { const result = await fetchBusinessZoneAgents({ ...query, keyword: query.keyword.trim() || undefined }); rows.value = result.list; total.value = result.total; } finally { loading.value = false; } }
function resetForm() { Object.assign(form, { uid: 0, name: '', phone: '', qualification: '', remark: '', payment_method: 0, payment_name: '', payment_account: '', payment_bank: '', payment_qr_img: '', type: 0, business_name: '', business_store_category: 0, business_store_type: 0 }); }
function add() { editingID.value = undefined; resetForm(); dialogOpen.value = true; }
function edit(row: BusinessZoneAgentRow) { if (row.status !== 0) { ElMessage.warning('已审核的代理资料不可编辑'); return; } editingID.value = row.circle_agent_id; Object.assign(form, { uid: row.uid, name: row.name, phone: row.phone, qualification: row.qualification, remark: row.remark, payment_method: row.payment_method, payment_name: row.payment_name, payment_account: row.payment_account, payment_bank: row.payment_bank, payment_qr_img: row.payment_qr_img, type: row.type, business_name: row.business_name, business_store_category: 0, business_store_type: 0 }); dialogOpen.value = true; }
async function save() { if (!form.name.trim() || !form.phone.trim()) { ElMessage.warning('代理姓名和手机号必填'); return; } if (editingID.value) await updateBusinessZoneAgent(editingID.value, form); else await createBusinessZoneAgent(form); ElMessage.success('已保存，新增申请需在代理审核中处理'); dialogOpen.value = false; await load(); }
function search() { query.page = 1; void load(); }
function reset() { query.keyword = ''; query.status = undefined; query.page = 1; void load(); }
onMounted(load);
</script>

<template>
  <Page title="代理人员" description="维护区域代理申请、结算资料和业务类型；审核通过后才能绑定至区域商圈。">
    <el-card shadow="never"><el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px" @submit.prevent="search"><el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="姓名/手机号/商户名" @keyup.enter="search" /></el-form-item><el-form-item label="审核状态"><el-select v-model="query.status" clearable class="w-full" placeholder="全部"><el-option label="待审核" :value="0" /><el-option label="已通过" :value="1" /><el-option label="已驳回" :value="-1" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button><el-button type="success" @click="add">新增代理</el-button></el-form-item></el-form>
      <el-table v-loading="loading" :data="rows" border><el-table-column prop="circle_agent_id" label="ID" width="72" /><el-table-column prop="name" label="代理姓名" min-width="110" /><el-table-column prop="phone" label="手机号" width="140" /><el-table-column prop="business_name" label="关联商户" min-width="140" /><el-table-column label="结算方式" width="100"><template #default="{ row }">{{ payText(row.payment_method) }}</template></el-table-column><el-table-column prop="balance" label="佣金余额" width="110" /><el-table-column label="审核状态" width="100"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column><el-table-column label="操作" width="90" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button></template></el-table-column></el-table>
      <div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" :page-sizes="[10,20,50,100]" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="dialogTitle" width="640px"><el-form label-width="110px"><el-form-item label="关联用户ID"><el-input-number v-model="form.uid" :min="0" /></el-form-item><el-form-item label="代理姓名" required><el-input v-model="form.name" /></el-form-item><el-form-item label="联系电话" required><el-input v-model="form.phone" /></el-form-item><el-form-item label="代理类型"><el-radio-group v-model="form.type"><el-radio :value="0">区域代理</el-radio><el-radio :value="1">商户型代理</el-radio></el-radio-group></el-form-item><el-form-item label="关联商户"><el-input v-model="form.business_name" placeholder="商户型代理填写" /></el-form-item><el-form-item label="结算方式"><el-select v-model="form.payment_method"><el-option label="银行卡" :value="0" /><el-option label="微信" :value="1" /><el-option label="支付宝" :value="2" /></el-select></el-form-item><el-form-item label="结算名称"><el-input v-model="form.payment_name" /></el-form-item><el-form-item label="结算账号"><el-input v-model="form.payment_account" /></el-form-item><el-form-item label="开户行"><el-input v-model="form.payment_bank" /></el-form-item><el-form-item label="资质材料"><el-input v-model="form.qualification" type="textarea" /></el-form-item><el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
