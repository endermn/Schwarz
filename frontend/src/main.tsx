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
import { SignIn, loginAction } from "./pages/SignIn.tsx";
import { SignUp, signUpAction } from "./pages/SignUp.tsx";
import { Products } from "@/pages/Products.tsx";
import { Map, loader as mapLoader, action as mapAction } from "./pages/Map.tsx";
import { loader as productsLoader } from "./pages/Products.tsx";
import { authProvider } from "./auth.ts";

const router = createBrowserRouter([
	{
		path: "/",
		element: <App />,
		loader: () => {
			let data = localStorage.getItem("user");
			if (data !== null) return JSON.parse(data);
			return { username: "" };
		},

		errorElement: <ErrorPage />,
		children: [
			{
				errorElement: <ErrorPage />,
				children: [
					{ index: true, element: <Home /> },
					{
						path: "signin/",
						element: <SignIn />,
						loader(data, stuff) {
							console.log(data, stuff);
							return localStorage.getItem("user") ? redirect("/") : null;
						},
						action: loginAction,
					},
					{
						path: "signup/",
						element: <SignUp />,
						loader() {
							return localStorage.getItem("user") ? redirect("/") : null;
						},
						action: signUpAction,
					},
					{
						path: "signout/",
						async loader() {
							localStorage.removeItem("user");
							return redirect("/");
						},
					},
					{ path: "products/", loader: productsLoader, element: <Products /> },
					{
						path: "map/",
						loader: mapLoader,
						action: mapAction,
						element: <Map />,
					},
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

					const loggedIn = await authProvider.signin(username, password);
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
		const res = await fetch("http://localhost:3000/api/user", {
			credentials: "include",
		});

		if (res.status === 200) {
			console.log("Session is valid");
			return redirect("/");
		} else if (res.status === 401) {
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
	</React.StrictMode>,
);
