"use client"

import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles,
} from "lucide-react"

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"

export function NavUser({
  user,
  isSuperadmin = false,
}: {
  user: {
    name: string
    email: string
    avatar: string
  }
  isSuperadmin?: boolean
}) {
  const { isMobile } = useSidebar()

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          {/* @ts-ignore */}
          <DropdownMenuTrigger
            render={(props: any) => (
              <SidebarMenuButton
                {...props}
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground p-2"
              >
                <Avatar className="h-9 w-9 rounded-full">
                  <AvatarImage src={user.avatar} alt={user.name} />
                  <AvatarFallback className="bg-primary/5 text-primary font-semibold tracking-wider">
                    {user.name.slice(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col flex-1 pl-1 text-left justify-center pb-0.5">
                  <span className="truncate text-sm font-semibold tracking-tight text-foreground">{user.name}</span>
                  <span className="truncate text-xs font-medium text-muted-foreground">{user.email}</span>
                </div>
                <ChevronsUpDown className="ml-auto size-4 text-muted-foreground/60" />
              </SidebarMenuButton>
            )}
          />
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuGroup>
              <DropdownMenuLabel className="p-0 font-normal">
                <div className="flex items-center gap-3 px-1.5 py-2 text-left text-sm">
                  <Avatar className="h-9 w-9 rounded-full">
                    <AvatarImage src={user.avatar} alt={user.name} />
                    <AvatarFallback className="bg-primary/5 text-primary font-semibold tracking-wider">
                      {user.name.slice(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col flex-1 text-left justify-center pb-0.5">
                    <span className="truncate text-sm font-semibold tracking-tight text-foreground">{user.name}</span>
                    <span className="truncate text-xs font-medium text-muted-foreground">{user.email}</span>
                  </div>
                </div>
              </DropdownMenuLabel>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            {!isSuperadmin && (
              <>
                <DropdownMenuGroup>
                  <DropdownMenuItem>
                    <Sparkles className="mr-2" />
                    Upgrade to Pro
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
              </>
            )}
            <DropdownMenuGroup>
              <a href={isSuperadmin ? "/superadmin/account" : "/account"} className="w-full">
                <DropdownMenuItem className="cursor-pointer">
                  <BadgeCheck className="mr-2" />
                  Account
                </DropdownMenuItem>
              </a>
              <a href={isSuperadmin ? "/superadmin/billing" : "/billing"} className="w-full">
                <DropdownMenuItem className="cursor-pointer">
                  <CreditCard className="mr-2" />
                  Billing
                </DropdownMenuItem>
              </a>
              <a href={isSuperadmin ? "/superadmin/notifications" : "/notifications"} className="w-full">
                <DropdownMenuItem className="cursor-pointer">
                  <Bell className="mr-2" />
                  Notifications
                </DropdownMenuItem>
              </a>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <LogOut />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
