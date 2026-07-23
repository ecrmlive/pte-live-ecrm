/**
 * Shared Playwright cleanup for merchant-admin E2E suites.
 * Use between flows / plugins / chained suite steps to drop overlays and poppers.
 */
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';

export async function dismissElPoppers(page) {
  await page.evaluate(() => {
    const removePoppers = (root) => {
      root.querySelectorAll('.el-picker__popper, .el-select__popper').forEach((el) => {
        el.remove();
      });
      root.querySelectorAll('.el-popper.is-pure').forEach((el) => {
        if (el.classList.contains('el-picker__popper') || el.classList.contains('el-select__popper')) {
          return;
        }
        if (el.closest('.el-dialog, [data-dismissable-modal], [role="dialog"]')) {
          return;
        }
        el.remove();
      });
    };
    removePoppers(document.body);
  });
}

export async function countOrphanPoppers(page) {
  return page.evaluate(() => {
    const isVisible = (el) => {
      const style = window.getComputedStyle(el);
      if (style.display === 'none' || style.visibility === 'hidden') return false;
      const rect = el.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    };
    return [
      ...document.querySelectorAll(
        '.el-picker__popper, .el-select__popper, body > .el-popper',
      ),
    ].filter(isVisible).length;
  });
}

export async function closeTopOverlays(page) {
  for (let attempt = 0; attempt < 4; attempt += 1) {
    const cancel = page
      .locator(
        '[role="dialog"] button, .vben-modal button, .el-dialog button, .el-message-box button',
      )
      .filter({ hasText: /^取消$/ })
      .first();
    if (await cancel.isVisible().catch(() => false)) {
      await cancel.click();
      await page.waitForTimeout(500);
    } else {
      await page.keyboard.press('Escape');
      await page.waitForTimeout(400);
    }

    const poppers = await countOrphanPoppers(page);
    const dialogOpen = await page
      .locator('.el-overlay:not([style*="display: none"])')
      .count()
      .catch(() => 0);
    if (poppers === 0 && dialogOpen === 0) {
      await dismissElPoppers(page);
      return;
    }
  }
  await dismissElPoppers(page);
}

/**
 * Navigate home, close dialogs, dismiss teleported poppers.
 * @param {import('playwright').Page} page
 * @param {{ reload?: boolean }} [options]
 */
export async function resetMerchantPage(page, options = {}) {
  const { reload = false } = options;

  await closeTopOverlays(page);
  await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' }).catch(() => {});
  await page.waitForTimeout(400);
  await closeTopOverlays(page);
  await dismissElPoppers(page);

  if (reload) {
    await page.reload({ waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(600);
    await closeTopOverlays(page);
    await dismissElPoppers(page);
  }
}
