export interface AccessNode {
  access_id: number;
  api_path?: string;
  children: AccessNode[];
  component?: string;
  create_time?: string | number;
  icon?: string;
  is_menu?: number;
  is_route: number;
  is_show: number;
  name: string;
  parent_id?: number;
  path: string;
  permission_code?: string;
  redirect_name?: string;
  remark?: string;
  sort: number;
}

export type AccessAddType = '' | 'child' | 'copy';

export interface AccessFormModel {
  access_id?: number;
  api_path?: string;
  component?: string;
  icon?: string;
  is_menu: number;
  is_route: number;
  is_show: number;
  name: string;
  parent_id: number;
  path: string;
  permission_code?: string;
  redirect_name?: string;
  remark?: string;
  sort: number;
}
