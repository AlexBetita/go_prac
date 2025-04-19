import { Link } from "react-router";
import { useAppSelector } from "@/lib/hooks/AppHooks";

import { Button } from "@/components/atoms/shadCN/button";
import { ProfileSheet } from "../molecules/ProfileSheet";

export default function Navbar() {
  const { user } = useAppSelector((s) => s.auth);
  
  return (
    <header className="w-full border-b border-border bg-background shadow-sm">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
        <Link to="/" className="text-lg font-semibold">
          MyApp
        </Link>

        <div className="flex items-center gap-4">
          <nav className="hidden md:flex gap-4 text-sm">
            <Link to="/" className="hover:underline">
              Home
            </Link>
          </nav>

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
  );
}
