"use client"

import { IconProgressCheck, type Icon } from "@tabler/icons-react"
import Link from "next/link";

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavMain({
  items,
}: {
  items: {
    title: string
    url: string
    icon?: Icon
  }[]
}) {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>Instances</SidebarGroupLabel>
      <SidebarGroupContent className="flex flex-col gap-2">
      <SidebarMenu>
      <Link href="/oauth" passHref>
  <SidebarMenuItem className="flex items-center gap-2 w-full">
      <SidebarMenuButton
        tooltip="Approve Access"
        className="w-full bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground duration-200 ease-linear"
      >
        <IconProgressCheck />
        <span>Approve Access</span>
      </SidebarMenuButton>
  </SidebarMenuItem>
  </Link>
</SidebarMenu>
        <SidebarMenu>
          {items.map((item) => (
            <SidebarMenuItem key={item.title}>
              <Link href={item.url} passHref>
              <SidebarMenuButton tooltip={item.title}>
                {item.icon && <item.icon />}
                <span>{item.title}</span>
              </SidebarMenuButton>
              </Link>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
