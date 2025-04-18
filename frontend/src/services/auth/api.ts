import axios from "axios";
import { User } from "@/lib/types/authTypes";

export async function login(
  email: string,
  password: string
): Promise<{ token: string; exp: number; user: User }> {
  const { data } = await axios.post("/api/auth/login", { email, password });
  return {
    token: data.token,
    exp: data.exp,
    user: data.user,
  };
}

export async function register(
  email: string,
  password: string
): Promise<string> {
  const { data } = await axios.post("/api/auth/register", { email, password });
  return data.token;
}

export async function fetchProfile(jwt: string): Promise<User> {
  const { data } = await axios.get("/api/auth/profile", {
    headers: { Authorization: `Bearer ${jwt}` },
  });
  return data;
}
