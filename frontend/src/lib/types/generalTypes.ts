import { ReactNode } from "react";

export interface AppProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  initialData?: any;
}

export type Props = {
  children: React.ReactNode;
};

export interface ExtendedAppProps extends AppProps {
  children?: ReactNode;
}

export type Theme = "dark" | "light" | "system";

export interface ThemeProviderProps {
  children: ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
}

export interface ThemeProviderState {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

export interface SidebarContextType {
  open: boolean;
  toggle: () => void;
  close: () => void;
}
