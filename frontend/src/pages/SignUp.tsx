import { fakeAuthProvider } from "@/auth";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePersistentStorageValue } from "@/lib/useLocalStorage";

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
		password as string,
	);
	if (!signedUp) {
		return {
			signingUp: "There was an error when signing up",
		} as Errors;
	}

	localStorage.setItem("user", JSON.stringify({ username }));
	return null;
}

export function SignUp() {
	let navigation = useNavigation();
	let isSigningUp = navigation.formData?.get("username") != null;

	let actionData = useActionData() as Errors;

	if (
		actionData &&
		!actionData.password &&
		!actionData.password &&
		!actionData.username
	) {
		usePersistentStorageValue("user", {
			username: navigation.formData?.get("username"),
		});
	}
	return (
		<div className="flex h-full items-center justify-center">
			<div className="w-full max-w-sm">
				<Form method="post" replace className="flex flex-col gap-y-7">
					<div>
						<Label>
							Потребителско име:{" "}
							<Input
								className="dark:bg-white dark:text-black"
								name="username"
							/>
							{actionData && actionData.username ? (
								<p style={{ color: "red" }}>{actionData.username}</p>
							) : null}
						</Label>
					</div>
					<div>
						<Label>
							Парола:
							<Input
								className="dark:bg-white dark:text-black"
								name="password"
							/>
							{actionData && actionData.password ? (
								<p style={{ color: "red" }}>{actionData.password}</p>
							) : null}
						</Label>{" "}
					</div>
					<div className="my-2 flex flex-col justify-center rounded-lg">
						<button
							className="my-2 w-full bg-blue-500 py-2 text-white"
							type="submit"
							disabled={isSigningUp}
						>
							{isSigningUp ? "Влизане..." : "Влез"}
						</button>
						<Label>
							{actionData && actionData.signingUp ? (
								<p style={{ color: "red" }}>{actionData.signingUp}</p>
							) : null}
						</Label>
					</div>

					<div className="text-center text-sm text-muted-foreground">
						<p>
							Имаш профил?{" "}
							<a href="/signin" className="underline">
								Влез
							</a>
						</p>
					</div>
				</Form>
			</div>
		</div>
	);
}
