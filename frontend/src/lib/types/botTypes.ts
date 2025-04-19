export interface BotState {
  entries: BotEntry[];
  loading: boolean;
  error?: string;
}

export interface BotResponse {
  type: string;
  response: any;
}

export interface BotEntry {
  userMessage: string;
  type: BotResponse["type"];
  response: BotResponse["response"];
}


export const BotInteractionResponseTypes = [
  "interaction",
  "related_posts",
] as const;

export type BotResponseType = (typeof BotInteractionResponseTypes)[number];