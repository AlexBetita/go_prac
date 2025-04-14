import Navbar from "@/components/organisms/NavBar";
import { Outlet } from "react-router";
import { ModeToggle } from "@/components/templates/mode-toggle";

export default function AppLayout() {
	return (
		<div className="min-h-screen">
			<Navbar />
			<main className="px-4 py-6">
				<Outlet />
			</main>

            <ModeToggle/>
		</div>
	);
}
