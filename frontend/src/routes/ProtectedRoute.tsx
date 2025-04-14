import { Navigate, Outlet } from "react-router";
import { useAppSelector } from "@/lib/hooks/AppHooks";

export default function ProtectedRoute() {
	const token = useAppSelector((s) => s.auth.token);
	return token ? <Outlet /> : <Navigate to="/login" replace />;
}
