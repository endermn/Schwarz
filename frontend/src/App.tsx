import { Outlet, useOutletContext } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";
import { UserProvider } from "./lib/UserContext";
import { useState } from "react";
import { ProductI, UserI } from "@/lib/types";

function App() {
	const [cart, setCart] = useState<ProductI[]>([]);

	const addToCart = (product: ProductI) => {
		setCart((prevCart) => [...prevCart, product]);
	};

	const removeFromCart = (id: number) => {
		setCart((prevCart) => prevCart.filter((p) => p.id !== id));
	};

	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<UserProvider>
				<div className="flex flex-col h-screen">
					<NavBar />
					<div className="flex-1">
						<Outlet context={{ cart, addToCart, removeFromCart }} />
					</div>
					<Footer />
				</div>
			</UserProvider>
		</ThemeProvider>
	);
}

export function getUser() {
	return useOutletContext<UserI>();
}

export default App;
