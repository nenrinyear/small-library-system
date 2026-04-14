"use client";

import { useAdminToken } from "@/components/AdminTokenContext";
import {
    fetchItemByQrId,
    fetchTags,
    registerItem,
    rentItem,
    returnItem,
    type ItemDetail,
    type Tag,
} from "@/lib/api";
import {
    Badge,
    Box,
    Button,
    Card,
    Heading,
    HStack,
    Input,
    Spinner,
    Stack,
    Text,
    Textarea,
} from "@chakra-ui/react";
import { useParams } from "next/navigation";
import { FormEvent, useEffect, useMemo, useState } from "react";

export default function QrDetailPage() {
    const { token } = useAdminToken();
    const params = useParams<{ qrId: string }>();
    const qrId = params?.qrId ?? "";

    const [itemDetail, setItemDetail] = useState<ItemDetail | null>(null);
    const [tags, setTags] = useState<Tag[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    const [title, setTitle] = useState("");
    const [publisher, setPublisher] = useState("");
    const [publishedDate, setPublishedDate] = useState("");
    const [description, setDescription] = useState("");
    const [selectedTagIds, setSelectedTagIds] = useState<number[]>([]);

    const [rentedBy, setRentedBy] = useState("");
    const [dueDate, setDueDate] = useState(defaultDueDate());
    const [returnedBy, setReturnedBy] = useState("");

    const isRegistered = useMemo(() => {
        return !!itemDetail?.item.title;
    }, [itemDetail]);

    const load = async () => {
        setLoading(true);
        setError("");
        try {
            const [itemPayload, tagPayload] = await Promise.all([
                fetchItemByQrId(qrId),
                fetchTags(),
            ]);
            setItemDetail(itemPayload);
            setTags(tagPayload.tags);
        } catch (err) {
            setError(err instanceof Error ? err.message : "読み込みに失敗しました");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!qrId) {
            return;
        }
        void load();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [qrId]);

    const handleRegister = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!token) {
            setError("Admin API Token を保存してください");
            return;
        }

        if (selectedTagIds.length === 0) {
            setError("タグを1つ以上選択してください");
            return;
        }

        setError("");
        try {
            const payload = await registerItem(
                qrId,
                {
                    title,
                    publisher,
                    publishedDate: publishedDate || undefined,
                    description: description || undefined,
                    tagIds: selectedTagIds,
                },
                token,
            );
            setItemDetail(payload);
        } catch (err) {
            setError(err instanceof Error ? err.message : "登録に失敗しました");
        }
    };

    const handleRent = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!token || !itemDetail) {
            setError("Admin API Token を保存してください");
            return;
        }

        setError("");
        try {
            await rentItem(itemDetail.item.id, { rentedBy, dueDate }, token);
            setRentedBy("");
            setDueDate(defaultDueDate());
            await load();
        } catch (err) {
            setError(err instanceof Error ? err.message : "貸出に失敗しました");
        }
    };

    const handleReturn = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!token || !itemDetail) {
            setError("Admin API Token を保存してください");
            return;
        }

        setError("");
        try {
            await returnItem(itemDetail.item.id, { returnedBy }, token);
            setReturnedBy("");
            await load();
        } catch (err) {
            setError(err instanceof Error ? err.message : "返却に失敗しました");
        }
    };

    return (
        <Stack gap={6}>
            <Heading size="xl">QR詳細: {qrId}</Heading>

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

            {!loading && itemDetail ? (
                <>
                    <Card.Root>
                        <Card.Header>
                            <HStack justify="space-between">
                                <Heading size="md">基本情報</Heading>
                                <Badge colorScheme={itemDetail.item.rentalStatus === "RENTED" ? "orange" : "green"}>
                                    {itemDetail.item.rentalStatus === "RENTED" ? "貸出中" : "利用可能"}
                                </Badge>
                            </HStack>
                        </Card.Header>
                        <Card.Body>
                            <Stack gap={2}>
                                <Text>タイトル: {itemDetail.item.title ?? "未登録"}</Text>
                                <Text>Publisher: {itemDetail.item.publisher ?? "-"}</Text>
                                <Text>公開日: {itemDetail.item.publishedDate ? formatDate(itemDetail.item.publishedDate) : "-"}</Text>
                                <Text>説明: {itemDetail.item.description ?? "-"}</Text>
                                <HStack wrap="wrap" gap={2}>
                                    {itemDetail.item.tags.map((tag) => (
                                        <Badge key={tag.id} colorScheme="purple">{tag.name}</Badge>
                                    ))}
                                </HStack>
                            </Stack>
                        </Card.Body>
                    </Card.Root>

                    {!isRegistered ? (
                        <Card.Root>
                            <Card.Header>
                                <Heading size="md">未登録アイテムの登録</Heading>
                            </Card.Header>
                            <Card.Body>
                                <form onSubmit={handleRegister}>
                                    <Stack gap={3}>
                                        <Input placeholder="タイトル" value={title} onChange={(event) => setTitle(event.target.value)} required />
                                        <Input placeholder="Publisher" value={publisher} onChange={(event) => setPublisher(event.target.value)} required />
                                        <Input type="month" value={publishedDate} onChange={(event) => setPublishedDate(event.target.value)} required />
                                        <Textarea placeholder="説明" value={description} onChange={(event) => setDescription(event.target.value)} />
                                        <Box>
                                            <Text fontWeight="bold" mb={2}>タグ（1つ以上）</Text>
                                            <Stack gap={1}>
                                                {tags.map((tag) => (
                                                    <label key={tag.id}>
                                                        <input
                                                            type="checkbox"
                                                            checked={selectedTagIds.includes(tag.id)}
                                                            onChange={(event) => {
                                                                if (event.target.checked) {
                                                                    setSelectedTagIds((prev) => [...prev, tag.id]);
                                                                } else {
                                                                    setSelectedTagIds((prev) => prev.filter((id) => id !== tag.id));
                                                                }
                                                            }}
                                                        />{" "}
                                                        {tag.name}
                                                    </label>
                                                ))}
                                            </Stack>
                                        </Box>
                                        <Button type="submit" colorScheme="blue">登録</Button>
                                    </Stack>
                                </form>
                            </Card.Body>
                        </Card.Root>
                    ) : (
                        <>
                            {itemDetail.item.rentalStatus === "AVAILABLE" ? (
                                <Card.Root>
                                    <Card.Header>
                                        <Heading size="md">貸出</Heading>
                                    </Card.Header>
                                    <Card.Body>
                                        <form onSubmit={handleRent}>
                                            <Stack gap={3}>
                                                <Input placeholder="貸出者" value={rentedBy} onChange={(event) => setRentedBy(event.target.value)} required />
                                                <Input type="date" value={dueDate} onChange={(event) => setDueDate(event.target.value)} required />
                                                <Button type="submit" colorScheme="blue">貸出を確定</Button>
                                            </Stack>
                                        </form>
                                    </Card.Body>
                                </Card.Root>
                            ) : (
                                <Card.Root>
                                    <Card.Header>
                                        <Heading size="md">返却</Heading>
                                    </Card.Header>
                                    <Card.Body>
                                        <Stack gap={3}>
                                            {itemDetail.activeRental ? (
                                                <Box bg="orange.50" p={3} rounded="md" borderWidth="1px" borderColor="orange.200">
                                                    <Text>貸出者: {itemDetail.activeRental.rentedBy}</Text>
                                                    <Text>貸出日: {formatDate(itemDetail.activeRental.rentedAt)}</Text>
                                                    <Text>返却予定: {formatDate(itemDetail.activeRental.dueDate)}</Text>
                                                </Box>
                                            ) : null}
                                            <form onSubmit={handleReturn}>
                                                <Stack gap={3}>
                                                    <Input placeholder="返却者（貸出者と同一）" value={returnedBy} onChange={(event) => setReturnedBy(event.target.value)} required />
                                                    <Button type="submit" colorScheme="orange">返却を確定</Button>
                                                </Stack>
                                            </form>
                                        </Stack>
                                    </Card.Body>
                                </Card.Root>
                            )}
                        </>
                    )}

                    <Card.Root>
                        <Card.Header>
                            <Heading size="md">貸出履歴</Heading>
                        </Card.Header>
                        <Card.Body>
                            <Stack gap={2}>
                                {itemDetail.rentalHistory.length === 0 ? (
                                    <Text color="gray.600">履歴はありません</Text>
                                ) : (
                                    itemDetail.rentalHistory.map((row) => (
                                        <Box key={row.id} borderWidth="1px" rounded="md" p={3}>
                                            <HStack justify="space-between" align="start">
                                                <Box>
                                                    <Text fontWeight="bold">{row.rentedBy}</Text>
                                                    <Text fontSize="sm" color="gray.600">貸出: {formatDate(row.rentedAt)}</Text>
                                                    <Text fontSize="sm" color="gray.600">返却予定: {formatDate(row.dueDate)}</Text>
                                                </Box>
                                                <Badge colorScheme={row.returnedAt ? "gray" : "green"}>
                                                    {row.returnedAt ? `返却済み (${formatDate(row.returnedAt)})` : "未返却"}
                                                </Badge>
                                            </HStack>
                                        </Box>
                                    ))
                                )}
                            </Stack>
                        </Card.Body>
                    </Card.Root>
                </>
            ) : null}
        </Stack>
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

function defaultDueDate() {
    const date = new Date();
    date.setDate(date.getDate() + 7);
    return date.toISOString().slice(0, 10);
}
