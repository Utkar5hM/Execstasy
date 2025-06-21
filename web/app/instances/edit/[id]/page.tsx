"use client"
 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Textarea } from "@/components/ui/textarea"
import { IconChevronLeft } from '@tabler/icons-react';

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
  } from "@/components/ui/select"
const formSchema = z.object({
  name: z.string().min(2, {
    message: "Name must be at least 2 characters.",
  }).max(255,{
	message: "Name must be at most 255 characters.",
  }),
  Description: z.string(),
  Status: z.enum(["active", "disabled"], ),
  HostAddress: z.string().max(255,{
	message: "Name must be at most 255 characters.",
  }).refine((value) => {
    if (value === "") return true; // Allow empty strings
    const hostnameRegex = /^(([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}|localhost|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(\[([0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}\]))$/;
    const ipv6Regex =  /((^\s*((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\s*$)|(^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$))/;
	return hostnameRegex.test(value) || ipv6Regex.test(value);
  }, {
    message: "Host address must be a valid hostname, IPv4, or IPv6 address, or empty.",
  }),
})
import { apiClient } from "@/utils/apiClient";
import { CheckCircle, CircleOff } from "lucide-react"
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTrigger } from "@/components/ui/dialog"
import { DialogDescription, DialogTitle } from "@radix-ui/react-dialog"
import Link from "next/link"
import { useParams } from "next/navigation"
import { ExecTable } from "@/components/execTable"
import { Skeleton } from "@/components/ui/skeleton"
import { Row } from "@tanstack/react-table"; // adjust import if needed
import { DefaultStatusResponse } from '@/utils/ResponseTypes';

import { decodeJwt } from "@/utils/userToken";

const decodedToken = decodeJwt();
export default function InstanceEditPage() {
	const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
	const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
	const [dialogStatus, setDialogStatus] = useState<"success" | "error" | null>(null);
	const [dialogClientID, setdialogClientID] = useState(""); // State to store the dialog description
	const [dialogInstanceID, setdialogInstanceID] = useState(""); // State to store the dialog description
  const [usersData, setUsersData] = useState<{ host_username: string }[]>([]);
  const [userDialogOpen, setUserDialogOpen] = useState(false);
  const [hostUserstatusDialogOpen, sethostUserstatusDialogOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedHostUsername, setSelectedHostUsername] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const UserFormSchema = z.object({
    host_username: z.string().min(1, {
      message: "Host Username must be at least 2 characters.",
    }),
  })

  const formUser = useForm<z.infer<typeof UserFormSchema>>({
    resolver: zodResolver(UserFormSchema),
    defaultValues: {
      host_username: "*",
    },
  })
  const id = useParams()?.id;
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/view/${id}`);
        if (response.status === 200) {
          const data = response.data as InstanceViewResponse;
          const status = data.Status === "active" || data.Status === "disabled" ? data.Status : "active";
          form.reset({
            name: data.Name || "",
            Description: data.Description || "",
            Status: status,
            HostAddress: data.HostAddress || "",
          });setUsersData(
            (data.HostUsers || [])
              .map((u: string) => ({ host_username: u }))
          );
          setLoading(false)
        }
      } catch (error) {
        console.error("Error fetching instance data:", error);
        setError("Failed to fetch instance data: " + error);
        // Optionally handle error
      }
    }
    if (id) fetchInstance();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			Description: "",
			Status: "active",
			HostAddress: "",
		},
	  })

  if (decodedToken?.role !== "admin") {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Access Denied</h1>
          <p className="text-gray-600">You do not have permission to view this page.</p>
        </div>
      </div>
    );
  }
  const Userscolumns = [
    {
      accessorKey: "host_username",
      header: "Host Username",
    },
    {
      accessorKey: "",
      header: "Actions",
      enableHiding: false,
      cell: ({ row }: { row: Row<{ host_username: string }> }) => {
        const host_username = row.original.host_username;
      
        const handleDelete = async (e: React.FormEvent) => {
          e.preventDefault();
          try {
            const response = await apiClient.delete<DefaultStatusResponse>(`/api/instances/hostUsers/${id}`, {
               host_username: selectedHostUsername 
              }
            );
            if (response.status === 200) {
              sethostUserstatusDialogOpen(true);
              setDialogDescription(`${response.data.status}: ${response.data.message}`);
              const updatedUsersData = usersData.filter(user => user.host_username !== selectedHostUsername);
              setUsersData(updatedUsersData);
            } else {
              sethostUserstatusDialogOpen(true);
              setDialogDescription(`${response.data.error}: ${response.data.error_description}`);
            }
            // Optionally refresh users list or show a dialog
          } catch (error) {
            console.error("Failed to delete host user:", error);
            setError("Failed to delete host user: " + error);
          }
          setDialogOpen(false);
          setSelectedHostUsername(null);
        };
      
        return (
          <>
            <Dialog
              open={dialogOpen && selectedHostUsername === host_username}
              onOpenChange={(open) => {
                setDialogOpen(open);
                if (open) setSelectedHostUsername(host_username);
                else setSelectedHostUsername(null);
              }}
            >
              <DialogTrigger asChild>
                <Button
                  variant="destructive"
                  size="sm"
                  className="hidden h-8 lg:flex w-[80px]"
                  onClick={() => {
                    setSelectedHostUsername(host_username);
                    setDialogOpen(true);
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
    type InstanceViewResponse = {
      Name: string;
      Description: string;
      Status: string;
      HostAddress?: string;
      HostUsers: Array<string>;
    };
    type InstanceResponse = {
      message: string;
      client_id: string;
      id: string;
      error_description?: string;
    };
      async function onSubmitUser(data: z.infer<typeof UserFormSchema>) {
        setUserDialogOpen(false);
        try{
        const response = await apiClient.post<DefaultStatusResponse>(`/api/instances/hostUsers/${id}`, data);
        console.log("API Response:", response);
        if (response.status === 200) {
          // Success: Open the dialog and set the message
          setDialogDescription(`${response.data.status}: ${response.data.message}`);
          setUsersData((prevData) => [...prevData, { host_username: data.host_username }]);
        } else {
          // Error: Open the dialog and set the error message
          setDialogDescription(`${response.data.error}: ${response.data.error_description}`);
        }
    
        sethostUserstatusDialogOpen(true); // Open the status dialog
      } catch (error) {
        console.error("Error submitting form:", error);
        setDialogDescription("An unexpected error occurred.");
        sethostUserstatusDialogOpen(true); // Open the status dialog for unexpected errors
      }
      }
	  async function onSubmit(values: z.infer<typeof formSchema>) {
		try {
		  const response = await apiClient.put(`/api/instances/edit/${id}`, values);
      const data = response.data as InstanceResponse;
		  if (response.status === 200) {
		    setDialogStatus("success");
		    setDialogDescription(data.message);
		    setdialogClientID(data.client_id);
			  setdialogInstanceID(data.id);
		  } else {
		    setDialogStatus("error");
		    setDialogDescription(data.error_description || "Unknown error");
		    setdialogClientID("");
		  }
		} catch (error) {
		  setDialogStatus("error");  
      setDialogDescription(
        error instanceof Error ? error.message : String(error)
      );
		  setdialogClientID("");
		} finally {
		  setStatusDialogOpen(true);
		}
	  }
  if (loading ) {
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
  return (
	<>
	<div className="p-8"><div className="flex items-baseline gap-x-4 pb-6">
	  <h1 className="text-2xl font-bold pb-6">Edit Instance</h1>
<Link className="" href={`/instances/view/${id}`}>
<Button variant="link"><IconChevronLeft stroke={2} />
Go Back to Instance</Button>
</Link>
</div>
	  <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Proxy Server 001" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="HostAddress"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Host Address</FormLabel>
              <FormControl>
                <Input placeholder="10.15.45.45" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
		<FormField
          control={form.control}
          name="Status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Status</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Set Instance Status" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="active"><CheckCircle className="w-5 h-5" />
				  <span>Active</span></SelectItem>
                  <SelectItem value="disabled"><CircleOff className="w-5 h-5" />
				  <span>Disabled</span></SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />
		<FormField
          control={form.control}
          name="Description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="A bit of information about the Instance"
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Update</Button>
      </form>
    </Form>
	</div>
            {/* Status Dialog */}
			<Dialog open={statusDialogOpen} onOpenChange={setStatusDialogOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Status</DialogTitle>
        {dialogStatus === "success" ? (
          <div>
            <div className="mb-4">{dialogDescription}</div>
            <div className="mb-2">
              <span className="block font-medium mb-1">Client ID:</span>
              <div className="flex items-center gap-2">
                <Input
                  disabled
                  className="flex-1"
                  placeholder="client ID"
                  value={dialogClientID}
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigator.clipboard.writeText(dialogClientID)}
                >
                  Copy
                </Button>
              </div>
            </div>
            <div className="flex w-full gap-4 mt-6">
				<Link className="w-1/2" href={`/instances/view/${dialogInstanceID}`}>
              <Button
                type="button"
                className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none"
                // onClick={...} // Add your view instance handler here
              >
                View Instance
              </Button>
			  </Link>
              <Button
                type="button"
                className="w-1/2 text-white focus:ring-2 focus:outline-none bg-gray-500 hover:bg-gray-600"
                onClick={() => setStatusDialogOpen(false)}
              >
                Close
              </Button>
            </div>
          </div>
        ) : (
          <DialogDescription>
          <>
            {dialogDescription}
            <Button
              type="button"
              className="w-full mt-4 text-white focus:ring-2 focus:outline-none bg-gray-500 hover:bg-gray-600"
              onClick={() => setStatusDialogOpen(false)}
            >
              Close
            </Button></>
      </DialogDescription>)}
    </DialogHeader>
  </DialogContent>
</Dialog>
<div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Host Users</h1>
      <ExecTable data={usersData}
        columns={Userscolumns}
        filterColumn="host_username"
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
          name="host_username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Host Username</FormLabel>
              <FormControl>
                <Input placeholder="root" {...field} />
              </FormControl>
              <FormDescription>
              This is the user on the host, &quot;*&quot; to allow all users on the host.
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
      <Dialog open={hostUserstatusDialogOpen}   onOpenChange={sethostUserstatusDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Status2</DialogTitle>
            <DialogDescription>{dialogDescription}</DialogDescription>
          </DialogHeader>

          <div className="flex w-full gap-4 mt-6">
				<Link className="w-1/2" href={`/instances/view/${id}`}>
              <Button
                type="button"
                className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none"
                // onClick={...} // Add your view instance handler here
              >
                View Instance
              </Button>
			  </Link>
              <Button
                type="button"
                className="w-1/2 text-white focus:ring-2 focus:outline-none bg-gray-500 hover:bg-gray-600"
                onClick={() => sethostUserstatusDialogOpen(false)}
              >
                Close
              </Button>
            </div>
        </DialogContent>
      </Dialog>
	  </>
  );
}