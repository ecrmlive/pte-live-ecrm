<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  createPlatformBrandApi,
  deletePlatformBrandApi,
  listPlatformBrandCategoriesApi,
  listPlatformBrandsApi,
  updatePlatformBrandApi,
  type PlatformBrand,
  type PlatformBrandCategory,
} from '#/api/core/platform-catalog';
import { EcrmFormDialog, EcrmListPage } from '#/components/ecrm';

const rows = ref<PlatformBrand[]>([]);
const categories = ref<PlatformBrandCategory[]>([]);
const loading = ref(false);
const open = ref(false);
const editing = ref<PlatformBrand>();
const filterCategoryID = ref<number>();
const form = reactive({ brand_name: '', category_id: 0, is_show: 1, sort: 0 });

const categoryOptions = computed(() => flatten(categories.value));
const categoryName = (id: number) => categoryOptions.value.find((x) => x.value === id)?.label || (id ? `#${id}` : '未分类');

function flatten(nodes: PlatformBrandCategory[], prefix = ''): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.brand_category_id },
    ...flatten(node.children || [], `${prefix}— `),
  ]);
}

async function load() {
  loading.value = true;
  try {
    const [brandPage, categoryPage] = await Promise.all([
      listPlatformBrandsApi(filterCategoryID.value ? { category_id: filterCategoryID.value } : undefined),
      listPlatformBrandCategoriesApi(),
    ]);
    rows.value = brandPage.list || [];
    categories.value = categoryPage.list || [];
  } finally {
    loading.value = false;
  }
}

function add() {
  editing.value = undefined;
  Object.assign(form, { brand_name: '', category_id: filterCategoryID.value || 0, is_show: 1, sort: 0 });
  open.value = true;
}

function edit(row: PlatformBrand) {
  editing.value = row;
  Object.assign(form, {
    brand_name: row.brand_name,
    category_id: row.category_id || 0,
    is_show: row.is_show,
    sort: row.sort,
  });
  open.value = true;
}

async function save() {
  if (!form.brand_name.trim()) {
    ElMessage.warning('请填写品牌名称');
    return;
  }
  const body = {
    brand_name: form.brand_name.trim(),
    category_id: form.category_id || 0,
    is_show: form.is_show,
    sort: form.sort,
  };
  if (editing.value) await updatePlatformBrandApi(editing.value.brand_id, body);
  else await createPlatformBrandApi(body);
  open.value = false;
  ElMessage.success('品牌已保存');
  await load();
}

async function remove(row: PlatformBrand) {
  try {
    await ElMessageBox.confirm(`删除品牌“${row.brand_name}”后不可恢复，是否继续？`, '删除品牌', { type: 'warning' });
    await deletePlatformBrandApi(row.brand_id);
    ElMessage.success('品牌已删除');
    await load();
  } catch {
    /* 取消 */
  }
}

onMounted(() => void load());
</script>

<template>
  <EcrmListPage title="品牌列表" description="维护平台品牌；可按品牌分类筛选。权限码 product.brand.manage。">
    <template #filters>
      <el-select v-model="filterCategoryID" clearable class="w-52" placeholder="全部分类" @change="load">
        <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </template>
    <template #actions>
      <el-button type="primary" @click="add">新增品牌</el-button>
    </template>
    <el-table v-loading="loading" :data="rows" row-key="brand_id">
      <el-table-column label="ID" prop="brand_id" width="90" />
      <el-table-column label="品牌名称" min-width="200" prop="brand_name" />
      <el-table-column label="分类" min-width="160">
        <template #default="{ row }">{{ categoryName(row.category_id) }}</template>
      </el-table-column>
      <el-table-column label="排序" prop="sort" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button link type="primary" @click="edit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </EcrmListPage>

  <EcrmFormDialog v-model="open" :title="editing ? '编辑品牌' : '新增品牌'">
    <el-form label-width="88px">
      <el-form-item label="品牌名称" required><el-input v-model="form.brand_name" /></el-form-item>
      <el-form-item label="品牌分类">
        <el-select v-model="form.category_id" clearable class="w-full" placeholder="未分类">
          <el-option label="未分类" :value="0" />
          <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item>
      <el-form-item label="显示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">取消</el-button>
      <el-button type="primary" @click="save">保存</el-button>
    </template>
  </EcrmFormDialog>
</template>
