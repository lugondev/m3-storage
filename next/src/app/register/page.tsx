'use client'

import {RegisterForm} from '@/components/auth/RegisterForm'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import Link from 'next/link'

export default function RegisterPage() {
	return (
		<div className='flex min-h-screen items-center justify-center bg-gray-100 dark:bg-gray-900'>
			<Card className='w-full max-w-md mx-4'>
				<CardHeader className='text-center'>
					<CardTitle>Create an Account</CardTitle>
					<CardDescription>Enter your details below to register.</CardDescription>
				</CardHeader>
				<CardContent>
					<RegisterForm />
					<div className='mt-4 text-center text-sm'>
						Already have an account?{' '}
						<Link href='/' className='underline'>
							Sign in
						</Link>
					</div>
				</CardContent>
			</Card>
		</div>
	)
}
