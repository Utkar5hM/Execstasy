import { Metadata } from "next"
import {OAuthForm} from "@/components/oauth"
import { Suspense } from "react";
import { Skeleton } from "@/components/ui/skeleton";

export const metadata: Metadata = {
  title: "Approve Access - Instance",
  description: "Approve access to your Instance",
}


function OauthPageInner() {
  return (
    <div className="bg-background flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="w-full max-w-sm">
        <OAuthForm />
      </div>
    </div>
  )
}

export default function OAuthPage() {
  return (
    <Suspense fallback={
      <div className="flex h-screen items-center justify-center">
        <div className="flex items-center space-x-4">
          <Skeleton className="h-12 w-12 rounded-full" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-[250px]" />
            <Skeleton className="h-4 w-[200px]" />
          </div>
        </div>
      </div>
    }>
      <OauthPageInner />
    </Suspense>
  );
}