import { useEffect, useRef, useState } from "react";
import { motion } from "framer-motion";
import { useAppDispatch } from "@/lib/hooks/AppHooks";
import { fetchProfileThunk } from "@/lib/store/slices/authSlice";
import { Button } from "@/components/atoms/shadCN/button";
import { Textarea } from "@/components/atoms/shadCN/textarea";
import { Badge } from "@/components/atoms/shadCN/badge";

const suggestions = ["Docs", "Examples", "Troubleshoot", "Optimize"];

export default function HomePage() {
  const dispatch = useAppDispatch();
  const [input, setInput] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    dispatch(fetchProfileThunk());
  }, [dispatch]);

  const handleInput = (e: React.FormEvent<HTMLTextAreaElement>) => {
    const ta = e.currentTarget;
    ta.style.height = "auto";
    ta.style.height = `${ta.scrollHeight}px`;
    setInput(ta.value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;
    console.log("send:", input);
    setInput("");
    if (textareaRef.current) textareaRef.current.style.height = "auto";
  };

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center px-4">
      <motion.h1
        initial={{ opacity: 0, y: -8 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="mb-6 text-2xl font-semibold text-center"
      >
        Good to see you, Alex.
      </motion.h1>

      <motion.div
        initial={{ opacity: 0, scale: 0.97 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ type: "spring", stiffness: 250, damping: 20 }}
        className="w-full max-w-3xl bg-card rounded-2xl shadow-md px-6 py-5"
      >
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
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
              border border-gray-300
              px-3 py-2
              focus:outline-none focus:ring-2 focus:ring-gray-400
              overflow-y-auto
              transition-[height] duration-200 ease-in-out
            "
          />

          <div className="flex items-center justify-between mt-4">
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

            <motion.div whileHover={{ scale: 1.03 }} whileTap={{ scale: 0.97 }}>
              <Button
                type="submit"
                disabled={!input.trim()}
                className="h-11 px-6"
              >
                Send
              </Button>
            </motion.div>
          </div>
        </form>
      </motion.div>
    </div>
  );
}
