import { NEXT_PUBLIC_SYSTEM_DESCRIPTION, NEXT_PUBLIC_SYSTEM_NAME } from "@/lib/environment-variables";

export default function Page() {
    return (
        <main className="w-full min-h-[calc(100vh-4rem)]">
            <header className="w-full min-h-[30vh] flex justify-start items-center bg-slate-400">
                <div className="flex flex-col ml-[5vw]">
                    <h1 className="text-3xl font-bold">{NEXT_PUBLIC_SYSTEM_NAME}</h1>
                    <p>{NEXT_PUBLIC_SYSTEM_DESCRIPTION}</p>
                </div>
            </header>

            
        </main>
    )
}