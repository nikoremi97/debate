"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useApiKey } from "@/lib/use-api-key"
import { Loader2 } from "lucide-react"

interface ProtectedRouteProps {
    children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
    const { apiKey, isLoading } = useApiKey()
    const router = useRouter()

    useEffect(() => {
        if (!isLoading && !apiKey) {
            router.push("/login")
        }
    }, [apiKey, isLoading, router])

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="flex items-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span>Loading...</span>
                </div>
            </div>
        )
    }

    if (!apiKey) {
        return null // Will redirect to login
    }

    return <>{children}</>
}
