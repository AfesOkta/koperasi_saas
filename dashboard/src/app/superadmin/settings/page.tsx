"use client"

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
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Save, Shield, Globe, BellRing } from "lucide-react"

export default function SettingsPage() {
  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="/superadmin">
                Platform Operator
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Platform Settings</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <div className="flex flex-1 flex-col gap-6 p-6 max-w-4xl">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Platform Settings</h2>
          <p className="text-muted-foreground">Configure global parameters and security for the Koperasi SaaS platform.</p>
        </div>

        <div className="grid gap-6">
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Globe className="w-5 h-5 text-primary" />
                <CardTitle>General Configuration</CardTitle>
              </div>
              <CardDescription>Universal settings for all tenants and platform branding.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="platform-name">Platform Name</Label>
                <Input id="platform-name" defaultValue="Koperasi SaaS Platform" />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="support-email">Global Support Email</Label>
                <Input id="support-email" defaultValue="support@koperasisaas.com" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Shield className="w-5 h-5 text-primary" />
                <CardTitle>Security & Access Control</CardTitle>
              </div>
              <CardDescription>Manage global authentication rules and security protocols.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 text-sm">
                <div className="flex items-center justify-between py-2 border-b">
                    <div>
                        <p className="font-medium">Multi-Factor Authentication</p>
                        <p className="text-xs text-muted-foreground">Require MFA for all superadmin accounts.</p>
                    </div>
                    <Button variant="outline" size="sm">Enable</Button>
                </div>
                <div className="flex items-center justify-between py-2">
                    <div>
                        <p className="font-medium">IP Whitelisting</p>
                        <p className="text-xs text-muted-foreground">Restrict access to superadmin dashboard by IP address.</p>
                    </div>
                    <Button variant="outline" size="sm">Configure</Button>
                </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <BellRing className="w-5 h-5 text-primary" />
                <CardTitle>Global Notifications</CardTitle>
              </div>
              <CardDescription>Broadcast messages to all Koperasi admins.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="broadcast-message">Broadcast Message</Label>
                <Input id="broadcast-message" placeholder="Maintenance announcement, new features, etc." />
              </div>
              <Button className="w-fit gap-2"><Save className="w-4 h-4" /> Send Broadcast</Button>
            </CardContent>
          </Card>
        </div>

        <div className="flex justify-end pt-4">
            <Button className="gap-2 px-8">
                <Save className="w-4 h-4" />
                Save Changes
            </Button>
        </div>
      </div>
    </>
  )
}
