import { useFetcher, useLoaderData } from "react-router-dom";
import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import { XIcon, ArrowRight, ArrowLeft } from "lucide-react";
import { getUser } from "@/App";
import { Button } from "@/components/ui/button";
import { PointI, DataI, SquareType } from "@/lib/types";
import { ScrollArea } from "@/components/ui/scroll-area";

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
	const user = getUser();
	const [pathStops, setPathStops] = useState(user.cart.length + 1);
	const [gridD, setGridD] = useState(gridData);

	const fetcher = useFetcher();

	useEffect(() => {
		const currentPath = fetcher.data?.dataPath.path as PointI[];
		if (currentPath) {
			const stops = currentPath.filter((poit) =>
				[
					SquareType.PRODUCT,
					SquareType.PRODUCT_VISITED,
					SquareType.CHECKOUT,
					SquareType.CHECKOUT_VISITED,
					SquareType.SELFCHECKOUT,
					SquareType.SELFCHECKOUT_VISITED,
					SquareType.EXIT,
				].includes(gridData[poit.y][poit.x].kind)
			);

			console.log(stops);
			console.log(currentPath);

			const upTo =
				pathStops === -1
					? 0
					: currentPath.findIndex(
							(p) => p.x === stops[pathStops].x && p.y === stops[pathStops].y
					  ) + 1;

			const gridCopy = JSON.parse(JSON.stringify(gridData));

			for (let i = 0; i < currentPath.slice(0, upTo).length; i++) {
				let el = gridCopy[currentPath[i].y][currentPath[i].x];
				if (i === 0) el.kind = SquareType.START;
				else if (el.kind === SquareType.PRODUCT) {
					el.kind = SquareType.PRODUCT_VISITED;
				} else if (el.kind === SquareType.CHECKOUT) {
					el.kind = SquareType.CHECKOUT_VISITED;
				} else if (el.kind === SquareType.SELFCHECKOUT) {
					el.kind = SquareType.SELFCHECKOUT_VISITED;
				} else {
					el.kind = SquareType.VISITED;
				}
			}
			setGridD(gridCopy);
		}
	}, [fetcher, pathStops]);

	console.log(gridData);

	const grid = gridD.map((row, rowIndex) => (
		<div key={rowIndex} className="flex w-full flex-1">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={`m-[1px] flex-1 shadow-md md:m-1 round-[${Math.floor(
						Math.random() * 20,
					)}] ${getColorFromKind(cell.kind)}`}
				/>
			))}
		</div>
	));

	return (
		<div className="flex h-full items-center justify-center">
			<div className="grid w-full grid-cols-1 md:min-h-[80vh] lg:grid-cols-4">
				<div className="col-span-1 flex items-center justify-center md:flex-col">
					<h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
						Продукти
					</h2>
					<ScrollArea className="my-6 ml-6 h-[50vh]">
						{user.cart.map((p) => {
							return (
								<li className="my-3">
									{p.name}
									<XIcon
										className="inline size-4 cursor-pointer"
										onClick={() => user.removeFromCart(p.id)}
									/>
								</li>
							);
						})}
					</ScrollArea>
					<div className="flex w-1/4 justify-between">
						<ArrowLeft
							onClick={() => {
								if (pathStops > -1) {
									setPathStops((prevState) => prevState - 1);
								}
							}}
							className="inline size-8 cursor-pointer font-bold"
						/>
						<ArrowRight
							onClick={() => {
								if (pathStops < user.cart.length + 1 + 1) {
									// 1 GOLDEN egg, 1 checkout
									setPathStops((prevState) => prevState + 1);
								}
							}}
							className="inline size-8 cursor-pointer font-bold"
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
				<div className="col-span-3 flex h-[60vw] max-h-[80vh] flex-col items-center justify-center p-5">
					{grid}
					<h1 className="hidden md:hidden">{0}</h1>
				</div>
			</div>
		</div>
	);
};

const getColorFromKind = (kind: number) => {
	switch (kind) {
		case SquareType.EMPTY:
			return "dark:bg-white dark:opacity-30 bg-transparent";
		case SquareType.EXIT:
			return `bg-red-500`;
		case SquareType.BLOCAKDE:
			return "bg-gray-500";
		case SquareType.PRODUCT:
			return "bg-yellow-500";
		case SquareType.CHECKOUT:
			return "bg-purple-500";
		case SquareType.SELFCHECKOUT:
			return "bg-pink-500";
		case SquareType.VISITED:
			return "bg-cyan-500";
		case SquareType.PRODUCT_VISITED:
			return "bg-yellow-200";
		case SquareType.CHECKOUT_VISITED:
			return "bg-purple-200";
		case SquareType.SELFCHECKOUT_VISITED:
			return "bg-pink-200";
		case SquareType.START:
			return "bg-green-500";
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
