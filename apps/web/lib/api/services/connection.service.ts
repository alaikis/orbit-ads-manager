import { apiClient } from "@/lib/api/client";

export interface ConnectionField {
  key: string;
  label: string;
  type: "text" | "password" | "number" | "select";
  required?: boolean;
  placeholder?: string;
  hint?: string;
  secret?: boolean;
  options?: Array<{ value: string; label: string }>;
}

export interface PlatformSchema {
  platform: string;
  label: string;
  description: string;
  auth_flow: "api_key" | "developer_token" | "oauth";
  oauth_provider?: string;
  authorize_url?: string;
  token_url?: string;
  client_id_env?: string;
  scopes?: string[];
  fields: ConnectionField[];
  post_auth_actions?: string[];
}

export interface Connection {
  id: number;
  tenant_id: number;
  workspace_id?: number;
  type: string;
  platform: string;
  name: string;
  status: "active" | "error" | "paused";
  scopes?: string[];
  last_refreshed_at?: string;
  last_error?: string;
  created_at: string;
  updated_at: string;
}

export interface OAuthStartResponse {
  authorize_url: string;
  state: string;
  expires_at: string;
}

export const connectionService = {
  list: () => apiClient.get<{ items: Connection[] }>("/connections"),
  get: (id: number) => apiClient.get<Connection>(`/connections/${id}`),
  create: (data: { platform: string; name: string; fields: Record<string, string>; scopes?: string[] }) =>
    apiClient.post<Connection>("/connections", data),
  update: (id: number, data: { name?: string; status?: string; fields?: Record<string, string> }) =>
    apiClient.patch<Connection>(`/connections/${id}`, data),
  delete: (id: number) => apiClient.delete<{ message: string }>(`/connections/${id}`),
  test: (id: number) => apiClient.post<{ ok: boolean; message?: string; [k: string]: unknown }>(`/connections/${id}/test`),
  sync: (id: number) => apiClient.post<{ ok: boolean; candidates?: unknown[]; [k: string]: unknown }>(`/connections/${id}/sync`),
  listPlatforms: () => apiClient.get<{ items: PlatformSchema[] }>("/connections/platforms"),
  oauthStart: (data: { platform: string; name: string; fields?: Record<string, string>; scopes?: string[]; redirect_to?: string }) =>
    apiClient.post<OAuthStartResponse>("/connections/oauth/start", data),
  googleShoppingMerchants: (id: number) =>
    apiClient.get<{ ok: boolean; merchants?: Array<{ id: string; name: string; type: string; admittance: string }> }>(
      `/connections/${id}/google-shopping/merchants`,
    ),
  googleShoppingProducts: (id: number, merchantId: string) =>
    apiClient.get<{ ok: boolean; products?: unknown }>(
      `/connections/${id}/google-shopping/products?merchant_id=${encodeURIComponent(merchantId)}`,
    ),
};
