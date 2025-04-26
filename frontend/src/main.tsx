import { StrictMode } from "react";
import { hydrateRoot, createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";

import App from "@/pages/App";
import AppRoutes from "@/routes/BaseRoutes";
import { store, persistor } from "@/lib/store";
import { ThemeProvider } from "@/components/templates/ThemeProvider";

import "@/styles/globals.css";
import Navbar from "./components/organisms/NavBar";

const initialData = (window as any).__INITIAL_DATA__ || null;

const rootElement = document.getElementById("root")!;
const AppTree = (
  <StrictMode>
    <Provider store={store}>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <PersistGate loading={null} persistor={persistor}>
          <BrowserRouter>
            <Navbar />
            <App initialData={initialData}>
              <AppRoutes />
            </App>
          </BrowserRouter>
        </PersistGate>
      </ThemeProvider>
    </Provider>
  </StrictMode>
);

if (rootElement.hasChildNodes()) {
  hydrateRoot(rootElement, AppTree);
} else {
  createRoot(rootElement).render(AppTree);
}
