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

const points = [
	{ x: 1, y: 2 },
	{ x: 23, y: 41 },
	{ x: 12, y: 22 },
	{ x: 7, y: 6 },
];

const Grid = () => {
	const gridColumns = 24;
	const gridRows = 42;
	const squareSize = 20; // Adjust square size as needed

	return (
		<div className="flex justify-center items-center h-screen bg-gray-100">
			<div
				className="relative bg-white border border-gray-300"
				style={{
					width: `${gridColumns * squareSize}px`,
					height: `${gridRows * squareSize}px`,
					display: "grid",
					gridTemplateColumns: `repeat(${gridColumns}, 1fr)`,
					gridTemplateRows: `repeat(${gridRows}, 1fr)`,
				}}
			>
				{points.map((point, index) => (
					<motion.div
						key={index}
						className="absolute w-5 h-5 bg-red-500"
						initial={{ opacity: 0, scale: 0 }}
						animate={{ opacity: 1, scale: 1 }}
						transition={{ duration: 0.5, delay: index * 0.2 }}
						style={{
							left: `${point.x * squareSize}px`,
							top: `${point.y * squareSize}px`,
						}}
					/>
				))}
			</div>
		</div>
	);
};

export function Map() {
	const { data } = useLoaderData() as any;
	console.log(data);

	return (
		<>
			<h1>Map</h1>
			<Grid />
			<canvas id="map"></canvas>
		</>
	);
}
