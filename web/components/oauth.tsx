"use client";

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
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
  } from "@/components/ui/input-otp"


import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { apiClient } from "@/utils/apiClient";


import React from "react";

interface InputOTPSlotProps extends React.InputHTMLAttributes<HTMLInputElement> {
  index: number;
}

export const InputOTPSlot: React.FC<InputOTPSlotProps> = ({ index, ...props }) => {
  return (
    <input
      {...props} // Forward all props (e.g., value, onChange)
      type="text"
      maxLength={1} // Limit to one character per slot
      className="w-10 h-10 text-center border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
      data-otp-slot={index} // Add a custom attribute for identification
    />
  );
};


export function OAuthForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const searchParams = useSearchParams(); // Access query parameters
  const otp = searchParams.get("user_code"); // Get the `user_code` query parameter
  const [otpValues, setOtpValues] = useState<string[]>(Array(8).fill("")); // State to manage OTP values
  const [loading, setLoading] = useState(false); // State to manage loading
  const [error, setError] = useState<string | null>(null); // State to manage errors
  const [statusDialogOpen, setStatusDialogOpen] = useState(false); // State to control the second dialog
  const [statusMessage, setStatusMessage] = useState<string>(""); // State to store the status message
  const [confirmationDialogOpen, setConfirmationDialogOpen] = useState(false); // State to control the confirmation dialog
  useEffect(() => {
    if (otp && otp.length === 9 && otp.includes("-")) {
      // Remove the dash and split the OTP into individual characters
      const otpArray = otp.replace("-", "").split("");
      setOtpValues(otpArray); // Update the state with the OTP values
    }
  }, [otp]);

  const handleChange = (index: number, value: string) => {
    const newOtpValues = [...otpValues];
    newOtpValues[index] = value;
    setOtpValues(newOtpValues);
  };

  const handleSubmit = async () => {
    setLoading(true);
    setError(null);

    try {
      // Combine OTP values into a single string with a dash
      const userCode = `${otpValues.slice(0, 4).join("")}-${otpValues.slice(4, 8).join("")}`;

      // Make the API call
      const response = await apiClient.post("/api/oauth", {
        user_code: userCode,
      });

      console.log("Response:", response);
      setStatusMessage("Access approved successfully!"); // Success message
    } catch (err: any) {
      console.error(err);
      setStatusMessage(err.message || "An unexpected error occurred"); // Error message
    } finally {
      setLoading(false);
      setStatusDialogOpen(true); // Open the status dialog
      setConfirmationDialogOpen(false); // Close the confirmation dialog

    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <form>
        <div className="flex flex-col gap-6">
          <div className="flex flex-col items-center gap-2">
            <h1 className="text-xl font-bold">Instance Access</h1>
            <div className="text-center text-sm">
            Enter the OTP given in your terminal
            </div>
          </div>
          <div className="flex flex-col gap-6">
            <div className="grid gap-3">
            <div className="flex flex-col items-center justify-center">
            <InputOTP maxLength={8} className="flex flex-col items-center gap-4">
                <InputOTPGroup>
                  {otpValues.slice(0, 4).map((value, index) => (
                    <InputOTPSlot
                      key={index}
                      index={index}
                      value={value} // Controlled value
                      onChange={(e) =>
                        handleChange(index, (e.target as HTMLInputElement).value) // Cast e.target
                      }
                    />
                  ))}
                </InputOTPGroup>
                <InputOTPSeparator />
                <InputOTPGroup>
                  {otpValues.slice(4, 8).map((value, index) => (
                    <InputOTPSlot
                          key={index + 4}
                          index={index + 4}
                          value={value} // Controlled value
                          onChange={(e) =>
                            handleChange(index + 4, (e.target as HTMLInputElement).value) // Cast e.target
                          }
                          />
                  ))}
                </InputOTPGroup>
                </InputOTP>
            </div>
            </div>
            <Dialog open={confirmationDialogOpen} onOpenChange={setConfirmationDialogOpen}>
            <DialogTrigger asChild>
            <Button type="button" className="w-full">
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
        type="button"
        className="w-full bg-green-500 text-white hover:bg-green-600 focus:ring-2 focus:ring-green-400 focus:outline-none mt-4"
        onClick={handleSubmit}
        disabled={loading}
      >
        {loading ? "Processing..." : "Sure"}
      </Button>
            </DialogContent>
          </Dialog>
                {/* Status Dialog */}
      <Dialog open={statusDialogOpen} onOpenChange={setStatusDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Status</DialogTitle>
            <DialogDescription>{statusMessage}</DialogDescription>
          </DialogHeader>
          <Button
            type="button"
            className="w-full bg-blue-500 text-white hover:bg-blue-600 focus:ring-2 focus:ring-blue-400 focus:outline-none mt-4"
            onClick={() => setStatusDialogOpen(false)}
          >
            Close
          </Button>
        </DialogContent>
      </Dialog>
          </div>
        </div>
      </form>
    </div>
  )
}
