import { useFetcher, useLoaderData } from "react-router-dom";
import { motion } from "framer-motion";
import { useState } from "react";
import { XIcon, ArrowRight, ArrowLeft } from "lucide-react";
import { getUser } from "@/App";
import { Button } from "@/components/ui/button";
import { PointI, DataI } from "@/lib/types";

// eslint-disable-next-line react-refresh/only-export-components
export async function loader() {
	const resMap = await fetch("http://localhost:12345/stores/0/layout");
	const dataMap = await resMap.json();
	return { dataMap };
}

export async function action({ request }: any) {
	const formData = await request.formData();
	let products = formData.get("products");
	console.log(products);
	const resPath = await fetch("http://localhost:12345/stores/0/find-route", {
		method: "POST",
		body: products,
	});
	const dataPath = await resPath.json();
	return { dataPath };
}

const GOLDEN = [170, 130, 240, 119, 239];

const Grid = ({ gridData }: { gridData: DataI[][] }) => {
	const [productsReached, setProductsReached] = useState(0);

	const user = getUser();

	const fetcher = useFetcher();
	const currentPath = fetcher.data?.dataPath.path as PointI[];

	// dear programmer, the code above is grotesque to say the least. There are 2 renders of the component and I can't seem to make it work with useEffect and thus the below code was born
	const checkouts = gridData.flatMap((row) =>
		row.filter((el) => el.kind === 4 || el.kind === 44)
	);

	const selfCheckouts = gridData.flatMap((row) =>
		row.filter((el) => el.kind === 5 || el.kind === 45)
	);

	for (let i = 0; i < currentPath?.length; i++) {
		let el = gridData[currentPath[i].y][currentPath[i].x];
		if (
			user.cart.find((p) => p.id === el.productId) ||
			GOLDEN.includes(el.productId)
		) {
			el.kind = 43; // visited product
		} else if (
			checkouts.map((checkout) => checkout.productId).includes(el.productId)
		) {
			el.kind = 44;
		} else if (
			selfCheckouts.map((checkout) => checkout.productId).includes(el.productId)
		) {
			el.kind = 45;
		} else {
			el.kind = 42; // part of path
		}
	}

	const grid = gridData.map((row, rowIndex) => (
		<div key={rowIndex} className="flex flex-1 w-full">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={` md:m-1 m-[1px] flex-1 shadow-md round-[${Math.floor(
						Math.random() * 20
					)}]  ${getColorFromKind(cell.kind)}`}
				/>
			))}
		</div>
	));

	return (
		<div className="flex justify-center items-center h-full">
			<div className="grid grid-cols-1 lg:grid-cols-4 w-full md:min-h-[80vh]">
				<div className="col-span-1 flex md:flex-col justify-center items-center">
					<h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
						Продукти
					</h2>
					<ul className="my-6 ml-6 list-disc [&>li]:mt-2">
						{user.cart.map((p) => {
							return (
								<li>
									{p.name}
									<XIcon
										className="inline size-4 cursor-pointer"
										onClick={() => user.removeFromCart(p.id)}
									/>
								</li>
							);
						})}
					</ul>
					<div className="flex justify-between w-1/4">
						<ArrowLeft
							onClick={() => {
								if (productsReached > 0) {
									setProductsReached((prevState) => prevState - 1);
								}
							}}
							className="inline size-8 font-bold cursor-pointer"
						/>
						<ArrowRight
							onClick={() => {
								if (productsReached < user.cart.length) {
									setProductsReached((prevState) => prevState + 1);
								}
							}}
							className="inline size-8 font-bold cursor-pointer"
						/>
					</div>
					<fetcher.Form method="post">
						<Button
							disabled={user.cart.length === 0}
							name="products"
							value={JSON.stringify({ products: user.cart.map((p) => p.id) })}
						>
							Намери пътя към светлината!
						</Button>
					</fetcher.Form>
				</div>
				<div className="col-span-3 flex flex-col items-center justify-center p-5 h-[60vw] max-h-[80vh]">
					{grid}
					<h1 className="hidden md:hidden">{0}</h1>
				</div>
			</div>
		</div>
	);
};

const getColorFromKind = (kind: number) => {
	switch (kind) {
		case 0:
			return "dark:bg-white dark:opacity-30 bg-transparent";
		case 1:
			return `bg-blue-500`;
		case 2:
			return "bg-green-500";
		case 3:
			return "bg-yellow-500";
		case 4:
			return "bg-purple-500";
		case 5:
			return "bg-pink-500";
		case 42:
			return "bg-cyan-500";
		case 43:
			return "bg-cyan-200";
		case 44:
			return "bg-purple-200";
		case 45:
			return "bg-pink-200";
		default:
			return "bg-gray-300";
	}
};

export function Map() {
	const { dataMap } = useLoaderData() as {
		dataMap: DataI[][];
	};

	return (
		<>
			<Grid gridData={dataMap} />
			<canvas id="map" className="hidden"></canvas>
		</>
	);
}
