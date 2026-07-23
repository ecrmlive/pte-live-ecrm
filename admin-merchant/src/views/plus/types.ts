export interface PlusPluginItem {
  access_id?: number;
  icon?: null | string;
  name: string;
  path: string;
  redirect_name?: string;
  remark?: string;
  upload_icon?: string;
}

export interface PlusCategory {
  children: PlusPluginItem[];
  name: string;
  plus_category_id?: number;
  sort?: number;
}
