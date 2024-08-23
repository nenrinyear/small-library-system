import { NEXT_PUBLIC_SYSTEM_NAME, NEXT_PUBLIC_SYSTEM_DESCRIPTION } from "@/lib/environment-variables";

export default function Header({
    title,
    description,
}: {
    title?: string,
    description?: string,
}) {
    return (
        <header className="w-full min-h-[30vh] flex justify-start items-center bg-slate-400">
            <div className="w-full max-w-[1024px] mx-auto">
                <div className="flex flex-col ml-[5vw]">
                    <h1 className="text-3xl font-bold">{title ?? NEXT_PUBLIC_SYSTEM_NAME}</h1>
                    <p>{description ?? NEXT_PUBLIC_SYSTEM_DESCRIPTION}</p>
                </div>
            </div>
        </header>
    )
}