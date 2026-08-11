<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElAlert,
  ElAvatar,
  ElButton,
  ElDatePicker,
  ElDescriptions,
  ElDescriptionsItem,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSkeleton,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  adjustPlatformUserAsset,
  adjustPlatformUserMemberLevel,
  assignPlatformUserGroup,
  assignPlatformUserGroups,
  assignPlatformUserLabel,
  assignPlatformUserLabels,
  changePlatformUserReferrer,
  createPlatformUser,
  fetchPlatformCouponTemplates,
  fetchPlatformMemberLevels,
  fetchPlatformUserDetail,
  fetchPlatformUserGroupOptions,
  fetchPlatformUserLabelOptions,
  fetchPlatformUsers,
  issuePlatformUserCoupon,
  resetPlatformUserPassword,
  setPlatformUserListSvip,
  updatePlatformUserProfile,
  type PlatformCouponTemplate,
  type PlatformMemberLevel,
  type PlatformUserDetail,
  type PlatformUserGroupOption,
  type PlatformUserLabelOption,
  type PlatformUserRow,
} from '#/api/core/ecrm';
import UserPickerModal, {
  type PickedPlatformUser,
} from '#/components/ecrm/user-picker-modal.vue';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import {
  listPrefixedKeywordFormField,
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';

type ChannelTab = '' | 'wechat' | 'mini_program' | 'h5' | 'app' | 'pc';
type MoreCommand =
  | 'balance'
  | 'points'
  | 'group'
  | 'label'
  | 'referrer'
  | 'password'
  | 'level'
  | 'coupon'
  | 'svip';

const CHANNEL_TABS: Array<{ key: ChannelTab; label: string }> = [
  { key: '', label: '全部用户' },
  { key: 'wechat', label: '微信' },
  { key: 'mini_program', label: '小程序' },
  { key: 'h5', label: 'H5' },
  { key: 'app', label: 'APP' },
  { key: 'pc', label: 'PC' },
];

const CHANNEL_LABELS: Record<string, string> = {
  wechat: '微信',
  mini_program: '小程序',
  h5: 'H5',
  pc: 'PC',
  ios: 'APP',
  android: 'APP',
  harmony: 'APP',
};

const assetLabels: Record<string, string> = {
  balance: '余额',
  points: '积分',
  commission: '佣金',
};
const orderLabels: Record<string, string> = {
  pending_pay: '待支付',
  paid: '已支付',
  awaiting_final: '待尾款',
  final_timeout: '尾款超时',
  fulfilling: '履约中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消',
  aftersale: '售后中',
};
const membershipLabels: Record<string, string> = {
  initial: '初始',
  upgrade: '升级',
  downgrade: '降级',
  manual: '人工调整',
};
const couponStatusLabels: Record<string, string> = {
  unused: '未使用',
  locked: '锁定',
  used: '已使用',
  expired: '已过期',
};

const detail = ref<PlatformUserDetail>();
const detailLoading = ref(false);
const canRead = ref(false);
const canCreate = ref(false);
const canProfile = ref(false);
const canAsset = ref(false);
const canGroup = ref(false);
const canLabel = ref(false);
const canReferrer = ref(false);
const canPassword = ref(false);
const canMember = ref(false);
const canCoupon = ref(false);
const canSvip = ref(false);
const isPlatform = ref(false);
const channelTab = ref<ChannelTab>('');
const selectedIds = ref<number[]>([]);
const labelOptions = ref<PlatformUserLabelOption[]>([]);
const groupOptions = ref<PlatformUserGroupOption[]>([]);
const levelOptions = ref<PlatformMemberLevel[]>([]);
const couponOptions = ref<PlatformCouponTemplate[]>([]);
const actionUser = ref<PlatformUserRow>();
const actionUserIds = ref<number[]>([]);

const assetForm = reactive({
  direction: 'inc' as 'inc' | 'dec',
  amount: undefined as number | undefined,
  reason: '',
});
const groupForm = reactive({ group_id: 0, reason: '' });
const labelForm = reactive({ label_ids: [] as number[], reason: '' });
const referrerForm = reactive({
  parent_user_id: undefined as number | undefined,
  reason: '',
});
const referrerParent = ref<PickedPlatformUser | null>(null);
const referrerPickerOpen = ref(false);
const passwordForm = reactive({
  password: '',
  confirm_password: '',
  reason: '',
});
const levelForm = reactive({ level_id: 0, reason: '' });
const couponForm = reactive({
  coupon_id: undefined as number | undefined,
  reason: '',
});
const svipForm = reactive({ is_svip: 0, svip_endtime: '' });
const createForm = reactive({
  account: '',
  password: '',
  confirm_password: '',
  nickname: '',
  avatar_url: '',
  real_name: '',
  phone: '',
  id_card: '',
  gender: 0 as 0 | 1 | 2,
  status: 1 as 0 | 1,
  is_promoter: 1 as 0 | 1,
});
const editForm = reactive({
  nickname: '',
  avatar_url: '',
  gender: 0 as 0 | 1 | 2,
  bio: '',
  reason: '',
});

const hasMoreActions = computed(
  () =>
    canAsset.value ||
    canGroup.value ||
    canLabel.value ||
    canReferrer.value ||
    canPassword.value ||
    canMember.value ||
    canCoupon.value ||
    canSvip.value,
);

const USER_LIST_SEARCH_OPTIONS = [
  { label: '全部', value: 'all' },
  { label: '昵称', value: 'nickname' },
  { label: '用户ID', value: 'uid' },
  { label: '手机号', value: 'phone' },
];

const FIELDS_SEARCH_OPTIONS = [
  { label: '姓名', value: 'real_name' },
  { label: '地址', value: 'address' },
  { label: '备注', value: 'mark' },
  { label: '身份证（实名认证）', value: 'id_card' },
];

const formOptions: VbenFormProps = listFormOptionsDefaults([
  listUserSearchFormField({
    defaultType: 'all',
    options: USER_LIST_SEARCH_OPTIONS,
    typeWidth: '96px',
  }),
  listPrefixedKeywordFormField({
    fieldName: 'fields_search',
    label: '信息补充',
    defaultType: 'real_name',
    options: FIELDS_SEARCH_OPTIONS,
    typeWidth: '150px',
  }),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'label_id',
    label: '用户标签',
  },
]);

const gridOptions: VxeGridProps<PlatformUserRow> = {
  checkboxConfig: { highlight: true, reserve: true },
  columns: [
    { type: 'checkbox', width: 48 },
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'avatar_url',
      slots: { default: 'avatar' },
      title: '头像',
      width: 72,
    },
    { field: 'nickname', minWidth: 140, showOverflow: false, title: '昵称' },
    {
      field: 'svip_label',
      minWidth: 110,
      showOverflow: false,
      title: '付费会员',
      formatter: ({ row }) => row.svip_label || '非会员',
    },
    { field: 'mobile', title: '手机号', width: 130 },
    {
      field: 'level_name',
      minWidth: 110,
      showOverflow: false,
      title: '等级',
      formatter: ({ cellValue }) => cellValue || '普通会员',
    },
    {
      field: 'group_name',
      minWidth: 110,
      showOverflow: false,
      title: '分组',
      formatter: ({ cellValue }) => cellValue || '未分组',
    },
    {
      field: 'parent_nickname',
      minWidth: 140,
      showOverflow: false,
      title: '推荐人',
      formatter: ({ row }) =>
        row.parent_user_id
          ? `${row.parent_nickname || '未知'} / ${row.parent_user_id}`
          : '—',
    },
    {
      field: 'source_channel',
      title: '用户类型',
      width: 100,
      formatter: ({ cellValue }) =>
        CHANNEL_LABELS[String(cellValue || '')] || '—',
    },
    {
      field: 'balance',
      title: '余额',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'points', title: '当前可用积分', width: 120 },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!isPlatform.value) {
          return { items: [], total: 0 };
        }
        const userSearch = parseUserSearch(formValues);
        const keyword = userSearch.keyword;
        const keywordType = userSearch.type || 'all';
        const fieldsSearch = parseUserSearch(formValues, 'fields_search');
        const fieldsType = fieldsSearch.type || 'real_name';
        const fieldsValue = fieldsSearch.keyword;
        const params: Parameters<typeof fetchPlatformUsers>[0] = {
          page: page.currentPage,
          limit: page.pageSize,
          label_id: formValues?.label_id
            ? Number(formValues.label_id)
            : undefined,
          fields_type: fieldsValue ? fieldsType : undefined,
          fields_value: fieldsValue || undefined,
          source_channel: channelTab.value || undefined,
        };
        if (keyword) {
          switch (keywordType) {
            case 'nickname':
              params.nickname = keyword;
              break;
            case 'phone':
              params.phone = keyword;
              break;
            case 'uid': {
              const uid = Number(keyword);
              if (!Number.isInteger(uid) || uid <= 0) {
                return { items: [], total: 0 };
              }
              params.id = uid;
              break;
            }
            default:
              params.keyword = keyword;
          }
        }
        const result = await fetchPlatformUsers(params);
        return { items: result.list || [], total: result.total || 0 };
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    checkboxAll: () => syncSelection(),
    checkboxChange: () => syncSelection(),
  },
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[960px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [CreateDrawer, createDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitCreate(),
});

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitEdit(),
});

const [AssetDrawer, assetDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitAsset(),
});

const [GroupDrawer, groupDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitGroup(),
});

const [LabelDrawer, labelDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitLabel(),
});

const [ReferrerDrawer, referrerDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitReferrer(),
});

const [PasswordDrawer, passwordDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitPassword(),
});

const [LevelDrawer, levelDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitLevel(),
});

const [CouponDrawer, couponDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitCoupon(),
});

const [SvipDrawer, svipDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submitSvip(),
});

const assetKind = ref<'balance' | 'points'>('balance');

function syncSelection() {
  const rows =
    (gridApi.grid?.getCheckboxRecords?.() as PlatformUserRow[] | undefined) ||
    [];
  selectedIds.value = rows.map((row) => row.id);
}

function channelLabel(channel?: string) {
  return CHANNEL_LABELS[String(channel || '')] || '—';
}

function setChannelTab(tab: ChannelTab) {
  channelTab.value = tab;
  gridApi.reload();
}

async function openDetail(row: PlatformUserRow) {
  detail.value = undefined;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: '用户详情', loading: true }).open();
  try {
    detail.value = await fetchPlatformUserDetail(row.id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function openEdit(row: PlatformUserRow) {
  actionUser.value = row;
  editDrawerApi.setState({ title: '编辑用户', loading: true }).open();
  try {
    const data = await fetchPlatformUserDetail(row.id);
    Object.assign(editForm, {
      nickname: data.profile.nickname || '',
      avatar_url: data.profile.avatar_url || '',
      gender: (data.profile.gender as 0 | 1 | 2) || 0,
      bio: data.profile.bio || '',
      reason: '平台后台编辑用户资料',
    });
  } finally {
    editDrawerApi.setState({ loading: false });
  }
}

function openCreate() {
  Object.assign(createForm, {
    account: '',
    password: '',
    confirm_password: '',
    nickname: '',
    avatar_url: '',
    real_name: '',
    phone: '',
    id_card: '',
    gender: 0,
    status: 1,
    is_promoter: 1,
  });
  createDrawerApi.setState({ title: '用户信息填写' }).open();
}

async function ensureGroups() {
  if (!groupOptions.value.length) {
    groupOptions.value = (await fetchPlatformUserGroupOptions()).list || [];
  }
}

async function ensureLevels() {
  if (!levelOptions.value.length) {
    levelOptions.value = (await fetchPlatformMemberLevels()).list || [];
  }
}

async function ensureCoupons() {
  if (!couponOptions.value.length) {
    couponOptions.value = (await fetchPlatformCouponTemplates()).list || [];
  }
}

async function onMoreCommand(cmd: MoreCommand, row: PlatformUserRow) {
  actionUser.value = row;
  actionUserIds.value = [row.id];
  switch (cmd) {
    case 'balance':
    case 'points': {
      assetKind.value = cmd;
      Object.assign(assetForm, {
        direction: 'inc',
        amount: undefined,
        reason: cmd === 'balance' ? '平台后台调整余额' : '平台后台调整积分',
      });
      assetDrawerApi
        .setState({
          title: cmd === 'balance' ? '修改用户余额' : '修改用户积分',
        })
        .open();
      break;
    }
    case 'group': {
      await ensureGroups();
      Object.assign(groupForm, {
        group_id: row.group_id || 0,
        reason: '平台后台设置分组',
      });
      groupDrawerApi.setState({ title: '修改用户分组' }).open();
      break;
    }
    case 'label': {
      Object.assign(labelForm, {
        label_ids: [],
        reason: '平台后台设置标签',
      });
      labelDrawerApi.setState({ title: '设置用户标签' }).open();
      break;
    }
    case 'referrer': {
      Object.assign(referrerForm, {
        parent_user_id: row.parent_user_id || undefined,
        reason: '平台后台修改推荐人',
      });
      referrerParent.value = row.parent_user_id
        ? {
            id: row.parent_user_id,
            nickname: row.parent_nickname || `用户#${row.parent_user_id}`,
            avatar_url: '',
            mobile: '',
          }
        : null;
      referrerDrawerApi.setState({ title: '修改推荐人' }).open();
      break;
    }
    case 'password': {
      Object.assign(passwordForm, {
        password: '',
        confirm_password: '',
        reason: '平台后台修改密码',
      });
      passwordDrawerApi.setState({ title: '修改密码' }).open();
      break;
    }
    case 'level': {
      await ensureLevels();
      Object.assign(levelForm, { level_id: 0, reason: '平台后台调整等级' });
      levelDrawerApi.setState({ title: '调整等级' }).open();
      break;
    }
    case 'coupon': {
      await ensureCoupons();
      Object.assign(couponForm, {
        coupon_id: undefined,
        reason: '平台后台发送优惠券',
      });
      couponDrawerApi.setState({ title: '发送优惠券' }).open();
      break;
    }
    case 'svip': {
      Object.assign(svipForm, {
        is_svip: row.is_svip || 0,
        svip_endtime: row.svip_expires_at
          ? formatShanghaiDateTime(row.svip_expires_at)
          : '',
      });
      svipDrawerApi.setState({ title: '付费会员设置' }).open();
      break;
    }
  }
}

async function openBatchGroup() {
  if (!selectedIds.value.length) return;
  await ensureGroups();
  actionUser.value = undefined;
  actionUserIds.value = [...selectedIds.value];
  Object.assign(groupForm, { group_id: 0, reason: '平台后台批量设置分组' });
  groupDrawerApi.setState({ title: '批量设置分组' }).open();
}

async function openBatchLabel() {
  if (!selectedIds.value.length) return;
  actionUser.value = undefined;
  actionUserIds.value = [...selectedIds.value];
  Object.assign(labelForm, {
    label_ids: [],
    reason: '平台后台批量设置标签',
  });
  labelDrawerApi.setState({ title: '批量设置标签' }).open();
}

async function submitCreate() {
  const account = createForm.account.trim();
  const nickname = createForm.nickname.trim();
  const password = createForm.password;
  if (account.length < 3 || account.length > 191) {
    ElMessage.warning('请填写手机号(账号)');
    return;
  }
  if (password.length < 12 || password.length > 72) {
    ElMessage.warning('请填写 12～72 位登录密码');
    return;
  }
  if (password !== createForm.confirm_password) {
    ElMessage.warning('两次输入的密码不一致');
    return;
  }
  if (nickname.length > 64) {
    ElMessage.warning('用户昵称不能超过 64 字');
    return;
  }
  createDrawerApi.lock();
  try {
    const result = await createPlatformUser({
      account,
      password,
      nickname,
      avatar_url: createForm.avatar_url.trim(),
      real_name: createForm.real_name.trim(),
      phone: createForm.phone.trim(),
      id_card: createForm.id_card.trim(),
      gender: createForm.gender,
      status: createForm.status,
      is_promoter: createForm.is_promoter,
      reason: '平台后台创建用户',
      idempotency_key: `user-create-${crypto.randomUUID()}`,
    });
    ElMessage.success(`用户 #${result.user_id} 已创建`);
    createDrawerApi.close();
    gridApi.reload();
  } finally {
    createDrawerApi.unlock();
  }
}

async function submitEdit() {
  if (!actionUser.value) return;
  const nickname = editForm.nickname.trim();
  const reason = editForm.reason.trim() || '平台后台编辑用户资料';
  if (nickname.length < 1 || nickname.length > 64 || reason.length < 2) {
    ElMessage.warning('请填写昵称与调整原因');
    return;
  }
  editDrawerApi.lock();
  try {
    await updatePlatformUserProfile(actionUser.value.id, {
      nickname,
      avatar_url: editForm.avatar_url.trim(),
      gender: editForm.gender,
      bio: editForm.bio.trim(),
      reason,
      idempotency_key: `user-profile-${actionUser.value.id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('用户资料已保存');
    editDrawerApi.close();
    gridApi.reload();
  } finally {
    editDrawerApi.unlock();
  }
}

async function submitAsset() {
  if (!actionUser.value || !assetForm.amount || assetForm.amount <= 0) {
    ElMessage.warning('请填写大于 0 的金额');
    return;
  }
  if (assetKind.value === 'points' && !Number.isInteger(assetForm.amount)) {
    ElMessage.warning('积分调整必须为整数');
    return;
  }
  const reason = assetForm.reason.trim() || '平台后台调整资产';
  const signed =
    assetForm.direction === 'inc' ? assetForm.amount : -assetForm.amount;
  assetDrawerApi.lock();
  try {
    await adjustPlatformUserAsset(actionUser.value.id, {
      asset_type: assetKind.value,
      amount: signed,
      reason,
      idempotency_key: `user-${assetKind.value}-${actionUser.value.id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('调整已生效');
    assetDrawerApi.close();
    gridApi.reload();
  } finally {
    assetDrawerApi.unlock();
  }
}

async function submitGroup() {
  const reason = groupForm.reason.trim() || '平台后台设置分组';
  const ids = actionUserIds.value;
  if (!ids.length || reason.length < 2) {
    ElMessage.warning('请选择用户并填写原因');
    return;
  }
  groupDrawerApi.lock();
  try {
    if (ids.length === 1 && ids[0]) {
      await assignPlatformUserGroup(ids[0], {
        group_id: groupForm.group_id,
        reason,
        idempotency_key: `user-group-${ids[0]}-${crypto.randomUUID()}`,
      });
    } else {
      await assignPlatformUserGroups({
        user_ids: ids,
        group_id: groupForm.group_id,
        reason,
        idempotency_key: `user-group-batch-${crypto.randomUUID()}`,
      });
    }
    ElMessage.success('分组已更新');
    groupDrawerApi.close();
    gridApi.reload();
  } finally {
    groupDrawerApi.unlock();
  }
}

async function submitLabel() {
  const reason = labelForm.reason.trim() || '平台后台设置标签';
  const ids = actionUserIds.value;
  if (!ids.length || reason.length < 2) {
    ElMessage.warning('请选择用户并填写原因');
    return;
  }
  labelDrawerApi.lock();
  try {
    if (ids.length === 1 && ids[0]) {
      await assignPlatformUserLabel(ids[0], {
        label_ids: labelForm.label_ids,
        reason,
        idempotency_key: `user-label-${ids[0]}-${crypto.randomUUID()}`,
      });
    } else {
      await assignPlatformUserLabels({
        user_ids: ids,
        label_ids: labelForm.label_ids,
        reason,
        idempotency_key: `user-label-batch-${crypto.randomUUID()}`,
      });
    }
    ElMessage.success('标签已更新');
    labelDrawerApi.close();
    gridApi.reload();
  } finally {
    labelDrawerApi.unlock();
  }
}

function openReferrerPicker() {
  referrerPickerOpen.value = true;
}

function onReferrerPicked(user: PickedPlatformUser) {
  if (actionUser.value && user.id === actionUser.value.id) {
    ElMessage.warning('上级用户不能是本人');
    return;
  }
  referrerParent.value = user;
  referrerForm.parent_user_id = user.id;
}

function clearReferrerParent() {
  referrerParent.value = null;
  referrerForm.parent_user_id = undefined;
}

async function submitReferrer() {
  if (!actionUser.value) return;
  const reason = referrerForm.reason.trim() || '平台后台修改推荐人';
  if (referrerForm.parent_user_id === actionUser.value.id) {
    ElMessage.warning('上级用户不能是本人');
    return;
  }
  referrerDrawerApi.lock();
  try {
    await changePlatformUserReferrer(actionUser.value.id, {
      parent_user_id: referrerForm.parent_user_id || 0,
      reason,
      idempotency_key: `referrer-${actionUser.value.id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('推荐人已更新');
    referrerDrawerApi.close();
    gridApi.reload();
  } finally {
    referrerDrawerApi.unlock();
  }
}

async function submitPassword() {
  if (!actionUser.value) return;
  const reason = passwordForm.reason.trim() || '平台后台修改密码';
  if (
    passwordForm.password.length < 12 ||
    passwordForm.password.length > 72 ||
    passwordForm.password !== passwordForm.confirm_password
  ) {
    ElMessage.warning('请填写两次一致的 12～72 位新密码');
    return;
  }
  passwordDrawerApi.lock();
  try {
    await resetPlatformUserPassword(actionUser.value.id, {
      password: passwordForm.password,
      reason,
      idempotency_key: `user-password-${actionUser.value.id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('密码已重置');
    passwordDrawerApi.close();
  } finally {
    passwordDrawerApi.unlock();
  }
}

async function submitLevel() {
  if (!actionUser.value) return;
  const reason = levelForm.reason.trim() || '平台后台调整等级';
  levelDrawerApi.lock();
  try {
    await adjustPlatformUserMemberLevel(actionUser.value.id, {
      level_id: levelForm.level_id,
      reason,
      idempotency_key: `member-level-${actionUser.value.id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('等级已调整');
    levelDrawerApi.close();
    gridApi.reload();
  } finally {
    levelDrawerApi.unlock();
  }
}

async function submitCoupon() {
  if (!actionUser.value || !couponForm.coupon_id) {
    ElMessage.warning('请选择优惠券');
    return;
  }
  const reason = couponForm.reason.trim() || '平台后台发送优惠券';
  couponDrawerApi.lock();
  try {
    await issuePlatformUserCoupon(actionUser.value.id, couponForm.coupon_id, {
      reason,
      idempotency_key: `coupon-issue-${actionUser.value.id}-${couponForm.coupon_id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('优惠券已发送');
    couponDrawerApi.close();
  } finally {
    couponDrawerApi.unlock();
  }
}

async function submitSvip() {
  if (!actionUser.value) return;
  if (svipForm.is_svip === 2 && !svipForm.svip_endtime) {
    ElMessage.warning('有效期会员必须填写到期时间');
    return;
  }
  svipDrawerApi.lock();
  try {
    await setPlatformUserListSvip(actionUser.value.id, {
      is_svip: svipForm.is_svip,
      ...(svipForm.is_svip === 2
        ? { svip_endtime: svipForm.svip_endtime }
        : {}),
    });
    ElMessage.success('付费会员已更新');
    svipDrawerApi.close();
    gridApi.reload();
  } finally {
    svipDrawerApi.unlock();
  }
}

onMounted(async () => {
  const [profile, codes, labels] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    fetchPlatformUserLabelOptions().catch(() => ({ list: [] })),
  ]);
  isPlatform.value = profile.roles.includes('platform');
  canRead.value = isPlatform.value && codes.includes('user.list.read');
  canCreate.value = isPlatform.value && codes.includes('user.create.execute');
  canProfile.value = isPlatform.value && codes.includes('user.profile.manage');
  canAsset.value = isPlatform.value && codes.includes('user.asset.adjust');
  canGroup.value = isPlatform.value && codes.includes('user.group.assign');
  canLabel.value = isPlatform.value && codes.includes('user.label.assign');
  canReferrer.value = isPlatform.value && codes.includes('user.referrer.manage');
  canPassword.value = isPlatform.value && codes.includes('user.password.reset');
  canMember.value = isPlatform.value && codes.includes('user.member.adjust');
  canCoupon.value = isPlatform.value && codes.includes('user.coupon.manage');
  canSvip.value = isPlatform.value && codes.includes('user.svip.manage');
  labelOptions.value = labels.list || [];
  gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'label_id',
      componentProps: {
        clearable: true,
        filterable: true,
        placeholder: '请选择',
        options: labelOptions.value.map((item) => ({
          label: item.label_name,
          value: item.label_id,
        })),
      },
    },
  ]);
  if (isPlatform.value) {
    gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="isPlatform && !canRead"
      class="mb-4"
      title="当前账号缺少 user.list.read 权限，列表仍尝试加载；敏感操作已按按钮权限受限。"
      type="warning"
      :closable="false"
    />
    <Grid>
      <template #toolbar-actions>
        <div class="user-list-toolbar">
          <div class="user-list-tabs" role="tablist">
            <button
              v-for="tab in CHANNEL_TABS"
              :key="tab.key || 'all'"
              type="button"
              role="tab"
              class="user-list-tabs__item"
              :aria-selected="channelTab === tab.key"
              :class="{ 'is-active': channelTab === tab.key }"
              @click="setChannelTab(tab.key)"
            >
              {{ tab.label }}
            </button>
          </div>
          <div class="user-list-toolbar__actions">
            <ElButton v-if="canCreate" type="primary" @click="openCreate">
              创建用户
            </ElButton>
            <ElButton
              v-if="canGroup"
              :disabled="!selectedIds.length"
              @click="openBatchGroup"
            >
              批量设置分组
            </ElButton>
            <ElButton
              v-if="canLabel"
              :disabled="!selectedIds.length"
              @click="openBatchLabel"
            >
              批量设置标签
            </ElButton>
          </div>
        </div>
      </template>

      <template #avatar="{ row }">
        <ElImage
          v-if="row.avatar_url"
          :src="resolveCosMediaUrl(row.avatar_url)"
          fit="cover"
          class="user-avatar"
        >
          <template #error>
            <div class="user-avatar user-avatar--empty">无</div>
          </template>
        </ElImage>
        <div v-else class="user-avatar user-avatar--empty">—</div>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canProfile"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElDropdown
          v-if="hasMoreActions"
          trigger="click"
          @command="(cmd: MoreCommand) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem v-if="canAsset" command="balance">
                调整余额
              </ElDropdownItem>
              <ElDropdownItem v-if="canAsset" command="points">
                调整积分
              </ElDropdownItem>
              <ElDropdownItem v-if="canGroup" command="group">
                设置分组
              </ElDropdownItem>
              <ElDropdownItem v-if="canLabel" command="label">
                设置标签
              </ElDropdownItem>
              <ElDropdownItem v-if="canReferrer" command="referrer">
                修改推荐人
              </ElDropdownItem>
              <ElDropdownItem v-if="canPassword" command="password">
                修改密码
              </ElDropdownItem>
              <ElDropdownItem v-if="canMember" command="level">
                调整等级
              </ElDropdownItem>
              <ElDropdownItem v-if="canCoupon" command="coupon">
                发送优惠券
              </ElDropdownItem>
              <ElDropdownItem v-if="canSvip" command="svip">
                付费会员设置
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="detail">
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="用户 ID">
                {{ detail.profile.id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="昵称">
                {{ detail.profile.nickname }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="手机号（脱敏）">
                {{ detail.profile.mobile || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="用户类型">
                {{ channelLabel(detail.profile.source_channel) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="注册时间">
                {{ formatShanghaiDateTime(detail.profile.created_at) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="账户余额">
                ¥{{ Number(detail.profile.balance).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="积分 / 佣金">
                {{ detail.profile.points }} / ¥{{
                  Number(detail.profile.commission).toFixed(2)
                }}
              </ElDescriptionsItem>
            </ElDescriptions>
            <ElTabs class="mt-4">
              <ElTabPane label="最近订单">
                <ElTable :data="detail.orders" max-height="360">
                  <ElTableColumn label="订单号" min-width="180" prop="order_no" />
                  <ElTableColumn label="店铺" min-width="130" prop="store_name" />
                  <ElTableColumn label="实付" width="100">
                    <template #default="{ row }">
                      ¥{{ Number(row.pay_amount).toFixed(2) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="件数" prop="total_quantity" width="70" />
                  <ElTableColumn label="状态" width="100">
                    <template #default="{ row }">
                      {{ orderLabels[row.status] || row.status }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="创建时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="资产流水">
                <ElTable :data="detail.assets" max-height="360">
                  <ElTableColumn prop="id" label="ID" width="80" />
                  <ElTableColumn label="类型" width="90">
                    <template #default="{ row }">
                      {{ assetLabels[row.asset_type] || row.asset_type }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="变动" width="110">
                    <template #default="{ row }">
                      <span
                        :class="
                          row.amount < 0 ? 'text-red-500' : 'text-green-600'
                        "
                      >
                        {{ row.amount > 0 ? '+' : ''
                        }}{{ Number(row.amount).toFixed(2) }}
                      </span>
                    </template>
                  </ElTableColumn>
                  <ElTableColumn
                    prop="reference_type"
                    label="业务来源"
                    min-width="120"
                  />
                  <ElTableColumn
                    prop="reference_id"
                    label="业务引用"
                    min-width="120"
                  />
                  <ElTableColumn label="创建时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="会员变更">
                <ElTable :data="detail.membership_logs" max-height="360">
                  <ElTableColumn
                    prop="previous_level_name"
                    label="原等级"
                    min-width="110"
                  />
                  <ElTableColumn prop="level_name" label="当前等级" min-width="110" />
                  <ElTableColumn label="变更类型" width="100">
                    <template #default="{ row }">
                      {{ membershipLabels[row.change_type] || row.change_type }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn prop="note" label="说明" min-width="220" />
                  <ElTableColumn label="时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="签到记录">
                <ElTable :data="detail.signs" max-height="360">
                  <ElTableColumn prop="sign_date" label="签到日期" min-width="140" />
                  <ElTableColumn prop="points" label="获得积分" width="110" />
                  <ElTableColumn
                    prop="continuous_days"
                    label="连续天数"
                    width="110"
                  />
                  <ElTableColumn label="记录时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="浏览记录">
                <ElTable :data="detail.browse_history" max-height="360">
                  <ElTableColumn prop="product_id" label="商品 ID" width="100" />
                  <ElTableColumn prop="title" label="商品" min-width="180" />
                  <ElTableColumn prop="store_name" label="店铺" min-width="140" />
                  <ElTableColumn label="浏览时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.viewed_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="持有优惠券">
                <ElTable :data="detail.coupons" max-height="360">
                  <ElTableColumn prop="coupon_id" label="券 ID" width="90" />
                  <ElTableColumn prop="name" label="优惠券" min-width="160" />
                  <ElTableColumn label="优惠" width="100">
                    <template #default="{ row }">
                      {{
                        row.discount_type === 'rate'
                          ? `${row.discount_value / 10} 折`
                          : `¥${Number(row.discount_value).toFixed(2)}`
                      }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="门槛" width="100">
                    <template #default="{ row }">
                      ¥{{ Number(row.min_amount).toFixed(2) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="状态" width="100">
                    <template #default="{ row }">
                      {{ couponStatusLabels[row.status] || row.status }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="领取时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.obtained_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="推荐关系">
                <ElDescriptions :column="1" border>
                  <ElDescriptionsItem label="上级用户">
                    {{
                      detail.distribution.parent_user_id
                        ? `${detail.distribution.parent_nickname || '未知用户'}（#${detail.distribution.parent_user_id}）`
                        : '未绑定'
                    }}
                  </ElDescriptionsItem>
                  <ElDescriptionsItem label="直推用户数">
                    {{ detail.distribution.direct_user_count }}
                  </ElDescriptionsItem>
                  <ElDescriptionsItem label="推广员资格">
                    <ElTag
                      :type="
                        detail.distribution.promoter_status === 1
                          ? 'success'
                          : 'info'
                      "
                    >
                      {{
                        detail.distribution.promoter_status === 1
                          ? '启用'
                          : '未开通或已停用'
                      }}
                    </ElTag>
                  </ElDescriptionsItem>
                </ElDescriptions>
              </ElTabPane>
            </ElTabs>
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>

    <CreateDrawer>
      <ElForm label-width="120px" class="user-create-form">
        <ElFormItem label="手机号(账号)" required>
          <ElInput
            v-model="createForm.account"
            maxlength="191"
            placeholder="请输入手机号(账号)"
            autocomplete="off"
          />
        </ElFormItem>
        <ElFormItem label="登录密码" required>
          <ElInput
            v-model="createForm.password"
            type="password"
            show-password
            autocomplete="new-password"
            placeholder="请输入登录密码（12～72 位）"
          />
        </ElFormItem>
        <ElFormItem label="确认密码" required>
          <ElInput
            v-model="createForm.confirm_password"
            type="password"
            show-password
            autocomplete="new-password"
            placeholder="请再次输入密码"
          />
        </ElFormItem>
        <ElFormItem label="用户昵称">
          <ElInput
            v-model="createForm.nickname"
            maxlength="64"
            placeholder="请输入用户昵称"
          />
        </ElFormItem>
        <ElFormItem label="头像">
          <ImageField v-model="createForm.avatar_url" :preview-size="72" />
        </ElFormItem>
        <ElFormItem label="真实姓名">
          <ElInput
            v-model="createForm.real_name"
            maxlength="64"
            placeholder="请输入真实姓名"
          />
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput
            v-model="createForm.phone"
            maxlength="32"
            placeholder="请输入手机号"
          />
        </ElFormItem>
        <ElFormItem label="身份证">
          <ElInput
            v-model="createForm.id_card"
            maxlength="32"
            placeholder="请输入身份证"
          />
        </ElFormItem>
        <ElFormItem label="性别">
          <ElRadioGroup v-model="createForm.gender">
            <ElRadio :value="0">保密</ElRadio>
            <ElRadio :value="1">男</ElRadio>
            <ElRadio :value="2">女</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="状态" required>
          <div>
            <ElRadioGroup v-model="createForm.status">
              <ElRadio :value="0">禁用</ElRadio>
              <ElRadio :value="1">正常</ElRadio>
            </ElRadioGroup>
            <div class="form-tip">禁用之后该用户不可登录商城</div>
          </div>
        </ElFormItem>
        <ElFormItem label="推广员" required>
          <ElRadioGroup v-model="createForm.is_promoter">
            <ElRadio :value="0">关闭</ElRadio>
            <ElRadio :value="1">开启</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </CreateDrawer>

    <EditDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="昵称" required>
          <ElInput v-model="editForm.nickname" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="头像">
          <ImageField v-model="editForm.avatar_url" :preview-size="72" />
        </ElFormItem>
        <ElFormItem label="性别">
          <ElRadioGroup v-model="editForm.gender">
            <ElRadio :value="0">保密</ElRadio>
            <ElRadio :value="1">男</ElRadio>
            <ElRadio :value="2">女</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="个人简介">
          <ElInput
            v-model="editForm.bio"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="调整原因" required>
          <ElInput
            v-model="editForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </EditDrawer>

    <AssetDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="当前值">
          <span v-if="assetKind === 'balance'">
            ¥{{ Number(actionUser?.balance || 0).toFixed(2) }}
          </span>
          <span v-else>{{ actionUser?.points || 0 }}</span>
        </ElFormItem>
        <ElFormItem label="调整方式" required>
          <ElRadioGroup v-model="assetForm.direction">
            <ElRadio value="inc">增加</ElRadio>
            <ElRadio value="dec">减少</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem
          :label="assetKind === 'balance' ? '金额' : '积分'"
          required
        >
          <ElInputNumber
            v-model="assetForm.amount"
            :min="assetKind === 'points' ? 1 : 0.01"
            :precision="assetKind === 'points' ? 0 : 2"
            :step="assetKind === 'points' ? 1 : 0.01"
            class="w-full"
            controls-position="right"
          />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="assetForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </AssetDrawer>

    <GroupDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span v-if="actionUserIds.length === 1">
            {{ actionUser?.nickname || `用户 #${actionUserIds[0]}` }}
          </span>
          <span v-else>已选 {{ actionUserIds.length }} 位用户</span>
        </ElFormItem>
        <ElFormItem label="用户分组" required>
          <ElSelect v-model="groupForm.group_id" class="w-full" filterable>
            <ElOption :value="0" label="未分组" />
            <ElOption
              v-for="group in groupOptions"
              :key="group.group_id"
              :value="group.group_id"
              :label="group.group_name"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="groupForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </GroupDrawer>

    <LabelDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span v-if="actionUserIds.length === 1">
            {{ actionUser?.nickname || `用户 #${actionUserIds[0]}` }}
          </span>
          <span v-else>已选 {{ actionUserIds.length }} 位用户</span>
        </ElFormItem>
        <ElFormItem label="用户标签">
          <ElSelect
            v-model="labelForm.label_ids"
            class="w-full"
            multiple
            clearable
            filterable
            placeholder="留空即清除运营标签"
          >
            <ElOption
              v-for="label in labelOptions"
              :key="label.label_id"
              :value="label.label_id"
              :label="label.label_name"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="labelForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </LabelDrawer>

    <ReferrerDrawer>
      <ElForm label-width="108px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="上级/推荐人">
          <div class="referrer-user-row">
            <div v-if="referrerParent" class="referrer-user-summary">
              <ElAvatar
                v-if="referrerParent.avatar_url"
                :size="32"
                :src="referrerParent.avatar_url"
              />
              <ElAvatar v-else :size="32">
                {{ (referrerParent.nickname || '?').slice(0, 1) }}
              </ElAvatar>
              <div class="referrer-user-meta">
                <div class="referrer-user-name">
                  {{ referrerParent.nickname }}
                </div>
                <div class="referrer-user-id">
                  UID {{ referrerParent.id }}
                </div>
              </div>
              <ElButton link type="danger" @click="clearReferrerParent">
                清除
              </ElButton>
            </div>
            <span v-else class="referrer-user-empty">未选择上级</span>
            <ElButton type="primary" plain @click="openReferrerPicker">
              选择
            </ElButton>
          </div>
          <div class="mt-1 text-xs text-muted-foreground">
            清除后提交即解除上级推荐人
          </div>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="referrerForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </ReferrerDrawer>

    <UserPickerModal
      v-model:open="referrerPickerOpen"
      @select="onReferrerPicked"
    />

    <PasswordDrawer>
      <ElForm label-width="108px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="新密码" required>
          <ElInput
            v-model="passwordForm.password"
            type="password"
            show-password
            autocomplete="new-password"
          />
        </ElFormItem>
        <ElFormItem label="确认密码" required>
          <ElInput
            v-model="passwordForm.confirm_password"
            type="password"
            show-password
            autocomplete="new-password"
          />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="passwordForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </PasswordDrawer>

    <LevelDrawer>
      <ElForm label-width="108px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="目标等级" required>
          <ElSelect v-model="levelForm.level_id" class="w-full">
            <ElOption :value="0" label="普通会员（清除等级）" />
            <ElOption
              v-for="level in levelOptions"
              :key="level.id"
              :value="level.id"
              :label="`${level.name}（等级 ${level.rank}）`"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="levelForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </LevelDrawer>

    <CouponDrawer>
      <ElForm label-width="108px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="优惠券" required>
          <ElSelect
            v-model="couponForm.coupon_id"
            class="w-full"
            filterable
            placeholder="选择有效优惠券"
          >
            <ElOption
              v-for="item in couponOptions"
              :key="item.coupon_id"
              :value="item.coupon_id"
              :label="`${item.name}（#${item.coupon_id}）`"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="couponForm.reason"
            type="textarea"
            :rows="2"
            maxlength="500"
          />
        </ElFormItem>
      </ElForm>
    </CouponDrawer>

    <SvipDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span>{{
            actionUser?.nickname || `用户 #${actionUser?.id || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="会员类型" required>
          <ElSelect v-model="svipForm.is_svip" class="w-full">
            <ElOption label="普通用户" :value="0" />
            <ElOption label="体验会员" :value="1" />
            <ElOption label="有效期会员" :value="2" />
            <ElOption label="永久会员" :value="3" />
            <ElOption label="关闭会员" :value="-1" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="svipForm.is_svip === 2" label="到期时间" required>
          <ElDatePicker
            v-model="svipForm.svip_endtime"
            class="w-full"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            type="datetime"
          />
        </ElFormItem>
      </ElForm>
    </SvipDrawer>
  </Page>
</template>

<style scoped>
.user-list-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.user-list-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 18px;
  border-bottom: 1px solid hsl(var(--border));
  padding-bottom: 8px;
}

.user-list-tabs__item {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 6px 0;
  font-size: 13px;
  color: hsl(var(--muted-foreground));
  cursor: pointer;
}

.user-list-tabs__item.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
  box-shadow: inset 0 -2px 0 hsl(var(--primary));
}

.user-list-toolbar__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
}

.user-avatar--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.4);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.referrer-user-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  width: 100%;
}

.referrer-user-summary {
  display: flex;
  gap: 10px;
  align-items: center;
  min-width: 0;
}

.referrer-user-meta {
  min-width: 0;
}

.referrer-user-name {
  font-size: 14px;
  line-height: 1.3;
  color: hsl(var(--foreground));
}

.referrer-user-id {
  font-size: 12px;
  line-height: 1.3;
  color: hsl(var(--muted-foreground));
}

.referrer-user-empty {
  font-size: 13px;
  color: hsl(var(--muted-foreground));
}

.user-create-form .form-tip {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.4;
  color: hsl(var(--muted-foreground));
}
</style>
