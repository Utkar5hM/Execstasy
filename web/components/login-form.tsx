import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import pkceChallenge from 'pkce-challenge';
import { useEffect, useState } from "react";
import { Skeleton } from "./ui/skeleton";
import { env } from 'next-runtime-env';


export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {

  const baseURL = env("NEXT_PUBLIC_BACKEND_URL") || "";
  const GitlabEnabled = (env("NEXT_PUBLIC_ENABLE_GITLAB_OAUTH") || "").toLowerCase() === "true";
  const GoogleEnabled = (env("NEXT_PUBLIC_ENABLE_GOOGLE_OAUTH") || "").toLowerCase() === "true";
  const [codeChallenge, setCodeChallenge] = useState("");
  const [loading, setLoading] = useState(true);
  useEffect(() => {

    async function getLoginURL() {
    const { code_challenge, code_verifier } = await pkceChallenge();
    // Store the code_verifier in localStorage or sessionStorage
    sessionStorage.setItem('code_verifier', code_verifier);
    setCodeChallenge(encodeURIComponent(code_challenge));
    setLoading(false);
  }
    getLoginURL();
  }, [])

  if (loading) {
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
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            Login with your Google account
          </CardDescription>
        </CardHeader>
        <CardContent>
        <div className="grid gap-6">
  <div className="flex flex-col gap-4">
    { GoogleEnabled && (
    <form action={baseURL + "/api/users/oauth/google/login?challenge=" + codeChallenge}  method="GET">
    <input type="hidden" name="challenge" value={codeChallenge} />
      <Button type="submit" variant="outline" className="w-full">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path
            d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
            fill="currentColor"
          />
        </svg>
        Login with Google
      </Button>
    </form>
    )}
    { GitlabEnabled && (
    <form action={baseURL + "/api/users/oauth/gitlab/login?challenge=" + codeChallenge} method="GET">
    <input type="hidden" name="challenge" value={codeChallenge} />
      <Button type="submit" variant="outline" className="w-full">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path
            d="M22.548 13.436l-2.09-6.43a1.13 1.13 0 0 0-2.15-.06l-1.7 4.24h-4.216l-1.7-4.24a1.13 1.13 0 0 0-2.15.06l-2.09 6.43a1.13 1.13 0 0 0 .41 1.27l7.09 5.16a1.13 1.13 0 0 0 1.36 0l7.09-5.16a1.13 1.13 0 0 0 .41-1.27z"
            fill="currentColor"
          />
        </svg>
        Login with GitLab
      </Button>
    </form>
    )}
  </div>
</div>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        By clicking continue, you agree to our <a href="#">Terms of Service</a>{" "}
        and <a href="#">Privacy Policy</a>.
      </div>
    </div>
  )
}
