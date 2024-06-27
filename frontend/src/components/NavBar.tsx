"use client";

import { cn } from "@/lib/utils";
import {
	NavigationMenu,
	NavigationMenuItem,
	NavigationMenuLink,
	NavigationMenuList,
	navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { ModeToggle } from "./mode-toggle";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";

import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Link } from "./Link";
import { MenuIcon } from "lucide-react";
import { UserNav } from "./UserNav";
import { UserI } from "@/lib/types";
import { NavLink } from "react-router-dom";
import { useState, forwardRef } from "react";

export function NavBar({ user }: { user: UserI }) {
	console.log(user);
	const [open, setOpen] = useState(false);
	return (
		<div className="m-3 flex justify-between">
			<div className="hidden md:block">
				<ModeToggle />
			</div>

			{/* Big menu */}
			<NavigationMenu className="hidden md:block">
				<NavigationMenuList>
					<NavigationMenuItem>
						<Link href="/">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Начало
							</NavigationMenuLink>
						</Link>
						<Link href="/products">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Продукти
							</NavigationMenuLink>
						</Link>
						<Link href="/map">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Карта
							</NavigationMenuLink>
						</Link>
						<Link href="/map/editor">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Редактор
							</NavigationMenuLink>
						</Link>
					</NavigationMenuItem>
				</NavigationMenuList>
			</NavigationMenu>

			{/* Mobile Haburger */}
			<div className="order-3 md:hidden">
				<ModeToggle />
				<Sheet open={open}>
					<SheetTrigger onClick={() => setOpen(true)}>
						<Button variant={"ghost"}>
							<MenuIcon />
						</Button>
					</SheetTrigger>
					<SheetContent side={"right"} className="flex w-1/3 flex-col">
						<SheetHeader>
							<div className="flex items-center justify-between">
								<SheetTitle>Меню</SheetTitle>
							</div>
							<Separator />
							<SheetDescription className="flex flex-col">
								<NavLink to="/">
									<Button
										onClick={() => setOpen(false)}
										variant="ghost"
										className="m-1 w-full"
									>
										Начало
									</Button>
								</NavLink>
								<NavLink to="/products">
									<Button
										onClick={() => setOpen(false)}
										variant="ghost"
										className="m-1 w-full"
									>
										Продукти
									</Button>
								</NavLink>
								<NavLink to="/map">
									<Button
										onClick={() => setOpen(false)}
										variant="ghost"
										className="m-1 w-full"
									>
										Карта
									</Button>
								</NavLink>
								<NavLink to="/map/editor">
									<Button
										onClick={() => setOpen(false)}
										variant="ghost"
										className="m-1 w-full"
									>
										Редактор
									</Button>
								</NavLink>
							</SheetDescription>
						</SheetHeader>
					</SheetContent>
				</Sheet>
			</div>

			{/* Account or Sign in */}
			{user?.username ? (
				<UserNav user={user} />
			) : (
				<a href="/signin">
					<Button className="">Влез</Button>
				</a>
			)}
		</div>
	);
}

const ListItem = forwardRef<
	React.ElementRef<"a">,
	React.ComponentPropsWithoutRef<"a">
>(({ className, title, children, ...props }, ref) => {
	return (
		<li>
			<NavigationMenuLink asChild>
				<a
					ref={ref}
					className={cn(
						"block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
						className,
					)}
					{...props}
				>
					<div className="text-sm font-medium leading-none">{title}</div>
					<p className="line-clamp-2 text-sm leading-snug text-muted-foreground">
						{children}
					</p>
				</a>
			</NavigationMenuLink>
		</li>
	);
});
ListItem.displayName = "ListItem";
