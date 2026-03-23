"use client"

import * as React from "react"
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  LayoutDashboard,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
  Users,
  Wallet,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavProjects } from "@/components/nav-projects"
import { NavUser } from "@/components/nav-user"
import { TeamSwitcher } from "@/components/team-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

// This is sample data.
const data = {
  user: {
    name: "Admin",
    email: "admin@koperasisaas.com",
    avatar: "",
  },
  teams: [
    {
      name: "Koperasi SaaS",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
  ],
  navMain: [
    {
      title: "Dashboard",
      url: "/",
      icon: LayoutDashboard,
      isActive: true,
    },
    {
      title: "Members",
      url: "/members",
      icon: Users,
    },
    {
      title: "Loans",
      url: "/loans",
      icon: Wallet,
    },
    {
      title: "Savings",
      url: "/savings",
      icon: PieChart,
    },
    {
      title: "Reports",
      url: "/reports",
      icon: BookOpen,
      items: [
        {
          title: "Financial Statements",
          url: "#",
        },
        {
          title: "Transaction Logs",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "/settings",
      icon: Settings2,
    },
  ],
  projects: [
    {
      name: "Accounting",
      url: "#",
      icon: Frame,
    },
    {
      name: "Inventory",
      url: "#",
      icon: PieChart,
    },
    {
      name: "HRM",
      url: "#",
      icon: Map,
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher teams={data.teams} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavProjects projects={data.projects} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
