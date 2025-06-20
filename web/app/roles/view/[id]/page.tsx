"use client"
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiClient } from "@/utils/apiClient";
import {
  BadgeCheckIcon,
} from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import { ExecTable } from "@/components/execTable";
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton";
import * as React from "react"
import { CellContext, ColumnDef } from "@tanstack/react-table";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import {  z } from "zod"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
const UserFormSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  role: z.enum(["standard", "admin"], {
    required_error: "Role is required.",
    invalid_type_error: "Role must be either 'standard' or 'admin'.",
  }),
})


import Link from "next/link";
import { APIAddRoleUser, APIRole, APIRolesAccess, APIRoleUsers, DefaultStatusResponse } from "@/utils/ResponseTypes";
import { decodeJwt } from "@/utils/userToken";
import { IconChevronLeft } from "@tabler/icons-react";

const decodedToken = decodeJwt();

export default function RoleViewPage() {
  const router = useRouter();
	const params = useParams(); // Access the params object
	const id = params?.id; // Extract the `id` parameter
  const [roleInfo, setRoleInfo] = useState<APIRole | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [loading2, setLoading2] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [usersData, setUsersData] = useState<APIRoleUsers[]>([]);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
  const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
  const [userDialogOpen, setUserDialogOpen] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [userDeleteDialogOpen, setUserDeleteDialogOpen] = useState(false);
  const [instanceDeleted, setInstanceDeleted] = useState(false); // State to track if the instance was deleted
  const [deleteUsername, setDeleteUsername] = useState<string | null>(null);
  const [canUpdateRole, setCanUpdateRole] = React.useState(false)

  useEffect(() => {
    async function checkCanUpdateRole() {
      try {
        const response = await apiClient.get<APIRolesAccess | DefaultStatusResponse>(`/api/roles/access/${id}`);
        if (response.status === 200) {
          const data = response.data as APIRolesAccess;
          setCanUpdateRole(data.access);
        } else {
          const data = response.data as DefaultStatusResponse;
          setError("Failed to check role update permission: " + data.error);
        }
      } catch (error) {
        setError("Failed to check role update permission: " + error);
      }
    }
    checkCanUpdateRole();
  }
  , [id]);
    

  const formUser = useForm<z.infer<typeof UserFormSchema>>({
    resolver: zodResolver(UserFormSchema),
    defaultValues: {
      username: "",
      role: "standard",
    },
  })

  async function onSubmitUser(data: z.infer<typeof UserFormSchema>) {
    setUserDialogOpen(false);
    try{
    const response = await apiClient.post<APIAddRoleUser>(`/api/roles/users/${id}`, data);
    if (response.status === 200) {
      // Success: Open the dialog and set the message
      setDialogDescription(`${response.data.status}: ${response.data.message}`);
      const fetchedData = response.data.data || {
        name: "Newly Added User",
        username: data.username,
        role: data.role,
        user_id: 0, 
      };
      const newUser: APIRoleUsers = {
        name: fetchedData.name,
        username: data.username,
        role: data.role,
        id: fetchedData.user_id,
      };
      // Add the new user to the usersData state
      setUsersData((prevUsers) => [...prevUsers, newUser]);
    } else {
      // Error: Open the dialog and set the error message
      setDialogDescription(`${response.data.error}: ${response.data.error_description}`);
    }

    setStatusDialogOpen(true); // Open the status dialog
  } catch (error) {
    console.error("Error submitting form:", error);
    setDialogDescription("An unexpected error occurred.");
    setStatusDialogOpen(true); // Open the status dialog for unexpected errors
  }
  }
  useEffect(() => {
    if (!statusDialogOpen && instanceDeleted) {
      setTimeout(() => {
        router.push("/roles");
      }, 0);
    }
  }, [statusDialogOpen, instanceDeleted, router]);
  useEffect(() => {
    async function fetchRole() {
      try {
        const response = await apiClient.get<APIRole | DefaultStatusResponse>(`/api/roles/view/${id}`);
        if (response.status ===200){
          setRoleInfo(response.data as APIRole);
        } else {
          const data = response.data as DefaultStatusResponse;
          setError("Failed to fetch role: " + data.error + "\n" + data.error_description);
        }
      } catch (err) {
        console.error("Failed to fetch role:", err);
        setError("Failed to load role.");
      } finally {
        setLoading(false);
      }
    }

    fetchRole();
  }, [id]);

  async function deleteRole(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const response = await apiClient.delete<DefaultStatusResponse>(`/api/roles/${id}`, null);
      if (response.status === 200) {
        setDialogDescription("Role deleted successfully.");
        setInstanceDeleted(true); 
      } else {
        setDialogDescription(response.data.error + " : " + response.data.error_description);
      }
    } catch (error) {
      setDialogDescription("An unexpected error occurred while deleting the instance." + error);
    } finally {
      setDialogOpen(false);
      setStatusDialogOpen(true); // Optionally show a status dialog
    }
  }
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get<APIRoleUsers[] | DefaultStatusResponse>(`/api/roles/users/${id}`);
        if (response.status ===200) {
          setUsersData(response.data as APIRoleUsers[]);
        } else {
          const data = response.data as DefaultStatusResponse;
          setError("Failed to check role update permission: " + data.error);
        }
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading2(false);
      }
    }

    fetchInstance();
  }, [id]);

  const Userscolumns: ColumnDef<APIRoleUsers, unknown>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => (
        <div className="flex items-center space-x-2">
          <span>{row.getValue("name")}</span>
          {row.original.role === "admin" && ( // Check if the role is 'admin'
            <Badge
              variant="outline"
              className="rounded-full"
            >
              <BadgeCheckIcon />
              Role Admin
            </Badge>
          )}
        </div>
      ),
    },
    {
      accessorKey: "username",
      header: "username",
    },
    {
      accessorKey: "",
      header: "",
      id: "delete",
      enableHiding: false,
      cell: (props: CellContext<APIRoleUsers, unknown>) => {
        if (!canUpdateRole) {
          return null; // Hide the delete button if the user cannot update the role
        }
        const username = props.row.original.username;
        const handleDelete = async (e: React.FormEvent) => {
          e.preventDefault();
          try {
            const response = await apiClient.delete<DefaultStatusResponse>(`/api/roles/users/${id}`, {
               username: deleteUsername,
              }
            );
            if (response.status === 200) {
              setStatusDialogOpen(true);
              setDialogDescription(`${response.data.status}: ${response.data.message}`);
              const updatedUsers = usersData.filter(user => user.username !== deleteUsername);
              setUsersData(updatedUsers); // Update the usersData state to remove the deleted user
            } else {
              setStatusDialogOpen(true);
              setDialogDescription(`${response.data.error}: ${response.data.error_description}`);
            }
            // Optionally refresh users list or show a dialog
          } catch (error) {
            console.error("Failed to delete host user:", error);
          }
          setUserDeleteDialogOpen(false);
          setDeleteUsername(null);
        };
      
        return (
          <>
            <Dialog
              open={userDeleteDialogOpen && deleteUsername === username}
              onOpenChange={(open) => {
                setUserDeleteDialogOpen(open);
                if (open) {
                  setDeleteUsername(username);
                }
                else {
                  setDeleteUsername(null);
                }
              }}
            >
              <DialogTrigger asChild>
                <Button
                  variant="destructive"
                  size="sm"
                  className="hidden h-8 lg:flex w-[80px] ml-auto mr-2"
                  onClick={() => {
                    setDeleteUsername(username);
                    setUserDeleteDialogOpen(true);
                  }}
                >
                  Delete
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Are you absolutely sure?</DialogTitle>
                  <DialogDescription>
                    This action cannot be undone. If you are sure, click the button below to proceed.
                  </DialogDescription>
                </DialogHeader>
                <form className="flex flex-col items-center" onSubmit={handleDelete}>
                  <Button
                    type="submit"
                    variant="outline"
                    className="w-full focus:ring-2 mt-4"
                  >
                    Yes
                  </Button>
                </form>
              </DialogContent>
            </Dialog>
          </>
        );
      }
    }
    
  ];


  if (loading || loading2) {
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

  if (error || !roleInfo) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }


  return (<>
    <div className="p-8"><div className="flex items-baseline gap-x-4">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight text-balance">
          {roleInfo.Name}</h1>
          <Link className="" href={`/roles`}>
<Button variant="link"><IconChevronLeft stroke={2} />
Go Back to Roles</Button>
</Link>
</div>
      <p className="pt-4  text-muted-foreground text-xl">{roleInfo.Description}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Created By:&nbsp;</strong> {roleInfo.CreatedBy}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Last Updated At:&nbsp;</strong> {roleInfo.UpdatedAt}</p>
      {decodedToken?.role === "admin" && ( 
      <p className="pt-4 flex items-center space-x-2">
<Link href={`/roles/edit/${id}`}>
<Button
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Edit
      </Button>
      </Link>
                  <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                  <DialogTrigger asChild>

                  <Button
                variant="destructive"
                size="sm"
                className="hidden h-8 lg:flex w-[80px]"
              >Delete
              </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Are you absolutely sure?</DialogTitle>
                      <DialogDescription>
                        This action cannot be undone. If you are sure, click the button below to proceed.
                      </DialogDescription>
                    </DialogHeader>
                    <form onSubmit={deleteRole} className="flex flex-col items-center">
            <Button
              type="submit" 
              variant="outline"
              className="w-full focus:ring-2  mt-4"
            >Yes
            </Button>
            </form>
                  </DialogContent>
                </Dialog>
      </p>
      )}
      <Separator className="my-6" />
    </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Users</h1>
      <ExecTable data={usersData}
        columns={Userscolumns}
        filterColumn="name"
        headerContent={
          <>
          {canUpdateRole && (
          <Dialog open={userDialogOpen} onOpenChange={setUserDialogOpen}>
        <DialogTrigger asChild>
          <Button >Add User</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
        <Form {...formUser}>
      <form onSubmit={formUser.handleSubmit(onSubmitUser)} className="space-y-6">
      
      <DialogHeader>
            <DialogTitle>Add User to Role</DialogTitle>
            <DialogDescription>
            </DialogDescription>
          </DialogHeader>
        <FormField
          control={formUser.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="username" {...field} />
              </FormControl>
              <FormDescription>
                This is the user on this Site.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={formUser.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Role</FormLabel>
              <div className="w-full">
              <Select onValueChange={field.onChange} defaultValue={field.value} >
                <FormControl>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Select a role for the user" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="standard" className="w-full">standard</SelectItem>
                  <SelectItem value="admin" className="w-full">admin</SelectItem>
                </SelectContent>
              </Select>
              </div>
              <FormMessage />
            </FormItem>
          )}
        />
          <DialogFooter>
        <Button type="submit">Add</Button>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
        </DialogFooter>
      </form>
    </Form>
        </DialogContent>
    </Dialog>)}
          </>
        } />
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