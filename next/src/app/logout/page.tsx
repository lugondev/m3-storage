'use client'

import {useEffect} from 'react'
import {useRouter} from 'next/navigation'
import {useAuth} from '@/contexts/AuthContext'

export default function LogoutPage() {
	const {logout, loading: authLoading, user} = useAuth()
	const router = useRouter()

	useEffect(() => {
		const performLogout = async () => {
			if (user) {
				await logout()
			}
			// After logout, the AuthContext's useEffect should redirect to /login
			// If for some reason it doesn't, or if the user is already logged out and lands here,
			// we can add a fallback redirect.
			if (!authLoading && !user) {
				router.push('/login')
			}
		}

		performLogout()
	}, [logout, router, user, authLoading])

	return (
		<div className='flex items-center justify-center min-h-screen'>
			<p>Logging out...</p>
		</div>
	)
}
