import React from "react";
import ReactDOM from "react-dom/client";
import App from "@/App.tsx";
import "./index.css";
import {
	LoaderFunctionArgs,
	RouterProvider,
	createBrowserRouter,
	redirect,
} from "react-router-dom";
import ErrorPage from "./pages/Error.tsx";
import { Home } from "./pages/Home.tsx";
import { SignIn, loginAction, loginLoader } from "./pages/SignIn.tsx";
import { SignUp, signUpAction, signUpLoader } from "./pages/SignUp.tsx";
import { Products } from "@/pages/Products.tsx";
import { Map, loader as mapLoader, action as mapAction } from "./pages/Map.tsx";
import { loader as productsLoader } from "./pages/Products.tsx";
import { MapEditor } from "./pages/MapEditor.tsx";
import { fakeAuthProvider } from "./auth.ts";

const router = createBrowserRouter([
	{
		path: "/",
		element: <App />,
		errorElement: <ErrorPage />,
		children: [
			{
				errorElement: <ErrorPage />,
				children: [
					{ index: true, element: <Home /> },
					{
						path: "signin/",
						element: <SignIn />,
						loader: loginLoader,
						action: loginAction,
					},
					{
						path: "signup/",
						element: <SignUp />,
						loader: signUpLoader,
						action: signUpAction,
					},
					{
						path: "signout/",
						async loader() {
							console.log("signing out");
							return await fakeAuthProvider.signout();
						},
					},
					{ path: "products/", loader: productsLoader, element: <Products /> },
					{
						path: "map/",
						loader: mapLoader,
						action: mapAction,
						element: <Map />,
					},
					{ path: "map/editor", loader: mapLoader, element: <MapEditor /> },
				],
			},
		],
	},
	{
		path: "/dashboard",
		loader: protectedLoader,
		children: [{ path: "test", element: <h1>Hello</h1> }],
	},
	{
		path: "/user",
		children: [
			{
				path: "signin",
				async action({ request }) {
					const form = await request.formData();
					const username = form.get("username");
					if (!username || typeof username !== "string") {
						return "no username";
					}

					const password = form.get("password");
					if (!password || typeof password !== "string") {
						return "no password";
					}

					const loggedIn = await fakeAuthProvider.signin(username, password);
					if (!loggedIn) {
						redirect("/signin");
					}

					redirect("/dashboard");
				},
			},
		],
	},
]);

// eslint-disable-next-line no-empty-pattern
async function protectedLoader({}: LoaderFunctionArgs) {
	try {
		const res = await fetch("http://localhost:12345/check-session", {
			credentials: "include",
		});

		if (res.status === 200) {
			console.log("Session is valid");
		} else if (res.status === 400) {
			window.location.href = "/signin"; // Redirect to signin
		} else {
			console.error("Unexpected response:", res);
		}
	} catch (error) {
		console.error("Fetch error:", error);
	}

	return null;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>
);
