"use client";
import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Skeleton } from "@/components/ui/skeleton";
import { apiClient } from "@/utils/apiClient";
import { AuthLogin } from "@/utils/ResponseTypes";
import Cookies from 'js-cookie';
import { Suspense } from "react";

function OAuthCallbackInner() {
  const searchParams = useSearchParams();
  const router = useRouter();

  useEffect(() => {
    const checkCookieAndRedirect = (token: string) => {
      if (Cookies.get('jwt')== token) {
        setTimeout(()=>{
          router.replace("/users/login/redirect");
        }, 500);
      } else {
        setTimeout(checkCookieAndRedirect, 100); // check again in 100ms
      }
    };
    async function handleOAuthCallback() {
      const code = searchParams.get("code");
      const state = searchParams.get("state");
      const provider = searchParams.get("provider");
      const codeVerifier = sessionStorage.getItem("code_verifier");
      if (code && state && codeVerifier && provider) {
        apiClient
          .post<AuthLogin>(
            `/api/users/oauth/${provider}/exchange`,
            {
              code,
              state,
              code_verifier: codeVerifier,
            }
          )
          .then(async (res) => {
            if (res.ok) {
              Cookies.set('jwt', res.data.access_token || "", { 
                expires: new Date(res.data.expiry || 0),
                sameSite: "lax",
                secure: true,
                path: "/"
              });
              checkCookieAndRedirect(res.data.access_token || "invalid");
            } else {
              return router.replace("/users/login");
            }
          })
          .catch(() => {
            return router.replace("/users/login");
          })
      } else {
        return router.replace("/users/login?error=missing_params");
      }
  }
  handleOAuthCallback() 
  }, [router, searchParams]);

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

export default function OAuthCallback() {
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
      <OAuthCallbackInner />
    </Suspense>
  );
}