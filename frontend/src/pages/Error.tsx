import { useRouteError } from "react-router-dom";

export default function ErrorPage() {
	const error = useRouteError() as { statusText: string; message: string }; // fix that later

	return (
		// <div id="error-page">
		// 	<h1>Опа!</h1>
		// 	<p>Съжаляваме, изникна неочаквана грешка!</p>
		// 	<p>{error.statusText || error.message}</p>
		// </div>
		<>
			<main
				id="error-page"
				className="grid min-h-full place-items-center px-6 py-24 sm:py-32 lg:px-8]"
			>
				<img
					src="https://seeklogo.com/images/L/linux-tux-logo-8C1B4FB97E-seeklogo.com.png"
					alt="Готино пингвинче"
					className="m-20"
				/>
				<h1 className="scroll-m-20 text-6xl font-extrabold">Опа!</h1>
				<p className="text-xl text-muted-foreground">
					Нещо се счупи! ({error.statusText || error.message})
				</p>
			</main>
		</>
	);
}
