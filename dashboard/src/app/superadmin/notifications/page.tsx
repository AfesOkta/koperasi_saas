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
import { Bell, ShieldAlert, Rocket, Users } from "lucide-react"

export default function SuperadminNotificationsPage() {
  const notifications = [
    {
      id: 1,
      title: "New Koperasi Registered",
      description: "Koperasi Sejahtera Bersama has just registered for a Trial plan.",
      date: "10 minutes ago",
      icon: <Users className="text-emerald-500 w-5 h-5" />,
      read: false,
      systemWide: true
    },
    {
      id: 2,
      title: "Plan Upgrade Request",
      description: "Kud Binangun is requesting an upgrade to the Enterprise custom tier.",
      date: "1 hour ago",
      icon: <Rocket className="text-purple-500 w-5 h-5" />,
      read: false,
      systemWide: true
    },
    {
      id: 3,
      title: "System Maintenance Completed",
      description: "The database optimization script has finished running successfully.",
      date: "5 hours ago",
      icon: <Bell className="text-blue-500 w-5 h-5" />,
      read: true,
      systemWide: true
    },
    {
      id: 4,
      title: "Failed Invoice Payment",
      description: "Payment attempt for Koperasi Mandiri (Starter Plan) failed. Automatic retry in 24h.",
      date: "1 day ago",
      icon: <ShieldAlert className="text-amber-500 w-5 h-5" />,
      read: true,
      systemWide: true
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
              <BreadcrumbLink href="/superadmin" className="font-medium">
                Platform Operator
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Notifications & Logs</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6 max-w-4xl">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">System Notifications</h2>
          <p className="text-muted-foreground mt-1">Platform-wide events, billing alerts, and system logs.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Global Event Stream</CardTitle>
            <CardDescription>Viewing all system-wide notifications. You have 2 unread alerts.</CardDescription>
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
