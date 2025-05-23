'use client'

import {Inter} from 'next/font/google'
import './globals.css'
import {AuthProvider} from '@/components/providers/AuthProvider'
import {AuthStatus} from '@/components/auth/AuthStatus'
import {ThemeProvider} from '@/components/providers/theme-provider'
import {QueryProvider} from '@/components/providers/QueryProvider' // Added QueryProvider import
import {ModeToggle} from '@/components/ui/mode-toggle'
import Link from 'next/link' // Import Link for logo
import {Toaster as SonnerToaster} from '@/components/ui/sonner' // Correct import for sonner

const inter = Inter({
	subsets: ['latin'],
	variable: '--font-inter',
})

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode
}>) {
	return (
		<html lang='en' suppressHydrationWarning>
			<body className={`${inter.variable} font-sans antialiased`}>
				<ThemeProvider attribute='class' defaultTheme='dark' enableSystem disableTransitionOnChange>
					<QueryProvider>
						{' '}
						{/* Wrap AuthProvider with QueryProvider */}
						<AuthProvider>
							{/* Updated Header Structure */}
							<header className='sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60'>
								<div className='container flex h-14 max-w-screen-2xl items-center'>
									{/* Left side: Logo/Brand */}
									<div className='mr-4 flex'>
										{' '}
										{/* Removed hidden md:flex for now */}
										<Link href='/' className='mr-6 flex items-center space-x-2'>
											{/* Optional Icon: <AppWindow className='h-6 w-6' /> */}
											<span className='font-bold sm:inline-block'>
												{' '}
												{/* Removed hidden */}
												AuthSys
											</span>
										</Link>
										{/* Main navigation links can go here later */}
									</div>

									{/* Mobile Nav Toggle can go here later */}

									{/* Right side: Theme toggle and Auth Status */}
									<div className='flex flex-1 items-center justify-end space-x-2'>
										<ModeToggle />
										<AuthStatus /> {/* AuthStatus now includes the dropdown */}
									</div>
								</div>
							</header>
							{/* Main content area below the sticky header */}
							<main className='flex-1'>{children}</main> {/* Removed pt-16 */}
							<SonnerToaster position='bottom-right' /> {/* Use Sonner Toaster */}
						</AuthProvider>
					</QueryProvider>{' '}
					{/* Close QueryProvider */}
				</ThemeProvider>
			</body>
		</html>
	)
}
