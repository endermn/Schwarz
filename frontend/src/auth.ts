interface AuthProvider {
	signin(username: string, password: string): Promise<boolean>;
	signup(username: string, password: string): Promise<boolean>;
	signout(): Promise<boolean>;
}

export const authProvider: AuthProvider = {
	async signin(username, password) {
		const res = await fetch("http://localhost:3000/api/login", {
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ username, password }),
			method: "POST",
			credentials: "include",
		});

		if (res.status == 401) {
			return false;
		}

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

		return true;
	},

	async signup(username, password) {
		const res = await fetch("http://localhost:3000/api/register", {
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

		return true;
	},
};
