import * as React from "react"
import { cn } from "@/lib/utils"

interface DrawerProps {
  open?: boolean
  onOpenChange?: (open: boolean) => void
  children: React.ReactNode
  className?: string
}

interface DrawerContextValue {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const DrawerContext = React.createContext<DrawerContextValue | undefined>(undefined)

const Drawer = ({ open = false, onOpenChange, children, className }: DrawerProps) => {
  return (
    <DrawerContext.Provider value={{ open, onOpenChange: onOpenChange || (() => {}) }}>
      {children}
    </DrawerContext.Provider>
  )
}

const DrawerTrigger = ({ children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) => {
  const { onOpenChange } = React.useContext(DrawerContext)!

  return (
    <button
      onClick={() => onOpenChange(true)}
      {...props}
    >
      {children}
    </button>
  )
}

const DrawerContent = ({ children, className }: { children: React.ReactNode; className?: string }) => {
  const { open, onOpenChange } = React.useContext(DrawerContext)!

  React.useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onOpenChange(false)
      }
    }

    if (open) {
      document.addEventListener("keydown", handleEscape)
      document.body.style.overflow = "hidden"
    }

    return () => {
      document.removeEventListener("keydown", handleEscape)
      document.body.style.overflow = ""
    }
  }, [open, onOpenChange])

  if (!open) return null

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-50 bg-black/50 transition-opacity"
        onClick={() => onOpenChange(false)}
      />

      {/* Drawer */}
      <div className="fixed inset-y-0 right-0 z-50 w-full sm:w-[600px] lg:w-[800px] bg-white shadow-xl transition-transform">
        <div className={cn("h-full overflow-y-auto", className)}>
          {children}
        </div>
      </div>
    </>
  )
}

const DrawerHeader = ({ children, className }: { children: React.ReactNode; className?: string }) => {
  return (
    <div className={cn("flex items-center justify-between px-6 py-4 border-b", className)}>
      {children}
    </div>
  )
}

const DrawerTitle = ({ children, className }: { children: React.ReactNode; className?: string }) => {
  return (
    <h2 className={cn("text-lg font-semibold", className)}>
      {children}
    </h2>
  )
}

const DrawerBody = ({ children, className }: { children: React.ReactNode; className?: string }) => {
  return (
    <div className={cn("p-6", className)}>
      {children}
    </div>
  )
}

const DrawerFooter = ({ children, className }: { children: React.ReactNode; className?: string }) => {
  return (
    <div className={cn("px-6 py-4 border-t bg-gray-50", className)}>
      {children}
    </div>
  )
}

const DrawerClose = ({ children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) => {
  const { onOpenChange } = React.useContext(DrawerContext)!

  return (
    <button
      onClick={() => onOpenChange(false)}
      {...props}
    >
      {children}
    </button>
  )
}

export {
  Drawer,
  DrawerTrigger,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerBody,
  DrawerFooter,
  DrawerClose,
}
