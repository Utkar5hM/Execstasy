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
	InputOTPSlot,
  } from "@/components/ui/input-otp"
export function OAuthForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
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
            </div>
            </div>
            <Dialog>
            <DialogTrigger asChild>
            <Button type="button" className="w-full">
              Approve
            </Button></DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Are you absolutely sure?</DialogTitle>
                <DialogDescription>
                  This action cannot be undone. This will permanently delete your account
                  and remove your data from our servers.
                </DialogDescription>
              </DialogHeader>
            <Button type="submit" className="w-full">
              Sure
            </Button>
            </DialogContent>
          </Dialog>
          </div>
        </div>
      </form>
    </div>
  )
}
