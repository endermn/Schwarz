import { fakeAuthProvider } from "@/auth";
import { Label } from "@/components/ui/label";

import {
	Form,
	LoaderFunctionArgs,
	redirect,
	useActionData,
	useNavigation,
} from "react-router-dom";

export async function signUpLoader() {
	if (fakeAuthProvider.isAuthenticated) {
		return redirect("/dashboard");
	}

	return null;
}

type Errors = {
	password?: string;
	username?: string;
	signingUp?: string;
};

export async function signUpAction({ request }: LoaderFunctionArgs) {
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

	const signedUp = await fakeAuthProvider.signup(
		username as string,
		password as string
	);
	if (!signedUp) {
		return {
			signingUp: "There was an error when signing up",
		} as Errors;
	}

	return null;
}

export function SignUp() {
	let navigation = useNavigation();
	let isSigningUp = navigation.formData?.get("username") != null;

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
					{actionData && actionData.signingUp ? (
						<p style={{ color: "red" }}>{actionData.signingUp}</p>
					) : null}
					<button type="submit" disabled={isSigningUp}>
						{isSigningUp ? "Signing up..." : "Sign up"}
					</button>
				</Form>
			</div>
		</div>
	);
}
