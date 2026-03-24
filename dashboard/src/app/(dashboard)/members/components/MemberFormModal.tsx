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
import { useMemberStore, Member, MemberStatus } from "../store/useMemberStore"
import { Badge } from "@/components/ui/badge"

interface MemberFormModalProps {
  isOpen: boolean
  onClose: () => void
  memberToEdit?: Member | null
}

export function MemberFormModal({ isOpen, onClose, memberToEdit }: MemberFormModalProps) {
  const { addMember, updateMember } = useMemberStore()
  
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    phone: "",
    balance: 0,
    status: "active" as MemberStatus,
  })

  useEffect(() => {
    if (memberToEdit) {
      setFormData({
        name: memberToEdit.name,
        email: memberToEdit.email,
        phone: memberToEdit.phone,
        balance: memberToEdit.balance,
        status: memberToEdit.status,
      })
    } else {
      setFormData({
        name: "",
        email: "",
        phone: "",
        balance: 0,
        status: "active",
      })
    }
  }, [memberToEdit, isOpen])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (memberToEdit) {
      updateMember(memberToEdit.id, formData)
    } else {
      addMember(formData)
    }
    
    onClose()
  }

  const toggleStatus = () => {
    setFormData(prev => ({
      ...prev,
      status: prev.status === 'active' ? 'inactive' : 'active'
    }))
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>{memberToEdit ? 'Edit Member' : 'Add Member'}</DialogTitle>
            <DialogDescription>
              {memberToEdit 
                ? "Make changes to the member's profile here. Click save when you're done." 
                : "Add a new member to the cooperative system."}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Full Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({...formData, name: e.target.value})}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="email">Email address</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({...formData, email: e.target.value})}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="phone">Phone Number</Label>
              <Input
                id="phone"
                value={formData.phone}
                onChange={(e) => setFormData({...formData, phone: e.target.value})}
                required
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="balance">Initial Balance (IDR)</Label>
                <Input
                  id="balance"
                  type="number"
                  min="0"
                  value={formData.balance}
                  onChange={(e) => setFormData({...formData, balance: Number(e.target.value)})}
                  required
                />
              </div>
              <div className="grid gap-2 items-start">
                <Label>Status</Label>
                <div 
                  className="cursor-pointer inline-flex items-center mt-2" 
                  onClick={toggleStatus}
                >
                  <Badge variant={formData.status === 'active' ? 'default' : 'secondary'} className="px-4 py-1.5 cursor-pointer">
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
