import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
} from "@/components/ui/carousel";
import Autoplay from "embla-carousel-autoplay";
import { AspectRatio } from "@/components/ui/aspect-ratio";

type Partner = {
	name: string;
	image: string;
};

const partners: Partner[] = [
	{
		name: "Lidl",
		image: "/partners/Lidl.png",
	},
	{
		name: "Kaufland",
		image: "/partners/Kaufland.png",
	},
	{
		name: "Schwarz",
		image: "/partners/Schwarz.png",
	},
	{
		name: "Spge",
		image: "partners/Spge.png",
	},
	{
		name: "Gnu",
		image: "/partners/Gnu.png",
	},
];

export default function PartnersCarousel() {
	return (
		<div className="flex justify-center">
			<Carousel
				opts={{
					align: "center",
					loop: false,
				}}
				plugins={[
					Autoplay({
						delay: 3000,
					}),
				]}
				className="w-3/4"
			>
				<CarouselContent>
					{partners.map((partner) => (
						<CarouselItem
							key={partner.name}
							className="select-none md:basis-1/3"
						>
							<AspectRatio
								ratio={1 / 1}
								className="flex items-center justify-center"
							>
								<img
									src={partner.image}
									alt={partner.name}
									width="75%"
									height="75%"
								/>
							</AspectRatio>
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		</div>
	);
}
