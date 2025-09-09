"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { MessageSquare, Clock, Hash, Plus, Menu, X } from "lucide-react"
import Link from "next/link"
import { config } from "@/lib/config"
import { useApiKey } from "@/lib/use-api-key"

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

interface ChatSidebarProps {
    currentConversationId?: string | null
    onNewChat: () => void
    refreshTrigger?: number // This will trigger a refresh when changed
    isOpen?: boolean
    onToggle?: () => void
}

export default function ChatSidebar({ currentConversationId, onNewChat, refreshTrigger, isOpen = true, onToggle }: ChatSidebarProps) {
    const [conversations, setConversations] = useState<ConversationSummary[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const { apiKey } = useApiKey()

    const loadConversations = async () => {
        try {
            setLoading(true)
            setError(null)

            const headers: Record<string, string> = {}
            if (apiKey) {
                headers["X-API-Key"] = apiKey
            }

            const response = await fetch(
                `${config.apiUrl}/conversations?limit=50&offset=0`,
                { headers }
            )

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }

            const data: ListConversationsResponse = await response.json()
            setConversations(data.conversations || [])
        } catch (err) {
            console.error("Failed to load conversations:", err)
            setError(err instanceof Error ? err.message : "Failed to load conversations")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        loadConversations()
    }, [refreshTrigger]) // Also refresh when refreshTrigger changes

    const formatDate = (dateString: string) => {
        const date = new Date(dateString)
        const now = new Date()
        const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60)

        if (diffInHours < 24) {
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        } else if (diffInHours < 24 * 7) {
            return date.toLocaleDateString([], { weekday: 'short' })
        } else {
            return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
        }
    }

    return (
        <>
            {/* Mobile Toggle Button */}
            <div className="lg:hidden fixed top-4 left-4 z-50">
                <Button
                    onClick={onToggle}
                    size="sm"
                    variant="outline"
                    className="h-8 w-8 p-0 bg-white shadow-md"
                >
                    {isOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
                </Button>
            </div>

            {/* Sidebar */}
            <div className={`
                bg-slate-50 border-r border-slate-200 flex flex-col h-full transition-all duration-300 min-h-0
                ${isOpen
                    ? 'w-72 sm:w-80 lg:w-80'
                    : 'w-0 lg:w-80 overflow-hidden'
                }
                lg:relative lg:translate-x-0
                ${isOpen
                    ? 'fixed lg:relative inset-y-0 left-0 z-40 translate-x-0'
                    : 'fixed lg:relative inset-y-0 left-0 z-40 -translate-x-full'
                }
            `}>
                {/* Header */}
                <div className="p-4 border-b border-slate-200 flex-shrink-0">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-lg font-semibold text-slate-900">Chat History</h2>
                        <div className="flex items-center gap-2">
                            <Button
                                onClick={onNewChat}
                                size="sm"
                                className="h-8 px-3"
                            >
                                <Plus className="h-4 w-4 mr-1" />
                                <span className="hidden sm:inline">New Chat</span>
                            </Button>
                            {onToggle && (
                                <Button
                                    onClick={onToggle}
                                    size="sm"
                                    variant="ghost"
                                    className="h-8 w-8 p-0 lg:hidden"
                                >
                                    <X className="h-4 w-4" />
                                </Button>
                            )}
                        </div>
                    </div>
                </div>

                {/* Conversations List */}
                <ScrollArea className="flex-1 p-4 min-h-0">
                    {loading ? (
                        <div className="flex items-center justify-center py-8">
                            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-slate-600"></div>
                        </div>
                    ) : error ? (
                        <div className="text-center py-8">
                            <p className="text-sm text-red-600 mb-2">{error}</p>
                            <Button
                                onClick={loadConversations}
                                variant="outline"
                                size="sm"
                            >
                                Retry
                            </Button>
                        </div>
                    ) : conversations.length === 0 ? (
                        <div className="text-center py-8">
                            <MessageSquare className="h-8 w-8 text-slate-400 mx-auto mb-2" />
                            <p className="text-sm text-slate-500">No conversations yet</p>
                            <p className="text-xs text-slate-400 mt-1">Start a new chat to begin</p>
                        </div>
                    ) : (
                        <div className="space-y-2 pr-2">
                            {conversations.map((conversation) => (
                                <Link
                                    key={conversation.id}
                                    href={`/chat?conversation_id=${conversation.id}`}
                                    className="block"
                                >
                                    <Card
                                        className={`cursor-pointer transition-all duration-200 hover:shadow-md overflow-hidden ${currentConversationId === conversation.id
                                            ? "ring-2 ring-blue-500 bg-blue-50"
                                            : "hover:bg-white"
                                            }`}
                                    >
                                        <CardContent className="p-3 w-full min-w-0">
                                            <div className="space-y-2 w-full min-w-0">
                                                {/* Topic */}
                                                <div className="flex items-center gap-2 min-w-0">
                                                    <Badge variant="secondary" className="text-xs flex-shrink-0">
                                                        {conversation.topic_name}
                                                    </Badge>
                                                </div>

                                                {/* Title */}
                                                <h3 className="font-medium text-sm text-slate-900 line-clamp-2 break-words min-w-0">
                                                    {conversation.title}
                                                </h3>

                                                {/* Metadata */}
                                                <div className="flex items-center justify-between text-xs text-slate-500 min-w-0">
                                                    <div className="flex items-center gap-1 flex-shrink-0">
                                                        <MessageSquare className="h-3 w-3" />
                                                        <span className="whitespace-nowrap">{conversation.message_count} msgs</span>
                                                    </div>
                                                    <div className="flex items-center gap-1 flex-shrink-0">
                                                        <Clock className="h-3 w-3" />
                                                        <span className="whitespace-nowrap">{formatDate(conversation.updated_at)}</span>
                                                    </div>
                                                </div>
                                            </div>
                                        </CardContent>
                                    </Card>
                                </Link>
                            ))}
                        </div>
                    )}
                </ScrollArea>
            </div>

            {/* Mobile Overlay */}
            {isOpen && onToggle && (
                <div
                    className="fixed inset-0 bg-black bg-opacity-50 z-30 lg:hidden"
                    onClick={onToggle}
                />
            )}
        </>
    )
}
