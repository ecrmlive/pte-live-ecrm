/**
 * Remove teleported Element Plus poppers (date picker, select, etc.) left in body.
 * Call when closing dialogs/modals so orphaned overlays do not block the page.
 */
export function dismissElPoppers(scope?: ParentNode | null) {
  if (typeof document === 'undefined') return;

  const root = scope ?? document.body;

  root.querySelectorAll<HTMLElement>('.el-picker__popper, .el-select__popper').forEach((el) => {
    el.remove();
  });

  root.querySelectorAll<HTMLElement>('.el-popper.is-pure').forEach((el) => {
    if (el.classList.contains('el-picker__popper') || el.classList.contains('el-select__popper')) {
      return;
    }
    if (el.closest('.el-dialog, [data-dismissable-modal], [role="dialog"]')) {
      return;
    }
    el.remove();
  });
}

/** Blur focused inputs inside a container and dismiss any open poppers. */
export function cleanupFormOverlays(container?: ParentNode | null) {
  if (typeof document === 'undefined') return;

  const active = document.activeElement;
  if (active instanceof HTMLElement && container?.contains(active)) {
    active.blur();
  } else if (active instanceof HTMLElement && active.closest('.el-picker-panel, .el-popper')) {
    active.blur();
  }

  dismissElPoppers(container);
  dismissElPoppers(document.body);
}
