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
import { useBillingStore, SubscriptionPlan } from "../store/useBillingStore"

interface PlanFormModalProps {
  isOpen: boolean
  onClose: () => void
  planToEdit?: SubscriptionPlan | null
}

export function PlanFormModal({ isOpen, onClose, planToEdit }: PlanFormModalProps) {
  const { addPlan, updatePlan } = useBillingStore()
  
  const [formData, setFormData] = useState({
    name: "",
    priceAmount: 0,
    priceDisplay: "",
    interval: "/mo",
    description: "",
    featuresText: "",
    isPopular: false,
    colorTheme: "blue" as "blue" | "purple" | "amber",
  })

  useEffect(() => {
    if (planToEdit) {
      setFormData({
        name: planToEdit.name,
        priceAmount: planToEdit.priceAmount,
        priceDisplay: planToEdit.priceDisplay,
        interval: planToEdit.interval,
        description: planToEdit.description,
        featuresText: planToEdit.features.map(f => f.name).join('\n'),
        isPopular: planToEdit.isPopular,
        colorTheme: planToEdit.colorTheme,
      })
    } else {
      setFormData({
        name: "",
        priceAmount: 0,
        priceDisplay: "",
        interval: "/mo",
        description: "",
        featuresText: "",
        isPopular: false,
        colorTheme: "blue",
      })
    }
  }, [planToEdit, isOpen])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    // Parse features from text
    const features = formData.featuresText
      .split('\n')
      .map(line => line.trim())
      .filter(line => line.length > 0)
      .map((name, i) => ({ id: `new_${Date.now()}_${i}`, name }))

    const submittedData = {
      name: formData.name,
      priceAmount: formData.priceAmount,
      priceDisplay: formData.priceDisplay,
      interval: formData.interval,
      description: formData.description,
      isPopular: formData.isPopular,
      colorTheme: formData.colorTheme,
      features,
    }

    if (planToEdit) {
      updatePlan(planToEdit.id, submittedData)
    } else {
      addPlan(submittedData)
    }
    
    onClose()
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
            <div className="grid gap-2">
              <Label htmlFor="name">Plan Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({...formData, name: e.target.value})}
                required
                placeholder="e.g. Starter, Pro, Custom"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="priceDisplay">Display Price</Label>
                <Input
                  id="priceDisplay"
                  value={formData.priceDisplay}
                  onChange={(e) => setFormData({...formData, priceDisplay: e.target.value})}
                  required
                  placeholder="e.g. Rp 499k"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="interval">Interval</Label>
                <Input
                  id="interval"
                  value={formData.interval}
                  onChange={(e) => setFormData({...formData, interval: e.target.value})}
                  placeholder="e.g. /mo, /yr, or empty"
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Short Description</Label>
              <Input
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({...formData, description: e.target.value})}
                required
                placeholder="Up to 100 members..."
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="features">Features (One per line)</Label>
              <Textarea
                id="features"
                value={formData.featuresText}
                onChange={(e) => setFormData({...formData, featuresText: e.target.value})}
                required
                placeholder="Core Banking Module&#10;Member Management&#10;24/7 Support"
                className="min-h-[100px]"
              />
            </div>

            <div className="grid grid-cols-2 gap-4 mt-2">
              <div className="grid gap-2">
                <Label>Theme Color</Label>
                <select
                  className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  value={formData.colorTheme}
                  onChange={(e) => setFormData({...formData, colorTheme: e.target.value as any})}
                >
                  <option value="blue">Blue (Starter)</option>
                  <option value="purple">Purple (Business)</option>
                  <option value="amber">Amber (Enterprise)</option>
                </select>
              </div>

              <div className="flex items-center space-x-2 pt-8">
                <Checkbox 
                  id="isPopular" 
                  checked={formData.isPopular}
                  onCheckedChange={(c) => setFormData({...formData, isPopular: !!c})}
                />
                <Label htmlFor="isPopular" className="cursor-pointer">Mark as Popular</Label>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">Save Plan</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
