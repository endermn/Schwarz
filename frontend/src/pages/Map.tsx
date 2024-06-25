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
const pointsa = [
	{ x: 1, y: 2 },
	{ x: 3, y: 4 },
	{ x: 5, y: 1 },
	{ x: 7, y: 6 },
	{ x: 2, y: 24 },
];

interface PointI {
	x: number;
	y: number;
}

const kindColors = {
	0: "bg-green-300",
	1: "bg-red-500",
	2: "bg-slate-500",
	3: "bg-blue-500",
	4: "bg-yellow-500",
	// Add more colors if needed
};
interface DataI {
	kind: number;
	productId: number;
	checkoutName: string;
}

const Grid = ({ gridData }: { gridData: DataI[][] }) => {
	const gridColumns = 24;
	const gridRows = 42;
	const squareSize = 20; // Adjust square size as needed

	return (
		<div className="relative flex justify-center bg-gray-100 ">
			<div
				className="relative  bg-white border m-auto"
				style={{
					width: `${gridColumns * squareSize}px`,
					height: `${gridRows * squareSize}px`,
					display: "grid",
					gridTemplateColumns: `repeat(${gridColumns}, 1fr)`,
					gridTemplateRows: `repeat(${gridRows}, 1fr)`,
				}}
			>
				{gridData.map((row, y) =>
					row.map((item, x) => (
						<motion.div
							key={`${x}-${y}`}
							className={`absolute w-5 h-5  ${(kindColors as any)[item.kind]}`}
							initial={{ opacity: 0, scale: 0 }}
							animate={{ opacity: 1, scale: 1 }}
							transition={{
								duration: 0.5,
								delay: (x + y * gridColumns) * 0.00001,
							}}
							style={{
								left: `${x * squareSize}px`,
								top: `${y * squareSize}px`,
							}}
						/>
					))
				)}
			</div>
		</div>
	);
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
