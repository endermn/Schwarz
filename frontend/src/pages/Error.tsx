import { useRouteError } from "react-router-dom";

export default function ErrorPage() {
	const error = useRouteError() as { statusText: string; message: string }; // fix that later

	return (
		<div id="error-page">
			<h1>Опа!</h1>
			<p>Съжаляваме, изникна неочаквана грешка!</p>
			<p>{error.statusText || error.message}</p>
		</div>
	);
}
