import { Outlet } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";

function App() {
	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<div className="min-h-screen flex flex-col justify-between">
				<div>
					<NavBar />
					<Outlet />
				</div>
				<Footer />
			</div>
		</ThemeProvider>
	);
}

export default App;
