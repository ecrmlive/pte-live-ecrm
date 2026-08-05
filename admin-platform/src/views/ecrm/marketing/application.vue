<script setup lang="ts">
import { Page } from '@vben/common-ui';
import { useCacheListCrud } from '#/components/ecrm/useCacheListCrud';
import { fetchMarketingApplications, saveMarketingApplications } from '#/api/core/platform-marketing-cache';
const { dialog, editingIndex, form, loading, open, remove, rows, save, saving } = useCacheListCrud(fetchMarketingApplications, saveMarketingApplications);
</script>
<template>
  <Page title="活动报名" description="基于 setting_cache 的 stub 配置；真实报名业务待域 API 接入。">
    <template #extra><el-button type="primary" @click="open()">新增报名</el-button></template>
    <el-alert class="mb-4" title="当前为监管 stub 列表，写入经 typed setting_cache 保存，不含密钥。" type="info" :closable="false" />
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="id"><el-table-column label="标识" prop="id" min-width="120" /><el-table-column label="名称" prop="name" min-width="160" /><el-table-column label="备注" prop="remark" min-width="220" show-overflow-tooltip /><el-table-column label="启用" width="90"><template #default="{ row }"><el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row, $index }"><el-button link type="primary" @click="open(row, $index)">编辑</el-button><el-button link type="danger" @click="remove($index, row.name)">删除</el-button></template></el-table-column></el-table></el-card>
    <el-dialog v-model="dialog" title="活动报名" width="520px" destroy-on-close><el-form label-width="84px"><el-form-item label="名称" required><el-input v-model="form.name" maxlength="64" /></el-form-item><el-form-item label="标识"><el-input v-model="form.id" maxlength="64" /></el-form-item><el-form-item label="备注"><el-input v-model="form.remark" maxlength="255" type="textarea" :rows="3" /></el-form-item><el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item></el-form><template #footer><el-button @click="dialog = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
