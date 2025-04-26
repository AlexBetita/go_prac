import { StrictMode } from "react";
import { hydrateRoot, createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";

import { store, persistor } from "@/lib/store";
import { ThemeProvider } from "@/components/templates/ThemeProvider";

import "@/styles/globals.css";

import BaseRoutes from "./routes/BaseRoutes";
import { AuthProvider } from "./lib/providers/AuthProvider";
import { SidebarProvider } from "./lib/providers/SidebarProvider";

const rootElement = document.getElementById("root")!;
const AppTree = (
  <StrictMode>
    <Provider store={store}>
      <AuthProvider>
        <SidebarProvider>
          <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <PersistGate loading={null} persistor={persistor}>
              <BrowserRouter>
                <BaseRoutes />
              </BrowserRouter>
            </PersistGate>
          </ThemeProvider>
        </SidebarProvider>
      </AuthProvider>
    </Provider>
  </StrictMode>
);

if (rootElement.hasChildNodes()) {
  hydrateRoot(rootElement, AppTree);
} else {
  createRoot(rootElement).render(AppTree);
}
