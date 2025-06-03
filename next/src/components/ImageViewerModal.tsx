'use client'

import {useState, useEffect} from 'react'
import {Dialog, DialogContent, DialogHeader, DialogTitle} from '@/components/ui/dialog'
import {Button} from '@/components/ui/button'
import {Download, X, ZoomIn, ZoomOut, RotateCw} from 'lucide-react'
import {type MediaItem} from '@/lib/apiClient'

interface ImageViewerModalProps {
	item: MediaItem | null
	open: boolean
	onOpenChange: (open: boolean) => void
}

const ImageViewerModal = ({item, open, onOpenChange}: ImageViewerModalProps) => {
	// Don't render if no item is provided
	if (!item) {
		return null
	}
	const [zoom, setZoom] = useState(1)
	const [rotation, setRotation] = useState(0)

	const handleZoomIn = () => setZoom((prev) => Math.min(prev + 0.25, 3))
	const handleZoomOut = () => setZoom((prev) => Math.max(prev - 0.25, 0.25))
	const handleRotate = () => setRotation((prev) => (prev + 90) % 360)
	const handleDownload = () => {
		const link = document.createElement('a')
		link.href = item.public_url
		link.download = item.file_name
		link.click()
	}

	const resetTransforms = () => {
		setZoom(1)
		setRotation(0)
	}

	// Reset transforms when modal opens/closes
	useEffect(() => {
		if (open) {
			resetTransforms()
		}
	}, [open])

	// Reset transforms when modal opens
	const handleOpenChange = (newOpen: boolean) => {
		onOpenChange(newOpen)
	}

	// Format file size helper
	const formatFileSize = (bytes: number) => {
		if (bytes === 0) return '0 Bytes'
		const k = 1024
		const sizes = ['Bytes', 'KB', 'MB', 'GB']
		const i = Math.floor(Math.log(bytes) / Math.log(k))
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
	}

	const isImage = item.media_type?.startsWith('image')
	const isVideo = item.media_type?.startsWith('video')
	const isAudio = item.media_type?.startsWith('audio')
	const isPdf = item.media_type?.includes('pdf')

	return (
		<Dialog open={open} onOpenChange={handleOpenChange}>
			<DialogContent className='max-w-6xl max-h-[95vh] w-[95vw] p-0 overflow-hidden'>
				<DialogHeader className='p-4 pb-2 border-b'>
					<div className='flex items-center justify-between'>
						<DialogTitle className='text-lg font-semibold truncate pr-4'>{item.file_name}</DialogTitle>
						<div className='flex items-center gap-2'>
							{isImage && (
								<>
									<Button variant='outline' size='sm' onClick={handleZoomOut} disabled={zoom <= 0.25}>
										<ZoomOut className='h-4 w-4' />
									</Button>
									<span className='text-sm text-gray-600 dark:text-gray-300 min-w-[60px] text-center'>{Math.round(zoom * 100)}%</span>
									<Button variant='outline' size='sm' onClick={handleZoomIn} disabled={zoom >= 3}>
										<ZoomIn className='h-4 w-4' />
									</Button>
									<Button variant='outline' size='sm' onClick={handleRotate}>
										<RotateCw className='h-4 w-4' />
									</Button>
								</>
							)}
						</div>
					</div>
					<div className='text-sm text-gray-500 dark:text-gray-400'>
						{item.media_type} • {formatFileSize(item.file_size)} • {new Date(item.created_at).toLocaleDateString()}
					</div>
				</DialogHeader>

				<div className='flex-1 overflow-hidden bg-gray-50 dark:bg-gray-900' style={{height: 'calc(95vh - 120px)'}}>
					<div className='h-full flex items-center justify-center p-4'>
						{isImage ? (
							<div className='relative max-w-full max-h-full overflow-auto'>
								<img
									src={item.public_url}
									alt={item.file_name}
									className='max-w-none h-auto transition-transform duration-200'
									style={{
										transform: `scale(${zoom}) rotate(${rotation}deg)`,
										transformOrigin: 'center',
										maxHeight: zoom === 1 ? '100%' : 'none',
										maxWidth: zoom === 1 ? '100%' : 'none',
									}}
									onError={(e) => {
										console.error('Failed to load image:', item.public_url)
										e.currentTarget.style.display = 'none'
									}}
								/>
							</div>
						) : isVideo ? (
							<video
								controls
								className='max-w-full max-h-full'
								style={{maxHeight: 'calc(95vh - 160px)'}}
								onError={(e) => {
									console.error('Failed to load video:', item.public_url)
								}}>
								<source src={item.public_url} type={item.media_type} />
								Your browser does not support the video tag.
							</video>
						) : isAudio ? (
							<div className='text-center space-y-4'>
								<div className='text-6xl'>🎵</div>
								<h3 className='text-lg font-medium text-gray-900 dark:text-gray-100'>{item.file_name}</h3>
								<audio
									controls
									className='w-full max-w-md'
									onError={(e) => {
										console.error('Failed to load audio:', item.public_url)
									}}>
									<source src={item.public_url} type={item.media_type} />
									Your browser does not support the audio tag.
								</audio>
							</div>
						) : isPdf ? (
							<div className='w-full h-full'>
								<iframe
									src={item.public_url}
									className='w-full h-full border-0'
									title={item.file_name}
									onError={(e) => {
										console.error('Failed to load PDF:', item.public_url)
									}}
								/>
							</div>
						) : (
							<div className='text-center space-y-4'>
								<div className='text-6xl'>📄</div>
								<h3 className='text-lg font-medium text-gray-900 dark:text-gray-100'>{item.file_name}</h3>
								<p className='text-gray-600 dark:text-gray-400'>Preview not available for this file type</p>
								<Button onClick={handleDownload} className='mt-4'>
									<Download className='h-4 w-4 mr-2' />
									Download File
								</Button>
							</div>
						)}
					</div>
				</div>
			</DialogContent>
		</Dialog>
	)
}

export default ImageViewerModal
