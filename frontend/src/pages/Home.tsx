import PartnersCarousel from "@/components/PartnersCarousel";
import { Link, NavLink } from "react-router-dom";

export function Home() {
	return (
		<>
			<div className="container mx-auto px-4 py-16 sm:py-28 md:py-40 lg:py-48">
				<div id="hero" className="text-center">
					<h1 className="text-4xl font-extrabold sm:text-5xl md:text-6xl lg:text-7xl">
						Пазарувай лесно
					</h1>
					<p className="mt-6 text-lg leading-8 text-muted-foreground sm:text-xl md:text-2xl lg:text-3xl">
						Спри да губиш време в магазина!
					</p>
<<<<<<< HEAD
					<div className="mt-10 flex items-center justify-center gap-x-6">
						<a
							href="/products"
							className="rounded-md bg-indigo-600 px-4 py-3 font-semibold leading-5 text-white"
=======
					<div className="mt-10 flex items-center justify-center gap-4 sm:gap-6">
						<NavLink
							to="/products"
							className="rounded-md bg-indigo-600 px-4 py-3 font-semibold leading-5 text-white hover:bg-indigo-500 sm:text-lg md:text-xl"
>>>>>>> main
						>
							Пазарувай
						</NavLink>
						<a
							href="/whole_squad.jpg"
							target="_blank"
							rel="noopener noreferrer"
							className="text-sm font-semibold leading-6 hover:text-indigo-600 sm:text-lg md:text-xl"
						>
							Научи повече
						</a>
					</div>
				</div>
				<div id="partners-carousel" className="mt-10 hidden md:block">
					<h2 className="text-center text-3xl font-bold sm:text-4xl md:text-5xl">
						Нашите партньори
					</h2>
					<PartnersCarousel />
				</div>
			</div>
		</>
	);
}
