import { useState, useEffect } from "react";
import { NavLink } from "react-router";
import {
  Home,
  Brain,
  Sparkles,
  MessageSquare,
  FileText,
  History,
  Settings,
  X,
} from "lucide-react";
import { Button } from "@/components/atoms/shadCN/button";
import { useSidebar } from "@/lib/context/SideBarContext";

const menuItems = [
  { title: "Home", path: "/", icon: Home },
  { title: "Models", path: "/models", icon: Brain },
  { title: "Prompts", path: "/prompts", icon: Sparkles },
  { title: "Chat", path: "/chat", icon: MessageSquare },
  { title: "Posts", path: "/posts", icon: FileText },
  { title: "History", path: "/history", icon: History },
  { title: "Settings", path: "/settings", icon: Settings },
];

export default function AppSidebar() {
  const { open, close } = useSidebar();

  const [isMobile, setIsMobile] = useState(false);
  useEffect(() => {
    const mq = window.matchMedia("(max-width: 768px)");
    const onChange = (e: MediaQueryListEvent) => setIsMobile(e.matches);
    setIsMobile(mq.matches);
    mq.addEventListener("change", onChange);
    return () => mq.removeEventListener("change", onChange);
  }, []);

  const container = [
    "fixed left-0 z-50 bg-background text-foreground border-r shadow-md",
    "transform transition-transform duration-300 ease-in-out",
    open ? "translate-x-0" : "-translate-x-full",
    isMobile
      ? "top-0 w-screen h-screen"
      : "top-16 w-64 h-[calc(100vh-4rem)]",
  ].join(" ");

  const navPadding = isMobile ? "pt-16 px-4" : "py-8 px-4";
  const navGap = isMobile ? "gap-3" : "gap-4";
  const iconSize = isMobile ? "h-5 w-5" : "h-6 w-6";
  const textSize = isMobile ? "text-base" : "text-lg";

  return (
    <div className={container}>
      {isMobile && (
        <div className="absolute top-4 right-4">
          <Button variant="ghost" size="icon" onClick={close}>
            <X className="h-6 w-6" />
          </Button>
        </div>
      )}

      <nav className={["flex flex-col", navGap, navPadding].join(" ")}>
        {menuItems.map(({ title, path, icon: Icon }) => (
          <NavLink
            key={path}
            to={path}
            onClick={close}
            className={({ isActive }) =>
              [
                "flex items-center rounded-md px-3 py-2 transition-colors",
                textSize,
                isActive
                  ? "bg-accent text-accent-foreground"
                  : "hover:bg-accent hover:text-accent-foreground",
              ].join(" ")
            }
          >
            <Icon className={iconSize} />
            <span className="pl-2">{title}</span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}
