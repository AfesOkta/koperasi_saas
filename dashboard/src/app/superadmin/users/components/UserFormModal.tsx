"use client"

import { useEffect, useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { useSuperadminUserStore, Koperasi, KoperasiPlan, KoperasiStatus } from "../store/useSuperadminUserStore"
import { Badge } from "@/components/ui/badge"

interface UserFormModalProps {
  isOpen: boolean
  onClose: () => void
  koperasiToEdit?: Koperasi | null
}

export function UserFormModal({ isOpen, onClose, koperasiToEdit }: UserFormModalProps) {
  const { addKoperasi, updateKoperasi } = useSuperadminUserStore()
  
  const [formData, setFormData] = useState({
    name: "",
    location: "",
    email: "",
    plan: "Starter" as KoperasiPlan,
    status: "trial" as KoperasiStatus,
  })

  useEffect(() => {
    if (koperasiToEdit) {
      setFormData({
        name: koperasiToEdit.name,
        location: koperasiToEdit.location,
        email: koperasiToEdit.email,
        plan: koperasiToEdit.plan,
        status: koperasiToEdit.status,
      })
    } else {
      setFormData({
        name: "",
        location: "",
        email: "",
        plan: "Starter",
        status: "trial",
      })
    }
  }, [koperasiToEdit, isOpen])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (koperasiToEdit) {
      updateKoperasi(koperasiToEdit.id, formData)
    } else {
      addKoperasi(formData)
    }
    
    onClose()
  }

  const cycleStatus = () => {
    setFormData(prev => {
      const nextStatus: Record<KoperasiStatus, KoperasiStatus> = {
        'active': 'inactive',
        'inactive': 'trial',
        'trial': 'active'
      };
      return { ...prev, status: nextStatus[prev.status] }
    })
  }

  const cyclePlan = () => {
    setFormData(prev => {
      const nextPlan: Record<KoperasiPlan, KoperasiPlan> = {
        'Starter': 'Business',
        'Business': 'Enterprise',
        'Enterprise': 'Starter'
      };
      return { ...prev, plan: nextPlan[prev.plan] }
    })
  }

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
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>{koperasiToEdit ? 'Edit Koperasi' : 'Add New Koperasi'}</DialogTitle>
            <DialogDescription>
              {koperasiToEdit 
                ? "Update the configuration for this Koperasi tenant." 
                : "Register a new Koperasi tenant on the platform."}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Koperasi Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({...formData, name: e.target.value})}
                required
                placeholder="e.g. Koperasi Maju Jaya"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="location">Location</Label>
              <Input
                id="location"
                value={formData.location}
                onChange={(e) => setFormData({...formData, location: e.target.value})}
                required
                placeholder="e.g. Jakarta, Indonesia"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="email">Admin Email Address</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({...formData, email: e.target.value})}
                required
                placeholder="admin@koperasi.com"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2 items-start">
                <Label>Subscription Plan</Label>
                <div 
                  className="cursor-pointer inline-flex items-center mt-2" 
                  onClick={cyclePlan}
                >
                  <Badge className={`px-4 py-1.5 cursor-pointer hover:opacity-80 transition-opacity ${getPlanBadgeClasses(formData.plan)}`}>
                    {formData.plan}
                  </Badge>
                </div>
              </div>
              <div className="grid gap-2 items-start">
                <Label>Account Status</Label>
                <div 
                  className="cursor-pointer inline-flex items-center mt-2" 
                  onClick={cycleStatus}
                >
                  <Badge className={`px-4 py-1.5 cursor-pointer hover:opacity-80 transition-opacity ${getStatusBadgeClasses(formData.status)}`}>
                    {formData.status.charAt(0).toUpperCase() + formData.status.slice(1)}
                  </Badge>
                </div>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
