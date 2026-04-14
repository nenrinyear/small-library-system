import type { Metadata } from "next";
import "./globals.css";
import AppProviders from "@/components/AppProviders";
import { NEXT_PUBLIC_SYSTEM_DESCRIPTION, NEXT_PUBLIC_SYSTEM_NAME } from "@/lib/environment-variables";

export const metadata: Metadata = {
    title: NEXT_PUBLIC_SYSTEM_NAME,
    description: NEXT_PUBLIC_SYSTEM_DESCRIPTION,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="ja">
            <body>
                <AppProviders>{children}</AppProviders>
            </body>
        </html>
    );
}
