import "./globals.css";
import { PublicEnvScript } from 'next-runtime-env';
import RootLayoutInner from "./root-layout";




export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning={true}>
      <head suppressHydrationWarning={true}>
      <PublicEnvScript />
      <title>ExecStasy</title>
      </head>
      <RootLayoutInner>
        {children}
        </RootLayoutInner>
    </html>
  );
}