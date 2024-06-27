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
		name: "",
		image: "",
	},
	{
		name: "",
		image: "",
	},

	{
		name: "Lidl",
		image: "../../public/partners/Lidl.png",
	},
	{
		name: "Kaufland",
		image: "../../public/partners/Kaufland.png",
	},
	{
		name: "Schwarz",
		image: "../../public/partners/Schwarz.png",
	},
	{
		name: "Spge",
		image: "../../public/partners/Spge.jpeg",
	},
	{
		name: "Gnu",
		image: "../../public/partners/Gnu.png",
	},
];

export default function PartnersCarousel() {
	return (
		<div className="flex items-start justify-center">
			<Carousel
				opts={{
					align: "start",
					loop: false,
				}}
				plugins={[
					Autoplay({
						delay: 3000,
					}),
				]}
				className="w-[70%]"
			>
				<CarouselContent className="flex items-center justify-center">
					{partners.map((p, index) => (
						<CarouselItem key={index} className="md:basis-1/3">
							<AspectRatio
								ratio={1 / 1}
								className="flex items-center justify-center"
							>
								<img src={p.image} alt={p.name} width="200" height="200" />
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
