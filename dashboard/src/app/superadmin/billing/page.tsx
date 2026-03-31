"use client"

import { useEffect, useState } from "react"
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
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Plus, RefreshCcw } from "lucide-react"
import { PlanCards } from "./components/PlanCards"
import { InvoiceTable } from "./components/InvoiceTable"
import { PlanFormModal } from "./components/PlanFormModal"
import { useBillingStore, SubscriptionPlan } from "./store/useBillingStore"

export default function BillingPage() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [planToEdit, setPlanToEdit] = useState<SubscriptionPlan | null>(null)
  
  const { fetchPlans, isLoading, error } = useBillingStore()

  useEffect(() => {
    fetchPlans()
  }, [fetchPlans])

  const handleAddClick = () => {
    setPlanToEdit(null)
    setIsModalOpen(true)
  }

  const handleEditClick = (plan: SubscriptionPlan) => {
    setPlanToEdit(plan)
    setIsModalOpen(true)
  }

  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b bg-background/95 backdrop-blur-sm sticky top-0 z-10 px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="/superadmin" className="font-medium">
                Platform Operator
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
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Subscription Plans</h2>
            <p className="text-muted-foreground mt-1">Manage and configure available Koperasi SaaS plans.</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" size="icon" onClick={() => fetchPlans()} disabled={isLoading}>
              <RefreshCcw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
            </Button>
            <Button className="gap-2 shadow-sm" onClick={handleAddClick}>
              <Plus className="w-4 h-4" />
              Create New Plan
            </Button>
          </div>
        </div>

        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <PlanCards onEdit={handleEditClick} />

        <Card className="shadow-xs mt-4">
          <CardHeader>
            <CardTitle>Recent Subscription Invoices</CardTitle>
            <CardDescription>View and manage latest billing activities across all Koperasis.</CardDescription>
          </CardHeader>
          <CardContent>
            <InvoiceTable />
          </CardContent>
        </Card>
      </div>

      <PlanFormModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        planToEdit={planToEdit} 
      />
    </>
  )
}
