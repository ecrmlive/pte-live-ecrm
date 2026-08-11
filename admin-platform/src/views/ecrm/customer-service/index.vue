<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import {
  claimCustomerServiceThread,
  createCustomerServiceQuickReply,
  deleteCustomerServiceQuickReply,
  fetchCustomerServiceAgents,
  fetchCustomerServiceAgentUsers,
  fetchCustomerServiceAssignmentLogs,
  fetchCustomerServiceDelivery,
  fetchCustomerServiceEvents,
  fetchCustomerServiceOrder,
  fetchCustomerServiceProducts,
  fetchCustomerServiceQuickReplies,
  fetchCustomerServiceRefunds,
  fetchCustomerServiceThread,
  fetchCustomerServiceThreads,
  fetchCustomerServiceUser,
  transferCustomerServiceThread,
  updateCustomerServiceQuickReply,
  updateCustomerServiceUserNote,
  type CustomerServiceAgent,
  type CustomerServiceAgentUser,
  type CustomerServiceEvent,
  type CustomerServiceAssignmentLog,
  type CustomerServiceQuickReply,
  type CustomerServiceThread,
} from '#/api/core/customer-service';

const eventOpen = ref(false);
const eventLoading = ref(false);
const events = ref<CustomerServiceEvent[]>([]);
const assignmentLogs = ref<CustomerServiceAssignmentLog[]>([]);
const eventConversation = ref('');
const eventNote = ref('');
const supportOpen = ref(false);
const supportLoading = ref(false);
const supportOrder = ref<Record<string, unknown>>();
const supportDelivery = ref<Record<string, unknown>>();
const supportProducts = ref<Record<string, unknown>[]>([]);
const supportRefunds = ref<Record<string, unknown>[]>([]);
const supportUser = ref<Record<string, unknown>>();
const supportThreadID = ref<number>();
const supportNote = ref('');
const supportNoteSaving = ref(false);
const quickOpen = ref(false);
const quickEditOpen = ref(false);
const quickLoading = ref(false);
const quickSaving = ref(false);
const quickRows = ref<CustomerServiceQuickReply[]>([]);
const quickTotal = ref(0);
const editingQuickID = ref<number>();
const quickQuery = reactive({ page: 1, limit: 20, store_id: undefined as number | undefined });
const quickForm = reactive({ store_id: undefined as number | undefined, title: '', content: '', status: 'enabled' as 'disabled' | 'enabled' });
const transferOpen = ref(false);
const transferring = ref(false);
const transferThread = ref<CustomerServiceThread>();
const transferForm = reactive({ target_admin_id: undefined as number | undefined, reason: '' });
const serviceAgents = ref<CustomerServiceAgent[]>([]);
const agentsLoading = ref(false);
const agentRosterOpen = ref(false);
const agentUsersOpen = ref(false);
const agentUsersLoading = ref(false);
const agentUsers = ref<CustomerServiceAgentUser[]>([]);
const agentUsersTotal = ref(0);
const selectedAgent = ref<CustomerServiceAgent>();
const agentUserQuery = reactive({ page: 1, limit: 20 });
function time(value?: string | null) {
  if (!value) return '—';
  return formatShanghaiDateTime(value);
}

function assignment(row: CustomerServiceThread) {
  return row.assigned_admin_id ? `已领取（后台用户 #${row.assigned_admin_id}）` : '待领取';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待服务 / 服务中', value: 'open' },
        { label: '已关闭', value: 'closed' },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '会话状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '全部队列', value: 0 },
        { label: '仅看我的', value: 1 },
      ],
      placeholder: '全部队列',
    },
    fieldName: 'mine',
    label: '归属',
  },
]);

const gridOptions: VxeGridProps<CustomerServiceThread> = {
  columns: [
    { field: 'id', title: '会话编号', width: 96 },
    {
      field: 'store_name',
      formatter: ({ row }) => row.store_name || `店铺 #${row.store_id}`,
      minWidth: 150,
      title: '店铺',
    },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'order_id',
      formatter: ({ cellValue }) => cellValue || '—',
      title: '关联订单',
      width: 116,
    },
    {
      field: 'im_conversation_id',
      minWidth: 180,
      showOverflow: false,
      title: 'IM 会话 ID',
    },
    {
      field: 'im_sdk_app_id',
      formatter: ({ cellValue }) => cellValue || '未配置',
      minWidth: 150,
      title: '当前 IM SDK',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 104,
    },
    {
      field: 'assigned_admin_id',
      formatter: ({ row }) => assignment(row),
      minWidth: 164,
      title: '领取状态',
    },
    {
      field: 'updated_at',
      formatter: ({ cellValue }) => time(cellValue),
      minWidth: 172,
      title: '更新时间',
    },
    platformListActionColumn({ minWidth: 224 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const status = formValues?.status as 'closed' | 'open' | undefined;
        const mineRaw = formValues?.mine;
        const result = await fetchCustomerServiceThreads({
          page: page.currentPage,
          limit: page.pageSize,
          mine: mineRaw === 1,
          status: status || undefined,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function claim(row: CustomerServiceThread) {
  try {
    await ElMessageBox.confirm('领取后该会话归入“仅看我的”，是否继续？', '领取客服会话', {
      confirmButtonText: '领取',
      cancelButtonText: '取消',
      type: 'warning',
    });
    await claimCustomerServiceThread(row.id);
    ElMessage.success('客服会话已领取');
    gridApi.reload();
  } catch {
    // 请求客户端已经负责展示服务端错误；取消操作不提示。
  }
}

async function openTransfer(row: CustomerServiceThread) {
  transferThread.value = row;
  transferForm.target_admin_id = undefined;
  transferForm.reason = '';
  transferOpen.value = true;
  agentsLoading.value = true;
  try {
    serviceAgents.value = (await fetchCustomerServiceAgents({ limit: 100 })).list.filter((agent) => agent.status === 1 && (agent.service_store_ids || []).includes(row.store_id));
  } finally { agentsLoading.value = false; }
}

async function openAgentRoster() {
  agentRosterOpen.value = true;
  agentsLoading.value = true;
  try { serviceAgents.value = (await fetchCustomerServiceAgents({ limit: 100 })).list; } finally { agentsLoading.value = false; }
}

async function loadAgentUsers() {
  if (!selectedAgent.value) return;
  agentUsersLoading.value = true;
  try {
    const result = await fetchCustomerServiceAgentUsers(selectedAgent.value.id, agentUserQuery);
    agentUsers.value = result.list;
    agentUsersTotal.value = result.total;
  } finally { agentUsersLoading.value = false; }
}

function openAgentUsers(agent: CustomerServiceAgent) {
  selectedAgent.value = agent;
  agentUserQuery.page = 1;
  agentUsersOpen.value = true;
  void loadAgentUsers();
}

async function openAgentUserThread(row: CustomerServiceAgentUser) {
  const thread = await fetchCustomerServiceThread(row.binding_id);
  await openEvents(thread);
}

async function transfer() {
  if (!transferThread.value || !transferForm.target_admin_id || !transferForm.reason.trim()) {
    ElMessage.warning('请填写目标客服后台用户 ID 和转接原因');
    return;
  }
  transferring.value = true;
  try {
    await transferCustomerServiceThread(
      transferThread.value.id,
      { target_admin_id: transferForm.target_admin_id, reason: transferForm.reason },
      crypto.randomUUID(),
    );
    transferOpen.value = false;
    gridApi.reload();
    ElMessage.success('客服会话已转接并写入审计记录');
  } finally { transferring.value = false; }
}

async function openEvents(row: CustomerServiceThread) {
  eventOpen.value = true;
  eventLoading.value = true;
  events.value = [];
  assignmentLogs.value = [];
  eventConversation.value = row.im_conversation_id;
  try {
    const [eventResult, assignmentResult] = await Promise.all([
      fetchCustomerServiceEvents(row.id, { page: 1, limit: 100 }),
      fetchCustomerServiceAssignmentLogs(row.id, { page: 1, limit: 100 }),
    ]);
    events.value = eventResult.list;
    assignmentLogs.value = assignmentResult.list;
    eventConversation.value = eventResult.conversation_id;
    eventNote.value = eventResult.note;
  } finally { eventLoading.value = false; }
}

async function openSupport(row: CustomerServiceThread) {
  supportOpen.value = true;
  supportLoading.value = true;
  supportOrder.value = undefined;
  supportDelivery.value = undefined;
  supportProducts.value = [];
  supportRefunds.value = [];
  supportUser.value = undefined;
  supportThreadID.value = row.id;
  supportNote.value = '';
  try {
    const [order, delivery, products, refunds, user] = await Promise.allSettled([fetchCustomerServiceOrder(row.id), fetchCustomerServiceDelivery(row.id), fetchCustomerServiceProducts(row.id), fetchCustomerServiceRefunds(row.id), fetchCustomerServiceUser(row.id)]);
    if (order.status === 'fulfilled') supportOrder.value = order.value;
    if (delivery.status === 'fulfilled') supportDelivery.value = delivery.value;
    if (products.status === 'fulfilled') supportProducts.value = products.value.list;
    if (refunds.status === 'fulfilled') supportRefunds.value = refunds.value.list;
    if (user.status === 'fulfilled') {
      supportUser.value = user.value;
      supportNote.value = typeof user.value.service_note === 'string' ? user.value.service_note : '';
    }
  } finally { supportLoading.value = false; }
}

async function saveSupportNote() {
  if (!supportThreadID.value || !supportNote.value.trim()) {
    ElMessage.warning('请输入用户备注');
    return;
  }
  supportNoteSaving.value = true;
  try {
    const result = await updateCustomerServiceUserNote(supportThreadID.value, supportNote.value);
    supportNote.value = result.content;
    if (supportUser.value) supportUser.value.service_note = result.content;
    ElMessage.success('用户备注已保存');
  } finally { supportNoteSaving.value = false; }
}

async function loadQuickReplies() {
  quickLoading.value = true;
  try {
    const result = await fetchCustomerServiceQuickReplies(quickQuery);
    quickRows.value = result.list;
    quickTotal.value = result.total;
  } finally { quickLoading.value = false; }
}

function openQuickReplies() {
  quickOpen.value = true;
  void loadQuickReplies();
}

function newQuickReply() {
  editingQuickID.value = undefined;
  Object.assign(quickForm, { store_id: quickQuery.store_id, title: '', content: '', status: 'enabled' });
  quickEditOpen.value = true;
}

function editQuickReply(row: CustomerServiceQuickReply) {
  editingQuickID.value = row.id;
  Object.assign(quickForm, { store_id: row.store_id, title: row.title, content: row.content, status: row.status });
  quickEditOpen.value = true;
}

async function saveQuickReply() {
  if (!quickForm.store_id || !quickForm.title.trim() || !quickForm.content.trim()) {
    ElMessage.warning('请填写授权店铺 ID、标题和回复内容');
    return;
  }
  quickSaving.value = true;
  try {
    if (editingQuickID.value) {
      await updateCustomerServiceQuickReply(editingQuickID.value, { title: quickForm.title, content: quickForm.content, status: quickForm.status });
    } else {
      await createCustomerServiceQuickReply({ store_id: quickForm.store_id, title: quickForm.title, content: quickForm.content, status: quickForm.status });
    }
    ElMessage.success(editingQuickID.value ? '快捷回复已更新' : '快捷回复已创建');
    quickEditOpen.value = false;
    await loadQuickReplies();
  } finally { quickSaving.value = false; }
}

async function removeQuickReply(row: CustomerServiceQuickReply) {
  try {
    await ElMessageBox.confirm(`确定删除“${row.title}”吗？删除后将保留审计记录且不再显示。`, '删除快捷回复', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' });
    await deleteCustomerServiceQuickReply(row.id);
    ElMessage.success('快捷回复已删除');
    await loadQuickReplies();
  } catch {
    // 取消不提示；请求客户端负责展示服务端权限或业务错误。
  }
}

</script>

<template>
  <Page auto-content-height>
    <ElAlert
      class="mb-4"
      type="info"
      :closable="false"
      title="仅展示已获授权店铺的会话；领取不会改变 IM 消息内容或 IM SDK AppID。"
    />
    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" @click="openQuickReplies">快捷回复管理</ElButton>
        <ElButton @click="openAgentRoster">客服人员与用户</ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 'open' ? 'warning' : 'info'">
          {{ row.status === 'open' ? '进行中' : '已关闭' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openSupport(row)">订单辅助</ElButton>
        <ElButton link type="primary" @click="openEvents(row)">事件</ElButton>
        <ElButton
          v-if="row.status === 'open' && !row.assigned_admin_id"
          link
          type="primary"
          @click="claim(row)"
        >
          领取
        </ElButton>
        <ElButton
          v-if="row.status === 'open' && row.assigned_admin_id"
          link
          type="warning"
          @click="openTransfer(row)"
        >
          转接
        </ElButton>
      </template>
    </Grid>
    <el-dialog v-model="transferOpen" title="转接客服会话" width="520px" destroy-on-close>
      <el-alert class="mb-4" type="warning" :closable="false" title="仅当前领取客服或平台可转接；目标客服必须启用且拥有同一店铺服务范围。" />
      <el-form label-width="120px"><el-form-item label="目标客服" required><el-select v-model="transferForm.target_admin_id" v-loading="agentsLoading" class="w-full" filterable placeholder="选择拥有本店铺服务范围的启用客服"><el-option v-for="agent in serviceAgents" :key="agent.id" :label="`${agent.display_name || `客服 #${agent.id}`}（ID ${agent.id}）`" :value="agent.id"><span>{{ agent.display_name || `客服 #${agent.id}` }}</span><span class="float-right text-gray-400">店铺 {{ agent.service_store_ids.join(',') }}</span></el-option></el-select></el-form-item><el-form-item label="转接原因" required><el-input v-model="transferForm.reason" :rows="3" maxlength="500" show-word-limit type="textarea" placeholder="例如：虚构演示，转交对应商品咨询队列" /></el-form-item></el-form>
      <template #footer><el-button @click="transferOpen = false">取消</el-button><el-button :loading="transferring" type="primary" @click="transfer">确认转接</el-button></template>
    </el-dialog>
    <el-dialog v-model="agentRosterOpen" title="客服人员与服务用户" width="760px" destroy-on-close><el-alert class="mb-3" type="info" :closable="false" title="仅展示当前授权店铺可见的客服人员；服务用户列表仅返回脱敏手机号和已分配会话。" /><el-table v-loading="agentsLoading" :data="serviceAgents" max-height="420"><el-table-column prop="id" label="客服 ID" width="100" /><el-table-column prop="display_name" label="姓名" min-width="150"><template #default="{ row }">{{ row.display_name || `客服 #${row.id}` }}</template></el-table-column><el-table-column label="服务店铺" min-width="180"><template #default="{ row }">{{ (row.service_store_ids || []).join(',') || '—' }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="120"><template #default="{ row }"><el-button link type="primary" @click="openAgentUsers(row)">服务用户</el-button></template></el-table-column></el-table></el-dialog>
    <el-dialog v-model="agentUsersOpen" :title="`${selectedAgent?.display_name || '客服'}的服务用户`" width="820px" destroy-on-close><el-table v-loading="agentUsersLoading" :data="agentUsers" max-height="420"><el-table-column prop="binding_id" label="会话编号" width="96" /><el-table-column prop="user_id" label="用户 ID" width="96" /><el-table-column prop="nickname" label="用户昵称" min-width="150"><template #default="{ row }">{{ row.nickname || '未设置昵称' }}</template></el-table-column><el-table-column prop="mobile" label="手机号" width="130"><template #default="{ row }">{{ row.mobile || '—' }}</template></el-table-column><el-table-column label="店铺" min-width="150"><template #default="{ row }">{{ row.store_name || `店铺 #${row.store_id}` }}</template></el-table-column><el-table-column label="会话状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 'open' ? 'warning' : 'info'">{{ row.status === 'open' ? '进行中' : '已关闭' }}</el-tag></template></el-table-column><el-table-column label="操作" width="110"><template #default="{ row }"><el-button link type="primary" @click="openAgentUserThread(row)">事件记录</el-button></template></el-table-column></el-table><div class="mt-3 flex justify-end"><el-pagination :current-page="agentUserQuery.page" :page-size="agentUserQuery.limit" :total="agentUsersTotal" background layout="total, prev, pager, next" @current-change="(page) => { agentUserQuery.page = page; loadAgentUsers(); }" /></div></el-dialog>
    <el-dialog v-model="eventOpen" title="客服会话记录" width="760px" destroy-on-close>
      <el-alert class="mb-3" type="info" :closable="false" :title="eventNote || '聊天正文由 pte-live-im 提供'" />
      <div class="mb-3 text-sm text-gray-500">IM 会话：{{ eventConversation || '—' }}</div>
      <el-tabs v-loading="eventLoading">
        <el-tab-pane label="订单 / 系统事件">
          <el-table :data="events" max-height="360"><el-table-column prop="id" label="事件 ID" width="90" /><el-table-column prop="msg_type" label="类型" width="90" /><el-table-column prop="sender_role" label="来源" width="90" /><el-table-column prop="content" label="内容" min-width="220" show-overflow-tooltip /><el-table-column label="时间" min-width="170"><template #default="{ row }">{{ time(row.created_at) }}</template></el-table-column></el-table>
        </el-tab-pane>
        <el-tab-pane label="转接审计">
          <el-table :data="assignmentLogs" empty-text="暂无转接记录" max-height="360"><el-table-column prop="id" label="审计 ID" width="90" /><el-table-column label="原客服" width="106"><template #default="{ row }">{{ row.from_admin_id ? `#${row.from_admin_id}` : '未领取' }}</template></el-table-column><el-table-column label="目标客服" width="106"><template #default="{ row }">#{{ row.target_admin_id }}</template></el-table-column><el-table-column label="操作人" width="96"><template #default="{ row }">#{{ row.operator_admin_id }}</template></el-table-column><el-table-column prop="reason" label="转接原因" min-width="230" show-overflow-tooltip /><el-table-column label="时间" min-width="170"><template #default="{ row }">{{ time(row.created_at) }}</template></el-table-column></el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
    <el-dialog v-model="supportOpen" title="订单与售后辅助" width="760px" destroy-on-close>
      <div v-loading="supportLoading">
        <el-descriptions v-if="supportOrder" :column="2" border title="关联订单">
          <el-descriptions-item label="订单号">{{ supportOrder.order_no }}</el-descriptions-item><el-descriptions-item label="状态">{{ supportOrder.status }}</el-descriptions-item><el-descriptions-item label="实付金额">{{ supportOrder.pay_amount }}</el-descriptions-item><el-descriptions-item label="商品数量">{{ supportOrder.total_quantity }}</el-descriptions-item>
        </el-descriptions>
        <el-descriptions v-if="supportDelivery" class="mt-4" :column="2" border title="订单配送">
          <el-descriptions-item label="配送方式">{{ supportDelivery.delivery_type || '—' }}</el-descriptions-item><el-descriptions-item label="配送状态">{{ supportDelivery.status || '—' }}</el-descriptions-item><el-descriptions-item label="承运商">{{ supportDelivery.carrier_code || '—' }}</el-descriptions-item><el-descriptions-item label="运单号">{{ supportDelivery.tracking_no || '—' }}</el-descriptions-item>
        </el-descriptions>
        <el-table class="mt-4" :data="supportProducts" empty-text="暂无关联商品">
          <el-table-column prop="product_id" label="商品 ID" width="96" /><el-table-column prop="title_snapshot" label="下单商品" min-width="190" /><el-table-column prop="sku_key" label="规格" min-width="120" /><el-table-column prop="quantity" label="数量" width="80" /><el-table-column prop="unit_price" label="下单单价" width="110" /><el-table-column prop="current_stock" label="当前库存" width="96" />
        </el-table>
        <el-descriptions v-if="supportUser" class="mt-4" :column="2" border title="用户资料">
          <el-descriptions-item label="昵称">{{ supportUser.nickname || '—' }}</el-descriptions-item><el-descriptions-item label="手机号">{{ supportUser.mobile || '—' }}</el-descriptions-item><el-descriptions-item label="来源">{{ supportUser.source_channel || '—' }}</el-descriptions-item><el-descriptions-item label="简介">{{ supportUser.bio || '—' }}</el-descriptions-item>
        </el-descriptions>
        <el-form v-if="supportUser" class="mt-4" label-width="76px"><el-form-item label="用户备注"><el-input v-model="supportNote" :rows="3" maxlength="500" placeholder="仅对当前授权店铺的该用户可见" show-word-limit type="textarea" /></el-form-item><el-form-item><el-button :loading="supportNoteSaving" type="primary" @click="saveSupportNote">保存备注</el-button></el-form-item></el-form>
        <el-table class="mt-4" :data="supportRefunds" empty-text="暂无关联退款"><el-table-column prop="refund_no" label="退款单号" /><el-table-column prop="amount" label="金额" width="120" /><el-table-column prop="status" label="状态" width="150" /><el-table-column prop="reason" label="原因" min-width="180" show-overflow-tooltip /></el-table>
      </div>
    </el-dialog>
    <el-dialog v-model="quickOpen" title="快捷回复管理" width="820px" destroy-on-close><el-form inline class="mb-3"><el-form-item label="店铺 ID"><el-input-number v-model="quickQuery.store_id" :min="1" controls-position="right" clearable /></el-form-item><el-form-item><el-button type="primary" @click="quickQuery.page = 1; loadQuickReplies()">查询</el-button><el-button @click="quickQuery.store_id = undefined; quickQuery.page = 1; loadQuickReplies()">重置</el-button><el-button type="success" @click="newQuickReply">新增快捷回复</el-button></el-form-item></el-form><el-alert class="mb-3" type="info" :closable="false" title="仅可管理服务端已授权店铺的数据；删除为软删除并保留审计记录。" /><el-table v-loading="quickLoading" :data="quickRows" max-height="400"><el-table-column prop="store_id" label="店铺 ID" width="92" /><el-table-column prop="title" label="标题" min-width="130" /><el-table-column prop="content" label="回复内容" min-width="260" show-overflow-tooltip /><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 'enabled' ? 'success' : 'info'">{{ row.status === 'enabled' ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="更新时间" min-width="160"><template #default="{ row }">{{ time(row.updated_at) }}</template></el-table-column><el-table-column label="操作" fixed="right" width="120"><template #default="{ row }"><el-button link type="primary" @click="editQuickReply(row)">编辑</el-button><el-button link type="danger" @click="removeQuickReply(row)">删除</el-button></template></el-table-column></el-table><div class="mt-3 flex justify-end"><el-pagination :current-page="quickQuery.page" :page-size="quickQuery.limit" :total="quickTotal" background layout="total, prev, pager, next" @current-change="(page) => { quickQuery.page = page; loadQuickReplies(); }" /></div></el-dialog>
    <el-dialog v-model="quickEditOpen" :title="editingQuickID ? '编辑快捷回复' : '新增快捷回复'" width="560px" destroy-on-close><el-form label-width="96px"><el-form-item label="店铺 ID" required><el-input-number v-model="quickForm.store_id" :disabled="!!editingQuickID" :min="1" /></el-form-item><el-form-item label="标题" required><el-input v-model="quickForm.title" maxlength="64" show-word-limit /></el-form-item><el-form-item label="回复内容" required><el-input v-model="quickForm.content" :rows="5" maxlength="2000" show-word-limit type="textarea" /></el-form-item><el-form-item label="状态"><el-radio-group v-model="quickForm.status"><el-radio value="enabled">启用</el-radio><el-radio value="disabled">停用</el-radio></el-radio-group></el-form-item></el-form><template #footer><el-button @click="quickEditOpen = false">取消</el-button><el-button :loading="quickSaving" type="primary" @click="saveQuickReply">保存</el-button></template></el-dialog>
  </Page>
</template>
