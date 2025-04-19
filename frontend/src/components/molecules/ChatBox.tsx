import { useRef, useState } from "react";
import { motion } from "framer-motion";

import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";
import { Button } from "@/components/atoms/shadCN/button";
import { Textarea } from "@/components/atoms/shadCN/textarea";
import { Badge } from "@/components/atoms/shadCN/badge";
import { chatWithBot } from "@/services/bot/thunks";

const suggestions = ["Docs", "Examples", "Troubleshoot", "Optimize"];

export default function ChatBox() {
  const dispatch = useAppDispatch();

  const user = useAppSelector((state) => state.auth.user);
  const loading = useAppSelector((state) => state.bot.loading);
  const interactions = useAppSelector((state)=> state.bot.entries);

  const [input, setInput] = useState("");
  const [hide, setHide] = useState(false)

  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleInput = (e: React.FormEvent<HTMLTextAreaElement>) => {
    const ta = e.currentTarget;
    ta.style.height = "auto";
    ta.style.height = `${ta.scrollHeight}px`;
    setInput(ta.value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const text = input.trim();
    if (!text) return;

    dispatch(chatWithBot(text));
    setHide(true)
    setInput("");
    if (textareaRef.current) textareaRef.current.style.height = "auto";
  };

  return (
    <div
      className="z-1
    flex items-center justify-center px-4 md:overflow-hidden h-full"
    >
      <div className="w-full max-w-5xl flex flex-col items-center">
        {!hide ||
          (interactions.length > 0 && (
            <motion.h1
              initial={{ opacity: 0, y: -8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4 }}
              className="mb-6 text-2xl font-semibold text-center"
            >
              Good to see you, {user?.email.split("@")[0]}
            </motion.h1>
          ))}

        <motion.div
          initial={{ opacity: 0, scale: 0.97 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ type: "spring", stiffness: 250, damping: 20 }}
          className="w-full h-full"
        >
          <form onSubmit={handleSubmit} className="flex flex-col">
            <Textarea
              ref={textareaRef}
              value={input}
              onInput={handleInput}
              placeholder="Ask anything…"
              className="
              flex-1
              min-h-[44px]
              max-h-[240px]
              resize-none
              rounded-lg
              border border-gray-400
              px-3 py-2
              focus:outline-none focus:ring-2 focus:ring-gray-400
              overflow-y-auto
              transition-[height] duration-200 ease-in-out
            "
            />

            <div className="flex items-center justify-between mt-2">
              <div className="flex flex-wrap gap-2">
                {suggestions.map((text) => (
                  <motion.div
                    key={text}
                    initial={{ opacity: 0, y: 6 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.2 }}
                  >
                    <Badge
                      variant="outline"
                      onClick={() => {
                        setInput(text);
                        setTimeout(() => {
                          if (textareaRef.current) {
                            textareaRef.current.style.height = "auto";
                            textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
                          }
                        }, 0);
                      }}
                      className="cursor-pointer"
                    >
                      {text}
                    </Badge>
                  </motion.div>
                ))}
              </div>

              <motion.div
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.97 }}
              >
                <Button
                  type="submit"
                  disabled={!input.trim()}
                  className="h-10 px-6"
                >
                  {loading ? "Sending…" : "Send"}
                </Button>
              </motion.div>
            </div>
          </form>
        </motion.div>
      </div>
    </div>
  );
}
