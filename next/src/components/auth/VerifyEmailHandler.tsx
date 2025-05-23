'use client'

import {verifyEmail} from '@/services/authService' // Import directly from service
import {useState, useEffect} from 'react'
import {toast} from 'sonner'
import Link from 'next/link'
import {Button} from '@/components/ui/button'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import {Skeleton} from '@/components/ui/skeleton' // For loading state

interface VerifyEmailHandlerProps {
	token: string | null // Token received from URL param
}

export function VerifyEmailHandler({token}: VerifyEmailHandlerProps) {
	const [loading, setLoading] = useState(true)
	const [message, setMessage] = useState<string | null>(null)
	const [isError, setIsError] = useState(false)

	useEffect(() => {
		const handleVerification = async () => {
			if (!token) {
				setMessage('Invalid or missing verification token.')
				setIsError(true)
				setLoading(false)
				return
			}

			setLoading(true)
			setIsError(false)
			try {
				const response = await verifyEmail(token)
				setMessage(response.message || 'Email verified successfully!') // Use message from backend if provided
				toast.success(response.message || 'Email verified successfully!')
			} catch (err) {
				console.error('Email verification error:', err)
				// Use const as errorMessage is not reassigned here
				const errorMessage = 'Failed to verify email. The link may be invalid or expired.'
				// You could potentially parse err.response.data for more specific backend errors
				if (err instanceof Error) {
					// Avoid showing generic backend errors
				}
				setMessage(errorMessage)
				setIsError(true)
				toast.error(errorMessage)
			} finally {
				setLoading(false)
			}
		}

		handleVerification()
	}, [token]) // Re-run if token changes (though it shouldn't on a static page load)

	return (
		<Card className='w-full max-w-md'>
			<CardHeader className='text-center'>
				<CardTitle>Email Verification</CardTitle>
				<CardDescription>{loading ? 'Verifying your email address...' : message || 'Processing verification...'}</CardDescription>
			</CardHeader>
			<CardContent className='text-center space-y-4'>
				{loading ? (
					<div className='space-y-2'>
						<Skeleton className='h-4 w-3/4 mx-auto' />
						<Skeleton className='h-4 w-1/2 mx-auto' />
					</div>
				) : (
					<>
						{message && <p className={`text-sm ${isError ? 'text-red-500' : 'text-green-600'}`}>{message}</p>}
						{!isError && <p className='text-sm text-muted-foreground'>You can now sign in with your verified email address.</p>}
						<Button asChild>
							<Link href='/login'>Go to Sign In</Link>
						</Button>
					</>
				)}
			</CardContent>
		</Card>
	)
}
