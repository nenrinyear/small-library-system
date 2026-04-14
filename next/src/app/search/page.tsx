"use client";

import { fetchTags, searchItems, type ItemSummary, type Tag } from "@/lib/api";
import {
    Box,
    Button,
    Card,
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

export default function SearchPage() {
    const [tags, setTags] = useState<Tag[]>([]);
    const [query, setQuery] = useState("");
    const [publisher, setPublisher] = useState("");
    const [tagId, setTagId] = useState("");
    const [page, setPage] = useState(1);
    const [pageSize] = useState(20);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [items, setItems] = useState<ItemSummary[]>([]);
    const [totalCount, setTotalCount] = useState(0);

    const totalPages = useMemo(() => Math.max(1, Math.ceil(totalCount / pageSize)), [totalCount, pageSize]);

    const loadTags = async () => {
        try {
            const payload = await fetchTags();
            setTags(payload.tags);
        } catch {
            // 検索本体は続行可能なので握りつぶす
        }
    };

    const loadItems = async (nextPage: number) => {
        setLoading(true);
        setError("");
        try {
            const payload = await searchItems({
                q: query || undefined,
                publisher: publisher || undefined,
                tagId: tagId ? Number(tagId) : undefined,
                page: nextPage,
                pageSize,
            });
            setItems(payload.items);
            setTotalCount(payload.totalCount);
            setPage(payload.page);
        } catch (err) {
            setError(err instanceof Error ? err.message : "検索に失敗しました");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void loadTags();
        void loadItems(1);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <Stack gap={6}>
            <Heading size="xl">検索</Heading>

            <Card.Root>
                <Card.Header>
                    <Heading size="md">条件</Heading>
                </Card.Header>
                <Card.Body>
                    <SimpleGrid columns={{ base: 1, md: 2 }} gap={3}>
                        <Input
                            placeholder="キーワード（タイトル・説明）"
                            value={query}
                            onChange={(event) => setQuery(event.target.value)}
                        />
                        <Input
                            placeholder="Publisher"
                            value={publisher}
                            onChange={(event) => setPublisher(event.target.value)}
                        />
                        <select
                            value={tagId}
                            onChange={(event) => setTagId(event.target.value)}
                            style={{
                                width: "100%",
                                height: "2.5rem",
                                borderRadius: "0.375rem",
                                border: "1px solid #CBD5E0",
                                padding: "0 0.75rem",
                                background: "#fff",
                            }}
                        >
                            <option value="">タグ指定なし</option>
                            {tags.map((tag) => (
                                <option key={tag.id} value={tag.id}>{tag.name}</option>
                            ))}
                        </select>
                    </SimpleGrid>
                    <HStack mt={4}>
                        <Button colorScheme="blue" onClick={() => void loadItems(1)}>
                            検索
                        </Button>
                        <Text color="gray.600" fontSize="sm">{totalCount} 件</Text>
                    </HStack>
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
                    <Text>検索中...</Text>
                </HStack>
            ) : null}

            <Stack gap={3}>
                {items.map((item) => (
                    <Box key={item.id} borderWidth="1px" rounded="md" p={3}>
                        <ChakraLink as={Link} href={`/qr/${item.qrId}`} color="blue.600" fontWeight="bold">
                            {item.title ?? "(未登録)"}
                        </ChakraLink>
                        <Text color="gray.600" fontSize="sm">Publisher: {item.publisher ?? "-"}</Text>
                        <Text mt={2}>{item.description ?? ""}</Text>
                    </Box>
                ))}
                {!loading && items.length === 0 ? (
                    <Text color="gray.600">該当データがありません。</Text>
                ) : null}
            </Stack>

            <HStack justify="space-between">
                <Button onClick={() => void loadItems(Math.max(1, page - 1))} disabled={page <= 1 || loading}>
                    前へ
                </Button>
                <Text>
                    {page} / {totalPages}
                </Text>
                <Button onClick={() => void loadItems(Math.min(totalPages, page + 1))} disabled={page >= totalPages || loading}>
                    次へ
                </Button>
            </HStack>
        </Stack>
    );
}
