import ChatBox from "../molecules/ChatBox";
import ChatResponse from "../molecules/ChatResponse";
import { useAppSelector } from "@/lib/hooks/AppHooks";

export default function ChatPage() {
  const messages = useAppSelector((state) => state.bot.entries);

  return (
    <div className="flex flex-col h-screen">
      {messages.length > 0 && <ChatResponse />}
      <ChatBox />
    </div>
  );
}
