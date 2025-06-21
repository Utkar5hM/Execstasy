import RootLayoutClient from "./root-layout";
import { AuthProvider } from "@/components/auth";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {


  return (
    <AuthProvider>
      <RootLayoutClient>
        {children}
      </RootLayoutClient>
    </AuthProvider>
  );
}