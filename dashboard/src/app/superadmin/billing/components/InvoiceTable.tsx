"use client"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

const invoices = [
  { id: '1', invoiceId: 'INV-2024-001', koperasi: 'Koperasi Maju Jaya', plan: 'Enterprise', amount: 'Rp 5.250.000', status: 'Paid', date: 'Mar 24, 2024' },
  { id: '2', invoiceId: 'INV-2024-002', koperasi: 'Kud Binangun', plan: 'Business', amount: 'Rp 1.490.000', status: 'Paid', date: 'Mar 22, 2024' },
  { id: '3', invoiceId: 'INV-2024-003', koperasi: 'Kopka Mandiri', plan: 'Starter', amount: 'Rp 499.000', status: 'Pending', date: 'Mar 20, 2024' },
]

export function InvoiceTable() {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Invoice ID</TableHead>
          <TableHead>Koperasi</TableHead>
          <TableHead>Plan</TableHead>
          <TableHead>Amount</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Date</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {invoices.map((inv) => (
          <TableRow key={inv.id}>
            <TableCell className="font-mono text-xs">{inv.invoiceId}</TableCell>
            <TableCell className="font-medium">{inv.koperasi}</TableCell>
            <TableCell>{inv.plan}</TableCell>
            <TableCell>{inv.amount}</TableCell>
            <TableCell>
              {inv.status === 'Paid' ? (
                <Badge className="bg-emerald-500">Paid</Badge>
              ) : (
                <Badge variant="outline" className="text-amber-500 border-amber-200 animate-pulse">Pending</Badge>
              )}
            </TableCell>
            <TableCell className="text-right">{inv.date}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
