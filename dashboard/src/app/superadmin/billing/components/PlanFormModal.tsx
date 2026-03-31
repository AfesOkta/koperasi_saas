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
import { Textarea } from "@/components/ui/textarea"
import { Checkbox } from "@/components/ui/checkbox"
import { Loader2 } from "lucide-react"
import { useBillingStore, SubscriptionPlan } from "../store/useBillingStore"

interface PlanFormModalProps {
  isOpen: boolean
  onClose: () => void
  planToEdit?: SubscriptionPlan | null
}

export function PlanFormModal({ isOpen, onClose, planToEdit }: PlanFormModalProps) {
  const { addPlan, updatePlan, isLoading } = useBillingStore()
  
  const [formData, setFormData] = useState({
    name: "",
    code: "",
    price: 0,
    description: "",
    max_users: 10,
    max_members: 100,
    is_popular: false,
  })

  useEffect(() => {
    if (planToEdit) {
      setFormData({
        name: planToEdit.name,
        code: planToEdit.code,
        price: planToEdit.price,
        description: planToEdit.description,
        max_users: planToEdit.max_users,
        max_members: planToEdit.max_members,
        is_popular: planToEdit.is_popular,
      })
    } else {
      setFormData({
        name: "",
        code: "",
        price: 0,
        description: "",
        max_users: 10,
        max_members: 100,
        is_popular: false,
      })
    }
  }, [planToEdit, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      if (planToEdit) {
        await updatePlan(planToEdit.id, formData)
      } else {
        await addPlan(formData)
      }
      onClose()
    } catch (err) {
      // Error is handled in the store
      console.error("Failed to save plan", err)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[475px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>{planToEdit ? 'Edit Subscription Plan' : 'Create New Plan'}</DialogTitle>
            <DialogDescription>
              {planToEdit 
                ? "Modify the features and pricing of this subscription tier." 
                : "Define a new subscription tier for Koperasis."}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4 min-h-[300px] max-h-[60vh] overflow-y-auto px-1">
            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Plan Name</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  required
                  placeholder="e.g. Starter, Pro"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="code">Plan Code</Label>
                <Input
                  id="code"
                  value={formData.code}
                  onChange={(e) => setFormData({...formData, code: e.target.value.toUpperCase()})}
                  required
                  placeholder="e.g. STARTER, PRO"
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="price">Price (IDR)</Label>
              <Input
                id="price"
                type="number"
                value={formData.price}
                onChange={(e) => setFormData({...formData, price: parseInt(e.target.value) || 0})}
                required
                placeholder="e.g. 500000"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="max_users">Max Users</Label>
                <Input
                  id="max_users"
                  type="number"
                  value={formData.max_users}
                  onChange={(e) => setFormData({...formData, max_users: parseInt(e.target.value) || 0})}
                  required
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="max_members">Max Members</Label>
                <Input
                  id="max_members"
                  type="number"
                  value={formData.max_members}
                  onChange={(e) => setFormData({...formData, max_members: parseInt(e.target.value) || 0})}
                  required
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({...formData, description: e.target.value})}
                required
                placeholder="Description of the plan features..."
                className="min-h-[100px]"
              />
            </div>

            <div className="flex items-center space-x-2 pt-2">
              <Checkbox 
                id="is_popular" 
                checked={formData.is_popular}
                onCheckedChange={(checked) => setFormData({...formData, is_popular: !!checked})}
              />
              <Label 
                htmlFor="is_popular" 
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
              >
                Mark as Popular
              </Label>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Plan
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
