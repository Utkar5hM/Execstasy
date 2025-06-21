"use client"
import { redirect } from "next/navigation";
import { useAuth } from "@/components/auth";
import { useEffect } from "react";
import { Skeleton } from "@/components/ui/skeleton";

export default function RedirectToHome() {
  const { refreshUser } = useAuth();

  useEffect(() => {
    refreshUser();
    redirect("/");
  }, [refreshUser]);
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