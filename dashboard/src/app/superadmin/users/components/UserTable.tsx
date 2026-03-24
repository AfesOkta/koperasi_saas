"use client"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuGroup,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal, Pencil, Trash, Mail, MapPin } from "lucide-react"
import { useSuperadminUserStore, Koperasi } from "../store/useSuperadminUserStore"
import { useMemo } from "react"

interface UserTableProps {
  onEdit: (koperasi: Koperasi) => void
}

export function UserTable({ onEdit }: UserTableProps) {
  const { koperasis, searchQuery, deleteKoperasi } = useSuperadminUserStore()

  const filteredKoperasis = useMemo(() => {
    return koperasis.filter((koperasi) => {
      return koperasi.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
             koperasi.email.toLowerCase().includes(searchQuery.toLowerCase())
    })
  }, [koperasis, searchQuery])

  const getPlanBadgeClasses = (plan: string) => {
    switch(plan) {
      case 'Enterprise': return "bg-blue-500 font-medium"
      case 'Business': return "bg-purple-500 font-medium"
      case 'Starter': return "bg-emerald-500 font-medium"
      default: return ""
    }
  }

  const getStatusBadgeClasses = (status: string) => {
    switch(status) {
      case 'active': return "bg-emerald-500 font-medium"
      case 'inactive': return "bg-slate-500 font-medium"
      case 'trial': return "bg-amber-500 font-medium"
      default: return ""
    }
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Koperasi Name</TableHead>
          <TableHead>Contact info</TableHead>
          <TableHead>Plan</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Joined Date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {filteredKoperasis.length === 0 ? (
          <TableRow>
            <TableCell colSpan={6} className="h-24 text-center">
              No koperasis found.
            </TableCell>
          </TableRow>
        ) : (
          filteredKoperasis.map((k) => (
            <TableRow key={k.id} className="group cursor-default hover:bg-muted/50 transition-colors">
              <TableCell>
                <div className="font-semibold text-lg tracking-tight">{k.name}</div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground mt-0.5">
                  <MapPin className="w-3 h-3" /> {k.location}
                </div>
              </TableCell>
              <TableCell>
                <div className="flex flex-col gap-1 text-sm">
                  <div className="flex items-center gap-1.5 text-muted-foreground">
                    <Mail className="w-3.5 h-3.5" /> {k.email}
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <Badge className={getPlanBadgeClasses(k.plan)}>{k.plan}</Badge>
              </TableCell>
              <TableCell>
                <Badge className={getStatusBadgeClasses(k.status)}>
                  {k.status.charAt(0).toUpperCase() + k.status.slice(1)}
                </Badge>
              </TableCell>
              <TableCell className="text-sm text-muted-foreground">{k.joinDate}</TableCell>
              <TableCell className="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground h-8 w-8 p-0 opacity-0 group-hover:opacity-100 data-[state=open]:opacity-100">
                    <span className="sr-only">Open menu</span>
                    <MoreHorizontal className="h-4 w-4" />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuGroup>
                      <DropdownMenuLabel>Actions</DropdownMenuLabel>
                      <DropdownMenuItem onClick={() => onEdit(k)} className="cursor-pointer">
                        <Pencil className="mr-2 h-4 w-4" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem 
                        className="text-red-600 focus:text-red-600 cursor-pointer"
                        onClick={() => deleteKoperasi(k.id)}
                      >
                        <Trash className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuGroup>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  )
}
