"use client"

import Link from "next/link"
import { Button } from "@/components/ui/button"
import { MessageSquare, Home, LogIn } from "lucide-react"
import { useApiKey } from "@/lib/use-api-key"

export function Navigation() {
    const { apiKey } = useApiKey()

    return (
        <nav className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="container mx-auto px-4 py-3">
                <div className="flex items-center justify-between">
                    <Link href="/" className="flex items-center gap-2">
                        <MessageSquare className="h-6 w-6" />
                        <span className="font-bold text-lg">Debate Chatbot</span>
                    </Link>

                    <div className="flex items-center gap-2">
                        <Link href="/">
                            <Button variant="ghost" size="sm">
                                <Home className="h-4 w-4 mr-2" />
                                Home
                            </Button>
                        </Link>
                        {apiKey ? (
                            <Link href="/chat">
                                <Button size="sm">
                                    <MessageSquare className="h-4 w-4 mr-2" />
                                    Chat
                                </Button>
                            </Link>
                        ) : (
                            <Link href="/login">
                                <Button size="sm">
                                    <LogIn className="h-4 w-4 mr-2" />
                                    Login
                                </Button>
                            </Link>
                        )}
                    </div>
                </div>
            </div>
        </nav>
    )
}
