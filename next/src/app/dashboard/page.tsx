'use client'

import React from 'react'
import {useAuth} from '@/contexts/AuthContext'
import Link from 'next/link'
import {Button} from '@/components/ui/button'
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card'
import {Skeleton} from '@/components/ui/skeleton'

export default function UserDashboardPage() {
	const {user, isSystemAdmin, loading} = useAuth()

	if (loading) {
		return (
			<div className='container mx-auto p-4 md:p-6'>
				<Skeleton className='mb-8 h-10 w-48' />
				<Card className='mb-8'>
					<CardHeader>
						<Skeleton className='h-8 w-40' />
					</CardHeader>
					<CardContent className='space-y-4'>
						<Skeleton className='h-4 w-full' />
						<Skeleton className='h-4 w-3/4' />
						<Skeleton className='h-4 w-1/2' />
					</CardContent>
				</Card>
			</div>
		)
	}

	// ProtectedWrapper ensures user is available
	if (!user) return null

	return (
		<div className='container mx-auto p-4 md:p-6'>
			<h1 className='mb-8 text-3xl font-bold text-gray-800 dark:text-gray-100'>User Dashboard</h1>

			{/* Personal Information Section */}
			<Card className='mb-8'>
				<CardHeader>
					<CardTitle className='text-2xl font-semibold text-gray-700 dark:text-gray-200'>Personal Information</CardTitle>
				</CardHeader>
				<CardContent>
					<div className='space-y-2 text-gray-700 dark:text-gray-300'>
						<p>
							<strong>Name:</strong> {user.first_name && user.last_name ? `${user.first_name} ${user.last_name}` : user.first_name || user.last_name || user.email || 'N/A'}
						</p>
						<p>
							<strong>Email:</strong> {user.email || 'N/A'}
						</p>
						<p>
							<strong>User ID:</strong> {user.id}
						</p>
						{isSystemAdmin && <p className='mt-2 rounded-md bg-blue-100 p-2 text-sm text-blue-700 dark:bg-blue-900 dark:text-blue-300'>You have System Administrator privileges.</p>}
					</div>
				</CardContent>
			</Card>

			{/* System Administration Section (Conditional) */}
			{isSystemAdmin && (
				<Card className='mb-8'>
					<CardHeader>
						<CardTitle className='text-2xl font-semibold text-gray-700 dark:text-gray-200'>System Administration</CardTitle>
					</CardHeader>
					<CardContent>
						<p className='mb-4 text-gray-600 dark:text-gray-400'>Access global system settings and management tools.</p>
						<Button asChild variant='default' size='lg'>
							<Link href='/dashboard/admin/dashboard'>Go to System Admin</Link>
						</Button>
					</CardContent>
				</Card>
			)}
		</div>
	)
}
