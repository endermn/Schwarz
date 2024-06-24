import * as NavigationMenu from "@radix-ui/react-navigation-menu";
import { NavLink, useLocation } from "react-router-dom";

export const Link = ({ href, ...props }: any) => {
	const location = useLocation();
	const isActive = href === location.pathname;

	return (
		<NavigationMenu.Link asChild active={isActive}>
			<NavLink to={href} className="NavigationMenuLink" {...props} />
		</NavigationMenu.Link>
	);
};
