import { useAppSelector } from "@/lib/hooks/AppHooks";
import { selectBotChatMessages } from "@/lib/store/selectors/botChatSelectors";
import { motion } from "framer-motion";

export default function ChatResponse() {
  const messages = useAppSelector(selectBotChatMessages);

  return (
    <div className="flex-1 overflow-y-auto px-4 py-6 space-y-3">
      {messages.map((msg, i) => (
        <motion.div
          key={i}
          initial={{ opacity: 0, y: 6 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: i * 0.03 }}
          className={`flex ${msg.isUser ? "justify-end" : "justify-start"}`}
        >
          <div
            className={`
              max-w-[75%] px-4 py-2 rounded-lg
              ${
                msg.isUser ? "bg-blue-500 text-white" : "bg-gray-200 text-black"
              }
            `}
          >
            {msg.text}
          </div>
        </motion.div>
      ))}
    </div>
  );
}
