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
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Search, Plus, Filter } from "lucide-react"
import { useMemberStore, Member, MemberStatus } from "./store/useMemberStore"
import { MemberTable } from "./components/MemberTable"
import { MemberFormModal } from "./components/MemberFormModal"

export default function MembersPage() {
  const { searchQuery, setSearchQuery, statusFilter, setStatusFilter } = useMemberStore()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [memberToEdit, setMemberToEdit] = useState<Member | null>(null)

  const handleAddClick = () => {
    setMemberToEdit(null)
    setIsModalOpen(true)
  }

  const handleEditClick = (member: Member) => {
    setMemberToEdit(member)
    setIsModalOpen(true)
  }

  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="#">
                Koperasi SaaS
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Members</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      
      <div className="flex flex-1 flex-col gap-4 p-4 md:p-6 pt-6 bg-muted/20">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Member Management</h1>
            <p className="text-sm text-muted-foreground mt-1">Manage your cooperative members and their profiles.</p>
          </div>
          <Button onClick={handleAddClick} className="w-full sm:w-auto shadow-sm">
            <Plus className="mr-2 h-4 w-4" /> Add Member
          </Button>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 items-center justify-between mt-4">
          <div className="relative w-full sm:max-w-xs shadow-sm rounded-md">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search members..."
              className="w-full pl-8 bg-background"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <div className="flex items-center gap-2 w-full sm:w-auto shadow-sm rounded-md">
            <div className="flex items-center bg-background border border-input rounded-md px-3 h-10 w-full sm:w-auto overflow-hidden text-sm ring-offset-background">
              <Filter className="h-4 w-4 text-muted-foreground mr-2 shrink-0" />
              <select
                className="w-full sm:w-[130px] bg-transparent outline-none focus:outline-none focus:ring-0 cursor-pointer text-sm"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value as 'all' | MemberStatus)}
              >
                <option value="all">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
              </select>
            </div>
          </div>
        </div>

        <div className="mt-2 bg-background rounded-lg shadow-sm border">
          <MemberTable onEdit={handleEditClick} />
        </div>
      </div>

      <MemberFormModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        memberToEdit={memberToEdit} 
      />
    </>
  )
}
