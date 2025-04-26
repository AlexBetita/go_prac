import { Props } from "../types/generalTypes";
import { InitialDataContext } from "../context/InitialDataContext";

export function InitialDataContextProvider({
  children,
  initialData,
// eslint-disable-next-line @typescript-eslint/no-explicit-any
}: Props & { initialData: any }) {
  return (
    <InitialDataContext.Provider value={initialData}>
      {children}
    </InitialDataContext.Provider>
  );
}
