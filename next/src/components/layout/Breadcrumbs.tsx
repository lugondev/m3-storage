'use client'

import React from 'react'
import Link from 'next/link'
import {usePathname} from 'next/navigation'
import {ChevronRight, Home} from 'lucide-react'
// cn import removed as it was unused

// Define the shape of a breadcrumb item
interface NavItem {
	href: string | null // Allow null for non-link items (like the last one)
	label: string
	icon: React.ElementType | null // Corrected: Allow icon to be null
}

// Helper function to capitalize first letter
const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1)

// Helper function to generate breadcrumb items from path
const generateBreadcrumbs = (pathname: string): NavItem[] => {
	// Added return type
	const pathSegments = pathname.split('/').filter((segment) => segment) // Remove empty strings
	const breadcrumbs: NavItem[] = [{href: '/', label: 'Home', icon: Home}] // Always start with Home, specify type

	let currentPath = ''
	pathSegments.forEach((segment) => {
		// Removed unused 'index'
		currentPath += `/${segment}`
		// Basic label generation, might need more complex logic for dynamic routes (e.g., venue IDs)
		const label = capitalize(segment.replace(/-/g, ' ')) // Replace hyphens and capitalize

		// TODO: Add logic to fetch dynamic labels (e.g., venue name for /venues/[venueId])
		// if (segment === '[venueId]' && venueData) { label = venueData.name; }

		// Modify the condition to create links for all segments for now
		// The rendering logic below handles not linking the last item.
		breadcrumbs.push({href: currentPath, label: label, icon: null})
	})

	// If the root is accessed directly, just show Home
	if (breadcrumbs.length === 1 && pathname === '/') {
		return breadcrumbs
	}
	// If only one segment after home (like /dashboard), show only that segment label
	if (breadcrumbs.length === 2) {
		return [{href: breadcrumbs[1].href, label: breadcrumbs[1].label, icon: Home}] // Show first section with Home icon
	}
	// If more than one segment, show Home > Segment1 > Segment2 ...
	// Keep the original breadcrumbs array structure

	return breadcrumbs
}

const Breadcrumbs: React.FC = () => {
	const pathname = usePathname()
	const breadcrumbItems = generateBreadcrumbs(pathname)

	// Handle the case where generateBreadcrumbs returns only one item for top-level sections
	if (breadcrumbItems.length === 1 && pathname !== '/') {
		const item = breadcrumbItems[0]
		return (
			<nav aria-label='Breadcrumb' className='flex items-center space-x-1 text-sm text-gray-500 dark:text-gray-400'>
				<div className='flex items-center'>
					{item.icon && <item.icon className='mr-1.5 h-4 w-4 flex-shrink-0' aria-hidden='true' />}
					<span className='font-medium text-gray-700 dark:text-gray-200'>{item.label}</span>
				</div>
			</nav>
		)
	}

	if (breadcrumbItems.length <= 1) return null // Don't show if only 'Home' and on root page

	return (
		<nav aria-label='Breadcrumb' className='flex items-center space-x-1 text-sm text-gray-500 dark:text-gray-400'>
			{breadcrumbItems.map((item, index) => (
				<React.Fragment key={item.href || item.label}>
					{index > 0 && <ChevronRight className='h-4 w-4 flex-shrink-0' aria-hidden='true' />}
					<div className='flex items-center'>
						{item.icon &&
							index === 0 && ( // Show Home icon only for the first item if it's Home
								<item.icon className='mr-1.5 h-4 w-4 flex-shrink-0' aria-hidden='true' />
							)}
						{item.href && index < breadcrumbItems.length - 1 ? ( // Make all but the last item links
							<Link href={item.href} className='hover:text-gray-700 dark:hover:text-gray-200'>
								{item.label}
							</Link>
						) : (
							// Last item is not a link, display as text
							<span className='font-medium text-gray-700 dark:text-gray-200'>{item.label}</span>
						)}
					</div>
				</React.Fragment>
			))}
		</nav>
	)
}

export default Breadcrumbs
