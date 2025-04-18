import { createSlice } from "@reduxjs/toolkit";
import { chatWithBot } from "@/services/bot/thunks";
import { BotState } from "@/lib/types/botTypes";

const initialState: BotState = {
  entries: [],
  loading: false,
  error: undefined,
};

const botSlice = createSlice({
  name: "bot",
  initialState,
  reducers: {
    clearHistory(state) {
      state.entries = [];
      state.error = undefined;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(chatWithBot.pending, (state) => {
        state.loading = true;
        state.error = undefined;
      })
      .addCase(chatWithBot.fulfilled, (state, { payload }) => {
        state.loading = false;
        state.entries.push(payload);
      })
      .addCase(chatWithBot.rejected, (state, { payload }) => {
        state.loading = false;
        state.error = payload;
      });
  },
});

export const { clearHistory } = botSlice.actions;
export default botSlice.reducer;
