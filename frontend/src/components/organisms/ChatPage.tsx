import ChatBox from "../molecules/chat/ChatBox";
import ChatResponse from "../molecules/chat/ChatResponse";
import { useAppSelector } from "@/lib/hooks/AppHooks";

export default function ChatPage() {
  const messages = useAppSelector((state) => state.bot.entries);

  return (
    <div className="flex flex-col justify-center h-full">
      {messages.length > 0 && <ChatResponse />}
      <ChatBox />
    </div>
  );
}
