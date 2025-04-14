import { BrowserRouter, Routes, Route } from "react-router";
import ProtectedRoute from "./ProtectedRoute";
import LoginPage from "@/pages/LoginPage";
import RegisterPage from "@/pages/RegisterPage";
import ProfilePage from "@/pages/ProfilePage";
import AppLayout from "@/pages/App";

export default function AppRoutes() {
	return (
		<BrowserRouter>
			<Routes>
                <Route element={<AppLayout />}>
				    <Route path="/login" element={<LoginPage />} />
                    <Route path="/register" element={<RegisterPage />} />

                    <Route element={<ProtectedRoute />}>
                        <Route path="/profile" element={<ProfilePage />} />
                        <Route path="/" element={<ProfilePage />} />
                    </Route>
                </Route>
			</Routes>
		</BrowserRouter>
	);
}
