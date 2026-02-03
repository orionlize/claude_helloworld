import { useMemo } from "react";

interface SeriesPoint {
  label: string;
  primary: number;
  secondary?: number;
}

interface LineChartProps {
  data: SeriesPoint[];
}

function buildPath(values: number[], width: number, height: number, padding: number) {
  if (values.length === 0) {
    return { line: "", area: "" };
  }
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;
  const step = width / Math.max(values.length - 1, 1);
  const points = values.map((value, index) => {
    const x = index * step;
    const y = height - padding - ((value - min) / range) * (height - padding * 2);
    return { x, y };
  });
  const line = points.map((point, idx) => `${idx === 0 ? "M" : "L"}${point.x.toFixed(2)} ${point.y.toFixed(2)}`).join(" ");
  const area = `${line} L ${width} ${height - padding} L 0 ${height - padding} Z`;
  return { line, area };
}

export function LineChart({ data }: LineChartProps) {
  const width = 100;
  const height = 40;
  const padding = 4;

  const { primaryPath, primaryArea, secondaryPath } = useMemo(() => {
    const primaryValues = data.map((item) => item.primary);
    const secondaryValues = data.map((item) => item.secondary ?? 0);
    const primary = buildPath(primaryValues, width, height, padding);
    const secondary = buildPath(secondaryValues, width, height, padding);
    return {
      primaryPath: primary.line,
      primaryArea: primary.area,
      secondaryPath: secondary.line,
    };
  }, [data]);

  return (
    <div className="relative w-full">
      <svg viewBox={`0 0 ${width} ${height}`} className="h-64 w-full">
        <defs>
          <linearGradient id="areaGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="rgba(255,255,255,0.35)" />
            <stop offset="70%" stopColor="rgba(255,255,255,0.08)" />
            <stop offset="100%" stopColor="rgba(255,255,255,0)" />
          </linearGradient>
        </defs>
        <g className="text-white/10">
          {[8, 16, 24, 32].map((y) => (
            <line key={y} x1={0} y1={y} x2={width} y2={y} stroke="currentColor" strokeWidth={0.2} />
          ))}
        </g>
        <path d={primaryArea} fill="url(#areaGradient)" />
        <path d={secondaryPath} fill="none" stroke="rgba(148,163,184,0.5)" strokeWidth={0.6} />
        <path d={primaryPath} fill="none" stroke="rgba(255,255,255,0.9)" strokeWidth={0.8} />
      </svg>
      <div className="mt-2 flex justify-between text-xs text-white/40">
        {data.map((item, index) => (
          <span key={item.label} className={index % 3 === 0 ? "" : "opacity-0 md:opacity-100"}>
            {item.label}
          </span>
        ))}
      </div>
    </div>
  );
}

export type { SeriesPoint };
