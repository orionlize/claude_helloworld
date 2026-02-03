import type { Summary, TrafficPoint, InterfaceItem, InterfaceDetail, SyncLog } from "@/lib/api";
import type { FieldNode } from "@/components/FieldTree";

const requestFields: FieldNode[] = [
  {
    name: "tenantId",
    type: "string",
    required: true,
    description: "租户标识",
  },
  {
    name: "filters",
    type: "object",
    children: [
      {
        name: "env",
        type: "string",
        description: "环境",
      },
      {
        name: "status",
        type: "string",
        description: "状态",
      },
      {
        name: "timeRange",
        type: "object",
        children: [
          { name: "from", type: "string" },
          { name: "to", type: "string" },
        ],
      },
    ],
  },
  {
    name: "pagination",
    type: "object",
    children: [
      { name: "page", type: "number" },
      { name: "pageSize", type: "number" },
    ],
  },
];

const responseFields: FieldNode[] = [
  {
    name: "success",
    type: "boolean",
  },
  {
    name: "data",
    type: "object",
    children: [
      {
        name: "items",
        type: "array",
        children: [
          {
            name: "interfaceId",
            type: "string",
          },
          {
            name: "name",
            type: "string",
          },
          {
            name: "method",
            type: "string",
          },
          {
            name: "status",
            type: "string",
          },
        ],
      },
      {
        name: "total",
        type: "number",
      },
    ],
  },
  {
    name: "traceId",
    type: "string",
  },
];

export const mockSummary: Summary = {
  metrics: [
    {
      id: "interfaces",
      title: "同步接口数",
      value: "1,284",
      description: "本周新增 38",
      trend: "+4.6%",
      trendType: "up",
    },
    {
      id: "sync-rate",
      title: "同步成功率",
      value: "99.92%",
      description: "过去 7 天",
      trend: "+0.12%",
      trendType: "up",
    },
    {
      id: "test-pass",
      title: "测试通过率",
      value: "92.4%",
      description: "最近 120 组",
      trend: "-1.4%",
      trendType: "down",
    },
    {
      id: "doc-visits",
      title: "文档访问量",
      value: "45,678",
      description: "30 天内",
      trend: "+12.5%",
      trendType: "up",
    },
  ],
  status: {
    project: "API 平台核心服务",
    environment: "staging",
    syncState: "healthy",
    lastSync: "2026-02-03 16:42",
  },
};

export const mockTraffic: TrafficPoint[] = [
  { label: "Jan 7", primary: 120, secondary: 90 },
  { label: "Jan 11", primary: 180, secondary: 140 },
  { label: "Jan 15", primary: 160, secondary: 120 },
  { label: "Jan 19", primary: 210, secondary: 170 },
  { label: "Jan 23", primary: 260, secondary: 200 },
  { label: "Jan 27", primary: 230, secondary: 190 },
  { label: "Jan 31", primary: 280, secondary: 220 },
  { label: "Feb 4", primary: 300, secondary: 240 },
  { label: "Feb 8", primary: 270, secondary: 210 },
  { label: "Feb 12", primary: 320, secondary: 250 },
  { label: "Feb 16", primary: 340, secondary: 270 },
  { label: "Feb 20", primary: 360, secondary: 300 },
];

export const mockInterfaces: InterfaceItem[] = [
  {
    id: "intf_001",
    name: "同步接口列表",
    method: "GET",
    path: "/api/interfaces",
    category: "同步层",
    status: "active",
    updatedAt: "2026-02-03 16:40",
  },
  {
    id: "intf_002",
    name: "创建测试套件",
    method: "POST",
    path: "/api/tests/suites",
    category: "测试层",
    status: "active",
    updatedAt: "2026-02-03 15:28",
  },
  {
    id: "intf_003",
    name: "Mock 规则列表",
    method: "GET",
    path: "/api/mock/rules",
    category: "Mock 层",
    status: "active",
    updatedAt: "2026-02-03 14:16",
  },
  {
    id: "intf_004",
    name: "文档分享设置",
    method: "PATCH",
    path: "/api/docs/share",
    category: "文档层",
    status: "deprecated",
    updatedAt: "2026-02-02 19:40",
  },
];

export const mockInterfaceDetail: InterfaceDetail = {
  ...mockInterfaces[0],
  description: "按条件分页查询当前项目同步的接口列表。",
  requestFields,
  responseFields,
};

export const mockSyncLogs: SyncLog[] = [
  {
    id: "log_01",
    time: "16:42",
    action: "接口更新",
    status: "success",
    detail: "同步 12 条接口，2 条分类",
    author: "easy-yapi",
  },
  {
    id: "log_02",
    time: "16:17",
    action: "字段变更",
    status: "warning",
    detail: "检测到响应字段缺少描述",
    author: "easy-yapi",
  },
  {
    id: "log_03",
    time: "15:58",
    action: "接口删除",
    status: "success",
    detail: "已移除废弃接口 1 条",
    author: "easy-yapi",
  },
  {
    id: "log_04",
    time: "15:32",
    action: "同步失败",
    status: "failed",
    detail: "鉴权失败，已重试",
    author: "easy-yapi",
  },
];
