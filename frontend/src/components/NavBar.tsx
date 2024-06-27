"use client";

import * as React from "react";

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

export function NavBar() {
	return (
		<div className="flex justify-between m-3">
			<div className="hidden md:block">
				<ModeToggle />
			</div>

			{/* Big menu */}
			<NavigationMenu className="hidden  md:block">
				<NavigationMenuList>
					<NavigationMenuItem>
						<Link href="/">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Начало
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
						<Link href="/products">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Продукти
							</NavigationMenuLink>
						</Link>
					</NavigationMenuItem>
				</NavigationMenuList>
			</NavigationMenu>

			{/* Mobile Haburger */}
			<div className="md:hidden order-3">
				<ModeToggle />
				<Sheet>
					<SheetTrigger>
						<Button variant={"ghost"}>
							<MenuIcon />
						</Button>
					</SheetTrigger>
					<SheetContent side={"right"} className="flex flex-col w-1/3">
						<SheetHeader>
							<div className="flex justify-between items-center">
								<SheetTitle>Меню</SheetTitle>
							</div>
							<Separator />
							<SheetDescription className="flex flex-col">
								<a href="/">
									<Button variant="ghost" className="w-full m-1">
										Начало
									</Button>
								</a>
								<a href="/products">
									<Button variant="ghost" className="w-full m-1">
										Продукти
									</Button>
								</a>
								<a href="/map">
									<Button variant="ghost" className="w-full m-1">
										Карта
									</Button>
								</a>
								<a href="/map/editor">
									<Button variant="ghost" className="w-full m-1">
										Редактор
									</Button>
								</a>
							</SheetDescription>
						</SheetHeader>
					</SheetContent>
				</Sheet>
			</div>

			{/* Account or Sign in */}
			<a href="/signin">
				<Button className="">Влез</Button>
			</a>
		</div>
	);
}

const ListItem = React.forwardRef<
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
