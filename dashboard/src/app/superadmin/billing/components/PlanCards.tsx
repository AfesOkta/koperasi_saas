"use client"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Edit2, ShieldCheck, Zap, Rocket, Users, Trash2, Loader2, CreditCard } from "lucide-react"
import { useBillingStore, SubscriptionPlan } from "../store/useBillingStore"

interface PlanCardsProps {
  onEdit?: (plan: SubscriptionPlan) => void
  isSuperadmin?: boolean
}

export function PlanCards({ onEdit, isSuperadmin = true }: PlanCardsProps) {
  const { plans, deletePlan, isLoading, error } = useBillingStore()

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price)
  }

  const getThemeStyles = (theme?: string) => {
    switch(theme) {
      case 'purple': return {
        badge: "text-purple-500 border-purple-200 bg-purple-50",
        icon: <Rocket className="w-4 h-4 text-purple-500" />,
        cardClass: "border-primary/20 shadow-lg shadow-primary/5"
      }
      case 'amber': return {
        badge: "text-amber-500 border-amber-200 bg-amber-50",
        icon: <ShieldCheck className="w-4 h-4 text-amber-500" />,
        cardClass: ""
      }
      case 'blue':
      default: return {
        badge: "text-blue-500 border-blue-200 bg-blue-50",
        icon: <Zap className="w-4 h-4 text-blue-500" />,
        cardClass: ""
      }
    }
  }

  if (isLoading && plans.length === 0) {
    return (
      <div className="flex items-center justify-center h-48">
        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error && plans.length === 0) {
    return (
      <div className="text-center p-8 border border-dashed rounded-lg">
        <p className="text-destructive font-medium">{error}</p>
        <Button variant="link" onClick={() => useBillingStore.getState().fetchPlans()}>Try again</Button>
      </div>
    )
  }

  return (
    <div className="grid gap-6 md:grid-cols-3">
      {plans.map((plan) => {
        const theme = getThemeStyles(plan.colorTheme)
        return (
          <Card key={plan.id} className={`relative overflow-hidden ${theme.cardClass}`}>
            {plan.is_popular && (
              <div className="absolute top-0 right-0 p-2">
                <Badge className="bg-primary text-primary-foreground">Popular</Badge>
              </div>
            )}
            <CardHeader>
              <div className="flex items-center justify-between">
                <Badge variant="outline" className={theme.badge}>{plan.name}</Badge>
                {theme.icon}
              </div>
              <CardTitle className="text-2xl mt-4">
                {formatPrice(plan.price)}
                <span className="text-sm font-normal text-muted-foreground ml-1">/mo</span>
              </CardTitle>
              <CardDescription className="line-clamp-2 min-h-[40px]">{plan.description}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-6">
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Badge variant="secondary" className="font-mono text-[10px] h-5">{plan.code}</Badge>
                  <span className="text-xs italic">Identifier Code</span>
                </div>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="flex items-center gap-2 px-3 py-2 bg-muted/50 rounded-md">
                    <Users className="w-4 h-4 text-blue-500" />
                    <span className="font-medium">{plan.max_users} Users</span>
                  </div>
                  <div className="flex items-center gap-2 px-3 py-2 bg-muted/50 rounded-md">
                    <CreditCard className="w-4 h-4 text-emerald-500" />
                    <span className="font-medium">{plan.max_members} Members</span>
                  </div>
                </div>
              </div>
              
              <div className="flex gap-2">
                {isSuperadmin ? (
                  <>
                    <Button 
                      variant={plan.is_popular ? "default" : "outline"} 
                      className="flex-1 gap-2"
                      onClick={() => onEdit?.(plan)}
                    >
                      <Edit2 className="w-4 h-4" /> Edit
                    </Button>
                    <Button 
                      variant="outline" 
                      className="px-3 text-destructive hover:bg-destructive hover:text-destructive-foreground active:scale-95 transition-all"
                      onClick={() => {
                        if (confirm(`Are you sure you want to delete the ${plan.name} plan?`)) {
                          deletePlan(plan.id)
                        }
                      }}
                      disabled={isLoading}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </>
                ) : (
                  <Button 
                    variant={plan.is_popular ? "default" : "outline"} 
                    className="w-full"
                  >
                    Subscribe Now
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
