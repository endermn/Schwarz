enum SquareType {
	EMPTY = 0,
	EXIT,
	BLOCAKDE,
	PRODUCT,
	CHECKOUT,
	SELFCHECKOUT,
	VISITED,
	PRODUCT_VISITED,
	CHECKOUT_VISITED,
	SELFCHECKOUT_VISITED,
	START,
}

export interface DataI {
	kind: number;
	productId: number;
	checkoutNumber: SquareType;
}

export interface PointI {
	x: number;
	y: number;
}

export interface PathI {
	path: PointI[];
}

export interface ProductI {
	id: number;
	name: string;
	category: string;
}

export interface UserI {
	cart: ProductI[];
	path: PathI | null;
	addToCart: (product: ProductI) => void;
	removeFromCart: (productId: number) => void;
}
