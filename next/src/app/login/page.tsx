'use client'
import {LoginContent} from '@/components/auth/LoginContent'
import {Suspense} from 'react'

export default function LoginPage() {
	return (
		<Suspense>
			<LoginContent />
		</Suspense>
	)
}
