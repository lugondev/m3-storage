'use client'

import {useAuth} from '@/contexts/AuthContext'
import {Button} from '@/components/ui/button'
import {Skeleton} from '@/components/ui/skeleton'
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger} from '@/components/ui/dropdown-menu'
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar'
import Link from 'next/link' // Import Link for navigation

export function AuthStatus() {
	const {user, isAuthenticated, logout, loading} = useAuth()

	if (loading) {
		// Show skeleton loaders while checking auth status
		return (
			<div className='flex items-center gap-2'>
				<Skeleton className='h-8 w-8 rounded-full' /> {/* Skeleton for Avatar */}
			</div>
		)
	}

	if (!isAuthenticated || !user) {
		// User is not logged in, show login buttons
		return
	}

	// User is logged in, show dropdown menu with user info and actions
	const displayName = user.first_name || user.last_name ? `${user.first_name || ''} ${user.last_name || ''}`.trim() : user.email
	// Use first character of display name or 'U' as fallback
	const fallbackInitials = displayName ? displayName.charAt(0).toUpperCase() : 'U'

	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant='ghost' className='relative h-8 w-8 rounded-full'>
					<Avatar className='h-8 w-8'>
						<AvatarImage src={user.avatar || ''} alt={displayName || 'User Avatar'} />
						<AvatarFallback>{fallbackInitials}</AvatarFallback>
					</Avatar>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className='w-56' align='end' forceMount>
				<DropdownMenuLabel className='font-normal'>
					<div className='flex flex-col space-y-1'>
						<p className='text-sm font-medium leading-none'>{displayName}</p>
						<p className='text-xs leading-none text-muted-foreground'>{user.email}</p>
					</div>
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<Link href='/dashboard/profile' passHref>
					<DropdownMenuItem>Profile</DropdownMenuItem>
				</Link>
				<DropdownMenuSeparator />
				<DropdownMenuItem onClick={logout} disabled={loading}>
					{loading ? 'Signing out...' : 'Sign out'}
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	)
}
