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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
  DropdownMenuGroup,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal, Pencil, Mail, MapPin, Loader2 } from "lucide-react"
import { useSuperadminUserStore, Koperasi } from "../store/useSuperadminUserStore"
import { useMemo } from "react"

interface UserTableProps {
  onEdit: (koperasi: Koperasi) => void
}

export function UserTable({ onEdit }: UserTableProps) {
  const { koperasis, searchQuery, isLoading } = useSuperadminUserStore()

  const filteredKoperasis = useMemo(() => {
    return koperasis.filter((koperasi) => {
      const nameMatch = koperasi.name?.toLowerCase().includes(searchQuery.toLowerCase()) || false
      const emailMatch = koperasi.email?.toLowerCase().includes(searchQuery.toLowerCase()) || false
      return nameMatch || emailMatch
    })
  }, [koperasis, searchQuery])

  const getPlanBadgeClasses = (plan: string) => {
    switch(plan?.toLowerCase()) {
      case 'enterprise': return "bg-blue-500 font-medium"
      case 'premium': case 'business': return "bg-purple-500 font-medium"
      case 'basic': case 'starter': return "bg-emerald-500 font-medium"
      default: return "bg-slate-500"
    }
  }

  const getStatusBadgeClasses = (status: string) => {
    switch(status?.toLowerCase()) {
      case 'active': return "bg-emerald-500 font-medium"
      case 'inactive': return "bg-slate-500 font-medium"
      case 'trial': case 'testing': return "bg-amber-500 font-medium"
      default: return "bg-slate-400"
    }
  }

  const formatDate = (dateStr: string) => {
    try {
      if (!dateStr) return '-'
      return new Intl.DateTimeFormat('en-US', {
        month: 'short',
        day: '2-digit',
        year: 'numeric'
      }).format(new Date(dateStr))
    } catch {
      return dateStr
    }
  }

  if (isLoading && koperasis.length === 0) {
    return (
      <div className="flex justify-center items-center h-48">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
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
                  <MapPin className="w-3 h-3" /> {k.address || 'No address provided'}
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
                <Badge className={getPlanBadgeClasses(k.plan)}>
                  {k.plan ? k.plan.charAt(0).toUpperCase() + k.plan.slice(1) : 'Unknown'}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge className={getStatusBadgeClasses(k.status)}>
                  {k.status ? k.status.charAt(0).toUpperCase() + k.status.slice(1) : 'Unknown'}
                </Badge>
              </TableCell>
              <TableCell className="text-sm text-muted-foreground">{formatDate(k.created_at)}</TableCell>
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
                        Edit User
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
