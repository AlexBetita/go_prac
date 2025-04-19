import { useEffect, useRef, useState } from "react";
import { useAppSelector } from "@/lib/hooks/AppHooks";
import { selectBotChatMessages } from "@/lib/store/selectors/botChatSelectors";
import { Label } from "@radix-ui/react-label";
import { motion } from "framer-motion";
import { Copy, ThumbsUp, ThumbsDown, ChevronDown } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/atoms/shadCN/tooltip";

export default function ChatResponse() {
  const messages = useAppSelector(selectBotChatMessages);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [showArrow, setShowArrow] = useState(false);

  const lastBotMsgIndex = [...messages].reverse().findIndex((m) => !m.isUser);
  const lastBotIndex =
    lastBotMsgIndex !== -1 ? messages.length - 1 - lastBotMsgIndex : -1;

  useEffect(() => {
    const el = scrollRef.current;
    if (!el) return;

    const checkScrollable = () => {
      const threshold = 4;
      const canScrollMore =
        el.scrollHeight - el.scrollTop > el.clientHeight + threshold;
      setShowArrow(canScrollMore);
    };

    checkScrollable();
    el.addEventListener("scroll", checkScrollable);

    const resizeObserver = new ResizeObserver(checkScrollable);
    resizeObserver.observe(el);

    return () => {
      el.removeEventListener("scroll", checkScrollable);
      resizeObserver.disconnect();
    };
  }, [messages]);

  const handleScrollBottom = () => {
    scrollRef.current?.scrollTo({
      top: scrollRef.current.scrollHeight,
      behavior: "smooth",
    });
  };

  return (
    <TooltipProvider skipDelayDuration={0}>
      <div className="relative w-full mx-auto h-full">
        <div
          ref={scrollRef}
          className="relative p-4 max-h-[80vh] overflow-y-auto -mb-6 pr-6"
          style={{
            scrollbarWidth: "thin",
            scrollbarColor: "#cbd5e1 transparent",
          }}
        >
          <div className="space-y-3">
            {messages.map((msg, i) => (
              <div key={i} className="space-y-1 group">
                <motion.div
                  initial={{ opacity: 0, y: 6 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.2, delay: i * 0.03 }}
                  className={`flex ${
                    msg.isUser ? "justify-end" : "justify-start"
                  }`}
                >
                  <Label
                    className={`px-4 py-2 rounded-lg text-sm ${
                      msg.isUser
                        ? "bg-gray-100 dark:bg-gray-800 text-black dark:text-white shadow-sm max-w-[85%] ml-12"
                        : "text-foreground max-w-[75%] mr-12"
                    }`}
                  >
                    {msg.text}
                  </Label>
                </motion.div>

                {!msg.isUser && (
                  <div
                    className={`
                      flex text-muted-foreground pl-2 ml-2 pb-4 transition-opacity duration-300
                      ${
                        i !== lastBotIndex
                          ? "opacity-0 group-hover:opacity-100 hover:opacity-100 delay-[0ms] group-hover:delay-[0ms] hover:delay-[0ms] group-hover:transition-opacity"
                          : ""
                      }
                    `}
                    style={{
                      transitionDelay: i !== lastBotIndex ? "0s" : undefined,
                    }}
                    onMouseLeave={(e) => {
                      const target = e.currentTarget;
                      target.style.transitionDelay = "500ms";
                    }}
                    onMouseEnter={(e) => {
                      const target = e.currentTarget;
                      target.style.transitionDelay = "0ms";
                    }}
                  >
                    <Tooltip disableHoverableContent delayDuration={300}>
                      <TooltipTrigger asChild>
                        <div className="w-8 h-6 flex items-center justify-center cursor-pointer">
                          <Copy className="w-4 h-4" />
                        </div>
                      </TooltipTrigger>
                      <TooltipContent side="bottom" sideOffset={1}>
                        Copy
                      </TooltipContent>
                    </Tooltip>

                    <Tooltip disableHoverableContent delayDuration={300}>
                      <TooltipTrigger asChild>
                        <div className="w-8 h-6 flex items-center justify-center cursor-pointer">
                          <ThumbsUp className="w-4 h-4" />
                        </div>
                      </TooltipTrigger>
                      <TooltipContent side="bottom" sideOffset={1}>
                        Like
                      </TooltipContent>
                    </Tooltip>

                    <Tooltip disableHoverableContent delayDuration={300}>
                      <TooltipTrigger asChild>
                        <div className="w-8 h-6 flex items-center justify-center cursor-pointer">
                          <ThumbsDown className="w-4 h-4" />
                        </div>
                      </TooltipTrigger>
                      <TooltipContent side="bottom" sideOffset={1}>
                        Dislike
                      </TooltipContent>
                    </Tooltip>
                  </div>
                )}
              </div>
            ))}

            
              <div className={`sticky bottom-4 pt-2 flex justify-center z-10 ${showArrow ? '': 'opacity-0'}`}>
                <ChevronDown
                  onClick={handleScrollBottom}
                  className="w-8 h-8 text-muted-foreground cursor-pointer hover:text-primary transition"
                />
              </div>
          </div>
        </div>
      </div>
    </TooltipProvider>
  );
}
