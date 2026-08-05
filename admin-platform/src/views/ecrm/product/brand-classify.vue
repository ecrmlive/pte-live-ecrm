<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  createPlatformBrandCategoryApi,
  deletePlatformBrandCategoryApi,
  listPlatformBrandCategoriesApi,
  updatePlatformBrandCategoryApi,
  type PlatformBrandCategory,
} from '#/api/core/platform-catalog';
import { EcrmFormDialog, EcrmListPage } from '#/components/ecrm';

const rows = ref<PlatformBrandCategory[]>([]);
const loading = ref(false);
const open = ref(false);
const editing = ref<PlatformBrandCategory>();
const form = reactive({ cate_name: '', is_show: 1, pid: 0, sort: 0 });

const flatParents = computed(() => flatten(rows.value));

function flatten(nodes: PlatformBrandCategory[], prefix = ''): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.brand_category_id },
    ...flatten(node.children || [], `${prefix}— `),
  ]);
}

async function load() {
  loading.value = true;
  try {
    rows.value = (await listPlatformBrandCategoriesApi()).list || [];
  } finally {
    loading.value = false;
  }
}

function add(parentID = 0) {
  editing.value = undefined;
  Object.assign(form, { cate_name: '', is_show: 1, pid: parentID, sort: 0 });
  open.value = true;
}

function edit(row: PlatformBrandCategory) {
  editing.value = row;
  Object.assign(form, { cate_name: row.cate_name, is_show: row.is_show, pid: row.pid, sort: row.sort });
  open.value = true;
}

async function save() {
  if (!form.cate_name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const body = { cate_name: form.cate_name.trim(), is_show: form.is_show, pid: form.pid, sort: form.sort };
  if (editing.value) {
    await updatePlatformBrandCategoryApi(editing.value.brand_category_id, body);
  } else {
    await createPlatformBrandCategoryApi(body);
  }
  open.value = false;
  ElMessage.success('品牌分类已保存');
  await load();
}

async function remove(row: PlatformBrandCategory) {
  try {
    await ElMessageBox.confirm(`删除分类“${row.cate_name}”？子分类与品牌须先清空。`, '删除品牌分类', { type: 'warning' });
    await deletePlatformBrandCategoryApi(row.brand_category_id);
    ElMessage.success('品牌分类已删除');
    await load();
  } catch {
    /* 取消 */
  }
}

onMounted(() => void load());
</script>

<template>
  <EcrmListPage title="品牌分类" description="维护品牌分类树；品牌列表可按分类筛选。权限码 product.brand.manage。">
    <template #actions>
      <el-button type="primary" @click="add(0)">新增分类</el-button>
    </template>
    <el-table v-loading="loading" :data="rows" row-key="brand_category_id" default-expand-all>
      <el-table-column label="分类名称" min-width="220" prop="cate_name" />
      <el-table-column label="排序" prop="sort" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template #default="{ row }">
          <el-button link type="primary" @click="add(row.brand_category_id)">加子类</el-button>
          <el-button link type="primary" @click="edit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </EcrmListPage>

  <EcrmFormDialog v-model="open" :title="editing ? '编辑品牌分类' : '新增品牌分类'">
    <el-form label-width="88px">
      <el-form-item label="上级分类">
        <el-select v-model="form.pid" clearable class="w-full" placeholder="顶级分类">
          <el-option label="顶级分类" :value="0" />
          <el-option
            v-for="item in flatParents.filter((x) => x.value !== editing?.brand_category_id)"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="分类名称" required><el-input v-model="form.cate_name" /></el-form-item>
      <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item>
      <el-form-item label="显示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">取消</el-button>
      <el-button type="primary" @click="save">保存</el-button>
    </template>
  </EcrmFormDialog>
</template>
