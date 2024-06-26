import { Outlet } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";
import { UserProvider } from "./lib/UserContext";

function App() {
	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<UserProvider>
				<div className="flex flex-col h-screen">
					<NavBar />
					<div className="flex-1">
						<Outlet />
					</div>
					<Footer />
				</div>
			</UserProvider>
		</ThemeProvider>
	);
}

export default App;
