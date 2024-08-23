import Header from "@/components/Header";
import { NEXT_PUBLIC_MAX_HISTORY_COUNT, NEXT_PUBLIC_SYSTEM_MANAGE_ITEM } from "@/lib/environment-variables";

export default function Page() {
    return (
        <main className="w-full min-h-[calc(100vh-4rem)]">
            <Header />

            {/* history */}
            <section className="w-full flex justify-start items-center">
                <div className="w-full p-4">
                    <h2 className="text-2xl font-bold">最新の貸出履歴(最大{NEXT_PUBLIC_MAX_HISTORY_COUNT}件)</h2>
                    {/* TODO: 貸出履歴実装 */}
                    <div className="w-full p-4">
                        <p>貸出履歴はありません。</p>
                    </div>
                </div>
            </section>

            {/* items */}
            <section className="w-full flex justify-start items-center">
                <div className="w-full p-4">
                    <h2 className="text-2xl font-bold">{NEXT_PUBLIC_SYSTEM_MANAGE_ITEM}一覧</h2>
                    {/* TODO: item一覧実装 */}
                    <div className="w-full p-4">
                        <p>{NEXT_PUBLIC_SYSTEM_MANAGE_ITEM}はありません。</p>
                    </div>
                </div>
            </section>
        </main>
    )
}