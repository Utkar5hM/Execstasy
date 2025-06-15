"use client";


import * as React from "react"
import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table"
import { ArrowUpDown, BadgeCheckIcon, ChevronDown,  MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import Link from "next/link"
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { apiClient } from "@/utils/apiClient";
import { Skeleton } from "@/components/ui/skeleton";
import { DataTableColumnHeader } from "../instances/components/data-table-column-header";

export type User = {
  id: number;
  created_at: string;
  updated_at: string;
  created_by: string;
  description: string;
  name: string;
  users: string;
  role: string
};
import { Badge } from "@/components/ui/badge"

import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DefaultStatusResponse } from "@/utils/ResponseTypes";

import {decodeJwt} from "@/utils/userToken"
const decodedToken = decodeJwt();

export default function UsersPage() {
  const [data, setData] = React.useState<User[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [statusDialogOpen, setStatusDialogOpen] = React.useState(false);
  const [dialogDescription, setDialogDescription] = React.useState("");

const columns: ColumnDef<User>[] = [
  {
    accessorKey: "role",
    header: "",
    cell: "",
    enableHiding: false
  },
  {
    accessorKey: "name",
    header: "Name",
    enableColumnFilter: true, // Enable filtering for this column
    cell: ({ row }) => (
      <div className="flex items-center space-x-2">
        <span>{row.getValue("name")}</span>
        {row.original.role === "admin" && ( // Check if the role is 'admin'
          <Badge
            variant="outline"
            className="rounded-full"
          >
            <BadgeCheckIcon />
            Admin
          </Badge>
        )}
      </div>
    ),
  },
  {
    accessorKey: "username",
    header: "Username",
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "id",
    header: "",
    enableHiding: false,
    cell: ({ row }) => {
      const handleRoleChange = async () => {
        try {
          const response = await apiClient.put<DefaultStatusResponse>('/api/users/role', {
            id: row.original.id
          })
          if (response.status === 200) {
            setDialogDescription(`${response.data.status}: ${response.data.message}`);
          } else {
            setDialogDescription(`${response.data.error}: ${response.data.error_description}`);
          }
        }catch (error) {
          setDialogDescription("An unexpected error occurred.");
          setStatusDialogOpen(true);
        } finally {
          setStatusDialogOpen(true);
        }
      } 
      if (decodedToken?.role !== "admin" || row.original.id === decodedToken?.id) {
        return null; // Hide the button if the user is not an admin
      }
      return (
        <Button
          variant="outline"
          size="sm"
          onClick={handleRoleChange}
          className="ml-auto hidden h-8 lg:flex w-[120px] mr-6"
        >
          {row.original.role === "admin" ? "Demote" : "Make Admin"}
        </Button>)
    }
}
//   {
//     accessorKey: "id",
//     header: "",
//     enableHiding: false,
//     cell: ({ row }) => (
//       <Link href={`/Users/view/${row.getValue("id")}`}>
// <Button
//         variant="outline"
//         size="sm"
//         className="ml-auto hidden h-8 lg:flex w-[80px]"
//       >View
//       </Button>
//       </Link>
//     ),
//   }
];
  React.useEffect(() => {
    async function fetchUsers() {
      try {
        const response = await apiClient.get<User[]>("/api/users");
        setData(response.data); // Update the data state with fetched Users
      } catch (err) {
        console.error("Failed to fetch users:", err);
        setError("Failed to load users.");
      } finally {
        setLoading(false);
      }
    }
    fetchUsers();
  }, []);

const [sorting, setSorting] = React.useState<SortingState>([]);

const table = useReactTable({
  data,
  columns,
  state: {
    columnFilters, // Add columnFilters state
  },
  onColumnFiltersChange: setColumnFilters, // Handle column filter changes
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
  getFilteredRowModel: getFilteredRowModel(), // Add filtered row model
});

if (loading) {
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="flex items-center space-x-4">
        <Skeleton className="h-12 w-12 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-4 w-[250px]" />
          <Skeleton className="h-4 w-[200px]" />
        </div>
      </div>
    </div>
      );
}
  return (
    <>
    <div className="p-8">
      <h1 className="text-2xl font-bold">Users</h1>
      <div className="w-full">
        <div className="flex items-center py-4">
          <Input
            placeholder="Filter Users by Name..."
            value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
            onChange={(event) =>
              table.getColumn("name")?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="ml-auto">
                Columns <ChevronDown />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {table
                .getAllColumns()
                .filter((column) => column.getCanHide())
                .map((column) => {
                  return (
                    <DropdownMenuCheckboxItem
                      key={column.id}
                      className="capitalize"
                      checked={column.getIsVisible()}
                      onCheckedChange={(value) =>
                        column.toggleVisibility(!!value)
                      }
                    >
                      {column.columnDef.header}
                    </DropdownMenuCheckboxItem>
                  );
                })}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHead key={header.id}>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                      </TableHead>
                    );
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    No results.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
        <div className="flex items-center justify-end space-x-2 py-4">
          <div className="text-muted-foreground flex-1 text-sm">
            {table.getFilteredSelectedRowModel().rows.length} of{" "}
            {table.getFilteredRowModel().rows.length} row(s) selected.
          </div>
          <div className="space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              Next
            </Button>
          </div>
        </div>
      </div>
    </div>
{/* Status Dialog */}
<Dialog open={statusDialogOpen}   onOpenChange={setStatusDialogOpen}>
<DialogContent>
  <DialogHeader>
    <DialogTitle>Status</DialogTitle>
    <DialogDescription>{dialogDescription}</DialogDescription>
  </DialogHeader>
  <Button
    type="button"
    className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none mt-4"
    onClick={() => setStatusDialogOpen(false)} // Close the dialog
  >
    Close
  </Button>
</DialogContent>
</Dialog> 
</>    
  );
}