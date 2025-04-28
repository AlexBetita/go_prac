import fs from "fs";
import path from "path";
import http from "http";
import { fileURLToPath } from "url";
import { createServer as createViteServer } from "vite";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

async function startServer() {
  const port = Number(process.env.PORT || 5174);
  const isProd = process.env.NODE_ENV === "production";

  let vite: any;
  if (!isProd) {
    vite = await createViteServer({
      root: process.cwd(),
      logLevel: "info",
      server: { middlewareMode: true, hmr: false },
      appType: "custom",
    });
  }

  const server = http.createServer(async (req, res) => {
    try {
      const url = req.url || "/";

      if (!isProd && vite) {
        await new Promise<void>((resolve, reject) =>
          vite.middlewares(req, res, (err: any) =>
            err ? reject(err) : resolve()
          )
        );
        if (res.writableEnded) return;
      }

      if (url.startsWith("/posts/")) {
        const indexPath = path.resolve(
          __dirname,
          isProd ? "dist/client/index.html" : "index.html"
        );
        const indexHtml = fs.readFileSync(indexPath, "utf-8");
        const template = isProd
          ? indexHtml
          : await vite.transformIndexHtml(url, indexHtml);

        const { render } = isProd
          ? await import("./dist/server/entry-server.js")
          : await vite.ssrLoadModule("/src/entry-server.tsx");

        const apiRes = await fetch(`http://localhost:8080/api${url}`);
        const initialData = apiRes.ok ? await apiRes.json() : null;

        const appHtml = await render(url, initialData);
        const html = template.replace("<!--app-html-->", appHtml);

        res.writeHead(200, { "Content-Type": "text/html" });
        res.end(html);
        return;
      }

      if (isProd) {
        const filePath = path.join(__dirname, "dist/client", url);
        if (fs.existsSync(filePath) && fs.statSync(filePath).isFile()) {
          const stream = fs.createReadStream(filePath);
          stream.pipe(res);
          return;
        }
      }

      res.writeHead(404, { "Content-Type": "text/plain" });
      res.end("Not found");
    } catch (e: any) {
      if (vite) vite.ssrFixStacktrace(e);
      console.error(e);
      if (!res.writableEnded) {
        res.writeHead(500, { "Content-Type": "text/plain" });
        res.end(e.stack || e.message);
      }
    }
  });

  server.listen(port, () =>
    console.log(`SSR server running at http://localhost:${port}`)
  );
}

startServer();
