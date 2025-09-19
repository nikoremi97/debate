"use client"

import { useState, useEffect, useRef, Suspense } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Loader2, Send, MessageSquare, LogOut } from "lucide-react"
import { config } from "@/lib/config"
import ChatSidebar from "@/components/chat-sidebar"
import { useApiKey } from "@/lib/use-api-key"
import { ProtectedRoute } from "@/components/protected-route"

// Types matching your backend API
type Message = {
    role: "user" | "bot"
    message: string
    ts: number
}

type ChatResponse = {
    conversation_id: string
    message: Message[]
    topic?: string
    stance?: string
}

type ChatRequest = {
    conversation_id?: string
    message: string
    topic?: string
}

function ChatContent() {
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
    const { apiKey, clearApiKey } = useApiKey()

    // Check for conversation_id in URL params
    useEffect(() => {
        const urlConversationId = searchParams.get('conversation_id')
        if (urlConversationId) {
            setConversationId(urlConversationId)
            // Load existing conversation
            loadConversation(urlConversationId)
            // Close sidebar on mobile when conversation is loaded
            setSidebarOpen(false)
        } else {
            // Clear conversation state when no conversation_id in URL
            setConversationId(null)
            setMessages([])
            setTopic("")
            setStance("")
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

            const headers: Record<string, string> = {}
            // Only send API key if we have one and we're not running locally
            if (apiKey && !config.apiUrl.includes('localhost')) {
                headers["X-API-Key"] = apiKey
            }

            const response = await fetch(`${config.apiUrl}/conversations/${id}`, {
                headers
            })

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

        // Clear topic and stance for new conversations
        if (!conversationId) {
            setTopic("")
            setStance("")
        }

        try {
            const requestBody: ChatRequest = {
                message: userMessage.message
            }

            // Include conversation_id if we have one
            if (conversationId) {
                requestBody.conversation_id = conversationId
            }

            // For new conversations, use the input message as the topic
            if (!conversationId || conversationId === "") {
                requestBody.topic = userMessage.message
                console.log("Debug - Using input message as topic for new conversation:", userMessage.message)
            } else {
                // For existing conversations, don't send topic - let backend use existing topic
                console.log("Debug - Existing conversation, not sending topic. Backend should use existing topic.")
            }

            console.log("Debug - Final request body:", requestBody)

            const headers: Record<string, string> = { "Content-Type": "application/json" }
            // Only send API key if we have one and we're not running locally
            if (apiKey && !config.apiUrl.includes('localhost')) {
                headers["X-API-Key"] = apiKey
            }

            const res = await fetch(`${config.apiUrl}${config.endpoints.chat}`, {
                method: "POST",
                headers,
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

            // Set topic and stance from the response
            if (data.topic) {
                setTopic(data.topic)
            }
            if (data.stance) {
                setStance(data.stance)
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
        <ProtectedRoute>
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
                            <CardTitle className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <MessageSquare className="h-5 w-5" />
                                    Debate Chatbot
                                </div>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => {
                                        clearApiKey()
                                        router.push('/login')
                                    }}
                                    className="flex items-center gap-2"
                                >
                                    <LogOut className="h-4 w-4" />
                                    Logout
                                </Button>
                            </CardTitle>
                            {topic && (
                                <div className="flex gap-2 items-center">
                                    <Badge variant="outline">Topic: {topic}</Badge>
                                    {stance && <Badge variant="secondary">Bot Stance: {stance}</Badge>}
                                </div>
                            )}
                        </CardHeader>
                        <CardContent className="flex-1 flex flex-col min-h-0 p-4">
                            <ScrollArea ref={scrollAreaRef} className="flex-1 pr-4 min-h-0">
                                <div className="space-y-4">
                                    {messages.length === 0 && (
                                        <div className="text-center text-muted-foreground py-8">
                                            <MessageSquare className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                            <p>Start a debate! Enter a topic above (optional) and send your first message to begin.</p>
                                            <p className="text-sm mt-2">The bot will take the opposite stance to challenge your perspective.</p>
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


                            {/* Topic input - only show for new conversations */}
                            {!conversationId && (
                                <div className="mb-4">
                                    <p className="text-sm text-muted-foreground mb-2">
                                        💡 <strong>Tip:</strong> Type your debate topic in the chat below (e.g., "Electric cars vs gas cars, I prefer electric cars")
                                    </p>
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
        </ProtectedRoute>
    )
}

export default function ChatPage() {
    return (
        <Suspense fallback={
            <div className="flex h-[calc(100vh-73px)] items-center justify-center">
                <div className="flex items-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span>Loading chat...</span>
                </div>
            </div>
        }>
            <ChatContent />
        </Suspense>
    )
}
