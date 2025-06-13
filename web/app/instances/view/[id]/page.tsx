"use client"
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
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
import { Label } from "@/components/ui/label"



import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
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
const FormSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  host_username: z.string().min(1, {
    message: "Host Username must be at least 2 characters.",
  }),
})

export default function InstanceViewPage() {
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
  const [dialogOpen, setDialogOpen] = useState(false); // State to control the dialog visibility
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      username: "",
      host_username: "*",
    },
  })
  async function onSubmit(data: z.infer<typeof FormSchema>) {
    setDialogOpen(false);
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

  const data = [
    { id: 1, name: "John Doe", email: "john@example.com", role: "admin" },
    { id: 2, name: "Jane Smith", email: "jane@example.com", role: "user" },
  ];
  
  const columns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "email",
      header: "Email",
    },
    {
      accessorKey: "role",
      header: "Role",
    },
  ];
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
      <p className="pt-4 flex items-center space-x-2"><strong>Created By:&nbsp;</strong> {instance.CreatedBy}</p>
      <p className="pt-4 flex items-center space-x-2">
<Button
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Edit
      </Button>
      <Button
        variant="destructive"
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Delete
      </Button>
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
          <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogTrigger asChild>
          <Button >Add User</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
      
      <DialogHeader>
            <DialogTitle>Allow User to Access Instance</DialogTitle>
            <DialogDescription>
            </DialogDescription>
          </DialogHeader>
        <FormField
          control={form.control}
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
          control={form.control}
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
            <Button>Add Role</Button>
          </>
        } />
      </div>
            {/* Status Dialog */}
            <Dialog open={statusDialogOpen} onOpenChange={setStatusDialogOpen}>
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