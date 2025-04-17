import { Button } from "@/components/atoms/shadCN/button";
import { Input } from "@/components/atoms/shadCN/input";
import { Label } from "@/components/atoms/shadCN/label";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/atoms/shadCN/sheet";

import { useAppSelector } from "@/lib/hooks/AppHooks";
import { Avatar, AvatarFallback } from "../atoms/shadCN/avatar";

export function ProfileSheet() {
  const { user } = useAppSelector((s) => s.auth);
  const initials = user?.email?.charAt(0).toUpperCase() ?? "?";
    return (
      <Sheet>
        <SheetTrigger asChild>
          <Avatar className="h-8 w-8 cursor-pointer hover:opacity-80">
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
        </SheetTrigger>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>User profile</SheetTitle>
            <SheetDescription>
              Profile Information.
            </SheetDescription>
          </SheetHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">
                Id
              </Label>
              <span>{user?.id}</span>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="username" className="text-right">
                Email
              </Label>
              <span>{user?.email}</span>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="username" className="text-right">
                Provider
              </Label>
              <span className="capitalize">{user?.provider}</span>
            </div>
          </div>
          <SheetFooter>

          </SheetFooter>
        </SheetContent>
      </Sheet>
    );
}
