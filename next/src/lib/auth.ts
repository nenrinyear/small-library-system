import NextAuth from "next-auth"
import credentials from "next-auth/providers/credentials"
import { ADMIN_PASSWORD } from "./environment-variables"

export const { handlers, signIn, signOut, auth } = NextAuth({
    providers: [
        credentials({
            credentials: {
                password: { label: "パスワード", type: "password" },
            },
            authorize: async (credentials) => {
                try {
                    if (credentials.password === ADMIN_PASSWORD) {
                        return Promise.resolve({ id: "1", name: "Admin" })
                    }
                    return Promise.reject(new Error("認証に失敗しました。"))
                } catch (error) {
                    return Promise.reject(new Error("認証に失敗しました。"))
                }
            }
        })
    ],
    secret: process.env.SECRET,
    session: {
        strategy: "jwt",
    },
    pages: {
        signIn: "/login",
        signOut: "/logout",
        error: "/login",
    },
    logger: {
        error: console.error,
        warn: console.warn,
        info: console.info,
        debug: console.debug,
    },
    callbacks: {
        jwt({ token, user }) {
            if (user) {
                token.id = user.id
            }
            return token
        },
        session({ session, token }) {
            session.user.id = token.id as any
            return session
        }
    }
})