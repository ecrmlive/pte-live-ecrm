package com.qixi.ecrm.android.presentation.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable

private val LightColors = lightColorScheme(
    primary = BrandBlue,
    secondary = BrandBlue,
    background = SurfaceBackground,
    surface = SurfaceBackground,
)

private val DarkColors = darkColorScheme(
    primary = BrandBlue,
    secondary = BrandBlue,
)

@Composable
fun EcrmTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = LightColors,
        typography = EcrmTypography,
        content = content,
    )
}

