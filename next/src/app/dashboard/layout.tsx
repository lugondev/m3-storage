'use client'

import {useAuth} from '@/contexts/AuthContext'
import {useRouter} from 'next/navigation'
import {useEffect} from 'react'
import AppShell from '@/components/layout/AppShell' // Import the new AppShell

export default function ProtectedLayout({children}: {children: React.ReactNode}) {
	const {user, loading} = useAuth()
	const router = useRouter()

	useEffect(() => {
		if (!loading && !user) {
			router.push('/')
		}
	}, [user, loading, router])

	if (loading) {
		return (
			<div className='flex h-screen items-center justify-center'>
				<div className='animate-pulse space-y-4'>
					<div className='h-4 bg-muted rounded w-[200px]'></div>
					<div className='h-4 bg-muted rounded w-[160px]'></div>
				</div>
			</div>
		)
	}

	if (!user) return null

	// Use the AppShell to wrap the children
	return <AppShell>{children}</AppShell>
}
