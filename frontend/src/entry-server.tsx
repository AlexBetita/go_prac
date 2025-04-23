import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import { StaticRouter } from "react-router";
import { Provider } from "react-redux";
import { store } from "@/lib/store";
import { ThemeProvider } from "@/components/templates/theme-provider";
import App from "@/pages/App";

export function render(url: string, initialData: any) {
  const dataScript = `<script>window.__INITIAL_DATA__ = ${JSON.stringify(
    initialData
  )}</script>`;

  const appHtml = renderToString(
    <StrictMode>
      <Provider store={store}>
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
          <StaticRouter location={url}>
            <App initialData={initialData} />
          </StaticRouter>
        </ThemeProvider>
      </Provider>
    </StrictMode>
  );

  return dataScript + appHtml;
}
