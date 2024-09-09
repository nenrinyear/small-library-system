import Header from "@/components/Header";
import { Button } from "@nextui-org/react";
import Link from "next/link";

export default function ManageDashboard() {
    return (
        <main className="w-full min-h-[calc(100vh-4rem)]">
            <Header title="管理画面" description="蔵書の管理を行います。" />

            <section className="w-full flex justify-start items-center">
                <div className="w-full max-w-[1024px] mx-auto p-4">
                    <div className="w-full flex justify-center items-center gap-4">
                        <Button as={Link} color="primary" href="/manage/books" variant="flat">蔵書管理</Button>
                        <Button as={Link} color="primary" href="/manage/tags" variant="flat">タグ管理</Button>
                    </div>
                </div>
            </section>
        </main>
    )
}