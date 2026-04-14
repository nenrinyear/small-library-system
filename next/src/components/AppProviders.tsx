"use client";

import { ChakraProvider, defaultSystem } from "@chakra-ui/react";
import AppShell from "@/components/AppShell";
import { AdminTokenProvider } from "@/components/AdminTokenContext";

export default function AppProviders({ children }: { children: React.ReactNode }) {
    return (
        <ChakraProvider value={defaultSystem}>
            <AdminTokenProvider>
                <AppShell>{children}</AppShell>
            </AdminTokenProvider>
        </ChakraProvider>
    );
}
