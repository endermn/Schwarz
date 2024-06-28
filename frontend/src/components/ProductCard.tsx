import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { getContext } from "@/App";
import { ProductI } from "@/lib/types";

export function ProductCard(props: ProductI) {
	const { name, category, id, imageURL } = props;
	const user = getContext();

	const inCart = user.cart.find((p) => p.id === id);

	return (
		<Card className="flex w-full flex-col items-center justify-evenly">
			<CardHeader className="text-center">
				<CardTitle>{name}</CardTitle>
				{/* Uncomment for images */}
				<img src={imageURL} alt="image" />
			</CardHeader>
			<CardContent>
				<CardDescription>{category}</CardDescription>
				<CardDescription className="text-center">Id: {id}</CardDescription>
			</CardContent>
			<CardFooter>
				{inCart ? (
					<Button
						onClick={() => {
							user.removeFromCart(id);
						}}
						variant="destructive"
					>
						Премахни
					</Button>
				) : (
					<Button
						onClick={() => {
							user.addToCart({ category, id, name, imageURL });
						}}
						variant="default"
					>
						Добави към количката
					</Button>
				)}
			</CardFooter>
		</Card>
	);
}
