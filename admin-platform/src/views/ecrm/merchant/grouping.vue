<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch,
} from 'vue';

import { Page, confirm, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElPagination,
  ElRadio,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
  ElTreeSelect,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteStoreGroup,
  fetchPlatformDiyPages,
  fetchStoreGroup,
  fetchStoreGroupMerchants,
  fetchStoreGroups,
  saveStoreGroup,
  setStoreGroupStatus,
  setStoreGroupTemplate,
  type DiyPageOption,
  type StoreGroupRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { fetchMapClientConfig } from '#/api/core/cloud-config';
import StorePickerModal, {
  type PickedStore,
} from '#/components/ecrm/store-picker-modal.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'view';

interface LinkedMerchant {
  mer_id: number;
  mer_name: string;
  real_name: string;
  mer_phone: string;
}

interface ParentOption {
  id: number;
  name: string;
  level: number;
  children?: ParentOption[];
}

declare global {
  interface Window {
    AMap?: any;
    _AMapSecurityConfig?: { securityJsCode: string };
  }
}

const MAX_GROUP_LEVEL = 2; // 0/1/2 → 最多三级
const NAME_MAX = 20;
const DEFAULT_LAT = 23.1291;
const DEFAULT_LNG = 113.2644;

const canManage = ref(false);
const treeRows = ref<StoreGroupRow[]>([]);
const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const activeTab = ref('basic');
const linkedMerchants = ref<LinkedMerchant[]>([]);
const storePickerOpen = ref(false);
const amapKey = ref('');
const amapSecurityCode = ref('');
const amapConfigured = ref(false);
const amapConfigLoaded = ref(false);
const amapConfigError = ref('');

const form = reactive({
  parent_id: 0,
  name: '',
  sort: 0,
  status: true,
  diy_page_id: 0,
  diy_page_name: '',
  positioning_status: false,
  longitude: undefined as number | undefined,
  latitude: undefined as number | undefined,
  address: '',
});

const isReadonly = computed(() => drawerMode.value === 'view');

const parentTreeOptions = computed<ParentOption[]>(() => {
  const exclude = new Set<number>();
  if (editingId.value) {
    const collect = (nodes: StoreGroupRow[], underSelf = false) => {
      for (const n of nodes) {
        const hit = underSelf || n.id === editingId.value;
        if (hit) exclude.add(n.id);
        if (n.children?.length) collect(n.children, hit);
      }
    };
    collect(treeRows.value);
  }
  const mapNode = (n: StoreGroupRow): ParentOption | null => {
    if (exclude.has(n.id) || n.level >= MAX_GROUP_LEVEL) return null;
    const children = (n.children || [])
      .map(mapNode)
      .filter((x): x is ParentOption => !!x);
    return { id: n.id, name: n.name, level: n.level, children };
  };
  return treeRows.value
    .map(mapNode)
    .filter((x): x is ParentOption => !!x);
});

function filterTree(
  nodes: StoreGroupRow[],
  keyword: string,
  status?: number,
): StoreGroupRow[] {
  return nodes
    .map((node) => {
      const children = node.children
        ? filterTree(node.children, keyword, status)
        : undefined;
      const nameMatch =
        !keyword || node.name.toLowerCase().includes(keyword);
      const statusMatch =
        status !== 0 && status !== 1 ? true : node.status === status;
      if ((nameMatch && statusMatch) || (children && children.length)) {
        return { ...node, children };
      }
      return null;
    })
    .filter((node): node is StoreGroupRow => node !== null);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入分组名称' },
    fieldName: 'keyword',
    label: '分组名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '关闭', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '开启状态',
  },
]);

const gridOptions: VxeGridProps<StoreGroupRow> = {
  columns: [
    {
      field: 'name',
      minWidth: 220,
      showOverflow: false,
      title: '分组名称',
      treeNode: true,
    },
    {
      field: 'merchant_count',
      title: '关联店铺',
      width: 100,
    },
    {
      field: 'diy_page_id',
      formatter: ({ cellValue }) =>
        cellValue === 0 || cellValue == null || cellValue === ''
          ? '—'
          : String(cellValue),
      title: '模版ID',
      width: 100,
    },
    {
      field: 'diy_page_name',
      formatter: ({ cellValue }) => cellValue || '—',
      title: '模版名称',
      minWidth: 140,
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '开启状态',
      width: 120,
    },
    platformListActionColumn({ minWidth: 360, width: 380 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        const status =
          statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined;
        // 全量拉树再前端过滤，保留「子命中、父仍展示」的树形搜索体验
        let list = (await fetchStoreGroups()).list || [];
        treeRows.value = list;
        if (keyword || status !== undefined) {
          list = filterTree(list, keyword, status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  treeConfig: {
    childrenField: 'children',
    expandAll: false,
  },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [GroupDrawer, groupDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (drawerMode.value === 'view') {
      if (canManage.value) await openEditById(editingId.value);
      else groupDrawerApi.close();
      return;
    }
    await saveGroup();
  },
  onOpenChange(isOpen) {
    if (!isOpen) {
      destroyMap();
    } else if (form.positioning_status) {
      void nextTick(() => initMap());
    }
  },
});

/* ---------- 选择模版（装修） ---------- */
/** list：列表操作立即绑定；form：写入抽屉表单，随 create/update 提交 */
const templatePickerMode = ref<'list' | 'form'>('list');
const templateTargetId = ref(0);
const templateKeyword = ref('');
const templateLoading = ref(false);
const templateRows = ref<DiyPageOption[]>([]);
const templateTotal = ref(0);
const templatePage = ref(1);
const templateLimit = ref(10);
const templateSelectedId = ref<number | undefined>(undefined);

const [TemplateModal, templateModalApi] = useVbenModal({
  title: '选择模版',
  class: 'w-[760px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  onConfirm: async () => {
    if (!templateSelectedId.value) {
      ElMessage.warning('请选择模版');
      return;
    }
    if (templatePickerMode.value === 'form') {
      const picked = templateRows.value.find(
        (row) => row.id === templateSelectedId.value,
      );
      form.diy_page_id = templateSelectedId.value;
      form.diy_page_name = picked?.name || '';
      templateModalApi.close();
      return;
    }
    templateModalApi.lock();
    try {
      await setStoreGroupTemplate(templateTargetId.value, templateSelectedId.value);
      ElMessage.success('装修模版已绑定');
      templateModalApi.close();
      gridApi.reload();
    } finally {
      templateModalApi.unlock();
    }
  },
});

/* ---------- 地图（高德 Web JS API 2.0） ---------- */
const mapEl = ref<HTMLDivElement>();
let mapInstance: any;
let mapMarker: any;
let mapScriptPromise: Promise<void> | null = null;
let syncingFromMap = false;

function resetForm(partial?: Partial<typeof form>) {
  Object.assign(form, {
    parent_id: 0,
    name: '',
    sort: 0,
    status: true,
    diy_page_id: 0,
    diy_page_name: '',
    positioning_status: false,
    longitude: undefined,
    latitude: undefined,
    address: '',
    ...partial,
  });
  linkedMerchants.value = [];
  activeTab.value = 'basic';
}

function openCreate() {
  drawerMode.value = 'create';
  editingId.value = 0;
  resetForm();
  groupDrawerApi
    .setState({ title: '新增店铺分组', confirmText: '保存', showConfirmButton: true })
    .open();
}

function openCreateChild(row: StoreGroupRow) {
  if (row.level >= MAX_GROUP_LEVEL) {
    ElMessage.warning('店铺分组最多支持三级');
    return;
  }
  drawerMode.value = 'create';
  editingId.value = 0;
  resetForm({ parent_id: row.id });
  groupDrawerApi
    .setState({ title: '新增下级分组', confirmText: '保存', showConfirmButton: true })
    .open();
}

async function fillFromGroup(row: StoreGroupRow) {
  resetForm({
    parent_id: row.parent_id || 0,
    name: row.name || '',
    sort: row.sort || 0,
    status: row.status === 1,
    diy_page_id: row.diy_page_id || 0,
    diy_page_name: row.diy_page_name || '',
    positioning_status: row.positioning_status === 1,
    longitude: row.longitude ?? undefined,
    latitude: row.latitude ?? undefined,
    address: row.address || '',
  });
  const merchants = (await fetchStoreGroupMerchants(row.id)).list || [];
  linkedMerchants.value = merchants.map((m) => ({
    mer_id: m.merchant_id,
    mer_name: m.merchant_name,
    real_name: m.contact_name || '',
    mer_phone: m.contact_mobile || '',
  }));
}

async function openEditById(id: number) {
  const detail = await fetchStoreGroup(id);
  drawerMode.value = 'edit';
  editingId.value = id;
  await fillFromGroup(detail);
  groupDrawerApi
    .setState({ title: '编辑店铺分组', confirmText: '保存', showConfirmButton: true })
    .open();
  if (form.positioning_status) {
    await nextTick();
    await initMap();
  }
}

async function openEdit(row: StoreGroupRow) {
  await openEditById(row.id);
}

async function openDetail(row: StoreGroupRow) {
  const detail = await fetchStoreGroup(row.id);
  drawerMode.value = 'view';
  editingId.value = row.id;
  await fillFromGroup(detail);
  groupDrawerApi
    .setState({
      title: '店铺分组详情',
      confirmText: canManage.value ? '编辑' : '关闭',
      showConfirmButton: true,
    })
    .open();
}

async function saveGroup() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分组名称');
    activeTab.value = 'basic';
    return;
  }
  if (form.name.trim().length > NAME_MAX) {
    ElMessage.warning(`分组名称最多 ${NAME_MAX} 个字符`);
    activeTab.value = 'basic';
    return;
  }
  if (form.positioning_status && !form.address.trim()) {
    ElMessage.warning('请填写区域中心');
    activeTab.value = 'basic';
    return;
  }
  if (
    (form.longitude === undefined) !== (form.latitude === undefined)
  ) {
    ElMessage.warning('经度和纬度需同时填写，或同时留空');
    activeTab.value = 'basic';
    return;
  }
  groupDrawerApi.lock();
  try {
    await saveStoreGroup(
      drawerMode.value === 'edit' ? editingId.value : undefined,
      {
        parent_id: form.parent_id || 0,
        name: form.name.trim(),
        sort: form.sort || 0,
        status: form.status,
        diy_page_id: form.diy_page_id || 0,
        positioning_status: form.positioning_status,
        longitude: form.longitude,
        latitude: form.latitude,
        address: form.address.trim(),
        merchant_ids: linkedMerchants.value.map((m) => m.mer_id),
      },
    );
    groupDrawerApi.close();
    ElMessage.success(drawerMode.value === 'edit' ? '店铺分组已更新' : '店铺分组已创建');
    gridApi.reload();
  } finally {
    groupDrawerApi.unlock();
  }
}

async function toggleStatus(row: StoreGroupRow, enabled: boolean) {
  if (!canManage.value) return;
  await setStoreGroupStatus(row.id, enabled);
  ElMessage.success('分组状态已更新');
  gridApi.reload();
}

async function remove(row: StoreGroupRow) {
  try {
    await confirm({
      title: '提示',
      content: `删除“${row.name}”后不可恢复；含子分组时系统会拒绝删除。是否继续？`,
      icon: 'warning',
    });
    await deleteStoreGroup(row.id);
    ElMessage.success('店铺分组已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

async function openTemplatePicker(row: StoreGroupRow) {
  templatePickerMode.value = 'list';
  templateTargetId.value = row.id;
  templateSelectedId.value = row.diy_page_id || undefined;
  templateKeyword.value = '';
  templatePage.value = 1;
  await loadTemplates();
  templateModalApi.open();
}

async function openFormTemplatePicker() {
  templatePickerMode.value = 'form';
  templateTargetId.value = 0;
  templateSelectedId.value = form.diy_page_id || undefined;
  templateKeyword.value = '';
  templatePage.value = 1;
  await loadTemplates();
  templateModalApi.open();
}

function clearFormTemplate() {
  form.diy_page_id = 0;
  form.diy_page_name = '';
}

async function loadTemplates() {
  templateLoading.value = true;
  try {
    const result = await fetchPlatformDiyPages({
      page: templatePage.value,
      limit: templateLimit.value,
      name: templateKeyword.value.trim() || undefined,
    });
    templateRows.value = result.list || [];
    templateTotal.value = result.total || 0;
  } finally {
    templateLoading.value = false;
  }
}

function searchTemplates() {
  templatePage.value = 1;
  void loadTemplates();
}

function onTemplatePageChange(page: number) {
  templatePage.value = page;
  void loadTemplates();
}

function removeLinked(merId: number) {
  linkedMerchants.value = linkedMerchants.value.filter((m) => m.mer_id !== merId);
}

function openStorePicker() {
  storePickerOpen.value = true;
}

function onStoresPicked(stores: PickedStore[]) {
  linkedMerchants.value = stores.map((item) => ({
    mer_id: item.mer_id,
    mer_name: item.mer_name,
    real_name: item.real_name,
    mer_phone: item.mer_phone,
  }));
}

function canAddChild(row: StoreGroupRow) {
  return row.level < MAX_GROUP_LEVEL;
}

async function ensureAmapConfig() {
  if (amapConfigLoaded.value) return amapConfigured.value;
  try {
    const cfg = await fetchMapClientConfig();
    amapKey.value = String(cfg.amap_web_js_key || '').trim();
    amapSecurityCode.value = String(cfg.amap_web_js_security_code || '').trim();
    amapConfigured.value = Boolean(
      cfg.configured && amapKey.value && amapSecurityCode.value,
    );
    amapConfigError.value = amapConfigured.value
      ? ''
      : '未配置高德地图 Key，请在「云服务配置 → 高德地图」或本地 init_key.sql 中填写';
  } catch {
    amapConfigured.value = false;
    amapConfigError.value =
      '读取高德地图配置失败，请确认已登录且 api-platform 可用';
  } finally {
    amapConfigLoaded.value = true;
  }
  return amapConfigured.value;
}

async function loadAmapScript() {
  const ok = await ensureAmapConfig();
  if (!ok) {
    throw new Error(amapConfigError.value || '未配置高德地图');
  }
  if (window.AMap) return;
  window._AMapSecurityConfig = {
    securityJsCode: amapSecurityCode.value,
  };
  if (!mapScriptPromise) {
    mapScriptPromise = new Promise<void>((resolve, reject) => {
      const script = document.createElement('script');
      script.src = `https://webapi.amap.com/maps?v=2.0&key=${encodeURIComponent(amapKey.value)}`;
      script.async = true;
      script.onload = () => resolve();
      script.onerror = () => {
        mapScriptPromise = null;
        reject(new Error('高德地图脚本加载失败'));
      };
      document.head.appendChild(script);
    });
  }
  await mapScriptPromise;
}

function destroyMap() {
  if (mapInstance) {
    mapInstance.destroy?.();
  }
  mapMarker = undefined;
  mapInstance = undefined;
  if (mapEl.value) {
    mapEl.value.innerHTML = '';
  }
}

function applyLatLngToMarker(lat: number, lng: number, moveCenter = true) {
  if (!window.AMap || !mapInstance || !mapMarker) return;
  const pos = [lng, lat];
  mapMarker.setPosition(pos);
  if (moveCenter) mapInstance.setCenter(pos);
}

async function initMap() {
  if (!form.positioning_status || !mapEl.value) return;
  try {
    await loadAmapScript();
  } catch (error) {
    ElMessage.warning(
      error instanceof Error ? error.message : '地图加载失败',
    );
    return;
  }
  const AMap = window.AMap;
  if (!AMap) return;
  const lat = form.latitude ?? DEFAULT_LAT;
  const lng = form.longitude ?? DEFAULT_LNG;
  destroyMap();
  mapInstance = new AMap.Map(mapEl.value, {
    center: [lng, lat],
    zoom: 15,
    viewMode: '2D',
  });
  mapMarker = new AMap.Marker({
    position: [lng, lat],
    map: mapInstance,
  });
  mapInstance.on('click', (event: any) => {
    if (isReadonly.value) return;
    const lnglat = event.lnglat;
    if (!lnglat) return;
    syncingFromMap = true;
    form.latitude = Number(Number(lnglat.getLat()).toFixed(7));
    form.longitude = Number(Number(lnglat.getLng()).toFixed(7));
    mapMarker?.setPosition([form.longitude, form.latitude]);
    reverseGeocode(form.longitude, form.latitude);
    void nextTick(() => {
      syncingFromMap = false;
    });
  });
}

function reverseGeocode(lng?: number, lat?: number) {
  if (lng === undefined || lat === undefined || !window.AMap) return;
  window.AMap.plugin('AMap.Geocoder', () => {
    const geocoder = new window.AMap.Geocoder({ radius: 1000 });
    geocoder.getAddress([lng, lat], (status: string, result: any) => {
      if (status !== 'complete' || !result?.regeocode) return;
      const addr = result.regeocode.formattedAddress;
      if (addr) form.address = addr;
    });
  });
}

async function searchLocation() {
  if (!form.address.trim()) {
    ElMessage.warning('请输入区域中心地址');
    return;
  }
  try {
    await loadAmapScript();
    await initMap();
  } catch (error) {
    ElMessage.warning(
      error instanceof Error ? error.message : '地图加载失败',
    );
    return;
  }
  window.AMap.plugin('AMap.Geocoder', () => {
    const geocoder = new window.AMap.Geocoder({});
    geocoder.getLocation(form.address.trim(), (status: string, result: any) => {
      const loc = result?.geocodes?.[0]?.location;
      if (status !== 'complete' || !loc) {
        ElMessage.warning('未找到该地址对应位置');
        return;
      }
      syncingFromMap = true;
      form.latitude = Number(Number(loc.lat).toFixed(7));
      form.longitude = Number(Number(loc.lng).toFixed(7));
      applyLatLngToMarker(form.latitude, form.longitude);
      ElMessage.success('位置已更新');
      void nextTick(() => {
        syncingFromMap = false;
      });
    });
  });
}

watch(
  () => form.positioning_status,
  async (enabled) => {
    if (!enabled) {
      destroyMap();
      return;
    }
    await ensureAmapConfig();
    await nextTick();
    await initMap();
  },
);

watch(
  () => [form.longitude, form.latitude] as const,
  ([lng, lat]) => {
    if (syncingFromMap || lng === undefined || lat === undefined) return;
    applyLatLngToMarker(lat, lng);
  },
);

onMounted(async () => {
  const [permissions] = await Promise.all([
    getAccessCodesApi(),
    ensureAmapConfig(),
  ]);
  canManage.value = permissions.includes('merchant.group.manage');
});

onBeforeUnmount(() => {
  destroyMap();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="grouping-toolbar">
          <ElAlert
            class="grouping-tip"
            type="warning"
            show-icon
            :closable="false"
            :title="
              '用于标识店铺所属商圈，主要用于移动端的商圈切换，最多支持三级结构，如“陕西/西安/雁塔区”'
            "
          />
          <div class="grouping-toolbar__actions">
            <ElButton
              v-if="canManage"
              :icon="Plus"
              type="primary"
              @click="openCreate"
            >
              新增分组
            </ElButton>
          </div>
        </div>
      </template>

      <template #status="{ row }">
        <ElSwitch
          v-if="canManage"
          :model-value="row.status === 1"
          inline-prompt
          active-text="开启"
          inactive-text="关闭"
          @change="(enabled: string | number | boolean) => toggleStatus(row, Boolean(enabled))"
        />
        <ElTag v-else :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '开启' : '关闭' }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <div class="grouping-ops">
          <ElButton
            v-if="canManage && canAddChild(row)"
            link
            type="primary"
            @click="openCreateChild(row)"
          >
            新增下级
          </ElButton>
          <ElButton
            v-if="canManage"
            link
            type="primary"
            @click="openEdit(row)"
          >
            编辑
          </ElButton>
          <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
          <ElButton
            v-if="canManage"
            link
            type="primary"
            @click="openTemplatePicker(row)"
          >
            模版
          </ElButton>
          <ElButton
            v-if="canManage"
            link
            type="danger"
            @click="remove(row)"
          >
            删除
          </ElButton>
        </div>
      </template>
    </Grid>

    <GroupDrawer>
      <div class="group-drawer">
        <ElTabs v-model="activeTab" class="group-drawer__tabs">
          <ElTabPane label="基础信息" name="basic">
            <ElForm label-width="112px" class="group-form">
              <ElFormItem label="上级分组" required>
                <ElTreeSelect
                  v-model="form.parent_id"
                  :data="[
                    { id: 0, name: '顶级分组', children: parentTreeOptions },
                  ]"
                  node-key="id"
                  :props="{ label: 'name', value: 'id', children: 'children' }"
                  check-strictly
                  :disabled="isReadonly"
                  class="w-full max-w-md"
                  placeholder="顶级分组"
                />
              </ElFormItem>
              <ElFormItem label="分组名称" required>
                <ElInput
                  v-model="form.name"
                  :maxlength="NAME_MAX"
                  show-word-limit
                  :disabled="isReadonly"
                  class="max-w-md"
                  placeholder="请输入分组名称"
                />
              </ElFormItem>
              <ElFormItem label="排序">
                <div>
                  <ElInputNumber
                    v-model="form.sort"
                    :min="0"
                    :disabled="isReadonly"
                  />
                  <div class="field-help">数字越大越靠前</div>
                </div>
              </ElFormItem>
              <ElFormItem label="开启状态">
                <div>
                  <ElSwitch
                    v-model="form.status"
                    inline-prompt
                    active-text="开启"
                    inactive-text="关闭"
                    :disabled="isReadonly"
                  />
                  <div class="field-help">
                    关闭后，移动端商城不展示该分组
                  </div>
                </div>
              </ElFormItem>
              <ElFormItem label="勾选模版">
                <div class="template-field">
                  <span class="template-field__name">
                    {{
                      form.diy_page_id
                        ? form.diy_page_name || `模版ID ${form.diy_page_id}`
                        : '未选择'
                    }}
                  </span>
                  <ElButton
                    v-if="!isReadonly"
                    type="primary"
                    @click="openFormTemplatePicker"
                  >
                    选择模版
                  </ElButton>
                  <ElButton
                    v-if="!isReadonly && form.diy_page_id"
                    @click="clearFormTemplate"
                  >
                    清空
                  </ElButton>
                </div>
              </ElFormItem>
              <ElFormItem label="区域定位">
                <div>
                  <ElSwitch
                    v-model="form.positioning_status"
                    inline-prompt
                    active-text="开启"
                    inactive-text="关闭"
                    :disabled="isReadonly"
                  />
                  <div class="field-help">
                    开启定位后，用户进入商城选择区域，会自动根据用户定位距离区域中心进行推荐
                  </div>
                </div>
              </ElFormItem>
              <template v-if="form.positioning_status">
                <ElFormItem label="区域中心" required>
                  <div class="address-row">
                    <ElInput
                      v-model="form.address"
                      :disabled="isReadonly"
                      class="address-input"
                      placeholder="请输入区域中心地址"
                    />
                    <ElButton
                      v-if="!isReadonly"
                      type="primary"
                      @click="searchLocation"
                    >
                      查找位置
                    </ElButton>
                  </div>
                </ElFormItem>
                <ElFormItem label=" ">
                  <div class="map-panel">
                    <ElAlert
                      v-if="amapConfigLoaded && !amapConfigured"
                      type="warning"
                      :closable="false"
                      show-icon
                      :title="
                        amapConfigError ||
                        '未配置高德地图，仍可手工填写区域中心地址与经纬度。请到「系统设置 → 云服务配置 → 高德地图」或本地 sql/admin/init_key.sql 配置。'
                      "
                    />
                    <div
                      v-show="amapConfigured"
                      ref="mapEl"
                      class="map-canvas"
                    />
                    <div v-if="!isReadonly" class="coord-row">
                      <span>经度</span>
                      <ElInputNumber
                        v-model="form.longitude"
                        :precision="7"
                        :min="-180"
                        :max="180"
                        controls-position="right"
                      />
                      <span>纬度</span>
                      <ElInputNumber
                        v-model="form.latitude"
                        :precision="7"
                        :min="-90"
                        :max="90"
                        controls-position="right"
                      />
                    </div>
                    <div v-else class="field-help">
                      经度 {{ form.longitude ?? '—' }} / 纬度
                      {{ form.latitude ?? '—' }}
                    </div>
                  </div>
                </ElFormItem>
              </template>
            </ElForm>
          </ElTabPane>

          <ElTabPane label="关联店铺" name="stores">
            <div class="stores-pane">
              <ElButton
                v-if="!isReadonly"
                type="primary"
                @click="openStorePicker"
              >
                选择店铺
              </ElButton>
              <ElTable :data="linkedMerchants" class="mt-3" border>
                <ElTableColumn label="店铺ID" prop="mer_id" width="90" />
                <ElTableColumn label="店铺名称" prop="mer_name" min-width="160" />
                <ElTableColumn label="联系人" prop="real_name" width="120" />
                <ElTableColumn label="联系电话" prop="mer_phone" width="140" />
                <ElTableColumn v-if="!isReadonly" label="操作" width="90" fixed="right">
                  <template #default="{ row }">
                    <ElButton
                      link
                      type="primary"
                      @click="removeLinked(row.mer_id)"
                    >
                      删除
                    </ElButton>
                  </template>
                </ElTableColumn>
              </ElTable>
            </div>
          </ElTabPane>
        </ElTabs>
      </div>
    </GroupDrawer>

    <StorePickerModal
      v-model:open="storePickerOpen"
      :selected="linkedMerchants"
      @confirm="onStoresPicked"
    />

    <TemplateModal>
      <div class="picker-filter">
        <ElForm inline @submit.prevent>
          <ElFormItem label="模版名称">
            <ElInput
              v-model="templateKeyword"
              clearable
              placeholder="请输入模版名称"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="searchTemplates">搜索</ElButton>
          </ElFormItem>
        </ElForm>
      </div>
      <ElTable
        v-loading="templateLoading"
        :data="templateRows"
        border
        highlight-current-row
      >
        <ElTableColumn label="" width="56" align="center">
          <template #default="{ row }">
            <ElRadio
              v-model="templateSelectedId"
              :value="row.id"
            />
          </template>
        </ElTableColumn>
        <ElTableColumn label="模版ID" prop="id" width="100" />
        <ElTableColumn label="模版名称" prop="name" min-width="200" />
        <ElTableColumn label="创建时间" min-width="170">
          <template #default="{ row }">
            {{ formatShanghaiDateTime(row.add_time) || '—' }}
          </template>
        </ElTableColumn>
      </ElTable>
      <div class="picker-pager">
        <ElPagination
          background
          layout="prev, pager, next, jumper"
          :current-page="templatePage"
          :page-size="templateLimit"
          :total="templateTotal"
          @current-change="onTemplatePageChange"
        />
      </div>
    </TemplateModal>
  </Page>
</template>

<style scoped>
.grouping-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.grouping-tip {
  width: 100%;
}

.grouping-toolbar__actions {
  display: flex;
  justify-content: flex-start;
}

.grouping-ops {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
}

.group-drawer__tabs :deep(.el-tabs__content) {
  padding-top: 8px;
}

.group-form {
  max-width: 860px;
}

.field-help {
  margin-top: 6px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.template-field {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.template-field__name {
  min-width: 120px;
  color: hsl(var(--foreground));
}

.address-row {
  display: flex;
  gap: 8px;
  width: 100%;
  max-width: 640px;
}

.address-input {
  flex: 1;
}

.map-panel {
  width: 100%;
  max-width: 800px;
}

.map-canvas {
  width: 100%;
  height: 360px;
  margin-top: 8px;
  border: 1px solid hsl(var(--border));
  border-radius: 4px;
  overflow: hidden;
}

.coord-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
}

.stores-pane {
  min-height: 240px;
}

.picker-filter {
  margin-bottom: 12px;
}

.picker-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
