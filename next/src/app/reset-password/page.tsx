'use client'

import {ResetPasswordForm} from '@/components/auth/ResetPasswordForm'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import {useSearchParams} from 'next/navigation'
import {Suspense} from 'react' // Needed for useSearchParams
import Link from 'next/link'

function ResetPasswordContent() {
	const searchParams = useSearchParams()
	const token = searchParams.get('token')

	return (
		<div className='flex min-h-screen items-center justify-center bg-gray-100 dark:bg-gray-900'>
			<Card className='w-full max-w-md mx-4'>
				<CardHeader className='text-center'>
					<CardTitle>Reset Your Password</CardTitle>
					<CardDescription>Enter and confirm your new password below.</CardDescription>
				</CardHeader>
				<CardContent>
					<ResetPasswordForm token={token} />
					<div className='mt-4 text-center text-sm'>
						Return to{' '}
						<Link href='/login' className='underline'>
							Sign in
						</Link>
					</div>
				</CardContent>
			</Card>
		</div>
	)
}

// Wrap the component using Suspense for useSearchParams
export default function ResetPasswordPage() {
	return (
		<Suspense fallback={<div>Loading...</div>}>
			<ResetPasswordContent />
		</Suspense>
	)
}
