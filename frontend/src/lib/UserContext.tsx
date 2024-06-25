/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { createContext, useState, useContext, useEffect } from "react";

// Define types for User and Product
interface User {
	role: string;
	isAuthenticated: boolean;
	cart: Product[];
}

export interface Product {
	id: number;
	name: string;
	price: number;
}

// Define context type
interface UserContextType {
	user: User;
	login: () => Promise<void>;
	logout: () => Promise<void>;
	addToCart: (product: Product) => void;
	removeFromCart: (productId: number) => void;
}

const initialUser: User = {
	role: "user",
	isAuthenticated: false,
	cart: [],
};

const UserContext = createContext<UserContextType>({
	user: initialUser,
	login: async () => {},
	logout: async () => {},
	addToCart: (_product) => {},
	removeFromCart: (_productId) => {},
});

interface Props {
	children: React.ReactNode; // Define the type of children prop
}

export const UserProvider: React.FC<Props> = ({ children }) => {
	const [user, setUser] = useState<User>(initialUser);

	useEffect(() => {
		// Check session on initial load
		checkSession();
	}, []);

	const checkSession = async () => {
		try {
			const response = await fetch("/api/check-session"); // Replace with your endpoint
			const data = await response.json();

			if (response.ok) {
				setUser({
					...data,
					isAuthenticated: true,
				});
			} else {
				setUser(initialUser);
			}
		} catch (error) {
			console.error("Error checking session:", error);
			setUser(initialUser);
		}
	};

	const login = async () => {
		try {
			const response = await fetch("/api/login", {
				method: "POST",
				credentials: "include", // Include cookies in the request
			});

			if (response.ok) {
				const data = await response.json();
				setUser({
					...data,
					isAuthenticated: true,
					cart: [],
				});
			} else {
				setUser(initialUser);
			}
		} catch (error) {
			console.error("Error logging in:", error);
			setUser(initialUser);
		}
	};

	const logout = async () => {
		try {
			const response = await fetch("/api/logout", {
				method: "POST",
				credentials: "include", // Include cookies in the request
			});

			if (response.ok) {
				setUser(initialUser);
			} else {
				console.error("Logout failed:", response.statusText);
			}
		} catch (error) {
			console.error("Error logging out:", error);
		}
	};

	const addToCart = (product: Product) => {
		setUser((prevUser) => ({
			...prevUser,
			cart: [...prevUser.cart, product],
		}));
	};

	const removeFromCart = (productId: number) => {
		setUser((prevUser) => ({
			...prevUser,
			cart: prevUser.cart.filter((item) => item.id !== productId),
		}));
	};

	return (
		<UserContext.Provider
			value={{ user, login, logout, addToCart, removeFromCart }}
		>
			{children}
		</UserContext.Provider>
	);
};

export const useUser = (): UserContextType => useContext(UserContext);
