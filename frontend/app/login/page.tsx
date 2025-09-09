"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Loader2, Key, Shield } from "lucide-react"
import { config } from "@/lib/config"

export default function LoginPage() {
    const [apiKey, setApiKey] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const router = useRouter()

    async function handleLogin() {
        if (!apiKey.trim()) {
            setError("Please enter an API key")
            return
        }

        setIsLoading(true)
        setError(null)

        try {
            // Test the API key by making a request to the health endpoint
            const response = await fetch(`${config.apiUrl}/health`, {
                method: "GET",
                headers: {
                    "X-API-Key": apiKey.trim()
                }
            })

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }

            // Store the API key in localStorage
            localStorage.setItem("debate_api_key", apiKey.trim())

            // Redirect to chat page
            router.push("/chat")
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to authenticate")
        } finally {
            setIsLoading(false)
        }
    }

    function handleKeyPress(e: React.KeyboardEvent) {
        if (e.key === "Enter") {
            handleLogin()
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="flex justify-center mb-4">
                        <div className="p-3 bg-primary/10 rounded-full">
                            <Shield className="h-8 w-8 text-primary" />
                        </div>
                    </div>
                    <CardTitle className="text-2xl font-bold">API Key Authentication</CardTitle>
                    <p className="text-muted-foreground mt-2">
                        Enter your API key to access the Debate Chatbot
                    </p>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="api-key">API Key</Label>
                        <div className="relative">
                            <Key className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                            <Input
                                id="api-key"
                                type="password"
                                placeholder="Enter your API key"
                                value={apiKey}
                                onChange={(e) => setApiKey(e.target.value)}
                                onKeyPress={handleKeyPress}
                                disabled={isLoading}
                                className="pl-10"
                            />
                        </div>
                    </div>

                    {error && (
                        <Alert variant="destructive">
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}

                    <Button
                        onClick={handleLogin}
                        disabled={isLoading || !apiKey.trim()}
                        className="w-full"
                    >
                        {isLoading ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                Authenticating...
                            </>
                        ) : (
                            "Login"
                        )}
                    </Button>

                    <div className="text-center text-sm text-muted-foreground">
                        <p>Don't have an API key?</p>
                        <p>Contact your administrator to get access.</p>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
