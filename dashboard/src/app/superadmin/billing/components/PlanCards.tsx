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
import { Edit2, ShieldCheck, Zap, Rocket } from "lucide-react"
import { useBillingStore, SubscriptionPlan } from "../store/useBillingStore"

interface PlanCardsProps {
  onEdit?: (plan: SubscriptionPlan) => void
  isSuperadmin?: boolean
}

export function PlanCards({ onEdit, isSuperadmin = true }: PlanCardsProps) {
  const { plans } = useBillingStore()

  const getThemeStyles = (theme: string) => {
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

  return (
    <div className="grid gap-6 md:grid-cols-3">
      {plans.map((plan) => {
        const theme = getThemeStyles(plan.colorTheme)
        return (
          <Card key={plan.id} className={`relative overflow-hidden ${theme.cardClass}`}>
            {plan.isPopular && (
              <div className="absolute top-0 right-0 p-2">
                <Badge className="bg-primary text-primary-foreground">Popular</Badge>
              </div>
            )}
            <CardHeader>
              <div className="flex items-center justify-between">
                <Badge variant="outline" className={theme.badge}>{plan.name}</Badge>
                {theme.icon}
              </div>
              <CardTitle className="text-2xl">
                {plan.priceDisplay}
                {plan.interval && <span className="text-sm font-normal text-muted-foreground">{plan.interval}</span>}
              </CardTitle>
              <CardDescription>{plan.description}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
              <div className="text-sm text-muted-foreground min-h-[80px]">
                <ul className="grid gap-2">
                  {plan.features.map((feature) => (
                    <li key={feature.id} className="flex items-center gap-2">
                      <ShieldCheck className="w-4 h-4 text-emerald-500 shrink-0" /> 
                      {feature.name}
                    </li>
                  ))}
                </ul>
              </div>
                {isSuperadmin ? (
                  <Button 
                    variant={plan.isPopular ? "default" : "outline"} 
                    className="w-full gap-2"
                    onClick={() => onEdit?.(plan)}
                  >
                    <Edit2 className="w-4 h-4" /> Edit Plan
                  </Button>
                ) : (
                  <Button 
                    variant={plan.isPopular ? "default" : "outline"} 
                    className="w-full"
                  >
                    Update Planning
                  </Button>
                )}
              </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
