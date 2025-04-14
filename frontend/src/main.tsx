import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";

import AppRoutes from "@/routes/BaseRoutes";
import { store } from "@/lib/store";
import { ThemeProvider } from "@/components/templates/theme-provider"

import "@/styles/globals.css";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<Provider store={store}>
            <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			    <AppRoutes />
            </ThemeProvider>
		</Provider>
	</StrictMode>
);
