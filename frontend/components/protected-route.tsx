"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useApiKey } from "@/lib/use-api-key"
import { Loader2 } from "lucide-react"
import { config } from "@/lib/config"

interface ProtectedRouteProps {
    children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
    const { apiKey, isLoading } = useApiKey()
    const router = useRouter()

    // Allow access when running locally (no API key required)
    const isLocalDevelopment = config.apiUrl.includes('localhost')

    useEffect(() => {
        if (!isLoading && !apiKey && !isLocalDevelopment) {
            router.push("/login")
        }
    }, [apiKey, isLoading, router, isLocalDevelopment])

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

    if (!apiKey && !isLocalDevelopment) {
        return null // Will redirect to login
    }

    return <>{children}</>
}
