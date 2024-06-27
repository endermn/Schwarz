import { Outlet, useOutletContext } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";
import { useState } from "react";
import { ContextI, ProductI, UserI } from "@/lib/types";
import { usePersistentStorageValue } from "./lib/useLocalStorage";

function App() {
	const [cart, setCart] = useState<ProductI[]>([]);
	const [user, setUser] = usePersistentStorageValue<UserI>("user");

	const addToCart = (product: ProductI) => {
		setCart((prevCart) => [...prevCart, product]);
	};

	const removeFromCart = (id: number) => {
		setCart((prevCart) => prevCart.filter((p) => p.id !== id));
	};

	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<div className="flex h-screen flex-col">
				<NavBar user={user} />
				<div className="flex-1">
					<Outlet
						context={{
							cart,
							addToCart,
							removeFromCart,
							paht: null,
							user,
							setUser,
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
