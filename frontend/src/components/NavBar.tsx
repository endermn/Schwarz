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
import { UserNav } from "./UserNav";
import { ModeToggle } from "./mode-toggle";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";

import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/components/ui/accordion";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Link } from "./Link";
const components: { title: string; href: string; description: string }[] = [
	{
		title: "Sign In",
		href: "/signin",
		description: "Sign in to use all the features of our app!",
	},
	{
		title: "Home",
		href: "/",
		description: "Home sweat home",
	},
	{
		title: "404",
		href: "/404",
		description: "Going to void",
	},
];

export function NavBar() {
	return (
		<div className="flex justify-between m-3">
			<ModeToggle />
			<NavigationMenu className="hidden md:block">
				<NavigationMenuList>
					<NavigationMenuItem>
						<Link href="/">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Home
							</NavigationMenuLink>
						</Link>
						<Link href="/map">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Map
							</NavigationMenuLink>
						</Link>
						<Link href="/map/editor">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Map Editor
							</NavigationMenuLink>
						</Link>
					</NavigationMenuItem>
					<NavigationMenuItem>
						<Link href="/products">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Product List
							</NavigationMenuLink>
						</Link>
					</NavigationMenuItem>
				</NavigationMenuList>
			</NavigationMenu>
			<div className="md:hidden">
				<Sheet>
					<SheetTrigger>
						<Button variant={"outline"}>Menu</Button>
					</SheetTrigger>
					<SheetContent>
						<SheetHeader>
							<SheetTitle>Links</SheetTitle>
							<Separator />

							<SheetDescription>
								<Accordion type="single" collapsible className="w-full">
									{components.map((component) => (
										<AccordionItem value={component.title}>
											<AccordionTrigger>{component.title}</AccordionTrigger>
											<AccordionContent>
												{component.description}
											</AccordionContent>
										</AccordionItem>
									))}
								</Accordion>
							</SheetDescription>
							<Button variant={"outline"}>Documentation</Button>
						</SheetHeader>
					</SheetContent>
				</Sheet>
			</div>
			<UserNav />
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
						className
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
