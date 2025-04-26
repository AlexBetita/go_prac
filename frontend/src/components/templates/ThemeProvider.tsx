import { ThemeProviderContext } from "@/lib/context/ThemeContext";
import { ThemeProviderProps, Theme } from "@/lib/types/generalTypes";
import { useEffect, useState } from "react";

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "vite-ui-theme",
}: ThemeProviderProps) {
  const isBrowser = typeof window !== "undefined";

  const [theme, setThemeState] = useState<Theme>(() => {
    if (isBrowser) {
      try {
        const stored = window.localStorage.getItem(storageKey) as Theme | null;
        return stored || defaultTheme;
      } catch {
        return defaultTheme;
      }
    }
    return defaultTheme;
  });

  useEffect(() => {
    if (!isBrowser) return;

    try {
      window.localStorage.setItem(storageKey, theme);
    } catch (e) {
      console.log(e);
    }

    const root = window.document.documentElement;
    root.classList.remove("light", "dark");
    if (theme === "system") {
      const systemDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      root.classList.add(systemDark ? "dark" : "light");
    } else {
      root.classList.add(theme);
    }
  }, [theme, storageKey, isBrowser]);

  const setTheme = (newTheme: Theme) => {
    setThemeState(newTheme);
  };

  return (
    <ThemeProviderContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeProviderContext.Provider>
  );
}
