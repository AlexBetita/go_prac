import { Routes, Route, Navigate } from "react-router";
import ProtectedRoute from "./ProtectedRoute";
import LoginPage from "@/pages/LoginPage";
import RegisterPage from "@/pages/RegisterPage";
import HomePage from "@/pages/HomePage";
import PostPage from "@/pages/PostPage";

import App from "@/pages/App";

// import ModelsPage from "@/pages/ModelsPage";
// import PromptsPage from "@/pages/PromptsPage";
// import ChatPage from "@/pages/ChatPage";
// import PostsPage from "@/pages/PostsPage";
// import HistoryPage from "@/pages/HistoryPage";
// import SettingsPage from "@/pages/SettingsPage";

export default function BaseRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route element={<App />}>
        <Route path="/posts/:slug" element={<PostPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<HomePage />} />
          {/* <Route path="/models" element={<ModelsPage />} />
          <Route path="/prompts" element={<PromptsPage />} />
          <Route path="/chat" element={<ChatPage />} />
          <Route path="/posts" element={<PostsPage />} /> */}
          {/* <Route path="/history" element={<HistoryPage />} />
          <Route path="/settings" element={<SettingsPage />} /> */}
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
