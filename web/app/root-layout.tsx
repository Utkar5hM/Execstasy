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
import Cookies from 'js-cookie';
import { AuthProvider } from "@/components/auth";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});



export default function RootLayoutInner({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname(); // Get the current route
  const [sidebarState, setSidebarState] = useState<boolean>(true);
  const isLoginPage = pathname === "/users/login"; // Check if the current route is `/users/login`
  const isLoginCallback = pathname === "/users/login/callback"; // Check if the current route is a login callback
  const router = useRouter(); // Initialize useRouter for redirection
  const [isLoggedIn, setIsLoggedIn] = useState(false); // State to track login status
  const [loading, setLoading] = useState(true); // State to track loading status
  const [loadingSideBar, setLoadingSideBar] = useState(true); // State to track loading status

  useEffect(() => {
    setSidebarState(Cookies.get("sidebar_state") === "true");
    setLoadingSideBar(false); // Set loading to false after initializing sidebar state
  }, []);

  useEffect(() => {
    // Check for the `jwt` token in cookies
  
    const token = Cookies.get('jwt'); // Get the token from cookies
  
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
    if (!loading && !isLoggedIn && !isLoginPage && !isLoginCallback &&!loadingSideBar) {
      if (window && window.location){
      }
      router.push("/users/login");
    } else if( !loading && isLoggedIn && !loadingSideBar && (isLoginPage || isLoginCallback)) {
      router.push("/"); 
    }
  }, [loading, isLoggedIn, isLoginPage, router, isLoginCallback, loadingSideBar]);
  
  return (
      <AuthProvider>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        suppressHydrationWarning={true}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          {/* Conditionally render SidebarProvider */}
          {!isLoginPage && !loading && !loadingSideBar ? (
            <SidebarProvider
              defaultOpen={sidebarState}
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
      </body></AuthProvider>
  );
}