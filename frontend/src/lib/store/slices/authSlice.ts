import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { AuthState } from "@/lib/types/authTypes";
import {
  fetchProfileThunk,
  loginThunk,
  registerThunk,
} from "@/services/auth/thunks";

const initialState: AuthState = {
  token: null,
  tokenExp: null,
  user: null,
  status: "idle",
};

const slice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    logout(state) {
      state.token = null;
      state.tokenExp = null;
      state.user = null;
      state.status = "idle";
    },
    setToken(state, action: PayloadAction<{ token: string; exp: number }>) {
      state.token = action.payload.token;
      state.tokenExp = action.payload.exp;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginThunk.fulfilled, (st, { payload }) => {
        st.token = payload.token;
        st.tokenExp = payload.exp;
        st.user = payload.user;
        st.status = "succeeded";
      })
      .addCase(registerThunk.fulfilled, (st, { payload }) => {
        st.token = payload;
        st.status = "succeeded";
      })
      .addCase(fetchProfileThunk.fulfilled, (st, { payload }) => {
        st.user = payload;
        st.status = "succeeded";
      })
      .addCase(fetchProfileThunk.rejected, (st) => {
        st.token = null;
        st.tokenExp = null;
        st.user = null;
        st.status = "failed";
      });
  },
});

export const { logout, setToken } = slice.actions;
export default slice.reducer;
