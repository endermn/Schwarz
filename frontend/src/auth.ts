interface AuthProvider {
	isAuthenticated: boolean;
	role: string;
	username: string | null;
	signin(username: string, password: string): Promise<boolean>;
	signup(username: string, password: string): Promise<boolean>;
	signout(): Promise<void>;
}

export const fakeAuthProvider: AuthProvider = {
	isAuthenticated: false,
	username: null,
	role: "user",

	async signin(username, password) {
		const res = await fetch("http://localhost:12345/users/login", {
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ username, password }),
			method: "POST",
		});

		if (res.status == 400) {
			return false;
		}

		fakeAuthProvider.isAuthenticated = true;
		fakeAuthProvider.username = username;

		return true;
	},

	async signout() {
		await new Promise((r) => setTimeout(r, 500));
		fakeAuthProvider.isAuthenticated = false;
		fakeAuthProvider.username = null;
	},

	async signup(username, password) {
		const res = await fetch("http://localhost:12345/users", {
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ username, password }),
			method: "POST",
		});

		if (res.status == 400) {
			return false;
		}

		fakeAuthProvider.isAuthenticated = true;
		fakeAuthProvider.username = username;

		return true;
	},
};
