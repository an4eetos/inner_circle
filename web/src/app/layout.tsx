import type { Metadata } from "next";
import { Geist } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL(process.env.SITE_URL ?? "http://localhost:3000"),
  title: {
    default: "Inner Circle: see inside, bring your inner circle",
    template: "%s · Inner Circle",
  },
  description:
    "A social health, biohacking and longevity app. Your body mapped in 3D, a program for every organ you are working on, and the friends you choose keeping you on track.",
  openGraph: {
    title: "Inner Circle",
    description: "See inside. Bring your inner circle.",
    images: ["/screens/00-overview.png"],
  },
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={`${geistSans.variable} h-full antialiased`}>
      <body className="flex min-h-full flex-col font-sans">{children}</body>
    </html>
  );
}
