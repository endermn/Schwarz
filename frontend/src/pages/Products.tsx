import { ProductCard } from "@/components/ProductCard";
import { ProductView } from "@/components/ProductView";
import { Input } from "@/components/ui/input";
import { Label } from "@radix-ui/react-label";
import { useState } from "react";
import { useLoaderData } from "react-router-dom";

interface ProductI {
	id: number;
	name: string;
	category: string;
}

export async function loader() {
	const resProducts = await fetch("http://localhost:12345/products");
	const products = await resProducts.json();
	console.log(products);

	return { products };
}

export function Products() {
	const { products } = useLoaderData() as { products: ProductI[] };
	const [search, setSearch] = useState("");

	const filteredProducts = products.filter((p) =>
		p.name.toLocaleLowerCase().startsWith(search),
	);
	console.log(search);

	return (
		<>
			<div className="flex justify-center flex-col gap-3 m-4">
				<h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
					Нашите продукти
				</h1>
				<div className="max-w-96">
					<Label htmlFor="name">Всичко от кето имаш нужда</Label>
					<Input
						type="text"
						placeholder="Плодове"
						id="name"
						value={search}
						onChange={(e) => {
							setSearch(e.target.value.toLowerCase());
						}}
					/>
				</div>
			</div>
			<ProductView>
				{filteredProducts.length !== 0 ? (
					filteredProducts.map((p) => (
						<ProductCard
							id={p.id}
							category={p.category}
							name={p.name}
						></ProductCard>
					))
				) : (
					<p>Няма намерени продукти!</p>
				)}
			</ProductView>
		</>
	);
}
