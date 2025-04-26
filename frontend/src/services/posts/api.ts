import axios from "axios";
import { Post } from "@/lib/types/postTypes";

export async function fetchPosts(): Promise<Post[]> {
  const { data } = await axios.get<Post[]>("/api/posts");
  return data;
}

export async function fetchPost(identifier: string): Promise<Post> {
  const { data } = await axios.get<Post>(`/api/posts/${identifier}`);
  return data;
}