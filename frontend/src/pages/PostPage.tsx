import { Suspense, lazy } from "react";
import { useParams, useNavigate } from "react-router";
import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";
import { RootState } from "@/lib/store";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/atoms/shadCN/card";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { Badge } from "@/components/atoms/shadCN/badge";
import { Button } from "@/components/atoms/shadCN/button";

import { Skeleton } from "@/components/atoms/shadCN/skeleton";
import { PostSkeleton } from "@/components/molecules/posts/PostSkeleton";
import { fetchPostThunk } from "@/services/posts/thunks";

const ReactMarkdown = lazy(() => import("react-markdown"));

export default function PostPage() {
  const { slug = "" } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const { posts, loading, error } = useAppSelector((s: RootState) => s.posts);
  const post = posts.find((p) => p.slug === slug);

  useEffect(() => {
    if (!post && slug) dispatch(fetchPostThunk(slug));
  }, [post, slug, dispatch]);

  useEffect(() => {
    if (!post && error) {
      const t = setTimeout(() => navigate("/", { replace: true }), 3000);
      return () => clearTimeout(t);
    }
  }, [post, error, navigate]);

  if (loading && !post) return <PostSkeleton />;

  if (!post) {
    return (
      <div className="flex flex-col items-center justify-center h-full py-10 space-y-4">
        <h2 className="text-2xl font-semibold">
          {error ?? "404 – Post not found"}
        </h2>
        <Button onClick={() => navigate(-1)}>Go Back</Button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <Card className="shadow-xl">
        <CardHeader>
          <CardTitle className="text-3xl md:text-4xl font-bold">
            {post.title}
          </CardTitle>
          <CardDescription className="flex flex-wrap items-center space-x-2 text-sm text-muted-foreground">
            <span>{new Date(post.created_at).toLocaleDateString()}</span>
            <Separator className="h-4 w-px bg-muted" />
            <span>{post.views} views</span>
          </CardDescription>
          <div className="mt-4 flex flex-wrap gap-2">
            {post.tags.map((tag) => (
              <Badge key={tag} className="uppercase tracking-wide">
                {tag}
              </Badge>
            ))}
          </div>
        </CardHeader>

        <CardContent className="prose prose-slate dark:prose-invert max-w-none">
          <Suspense fallback={<SkeletonLines lines={8} />}>
            <ReactMarkdown>{post.content}</ReactMarkdown>
          </Suspense>
        </CardContent>
      </Card>

      <section className="mt-8">
        <h3 className="text-lg font-semibold mb-2">Keywords</h3>
        <div className="flex flex-wrap gap-2">
          {post.keywords.map((kw) => (
            <Badge variant="secondary" key={kw} className="cursor-default">
              {kw}
            </Badge>
          ))}
        </div>
      </section>
    </div>
  );
}

function SkeletonLines({ lines }: { lines: number }) {
  return (
    <div className="space-y-2">
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton key={i} className="h-4 w-full rounded" />
      ))}
    </div>
  );
}
