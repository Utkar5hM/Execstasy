"use client";

import { cn } from "@/lib/utils"
import React from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp"

// import { apiClient } from "@/utils/apiClient";

import { useSearchParams } from "next/navigation";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form"

const FormSchema = z.object({
  user_code: z.string().min(8, {
    message: "Your one-time password must be 8 characters.",
  }).regex(/^[A-Z0-9]+$/, { message: "Your one-time password must be alphanumeric." }),
})
export function OAuthForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const searchParams = useSearchParams(); // Access query parameters
  const userCode = searchParams.get("user_code") || ""; // Get the `user_code` query parameter 
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      user_code: userCode ? userCode : "",
    },
  })
  function onSubmit(data: z.infer<typeof FormSchema>) {
    toast("You submitted the following values", {
      description: (
        <pre className="mt-2 w-[320px] rounded-md bg-neutral-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    })
    console.log("Form submitted with data:", data);
  }
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>

<Form {...form}>
        <div className="flex flex-col gap-6">
          <div className="flex flex-col items-center gap-2">
            <h1 className="text-xl font-bold">Instance Access</h1>
            <div className="text-center text-sm">
            Enter the OTP given in your terminal
            </div>
          </div>
            <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col items-center">
            <FormField
          control={form.control}
          name="user_code"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <InputOTP maxLength={8} {...field}
                
          onInput={(e) => {
            const input = e.target as HTMLInputElement;
            input.value = input.value.toUpperCase().replace(/[^A-Z0-9]/g, ""); // Allow only A-Z and 0-9
            field.onChange(input.value); // Update the form state
          }}>
                  <InputOTPGroup>
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                    <InputOTPSlot index={3} />
                    </InputOTPGroup>
                    <InputOTPSeparator />
                    <InputOTPGroup>
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                    <InputOTPSlot index={6} />
                    <InputOTPSlot index={7} />
                  </InputOTPGroup>
                </InputOTP>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
            <Dialog >
            <DialogTrigger asChild>
            <Button type="button" className="w-full mt-6">
              Approve
            </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Are you absolutely sure?</DialogTitle>
                <DialogDescription>
                  This action cannot be undone. If you are sure, click the button below to proceed.
                </DialogDescription>
              </DialogHeader>
      <Button
        type="submit" onClick={form.handleSubmit(onSubmit)}
        className="w-full bg-green-500 text-white hover:bg-green-600 focus:ring-2 focus:ring-green-400 focus:outline-none mt-4"
      >Submit
      </Button>
            </DialogContent>
          </Dialog>

      {/* Status Dialog */}
      <Dialog >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Status</DialogTitle>
            <DialogDescription>Status Message</DialogDescription>
          </DialogHeader>
          <Button
            type="button"
            className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none mt-4"
          >
            Close
          </Button>
        </DialogContent>
      </Dialog>
      </form>
      </div>
</Form>
</div>
  )
}