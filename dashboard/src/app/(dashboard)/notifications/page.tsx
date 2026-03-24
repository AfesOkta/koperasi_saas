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
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Bell, CreditCard, ShieldAlert } from "lucide-react"

export default function NotificationsPage() {
  const notifications = [
    {
      id: 1,
      title: "Subscription Renewed",
      description: "Your Business Plan subscription has been successfully renewed for another month.",
      date: "2 hours ago",
      icon: <CreditCard className="text-emerald-500 w-5 h-5" />,
      read: false
    },
    {
      id: 2,
      title: "New Feature Available",
      description: "Detailed financial reporting is now available in your dashboard.",
      date: "1 day ago",
      icon: <Bell className="text-blue-500 w-5 h-5" />,
      read: true
    },
    {
      id: 3,
      title: "Security Alert",
      description: "A new login was detected from a new IP address.",
      date: "3 days ago",
      icon: <ShieldAlert className="text-amber-500 w-5 h-5" />,
      read: true
    }
  ]

  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b bg-background/95 backdrop-blur-sm sticky top-0 z-10 px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="/" className="font-medium">
                Koperasi SaaS
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Notifications</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6 max-w-4xl">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Your Notifications</h2>
          <p className="text-muted-foreground mt-1">Updates and alerts regarding your Koperasi account.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>You have 1 unread notification.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            {notifications.map((notification) => (
              <div 
                key={notification.id} 
                className={`flex items-start gap-4 p-4 rounded-lg border transition-colors ${notification.read ? 'bg-background' : 'bg-muted/50 border-primary/20'}`}
              >
                <div className="mt-0.5 shrink-0 bg-background p-2 rounded-full shadow-sm border">
                  {notification.icon}
                </div>
                <div className="flex-1 space-y-1">
                  <div className="flex items-center justify-between">
                    <p className={`font-medium ${!notification.read && 'text-primary'}`}>{notification.title}</p>
                    <span className="text-xs text-muted-foreground">{notification.date}</span>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    {notification.description}
                  </p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </>
  )
}
