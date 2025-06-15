"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    // Clear the JWT cookie
    document.cookie = "jwt=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;";

    // Redirect to the login page
    router.push("/users/login");
  }, []);

  return null; // Render nothing since this page is just for logout
}