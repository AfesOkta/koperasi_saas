"use client"

import * as React from "react"
import {
  LayoutDashboard,
  Users,
  CreditCard,
  Settings2,
  Command,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavUser } from "@/components/nav-user"
import { TeamSwitcher } from "@/components/team-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

const data = {
  user: {
    name: "Superadmin",
    email: "superadmin@koperasisaas.com",
    avatar: "",
  },
  teams: [
    {
      name: "Koperasi SaaS Plataform",
      logo: Command,
      plan: "Operator",
    },
  ],
  navMain: [
    {
      title: "Dashboard",
      url: "/superadmin",
      icon: LayoutDashboard,
      isActive: true,
    },
    {
      title: "Billing",
      url: "/superadmin/billing",
      icon: CreditCard,
    },
    {
      title: "Users (Koperasi)",
      url: "/superadmin/users",
      icon: Users,
    },
    {
      title: "Settings",
      url: "/superadmin/settings",
      icon: Settings2,
    },
  ],
}

export function SuperadminSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher teams={data.teams} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} isSuperadmin={true} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
