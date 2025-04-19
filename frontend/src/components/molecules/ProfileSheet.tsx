import { Button } from "@/components/atoms/shadCN/button";
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

import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";
import { Avatar, AvatarFallback } from "../atoms/shadCN/avatar";
import { logout } from "@/lib/store/slices/authSlice";

export function ProfileSheet() {
    const dispatch = useAppDispatch()
  
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
            <SheetDescription>Profile Information.</SheetDescription>
          </SheetHeader>
          <div className="grid gap-4 py-4">
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
            <SheetClose asChild>
              <Button type="submit" onClick={() => {
                dispatch(logout())
              }}>Logout</Button>
            </SheetClose>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    );
}
