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
  username: z.string().min(2, {
    message: "Name must be at least 2 characters.",
  }).max(255,{
	message: "Name must be at most 255 characters.",
  }),
})
import { apiClient } from "@/utils/apiClient";
import { Dialog, DialogContent, DialogHeader } from "@/components/ui/dialog"
import { DialogDescription, DialogTitle } from "@radix-ui/react-dialog"
import Link from "next/link"
import { Skeleton } from "@/components/ui/skeleton"
export default function InstanceEditPage() {
	const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the status dialog
	const [dialogDescription, setDialogDescription] = useState(""); // State to store the dialog description
	const [dialogStatus, setDialogStatus] = useState<"success" | "error" | null>(null);
  const [loading, setLoading] = useState(true);

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			username: "",
		},
	  })
    type InstanceViewResponse = {
      name: string;
      username: string;
    };
    useEffect(() => {
      async function fetchInstance() {
        try {
          const response = await apiClient.get(`/api/users/me`);
          if (response.status === 200) {
            const data = response.data as InstanceViewResponse;
            form.reset({
              name: data.name || "",
              username: data.username || ""
            });
            setLoading(false)
          }
        } catch (error) {
          // Optionally handle error
        }
      }
      fetchInstance();
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    type ProfileResponse = {
      message: string;
      error_description?: string;
    };
	  async function onSubmit(values: z.infer<typeof formSchema>) {
		try {
		  const response = await apiClient.put(`/api/users/me`, values);
      const data = response.data as ProfileResponse;
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
<Link className="" href={`/me`}>
<Button variant="link"><IconChevronLeft stroke={2} />
Go Back to Profile</Button>
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
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="Your username" {...field} />
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
				<Link className="w-1/2" href={`/me`}>
              <Button
                type="button"
                className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none"
                // onClick={...} // Add your view instance handler here
              >
                View Profile
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