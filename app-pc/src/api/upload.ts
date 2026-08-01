import { http, ApiError } from "@/utils/request";

export interface COSPutIntent {
  object_key: string;
  upload_url: string;
  method: "PUT";
  headers: Record<string, string>;
  expires_at: string;
}

const LICENSE_PURPOSE = "merchant_application_license";

export async function uploadMerchantLicense(file: File): Promise<string> {
  const intent = await http.post<COSPutIntent>("/uploads/presign", {
    filename: file.name,
    content_type: file.type,
    size: file.size,
    purpose: LICENSE_PURPOSE,
  });
  const upload = await fetch(intent.upload_url, {
    method: intent.method,
    headers: intent.headers,
    body: file,
  });
  if (!upload.ok) {
    throw new ApiError(upload.status, "营业执照上传至 COS 失败");
  }
  const completed = await http.post<{ object_key: string }>("/uploads/complete", {
    object_key: intent.object_key,
    purpose: LICENSE_PURPOSE,
  });
  return completed.object_key;
}
