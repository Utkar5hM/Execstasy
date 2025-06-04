"use client";

import { ThemeProvider } from "@/components/theme-provider";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AppSidebar } from "@/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { usePathname, useRouter } from "next/navigation";
import { SiteHeader } from "@/components/site-header"
import { useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});



export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname(); // Get the current route
  const isLoginPage = pathname === "/users/login"; // Check if the current route is `/users/login`
  const router = useRouter(); // Initialize useRouter for redirection
  const [isLoggedIn, setIsLoggedIn] = useState(false); // State to track login status
  const [loading, setLoading] = useState(true); // State to track loading status

  useEffect(() => {
    // Check for the `jwt` token in cookies
    const cookies = document.cookie.split(";").reduce((acc, cookie) => {
      const [key, value] = cookie.trim().split("=");
      acc[key] = value;
      return acc;
    }, {} as Record<string, string>);

    const token = cookies.jwt;

    if (token) {
      try {
        // Decode the token
        const decoded: { exp: number } = jwtDecode(token);

        // Check if the token is expired
        const isExpired = decoded.exp * 1000 < Date.now();
        if (!isExpired) {
          setIsLoggedIn(true); // Token is valid and not expired
        } else {
          console.warn("Token is expired");
          setIsLoggedIn(false);
        }
      } catch (error) {
        console.error("Invalid token:", error);
        setIsLoggedIn(false);
      }
    } else {
      setIsLoggedIn(false); // No token found
    }

    setLoading(false); // Validation is complete
  }, []);

  useEffect(() => {
    // Redirect to login page if not logged in and not already on the login page
    if (!loading && !isLoggedIn && !isLoginPage) {
      router.push("/users/login");
    }
  }, [loading, isLoggedIn, isLoginPage, router]);
  
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {/* Conditionally render SidebarProvider */}
          {!isLoginPage ? (
            <SidebarProvider
              style={
                {
                  "--sidebar-width": "calc(var(--spacing) * 72)",
                  "--header-height": "calc(var(--spacing) * 12)",
                } as React.CSSProperties
              }
            >
              <AppSidebar variant="inset" />
              <SidebarInset>
                  <SiteHeader />
                <div className="flex flex-1 flex-col">
                  <div className="@container/main flex flex-1 flex-col gap-2">
                    {children}
                  </div>
                </div>
        </SidebarInset>
            </SidebarProvider>
          ) : (
            // Render children directly for login page
            <>{children}</>
          )}
        </ThemeProvider>
      </body>
    </html>
  );
}