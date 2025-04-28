import { Post, PostsState } from "@/lib/types/postTypes";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { fetchPostsThunk, fetchPostThunk } from "@/services/posts/thunks";

import { mockPosts } from "../mock/postMockData";

const initialState: PostsState = {
  posts: import.meta.env.VITE_ENABLE_MOCK === "yeah" ? mockPosts : [],
  loading: false,
  error: undefined,
};

const postSlice = createSlice({
  name: "posts",
  initialState,
  reducers: {
    addPosts(state, action: PayloadAction<Post[]>) {
      state.posts = action.payload;
    },
    addPost(state, action: PayloadAction<Post>) {
      state.posts.push(action.payload);
    },
    clearPosts(state) {
      state.posts = [];
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchPostsThunk.pending, (state) => {
        state.loading = true;
        state.error = undefined;
      })
      .addCase(fetchPostsThunk.fulfilled, (state, action) => {
        state.loading = false;
        state.posts = action.payload;
      })
      .addCase(fetchPostsThunk.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      });

    builder
      .addCase(fetchPostThunk.pending, (state) => {
        state.loading = true;
        state.error = undefined;
      })
      .addCase(fetchPostThunk.fulfilled, (state, action) => {
        state.loading = false;
        const idx = state.posts.findIndex(
          (p) => p.slug === action.payload.slug
        );
        if (idx >= 0) {
          state.posts[idx] = action.payload;
        } else {
          state.posts.push(action.payload);
        }
      })
      .addCase(fetchPostThunk.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const { addPosts, addPost, clearPosts } = postSlice.actions;
export default postSlice.reducer;
