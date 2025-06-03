import React from 'react'
import Header from './Header' // Reverted: Removed .tsx extension
import {Sidebar} from './Sidebar' // Changed to named import

interface AppShellProps {
	children: React.ReactNode
}

const AppShell: React.FC<AppShellProps> = ({children}) => {
	// TODO: Add state for sidebar open/closed (for mobile)
	const [isSidebarOpen, setIsSidebarOpen] = React.useState(false) // Basic state for mobile

	const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen)

	return (
		<div className='flex h-screen bg-gray-100 dark:bg-gray-900'>
			<Sidebar isOpen={isSidebarOpen} />
			{/* Optional: Add overlay for mobile when sidebar is open */}
			{isSidebarOpen && <div className='fixed inset-0 z-20 bg-black opacity-50 md:hidden' onClick={toggleSidebar}></div>}
			<div className='flex flex-1 flex-col overflow-hidden'>
				<Header onMenuButtonClick={toggleSidebar} />
				<main className='flex-1 overflow-y-auto overflow-x-hidden bg-gray-100 p-4 dark:bg-gray-900'>{children}</main>
			</div>
		</div>
	)
}

export default AppShell
