"use client"
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiClient } from "@/utils/apiClient";
import {
  CheckCircle,
  CircleOff,
} from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import { ExecTable } from "@/components/execTable";
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton";
import * as React from "react"
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



import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
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
  host_username: z.string().min(1, {
    message: "Host Username must be at least 2 characters.",
  }),
})


const RoleFormSchema = z.object({
  id: z.number().min(1, {
    message: "Role ID must be at least 1.",
  }),
  host_username: z.string().min(1, {
    message: "Host Username must be at least 2 characters.",
  }),
})

import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import Link from "next/link";


export default function InstanceViewPage() {
  const router = useRouter();
	const params = useParams(); // Access the params object
	const id = params?.id; // Extract the `id` parameter
  const [instance, setInstance] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [loading2, setLoading2] = useState(true);
  const [loading3, setLoading3] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [usersData, setUsersData] = useState([]);
  const [rolesData, setRolesData] = useState([]);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
  const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
  const [open, setOpen] = React.useState(false)
  const [rolesFetched, setRolesFetched] = useState<{ value: string; label: string }[]>([]);
  const [userDialogOpen, setUserDialogOpen] = useState(false);
  const [roleDialogOpen, setRoleDialogOpen] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [instanceDeleted, setInstanceDeleted] = useState(false); // State to track if the instance was deleted
  useEffect(() => {
    async function fetchRoles() {
      try {
        const response = await apiClient.get("/api/roles");
        // Transform the API data to { value, label }
        const roles = (response.data).map((role: any) => ({
          value: String(role.id),
          label: role.name,
        }));
        setRolesFetched(roles);
      } catch (err) {
        setRolesFetched([]);
      }
    }
    fetchRoles();
  }, []);
  const formUser = useForm<z.infer<typeof UserFormSchema>>({
    resolver: zodResolver(UserFormSchema),
    defaultValues: {
      username: "",
      host_username: "*",
    },
  })
  const formRole = useForm<z.infer<typeof RoleFormSchema>>({
    resolver: zodResolver(RoleFormSchema),
    defaultValues: {
      id: 0,
      host_username: "*",
    },
  })

  async function onSubmitUser(data: z.infer<typeof UserFormSchema>) {
    setUserDialogOpen(false);
    try{
    const response = await apiClient.post(`/api/instances/users/${id}`, data);
    console.log("API Response:", response);
    if (response.status === 200) {
      // Success: Open the dialog and set the message
      setDialogDescription(`${response.data.status}: ${response.data.message}`);
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
        router.push("/instances");
      }, 0);
    }
  }, [statusDialogOpen, instanceDeleted]);
  async function onSubmitRole(data: z.infer<typeof RoleFormSchema>) {
    setRoleDialogOpen(false);
    try{
    const response = await apiClient.post(`/api/instances/roles/${id}`, data);
    console.log("API Response:", response);
    if (response.status === 200) {
      // Success: Open the dialog and set the message
      setDialogDescription(`${response.data.status}: ${response.data.message}`);
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
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/view/${id}`);
        setInstance(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading(false);
      }
    }

    fetchInstance();
  }, [id]);

  async function deleteInstance(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const response = await apiClient.delete(`/api/instances/${id}`);
      if (response.status === 200) {
        setDialogDescription("Instance deleted successfully.");
        setInstanceDeleted(true); // Set instanceDeleted to true
      } else {
        setDialogDescription("Failed to delete instance.");
      }
    } catch (error) {
      setDialogDescription("An unexpected error occurred while deleting the instance.");
    } finally {
      setDialogOpen(false);
      setStatusDialogOpen(true); // Optionally show a status dialog
    }
  }
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/users/${id}`);
        setUsersData(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading2(false);
      }
    }

    fetchInstance();
  }, [id]);
  
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/roles/${id}`);
        setRolesData(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading3(false);
      }
    }

    fetchInstance();
  }, [id]);

  const RolesColumns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "host_username",
      header: "Host Username",
    },
  ]
  const Userscolumns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "username",
      header: "username",
    },
    {
      accessorKey: "host_username",
      header: "Host Username",
    },
  ];


  if (loading || loading2 || loading3) {
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

  if (error) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }


  return (<>
    <div className="p-8">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight text-balance">
          {instance.Name}</h1>

      <p className="pt-4  text-muted-foreground text-xl">{instance.Description}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Host Address:</strong>     <code className="bg-muted relative rounded px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
      {instance.HostAddress}
    </code>
      </p>
      <p className="pt-4 flex items-center space-x-2 ">
      <strong>Status:</strong>
  {instance.Status === "active" ? (
    <span className="text-muted-foreground flex items-center space-x-1">
      <CheckCircle className="w-5 h-5" />
      <span>Active</span>
    </span>
  ) : (
    <span className="flex items-center space-x-1">
      <CircleOff className="w-5 h-5" />
      <span>Disabled</span>
    </span>
  )}
</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Host Users:&nbsp;</strong> {instance.HostUsers && instance.HostUsers.join(", ")}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Created By:&nbsp;</strong> {instance.CreatedBy}</p>
      <p className="pt-4 flex items-center space-x-2">
<Link href={`/instances/edit/${id}`}>
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
                    <form onSubmit={deleteInstance} className="flex flex-col items-center">
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
      <Separator className="my-6" />
    </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Users</h1>
      <ExecTable data={usersData}
        columns={Userscolumns}
        filterColumn="name"
        headerContent={
          <>
          <Dialog open={userDialogOpen} onOpenChange={setUserDialogOpen}>
        <DialogTrigger asChild>
          <Button >Add User</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
        <Form {...formUser}>
      <form onSubmit={formUser.handleSubmit(onSubmitUser)} className="space-y-6">
      
      <DialogHeader>
            <DialogTitle>Allow User to Access Instance</DialogTitle>
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
          name="host_username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Host Username</FormLabel>
              <FormControl>
                <Input placeholder="root" {...field} />
              </FormControl>
              <FormDescription>
                This is the user on the host, type "*" to allow all users on the host.
              </FormDescription>
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
    </Dialog>
          </>
        } />
      </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Roles</h1>
      <ExecTable data={rolesData} 
      columns={RolesColumns} 
      filterColumn="name" 
        headerContent={
          <>
          <Dialog open={roleDialogOpen} onOpenChange={setRoleDialogOpen}>
        <DialogTrigger asChild>
          <Button >Add Role</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
        <Form {...formRole}>
      <form onSubmit={formRole.handleSubmit(onSubmitRole)} className="space-y-6">
      
      <DialogHeader>
            <DialogTitle>Allow Role to Access Instance</DialogTitle>
            <DialogDescription>
            </DialogDescription>
          </DialogHeader>
          <FormField
          control={formRole.control}
          name="id"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Role</FormLabel>
              <FormControl>
                <Popover open={open} onOpenChange={setOpen}>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      role="combobox"
                      aria-expanded={open}
                      className="justify-between"
                    >
                      {rolesFetched.find((role) => role.value === String(field.value))?.label || "Select Role..."}
                      <ChevronsUpDown className="opacity-50" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="p-0">
                    <Command>
                      <CommandInput placeholder="Search Role..." className="h-9" />
                      <CommandList>
                        <CommandEmpty>No Role found.</CommandEmpty>
                        <CommandGroup>
                          {rolesFetched.map((role) => (
                            <CommandItem
                              key={role.value}
                              value={role.value}
                              onSelect={() => {
                                field.onChange(Number(role.value));
                                setOpen(false);
                              }}
                            >
                              {role.label}
                              <Check
                                className={cn(
                                  "ml-auto",
                                  String(field.value) === role.value ? "opacity-100" : "opacity-0"
                                )}
                              />
                            </CommandItem>
                          ))}
                        </CommandGroup>
                      </CommandList>
                    </Command>
                  </PopoverContent>
                </Popover>
              </FormControl>
              <FormDescription>
                This is the role on this Site.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={formRole.control}
          name="host_username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Host Username</FormLabel>
              <FormControl>
                <Input placeholder="root" {...field} />
              </FormControl>
              <FormDescription>
                This is the user on the host, type "*" to allow all users on the host.
              </FormDescription>
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
    </Dialog>
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