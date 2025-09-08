import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { MessageSquare, Brain, Zap, Users } from "lucide-react";

export default function Home() {
  return (
    <div className="min-h-[calc(100vh-73px)] bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <div className="flex justify-center mb-6">
            <div className="p-4 bg-primary rounded-full">
              <MessageSquare className="h-12 w-12 text-primary-foreground" />
            </div>
          </div>
          <h1 className="text-5xl font-bold text-gray-900 dark:text-white mb-4">
            Debate Chatbot
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto mb-8">
            Engage in intelligent debates with an AI-powered chatbot that takes a stance and argues persuasively.
            Challenge your thinking and sharpen your argumentation skills.
          </p>
          <Link href="/chat">
            <Button size="lg" className="text-lg px-8 py-6">
              Start Debating
              <MessageSquare className="ml-2 h-5 w-5" />
            </Button>
          </Link>
        </div>

        <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <Card className="text-center">
            <CardHeader>
              <Brain className="h-12 w-12 mx-auto text-blue-600 mb-4" />
              <CardTitle>AI-Powered</CardTitle>
              <CardDescription>
                Powered by advanced language models for intelligent and contextual responses
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="text-center">
            <CardHeader>
              <Zap className="h-12 w-12 mx-auto text-yellow-600 mb-4" />
              <CardTitle>Dynamic Topics</CardTitle>
              <CardDescription>
                The bot randomly selects debate topics and takes a stance to challenge your perspective
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="text-center">
            <CardHeader>
              <Users className="h-12 w-12 mx-auto text-green-600 mb-4" />
              <CardTitle>Persuasive Arguments</CardTitle>
              <CardDescription>
                Experience compelling arguments that will make you think critically about various topics
              </CardDescription>
            </CardHeader>
          </Card>
        </div>

        <div className="mt-16 text-center">
          <Card className="max-w-4xl mx-auto">
            <CardHeader>
              <CardTitle>How It Works</CardTitle>
              <CardDescription>
                Get started with your debate in just a few simple steps
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-3 gap-6 text-left">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">1</div>
                    <h3 className="font-semibold">Start a Conversation</h3>
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Send your first message to begin a new debate session
                  </p>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">2</div>
                    <h3 className="font-semibold">Bot Takes a Stance</h3>
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    The AI selects a topic and takes either a PRO or CON position
                  </p>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-bold">3</div>
                    <h3 className="font-semibold">Debate & Learn</h3>
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Engage in back-and-forth arguments to explore different perspectives
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
