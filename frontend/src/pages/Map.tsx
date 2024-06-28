import { useFetcher, useLoaderData } from "react-router-dom";
import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { XIcon, ArrowRight, ArrowLeft } from "lucide-react";
import { getContext } from "@/App";
import { Button } from "@/components/ui/button";
import { PointI, DataI, SquareType } from "@/lib/types";
import { ScrollArea } from "@/components/ui/scroll-area";

import {
	AlertDialog,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

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

function deduplicateArray(arr: PointI[]) {
	const seen = new Set();
	return arr.filter((item) => {
		const key = `${item.x},${item.y}`;
		if (seen.has(key)) {
			return false;
		} else {
			seen.add(key);
			return true;
		}
	});
}

const Grid = ({ gridData }: { gridData: DataI[][] }) => {
	const user = getContext();
	const [pathStops, setPathStops] = useState(user.cart.length + 1 + 1); // 1 gold egg and 1 exit
	const [itemRemoved, setItemRemoved] = useState(false);

	const fetcher = useFetcher();

	const divVariants = (delay: number) => ({
		hidden: { scale: 0 },
		visible: { scale: 1, transition: { duration: 0.1, delay } },
	});

	const currentPath = fetcher.data?.dataPath.path as PointI[];

	const pathSlice = useMemo(() => {
		const currentPath = fetcher.data?.dataPath.path as PointI[];
		if (currentPath) {
			const dirtyStops = currentPath.filter((poit) =>
				[
					SquareType.PRODUCT,
					SquareType.PRODUCT_VISITED,
					SquareType.CHECKOUT,
					SquareType.CHECKOUT_VISITED,
					SquareType.SELFCHECKOUT,
					SquareType.SELFCHECKOUT_VISITED,
					SquareType.EXIT,
				].includes(gridData[poit.y][poit.x].kind),
			);

			const stops = deduplicateArray(dirtyStops);
			if (pathStops >= stops.length || pathStops < 0) return [];

			const upTo =
				pathStops === -1
					? 0
					: currentPath.findIndex(
							(p) => p.x === stops[pathStops].x && p.y === stops[pathStops].y,
						) + 1;

			return currentPath.slice(0, upTo);
		}
		return [];
	}, [fetcher, gridData, pathStops]);

	const [prevPath, setPrevPath] = useState<PointI[]>(pathSlice);

	const gridMemo = useMemo<DataI[][]>(() => {
		const gridCopy = JSON.parse(JSON.stringify(gridData));
		if (currentPath && !itemRemoved && fetcher.state === "idle") {
			for (let i = 0; i < pathSlice.length; i++) {
				let el = gridCopy[currentPath[i].y][currentPath[i].x];
				if (i === 0) el.kind = SquareType.START;
				else if (el.kind == SquareType.PRODUCT_VISITED) continue;
				else if (el.kind === SquareType.PRODUCT) {
					el.kind = SquareType.PRODUCT_VISITED;
				} else if (el.kind === SquareType.CHECKOUT) {
					el.kind = SquareType.CHECKOUT_VISITED;
				} else if (el.kind === SquareType.SELFCHECKOUT) {
					el.kind = SquareType.SELFCHECKOUT_VISITED;
				} else if (el.kind === SquareType.EXIT) {
					el.kind = SquareType.EXIT_VISITED;
				} else {
					el.kind = SquareType.VISITED;
				}
			}
		}
		return gridCopy;
	}, [fetcher, pathStops, user.cart, itemRemoved]);

	const grid = gridMemo.map((row, rowIndex) => (
		<div key={rowIndex} className="flex w-full flex-1">
			{row.map((cell, colIndex) => {
				const isAnimated = Object.values(SquareType)
					.filter((t) => t != SquareType.EMPTY)
					.includes(cell.kind);

				const pointIndex = currentPath?.findIndex(
					(square) => square.x === colIndex && square.y === rowIndex,
				);

				let delay = 0;
				if (pointIndex !== undefined && pointIndex !== -1) {
					delay = (pointIndex - prevPath.length) * 0.05;
				}

				return (
					<motion.div
						key={colIndex}
						animate={isAnimated ? "visible" : "hidden"}
						variants={divVariants(delay)}
						initial="hidden"
						className={`m-[1px] flex-1 shadow-md md:m-1 round-[${Math.floor(
							Math.random() * 20,
						)}] ${getColorFromKind(cell.kind)}`}
					/>
				);
			})}
		</div>
	));

	const productLegend = [
		{
			kind: SquareType.START,
			name: "Вход",
		},
		{
			kind: SquareType.BLOCAKDE,
			name: "Стена",
		},
		{
			kind: SquareType.EMPTY,
			name: "Празен път",
		},

		{
			kind: SquareType.VISITED,
			name: "Изминат път",
		},
		{
			kind: SquareType.PRODUCT,
			name: "Продукт",
		},
		{
			kind: SquareType.PRODUCT_VISITED,
			name: "Посетен продукт",
		},
		{
			kind: SquareType.CHECKOUT,
			name: "Каса",
		},
		{
			kind: SquareType.CHECKOUT_VISITED,
			name: "Посетена каса",
		},

		{
			kind: SquareType.SELFCHECKOUT,
			name: "Каса на самообслужване",
		},
		{
			kind: SquareType.SELFCHECKOUT_VISITED,
			name: "Посетена каса на самообслужване",
		},

		{
			kind: SquareType.EXIT,
			name: "Изход",
		},
		{
			kind: SquareType.EXIT_VISITED,
			name: "Изход посетен",
		},
	];

	const [legendOpen, setLegendOpen] = useState(false);
	useEffect(() => {
		const legendSeen = localStorage.getItem("legendSeen");

		if (!legendSeen) setLegendOpen(true);
		localStorage.setItem("legendSeen", "true");
	}, []);

	return (
		<div className="m-5 flex h-full items-center justify-center">
			<div className="grid w-full grid-cols-1 md:min-h-[80vh] lg:grid-cols-4">
				<div className="col-span-3 flex h-[60vw] max-h-[80vh] flex-col items-center justify-center p-5">
					{grid.reverse()}
					<h1 className="hidden md:hidden">{0}</h1>
				</div>
				<div className="col-span-1 flex flex-col items-center justify-between">
					<div className="w-full">
						<h2 className="scroll-m-20 border-b pb-2 text-center text-3xl font-semibold tracking-tight first:mt-0">
							Продукти
						</h2>
						<ScrollArea className="my-6 ml-6 h-[25vh] md:h-[50vh]">
							<AnimatePresence mode="popLayout">
								{user.cart.map((p) => (
									<motion.div
										key={p.id}
										layout
										initial={{ opacity: 0, x: -400, scale: 0.5 }}
										animate={{ opacity: 1, x: 0, scale: 1 }}
										exit={{ opacity: 0, x: 200, scale: 1.2 }}
										transition={{ duration: 0.6, type: "spring" }}
										className="mb-3 flex items-center justify-between rounded-lg border-2 border-black/10 px-4 py-2 dark:border-white/70"
									>
										{p.name}
										<XIcon
											className="inline size-5 cursor-pointer rounded-xl bg-red-500 p-1 text-white"
											onClick={() => {
												user.removeFromCart(p.id);
												setItemRemoved(true);
												setPathStops((prevPath) => prevPath - 1);
											}}
										/>
									</motion.div>
								))}
							</AnimatePresence>
						</ScrollArea>
					</div>

					<div className="flex max-w-[80vw] flex-col gap-y-5">
						<div className="flex gap-4">
							<AlertDialog open={legendOpen}>
								<AlertDialogTrigger>
									<Button
										onClick={() => setLegendOpen(true)}
										variant={"secondary"}
									>
										Помощ?
									</Button>
								</AlertDialogTrigger>
								<AlertDialogContent>
									<AlertDialogHeader>
										<AlertDialogTitle>Легенда на картата</AlertDialogTitle>
										<AlertDialogDescription>
											<div className="grid grid-cols-1 items-center md:grid-cols-2">
												{productLegend.map((product) => {
													return (
														<div
															key={product.name}
															className="flex h-full text-balance border-b-2 border-b-black/20 py-2 dark:border-b-white/50 dark:text-white"
														>
															<div
																className={`mr-3 h-4 w-4 border-2 border-black dark:border-white ${getColorFromKind(product.kind)}`}
															></div>
															{product.name}
														</div>
													);
												})}
											</div>
										</AlertDialogDescription>
									</AlertDialogHeader>
									<AlertDialogFooter>
										<AlertDialogCancel
											className="bg-blue-500 text-white"
											onClick={() => setLegendOpen(false)}
										>
											Разбрах!
										</AlertDialogCancel>
									</AlertDialogFooter>
								</AlertDialogContent>
							</AlertDialog>
							<fetcher.Form method="post">
								<Button
									disabled={user.cart.length === 0}
									name="products"
									value={JSON.stringify({
										products: user.cart.map((p) => p.id),
									})}
									className="bg-green-500 disabled:bg-slate-500"
									onClick={() => {
										setItemRemoved(false);
										setPathStops(user.cart.length + 2);
									}}
								>
									Намери пътя!
								</Button>
							</fetcher.Form>
						</div>

						<div className="flex w-full justify-around">
							<ArrowLeft
								onClick={() => {
									if (pathStops > -1) {
										setPathStops((prevState) => prevState - 1);
									}
								}}
								className="inline size-8 cursor-pointer font-bold"
							/>
							<span>
								{user.cart.length !== 0 && (
									<span>
										{fetcher.data && !itemRemoved ? pathStops + 1 : 0}/{" "}
										{user.cart.length + 3}
									</span>
								)}
							</span>
							<ArrowRight
								onClick={() => {
									if (pathStops < user.cart.length + 1 + 1) {
										// 1 GOLDEN egg, 1 checkout
										setPathStops((prevState) => prevState + 1);
										setPrevPath(pathSlice);
									}
								}}
								className="inline size-8 cursor-pointer font-bold"
							/>
						</div>
					</div>
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
			return "bg-yellow-700";
		case SquareType.CHECKOUT_VISITED:
			return "bg-purple-700";
		case SquareType.SELFCHECKOUT_VISITED:
			return "bg-pink-700";
		case SquareType.START:
			return "bg-green-500";
		case SquareType.START:
			return "bg-green-500";
		case SquareType.EXIT_VISITED:
			return "bg-red-700";
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
