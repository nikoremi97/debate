// Configuration for the frontend application
export const config = {
    // Backend API URL - can be overridden by environment variable
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',

    // API endpoints
    endpoints: {
        chat: '/chat',
        health: '/health',
        ready: '/ready',
        conversations: '/conversations'
    }
} as const
