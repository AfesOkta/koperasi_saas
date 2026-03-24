"use client"

import { useState } from "react"
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
import { Input } from "@/components/ui/input"
import { Search, UserPlus } from "lucide-react"
import { UserTable } from "./components/UserTable"
import { UserFormModal } from "./components/UserFormModal"
import { useSuperadminUserStore, Koperasi } from "./store/useSuperadminUserStore"

export default function UsersPage() {
  const { koperasis, searchQuery, setSearchQuery } = useSuperadminUserStore()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [koperasiToEdit, setKoperasiToEdit] = useState<Koperasi | null>(null)

  const handleAddClick = () => {
    setKoperasiToEdit(null)
    setIsModalOpen(true)
  }

  const handleEditClick = (koperasi: Koperasi) => {
    setKoperasiToEdit(koperasi)
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
              <BreadcrumbPage>Koperasi Users</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Koperasi Management</h2>
            <p className="text-muted-foreground mt-1">Monitor and manage all Koperasis on the platform.</p>
          </div>
          <Button className="w-full sm:w-auto shadow-sm gap-2" onClick={handleAddClick}>
            <UserPlus className="w-4 h-4" />
            Add New Koperasi
          </Button>
        </div>

        <Card className="shadow-xs">
          <CardHeader className="pb-3">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
              <div>
                <CardTitle>Registered Koperasis</CardTitle>
                <CardDescription>A total of {koperasis.length} Koperasis are registered on the platform.</CardDescription>
              </div>
              <div className="flex items-center gap-2 w-full md:w-auto">
                <div className="relative w-full md:w-64">
                  <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="search"
                    placeholder="Search Koperasi..."
                    className="pl-8 h-9 shadow-sm"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <UserTable onEdit={handleEditClick} />
          </CardContent>
        </Card>
      </div>

      <UserFormModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        koperasiToEdit={koperasiToEdit} 
      />
    </>
  )
}
