"use client"

import { useState, useEffect } from "react"

export function useApiKey() {
    const [apiKey, setApiKey] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        // Get API key from localStorage on mount
        const storedApiKey = localStorage.getItem("debate_api_key")
        setApiKey(storedApiKey)
        setIsLoading(false)
    }, [])

    const setApiKeyAndStore = (key: string) => {
        localStorage.setItem("debate_api_key", key)
        setApiKey(key)
    }

    const clearApiKey = () => {
        localStorage.removeItem("debate_api_key")
        setApiKey(null)
    }

    return {
        apiKey,
        isLoading,
        setApiKey: setApiKeyAndStore,
        clearApiKey
    }
}
