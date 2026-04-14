"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";

const STORAGE_KEY = "sls_admin_token";

type AdminTokenContextValue = {
    token: string;
    setToken: (value: string) => void;
};

const AdminTokenContext = createContext<AdminTokenContextValue | null>(null);

export function AdminTokenProvider({ children }: { children: React.ReactNode }) {
    const [token, setTokenState] = useState("");

    useEffect(() => {
        const saved = globalThis.localStorage?.getItem(STORAGE_KEY) ?? "";
        if (saved) {
            setTokenState(saved);
        }
    }, []);

    const setToken = (value: string) => {
        setTokenState(value);
        globalThis.localStorage?.setItem(STORAGE_KEY, value);
    };

    const contextValue = useMemo(
        () => ({ token, setToken }),
        [token],
    );

    return (
        <AdminTokenContext.Provider value={contextValue}>
            {children}
        </AdminTokenContext.Provider>
    );
}

export function useAdminToken() {
    const context = useContext(AdminTokenContext);
    if (!context) {
        throw new Error("useAdminToken must be used inside AdminTokenProvider");
    }
    return context;
}
