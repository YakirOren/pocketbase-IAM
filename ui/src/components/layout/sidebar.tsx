import { Link, useLocation } from "react-router";
import { useGetIdentity, useLogout } from "@refinedev/core";
import {
  Shield,
  UserCheck,
  Users,
  Contact,
  Database,
  FlaskConical,
  EllipsisVertical,
  LogOut,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const navItems = [
  { label: "Policies", path: "/policies", icon: Shield },
  { label: "Roles", path: "/roles", icon: UserCheck },
  { label: "Groups", path: "/groups", icon: Users },
  { label: "Users", path: "/users", icon: Contact },
  { label: "Managed Collections", path: "/managed-collections", icon: Database },
];

const toolItems = [
  { label: "Policy Simulator", path: "/simulator", icon: FlaskConical },
];

export function Sidebar() {
  const location = useLocation();
  const { data: identity } = useGetIdentity<{ email?: string }>();
  const { mutate: logout } = useLogout();

  const renderItem = (item: (typeof navItems)[0]) => {
    const isActive = location.pathname.startsWith(item.path);
    return (
      <Link
        key={item.path}
        to={item.path}
        className={cn(
          "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
          isActive
            ? "bg-accent text-accent-foreground"
            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        )}
      >
        <item.icon className="h-4 w-4" />
        {item.label}
      </Link>
    );
  };

  return (
    <aside className="flex h-screen w-60 flex-col border-r bg-background">
      <div className="flex h-14 items-center border-b px-4">
        <Shield className="mr-2 h-5 w-5" />
        <span className="text-lg font-semibold">IAM</span>
      </div>
      <nav className="flex-1 space-y-1 p-3">
        {navItems.map(renderItem)}
        <Separator className="my-3" />
        {toolItems.map(renderItem)}
      </nav>
      <div className="flex items-center gap-3 border-t px-4 py-3">
        <Avatar className="h-8 w-8">
          <AvatarFallback className="text-xs">
            {identity?.email?.charAt(0).toUpperCase() ?? "?"}
          </AvatarFallback>
        </Avatar>
        <span className="min-w-0 flex-1 truncate text-sm text-muted-foreground">
          {identity?.email ?? ""}
        </span>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="rounded p-1 hover:bg-accent">
              <EllipsisVertical className="h-4 w-4 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent side="top" align="end">
            <DropdownMenuItem onClick={() => logout()} className="text-destructive focus:text-destructive">
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>
  );
}
