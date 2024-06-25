export function ProductView({ children }: { children: React.ReactNode }) {
	return (
		<div className="m-4 gap-8 justify-items-center grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4">
			{children}
		</div>
	);
}
