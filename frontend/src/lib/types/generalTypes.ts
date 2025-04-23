import { ReactNode } from "react";

export interface AppProps {
  initialData?: any;
}

export type Props = {
  children: React.ReactNode;
};

export interface ExtendedAppProps extends AppProps {
  children?: ReactNode;
}