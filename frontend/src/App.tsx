import { Outlet, useLoaderData, useOutletContext } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";
import { useState } from "react";
import { ContextI, ProductI, UserI } from "@/lib/types";

function App() {
	const localStorageCart = localStorage.getItem("products");
	let cartData = [] as ProductI[];
	if (localStorageCart) cartData = JSON.parse(localStorageCart);
	const [cart, setCart] = useState<ProductI[]>(cartData);
	const user = useLoaderData() as UserI; // no username -> no user

	const addToCart = (product: ProductI) => {
		const data = localStorage.getItem("products") as any;
		let cart: ProductI[] = [];
		if (data != null) {
			cart = JSON.parse(data) as ProductI[];
		}

		cart.push(product);
		localStorage.setItem("products", JSON.stringify(cart));
		setCart(cart);
	};

	const removeFromCart = (id: number) => {
		let cardData: ProductI[] = [];
		const localStorageCart = localStorage.getItem("products");
		if (localStorageCart) cardData = JSON.parse(localStorageCart);

		const cleanData = cardData.filter((product) => product.id !== id);
		localStorage.setItem("products", JSON.stringify(cleanData));

		setCart(cleanData);
	};

	const clearCart = () => {
		const newCart: ProductI[] = [];
		localStorage.setItem("products", JSON.stringify(newCart));
		setCart(newCart);
	};

	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<div className="flex h-screen flex-col">
				<NavBar user={user} cart={cart} />
				<div className="flex-1">
					<Outlet
						context={{
							cart,
							addToCart,
							removeFromCart,
							clearCart,
							user,
						}}
					/>
				</div>
				<Footer />
			</div>
		</ThemeProvider>
	);
}

export function getContext() {
	return useOutletContext<ContextI>();
}

export default App;
