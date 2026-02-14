import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Foresee",
  description: "Foresee - Transaction Management",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className="antialiased">{children}</body>
    </html>
  );
}
