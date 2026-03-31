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
import { Badge } from "@/components/ui/badge"
import { Loader2 } from "lucide-react"
import { useSuperadminUserStore, Koperasi } from "../store/useSuperadminUserStore"

interface UserFormModalProps {
  isOpen: boolean
  onClose: () => void
  koperasiToEdit?: Koperasi | null
}

export function UserFormModal({ isOpen, onClose, koperasiToEdit }: UserFormModalProps) {
  const { addKoperasi, updateKoperasi, isLoading } = useSuperadminUserStore()
  
  const [formData, setFormData] = useState({
    name: "",
    address: "",
    email: "",
    phone: "",
    plan: "basic",
    status: "active",
    // Admin setup fields (Create only)
    admin_name: "",
    admin_email: "",
    admin_password: "",
  })

  useEffect(() => {
    if (koperasiToEdit) {
      setFormData({
        name: koperasiToEdit.name || "",
        address: koperasiToEdit.address || "",
        email: koperasiToEdit.email || "",
        phone: koperasiToEdit.phone || "",
        plan: koperasiToEdit.plan || "basic",
        status: koperasiToEdit.status || "active",
        admin_name: "",
        admin_email: "",
        admin_password: "",
      })
    } else {
      setFormData({
        name: "",
        address: "",
        email: "",
        phone: "",
        plan: "basic",
        status: "active",
        admin_name: "",
        admin_email: "",
        admin_password: "",
      })
    }
  }, [koperasiToEdit, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      if (koperasiToEdit) {
        // Edit Mode: Send update payload
        await updateKoperasi(koperasiToEdit.id, {
          name: formData.name,
          email: formData.email,
          phone: formData.phone,
          address: formData.address,
          plan: formData.plan,
          // assuming API might support updating status, otherwise it just ignores
          status: formData.status, 
        })
      } else {
        // Create Mode: Send /onboard payload
        await addKoperasi({
          organization_name: formData.name,
          email: formData.email,
          phone: formData.phone,
          address: formData.address,
          admin_name: formData.admin_name,
          admin_email: formData.admin_email,
          admin_password: formData.admin_password,
        })
      }
      onClose()
    } catch (err: any) {
      // Handled in store/console, but we could add a toast here
      console.error("Submission failed", err)
    }
  }

  const cyclePlan = () => {
    setFormData(prev => {
      const nextPlan: Record<string, string> = {
        'basic': 'premium',
        'premium': 'enterprise',
        'enterprise': 'basic'
      };
      return { ...prev, plan: nextPlan[prev.plan?.toLowerCase()] || 'basic' }
    })
  }

  const getPlanBadgeClasses = (plan: string) => {
    switch(plan?.toLowerCase()) {
      case 'enterprise': return "bg-blue-500 hover:bg-blue-600 text-white border-0"
      case 'premium': case 'business': return "bg-purple-500 hover:bg-purple-600 text-white border-0"
      case 'basic': case 'starter': return "bg-emerald-500 hover:bg-emerald-600 text-white border-0"
      default: return "bg-slate-500 hover:bg-slate-600 text-white border-0"
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[475px] max-h-[90vh] overflow-y-auto">
        <form onSubmit={handleSubmit}>
          <DialogHeader className="mb-4">
            <DialogTitle>{koperasiToEdit ? 'Edit Koperasi Settings' : 'Onboard New Koperasi'}</DialogTitle>
            <DialogDescription>
              {koperasiToEdit 
                ? "Update organization details and subscription plan." 
                : "Register a new Koperasi tenant and setup their admin account."}
            </DialogDescription>
          </DialogHeader>
          
          <div className="grid gap-4 py-2">
            {!koperasiToEdit && (
              <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground w-full border-b pb-1 mb-1">
                Organization Details
              </div>
            )}
            
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

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="email">Org Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({...formData, email: e.target.value})}
                  required
                  placeholder="hello@koperasi.com"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="phone">Phone</Label>
                <Input
                  id="phone"
                  value={formData.phone}
                  onChange={(e) => setFormData({...formData, phone: e.target.value})}
                  placeholder="081234..."
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="address">Address</Label>
              <Input
                id="address"
                value={formData.address}
                onChange={(e) => setFormData({...formData, address: e.target.value})}
                required
                placeholder="e.g. Jakarta, Indonesia"
              />
            </div>

            {koperasiToEdit && (
              <div className="grid gap-2 items-start mt-2">
                <Label>Subscription Plan</Label>
                <p className="text-xs text-muted-foreground mb-1">Click to change target subscription tier</p>
                <div className="cursor-pointer inline-flex items-center" onClick={cyclePlan}>
                  <Badge variant="outline" className={`px-4 py-1.5 cursor-pointer transition-colors ${getPlanBadgeClasses(formData.plan)}`}>
                    {formData.plan ? formData.plan.charAt(0).toUpperCase() + formData.plan.slice(1) : 'Unknown'}
                  </Badge>
                </div>
              </div>
            )}

            {!koperasiToEdit && (
              <>
                <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground w-full border-b pb-1 mt-4 mb-1">
                  Initial Admin Setup
                </div>
                
                <div className="grid gap-2">
                  <Label htmlFor="admin_name">Admin Name</Label>
                  <Input
                    id="admin_name"
                    value={formData.admin_name}
                    onChange={(e) => setFormData({...formData, admin_name: e.target.value})}
                    required={!koperasiToEdit}
                    placeholder="e.g. John Doe"
                  />
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                   <div className="grid gap-2">
                    <Label htmlFor="admin_email">Admin Email (Login)</Label>
                    <Input
                      id="admin_email"
                      type="email"
                      value={formData.admin_email}
                      onChange={(e) => setFormData({...formData, admin_email: e.target.value})}
                      required={!koperasiToEdit}
                      placeholder="admin@koperasi.com"
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="admin_password">Admin Password</Label>
                    <Input
                      id="admin_password"
                      type="password"
                      value={formData.admin_password}
                      onChange={(e) => setFormData({...formData, admin_password: e.target.value})}
                      required={!koperasiToEdit}
                    />
                  </div>
                </div>
              </>
            )}
          </div>
          
          <DialogFooter className="mt-6">
            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
              {koperasiToEdit ? 'Save Changes' : 'Complete Onboarding'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
