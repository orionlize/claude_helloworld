import { useMemo, useState } from "react";
import { ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";

export interface FieldNode {
  name: string;
  type: string;
  required?: boolean;
  description?: string;
  children?: FieldNode[];
}

interface FieldTreeProps {
  title: string;
  nodes: FieldNode[];
}

function buildPaths(nodes: FieldNode[], prefix = ""): string[] {
  const paths: string[] = [];
  nodes.forEach((node) => {
    const path = prefix ? `${prefix}.${node.name}` : node.name;
    paths.push(path);
    if (node.children?.length) {
      paths.push(...buildPaths(node.children, path));
    }
  });
  return paths;
}

export function FieldTree({ title, nodes }: FieldTreeProps) {
  const allPaths = useMemo(() => buildPaths(nodes), [nodes]);
  const [open, setOpen] = useState<Set<string>>(new Set(allPaths));
  const [hoveredPath, setHoveredPath] = useState<string | null>(null);

  const toggle = (path: string) => {
    setOpen((prev) => {
      const next = new Set(prev);
      if (next.has(path)) {
        next.delete(path);
      } else {
        next.add(path);
      }
      return next;
    });
  };

  const renderNode = (node: FieldNode, depth: number, parentPath: string) => {
    const path = parentPath ? `${parentPath}.${node.name}` : node.name;
    const hasChildren = Boolean(node.children?.length);
    const isOpen = open.has(path);
    const isRoot = depth === 0;
    const isAncestor = hoveredPath ? hoveredPath.startsWith(path) : false;
    const isDescendant = hoveredPath ? path.startsWith(hoveredPath) : false;
    const isHighlighted = hoveredPath === path;
    const lineTone = "bg-white/10";
    const dotTone = "bg-white/40";

    return (
      <div key={path} className="space-y-1">
        <div
          onMouseEnter={() => setHoveredPath(path)}
          onMouseLeave={() => setHoveredPath(null)}
          data-highlighted={isHighlighted}
          className={cn(
            "group relative flex items-center gap-2 rounded-lg px-2 py-1 text-xs transition",
            "before:absolute before:inset-0 before:rounded-lg before:bg-gradient-to-r before:from-sky-400/15 before:via-white/10 before:to-transparent",
            "before:opacity-0 before:transition hover:before:opacity-100",
            "hover:shadow-[0_0_24px_rgba(56,189,248,0.16)]",
            "text-white/75 hover:bg-white/5 hover:text-white",
            !isRoot && "pl-6"
          )}
          style={{ marginLeft: depth * 12 }}
        >
          {!isRoot && (
            <>
              <span
                className={cn(
                  "absolute left-2 top-0 bottom-0 w-px",
                  lineTone
                )}
              />
              <span className="absolute left-2 top-0 bottom-0 w-[3px] bg-sky-400/35 blur-[2px] opacity-0 group-hover:opacity-100" />
              <span
                className={cn(
                  "absolute left-2 top-1/2 h-px w-3 -translate-y-1/2",
                  lineTone
                )}
              />
              <span
                className={cn(
                  "absolute left-2 top-1/2 h-1.5 w-1.5 -translate-y-1/2 rounded-full",
                  dotTone
                )}
              />
            </>
          )}
          <div className="relative z-10 flex w-full items-center gap-2">
            <button
              type="button"
              onClick={() => hasChildren && toggle(path)}
              className={cn(
                "flex h-5 w-5 items-center justify-center rounded-full border border-white/10 transition",
                "group-hover:border-sky-300/40 group-hover:bg-white/10",
                !hasChildren && "opacity-40"
              )}
            >
              {hasChildren ? (
                <ChevronRight
                  className={cn(
                    "h-3 w-3 transition-transform duration-200",
                    isOpen && "rotate-90"
                  )}
                />
              ) : (
                <span className="h-1.5 w-1.5 rounded-full bg-white/40" />
              )}
            </button>
            <span className="font-medium text-white">{node.name}</span>
            <span className="rounded-full bg-white/10 px-2 py-0.5 text-[10px] uppercase tracking-wide text-white/50">
              {node.type}
            </span>
            {node.required && (
              <span className="text-[10px] text-rose-300">必填</span>
            )}
            {node.description && (
              <span className="ml-auto text-[10px] text-white/40">
                {node.description}
              </span>
            )}
          </div>
        </div>
        {hasChildren && (
          <div
            className={cn(
              "grid transition-[grid-template-rows,opacity,transform] duration-300 ease-out",
              isOpen
                ? "grid-rows-[1fr] opacity-100 translate-y-0"
                : "grid-rows-[0fr] opacity-0 -translate-y-2"
            )}
          >
            <div className="min-h-0">
              <div
                className={cn(
                  "ml-5 border-l pl-4",
                  "border-white/10"
                )}
              >
                <div className="space-y-1">
                  {node.children?.map((child) => renderNode(child, depth + 1, path))}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-white">{title}</div>
          <div className="text-xs text-white/50">字段以树状结构展示，来源于同步数据</div>
        </div>
        <span className="rounded-full border border-white/15 bg-white/5 px-3 py-1 text-[10px] text-white/60">
          只读
        </span>
      </div>
      <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
        <div className="space-y-2">
          {nodes.map((node) => renderNode(node, 0, ""))}
        </div>
      </div>
    </div>
  );
}
