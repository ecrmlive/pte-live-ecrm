import { requestClient } from '#/api/request';

export interface SignupFormField {
  key: string;
  label: string;
  type: string;
}

export interface SignupFormOption {
  fields: SignupFormField[];
  id: number;
  name: string;
}

export interface SignupActivity {
  activity_status: number;
  activity_status_text: string;
  color: string;
  cover_url: string;
  created_at: string;
  ends_at: string;
  form_fields?: SignupFormField[];
  form_id: number;
  form_name: string;
  id: number;
  info: string;
  name: string;
  poster_url: string;
  quota: number;
  signup_count_text: string;
  sort: number;
  starts_at: string;
  status: number;
  total: number;
  updated_at: string;
}

export interface SignupActivityPage {
  limit: number;
  list: SignupActivity[];
  page: number;
  total: number;
}

export interface SignupActivityInput {
  color?: string;
  cover_url: string;
  ends_at: string;
  form_id: number;
  info?: string;
  name: string;
  poster_url: string;
  quota?: number;
  sort?: number;
  starts_at: string;
  status?: number;
}

export interface SignupActivityListParams {
  activity_status?: number | string;
  form_id?: number;
  keyword?: string;
  limit: number;
  page: number;
}

export interface SignupRecord {
  activity_id: number;
  avatar: string;
  created_at: string;
  form_cols: Record<string, string>;
  form_value: Record<string, unknown>;
  id: number;
  index: number;
  mobile: string;
  nickname: string;
  user_id: number;
}

export interface SignupRecordPage {
  fields: SignupFormField[];
  limit: number;
  list: SignupRecord[];
  page: number;
  total: number;
}

export interface SignupExportResult {
  content: string;
  file_name: string;
}

export function listSignupFormOptionsApi() {
  return requestClient.get<{ list: SignupFormOption[] }>(
    '/marketing/applications/form-options',
  );
}

export function listSignupActivitiesApi(params: SignupActivityListParams) {
  return requestClient.get<SignupActivityPage>('/marketing/applications', {
    params,
  });
}

export function getSignupActivityApi(id: number) {
  return requestClient.get<SignupActivity>(`/marketing/applications/${id}`);
}

export function createSignupActivityApi(body: SignupActivityInput) {
  return requestClient.post<SignupActivity>('/marketing/applications', body);
}

export function updateSignupActivityApi(id: number, body: SignupActivityInput) {
  return requestClient.put<SignupActivity>(`/marketing/applications/${id}`, body);
}

export function setSignupActivityStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(
    `/marketing/applications/${id}/status`,
    { status },
  );
}

export function deleteSignupActivityApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/marketing/applications/${id}`);
}

export function listSignupRecordsApi(
  id: number,
  params: { keyword?: string; limit: number; page: number },
) {
  return requestClient.get<SignupRecordPage>(
    `/marketing/applications/${id}/users`,
    { params },
  );
}

export function exportSignupRecordsApi(id: number, params?: { keyword?: string }) {
  return requestClient.get<SignupExportResult>(
    `/marketing/applications/${id}/users/export`,
    { params },
  );
}
