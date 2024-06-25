import { useLoaderData } from "react-router-dom";
import { motion } from "framer-motion";

export async function loader() {
	const res = await fetch("http://localhost:12345/store-layout");
	const data = await res.json();
	return { data };
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

const Grid = ({ gridData }: { gridData: DataI[][] }) => {
	const numRows = gridData.length;
	const numCols = gridData[0].length;

	return (
		<div className="flex flex-col items-center justify-center h-screen">
			{gridData.map((row, rowIndex) => (
				<div key={rowIndex} className="flex">
					{row.map((cell, colIndex) => (
						<motion.div
							key={colIndex}
							className={`w-4 h-4 m-1 shadow-md ${getColorFromKind(cell.kind)}`}
							initial={{ scale: 0 }}
							animate={{ scale: 1 }}
							transition={{
								duration: 0.2,
								delay:
									rowIndex * 0.08 + colIndex * 0.08 * (cell.kind !== 0 ? 0 : 1),
							}}
							whileTap={
								cell.kind === 3
									? {
											scale: 4,
											transition: { duration: 0.2 },
									  }
									: null
							}
							onTap
						/>
					))}
				</div>
			))}
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
		default:
			return "bg-gray-300";
	}
};

export function Map() {
	const { data } = useLoaderData() as { data: DataI[][] };
	console.log(data);

	return (
		<>
			<h1>Map</h1>
			<Grid gridData={data} />
			<canvas id="map"></canvas>
		</>
	);
}
