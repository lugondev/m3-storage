import React from 'react'
import {Button} from '@/components/ui/button'
import {Menu, User, LogOut, Settings} from 'lucide-react'
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger} from '@/components/ui/dropdown-menu'
import Link from 'next/link'
import Breadcrumbs from './Breadcrumbs' // Import Breadcrumbs
// TODO: Add sign out functionality from AuthContext

interface HeaderProps {
	onMenuButtonClick: () => void
}

const Header: React.FC<HeaderProps> = ({onMenuButtonClick}) => {
	return (
		<header className='sticky top-0 z-10 flex h-16 flex-shrink-0 items-center justify-between border-b border-gray-200 bg-white px-4 dark:border-gray-700 dark:bg-gray-800 sm:px-6 lg:px-8'>
			{/* Left side: Mobile Menu Button and Breadcrumbs/Title */}
			<div className='flex items-center'>
				<Button variant='ghost' size='icon' className='mr-4 md:hidden' onClick={onMenuButtonClick} aria-label='Toggle menu'>
					<Menu className='h-6 w-6' />
				</Button>

				{/* Breadcrumbs */}
				<div className='hidden md:block'>
					<Breadcrumbs />
				</div>
			</div>
			{/* Right side: User Profile Section */}
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant='ghost' size='icon' aria-label='User profile'>
						<User className='h-6 w-6 text-gray-600 dark:text-gray-300' />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align='end'>
					<DropdownMenuLabel>My Account</DropdownMenuLabel>
					<DropdownMenuSeparator />
					<DropdownMenuItem asChild>
						<Link href='/profile'>
							<User className='mr-2 h-4 w-4' />
							<span>Profile</span>
						</Link>
					</DropdownMenuItem>
					<DropdownMenuItem asChild>
						<Link href='/settings'>
							<Settings className='mr-2 h-4 w-4' />
							<span>Settings</span>
						</Link>
					</DropdownMenuItem>
					<DropdownMenuSeparator />
					<DropdownMenuItem onClick={() => alert('Sign out clicked!')}>
						{' '}
						{/* TODO: Implement actual sign out */}
						<LogOut className='mr-2 h-4 w-4' />
						<span>Sign out</span>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</header>
	)
}

export default Header
