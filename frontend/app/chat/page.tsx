"use client"

import { useState, useEffect, useRef } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Loader2, Send, MessageSquare } from "lucide-react"
import { config } from "@/lib/config"
import ChatSidebar from "@/components/chat-sidebar"

// Types matching your backend API
type Message = {
    role: "user" | "bot"
    message: string
    ts: number
}

type ChatResponse = {
    conversation_id: string
    message: Message[]
}

type ChatRequest = {
    conversation_id?: string
    message: string
}

export default function ChatPage() {
    const searchParams = useSearchParams()
    const router = useRouter()
    const [messages, setMessages] = useState<Message[]>([])
    const [input, setInput] = useState("")
    const [conversationId, setConversationId] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [topic, setTopic] = useState<string>("")
    const [stance, setStance] = useState<string>("")
    const [sidebarRefreshTrigger, setSidebarRefreshTrigger] = useState(0)
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const scrollAreaRef = useRef<HTMLDivElement>(null)

    // Check for conversation_id in URL params
    useEffect(() => {
        const urlConversationId = searchParams.get('conversation_id')
        if (urlConversationId) {
            setConversationId(urlConversationId)
            // Load existing conversation
            loadConversation(urlConversationId)
            // Close sidebar on mobile when conversation is loaded
            setSidebarOpen(false)
        }
    }, [searchParams])

    // Auto-scroll to bottom when messages change
    useEffect(() => {
        if (scrollAreaRef.current) {
            const scrollContainer = scrollAreaRef.current.querySelector('[data-radix-scroll-area-viewport]')
            if (scrollContainer) {
                scrollContainer.scrollTop = scrollContainer.scrollHeight
            }
        }
    }, [messages])

    const loadConversation = async (id: string) => {
        try {
            setIsLoading(true)
            setError(null)

            const response = await fetch(`${config.apiUrl}/conversations/${id}`)

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }

            const conversation = await response.json()

            // Update messages with the conversation history
            setMessages(conversation.messages || [])

            // Set the conversation ID
            setConversationId(id)

            // Set topic and stance from the conversation
            if (conversation.topic) {
                setTopic(conversation.topic)
            }
            if (conversation.stance) {
                setStance(conversation.stance)
            }
        } catch (err) {
            console.error("Error loading conversation:", err)
            setError(err instanceof Error ? err.message : "Failed to load conversation")
        } finally {
            setIsLoading(false)
        }
    }

    async function sendMessage() {
        if (!input.trim() || isLoading) return

        const userMessage: Message = {
            role: "user",
            message: input,
            ts: Date.now()
        }

        setMessages(prev => [...prev, userMessage])
        setInput("")
        setIsLoading(true)
        setError(null)

        try {
            const requestBody: ChatRequest = {
                message: userMessage.message
            }

            // Include conversation_id if we have one
            if (conversationId) {
                requestBody.conversation_id = conversationId
            }

            const res = await fetch(`${config.apiUrl}${config.endpoints.chat}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(requestBody),
            })

            if (!res.ok) {
                const errorText = await res.text()
                throw new Error(`HTTP error! status: ${res.status}, response: ${errorText}`)
            }

            const data: ChatResponse = await res.json()

            // Update conversation ID if this is a new conversation
            if (!conversationId) {
                setConversationId(data.conversation_id)
                // Trigger sidebar refresh to show the new conversation
                setSidebarRefreshTrigger(prev => prev + 1)
            }

            // Update messages with the full conversation history
            setMessages(data.message)

            // Extract topic and stance from the first bot message if available
            if (data.message.length > 0) {
                const firstBotMessage = data.message.find(m => m.role === "bot")
                if (firstBotMessage && firstBotMessage.message.includes("Topic:")) {
                    const topicMatch = firstBotMessage.message.match(/Topic:\s*([^\n]+)/)
                    const stanceMatch = firstBotMessage.message.match(/Stance:\s*([^\n]+)/)
                    if (topicMatch) setTopic(topicMatch[1].trim())
                    if (stanceMatch) setStance(stanceMatch[1].trim())
                }
            }

        } catch (err) {
            console.error("Full error details:", err)
            const errorMessage = err instanceof Error ? err.message : "Failed to send message"
            setError(`Error: ${errorMessage}`)
            // Remove the user message if the request failed
            setMessages(prev => prev.slice(0, -1))
        } finally {
            setIsLoading(false)
        }
    }

    const handleNewChat = () => {
        // Reset all state for a new chat
        setMessages([])
        setInput("")
        setConversationId(null)
        setError(null)
        setTopic("")
        setStance("")

        // Close sidebar on mobile after starting new chat
        setSidebarOpen(false)

        // Navigate to chat without conversation_id
        router.push('/chat')
    }

    const toggleSidebar = () => {
        setSidebarOpen(!sidebarOpen)
    }


    return (
        <div className="flex h-[calc(100vh-73px)]">
            {/* Sidebar */}
            <ChatSidebar
                currentConversationId={conversationId}
                onNewChat={handleNewChat}
                refreshTrigger={sidebarRefreshTrigger}
                isOpen={sidebarOpen}
                onToggle={toggleSidebar}
            />

            {/* Main Chat Area */}
            <div className="flex-1 flex flex-col m-4 min-h-0 lg:ml-4">
                <Card className="flex-1 flex flex-col min-h-0">
                    <CardHeader className="pb-3 flex-shrink-0">
                        <CardTitle className="flex items-center gap-2">
                            <MessageSquare className="h-5 w-5" />
                            Debate Chatbot
                        </CardTitle>
                        {topic && (
                            <div className="flex gap-2 items-center">
                                <Badge variant="outline">Topic: {topic}</Badge>
                            </div>
                        )}
                    </CardHeader>
                    <CardContent className="flex-1 flex flex-col min-h-0 p-4">
                        <ScrollArea ref={scrollAreaRef} className="flex-1 pr-4 min-h-0">
                            <div className="space-y-4">
                                {messages.length === 0 && (
                                    <div className="text-center text-muted-foreground py-8">
                                        <MessageSquare className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                        <p>Start a debate! Send your first message to begin.</p>
                                    </div>
                                )}
                                {messages.map((message, i) => (
                                    <div
                                        key={`${message.ts}-${i}`}
                                        className={`flex ${message.role === "user" ? "justify-end" : "justify-start"}`}
                                    >
                                        <div
                                            className={`max-w-[80%] rounded-lg p-3 ${message.role === "user"
                                                ? "bg-primary text-primary-foreground"
                                                : "bg-muted"
                                                }`}
                                        >
                                            <div className="text-sm font-medium mb-1">
                                                {message.role === "user" ? "You" : "Debate Bot"}
                                            </div>
                                            <div className="whitespace-pre-wrap">{message.message}</div>
                                        </div>
                                    </div>
                                ))}
                                {isLoading && (
                                    <div className="flex justify-start">
                                        <div className="bg-muted rounded-lg p-3 flex items-center gap-2">
                                            <Loader2 className="h-4 w-4 animate-spin" />
                                            <span className="text-sm text-muted-foreground">Bot is thinking...</span>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </ScrollArea>

                        {error && (
                            <div className="mt-4 p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
                                <p className="text-sm text-destructive">{error}</p>
                            </div>
                        )}


                        <form
                            onSubmit={e => {
                                e.preventDefault()
                                sendMessage()
                            }}
                            className="flex gap-2 mt-4 flex-shrink-0"
                        >
                            <Input
                                value={input}
                                onChange={e => setInput(e.target.value)}
                                placeholder="Type your debate argument..."
                                disabled={isLoading}
                                className="flex-1"
                            />
                            <Button
                                type="submit"
                                disabled={isLoading || !input.trim()}
                            >
                                {isLoading ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                ) : (
                                    <Send className="h-4 w-4" />
                                )}
                            </Button>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
