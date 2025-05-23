'use client'

import React from 'react'
import {useRouter} from 'next/navigation'
import {Button} from '@/components/ui/button'
import {ArrowLeft} from 'lucide-react' // Or ArrowLeftIcon from radix

interface PageHeaderProps {
	title: string
	description?: string
	backButton?: {
		href?: string
		text: string
		onClick?: () => void
	}
	actions?: React.ReactNode
}

export function PageHeader({title, description, backButton, actions}: PageHeaderProps) {
	const router = useRouter()

	const handleBackClick = () => {
		if (backButton?.onClick) {
			backButton.onClick()
		} else if (backButton?.href) {
			router.push(backButton.href)
		} else {
			router.back() // Default fallback
		}
	}

	return (
		<div className='mb-6'>
			<div className='flex items-center justify-between mb-2'>
				<div className='flex items-center gap-3'>
					{backButton && (
						<Button variant='outline' size='sm' onClick={handleBackClick}>
							<ArrowLeft className='mr-2 h-4 w-4' />
							{backButton.text}
						</Button>
					)}
					<h1 className='text-3xl font-bold'>{title}</h1>
				</div>
				{actions && <div className='flex items-center gap-2'>{actions}</div>}
			</div>
			{description && <p className='text-muted-foreground'>{description}</p>}
		</div>
	)
}
