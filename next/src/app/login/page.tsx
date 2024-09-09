import Header from "@/components/Header";
import { auth } from "@/lib/auth";
import { signIn } from "@/lib/auth";
import { Button, Input } from "@nextui-org/react";
import { redirect } from "next/navigation";

export default async function SignInPage({ searchParams }: Readonly<{ searchParams: { [key: string]: string | string[] | undefined } }>) {
    const session = await auth();

    if (session?.user?.id) {
        redirect("/manage");
    }

    const { message, error } = searchParams;

    return (
        <main className="w-full min-h-[calc(100vh-4rem)]">
            <Header title="ログイン" description="管理者パスワードを入力してください。" />

            {message && <div className="w-full p-4 bg-green-200">{message}</div>}
            {error && <div className="w-full p-4 bg-red-200">{error}</div>}

            <section className="w-full flex justify-start items-center">
                <div className="w-full max-w-[1024px] mx-auto p-4">
                    <form
                        action={async formData => {
                            "use server";
                            await signIn("credentials", {
                                password: formData.get("password") as string,
                                redirectTo: "/manage",
                            });
                        }}
                        className="w-full flex justify-center items-center gap-4"
                    >
                        <Input type="password" name="password" label="パスワード" required />
                        <Button type="submit">ログイン</Button>
                    </form>
                </div>
            </section>
        </main>
    );
}