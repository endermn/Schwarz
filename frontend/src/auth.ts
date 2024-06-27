interface AuthProvider {
	isAuthenticated: boolean;
	role: string;
	username: string | null;
	signin(username: string, password: string): Promise<boolean>;
	signup(username: string, password: string): Promise<boolean>;
	signout(): Promise<boolean>;
}

export const fakeAuthProvider: AuthProvider = {
	isAuthenticated: false,
	username: null,
	role: "user",

	async signin(username, password) {
		const res = await fetch("http://localhost:3000/api/login", {
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ username, password }),
			method: "POST",
			credentials: "include",
		});

		if (res.status == 400) {
			return false;
		}

		const data = await res.json();
		this.username = data["username"];
		this.role = data["role"];

		return true;
	},

	async signout() {
		const res = await fetch("http://localhost:3000/api/logout", {
			method: "POST",
			credentials: "include",
		});

		if (res.status != 202) {
			return false;
		}

		this.isAuthenticated = false;
		this.username = null;
		this.role = "user";

		return true;
	},

	async signup(username, password) {
		const res = await fetch("http://localhost:12345/users", {
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ username, password }),
			method: "POST",
			credentials: "include",
		});

		if (res.status == 400) {
			return false;
		}

		const data = await res.json();
		this.username = data["username"];
		this.role = data["role"];

		return true;
	},
};
