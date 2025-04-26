import { Routes, Route, Navigate } from "react-router";
import ProtectedRoute from "./ProtectedRoute";
import LoginPage from "@/pages/LoginPage";
import RegisterPage from "@/pages/RegisterPage";
// import ProfilePage from "@/pages/ProfilePage";
import AppLayout from "@/pages/App";
import HomePage from "@/pages/HomePage";
import PostPage from "@/pages/PostPage";

export default function BaseRoutes() {
	return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/posts/:slug" element={<PostPage />} />

        <Route element={<ProtectedRoute />}>
          {/* <Route path="/profile" element={<ProfilePage />} /> */}
          <Route path="/" element={<HomePage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
