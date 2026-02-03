import { ReactNode } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface StatCardProps {
  title: string;
  value: string;
  description: string;
  badge: string;
  badgeVariant?: "default" | "success" | "warning" | "info";
  icon?: ReactNode;
}

export function StatCard({
  title,
  value,
  description,
  badge,
  badgeVariant = "default",
  icon,
}: StatCardProps) {
  return (
    <Card className="relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-white/5 via-transparent to-transparent" />
      <CardContent className="relative flex h-full flex-col gap-3">
        <div className="flex items-center justify-between">
          <p className="text-sm text-white/60">{title}</p>
          <div className="text-white/50">{icon}</div>
        </div>
        <div className="text-2xl font-semibold text-white">{value}</div>
        <div className="flex items-center justify-between">
          <p className="text-xs text-white/60">{description}</p>
          <Badge variant={badgeVariant}>{badge}</Badge>
        </div>
      </CardContent>
    </Card>
  );
}
