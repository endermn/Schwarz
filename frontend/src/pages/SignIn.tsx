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

	const signedIn = await authProvider.signin(
		username as string,
		password as string,
	);
	if (!signedIn) {
		return {
			signingIn: "Грешно потребителско име или парола",
		} as Errors;
	}

	localStorage.setItem("user", JSON.stringify({ username }));
	return redirect("/");
}

export function SignIn() {
	let navigation = useNavigation();
	let isLoggingIn = navigation.formData?.get("username") != null;

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
						</Label>{" "}
					</div>
					<div>
						<Label>
							Парола:{" "}
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
						<Button
							className="mb-1 w-full rounded bg-blue-500 py-2 text-white"
							type="submit"
							disabled={isLoggingIn}
						>
							{isLoggingIn ? "Влизане..." : "Влез"}
						</Button>
						<Label>
							{actionData && actionData.signingIn ? (
								<p style={{ color: "red" }}>{actionData.signingIn}</p>
							) : null}
						</Label>
					</div>

					<div className="text-center text-sm text-muted-foreground">
						<p>
							Нямаш профил?{" "}
							<a href="/signup" className="underline">
								Регистрирай се
							</a>
						</p>
					</div>
				</Form>
			</div>
		</div>
	);
}
