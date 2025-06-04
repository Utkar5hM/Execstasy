import { promises as fs } from "fs"
import path from "path"
import { Metadata } from "next"
import { z } from "zod"
import {OAuthForm} from "@/components/oauth"

export const metadata: Metadata = {
  title: "Approve Access - Instance",
  description: "Approve access to your Instance",
}


export default async function TaskPage() {
  return (
    <div className="bg-background flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="w-full max-w-sm">
        <OAuthForm />
      </div>
    </div>
  )
}