import { useMemo, useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
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

    return (
      <div key={path} className="space-y-1">
        <div
          className={cn(
            "flex items-center gap-2 rounded-lg px-2 py-1 text-xs",
            "text-white/75 hover:bg-white/5"
          )}
          style={{ marginLeft: depth * 12 }}
        >
          <button
            type="button"
            onClick={() => hasChildren && toggle(path)}
            className={cn(
              "flex h-5 w-5 items-center justify-center rounded-full border border-white/10",
              !hasChildren && "opacity-40"
            )}
          >
            {hasChildren ? (
              isOpen ? (
                <ChevronDown className="h-3 w-3" />
              ) : (
                <ChevronRight className="h-3 w-3" />
              )
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
        {hasChildren && isOpen && (
          <div className="ml-3 border-l border-white/10 pl-3">
            <div className="space-y-1">
              {node.children?.map((child) => renderNode(child, depth + 1, path))}
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
