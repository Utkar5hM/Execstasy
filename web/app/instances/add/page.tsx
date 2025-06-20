"use client"
 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Textarea } from "@/components/ui/textarea"

import { useState } from "react";
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
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
import { Dialog, DialogContent, DialogHeader } from "@/components/ui/dialog"
import { DialogDescription, DialogTitle } from "@radix-ui/react-dialog"
import Link from "next/link"
import { decodeJwt } from "@/utils/userToken";
import { CreateInstanceResponse, DefaultStatusResponse } from "@/utils/ResponseTypes"

const decodedToken = decodeJwt();
export default function InstanceAddPage() {
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
	const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
	const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
	const [dialogStatus, setDialogStatus] = useState<"success" | "error" | null>(null);
	const [dialogClientID, setdialogClientID] = useState(""); // State to store the dialog description
	const [dialogInstanceID, setdialogInstanceID] = useState<number | undefined>(undefined);
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			Description: "",
			Status: "active",
			HostAddress: "",
		},
	  })
	 
	  async function onSubmit(values: z.infer<typeof formSchema>) {
		try {
		  const response = await apiClient.post<CreateInstanceResponse | DefaultStatusResponse>(`/api/instances`, values);
		  if (response.status === 200) {
        const data = response.data as CreateInstanceResponse;
		    setDialogStatus("success");
		    setDialogDescription(data.message);
		    setdialogClientID(data.client_id);
			  setdialogInstanceID(data.id);
		  } else {
        const data = response.data as DefaultStatusResponse;
		    setDialogStatus("error");
		    setDialogDescription(data.error_description ?? "An unknown error occurred.");
		    setdialogClientID("");
		  }
		} catch (error) {
		  setDialogStatus("error");
		  setDialogDescription("An unexpected error occurred.");
		  setdialogClientID("");
		} finally {
		  setStatusDialogOpen(true);
		}
	  }

  return (
	<>
	<div className="p-8">
	  <h1 className="text-2xl font-bold pb-6">Add Instance</h1>
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
        <Button type="submit">Add</Button>
      </form>
    </Form>
	</div>
            {/* Status Dialog */}
			<Dialog open={statusDialogOpen} onOpenChange={setStatusDialogOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Status</DialogTitle>
      <DialogDescription>
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
          <div>
            <p className="mb-4">{dialogDescription}</p>
            <Button
              type="button"
              className="w-full mt-4 text-white focus:ring-2 focus:outline-none bg-gray-500 hover:bg-gray-600"
              onClick={() => setStatusDialogOpen(false)}
            >
              Close
            </Button>
          </div>
        )}
      </DialogDescription>
    </DialogHeader>
  </DialogContent>
</Dialog>
	  </>
  );
}