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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
const formSchema = z.object({
  name: z.string().min(2, {
    message: "Name must be at least 2 characters.",
  }).max(255,{
	message: "Name must be at most 255 characters.",
  }),
  Description: z.string(),
})
import { apiClient } from "@/utils/apiClient";
import { Dialog, DialogContent, DialogHeader } from "@/components/ui/dialog"
import { DialogDescription, DialogTitle } from "@radix-ui/react-dialog"
import Link from "next/link"
import { useParams } from "next/navigation"
import { Skeleton } from "@/components/ui/skeleton"
import { decodeJwt } from "@/utils/userToken";

const decodedToken = decodeJwt();
export default function RoleEditPage() {
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
  const params = useParams();
  const id = params?.id;
	const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
	const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
	const [dialogStatus, setDialogStatus] = useState<"success" | "error" | null>(null);
  const [loading, setLoading] = useState(true);

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			Description: "",
		},
	  })
    type InstanceViewResponse = {
      Name: string;
      Description: string;
    };
    useEffect(() => {
      async function fetchInstance() {
        try {
          const response = await apiClient.get(`/api/roles/view/${id}`);
          if (response.status === 200) {
            const data = response.data as InstanceViewResponse;
            form.reset({
              name: data.Name || "",
              Description: data.Description || ""
            });
            setLoading(false)
          }
        } catch (error) {
          // Optionally handle error
        }
      }
      if (id) fetchInstance();
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id]);
    type RoleResponse = {
      message: string;
      error_description?: string;
    };
	  async function onSubmit(values: z.infer<typeof formSchema>) {
		try {
		  const response = await apiClient.put(`/api/roles/edit/${id}`, values);
      const data = response.data as RoleResponse;
		  if (response.status === 200) {
		    setDialogStatus("success");
		    setDialogDescription(data.message);
		  } else {
		    setDialogStatus("error");
		    setDialogDescription(data.error_description || "Unknown error");
		  }
		} catch (error) {
		  setDialogStatus("error");
		  setDialogDescription("An unexpected error occurred.");
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
  return (
	<>
	<div className="p-8"><div className="flex items-baseline gap-x-4 pb-6">
	  <h1 className="text-2xl font-bold pb-6">Edit Role</h1>
<Link className="" href={`/roles/view/${id}`}>
<Button variant="link"><IconChevronLeft stroke={2} />
Go Back to Role</Button>
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
                <Input placeholder="Infra Team" {...field} />
              </FormControl>
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
                  placeholder="Provide a brief description of the role and its responsibilities"
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
            <div className="flex w-full gap-4 mt-6">
				<Link className="w-1/2" href={`/roles/view/${id}`}>
              <Button
                type="button"
                className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none"
                // onClick={...} // Add your view instance handler here
              >
                View Role
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
	  </>
  );
}