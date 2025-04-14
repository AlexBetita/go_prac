import { Link } from "react-router";
import { useAppSelector } from "@/lib/hooks/AppHooks";

import { Avatar, AvatarFallback } from "@/components/atoms/shadCNavatar";
import { Button } from "@/components/atoms/shadCNbutton";

export default function Navbar() {
	const { user } = useAppSelector((s) => s.auth);

	const initials = user?.email?.charAt(0).toUpperCase() ?? "?";

	return (
		<header className="w-full border-b border-border bg-background shadow-sm">
			<div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
                
				<Link to="/" className="text-lg font-semibold">
					MyApp
				</Link>

				<div className="flex items-center gap-4">
					<nav className="hidden md:flex gap-4 text-sm">
						<Link to="/" className="hover:underline">
							Home
						</Link>
						{user && (
							<Link to="/profile" className="hover:underline">
								Profile
							</Link>
						)}
					</nav>

					{user ? (
						<Link to="/profile">
							<Avatar className="h-8 w-8">
								<AvatarFallback>{initials}</AvatarFallback>
							</Avatar>
						</Link>
					) : (
						<Button asChild variant="outline" size="sm">
							<Link to="/login">Login</Link>
						</Button>
					)}
				</div>
			</div>
		</header>
	);
}
