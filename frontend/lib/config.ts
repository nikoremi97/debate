// Configuration for the frontend application
export const config = {
    // Backend API URL - can be overridden by environment variable
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'https://d13sbjy1c5yh6c.cloudfront.net',

    // API endpoints
    endpoints: {
        chat: '/chat',
        health: '/health',
        ready: '/ready',
        conversations: '/conversations'
    }
} as const
