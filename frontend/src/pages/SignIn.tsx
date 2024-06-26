import { fakeAuthProvider } from "@/auth";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Form,
	LoaderFunctionArgs,
	redirect,
	useActionData,
	useNavigation,
} from "react-router-dom";

export async function loginLoader() {
	if (fakeAuthProvider.isAuthenticated) {
		return redirect("/dashboard");
	}

	return null;
}

type Errors = {
	password?: string;
	username?: string;
	signingIn?: string;
};

export async function loginAction({ request }: LoaderFunctionArgs) {
	let formData = await request.formData();
	let errors: Errors = {};

	let username = formData.get("username");
	if (!username || typeof username !== "string") {
		errors.username = "Моля въведете потребителско име";
	}

	let password = formData.get("password");
	if (!password || typeof password !== "string") {
		errors.password = "Моля въведете парола";
	}

	if (
		typeof errors.password === "string" ||
		typeof errors.username === "string"
	) {
		return errors;
	}

	const signedIn = await fakeAuthProvider.signin(
		username as string,
		password as string
	);
	if (!signedIn) {
		return {
			signingIn: "Yo shit aint right dawg",
		} as Errors;
	}

	return null;
}

export function SignIn() {
	let navigation = useNavigation();
	let isLoggingIn = navigation.formData?.get("username") != null;

	let actionData = useActionData() as Errors;

	return (
		<div className="flex justify-center items-center h-full">
			<div>
				<Form method="post" replace>
					<input type="hidden" name="redirectTo" />
					<Label>
						Username: <input name="username" />
						{actionData && actionData.username ? (
							<p style={{ color: "red" }}>{actionData.username}</p>
						) : null}
					</Label>{" "}
					<Label>
						Password: <input name="password" />
						{actionData && actionData.password ? (
							<p style={{ color: "red" }}>{actionData.password}</p>
						) : null}
					</Label>{" "}
					{actionData && actionData.signingIn ? (
						<p style={{ color: "red" }}>{actionData.signingIn}</p>
					) : null}
					<button type="submit" disabled={isLoggingIn}>
						{isLoggingIn ? "Logging in..." : "Login"}
					</button>
				</Form>
			</div>
		</div>
	);
}
