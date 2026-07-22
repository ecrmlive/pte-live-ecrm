<template>
  <a-card :bordered="false" class="page-card">
    <a-form layout="vertical" :model="form" @finish="submit">
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="商品名称" name="store_name" :rules="[{ required: true, message: '必填' }]">
            <a-input v-model:value="form.store_name" maxlength="64" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="平台分类" name="cate_id" :rules="[{ required: true, message: '请选择分类' }]">
            <a-tree-select
              v-model:value="form.cate_id"
              :tree-data="cateTree"
              :field-names="{ label: 'cate_name', value: 'store_category_id', children: 'children' }"
              placeholder="选择分类"
              tree-default-expand-all
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="售价" name="price" :rules="[{ required: true, message: '必填' }]">
            <a-input-number v-model:value="form.price" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="原价">
            <a-input-number v-model:value="form.ot_price" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="会员价类型">
            <a-select v-model:value="form.svip_price_type">
              <a-select-option :value="0">不参加</a-select-option>
              <a-select-option :value="1">默认比例(9折)</a-select-option>
              <a-select-option :value="2">自定义价</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col v-if="form.svip_price_type === 2" :span="12">
          <a-form-item label="会员价">
            <a-input-number v-model:value="form.svip_price" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="成本价">
            <a-input-number v-model:value="form.cost" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="库存" name="stock" :rules="[{ required: true, message: '必填' }]">
            <a-input-number v-model:value="form.stock" :min="0" :precision="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="单位">
            <a-input v-model:value="form.unit_name" placeholder="件" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="品牌">
            <a-select v-model:value="form.brand_id" allow-clear placeholder="可选">
              <a-select-option v-for="b in brands" :key="b.brand_id" :value="b.brand_id">
                {{ b.brand_name }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="24">
          <a-form-item label="简介">
            <a-textarea v-model:value="form.store_info" :rows="3" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="关键字">
            <a-input v-model:value="form.keyword" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="配送方式">
            <a-select v-model:value="form.delivery_way">
              <a-select-option value="1">快递</a-select-option>
              <a-select-option value="2">到店</a-select-option>
              <a-select-option value="1,2">快递+到店</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="上架（需审核通过后才对 C 端可见）">
            <a-switch v-model:checked="showOn" checked-children="上架" un-checked-children="下架" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="SVIP 价类型">
            <a-select v-model:value="form.svip_price_type">
              <a-select-option :value="0">不参加</a-select-option>
              <a-select-option :value="1">默认折扣(9折)</a-select-option>
              <a-select-option :value="2">自定义会员价</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="自定义会员价">
            <a-input-number v-model:value="form.svip_price" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>

      <a-space>
        <a-button type="primary" html-type="submit" :loading="saving">保存</a-button>
        <a-button @click="router.push('/product/list')">返回列表</a-button>
      </a-space>
    </a-form>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import {
  createProduct,
  fetchBrands,
  fetchCategories,
  fetchProduct,
  updateProduct,
  type Brand,
  type CategoryNode,
  type ProductSave,
} from '@/api/catalog';
import { useAuthStore } from '@/stores/auth';

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const saving = ref(false);
const cateTree = ref<CategoryNode[]>([]);
const brands = ref<Brand[]>([]);
const productId = computed(() => Number(route.query.id || 0));
const showOn = ref(false);

const form = reactive<ProductSave>({
  store_name: '',
  store_info: '',
  keyword: '',
  cate_id: undefined as unknown as number,
  brand_id: undefined,
  unit_name: '件',
  price: 0,
  ot_price: 0,
  cost: 0,
  stock: 0,
  image: '',
  slider_image: '',
  delivery_way: '2',
  type: 0,
  spec_type: 0,
  svip_price_type: 0,
  svip_price: 0,
  mer_svip_status: 1,
});

async function loadMeta() {
  const [c, b] = await Promise.all([fetchCategories(), fetchBrands()]);
  cateTree.value = c.data.list || [];
  brands.value = b.data.list || [];
}

async function loadDetail() {
  if (!productId.value) return;
  const { data } = await fetchProduct(productId.value);
  form.store_name = data.store_name;
  form.store_info = data.store_info;
  form.keyword = data.keyword;
  form.cate_id = data.cate_id;
  form.brand_id = data.brand_id || undefined;
  form.unit_name = data.unit_name || '件';
  form.price = Number(data.price);
  form.ot_price = Number(data.ot_price);
  form.stock = data.stock;
  form.image = data.image;
  form.slider_image = data.slider_image;
  form.delivery_way = data.delivery_way || '2';
  form.type = data.type;
  form.spec_type = data.spec_type;
  form.svip_price_type = data.svip_price_type ?? 0;
  form.svip_price = Number(data.svip_price || 0);
  form.mer_svip_status = data.mer_svip_status ?? 1;
  showOn.value = data.is_show === 1;
}

async function submit() {
  if (!productId.value && !auth.hasPerm('product/create')) {
    message.warning('无发布商品权限');
    return;
  }
  saving.value = true;
  try {
    const body: ProductSave = {
      ...form,
      cate_id: Number(form.cate_id),
      brand_id: form.brand_id ? Number(form.brand_id) : 0,
      is_show: showOn.value ? 1 : 0,
    };
    if (productId.value) {
      await updateProduct(productId.value, body);
      message.success('已保存（若商户需审核，将回到待审）');
    } else {
      await createProduct(body);
      message.success('已创建');
    }
    router.push('/product/list');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  await loadMeta();
  await loadDetail();
});
</script>

<style scoped>
.page-card {
  border-radius: 14px;
  max-width: 960px;
}
</style>
