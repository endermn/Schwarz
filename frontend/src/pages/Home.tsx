export function Home() {
	return (
		<div>
			<div className="py-32 sm:py-48 md:py-56 max-w-2xl mx-auto">
				<div className="text-center">
					<h1 className="font-extrabold sm:text-6xl">Пазарувай лесно</h1>
					<p className="leading-8 text-xl mt-6 text-current/50">
						Лорем ипсум, да много як магазин... Индрустриалната революция и
						нейните последици.
					</p>
					<div className="mt-10 gap-x-6 flex items-center justify-center">
						<a
							href="/signup"
							className="font-semibold leading-5 py-3 px-4 bg-indigo-600 text-white rounded-md"
						>
							Започни
						</a>
						<a href="#" className="leading-6 font-semibold text-sm">
							Научи повече <span aria-hidden="true">→</span>
						</a>
					</div>
				</div>
			</div>
		</div>
	);
}
