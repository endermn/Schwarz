import { getContext } from "@/App";

export function Home() {
	const context = getContext();
	return (
		<div>
			<div className="mx-auto max-w-2xl py-32 sm:py-48 md:py-56">
				<div className="text-center">
					<h1 className="text-5xl font-extrabold">Пазарувай лесно</h1>
					<p className="text-current/50 mt-6 text-xl leading-8 text-muted-foreground">
						не губи време в магазина
					</p>
					<div className="mt-10 flex items-center justify-center gap-x-6">
						<a
							href="/signup"
							className="rounded-md bg-indigo-600 px-4 py-3 font-semibold leading-5 text-white"
						>
							Пазарувай
						</a>
						<a
							href="https://www.wikipedia.org/"
							target="_blank"
							rel="noopener noreferrer"
							className="text-sm font-semibold leading-6"
						>
							Научи повече
						</a>
					</div>
				</div>
			</div>
		</div>
	);
}
