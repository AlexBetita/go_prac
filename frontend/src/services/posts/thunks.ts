import { createAsyncThunk } from "@reduxjs/toolkit";
import { AxiosError } from "axios";

import { fetchPosts, fetchPost } from "./api";
import { Post } from "@/lib/types/postTypes";


export const fetchPostsThunk = createAsyncThunk<
  Post[],
  void,
  { rejectValue: string }
>("posts/fetchAll", async (_, { rejectWithValue }) => {
  try {
    return await fetchPosts();
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(
      err.response?.data?.message ?? "Failed to load posts"
    );
  }
});

export const fetchPostThunk = createAsyncThunk<
  Post,
  string,
  { rejectValue: string }
>("posts/fetchOne", async (identifier, { rejectWithValue }) => {
  try {
    return await fetchPost(identifier);
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(err.response?.data?.message ?? "Post not found");
  }
});
