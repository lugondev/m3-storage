'use client' // Added for using hooks

import React from 'react'
import Header from './Header'
import Sidebar from './Sidebar'
import {useAuth} from '@/contexts/AuthContext' // Import useAuth

interface AppShellProps {
	children: React.ReactNode
	sidebarType?: 'system' | 'user' // Added 'user'
	tenantId?: string
	tenantName?: string
}

const AppShell: React.FC<AppShellProps> = ({children, sidebarType: propSidebarType, tenantId, tenantName}) => {
	const [isSidebarOpen, setIsSidebarOpen] = React.useState(true)
	const {isAuthenticated, isSystemAdmin, loading: authLoading} = useAuth()

	const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen)

	// Determine the actual sidebar type
	let actualSidebarType: 'system' | 'user' | undefined = propSidebarType

	if (!propSidebarType && isAuthenticated) {
		if (isSystemAdmin === true) {
			actualSidebarType = 'system'
		} else if (isSystemAdmin === false) {
			actualSidebarType = 'user'
		}
	}

	// If auth is loading, or user is not authenticated, we might not want to show the sidebar,
	// or show a specific version. For now, Sidebar handles undefined type by showing user links.
	// This part can be refined based on exact requirements for loading/unauthenticated states.
	if (authLoading) {
		// Optionally return a loading spinner or a minimal layout
		// For now, let it proceed, Sidebar will show default links if type is undefined.
	}

	// Do not render sidebar if not authenticated and not loading
	// (assuming AppShell is primarily for authenticated sections)
	// However, if a page using AppShell is public and wants a sidebar, this logic might need adjustment.
	// For this task, we assume AppShell is used in protected routes.
	const showSidebar = isAuthenticated || authLoading // Show sidebar if authenticated or auth is still loading (to avoid flicker)

	return (
		<div className='flex h-screen bg-gray-100 dark:bg-gray-900'>
			{showSidebar && (
				<>
					{/* Sidebar: Hidden on mobile by default, visible on md and larger */}
					<div className={`fixed inset-y-0 left-0 z-30 w-64 transform bg-gray-800 text-white transition-transform duration-300 ease-in-out md:relative md:translate-x-0 ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
						<Sidebar type={actualSidebarType} tenantId={tenantId} tenantName={tenantName} />
					</div>

					{/* Overlay for mobile when sidebar is open */}
					{isSidebarOpen && <div className='fixed inset-0 z-20 bg-black opacity-50 md:hidden' onClick={toggleSidebar}></div>}
				</>
			)}

			<div className='flex flex-1 flex-col overflow-hidden'>
				<Header onMenuButtonClick={toggleSidebar} />
				<main className='flex-1 overflow-y-auto overflow-x-hidden bg-gray-100 p-4 dark:bg-gray-900'>{children}</main>
			</div>
		</div>
	)
}

export default AppShell
