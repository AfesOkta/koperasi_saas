"use client"

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { PlanCards } from "@/app/superadmin/billing/components/PlanCards"

export default function RegularBillingPage() {
  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b bg-background/95 backdrop-blur-sm sticky top-0 z-10 px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="/" className="font-medium">
                Koperasi SaaS
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Billing & Subscription</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Subscription Plans</h2>
          <p className="text-muted-foreground mt-1">Review your current plan and explore higher tiers.</p>
        </div>

        <PlanCards isSuperadmin={false} />
        
      </div>
    </>
  )
}
