import { ProductCard } from "@/components/ProductCard";
import { ProductView } from "@/components/ProductView";
import { Input } from "@/components/ui/input";
import { Label } from "@radix-ui/react-label";
import { useState } from "react";
import { useLoaderData } from "react-router-dom";
import { ProductI } from "@/lib/types";

import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";

export async function loader() {
	const resProducts = await fetch("http://localhost:12345/users/admin/products");
	const products = await resProducts.json();

	const resCategories = await fetch("http://localhost:12345/categories");
	const categories = await resCategories.json();
	console.log(products);

	return { products, categories };
}

export function Products() {
	const { products, categories } = useLoaderData() as {
		products: ProductI[];
		categories: string[];
	};

	const [search, setSearch] = useState("");
	const [category, setCategory] = useState("empty");

	const filteredProducts = products
		.filter((p) => p.category === category || category === "empty")
		.filter((p) => {
			const searchIsntId = isNaN(Number(search)) || search === "";

			if (searchIsntId) {
				return (
					p.name
						.toLocaleLowerCase()
						.split(" ")
						.some((word) => word.startsWith(search)) || p.name.includes(search)
				);
			} else {
				return p.id.toString() === search;
			}
		});

	return (
		<>
			<div className="container">
				<div className="m-4 flex flex-col items-center justify-center gap-3">
					<h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
						Нашите продукти
					</h1>
					<Label htmlFor="name">Нашите вкусни продукти!</Label>
					<div className="flex gap-2">
						<Input
							type="text"
							placeholder="Плодове"
							id="name"
							value={search}
							onChange={(e) => {
								setSearch(e.target.value.toLowerCase());
							}}
						/>

						<Select
							onValueChange={(value: any) => {
								setCategory(value);
								setSearch("");
							}}
						>
							<SelectTrigger className="w-64">
								<SelectValue placeholder="Категории" />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value={"empty"}>Всички</SelectItem>
								{categories.map((c) => (
									<SelectItem value={c}>{c}</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>
				</div>
			</div>
			<ProductView>
				{filteredProducts.length !== 0 ? (
					filteredProducts.map((p) => (
						<ProductCard
							key={p.id}
							id={p.id}
							category={p.category}
							name={p.name}
							imageURL={p.imageURL}
						></ProductCard>
					))
				) : (
					<p>Няма намерени продукти!</p>
				)}
			</ProductView>
		</>
	);
}
