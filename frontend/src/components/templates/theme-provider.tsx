import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

type Theme = "dark" | "light" | "system";

interface ThemeProviderProps {
  children: ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
}

interface ThemeProviderState {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => {},
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

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
	  console.log(e)
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

export function useTheme() {
  const context = useContext(ThemeProviderContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}
