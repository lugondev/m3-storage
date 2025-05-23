'use client'

import ProtectedWrapper from '@/components/auth/ProtectedWrapper'
import AppShell from '@/components/layout/AppShell'
import {useAuth} from '@/contexts/AuthContext'

export default function ProtectedLayout({children}: {children: React.ReactNode}) {
	const {isSystemAdmin} = useAuth()

	return (
		<ProtectedWrapper>
			<AppShell sidebarType={isSystemAdmin ? 'system' : 'user'}>{children}</AppShell> // User sidebar
		</ProtectedWrapper>
	)
}
