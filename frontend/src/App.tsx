import { useEffect, useMemo, useState } from "react";
import {
  Activity,
  ArrowUpRight,
  CheckCircle2,
  GitBranch,
  LineChart as LineChartIcon,
  RefreshCw,
  Sparkles,
} from "lucide-react";
import { Sidebar } from "@/components/Sidebar";
import { StatCard } from "@/components/StatCard";
import { LineChart, type SeriesPoint } from "@/components/LineChart";
import { FieldTree } from "@/components/FieldTree";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { api, type InterfaceDetail, type InterfaceItem, type SyncLog, type Summary, type TrafficPoint } from "@/lib/api";
import { mockInterfaceDetail, mockInterfaces, mockSummary, mockSyncLogs, mockTraffic } from "@/lib/mock";
import { cn } from "@/lib/utils";

const trendVariant = (trend: Summary["metrics"][number]["trendType"]) => {
  if (trend === "up") return "success";
  if (trend === "down") return "warning";
  return "default";
};

const methodTone = (method: string) => {
  const normalized = method.toUpperCase();
  if (normalized === "GET") return "info";
  if (normalized === "POST") return "success";
  if (normalized === "PATCH") return "warning";
  return "default";
};

function formatSyncState(state: Summary["status"]["syncState"]) {
  switch (state) {
    case "healthy":
      return { label: "同步健康", variant: "success" } as const;
    case "warning":
      return { label: "存在警告", variant: "warning" } as const;
    default:
      return { label: "同步异常", variant: "default" } as const;
  }
}

export default function App() {
  const [summary, setSummary] = useState<Summary>(mockSummary);
  const [traffic, setTraffic] = useState<TrafficPoint[]>(mockTraffic);
  const [interfaces, setInterfaces] = useState<InterfaceItem[]>(mockInterfaces);
  const [detail, setDetail] = useState<InterfaceDetail>(mockInterfaceDetail);
  const [syncLogs, setSyncLogs] = useState<SyncLog[]>(mockSyncLogs);
  const [range, setRange] = useState("90d");

  useEffect(() => {
    const load = async () => {
      try {
        const [summaryData, interfacesData, logsData] = await Promise.all([
          api.summary(),
          api.interfaces(),
          api.syncLogs(),
        ]);
        setSummary(summaryData);
        setInterfaces(interfacesData);
        setSyncLogs(logsData);
        if (interfacesData[0]) {
          const detailData = await api.interfaceDetail(interfacesData[0].id);
          setDetail(detailData);
        }
      } catch (error) {
        console.warn("API unavailable, fallback to mock data", error);
      }
    };
    load();
  }, []);

  useEffect(() => {
    const loadTraffic = async () => {
      try {
        const trafficData = await api.traffic(range);
        setTraffic(trafficData);
      } catch (error) {
        console.warn("Traffic API unavailable", error);
      }
    };
    loadTraffic();
  }, [range]);

  const seriesData: SeriesPoint[] = useMemo(
    () => traffic.map((item) => ({ label: item.label, primary: item.primary, secondary: item.secondary })),
    [traffic]
  );

  const syncState = formatSyncState(summary.status.syncState);

  return (
    <div className="relative min-h-screen px-4 py-6 md:px-8">
      <div className="pointer-events-none fixed inset-0 -z-10">
        <div className="absolute left-16 top-10 h-40 w-40 rounded-full bg-sky-500/20 blur-[90px]" />
        <div className="absolute right-20 top-24 h-40 w-40 rounded-full bg-indigo-500/10 blur-[100px]" />
        <div className="absolute bottom-10 left-1/3 h-56 w-56 rounded-full bg-emerald-400/10 blur-[120px]" />
        <div className="absolute inset-0 opacity-[0.2] soft-grid" />
      </div>

      <div className="grid gap-6 md:grid-cols-[260px_1fr]">
        <Sidebar />

        <main className="space-y-6">
          <section className="glass-panel rounded-3xl p-6">
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.3em] text-white/40">Overview</p>
                <h1 className="text-2xl font-semibold text-white">接口协作控制台</h1>
                <p className="mt-2 text-sm text-white/60">
                  {summary.status.project} · {summary.status.environment} 环境
                </p>
              </div>
              <div className="flex flex-wrap items-center gap-3">
                <Badge variant={syncState.variant}>{syncState.label}</Badge>
                <Badge variant="info">最近同步：{summary.status.lastSync}</Badge>
                <Button variant="outline">
                  <RefreshCw className="h-4 w-4" />
                  触发同步
                </Button>
              </div>
            </div>
          </section>

          <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            {summary.metrics.map((metric) => (
              <StatCard
                key={metric.id}
                title={metric.title}
                value={metric.value}
                description={metric.description}
                badge={metric.trend}
                badgeVariant={trendVariant(metric.trendType)}
                icon={<ArrowUpRight className="h-4 w-4" />}
              />
            ))}
          </section>

          <section>
            <Card className="p-6">
              <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div>
                  <p className="text-sm font-semibold text-white">同步流量趋势</p>
                  <p className="text-xs text-white/50">展示 easy-yapi 推送后的接口变更节奏</p>
                </div>
                <Tabs value={range} onValueChange={setRange}>
                  <TabsList>
                    <TabsTrigger value="90d">近 90 天</TabsTrigger>
                    <TabsTrigger value="30d">近 30 天</TabsTrigger>
                    <TabsTrigger value="7d">近 7 天</TabsTrigger>
                  </TabsList>
                </Tabs>
              </div>
              <div className="mt-6">
                <LineChart data={seriesData} />
              </div>
              <div className="mt-6 grid gap-4 sm:grid-cols-3">
                <div className="glass-card rounded-2xl p-4">
                  <div className="flex items-center gap-2 text-xs text-white/60">
                    <Activity className="h-4 w-4" />
                    今日推送
                  </div>
                  <div className="mt-2 text-2xl font-semibold text-white">128</div>
                  <div className="text-xs text-white/40">较昨日 +12%</div>
                </div>
                <div className="glass-card rounded-2xl p-4">
                  <div className="flex items-center gap-2 text-xs text-white/60">
                    <LineChartIcon className="h-4 w-4" />
                    同步失败
                  </div>
                  <div className="mt-2 text-2xl font-semibold text-white">3</div>
                  <div className="text-xs text-white/40">已自动重试 2 次</div>
                </div>
                <div className="glass-card rounded-2xl p-4">
                  <div className="flex items-center gap-2 text-xs text-white/60">
                    <CheckCircle2 className="h-4 w-4" />
                    规范检测
                  </div>
                  <div className="mt-2 text-2xl font-semibold text-white">92%</div>
                  <div className="text-xs text-white/40">字段描述完整度</div>
                </div>
              </div>
            </Card>
          </section>

          <section className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
            <Card className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-semibold text-white">接口文档 · 树状字段</p>
                  <p className="text-xs text-white/50">字段全部只读，来源于 easy-yapi 推送</p>
                </div>
                <Button variant="ghost" size="sm">
                  <Sparkles className="h-4 w-4" />
                  生成用例
                </Button>
              </div>
              <div className="mt-6 grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
                <div className="space-y-3">
                  <p className="text-xs uppercase tracking-[0.2em] text-white/40">接口列表</p>
                  <div className="space-y-2">
                    {interfaces.map((item) => (
                      <button
                        key={item.id}
                        onClick={async () => {
                          try {
                            const detailData = await api.interfaceDetail(item.id);
                            setDetail(detailData);
                          } catch (error) {
                            setDetail((prev) => (prev.id === item.id ? prev : mockInterfaceDetail));
                          }
                        }}
                        className={cn(
                          "flex w-full flex-col gap-2 rounded-2xl border border-white/10 bg-white/5 p-3 text-left transition",
                          detail.id === item.id ? "ring-2 ring-white/20" : "hover:bg-white/10"
                        )}
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium text-white">{item.name}</span>
                          <Badge variant={methodTone(item.method)}>{item.method}</Badge>
                        </div>
                        <div className="text-xs text-white/50">{item.path}</div>
                        <div className="flex items-center justify-between text-[11px] text-white/40">
                          <span>{item.category}</span>
                          <span>{item.updatedAt}</span>
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
                <div className="space-y-4">
                  <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="text-sm font-semibold text-white">{detail.name}</div>
                        <div className="text-xs text-white/50">
                          {detail.method} {detail.path}
                        </div>
                      </div>
                      <Badge variant={detail.status === "active" ? "success" : "warning"}>
                        {detail.status === "active" ? "Active" : "Deprecated"}
                      </Badge>
                    </div>
                    <div className="mt-3 text-xs text-white/60">{detail.description}</div>
                  </div>
                  <Tabs defaultValue="request">
                    <TabsList className="w-full justify-start">
                      <TabsTrigger value="request">请求字段</TabsTrigger>
                      <TabsTrigger value="response">响应字段</TabsTrigger>
                    </TabsList>
                    <TabsContent value="request">
                      <FieldTree title="Request Schema" nodes={detail.requestFields} />
                    </TabsContent>
                    <TabsContent value="response">
                      <FieldTree title="Response Schema" nodes={detail.responseFields} />
                    </TabsContent>
                  </Tabs>
                </div>
              </div>
            </Card>

            <Card className="p-6">
              <CardHeader className="p-0">
                <CardTitle>同步日志</CardTitle>
                <p className="text-xs text-white/50">最近 24 小时推送情况</p>
              </CardHeader>
              <CardContent className="mt-4 p-0">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>时间</TableHead>
                      <TableHead>动作</TableHead>
                      <TableHead>状态</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {syncLogs.map((log) => (
                      <TableRow key={log.id}>
                        <TableCell>{log.time}</TableCell>
                        <TableCell>
                          <div className="text-sm text-white">{log.action}</div>
                          <div className="text-xs text-white/50">{log.detail}</div>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              log.status === "success"
                                ? "success"
                                : log.status === "warning"
                                ? "warning"
                                : "default"
                            }
                          >
                            {log.status}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </section>

          <section className="grid gap-4 md:grid-cols-3">
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-white/10 p-2">
                  <GitBranch className="h-5 w-5 text-white/70" />
                </div>
                <div>
                  <div className="text-sm text-white">环境变量</div>
                  <div className="text-xs text-white/50">18 组全局变量</div>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-white/10 p-2">
                  <CheckCircle2 className="h-5 w-5 text-white/70" />
                </div>
                <div>
                  <div className="text-sm text-white">测试套件</div>
                  <div className="text-xs text-white/50">7 套正在执行</div>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-white/10 p-2">
                  <Sparkles className="h-5 w-5 text-white/70" />
                </div>
                <div>
                  <div className="text-sm text-white">AI 规范检测</div>
                  <div className="text-xs text-white/50">本周建议 12 条</div>
                </div>
              </div>
            </Card>
          </section>
        </main>
      </div>
    </div>
  );
}
