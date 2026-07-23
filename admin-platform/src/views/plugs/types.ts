export interface PlugItem {
  access_id: number;
  icon?: null | string;
  name: string;
  remark?: string;
}

export interface PlugCategory {
  children: PlugItem[];
  name: string;
  plus_category_id: number;
}

export interface PlugCandidate {
  access_id: number;
  name: string;
}
