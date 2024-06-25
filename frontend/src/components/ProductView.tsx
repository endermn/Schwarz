export function ProductView({ children }: { children: React.ReactNode }) {
	return (
		<div className="h-screen items-center justify-items-center grid grid-cols-1 md:grid-cols-4 sm:grid-cols-4 gap-0">
			{children}
		</div>
	);
}
