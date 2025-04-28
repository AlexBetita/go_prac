import { createContext, useContext } from "react";
import { ThemeProviderState } from "../types/generalTypes";

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => {},
};

export const ThemeProviderContext =
  createContext<ThemeProviderState>(initialState);

export function useTheme() {
  const context = useContext(ThemeProviderContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}
