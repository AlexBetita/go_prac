export interface User {
  email: string;
  provider: "local" | "google";
  created_at: number;
  updated_at: number;
}

export interface AuthState {
  token: string | null;
  tokenExp: number | null;
  user: User | null;
  status: "idle" | "loading" | "failed" | "succeeded";
}
