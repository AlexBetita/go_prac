import { Link } from "react-router";
import { useAppSelector } from "@/lib/hooks/AppHooks";
import { Button } from "@/components/atoms/shadCN/button";
import { ProfileSheet } from "@/components/molecules/profile/ProfileSheet";
import { Menu } from "lucide-react";
import { useSidebar } from "@/lib/context/SideBarContext";

export default function Navbar() {
  const { user } = useAppSelector((s) => s.auth);
  const { open, toggle } = useSidebar();

  return (
    <>
      <header className="w-full border-b border-border bg-background shadow-sm">
        <div className="mx-auto flex items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="icon" onClick={toggle}>
              <Menu className="h-6 w-6" />
            </Button>
            <Link to="/" className="text-lg font-semibold">
              MyApp
            </Link>
          </div>

          <div className="flex items-center gap-4">
            {user ? (
              <ProfileSheet />
            ) : (
              <Button asChild variant="outline" size="sm">
                <Link to="/login">Login</Link>
              </Button>
            )}
          </div>
        </div>
      </header>

      {open && (
        <div
          className="fixed left-0 top-16 z-40 h-[calc(100vh-4rem)]"
          onClick={toggle}
        />
      )}
    </>
  );
}
