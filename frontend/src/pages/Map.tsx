import { useLoaderData } from "react-router-dom";
import { motion } from "framer-motion";
import { useState } from "react";

export async function loader() {
	const resMap = await fetch("http://localhost:12345/store-layout");
	const dataMap = await resMap.json();
	const resPath = await fetch("http://localhost:12345/find-route", {
		method: "POST",
		body: JSON.stringify({
			products: [3, 12, 43, 51, 61, 123, 210, 89, 69, 101, 10, 44],
		}),
	});
	const dataPath = await resPath.json();
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

	console.log(selectedProductId);
	const handleTap = (productId: number) => {
		setSelectedProductId(productId);
	};

	const grid = gridData.map((row, rowIndex) => (
		<div key={rowIndex} className="flex">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={`w-4 h-4 m-1 shadow-md round-[${Math.floor(
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
				/>
			))}
		</div>
	));

	return (
		<div className="flex flex-col items-center justify-center h-screen">
			{grid}
			<h1>{selectedProductId}</h1>
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
				return `bg-gradient-to-r from-yellow-700 to-red-300`;
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
			<h1>Map</h1>
			<Grid gridData={dataMap} path={dataPath} />
			<canvas id="map"></canvas>
		</>
	);
}
