package com.qixi.ecrm.android.data.local

import androidx.room.ColumnInfo
import androidx.room.Entity

@Entity(tableName = "cached_store", primaryKeys = ["app_id"])
data class CachedStoreEntity(
    @ColumnInfo(name = "app_id") val appId: String,
    @ColumnInfo(name = "name") val name: String,
)

