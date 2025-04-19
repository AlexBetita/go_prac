import { Navigate, Outlet } from "react-router";
import { useAppSelector } from "@/lib/hooks/AppHooks";

export default function ProtectedRoute() {
	const auth = useAppSelector((s) => s.auth);
	return (auth.token && auth.tokenExp) ? <Outlet /> : <Navigate to="/login" replace />;
}
