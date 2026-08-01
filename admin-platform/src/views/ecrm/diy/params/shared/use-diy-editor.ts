import type { InjectionKey } from 'vue';

import { inject } from 'vue';

export interface DiyEditorContext {
  onEditorAddData: () => void;
  onEditorAddImg: () => void;
  onEditorDeleleData: (index: number, selectedIndex: number) => void;
  onEditorDeleleImg: (index: number, selectedIndex: number) => void;
  onEditorResetColor: (
    holder: Record<string, unknown>,
    attribute: string,
    color: string,
  ) => void;
  onEditorSelectImage: (target: Record<string, unknown>, imgUrl: string) => void;
  openProduct: (list: Array<{ product_id?: number }>, islist?: boolean) => void;
  openStore: (islist?: boolean) => void;
}

export const DIY_EDITOR_KEY: InjectionKey<DiyEditorContext> = Symbol('diyEditor');

export function useDiyEditor(): DiyEditorContext {
  const editor = inject(DIY_EDITOR_KEY);
  if (!editor) {
    throw new Error('diyEditor not provided');
  }
  return editor;
}

export function useDiyEditorOptional(): DiyEditorContext | undefined {
  return inject(DIY_EDITOR_KEY);
}
