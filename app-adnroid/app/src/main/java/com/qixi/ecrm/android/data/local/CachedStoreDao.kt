package com.qixi.ecrm.android.data.local

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow

@Dao
interface CachedStoreDao {
    @Query("SELECT * FROM cached_store WHERE app_id = :appId LIMIT 1")
    fun observe(appId: String): Flow<CachedStoreEntity?>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(store: CachedStoreEntity)
}

