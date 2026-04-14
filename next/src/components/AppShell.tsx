"use client";

import { NEXT_PUBLIC_SYSTEM_NAME_SHORT } from "@/lib/environment-variables";
import { Box, Button, Flex, HStack, Input, Link as ChakraLink, Text } from "@chakra-ui/react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { useAdminToken } from "./AdminTokenContext";

export default function AppShell({ children }: { children: React.ReactNode }) {
    const pathname = usePathname();
    const { token, setToken } = useAdminToken();
    const [input, setInput] = useState(token);

    return (
        <Box minH="100vh" bg="gray.50" color="gray.900">
            <Box as="header" borderBottomWidth="1px" bg="white" px={{ base: 4, md: 6 }} py={3} position="sticky" top={0} zIndex={30}>
                <Flex maxW="1200px" mx="auto" align="center" justify="space-between" gap={4} wrap="wrap">
                    <HStack gap={4}>
                        <ChakraLink as={Link} href="/" fontWeight="bold" fontSize="lg">
                            {NEXT_PUBLIC_SYSTEM_NAME_SHORT}
                        </ChakraLink>
                        <HStack gap={3}>
                            <NavLink href="/" active={pathname === "/"}>ホーム</NavLink>
                            <NavLink href="/search" active={pathname?.startsWith("/search") ?? false}>検索</NavLink>
                        </HStack>
                    </HStack>
                    <HStack gap={2} minW={{ base: "100%", md: "auto" }}>
                        <Input
                            type="password"
                            value={input}
                            onChange={(event) => setInput(event.target.value)}
                            placeholder="Admin API Token"
                            size="sm"
                            bg="white"
                        />
                        <Button
                            size="sm"
                            colorScheme="blue"
                            onClick={() => setToken(input.trim())}
                        >
                            保存
                        </Button>
                    </HStack>
                </Flex>
            </Box>
            <Box as="main" maxW="1200px" mx="auto" px={{ base: 4, md: 6 }} py={6}>
                {children}
            </Box>
            <Box as="footer" py={6} textAlign="center" color="gray.600">
                <Text fontSize="sm">Small Library System</Text>
            </Box>
        </Box>
    );
}

function NavLink({ href, active, children }: { href: string; active: boolean; children: React.ReactNode }) {
    return (
        <ChakraLink
            as={Link}
            href={href}
            color={active ? "blue.600" : "gray.700"}
            fontWeight={active ? "bold" : "medium"}
            _hover={{ textDecoration: "none", color: "blue.500" }}
        >
            {children}
        </ChakraLink>
    );
}
