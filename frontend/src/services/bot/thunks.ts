import { createAsyncThunk } from "@reduxjs/toolkit";
import { AxiosError } from "axios";

import { sendBotMessage } from "./api";
import { AuthState } from "@/lib/types/authTypes";
import { BotEntry, BotInteractionResponseTypes, BotResponseType } from "@/lib/types/botTypes";

export const chatWithBot = createAsyncThunk<
  BotEntry,
  string,
  { state: { auth: AuthState }; rejectValue: string }
>("bot/chat", async (message, { getState, rejectWithValue }) => {
  const jwt = getState().auth.token;
  if (!jwt) return rejectWithValue("No token in store");
  try {
    const botRes = await sendBotMessage(message, jwt);
    if (BotInteractionResponseTypes.includes(botRes.type as BotResponseType)) {
      return {
        userMessage: message,
        type: botRes.type as BotEntry["type"],
        response: botRes.response
      };
    }
    // handle post differently... WIP
      return {
        userMessage: message,
        type: botRes.type as BotEntry["type"],
        response: botRes.response,
      };
  } catch (e: unknown) {
    const err = e as AxiosError<{ message?: string }>;
    return rejectWithValue(err.response?.data?.message || err.message);
  }
});
