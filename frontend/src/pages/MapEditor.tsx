import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { motion, useAnimationControls } from "framer-motion";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";

export function MapEditor() {
	const [sliderX, setSliderX] = useState(10);
	const [sliderY, setSliderY] = useState(10);

	const [tool, setTool] = useState("product");
	const [clickedSquare, setClickedSquare] = useState(null);

	const controls = useAnimationControls();

	useEffect(() => {
		controls.start((i) => {
			console.log(i);

			return { scale: 1 };
		});
	}, [sliderY, sliderX]);

	const emptyMap = Array.from({ length: sliderX }).map(() =>
		Array.from({ length: sliderY }),
	);

	const grid = emptyMap.map((row, rowIndex) => (
		<div key={rowIndex} className="flex">
			{row.map((_, colIndex) => (
				<motion.div
					key={colIndex}
					className={`m-1 h-4 w-4 shadow-md round-[${Math.floor(
						Math.random() * 20,
					)}] bg-slate-500`}
					initial={{ scale: 0 }}
					animate={controls}
					onClick={(e) => {
						console.log("Clicked");
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

			<RadioGroup defaultValue="option-one">
				<div className="flex items-center space-x-2">
					<RadioGroupItem
						value="empty"
						id="empty"
						onChange={(e) => console.log(e)}
					/>
					<Label htmlFor="empty">Empty</Label>
				</div>
				<div className="flex items-center space-x-2">
					<RadioGroupItem
						value="blockade"
						id="blockade"
						onChange={(e) => console.log(e)}
					/>
					<Label htmlFor="blockade">Blockade</Label>
				</div>
				<div className="flex items-center space-x-2">
					<RadioGroupItem
						value="product"
						id="product"
						onChange={(e) => console.log(e)}
					/>
					<Label htmlFor="product">Product</Label>
				</div>
			</RadioGroup>

			{grid}
		</>
	);
}
