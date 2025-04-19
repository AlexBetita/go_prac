import { Post, PostsState } from "@/lib/types/postTypes";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

const initialState: PostsState = {
  posts: [],
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
});

export const { addPosts, addPost, clearPosts } = postSlice.actions;
export default postSlice.reducer;
