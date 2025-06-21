"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Cookies from 'js-cookie';

export default function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    // Clear the JWT cookie
    Cookies.remove('jwt');
    // redirect with timeout
    setTimeout(() => {
      router.push("/users/login");
    }, 500); 

  }, [router]);

  return null; // Render nothing since this page is just for logout
}