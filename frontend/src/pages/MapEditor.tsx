import { useState } from "react";
import { Input } from "@/components/ui/input";
import { motion } from "framer-motion";

export function MapEditor() {
	const [sliderX, setSliderX] = useState(10);
	const [sliderY, setSliderY] = useState(10);

	const emptyMap = Array.from({ length: sliderX }).map(() =>
		Array.from({ length: sliderY })
	);

	const grid = emptyMap.map((row, rowIndex) => (
		<div key={rowIndex} className="flex">
			{row.map((cell, colIndex) => (
				<motion.div
					key={colIndex}
					className={`w-4 h-4 m-1 shadow-md round-[${Math.floor(
						Math.random() * 20
					)}]  bg-slate-500`}
					initial={{ scale: 0 }}
					animate={{
						scale: 1,
					}}
					/*onHoverStart={() => {
						if (cell.kind === 3) handleTap(cell.productId);
					}}*/
				/>
			))}
		</div>
	));

	return (
		<>
			<Input
				className="w-1/2"
				value={sliderX}
				type="number"
				placeholder="Map size"
				onChange={(e) => {
					setSliderX(+e.target.value);
				}}
			/>
			<Input
				className="w-1/2"
				value={sliderY}
				type="number"
				placeholder="Map size"
				onChange={(e) => {
					setSliderY(+e.target.value);
				}}
			/>

			{grid}
		</>
	);
}
