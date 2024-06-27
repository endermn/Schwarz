import { useFetcher, useLoaderData } from "react-router-dom";
import { motion, useAnimationControls } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import { XIcon } from "lucide-react";
import { getUser } from "@/App";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AspectRatio } from "@/components/ui/aspect-ratio";

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

interface PointI {
	x: number;
	y: number;
}

interface DataI {
	kind: number;
	productId: number;
	checkoutName: string;
}
interface PointI {
	x: number;
	y: number;
}

interface PathI {
	path: PointI[];
}

interface DataI {
	kind: number;
	productId: number;
	checkoutName: string;
}

function rotate90(grid: DataI[][]): DataI[][] {
	// Step 1: Transpose the matrix (swap rows and columns)
	const transposedGrid = transpose(grid);

	// Step 2: Reverse each row of the transposed matrix
	const rotatedGrid = transposedGrid.map((row) => row.reverse());

	return rotatedGrid;
}

function transpose(matrix: DataI[][]): DataI[][] {
	// Create a new matrix where rows become columns (transpose operation)
	const rows = matrix.length;
	const cols = matrix[0].length;
	const transposedMatrix: DataI[][] = [];

	for (let j = 0; j < cols; j++) {
		const newRow: DataI[] = [];
		for (let i = 0; i < rows; i++) {
			newRow.push(matrix[i][j]);
		}
		transposedMatrix.push(newRow);
	}

	return transposedMatrix;
}

const Grid = ({ gridData }: { gridData: DataI[][] }) => {
	const [selectedProductId, setSelectedProductId] = useState(-1);
	const handleTap = (productId: number) => {
		setSelectedProductId(productId);
	};

	const [steps, setSteps] = useState(0);

	const user = getUser();
	const controls = useAnimationControls();

	const fetcher = useFetcher();
	let data = fetcher.data?.dataPath; // path repaint is on button submit... which becomes diabled when there are rzero elements

	const productCoords: PointI[] = [];
	const currentPath = data?.path as PointI[];

	for (let i = 0; i < currentPath?.length; i++) {
		let el = gridData[currentPath[i].y][currentPath[i].x];
		if (el.kind === 3) {
			productCoords.push(currentPath[i]);
			el.kind = 43;
		} else {
			el.kind = 42;
		}
	}

	const upTo =
		steps === productCoords.length
			? currentPath?.length
			: currentPath?.findIndex(
					(p) =>
						p.x === productCoords[steps].x && p.y === productCoords[steps].y
			  ) + 1;

	const grid = gridData.map((row, rowIndex) => (
		<div key={rowIndex} className="flex flex-1 w-full">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={` md:m-1 m-[1px] flex-1 shadow-md round-[${Math.floor(
						Math.random() * 20
					)}]  ${getColorFromKind(
						cell.kind,
						colIndex,
						rowIndex,
						currentPath?.slice(0, upTo)
					)}`}
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
					<Input
						disabled={user.cart.length === 0}
						className="w-1/2 mb-2"
						value={steps}
						type="number"
						placeholder="Map size"
						onChange={(e) => {
							if (
								+e.target.value <= productCoords.length &&
								+e.target.value >= 0
							) {
								setSteps(+e.target.value);
							}
						}}
					/>
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
					<h1 className="hidden md:hidden">{selectedProductId}</h1>
				</div>
			</div>
		</div>
	);
};

const getColorFromKind = (
	kind: number,
	x: number,
	y: number,
	path: PointI[] | null
) => {
	/*
	const good = path?.find((point) => point.x === x && point.y === y);
	if (good) {
		if (kind === 0) return "bg-red-300";
		switch (kind) {
			case 1:
				return `bg-gradient-to-r from-blue-500 to-red-300`;
			case 2:
				return `bg-gradient-to-r from-bg-green-500 to-red-300`;
			case 3:
				return `bg-gradient-to-r from-yellow-500 to-red-900`;
			case 4:
				return `bg-gradient-to-r from-purple-500 to-red-300`;
			default:
				return "bg-gray-300";
		}
	}*/
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
		case 42:
			return "bg-cyan-500";
		case 43:
			return "bg-cyan-200";
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
