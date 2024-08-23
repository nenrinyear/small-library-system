import type { Metadata } from "next";
import { IBM_Plex_Sans_JP } from "next/font/google";
import "./globals.css";
import Providers from "@/components/NextUIProvider";

const ibm_plex_sans = IBM_Plex_Sans_JP({
    display: "swap",
    subsets: ["latin-ext"],
    weight: "400",
    style: "normal",
});

export const metadata: Metadata = {
    title: {
        template: "%s | SmallLibrarySystem",
        default: "SmallLibrarySystem",
    },
    description: process.env.NEXT_PUBLIC_SYSTEM_DESCRIPTION ?? "QRコードを用いて物品を管理するシステムです。",
    openGraph: {
        locale: "ja_JP",
        type: "website",
        title: {
            template: "%s | SmallLibrarySystem",
            default: "SmallLibrarySystem",
        },
        description: process.env.NEXT_PUBLIC_SYSTEM_DESCRIPTION ?? "QRコードを用いて物品を管理するシステムです。",
        url: process.env.NEXT_PUBLIC_SYSTEM_HOST ?? "http://localhost:3000",
        images: [
            {
                url: `${process.env.NEXT_PUBLIC_SYSTEM_HOST}/ogp.png`,
                width: 1200,
                height: 630,
                alt: process.env.NEXT_PUBLIC_SYSTEM_NAME ?? "SmallLibrarySystem",
            },
        ],
    },
    twitter: {
        site: process.env.NEXT_PUBLIC_TWITTER_ACCOUNT ?? "@twitter",
        card: "summary_large_image",
        images: [
            {
                url: `${process.env.NEXT_PUBLIC_SYSTEM_HOST}/ogp.png`,
                width: 1200,
                height: 630,
                alt: process.env.NEXT_PUBLIC_SYSTEM_NAME ?? "SmallLibrarySystem",
            },
        ],
    },
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="ja" className="light">
            <body className={ibm_plex_sans.className}>
                <Providers>
                    {children}
                </Providers>
            </body>
        </html>
    );
}
