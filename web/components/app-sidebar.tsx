"use client"

import * as React from "react"
import {
  IconCamera,
  IconUserCheck,
  IconDeviceLaptop,
  IconFileAi,
  IconFileDescription,
  IconTerminal2,
  IconPackage,
  IconUsers,
  IconUser,
} from "@tabler/icons-react"

import { NavDocuments } from "@/components/nav-documents"
import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import { NavUser } from "@/components/nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

const data = {
  user: {
    name: "utkar5hm",
    email: "utkarshrm568@gmail.com",
    avatar: "https://gravatar.com/userimage/87626797/43ad4694a7e1b0538f1ed23bb37cb41d.jpeg?size=256",
  },
  navMain: [
    {
      title: "Instances",
      url: "/instances",
      icon: IconDeviceLaptop,
    },
  ],
  navClouds: [
    {
      title: "Capture",
      icon: IconCamera,
      isActive: true,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
    {
      title: "Proposal",
      icon: IconFileDescription,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
    {
      title: "Prompts",
      icon: IconFileAi,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
  ],
  navSecondary: [
    {
      title: "Get ExecStasy-PAM",
      url: "https://github.com/Utkar5hM/execstasy-pam",
      icon: IconPackage,
    },
  ],
  documents: [
    {
      name: "Roles",
      url: "/roles",
      icon: IconUserCheck,
    },
    {
      name: "Users",
      url: "/users",
      icon: IconUsers,
    },
    {
      name: "My Profile",
      url: "/me",
      icon: IconUser,
    },
  ],
}
import Link from "next/link";
import { useAuth } from "@/components/auth";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const { user } = useAuth();

  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <Link href="/">
                <IconTerminal2 className="!size-5" />
                <span className="text-base font-semibold">ExecStasy</span>
                </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavDocuments items={data.documents} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
    </Sidebar>
  )
}
