
import { Outlet } from "react-router";
import { ModeToggle } from "@/components/templates/ModeToggle";
import { AuthProvider } from "@/lib/providers/AuthProvider";
import { ExtendedAppProps } from "@/lib/types/generalTypes";
import { InitialDataContextProvider } from "@/lib/providers/AppProvider";

export default function App({ initialData, children }: ExtendedAppProps) {
  return (
    <InitialDataContextProvider initialData={initialData}>
      <AuthProvider>
        <div className="h-full">
          <main>
            {children || <Outlet />}
          </main>
          <ModeToggle />
        </div>
      </AuthProvider>
    </InitialDataContextProvider>
  );
}
