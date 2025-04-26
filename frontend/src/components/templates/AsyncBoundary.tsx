import { Suspense } from "react";
import { ErrorBoundary } from "react-error-boundary";
import { Button } from "@/components/atoms/shadCN/button";

export function AsyncBoundary({ children }: { children: React.ReactNode }) {
  return (
    <ErrorBoundary
      fallbackRender={({ error, resetErrorBoundary }) => (
        <div className="flex flex-col items-center justify-center h-full py-10 space-y-4">
          <h2 className="text-2xl font-semibold text-destructive">
            {error.message ?? "Something went wrong"}
          </h2>
          <Button onClick={resetErrorBoundary}>Retry</Button>
        </div>
      )}
    >
      <Suspense
        fallback={
          <div className="flex items-center justify-center h-full py-10">
            <h2 className="text-xl font-semibold">Loading…</h2>
          </div>
        }
      >
        {children}
      </Suspense>
    </ErrorBoundary>
  );
}
