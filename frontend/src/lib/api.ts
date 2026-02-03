import type { FieldNode } from "@/components/FieldTree";

export type TrendType = "up" | "down" | "neutral";

export interface Metric {
  id: string;
  title: string;
  value: string;
  description: string;
  trend: string;
  trendType: TrendType;
}

export interface Summary {
  metrics: Metric[];
  status: {
    project: string;
    environment: string;
    syncState: "healthy" | "warning" | "error";
    lastSync: string;
  };
}

export interface TrafficPoint {
  label: string;
  primary: number;
  secondary: number;
}

export interface InterfaceItem {
  id: string;
  name: string;
  method: string;
  path: string;
  category: string;
  status: "active" | "deprecated";
  updatedAt: string;
}

export interface InterfaceDetail extends InterfaceItem {
  description: string;
  requestFields: FieldNode[];
  responseFields: FieldNode[];
}

export interface SyncLog {
  id: string;
  time: string;
  action: string;
  status: "success" | "failed" | "warning";
  detail: string;
  author: string;
}

const API_BASE = import.meta.env.VITE_API_BASE ?? "/api";

async function fetchJson<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`);
  if (!response.ok) {
    throw new Error(`API request failed: ${response.status}`);
  }
  return (await response.json()) as T;
}

export const api = {
  summary: () => fetchJson<Summary>("/summary"),
  traffic: (range: string) => fetchJson<TrafficPoint[]>(`/traffic?range=${range}`),
  interfaces: () => fetchJson<InterfaceItem[]>("/interfaces"),
  interfaceDetail: (id: string) => fetchJson<InterfaceDetail>(`/interfaces/${id}`),
  syncLogs: () => fetchJson<SyncLog[]>("/sync-logs"),
};
