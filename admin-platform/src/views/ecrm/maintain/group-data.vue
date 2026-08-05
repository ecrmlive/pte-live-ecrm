<script setup lang="ts">
import { Page } from '@vben/common-ui';
import { useCacheListCrud } from '#/components/ecrm/useCacheListCrud';
import { fetchMaintainGroupData, saveMaintainGroupData } from '#/api/core/platform-maintain';
import { EcrmListPage } from '#/components/ecrm';

const { dialog, editingIndex, form, loading, open, remove, rows, save, saving } = useCacheListCrud(
  fetchMaintainGroupData,
  saveMaintainGroupData,
);
</script>

<template>
  <Page title="组合数据" description="组合数据 stub CRUD；真实 group_data 域待接入。">
    <EcrmListPage title="组合数据列表">
      <template #actions><el-button type="primary" @click="open()">新增组合</el-button></template>
      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column label="标识" prop="id" min-width="120" />
        <el-table-column label="名称" prop="name" min-width="160" />
        <el-table-column label="备注" prop="remark" min-width="220" show-overflow-tooltip />
        <el-table-column label="启用" width="90"><template #default="{ row }"><el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="150"><template #default="{ row, $index }"><el-button link type="primary" @click="open(row, $index)">编辑</el-button><el-button link type="danger" @click="remove($index, row.name)">删除</el-button></template></el-table-column>
      </el-table>
    </EcrmListPage>
    <el-dialog v-model="dialog" :title="`${editingIndex === undefined ? '新增' : '编辑'}组合数据`" width="520px" destroy-on-close>
      <el-form label-width="84px"><el-form-item label="名称" required><el-input v-model="form.name" maxlength="64" /></el-form-item><el-form-item label="标识"><el-input v-model="form.id" maxlength="64" /></el-form-item><el-form-item label="备注"><el-input v-model="form.remark" maxlength="255" type="textarea" :rows="3" /></el-form-item><el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item></el-form>
      <template #footer><el-button @click="dialog = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
