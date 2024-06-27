import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { getUser } from "@/App";

type Product = {
	id: number;
	name: string;
	category: string;
};

export function ProductCard(props: Product) {
	const { name, category, id } = props;
	const user = getUser();

	const inCart = user.cart.find((p) => p.id === id);

	return (
		<Card className="w-full flex flex-col justify-evenly items-center">
			<CardHeader className="text-center">
				<CardTitle>{name}</CardTitle>
			</CardHeader>
			<CardContent>
				<CardDescription>{category}</CardDescription>
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
							console.log("test");
							user.addToCart({ category, id, name });
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
