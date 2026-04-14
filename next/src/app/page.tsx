"use client";

import { fetchDashboard, generateItems, type DashboardResponse, type RentalEvent } from "@/lib/api";
import { useAdminToken } from "@/components/AdminTokenContext";
import {
    Badge,
    Box,
    Button,
    Card,
    Flex,
    Heading,
    HStack,
    Input,
    Link as ChakraLink,
    SimpleGrid,
    Spinner,
    Stack,
    Text,
} from "@chakra-ui/react";
import Link from "next/link";
import { useEffect, useMemo, useState } from "react";

export default function HomePage() {
    const { token } = useAdminToken();
    const [dashboard, setDashboard] = useState<DashboardResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [count, setCount] = useState("1");
    const [generated, setGenerated] = useState<string[]>([]);

    const loadDashboard = async () => {
        setLoading(true);
        setError("");
        try {
            const data = await fetchDashboard();
            setDashboard(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : "failed to load dashboard");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void loadDashboard();
    }, []);

    const overdueRentalIds = useMemo(() => {
        const now = new Date();
        const ids = new Set<number>();
        for (const row of dashboard?.rentedItems ?? []) {
            if (new Date(row.dueDate).getTime() < now.getTime()) {
                ids.add(row.id);
            }
        }
        return ids;
    }, [dashboard]);

    const onGenerate = async () => {
        const parsed = Number(count);
        if (!token) {
            setError("Admin API Token を保存してください");
            return;
        }
        if (!Number.isInteger(parsed) || parsed < 1 || parsed > 100) {
            setError("生成数は 1 から 100 の整数で指定してください");
            return;
        }

        setError("");
        try {
            const response = await generateItems(parsed, token);
            setGenerated(response.items.map((item) => item.qrId));
            await loadDashboard();
        } catch (err) {
            setError(err instanceof Error ? err.message : "failed to generate qr ids");
        }
    };

    return (
        <Stack gap={6}>
            <Heading size="xl">ダッシュボード</Heading>

            <Card.Root>
                <Card.Header>
                    <Heading size="md">QR ID 発行（未登録アイテム）</Heading>
                </Card.Header>
                <Card.Body>
                    <Stack gap={3}>
                        <HStack>
                            <Input
                                value={count}
                                onChange={(event) => setCount(event.target.value)}
                                maxW="160px"
                                type="number"
                                min={1}
                                max={100}
                            />
                            <Button colorScheme="blue" onClick={onGenerate}>発行</Button>
                        </HStack>
                        {generated.length > 0 ? (
                            <Box>
                                <Text fontWeight="bold" mb={2}>生成済み QR ID</Text>
                                <Flex wrap="wrap" gap={2}>
                                    {generated.map((qrId) => (
                                        <Badge key={qrId} colorScheme="purple" px={2} py={1}>
                                            {qrId}
                                        </Badge>
                                    ))}
                                </Flex>
                            </Box>
                        ) : null}
                    </Stack>
                </Card.Body>
            </Card.Root>

            {error ? (
                <Box bg="red.50" borderWidth="1px" borderColor="red.200" rounded="md" p={3}>
                    <Text color="red.700">{error}</Text>
                </Box>
            ) : null}

            {loading ? (
                <HStack>
                    <Spinner size="sm" />
                    <Text>読み込み中...</Text>
                </HStack>
            ) : null}

            <SimpleGrid columns={{ base: 1, lg: 2 }} gap={4}>
                <Card.Root>
                    <Card.Header>
                        <Heading size="md">貸出中</Heading>
                    </Card.Header>
                    <Card.Body>
                        <Stack gap={3}>
                            {(dashboard?.rentedItems ?? []).length === 0 ? (
                                <Text color="gray.600">貸出中のアイテムはありません</Text>
                            ) : (
                                dashboard?.rentedItems.map((item) => (
                                    <Box key={item.id} borderWidth="1px" rounded="md" p={3} bg={overdueRentalIds.has(item.id) ? "orange.50" : "white"}>
                                        <HStack justify="space-between" align="start">
                                            <Box>
                                                <ChakraLink as={Link} href={`/qr/${item.qrId}`} fontWeight="bold" color="blue.600">
                                                    {item.title ?? "(未登録)"}
                                                </ChakraLink>
                                                <Text fontSize="sm" color="gray.600">貸出者: {item.rentedBy}</Text>
                                                <Text fontSize="sm" color="gray.600">返却予定: {formatDate(item.dueDate)}</Text>
                                            </Box>
                                            <Badge colorScheme={overdueRentalIds.has(item.id) ? "red" : "green"}>
                                                {overdueRentalIds.has(item.id) ? "期限超過" : "貸出中"}
                                            </Badge>
                                        </HStack>
                                    </Box>
                                ))
                            )}
                        </Stack>
                    </Card.Body>
                </Card.Root>

                <Card.Root>
                    <Card.Header>
                        <Heading size="md">利用可能</Heading>
                    </Card.Header>
                    <Card.Body>
                        <Stack gap={3}>
                            {(dashboard?.availableItems ?? []).length === 0 ? (
                                <Text color="gray.600">利用可能なアイテムはありません</Text>
                            ) : (
                                dashboard?.availableItems.map((item) => (
                                    <Box key={item.id} borderWidth="1px" rounded="md" p={3}>
                                        <ChakraLink as={Link} href={`/qr/${item.qrId}`} fontWeight="bold" color="blue.600">
                                            {item.title ?? "(未登録)"}
                                        </ChakraLink>
                                        <Text fontSize="sm" color="gray.600">発行ID: {item.qrId}</Text>
                                    </Box>
                                ))
                            )}
                        </Stack>
                    </Card.Body>
                </Card.Root>
            </SimpleGrid>

            <Card.Root>
                <Card.Header>
                    <Heading size="md">最近の貸出履歴</Heading>
                </Card.Header>
                <Card.Body>
                    <Stack gap={3}>
                        {(dashboard?.recentRentals ?? []).length === 0 ? (
                            <Text color="gray.600">履歴はありません</Text>
                        ) : (
                            dashboard?.recentRentals.map((event) => <RentalHistoryCard key={event.rentalId} event={event} />)
                        )}
                    </Stack>
                </Card.Body>
            </Card.Root>
        </Stack>
    );
}

function RentalHistoryCard({ event }: { event: RentalEvent }) {
    return (
        <Box borderWidth="1px" rounded="md" p={3}>
            <HStack justify="space-between" align="start">
                <Box>
                    <ChakraLink as={Link} href={`/qr/${event.qrId}`} color="blue.600" fontWeight="bold">
                        {event.title ?? "(未登録)"}
                    </ChakraLink>
                    <Text fontSize="sm" color="gray.600">貸出者: {event.rentedBy}</Text>
                    <Text fontSize="sm" color="gray.600">
                        貸出: {formatDate(event.rentedAt)} / 返却予定: {formatDate(event.dueDate)}
                    </Text>
                </Box>
                <Badge colorScheme={event.returnedAt ? "gray" : "green"}>
                    {event.returnedAt ? `返却済み (${formatDate(event.returnedAt)})` : "未返却"}
                </Badge>
            </HStack>
        </Box>
    );
}

function formatDate(value: string) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return value;
    }
    return new Intl.DateTimeFormat("ja-JP", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
    }).format(date);
}
