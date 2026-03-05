import { useTable } from "@refinedev/react-table";
import { useNavigation } from "@refinedev/core";
import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Eye } from "lucide-react";

interface IUser {
  id: string;
  email: string;
}

export function UserList() {
  const { show } = useNavigation();

  const columns: ColumnDef<IUser>[] = [
    { accessorKey: "email", header: "Email" },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => show("users", row.original.id)}>
          <Eye className="h-4 w-4" />
        </Button>
      ),
    },
  ];

  const { reactTable: table } = useTable<IUser>({
    columns,
    refineCoreProps: { resource: "users" },
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div>
      <div className="mb-4">
        <h1 className="text-2xl font-bold">Users</h1>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((hg) => (
              <TableRow key={hg.id}>
                {hg.headers.map((h) => (
                  <TableHead key={h.id}>
                    {h.isPlaceholder ? null : flexRender(h.column.columnDef.header, h.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">No users found.</TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-end gap-2 py-4">
        <Button variant="outline" size="sm" onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()}>Previous</Button>
        <Button variant="outline" size="sm" onClick={() => table.nextPage()} disabled={!table.getCanNextPage()}>Next</Button>
      </div>
    </div>
  );
}
