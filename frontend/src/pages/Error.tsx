import { useRouteError } from "react-router-dom";

export default function ErrorPage() {
	const error = useRouteError() as { statusText: string; message: string }; // fix that later

	return (
		<>
			<main
				id="error-page"
				className="flex flex-col items-center justify-center min-h-full place-items-center px-6 py-24 sm:py-32 lg:px-8]"
			>
				<img
					src="https://seeklogo.com/images/L/linux-tux-logo-8C1B4FB97E-seeklogo.com.png"
					alt="Готино пингвинче"
					className="m-5 h-[200px] w-[200px]"
				/>
				<div className="flex flex-col items-center justify-center">
					<h1 className="scroll-m-20 text-6xl font-extrabold">Опа!</h1>
					<p className="text-xl m-3 text-muted-foreground">Нещо се счупи!</p>
					<p className="text-xl text-center text-muted-foreground">
						{error.statusText || error.message}
					</p>
				</div>
			</main>
		</>
	);
}
