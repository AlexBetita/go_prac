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