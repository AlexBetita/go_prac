import { createAsyncThunk } from "@reduxjs/toolkit";
import { sendBotMessage } from "./api";
import { AuthState } from "@/lib/types/authTypes";
import { BotEntry } from "@/lib/types/botTypes";

export const chatWithBot = createAsyncThunk<
  BotEntry,
  string,
  { state: { auth: AuthState }; rejectValue: string }
>("bot/chat", async (message, { getState, rejectWithValue }) => {
  const jwt = getState().auth.token;
  if (!jwt) return rejectWithValue("No token in store");
  try {
    const botRes = await sendBotMessage(message, jwt);
    return {
      userMessage: message,
      type: botRes.type as BotEntry["type"],
      response: botRes.response,
    };
  } catch (err: any) {
    return rejectWithValue(err.response?.data || err.message);
  }
});
