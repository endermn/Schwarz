export function ProductView({ children }: { children: React.ReactNode }) {
	return (
		<div className="m-10 grid justify-items-center gap-4 text-balance sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
			{children}
		</div>
	);
}
