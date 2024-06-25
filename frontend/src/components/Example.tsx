import React from "react";
import { Product, useUser } from "@/lib/UserContext";

const ExampleComponent: React.FC = () => {
	const { user, addToCart, removeFromCart } = useUser();

	const handleAddToCart = (product: Product) => {
		addToCart(product);
	};

	const handleRemoveFromCart = (productId: number) => {
		removeFromCart(productId);
	};

	return (
		<div>
			<p>User Role: {user.role}</p>
			<p>Is Authenticated: {user.isAuthenticated ? "Yes" : "No"}</p>
			<h2>User Cart</h2>
			<ul>
				{user.cart.map((product) => (
					<li key={product.id}>
						{product.name} - ${product.price}
						<button onClick={() => handleRemoveFromCart(product.id)}>
							Remove
						</button>
					</li>
				))}
			</ul>
			<button
				onClick={() => handleAddToCart({ id: 1, name: "Product A", price: 10 })}
			>
				Add Product A to Cart
			</button>
		</div>
	);
};

export default ExampleComponent;
