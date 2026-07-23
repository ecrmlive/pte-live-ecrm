/**
 * Dismiss stale Vben modal overlays left in the DOM after close animations.
 * Safe to call before opening another dialog on the same page.
 */
export function dismissStaleModalOverlays() {
  if (typeof document === 'undefined') return;

  const overlays = [
    ...document.querySelectorAll<HTMLElement>('[data-dismissable-modal].bg-overlay'),
  ];
  if (overlays.length <= 1) return;

  const top = overlays.reduce((best, el) => {
    const z = Number.parseInt(getComputedStyle(el).zIndex || '0', 10);
    const bestZ = Number.parseInt(getComputedStyle(best).zIndex || '0', 10);
    return z >= bestZ ? el : best;
  }, overlays[0]!);

  for (const overlay of overlays) {
    if (overlay !== top) {
      overlay.style.pointerEvents = 'none';
    }
  }
}
