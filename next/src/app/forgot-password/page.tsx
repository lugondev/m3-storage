'use client'

import {ForgotPasswordForm} from '@/components/auth/ForgotPasswordForm'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import Link from 'next/link'

export default function ForgotPasswordPage() {
	return (
		<div className='flex min-h-screen items-center justify-center bg-gray-100 dark:bg-gray-900'>
			<Card className='w-full max-w-md mx-4'>
				<CardHeader className='text-center'>
					<CardTitle>Forgot Password</CardTitle>
					<CardDescription>Enter your email address to receive reset instructions.</CardDescription>
				</CardHeader>
				<CardContent>
					<ForgotPasswordForm />
					<div className='mt-4 text-center text-sm'>
						Remembered your password?{' '}
						<Link href='/login' className='underline'>
							Sign in
						</Link>
					</div>
					<div className='mt-2 text-center text-sm'>
						Don{'\u0027'}t have an account?{' '}
						<Link href='/register' className='underline'>
							Register
						</Link>
					</div>
				</CardContent>
			</Card>
		</div>
	)
}
