'use client'

import React from 'react'

export default function ProfileLayout({children}: {children: React.ReactNode}) {
	// The parent layout (app/(protected)/layout.tsx) handles initial auth loading.
	// This layout provides the AppShell for the profile section.
	// AppShell internally uses AuthContext to determine the correct sidebar
	// (e.g., 'user' links for non-admins, 'system' links for system admins)
	// because no 'sidebarType' prop is explicitly passed here.
	return <>{children}</>
}
