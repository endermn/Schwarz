import { useLoaderData } from "react-router-dom";

export async function loader() {
	const res = await fetch("localhost:12345/store-layout");
	const data = await res.json();
	return { data };
}

export function Map() {
	const { data } = useLoaderData() as any;
	console.log(data);
	return <h1>Map</h1>;
}
