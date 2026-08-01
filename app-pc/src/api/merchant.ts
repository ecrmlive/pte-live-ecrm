import { http } from "@/utils/request";

export interface MerchantApplication {
  id: number;
  applicant_user_id: number;
  merchant_name: string;
  contact_name: string;
  contact_mobile: string;
  category_name: string;
  merchant_type: string;
  license_key: string;
  status: "pending" | "approved" | "rejected";
  created_at: string;
}

export function createMerchantApplication(data: Omit<MerchantApplication, "id" | "applicant_user_id" | "status" | "created_at">) {
  return http.post<MerchantApplication>("/merchant-applications", data);
}

export function fetchMyMerchantApplications() {
  return http.get<{ list: MerchantApplication[] }>("/merchant-applications/mine");
}
