'use client'
import {useAuth} from '@/contexts/AuthContext'
import {useRouter, usePathname} from 'next/navigation'
import {useEffect, useState} from 'react'

export default function ProtectedWrapper({children}: {children: React.ReactNode}) {
	const {user, loading, isAuthenticated} = useAuth()
	const router = useRouter()
	const pathname = usePathname()
	const [isCheckingAuth, setIsCheckingAuth] = useState(true)

	useEffect(() => {
		if (loading) {
			setIsCheckingAuth(true)
			return
		}

		if (!isAuthenticated || !user) {
			setIsCheckingAuth(false)
			if (pathname !== '/login' && pathname !== '/') {
				router.push('/login?redirect=' + encodeURIComponent(pathname))
			}
			return
		}

		setIsCheckingAuth(false)
	}, [user, loading, isAuthenticated, router, pathname])

	// During initial auth check, render children to maintain layout
	// The layout components will handle their own loading states
	if (isCheckingAuth) {
		return <>{children}</>
	}

	// Not authenticated and not loading anymore - layout components will handle redirection
	if (!isAuthenticated && !loading) {
		return null
	}

	return <>{children}</>
}
