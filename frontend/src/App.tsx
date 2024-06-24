import { Outlet } from "react-router-dom";
import { NavBar } from "./components/NavBar";
import { ThemeProvider } from "./components/theme-provider";

function App() {
	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<NavBar />
			<Outlet />
		</ThemeProvider>
	);
}

export default App;
