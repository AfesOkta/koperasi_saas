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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"

export default function SuperadminAccountPage() {
  const user = {
    name: "Superadmin",
    email: "superadmin@koperasisaas.com",
    role: "Platform Operator",
    joined: "System Initialization",
    avatar: ""
  }

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
              <BreadcrumbPage>Account Profile</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6 max-w-4xl">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Operator Profile</h2>
          <p className="text-muted-foreground mt-1">Superadmin platform details and security settings.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>System Identity</CardTitle>
            <CardDescription>Highest level administrative access to the SaaS Platform.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col md:flex-row gap-8">
            <div className="flex flex-col items-center gap-4">
              <Avatar className="h-24 w-24 rounded-full ring-4 ring-primary/20">
                <AvatarImage src={user.avatar} alt={user.name} />
                <AvatarFallback className="bg-primary/10 text-primary text-2xl font-semibold tracking-wider">
                  SU
                </AvatarFallback>
              </Avatar>
              <Badge className="bg-indigo-500 hover:bg-indigo-600">Superuser</Badge>
            </div>
            
            <div className="flex-1 grid gap-4 grid-cols-1 md:grid-cols-2">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Account Name</p>
                <p className="font-medium mt-1">{user.name}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Admin Email</p>
                <p className="font-medium mt-1 tracking-tight">{user.email}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Privilege Level</p>
                <p className="font-medium mt-1">{user.role}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Activated Since</p>
                <p className="font-medium mt-1 text-muted-foreground">{user.joined}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  )
}
