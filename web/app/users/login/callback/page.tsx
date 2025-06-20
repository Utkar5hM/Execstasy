"use client";
import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Skeleton } from "@/components/ui/skeleton";
import { apiClient } from "@/utils/apiClient";
import { AuthLogin } from "@/utils/ResponseTypes";

export default function OAuthCallback() {
  const searchParams = useSearchParams();
  const router = useRouter();

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    const codeVerifier = sessionStorage.getItem("code_verifier");
    if (code && state && codeVerifier) {
      apiClient
        .post<AuthLogin>(
          "/api/users/oauth/gitlab/exchange",
          {
            code,
            state,
            code_verifier: codeVerifier,
          }
        )
        .then((res) => {
          if (res.ok) {
            document.cookie = `jwt=${res.data.access_token}; path=/; max-age=${res.data.expiry}; Secure;`;
            setTimeout(() => {
              router.push("/");
            }, 1000);
          } else {
            router.replace("/users/login");
          }
        })
        .catch(() => {
          router.replace("/users/login");
        }).finally(() => {
          sessionStorage.removeItem("code_verifier");
    })
    } else {
      router.replace("/login?error=missing_params");
    }
  }, [router]);

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