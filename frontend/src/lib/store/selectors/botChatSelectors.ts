import { RootState } from "..";

export const selectBotChatMessages
    = (state: RootState) =>
  state.bot.entries.flatMap((e) => [
    { text: e.userMessage, isUser: true },
    { text: String(e.response), isUser: false },
  ]);
