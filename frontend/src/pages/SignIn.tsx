import { fakeAuthProvider } from "@/auth";
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
		password as string,
	);
	if (!signedIn) {
		return {
			signingIn: "Грешно потребителско име или парола",
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
			<div className="max-w-sm w-full">
				<Form method="post" replace className="flex flex-col gap-y-7 ">
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
						</Label>{" "}
					</div>
					<div>
						<Label>
							Парола:{" "}
							<Input
								className="dark:bg-white dark:text-black"
								name="password"
							/>
							{actionData && actionData.password ? (
								<p style={{ color: "red" }}>{actionData.password}</p>
							) : null}
						</Label>{" "}
					</div>
					<div className="my-2 w-full bg-blue-500 text-white rounded-lg flex justify-center py-2">
						<button type="submit" disabled={isLoggingIn}>
							{isLoggingIn ? "Влизане..." : "Влез"}
						</button>
						{actionData && actionData.signingIn ? (
							<p style={{ color: "red" }}>{actionData.signingIn}</p>
						) : null}
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
