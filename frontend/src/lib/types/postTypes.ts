export interface Post {
  id: string;
  title: string;
  content: string;
  summary: string;
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
