import { BotResponse } from "@/lib/types/botTypes";
import axios from "axios";


export async function sendBotMessage(
  message: string,
  jwt: string
): Promise<BotResponse> {
  const { data } = await axios.post<BotResponse>(
    "/api/bot/chat",
    { message },
    {
      headers: { Authorization: `Bearer ${jwt}` },
    }
  );
  return data;
}
