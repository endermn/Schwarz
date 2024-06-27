export function ProductView({ children }: { children: React.ReactNode }) {
	return (
		<div className="m-4 gap-4 justify-items-center grid m:grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
			{children}
		</div>
	);
}
