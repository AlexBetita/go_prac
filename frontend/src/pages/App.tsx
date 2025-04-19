import Navbar from "@/components/organisms/NavBar";
import { Outlet } from "react-router";
import { ModeToggle } from "@/components/templates/mode-toggle";
import { AuthProvider } from "@/lib/providers/AuthProvider";

export default function AppLayout() {
	return (
		<AuthProvider>
			<div className="min-h-screen">
				<Navbar />
				<main className="py-1">
					<Outlet />
				</main>

				<ModeToggle/>
			</div>
		</AuthProvider>
	);
}
