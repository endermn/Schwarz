import { authProvider } from "@/auth";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

import {
	Form,
	LoaderFunctionArgs,
	redirect,
	useActionData,
	useNavigation,
} from "react-router-dom";

import { Button } from "@/components/ui/button";
import { Eye, EyeOff } from "lucide-react";

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

	const signedUp = await authProvider.signup(
		username as string,
		password as string,
	);
	if (!signedUp) {
		return {
			signingUp: "Възникна проблем повреме на регистрацията",
		} as Errors;
	}

	localStorage.setItem("user", JSON.stringify({ username }));
	return redirect("/");
}

export function SignUp() {
	let navigation = useNavigation();
	let isSigningUp = navigation.formData?.get("username") != null;

	let actionData = useActionData() as Errors;

	let [showPassword, setShowPassword] = useState(true);

	return (
		<div className="flex h-full items-center justify-center">
			<div className="w-full max-w-sm">
				<Form method="post" replace className="flex flex-col gap-y-7">
					<div>
						<Label>
							Потребителско име:{" "}
							<Input
								className="mb-1 dark:bg-white dark:text-black"
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
							<div className="relative flex items-center">
								<Input
									className="mb-1 dark:bg-white dark:text-black"
									type={showPassword ? "text" : "password"}
									name="password"
								/>

								<Button
									type="button"
									variant="ghost"
									onClick={() => setShowPassword(!showPassword)}
									className="absolute right-0 ml-1 hover:bg-transparent dark:bg-transparent dark:text-black"
								>
									{showPassword ? <Eye /> : <EyeOff />}
								</Button>
							</div>
							{actionData && actionData.password ? (
								<p style={{ color: "red" }}>{actionData.password}</p>
							) : null}
						</Label>{" "}
					</div>
					<div className="flex flex-col justify-center rounded-lg">
						<button
							className="mb-1 w-full rounded bg-blue-500 py-2 text-white"
							type="submit"
							disabled={isSigningUp}
						>
							{isSigningUp ? "Регистриране..." : "Регистрирай се"}
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
