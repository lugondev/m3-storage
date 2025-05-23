// @/components/layout/Breadcrumb.tsx
'use client'

import React from 'react'
import Link from 'next/link'
import {usePathname} from 'next/navigation'

interface BreadcrumbItem {
	label: string
	href: string
}

interface BreadcrumbProps {
	items?: BreadcrumbItem[] // Optional, can be derived from pathname
}

const Breadcrumb: React.FC<BreadcrumbProps> = ({items}) => {
	const pathname = usePathname()

	// If items are not provided, derive them from the pathname
	const derivedItems = React.useMemo(() => {
		if (items) return items

		const pathSegments = pathname.split('/').filter((segment) => segment)
		const breadcrumbItems: BreadcrumbItem[] = [{label: 'Home', href: '/'}]

		let currentPath = ''
		pathSegments.forEach((segment) => {
			currentPath += `/${segment}`
			// Capitalize the first letter of the segment for display
			const label = segment.charAt(0).toUpperCase() + segment.slice(1)
			breadcrumbItems.push({label, href: currentPath})
		})
		return breadcrumbItems
	}, [items, pathname])

	return (
		<nav aria-label='breadcrumb'>
			<ol className='flex space-x-2 text-sm text-gray-500'>
				{derivedItems.map((item, index) => (
					<li key={item.href} className='flex items-center'>
						{index > 0 && <span className='mx-2'>/</span>}
						{index === derivedItems.length - 1 ? (
							<span className='font-medium text-gray-700'>{item.label}</span>
						) : (
							<Link href={item.href} className='hover:text-gray-700'>
								{item.label}
							</Link>
						)}
					</li>
				))}
			</ol>
		</nav>
	)
}

export default Breadcrumb
