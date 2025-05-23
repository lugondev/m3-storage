'use client'

import React, {useState} from 'react'
import {QueryClient, QueryClientProvider} from '@tanstack/react-query'
// Optionally, import ReactQueryDevtools for development
// import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

interface QueryProviderProps {
	children: React.ReactNode
}

export function QueryProvider({children}: QueryProviderProps) {
	// Use useState to ensure the QueryClient is only created once per component instance
	const [queryClient] = useState(
		() =>
			new QueryClient({
				defaultOptions: {
					queries: {
						// Global defaults for queries if needed
						// Example: staleTime: 5 * 60 * 1000, // 5 minutes
						// Example: refetchOnWindowFocus: false,
					},
					mutations: {
						// Global defaults for mutations if needed
					},
				},
			}),
	)

	return (
		<QueryClientProvider client={queryClient}>
			{children}
			{/* Optionally include DevTools for debugging in development */}
			{/* <ReactQueryDevtools initialIsOpen={false} /> */}
		</QueryClientProvider>
	)
}
