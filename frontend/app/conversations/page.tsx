"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { MessageSquare, Clock, Hash } from "lucide-react"
import Link from "next/link"
import { config } from "@/lib/config"

type ConversationSummary = {
    id: string
    topic_name: string
    bot_stance: string
    title: string
    message_count: number
    created_at: string
    updated_at: string
}

type ListConversationsResponse = {
    conversations: ConversationSummary[]
    total: number
    page: number
    limit: number
}

export default function ConversationsPage() {
    const [conversations, setConversations] = useState<ConversationSummary[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [page, setPage] = useState(1)
    const [hasMore, setHasMore] = useState(true)

    const loadConversations = async (pageNum: number = 1) => {
        try {
            setLoading(true)
            setError(null)

            const response = await fetch(
                `${config.apiUrl}/conversations?limit=20&offset=${(pageNum - 1) * 20}`
            )

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }

            const data: ListConversationsResponse = await response.json()

            if (pageNum === 1) {
                setConversations(data.conversations)
            } else {
                setConversations(prev => [...prev, ...data.conversations])
            }

            setHasMore(data.conversations.length === 20)
            setPage(pageNum)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load conversations")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        loadConversations(1)
    }, [])

    const loadMore = () => {
        if (!loading && hasMore) {
            loadConversations(page + 1)
        }
    }

    const formatDate = (dateString: string) => {
        const date = new Date(dateString)
        return date.toLocaleDateString() + " " + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }

    return (
        <div className="flex flex-col h-[calc(100vh-73px)] max-w-4xl mx-auto p-4">
            <div className="mb-6">
                <h1 className="text-3xl font-bold mb-2">Chat History</h1>
                <p className="text-muted-foreground">
                    Browse your past debate conversations and continue where you left off.
                </p>
            </div>

            {error && (
                <Card className="mb-4 border-destructive">
                    <CardContent className="pt-6">
                        <p className="text-destructive">{error}</p>
                        <Button
                            variant="outline"
                            onClick={() => loadConversations(1)}
                            className="mt-2"
                        >
                            Try Again
                        </Button>
                    </CardContent>
                </Card>
            )}

            <ScrollArea className="flex-1">
                <div className="space-y-4">
                    {conversations.map((conversation) => (
                        <Card key={conversation.id} className="hover:shadow-md transition-shadow">
                            <CardHeader>
                                <div className="flex items-start justify-between">
                                    <div className="flex-1">
                                        <CardTitle className="text-lg mb-2">
                                            {conversation.title}
                                        </CardTitle>
                                        <CardDescription className="flex items-center gap-4 text-sm">
                                            <span className="flex items-center gap-1">
                                                <Hash className="h-4 w-4" />
                                                {conversation.topic_name}
                                            </span>
                                            <span className="flex items-center gap-1">
                                                <MessageSquare className="h-4 w-4" />
                                                {conversation.message_count} messages
                                            </span>
                                            <span className="flex items-center gap-1">
                                                <Clock className="h-4 w-4" />
                                                {formatDate(conversation.updated_at)}
                                            </span>
                                        </CardDescription>
                                    </div>
                                    <div className="flex flex-col gap-2">
                                        <Badge
                                            variant={conversation.bot_stance === "PRO" ? "default" : "secondary"}
                                        >
                                            Bot: {conversation.bot_stance}
                                        </Badge>
                                        <Link href={`/chat?conversation_id=${conversation.id}`}>
                                            <Button size="sm" variant="outline">
                                                Continue Chat
                                            </Button>
                                        </Link>
                                    </div>
                                </div>
                            </CardHeader>
                        </Card>
                    ))}

                    {loading && (
                        <Card>
                            <CardContent className="pt-6">
                                <div className="flex items-center justify-center">
                                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                                    <span className="ml-2">Loading conversations...</span>
                                </div>
                            </CardContent>
                        </Card>
                    )}

                    {!loading && conversations.length === 0 && !error && (
                        <Card>
                            <CardContent className="pt-6 text-center">
                                <MessageSquare className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                <h3 className="text-lg font-semibold mb-2">No conversations yet</h3>
                                <p className="text-muted-foreground mb-4">
                                    Start your first debate to see it appear here.
                                </p>
                                <Link href="/chat">
                                    <Button>Start a Debate</Button>
                                </Link>
                            </CardContent>
                        </Card>
                    )}

                    {hasMore && !loading && conversations.length > 0 && (
                        <div className="text-center">
                            <Button variant="outline" onClick={loadMore}>
                                Load More
                            </Button>
                        </div>
                    )}
                </div>
            </ScrollArea>
        </div>
    )
}
