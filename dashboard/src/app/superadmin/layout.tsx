import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { SuperadminSidebar } from "@/components/superadmin/sidebar";
import { TooltipProvider } from "@/components/ui/tooltip";
import { DevToolsHider } from "@/components/dev-tools-hider";

export default function SuperadminLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <TooltipProvider>
      <SidebarProvider>
        <SuperadminSidebar />
        <SidebarInset>
          <main className="w-full">
            <DevToolsHider />
            {children}
          </main>
        </SidebarInset>
      </SidebarProvider>
    </TooltipProvider>
  );
}
