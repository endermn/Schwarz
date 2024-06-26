import { useLoaderData } from "react-router-dom";
import { motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import { HandIcon, PlusIcon } from "lucide-react";

// eslint-disable-next-line react-refresh/only-export-components
export async function loader() {
	const resMap = await fetch("http://localhost:12345/stores/0/layout");
	const dataMap = await resMap.json();
	console.log("MAP HERE \n", dataMap);
	const resPath = await fetch("http://localhost:12345/stores/0/find-route", {
		method: "POST",
		body: JSON.stringify({
			products: [3, 12, 43],
		}),
	});
	const dataPath = await resPath.json();
	console.log(dataPath);
	return { dataMap, dataPath };
}

export const MyComponent = () => (
	<motion.div
		className="w-56 h-56 bg-slate-700"
		animate={{
			scale: [1, 2, 2, 1, 1],
			rotate: [0, 0, 270, 270, 0],
			borderRadius: ["20%", "20%", "50%", "50%", "20%"],
		}}
	/>
);

// src/Grid.js
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

const Grid = ({ gridData, path }: { gridData: DataI[][]; path: PathI }) => {
	const [selectedProductId, setSelectedProductId] = useState(0);
	const [cellSize, setCellSize] = useState(20);
	const mapContainerRef = useRef<HTMLDivElement | null>(null);

	console.log(selectedProductId);
	const handleTap = (productId: number) => {
		setSelectedProductId(productId);
	};

	useEffect(() => {
		const updateCellSize = () => {
			let { clientWidth } = mapContainerRef.current!; // Use getBoundingClientRect for precise width
			if (clientWidth < 600) {
				clientWidth *= 1.2;
			} else {
				clientWidth *= 0.5;
			}
			const cols = gridData[0]?.length || 1;
			const cellWidth = Math.floor(clientWidth / cols); // Round cellWidth to an integer
			setCellSize(cellWidth);
		};

		updateCellSize();
		window.addEventListener("resize", updateCellSize);
		return () => window.removeEventListener("resize", updateCellSize);
	}, [gridData]);

	useEffect(() => console.log(cellSize), [cellSize]);

	const grid = gridData.map((row, rowIndex) => (
		<div key={rowIndex} className="flex">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={` md:m-1 m-[1px] shadow-md round-[${Math.floor(
						Math.random() * 20
					)}]  ${getColorFromKind(cell.kind, colIndex, rowIndex, path)}`}
					initial={{ scale: 0 }}
					animate={{
						scale: 1,
					}}
					transition={{
						duration: 0.2,
						delay:
							rowIndex * 0.04 + colIndex * 0.04 * (cell.kind !== 0 ? 0 : 1),
					}}
					onHoverStart={() => {
						if (cell.kind === 3) handleTap(cell.productId);
					}}
					style={{ width: cellSize, height: cellSize }}
				/>
			))}
		</div>
	));

	return (
		<div className="flex justify-center items-center h-full rotate-90 md:rotate-0">
			<div className="grid grid-cols-1 md:grid-cols-4 w-full md:min-h-[80vh]">
				<div className="col-span-1 flex md:flex-col justify-center items-center">
					<div className="bg-slate-700 justify-center items-center rounded-lg cursor-pointer size-10 m-2 hidden md:flex">
						<HandIcon className="text-white" />
					</div>
					<div className="bg-slate-700 hidden md:flex justify-center items-center rounded-lg cursor-pointer size-10 m-2">
						<PlusIcon className="text-white" />
					</div>
				</div>
				<div
					className="col-span-3 flex flex-col items-center justify-center"
					ref={mapContainerRef}
				>
					{grid}
					<h1 className="hidden md:hidden">{selectedProductId}</h1>
				</div>
			</div>
		</div>
	);
};

const getColorFromKind = (kind: number, x: number, y: number, path: PathI) => {
	const good = path.path.find((point) => point.x === x && point.y === y);
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
	}
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
		default:
			return "bg-gray-300";
	}
};

export function Map() {
	const { dataPath, dataMap } = useLoaderData() as {
		dataMap: DataI[][];
		dataPath: PathI;
	};

	console.log(dataPath);

	return (
		<>
			<Grid gridData={dataMap} path={dataPath} />
			<canvas id="map" className="hidden"></canvas>
		</>
	);
}
