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
import { Button } from "@/components/ui/button"

export default function AccountPage() {
  const user = {
    name: "Admin",
    email: "admin@koperasisaas.com",
    role: "Koperasi Administrator",
    joined: "January 2024",
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
              <BreadcrumbLink href="/" className="font-medium">
                Koperasi SaaS
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Account Overview</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6 max-w-4xl">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Account Profile</h2>
          <p className="text-muted-foreground mt-1">Manage your account settings and preferences.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
            <CardDescription>Your personal identity and contact details.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col md:flex-row gap-8">
            <div className="flex flex-col items-center gap-4">
              <Avatar className="h-24 w-24 rounded-full">
                <AvatarImage src={user.avatar} alt={user.name} />
                <AvatarFallback className="bg-primary/5 text-primary text-2xl font-semibold tracking-wider">
                  {user.name.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <Button variant="outline" size="sm">Change Avatar</Button>
            </div>
            
            <div className="flex-1 grid gap-4 grid-cols-1 md:grid-cols-2">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Full Name</p>
                <p className="font-medium mt-1">{user.name}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Email Address</p>
                <p className="font-medium mt-1">{user.email}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Role</p>
                <p className="font-medium mt-1">{user.role}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Joined At</p>
                <p className="font-medium mt-1">{user.joined}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  )
}
