// Configuration for the frontend application
export const config = {
    // Backend API URL - can be overridden by environment variable
    // This will be replaced with the actual CloudFront URL after deployment
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'https://d1234567890.cloudfront.net',

    // API endpoints
    endpoints: {
        chat: '/chat',
        health: '/health',
        ready: '/ready',
        conversations: '/conversations'
    }
} as const
