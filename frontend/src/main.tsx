import React from "react";
import ReactDOM from "react-dom/client";
import App from "@/App.tsx";
import "./index.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import ErrorPage from "./pages/Error.tsx";
import { Home } from "./pages/Home.tsx";
import { SignIn } from "./pages/SignIn.tsx";
import { SignUp } from "./pages/SignUp.tsx";
import { Products, loader as productsLoader } from "@/pages/Products.tsx";
import { Map, loader as mapLoader } from "./pages/Map.tsx";
import { MapEditor } from "./pages/MapEditor.tsx";

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
					{ path: "signin/", element: <SignIn /> },
					{ path: "products/", loader: productsLoader, element: <Products /> },
					{ path: "map/", loader: mapLoader, element: <Map /> },
					{ path: "map/editor", loader: mapLoader, element: <MapEditor /> },
					{ path: "signup/", element: <SignUp /> },
				],
			},
		],
	},
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>
);
