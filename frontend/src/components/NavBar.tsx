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
import { MenuIcon, ShoppingCart, XIcon } from "lucide-react";
import { UserNav } from "./UserNav";
import { ProductI, UserI } from "@/lib/types";
import { NavLink } from "react-router-dom";
import { useState, forwardRef } from "react";

const pages = [
	{
		title: "Начало",
		href: "/",
	},
	{
		title: "Продукти",
		href: "/products",
	},
	{
		title: "Карта",
		href: "/map",
	},
	{
		title: "Редактор",
		href: "/map/editor",
	},
];

export function NavBar({ user, cart }: { user: UserI; cart: ProductI[] }) {
	const [open, setOpen] = useState(false);
	return (
		<div className="m-2 flex justify-between p-2">
			<div className="hidden md:block">
				<ModeToggle />
			</div>

			<NavigationMenu className="hidden md:block">
				<NavigationMenuList>
					{pages.map((page) => (
						<NavigationMenuItem
							key={page.title + Math.floor(Math.random() * 20)}
						>
							<NavigationMenuLink
								href={page.href}
								className={navigationMenuTriggerStyle()}
							>
								{page.title}
							</NavigationMenuLink>
						</NavigationMenuItem>
					))}
				</NavigationMenuList>
			</NavigationMenu>

			{/* Mobile Haburger */}
			<div className="order-3 flex md:hidden">
				<ModeToggle />
				<Sheet open={open}>
					<SheetTrigger className="p-1" onClick={() => setOpen(true)}>
						<MenuIcon />
					</SheetTrigger>
					<SheetContent side={"right"} className="flex w-1/2 flex-col">
						<SheetHeader>
							<div className="flex items-center justify-between">
								<SheetTitle>Меню</SheetTitle>
								<XIcon onClick={() => setOpen(false)} className="size-6" />
							</div>
							<Separator />
							<SheetDescription className="flex flex-col">
								{pages.map((page) => (
									<NavLink key={page.title} to={page.href}>
										<Button
											onClick={() => setOpen(false)}
											variant="ghost"
											className="m-1 w-full"
										>
											{page.title}
										</Button>
									</NavLink>
								))}
							</SheetDescription>
						</SheetHeader>
					</SheetContent>
				</Sheet>
			</div>

			<div className="flex items-center gap-5">
				<NavLink to={cart.length === 0 ? "/products" : "/map"}>
					<div className="relative inline-block">
						<ShoppingCart size={28} />
						{cart.length > 0 && (
							<span className="absolute right-0 top-0 -translate-y-1/2 translate-x-1/2 transform rounded-full bg-green-500 px-1 text-xs font-bold text-white">
								{cart.length}
							</span>
						)}
					</div>
				</NavLink>
				{user?.username ? (
					<UserNav user={user} />
				) : (
					<a href="/signin">
						<Button className="">Влез</Button>
					</a>
				)}
			</div>
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
