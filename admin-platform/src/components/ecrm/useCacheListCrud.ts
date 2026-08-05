import { onMounted, reactive, ref, type Ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

export interface CacheListItem {
  id: string;
  name: string;
  enabled: boolean;
  remark: string;
}

export function useCacheListCrud(
  fetchList: () => Promise<{ list: CacheListItem[] }>,
  saveList: (list: CacheListItem[]) => Promise<{ list: CacheListItem[] }>,
) {
  const loading = ref(false);
  const saving = ref(false);
  const dialog = ref(false);
  const rows: Ref<CacheListItem[]> = ref([]);
  const editingIndex = ref<number>();
  const form = reactive<CacheListItem>({ id: '', name: '', enabled: true, remark: '' });

  function resetForm() {
    editingIndex.value = undefined;
    Object.assign(form, { id: '', name: '', enabled: true, remark: '' });
  }

  function open(row?: CacheListItem, index?: number) {
    resetForm();
    if (row !== undefined && index !== undefined) {
      editingIndex.value = index;
      Object.assign(form, row);
    }
    dialog.value = true;
  }

  async function load() {
    loading.value = true;
    try {
      rows.value = (await fetchList()).list || [];
    } finally {
      loading.value = false;
    }
  }

  async function save() {
    const name = form.name.trim();
    if (!name) {
      ElMessage.warning('请填写名称');
      return;
    }
    const next = rows.value.map((item) => ({ ...item }));
    const payload: CacheListItem = {
      id: form.id.trim() || name,
      name,
      enabled: form.enabled,
      remark: form.remark.trim(),
    };
    if (editingIndex.value === undefined) next.push(payload);
    else next[editingIndex.value] = payload;

    saving.value = true;
    try {
      rows.value = (await saveList(next)).list || [];
      dialog.value = false;
      ElMessage.success('已保存');
    } finally {
      saving.value = false;
    }
  }

  async function remove(index: number, name: string) {
    try {
      await ElMessageBox.confirm(`确认删除“${name}”？`, '删除确认', { type: 'warning' });
      rows.value = (await saveList(rows.value.filter((_, i) => i !== index))).list || [];
      ElMessage.success('已删除');
    } catch {
      // 用户取消
    }
  }

  onMounted(() => void load());

  return { dialog, editingIndex, form, load, loading, open, remove, rows, save, saving };
}
