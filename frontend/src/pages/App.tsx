import { Outlet } from "react-router";
import { ModeToggle } from "@/components/templates/ModeToggle";
import { ExtendedAppProps } from "@/lib/types/generalTypes";
import { InitialDataContextProvider } from "@/lib/providers/AppProvider";
import NavBar from "@/components/organisms/NavBar";
import AppSidebar from "@/components/organisms/AppSidebar";
import { useSidebar } from "@/lib/context/SideBarContext";

export default function App({ initialData, children }: ExtendedAppProps) {
  const { open } = useSidebar();
  return (
    <InitialDataContextProvider initialData={initialData}>
      <NavBar />
      <div className="h-full">
        <AppSidebar />
        <main
          className={`flex-1 overflow-auto transition-all duration-300 ${
            open ? "ml-64" : "ml-0"
          }`}
        >
          {children || <Outlet />}
        </main>
        <ModeToggle />
      </div>
    </InitialDataContextProvider>
  );
}
