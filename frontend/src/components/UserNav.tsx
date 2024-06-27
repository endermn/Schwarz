import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { UserI } from "@/lib/types";
import { UserIcon } from "lucide-react";
import { Link } from "react-router-dom";

export function UserNav({ user }: { user: UserI }) {
	return (
		<DropdownMenu>
			<DropdownMenuTrigger>
				<Button variant="ghost" className="relative h-10 w-10 rounded-full">
					<Avatar className="h-10 w-10">
						<AvatarFallback>
							<UserIcon />
						</AvatarFallback>
					</Avatar>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56" align="end" forceMount>
				<DropdownMenuLabel className="font-normal">
					<div className="flex flex-col space-y-1">
						<p className="text-sm font-medium leading-none">{user.username}</p>
					</div>
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuSeparator />
				<Link to={"/signout"}>
					<DropdownMenuItem>Излез</DropdownMenuItem>
				</Link>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
