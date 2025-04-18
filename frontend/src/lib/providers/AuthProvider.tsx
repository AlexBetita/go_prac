import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";
import { fetchProfileThunk } from "@/services/auth/thunks";
import { logout } from "@/lib/store/slices/authSlice";
import { useEffect } from "react";

type Props = {
  children: React.ReactNode;
};

export function AuthProvider({ children }: Props) {
  const dispatch = useAppDispatch();
  const token = useAppSelector((s) => s.auth.token);
  const exp = useAppSelector((s) => s.auth.tokenExp);
  const user = useAppSelector((s) => s.auth.user);

  useEffect(() => {
    if (!token || !exp) return;
    if (exp * 1000 < Date.now()) {
      dispatch(logout());
    } else if (!user) {
      dispatch(fetchProfileThunk());
    }
  }, [token, exp, user, dispatch]);

  return <>{children}</>;
}