export interface Post {
  id: string;
  topic: string;
  content: string;
  summary: string;
  message: string;
  keywords: string[];
  tags: string[];
  slug: string;
  views: number;
  created_by: string;
  created_at: string;
}

export interface PostsState {
  posts: Post[];
  loading: boolean;
  error?: string;
}
