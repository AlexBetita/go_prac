import { createAsyncThunk } from "@reduxjs/toolkit";
import { AxiosError } from "axios";

import { login, register, fetchProfile } from "./api";
import { User, AuthState } from "@/lib/types/authTypes";

export const loginThunk = createAsyncThunk<
  { token: string; exp: number; user: User },
  { email: string; password: string },
  { rejectValue: string }
  >("auth/login", async ({ email, password }, { rejectWithValue }) => {
    try {
      const response = await login(email, password);
    return response;
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(err.response?.data?.message ?? "Login failed");
  }
});

export const registerThunk = createAsyncThunk<
  string,
  { email: string; password: string },
  { rejectValue: string }
>("auth/register", async ({ email, password }, { rejectWithValue }) => {
  try {
    return await register(email, password);
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(err.response?.data?.message ?? "Register failed");
  }
});

export const fetchProfileThunk = createAsyncThunk<
  User,
  void,
  { state: { auth: AuthState }; rejectValue: string }
>("auth/profile", async (_, { getState, rejectWithValue }) => {
  const jwt = getState().auth.token;
  if (!jwt) return rejectWithValue("No token");
  try {
    return await fetchProfile(jwt);
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(err.response?.data?.message ?? "Unauthorized");
  }
});
