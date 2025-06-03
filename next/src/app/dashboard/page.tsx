import React from 'react'
import PageContainer from '@/components/layout/PageContainer'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card' // Assuming shadcn card

// TODO: Fetch actual data for these sections
const stats = [
	{title: 'Active Venues', value: '12', change: '+2 since last month'},
	{title: 'Upcoming Events', value: '5', change: ''},
	{title: 'Total Staff', value: '45', change: '+5 new hires'},
	{title: 'Open Alerts', value: '3', change: 'High priority'},
]

const recentActivity = [
	{id: 1, description: 'New event "Summer Gala" created for "Grand Hall"', time: '2 hours ago'},
	{id: 2, description: 'Table layout updated for "Riverside Cafe"', time: '5 hours ago'},
	{id: 3, description: 'User "Alice" added to staff at "Main Arena"', time: '1 day ago'},
	{id: 4, description: 'Product "Craft Beer" stock updated for "Downtown Pub"', time: '2 days ago'},
]

const DashboardPage = () => {
	return (
		<PageContainer>
			<h1 className='text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-6'>Dashboard Overview</h1>

			{/* Overview Statistics */}
			<div className='grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6'>
				{stats.map((stat) => (
					<Card key={stat.title}>
						<CardHeader className='pb-2'>
							<CardDescription>{stat.title}</CardDescription>
							<CardTitle className='text-4xl'>{stat.value}</CardTitle>
						</CardHeader>
						<CardContent>
							<p className='text-xs text-muted-foreground'>{stat.change}</p>
						</CardContent>
					</Card>
				))}
			</div>

			<div className='grid grid-cols-1 gap-6 lg:grid-cols-3'>
				{/* Recent Activity */}
				<Card className='lg:col-span-2'>
					<CardHeader>
						<CardTitle>Recent Activity</CardTitle>
						<CardDescription>Latest updates across your venues.</CardDescription>
					</CardHeader>
					<CardContent>
						<ul className='space-y-4'>
							{recentActivity.map((activity) => (
								<li key={activity.id} className='flex items-start space-x-3'>
									<div className='flex-shrink-0 pt-1'>
										{/* Placeholder Icon */}
										<div className='h-2 w-2 rounded-full bg-blue-500'></div>
									</div>
									<div>
										<p className='text-sm text-gray-800 dark:text-gray-200'>{activity.description}</p>
										<p className='text-xs text-muted-foreground'>{activity.time}</p>
									</div>
								</li>
							))}
						</ul>
						{/* TODO: Add link to full activity log */}
					</CardContent>
				</Card>

				{/* Quick Access & Alerts */}
				<div className='space-y-6'>
					{/* Quick Access */}
					<Card>
						<CardHeader>
							<CardTitle>Quick Access</CardTitle>
						</CardHeader>
						<CardContent className='flex flex-col space-y-2'>
							{/* TODO: Replace with actual links/buttons */}
							<button className='text-sm text-blue-600 hover:underline dark:text-blue-400'>Create New Venue</button>
							<button className='text-sm text-blue-600 hover:underline dark:text-blue-400'>Manage Staff</button>
							<button className='text-sm text-blue-600 hover:underline dark:text-blue-400'>View Upcoming Events</button>
							<button className='text-sm text-blue-600 hover:underline dark:text-blue-400'>Go to Settings</button>
						</CardContent>
					</Card>

					{/* Alerts */}
					<Card>
						<CardHeader>
							<CardTitle>Alerts & Notifications</CardTitle>
						</CardHeader>
						<CardContent>
							{/* TODO: Display actual alerts */}
							<p className='text-sm text-red-600 dark:text-red-400'>Low stock warning for Red Wine</p>
							<p className='text-sm text-orange-500 dark:text-orange-400'>Staff shift conflict detected</p>
							{/* TODO: Add link to view all alerts */}
						</CardContent>
					</Card>
				</div>
			</div>
		</PageContainer>
	)
}

export default DashboardPage
