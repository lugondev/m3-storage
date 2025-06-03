'use client' // Needed for usePathname hook

import React, {useState} from 'react' // Import useState
import Link from 'next/link'
import {usePathname} from 'next/navigation'
import {cn} from '@/lib/utils'
import {Building, Settings, BarChart, LayoutDashboard, ChevronDown, CreditCard} from 'lucide-react' // Removed Grid3x3 icon, Added ShieldCheck
import {Collapsible, CollapsibleContent, CollapsibleTrigger} from '@/components/ui/collapsible' // Import Collapsible

interface NavItem {
	href?: string // Optional href for parent items
	label: string
	icon: React.ElementType
	children?: NavItem[] // Optional children for collapsible sections
	disabled?: boolean // Optional disabled state
}

interface SidebarProps {
	isOpen: boolean // For mobile state
}

// TODO: Move this to a config file or context later for better management
const navItems: NavItem[] = [
	{href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard},
	{href: '/dashboard/venues', label: 'Venues', icon: Building},
	{href: '/dashboard/analytics', label: 'Analytics', icon: BarChart, disabled: true}, // Example disabled item
	{
		// Collapsible Settings section
		label: 'Settings',
		icon: Settings,
		children: [{href: '/dashboard/settings/billing', label: 'Billing', icon: CreditCard, disabled: true}],
	},
]

const Sidebar: React.FC<SidebarProps> = ({isOpen}) => {
	const pathname = usePathname()
	const [openSections, setOpenSections] = useState<Record<string, boolean>>(() => {
		// Initialize open sections based on current path
		const initialOpen: Record<string, boolean> = {}
		navItems.forEach((item) => {
			if (item.children && item.children.some((child) => child.href && pathname.startsWith(child.href))) {
				initialOpen[item.label] = true
			}
		})
		return initialOpen
	})

	const toggleSection = (label: string) => {
		setOpenSections((prev) => ({...prev, [label]: !prev[label]}))
	}

	return (
		<aside className={cn('fixed inset-y-0 left-0 z-30 w-64 transform border-r border-gray-200 bg-white p-4 transition-transform duration-300 ease-in-out dark:border-gray-700 dark:bg-gray-800 md:static md:z-auto md:translate-x-0', isOpen ? 'translate-x-0' : '-translate-x-full')} aria-label='Sidebar'>
			<div className='mb-8 flex items-center justify-center'>
				{/* Placeholder Logo/Title - Replace with actual logo if available */}
				<Link href='/dashboard' className='text-2xl font-semibold text-gray-800 dark:text-gray-100'>
					VenueApp
				</Link>
			</div>
			<nav className='flex-1 space-y-1'>
				{navItems.map((item) => {
					const isSectionOpen = openSections[item.label] ?? false
					// Check if the parent itself or any child is active
					const isParentActive = (item.href && pathname === item.href) || (item.href && item.href !== '/dashboard' && pathname.startsWith(item.href))
					const isChildActive = item.children?.some((child) => child.href && pathname.startsWith(child.href)) ?? false
					const isActive = isParentActive || isChildActive

					return item.children ? (
						<Collapsible key={item.label} open={isSectionOpen} onOpenChange={() => toggleSection(item.label)} className='space-y-1'>
							<CollapsibleTrigger asChild>
								<button className={cn('flex w-full items-center justify-between rounded-md px-3 py-2 text-sm font-medium transition-colors duration-150 ease-in-out', isActive ? 'bg-gray-100 text-gray-900 dark:bg-gray-700 dark:text-white' : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700', item.disabled && 'cursor-not-allowed opacity-50')} disabled={item.disabled}>
									<div className='flex items-center'>
										<item.icon className={cn('mr-3 h-5 w-5 flex-shrink-0', isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:text-gray-500 dark:group-hover:text-gray-400')} aria-hidden='true' />
										{item.label}
									</div>
									<ChevronDown className={cn('h-4 w-4 transform transition-transform duration-200', isSectionOpen ? 'rotate-180' : '')} />
								</button>
							</CollapsibleTrigger>
							<CollapsibleContent className='space-y-1 pl-7'>
								{item.children.map((child) => {
									const isChildLinkActive = child.href && pathname.startsWith(child.href)
									return (
										<Link
											key={child.label}
											// Use '#' if href is missing
											href={child.href ?? '#'}
											className={cn('flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors duration-150 ease-in-out', isChildLinkActive ? 'bg-gray-100 text-gray-900 dark:bg-gray-700 dark:text-white' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white', child.disabled && 'cursor-not-allowed opacity-50')}
											aria-disabled={child.disabled}
											tabIndex={child.disabled ? -1 : undefined}
											// Prevent navigation for disabled items
											onClick={(e) => child.disabled && e.preventDefault()}>
											<child.icon className={cn('mr-3 h-4 w-4 flex-shrink-0', isChildLinkActive ? 'text-indigo-500 dark:text-indigo-300' : 'text-gray-400 group-hover:text-gray-500 dark:text-gray-500 dark:group-hover:text-gray-300')} aria-hidden='true' />
											{child.label}
										</Link>
									)
								})}
							</CollapsibleContent>
						</Collapsible>
					) : (
						// Render non-collapsible item
						// Ensure href exists before rendering Link
						item.href && (
							<Link
								key={item.label}
								href={item.href}
								aria-disabled={item.disabled}
								tabIndex={item.disabled ? -1 : undefined}
								// Prevent navigation for disabled items
								onClick={(e) => item.disabled && e.preventDefault()}
								className={cn('flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors duration-150 ease-in-out', isActive ? 'bg-gray-100 font-semibold text-gray-900 dark:bg-gray-700 dark:text-white' : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700', item.disabled && 'cursor-not-allowed opacity-50')}>
								{/* Icon and label now direct children */}
								<item.icon className={cn('mr-3 h-5 w-5 flex-shrink-0', isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:text-gray-500 dark:group-hover:text-gray-400')} aria-hidden='true' />
								{item.label}
							</Link>
						)
					)
				})}
			</nav>
		</aside>
	)
}

// Exporting as named export to match the import suggestion
export {Sidebar}
