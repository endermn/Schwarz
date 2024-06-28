import { motion, useDragControls } from "framer-motion";

const controls = useDragControls();

function startDrag(event: any) {
	controls.start(event);
}

function ResizableGrid() {
	return (
		<>
			<div onPointerDown={startDrag} className="bg-red size-10" />
			<motion.div drag="x" dragControls={controls} className="bg-red size-10" />
		</>
	);
}

export default ResizableGrid;
