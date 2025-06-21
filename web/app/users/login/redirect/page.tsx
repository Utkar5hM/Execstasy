"use client"
import { redirect } from "next/navigation";
import { useAuth } from "@/components/auth";
import { useEffect } from "react";

export default function RedirectToHome() {
  const { refreshUser } = useAuth();

  useEffect(() => {
    refreshUser();
    redirect("/");
  }, [refreshUser]);
  return <div>Redirecting...</div>;
}