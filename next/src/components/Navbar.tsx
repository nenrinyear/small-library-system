"use client";

import { Button, Link, Navbar, NavbarBrand, NavbarContent, NavbarItem, NavbarMenu, NavbarMenuItem, NavbarMenuToggle } from "@nextui-org/react";
import { useState } from "react";

export default function HeaderNavigation() { 
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const navbarMenu: {
        label: string;
        href: string;
        color?: "primary" | "foreground" | "secondary" | "success" | "warning" | "danger";
    }[] = [
        {
            label: "検索",
            href: "/search",
        },
        {
            label: "タグ一覧",
            href: "/tags",
        },
        {
            label: "管理者ログイン",
            href: "/login",
            color: "primary",
        },
    ];
    return (
        <Navbar onMenuOpenChange={setIsMenuOpen}>
            <NavbarContent>
                <NavbarMenuToggle
                    aria-label={isMenuOpen ? "メニューを閉じる" : "メニューを開く"}
                    className="sm:hidden"
                />
                <NavbarBrand>
                    <Link color="foreground" href="/" className="font-bold">{process.env.NEXT_PUBLIC_SYSTEM_NAME_SHORT ?? "蔵書管理"}</Link>
                </NavbarBrand>
            </NavbarContent>

            <NavbarContent className="hidden sm:flex gap-4" justify="center">
                <NavbarItem>
                    <Link color="foreground" href="/search">検索</Link>
                </NavbarItem>
                <NavbarItem>
                    <Link color="foreground" href="/tags">タグ一覧</Link>
                </NavbarItem>
            </NavbarContent>

            <NavbarContent className="hidden sm:flex gap-4" justify="end">
                <NavbarItem>
                    <Button as={Link} color="primary" href="/login" variant="flat">管理者ログイン</Button>
                </NavbarItem>
            </NavbarContent>

            <NavbarMenu>
                {navbarMenu.map((item, index) => (
                    <NavbarMenuItem key={`${item.label}-${index}`}>
                        <Link
                            color={item.color ?? "foreground"}
                            className="w-full"
                            size="lg"
                            href={item.href}
                        >    
                            {item.label}
                        </Link>
                    </NavbarMenuItem>
                ))}
            </NavbarMenu>
        </Navbar>
    )
}