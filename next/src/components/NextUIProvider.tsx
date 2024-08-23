"use client";
import { NextUIProvider } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import HeaderNavigation from "./Navbar";
import { SessionProvider } from "next-auth/react";

export default function Providers({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const router = useRouter();

    return (
        <NextUIProvider navigate={router.push}>
            <SessionProvider>
                <HeaderNavigation />
                {children}
            </SessionProvider>
        </NextUIProvider>
    );
}