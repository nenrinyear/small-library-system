import { NEXT_PUBLIC_API_BASE_URL } from "@/lib/environment-variables";

export type Tag = {
    id: number;
    name: string;
};

export type ItemSummary = {
    id: number;
    qrId: string;
    title: string | null;
    publisher: string | null;
    publishedDate: string | null;
    description: string | null;
    rentalStatus: "AVAILABLE" | "RENTED";
    createdAt: string;
    tags: Tag[];
};

export type Rental = {
    id: number;
    itemId: number;
    rentedBy: string;
    rentedAt: string;
    dueDate: string;
    returnedAt: string | null;
    status: "ACTIVE" | "RETURNED";
    createdAt: string;
    updatedAt: string;
};

export type ItemDetail = {
    item: ItemSummary;
    activeRental: Rental | null;
    rentalHistory: Rental[];
};

export type RentedSummary = ItemSummary & {
    rentedBy: string;
    rentedAt: string;
    dueDate: string;
};

export type RentalEvent = {
    rentalId: number;
    itemId: number;
    qrId: string;
    title: string | null;
    rentedBy: string;
    rentedAt: string;
    dueDate: string;
    returnedAt: string | null;
    updatedAt: string;
};

export type DashboardResponse = {
    availableItems: ItemSummary[];
    rentedItems: RentedSummary[];
    recentRentals: RentalEvent[];
};

export type SearchResponse = {
    items: ItemSummary[];
    totalCount: number;
    page: number;
    pageSize: number;
};

type RequestOptions = {
    method?: "GET" | "POST";
    body?: unknown;
    token?: string;
};

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const headers: Record<string, string> = {
        "Content-Type": "application/json",
    };

    if (options.token) {
        headers.Authorization = `Bearer ${options.token}`;
    }

    const response = await fetch(`${NEXT_PUBLIC_API_BASE_URL}${path}`, {
        method: options.method ?? "GET",
        headers,
        body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
        cache: "no-store",
    });

    if (!response.ok) {
        let message = `request failed: ${response.status}`;
        try {
            const payload = (await response.json()) as { error?: string };
            if (payload.error) {
                message = payload.error;
            }
        } catch {
            // noop
        }
        throw new Error(message);
    }

    return (await response.json()) as T;
}

export async function fetchDashboard() {
    return request<DashboardResponse>("/api/dashboard");
}

export async function fetchTags() {
    return request<{ tags: Tag[] }>("/api/tags");
}

export async function fetchItemByQrId(qrId: string) {
    return request<ItemDetail>(`/api/items/${encodeURIComponent(qrId)}`);
}

export async function generateItems(count: number, token: string) {
    return request<{ items: ItemSummary[] }>("/api/items/generate", {
        method: "POST",
        body: { count },
        token,
    });
}

export async function registerItem(qrId: string, body: {
    title: string;
    publisher: string;
    publishedDate?: string;
    description?: string;
    tagIds: number[];
}, token: string) {
    return request<ItemDetail>(`/api/items/${encodeURIComponent(qrId)}/register`, {
        method: "POST",
        body,
        token,
    });
}

export async function rentItem(itemId: number, body: { rentedBy: string; dueDate: string }, token: string) {
    return request<{ status: string }>(`/api/items/${itemId}/rent`, {
        method: "POST",
        body,
        token,
    });
}

export async function returnItem(itemId: number, body: { returnedBy: string }, token: string) {
    return request<{ status: string }>(`/api/items/${itemId}/return`, {
        method: "POST",
        body,
        token,
    });
}

export async function searchItems(params: {
    q?: string;
    publisher?: string;
    tagId?: number;
    page?: number;
    pageSize?: number;
}) {
    const qp = new URLSearchParams();
    if (params.q) qp.set("q", params.q);
    if (params.publisher) qp.set("publisher", params.publisher);
    if (params.tagId) qp.set("tagId", String(params.tagId));
    if (params.page) qp.set("page", String(params.page));
    if (params.pageSize) qp.set("pageSize", String(params.pageSize));

    const query = qp.toString();
    return request<SearchResponse>(`/api/search${query ? `?${query}` : ""}`);
}
