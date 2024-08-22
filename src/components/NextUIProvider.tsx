"use client";
import { NextUIProvider } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import HeaderNavigation from "./Navbar";

export default function UIProvider({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const router = useRouter();

    return (
        <NextUIProvider navigate={router.push}>
            <HeaderNavigation />
            {children}
        </NextUIProvider>
    );
}