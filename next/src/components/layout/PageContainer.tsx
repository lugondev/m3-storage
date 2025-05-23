import React from 'react'
import {cn} from '@/lib/utils' // Import cn for conditional classes if needed later

interface PageContainerProps {
	children: React.ReactNode
	className?: string // Allow consumers to add extra classes
}

const PageContainer: React.FC<PageContainerProps> = ({children, className}) => {
	return <div className={cn('container mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8', className)}>{children}</div>
}

export default PageContainer
