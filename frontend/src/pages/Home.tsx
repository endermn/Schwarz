export function Home() {
	return (
		<div>
			<div className="py-32 sm:py-48 md:py-56 max-w-2xl mx-auto">
				<div className="hidden sm:mb-8 sm:flex sm:justify-center">
					<div className="ring-1 ring-offset-1 shadow-sm leading-6 py-1 px-3 rounded-full">
						Announcing our next round of funding.{" "}
						<a href="#" className="font-semibold">
							Read more
						</a>
					</div>
				</div>
				<div className="text-center">
					<h1 className="font-extrabold sm:text-6xl">
						Data to enrich your online business
					</h1>
					<p className="leading-8 text-xl mt-6 text-current/50">
						Anim aute id magna aliqua ad ad non deserunt sunt. Qui irure qui
						lorem cupidatat commodo. Elit sunt amet fugiat veniam occaecat
						fugiat aliqua.
					</p>
					<div className="mt-10 gap-x-6 flex items-center justify-center">
						<a
							href="/signup"
							className="font-semibold leading-5 py-3 px-4 bg-indigo-600 text-white rounded-md"
						>
							Get started
						</a>
						<a href="#" className="leading-6 font-semibold text-sm">
							Learn more <span aria-hidden="true">→</span>
						</a>
					</div>
				</div>
			</div>
		</div>
	);
}
