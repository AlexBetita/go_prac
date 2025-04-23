import { Props } from "../types/generalTypes";
import { InitialDataContext } from "./initialDataContext";

export function InitialDataContextProvider({
  children,
  initialData,
}: Props & { initialData: any }) {
  return (
    <InitialDataContext.Provider value={initialData}>
      {children}
    </InitialDataContext.Provider>
  );
}
