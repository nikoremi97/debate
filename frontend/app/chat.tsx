"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Loader2, Send, MessageSquare } from "lucide-react"
import { config } from "@/lib/config"

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
    const [messages, setMessages] = useState<Message[]>([])
    const [input, setInput] = useState("")
    const [conversationId, setConversationId] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [topic, setTopic] = useState<string>("")
    const [stance, setStance] = useState<string>("")

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
                throw new Error(`HTTP error! status: ${res.status}`)
            }

            const data: ChatResponse = await res.json()

            // Update conversation ID if this is a new conversation
            if (!conversationId) {
                setConversationId(data.conversation_id)
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
            setError(err instanceof Error ? err.message : "Failed to send message")
            // Remove the user message if the request failed
            setMessages(prev => prev.slice(0, -1))
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="flex flex-col h-[calc(100vh-73px)] max-w-4xl mx-auto p-4">
            <Card className="flex-1 overflow-hidden">
                <CardHeader className="pb-3">
                    <CardTitle className="flex items-center gap-2">
                        <MessageSquare className="h-5 w-5" />
                        Debate Chatbot
                    </CardTitle>
                    {topic && (
                        <div className="flex gap-2 items-center">
                            <Badge variant="outline">Topic: {topic}</Badge>
                            {stance && (
                                <Badge variant={stance === "PRO" ? "default" : "secondary"}>
                                    {stance}
                                </Badge>
                            )}
                        </div>
                    )}
                </CardHeader>
                <CardContent className="h-full overflow-hidden">
                    <ScrollArea className="h-full pr-4">
                        <div className="space-y-4">
                            {messages.length === 0 && (
                                <div className="text-center text-muted-foreground py-8">
                                    <MessageSquare className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                    <p>Start a debate! Send your first message to begin.</p>
                                </div>
                            )}
                            {messages.map((message, i) => (
                                <div
                                    key={i}
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
                </CardContent>
            </Card>

            {error && (
                <div className="mt-2 p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
                    <p className="text-sm text-destructive">{error}</p>
                </div>
            )}

            <form
                onSubmit={e => {
                    e.preventDefault()
                    sendMessage()
                }}
                className="flex gap-2 mt-4"
            >
                <Input
                    value={input}
                    onChange={e => setInput(e.target.value)}
                    placeholder="Type your debate argument..."
                    disabled={isLoading}
                    className="flex-1"
                />
                <Button type="submit" disabled={isLoading || !input.trim()}>
                    {isLoading ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                        <Send className="h-4 w-4" />
                    )}
                </Button>
            </form>
        </div>
    )
}
