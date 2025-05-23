'use client'

import {Button} from '@/components/ui/button'
import {Menu} from 'lucide-react'
import {Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger} from '@/components/ui/sheet'
import Sidebar from './Sidebar'

export function MobileHeader() {
	return (
		<div className='flex h-14 items-center border-b px-4 lg:hidden'>
			<Sheet>
				<SheetTrigger asChild>
					<Button variant='ghost' size='icon' className='mr-2'>
						<Menu className='h-5 w-5' />
						<span className='sr-only'>Toggle sidebar</span>
					</Button>
				</SheetTrigger>
				<SheetContent side='left' className='w-72 p-0'>
					<SheetHeader className='border-b p-4'>
						<SheetTitle>Menu</SheetTitle>
					</SheetHeader>
					<Sidebar />
				</SheetContent>
			</Sheet>
			<div className='flex-1'>
				<h1 className='text-lg font-semibold'>Auth System</h1>
			</div>
		</div>
	)
}
