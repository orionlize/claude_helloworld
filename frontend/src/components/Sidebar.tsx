import {
  Activity,
  Blocks,
  BookOpen,
  ClipboardCheck,
  DatabaseZap,
  Layers,
  Settings,
  ShieldCheck,
  Sparkles,
  TestTube2,
  Users,
} from "lucide-react";

const navItems = [
  { label: "总览", icon: Activity },
  { label: "接口调试", icon: Layers },
  { label: "测试套件", icon: TestTube2 },
  { label: "Mock 服务", icon: DatabaseZap },
  { label: "文档中心", icon: BookOpen },
  { label: "团队协作", icon: Users },
  { label: "合规审计", icon: ShieldCheck },
  { label: "AI 助手", icon: Sparkles },
];

const quickItems = [
  { label: "快速创建", icon: Blocks },
  { label: "同步日志", icon: ClipboardCheck },
];

export function Sidebar() {
  return (
    <aside className="glass-panel hidden h-[calc(100vh-48px)] w-full flex-col rounded-3xl p-6 md:flex">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-white/10">
          <span className="text-sm font-semibold text-white">AP</span>
        </div>
        <div>
          <div className="text-sm font-semibold text-white">API 管理平台</div>
          <div className="text-xs text-white/50">easy-yapi 同步</div>
        </div>
      </div>

      <div className="mt-8 space-y-2">
        {quickItems.map((item) => (
          <button
            key={item.label}
            className="flex w-full items-center gap-3 rounded-xl bg-white/10 px-3 py-2 text-left text-sm text-white/80 transition hover:bg-white/20"
          >
            <item.icon className="h-4 w-4" />
            {item.label}
          </button>
        ))}
      </div>

      <div className="mt-8 space-y-1">
        {navItems.map((item, index) => (
          <button
            key={item.label}
            className={`flex w-full items-center gap-3 rounded-xl px-3 py-2 text-left text-sm transition ${
              index === 0
                ? "bg-white/10 text-white"
                : "text-white/60 hover:bg-white/5 hover:text-white"
            }`}
          >
            <item.icon className="h-4 w-4" />
            {item.label}
          </button>
        ))}
      </div>

      <div className="mt-auto pt-8">
        <div className="flex items-center gap-3 rounded-2xl bg-white/5 p-3">
          <div className="h-10 w-10 rounded-full bg-gradient-to-br from-sky-400/70 to-purple-400/60" />
          <div>
            <div className="text-sm text-white">shadcn</div>
            <div className="text-xs text-white/50">owner@example.com</div>
          </div>
          <Settings className="ml-auto h-4 w-4 text-white/40" />
        </div>
      </div>
    </aside>
  );
}
