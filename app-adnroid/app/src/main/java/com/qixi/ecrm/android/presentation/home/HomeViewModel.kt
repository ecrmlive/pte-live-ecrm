package com.qixi.ecrm.android.presentation.home

import androidx.lifecycle.ViewModel
import com.qixi.ecrm.android.domain.model.HomeEntry
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

data class HomeUiState(
    val greeting: String = "欢迎来到七禧商城",
    val entries: List<HomeEntry> = listOf(
        HomeEntry("商品分类", "发现精选好物"),
        HomeEntry("店铺街", "浏览优质商户"),
        HomeEntry("限时活动", "领取专属优惠"),
        HomeEntry("客服帮助", "订单与售后咨询"),
    ),
)

@HiltViewModel
class HomeViewModel @Inject constructor() : ViewModel() {
    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()
}

