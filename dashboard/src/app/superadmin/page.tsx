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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts"
import { TrendingUp, CreditCard, Building2, CheckCircle2, Clock } from "lucide-react"
import { Badge } from "@/components/ui/badge"

const chartData = [
  { month: "Jan", revenue: 4500, users: 12 },
  { month: "Feb", revenue: 5200, users: 15 },
  { month: "Mar", revenue: 4800, users: 18 },
  { month: "Apr", revenue: 6100, users: 22 },
  { month: "May", revenue: 5900, users: 25 },
  { month: "Jun", revenue: 7200, users: 30 },
]

const chartConfig = {
  revenue: {
    label: "Revenue",
    color: "var(--chart-1)",
  },
  users: {
    label: "Koperasis",
    color: "var(--chart-2)",
  },
} satisfies ChartConfig

export default function SuperadminPage() {
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
              <BreadcrumbPage>Dashboard Overview</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      
      <div className="flex flex-1 flex-col gap-6 p-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Platform Overview</h1>
          <p className="text-muted-foreground mt-1">Monitor all Koperasi tenants and platform revenue.</p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card className="border-indigo-100 dark:border-indigo-900/50 shadow-xs">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-muted-foreground">Total Koperasis</CardTitle>
              <div className="p-2 bg-indigo-50 dark:bg-indigo-500/10 rounded-full">
                <Building2 className="w-4 h-4 text-indigo-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-foreground">42</div>
              <p className="text-xs font-medium text-emerald-600 dark:text-emerald-400 mt-1 flex items-center">
                <TrendingUp className="w-3 h-3 mr-1" />
                +3 from last month
              </p>
            </CardContent>
          </Card>

          <Card className="border-emerald-100 dark:border-emerald-900/50 shadow-xs">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-muted-foreground">Monthly Revenue</CardTitle>
              <div className="p-2 bg-emerald-50 dark:bg-emerald-500/10 rounded-full">
                <TrendingUp className="w-4 h-4 text-emerald-600" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-foreground">Rp 127.5M</div>
              <p className="text-xs font-medium text-emerald-600 dark:text-emerald-400 mt-1 flex items-center">
                <TrendingUp className="w-3 h-3 mr-1" />
                +15.2% from last month
              </p>
            </CardContent>
          </Card>

          <Card className="border-amber-100 dark:border-amber-900/50 shadow-xs">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-muted-foreground">Active Subscriptions</CardTitle>
              <div className="p-2 bg-amber-50 dark:bg-amber-500/10 rounded-full">
                <CreditCard className="w-4 h-4 text-amber-600" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-foreground">38</div>
              <p className="text-xs font-medium text-muted-foreground mt-1 flex items-center">
                <Clock className="w-3 h-3 mr-1" />
                4 in trial period
              </p>
            </CardContent>
          </Card>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
          <Card className="lg:col-span-4 shadow-xs">
            <CardHeader>
              <CardTitle>Platform Growth</CardTitle>
              <CardDescription>Revenue and User acquisition over time</CardDescription>
            </CardHeader>
            <CardContent className="pl-0 pb-0">
              <ChartContainer config={chartConfig} className="h-[320px] w-full mt-2">
                <BarChart accessibilityLayer data={chartData} margin={{ left: 12, right: 12, bottom: 0, top: 0 }}>
                  <CartesianGrid vertical={false} strokeDasharray="3 3" className="stroke-muted" />
                  <XAxis
                    dataKey="month"
                    tickLine={false}
                    tickMargin={10}
                    axisLine={false}
                    className="text-xs font-medium fill-muted-foreground"
                  />
                  <YAxis
                    tickLine={false}
                    axisLine={false}
                    tickMargin={10}
                    className="text-xs font-medium fill-muted-foreground"
                    hide
                  />
                  <ChartTooltip cursor={{ fill: 'var(--color-muted)' }} content={<ChartTooltipContent />} />
                  <Bar dataKey="revenue" fill="var(--color-revenue)" radius={[4, 4, 0, 0]} maxBarSize={40} />
                  <Bar dataKey="users" fill="var(--color-users)" radius={[4, 4, 0, 0]} maxBarSize={40} />
                </BarChart>
              </ChartContainer>
            </CardContent>
          </Card>

          <Card className="lg:col-span-3 shadow-xs">
            <CardHeader>
              <CardTitle>Recent Activations</CardTitle>
              <CardDescription>Latest Koperasi registrations</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="font-medium text-muted-foreground">Koperasi</TableHead>
                    <TableHead className="font-medium text-muted-foreground hidden sm:table-cell">Plan</TableHead>
                    <TableHead className="text-right font-medium text-muted-foreground">Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow className="bg-transparent hover:bg-muted/50 transition-colors">
                    <TableCell>
                      <div className="font-semibold text-foreground">Koperasi Maju Jaya</div>
                      <div className="text-xs text-muted-foreground mt-0.5">Jakarta</div>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className="font-medium text-indigo-600 dark:text-indigo-400 border-indigo-200 dark:border-indigo-900 bg-indigo-50 dark:bg-indigo-900/20">Enterprise</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <span className="inline-flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400 font-medium text-sm">
                        <CheckCircle2 className="w-3.5 h-3.5" />
                        Active
                      </span>
                    </TableCell>
                  </TableRow>
                  <TableRow className="bg-transparent hover:bg-muted/50 transition-colors">
                    <TableCell>
                      <div className="font-semibold text-foreground">Kud Binangun</div>
                      <div className="text-xs text-muted-foreground mt-0.5">Yogyakarta</div>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className="font-medium text-purple-600 dark:text-purple-400 border-purple-200 dark:border-purple-900 bg-purple-50 dark:bg-purple-900/20">Business</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <span className="inline-flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400 font-medium text-sm">
                        <CheckCircle2 className="w-3.5 h-3.5" />
                        Active
                      </span>
                    </TableCell>
                  </TableRow>
                  <TableRow className="bg-transparent hover:bg-muted/50 transition-colors border-b-0">
                    <TableCell>
                      <div className="font-semibold text-foreground">Kopka Mandiri</div>
                      <div className="text-xs text-muted-foreground mt-0.5">Bandung</div>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className="font-medium text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/20">Starter</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <span className="inline-flex items-center gap-1.5 text-amber-600 dark:text-amber-500 font-medium text-sm">
                        <Clock className="w-3.5 h-3.5" />
                        Trial
                      </span>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  )
}
