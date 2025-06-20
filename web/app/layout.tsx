import { cookies } from "next/headers";
import RootLayoutClient from "./root-layout";

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const cookieStore = await cookies();
  const sidebarState = cookieStore.get("sidebar_state")?.value === "true";

  return (
    <RootLayoutClient defaultOpen={sidebarState}>
      {children}
    </RootLayoutClient>
  );
}