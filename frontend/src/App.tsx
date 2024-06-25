import { Outlet } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";
import { Footer } from "./components/Footer";

function App() {
	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<div className="flex flex-col h-screen">
				<NavBar />
				<div className="flex-1">
					<Outlet />
				</div>
				<Footer />
			</div>
		</ThemeProvider>
	);
}

export default App;
