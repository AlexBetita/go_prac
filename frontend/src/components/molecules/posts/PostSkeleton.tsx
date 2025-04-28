import { Skeleton } from "@/components/atoms/shadCN/skeleton";

export function PostSkeleton() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-8 space-y-6">
      <div className="rounded-2xl border shadow-xl p-6 space-y-4">
        <Skeleton className="h-10 w-3/4 rounded" />
        <div className="flex space-x-4">
          <Skeleton className="h-4 w-24 rounded" />
          <Skeleton className="h-4 w-16 rounded" />
        </div>
        <div className="flex gap-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-6 w-20 rounded-full" />
          ))}
        </div>
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-4 w-full rounded" />
        ))}
      </div>
    </div>
  );
}
