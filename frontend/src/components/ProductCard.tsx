import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";

type Product = {
	id: number;
	name: string;
	category: string;
};

export function ProductCard(props: Product) {
	const { name, category } = props;
	return (
		<Card className="w-full flex flex-col justify-evenly items-center">
			<CardHeader>
				<CardTitle>{name}</CardTitle>
			</CardHeader>
			<CardContent>
				<CardDescription>{category}</CardDescription>
			</CardContent>
			<CardFooter>
				<Button variant="default">Add to cart</Button>
			</CardFooter>
		</Card>
	);
}
