"use client";

import * as React from "react";

import { cn } from "@/lib/utils";
import {
	NavigationMenu,
	NavigationMenuContent,
	NavigationMenuItem,
	NavigationMenuLink,
	NavigationMenuList,
	NavigationMenuTrigger,
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
						<NavigationMenuTrigger>Getting started</NavigationMenuTrigger>
						<NavigationMenuContent>
							<ul className="grid gap-3 p-4 md:w-[400px] lg:w-[500px] lg:grid-cols-[.75fr_1fr]">
								<li className="row-span-3">
									<NavigationMenuLink asChild>
										<a
											className="flex h-full w-full select-none flex-col justify-end rounded-md bg-gradient-to-b from-muted/50 to-muted p-6 no-underline outline-none focus:shadow-md"
											href="/"
										>
											<div className="mb-2 mt-4 text-lg font-medium">
												shadcn/ui
											</div>
											<p className="text-sm leading-tight text-muted-foreground">
												Beautifully designed components built with Radix UI and
												Tailwind CSS.
											</p>
										</a>
									</NavigationMenuLink>
								</li>
								<ListItem
									href="/docs/primitives/typography"
									title="Sample project"
								>
									Random text here
								</ListItem>
							</ul>
						</NavigationMenuContent>
					</NavigationMenuItem>
					<NavigationMenuItem>
						<NavigationMenuTrigger>Components</NavigationMenuTrigger>
						<NavigationMenuContent>
							<ul className="grid w-[400px] gap-3 p-4 md:w-[500px] md:grid-cols-2 lg:w-[600px] ">
								{components.map((component) => (
									<Link
										key={component.title}
										title={component.title}
										href={component.href}
									>
										<ListItem title={component.title}>
											{component.description}
										</ListItem>
									</Link>
								))}
							</ul>
						</NavigationMenuContent>
					</NavigationMenuItem>
					<NavigationMenuItem>
						<a href="/docs">
							<NavigationMenuLink className={navigationMenuTriggerStyle()}>
								Documentation
							</NavigationMenuLink>
						</a>
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
