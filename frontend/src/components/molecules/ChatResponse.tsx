import { useAppSelector } from "@/lib/hooks/AppHooks";
import { selectBotChatMessages } from "@/lib/store/selectors/botChatSelectors";
import { motion } from "framer-motion";

export default function ChatResponse() {
  const messages = useAppSelector(selectBotChatMessages);

  return (
    <div className="relative w-full max-w-3xl mx-auto">
      <div
        className="bg-card border border-gray-200 rounded-xl shadow-md p-4 max-h-[65vh] overflow-y-auto -mb-6 pr-6"
        style={{
          scrollbarWidth: "thin",
          scrollbarColor: "#cbd5e1 transparent",
        }}
      >
        <div className="space-y-3">
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
                  px-4 py-2 rounded-lg text-sm
                  ${msg.isUser
                    ? "bg-blue-500 text-white shadow-sm max-w-[85%] ml-12"
                    : "bg-gray-100 text-black shadow-sm max-w-[75%] mr-12"}
                `}
              >
                {msg.text}
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}
