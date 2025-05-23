import React from 'react'
import {Button} from '@/components/ui/button'
import {Menu} from 'lucide-react'
import Breadcrumbs from './Breadcrumbs' // Import Breadcrumbs

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
			{/* Right side: Tenant Switcher and User Profile Section */}
			<div className='flex items-center space-x-4'>{/* <TenantSwitcher /> Add Tenant Switcher here */}</div>
		</header>
	)
}

export default Header
