import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";
import { fetchProfileThunk, logout } from "@/lib/store/slices/authSlice";

import { Button } from "@/components/atoms/shadCNbutton";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/atoms/shadCNcard";
import { Skeleton } from "@/components/atoms/shadCNskeleton";

export default function ProfilePage() {
	const { user } = useAppSelector((s) => s.auth);
	const dispatch = useAppDispatch();

	useEffect(() => {
		dispatch(fetchProfileThunk());
	}, [dispatch]);

	if (!user) {
		return (
			<div className="p-10 max-w-lg mx-auto">
				<Skeleton className="h-10 w-32 mb-4" />
				<Skeleton className="h-40 w-full rounded" />
			</div>
		);
	}

	return (
		<div className="relative min-h-screen bg-background">

			{/* Main Profile Card */}
			<div className="flex justify-center items-center pt-28">
				<Card className="w-full max-w-md shadow-md">
					<CardHeader>
						<CardTitle className="text-center text-2xl font-semibold">
							User Profile
						</CardTitle>
					</CardHeader>
					<CardContent className="space-y-4 text-sm text-muted-foreground">
						<div className="grid grid-cols-2 gap-y-3">
							<span className="font-medium text-foreground">
								ID:
							</span>
							<span>{user.id}</span>
							<span className="font-medium text-foreground">
								Email:
							</span>
							<span>{user.email}</span>
							<span className="font-medium text-foreground">
								Provider:
							</span>
							<span className="capitalize">{user.provider}</span>
							<span className="font-medium text-foreground">
								Created:
							</span>
							<span>
								{new Date(
									user.created_at * 1000
								).toLocaleString()}
							</span>
							<span className="font-medium text-foreground">
								Updated:
							</span>
							<span>
								{new Date(
									user.updated_at * 1000
								).toLocaleString()}
							</span>
						</div>
						<Button
							className="w-full mt-6"
							variant="destructive"
							onClick={() => dispatch(logout())}
						>
							Logout
						</Button>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
