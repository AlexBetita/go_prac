import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { Post } from "@/lib/types/postTypes";

export default function PostPage() {
  const { identifier } = useParams<{ identifier: string }>();
  const [post, setPost] = useState<Post | null>(null);

  useEffect(() => {
    fetch(`/api/posts/${identifier}`)
      .then((res) => res.json())
      .then(setPost);
  }, [identifier]);

  if (!post) return <div>Loading…</div>;
  return (
    <article>
      <h1>{post.title}</h1>
      <div>{post.content}</div>
    </article>
  );
}
