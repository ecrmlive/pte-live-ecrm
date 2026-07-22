import { http } from "@/utils/request";

export interface SpreadMe {
  uid: number;
  spread_uid: number;
  is_promoter: number;
  spread_count: number;
}

export interface SpreadBill {
  bill_id: number;
  number: number;
  balance: number;
  title: string;
  mark?: string;
  create_time?: string;
}

export function fetchSpreadMe() {
  return http.get<SpreadMe>("/spread/me");
}

export function bindSpread(spread_uid: number) {
  return http.post<{ ok?: boolean }>("/spread/bind", { spread_uid });
}

export function fetchSpreadBills(page = 1, limit = 20) {
  return http.get<{ list: SpreadBill[]; total: number }>(
    `/spread/bills?page=${page}&limit=${limit}`
  );
}
