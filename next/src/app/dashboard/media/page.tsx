'use client'

import {useState, useEffect, useCallback} from 'react'
import {useDropzone} from 'react-dropzone'
import PageContainer from '@/components/layout/PageContainer'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {Badge} from '@/components/ui/badge'
import {Alert, AlertDescription} from '@/components/ui/alert'
import {Skeleton} from '@/components/ui/skeleton'
import {uploadMedia, listMedia, deleteMedia, type MediaItem, type UploadResponse, type MediaListResponse} from '@/lib/apiClient'
import {Upload, X, Download, Eye, Trash2, Search, Filter} from 'lucide-react'
import {toast} from 'sonner'

const MediaManagementPage = () => {
	const [mediaItems, setMediaItems] = useState<MediaItem[]>([])
	const [loading, setLoading] = useState(true)
	const [uploading, setUploading] = useState(false)
	const [searchTerm, setSearchTerm] = useState('')
	const [selectedFiles, setSelectedFiles] = useState<File[]>([])
	const [uploadProgress, setUploadProgress] = useState<Record<string, number>>({})

	// Fetch media items
	const fetchMedia = async () => {
		try {
			setLoading(true)
			const response: MediaListResponse = await listMedia()
			setMediaItems(response.media || [])
		} catch (error) {
			console.error('Failed to fetch media:', error)
			toast.error('Failed to load media items')
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		fetchMedia()
	}, [])

	// File upload with drag and drop
	const onDrop = useCallback((acceptedFiles: File[]) => {
		setSelectedFiles((prev) => [...prev, ...acceptedFiles])
	}, [])

	const {getRootProps, getInputProps, isDragActive} = useDropzone({
		onDrop,
		accept: {
			'image/*': ['.jpeg', '.jpg', '.png', '.gif', '.webp'],
			'video/*': ['.mp4', '.avi', '.mov', '.wmv', '.flv'],
			'audio/*': ['.mp3', '.wav', '.ogg', '.m4a'],
			'application/pdf': ['.pdf'],
			'text/*': ['.txt', '.md'],
		},
		multiple: true,
	})

	// Upload files
	const handleUpload = async () => {
		if (selectedFiles.length === 0) return

		setUploading(true)
		const newUploadProgress: Record<string, number> = {}

		try {
			for (const file of selectedFiles) {
				const fileKey = `${file.name}-${file.size}`
				newUploadProgress[fileKey] = 0
				setUploadProgress({...newUploadProgress})

				try {
					const response: UploadResponse = await uploadMedia(file)
					newUploadProgress[fileKey] = 100
					setUploadProgress({...newUploadProgress})

					toast.success(`${file.name} uploaded successfully`)
				} catch (error) {
					console.error(`Failed to upload ${file.name}:`, error)
					toast.error(`Failed to upload ${file.name}`)
				}
			}

			// Refresh media list
			await fetchMedia()
			setSelectedFiles([])
			setUploadProgress({})
		} catch (error) {
			console.error('Upload error:', error)
			toast.error('Upload failed')
		} finally {
			setUploading(false)
		}
	}

	// Delete media item
	const handleDelete = async (id: string, filename: string) => {
		if (!confirm(`Are you sure you want to delete ${filename}?`)) return

		try {
			await deleteMedia(id)
			setMediaItems((prev) => prev.filter((item) => item.id !== id))
			toast.success(`${filename} deleted successfully`)
		} catch (error) {
			console.error('Failed to delete media:', error)
			toast.error('Failed to delete media item')
		}
	}

	// Remove file from upload queue
	const removeSelectedFile = (index: number) => {
		setSelectedFiles((prev) => prev.filter((_, i) => i !== index))
	}

	// Filter media items
	const filteredMedia = mediaItems.filter((item) => item.filename.toLowerCase().includes(searchTerm.toLowerCase()) || item.mime_type.toLowerCase().includes(searchTerm.toLowerCase()))

	// Format file size
	const formatFileSize = (bytes: number) => {
		if (bytes === 0) return '0 Bytes'
		const k = 1024
		const sizes = ['Bytes', 'KB', 'MB', 'GB']
		const i = Math.floor(Math.log(bytes) / Math.log(k))
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
	}

	// Get file type badge color
	const getFileTypeBadge = (mimeType: string) => {
		if (mimeType.startsWith('image/')) return 'bg-green-100 text-green-800'
		if (mimeType.startsWith('video/')) return 'bg-blue-100 text-blue-800'
		if (mimeType.startsWith('audio/')) return 'bg-purple-100 text-purple-800'
		if (mimeType.includes('pdf')) return 'bg-red-100 text-red-800'
		return 'bg-gray-100 text-gray-800'
	}

	return (
		<PageContainer>
			<div className='space-y-6'>
				<div className='flex items-center justify-between'>
					<h1 className='text-2xl font-semibold text-gray-900 dark:text-gray-100'>Media Management</h1>
					<Badge variant='outline'>{mediaItems.length} items</Badge>
				</div>

				{/* Upload Section */}
				<Card>
					<CardHeader>
						<CardTitle className='flex items-center gap-2'>
							<Upload className='h-5 w-5' />
							Upload Media
						</CardTitle>
						<CardDescription>Drag and drop files here or click to browse. Supports images, videos, audio, PDFs, and text files.</CardDescription>
					</CardHeader>
					<CardContent className='space-y-4'>
						{/* Dropzone */}
						<div {...getRootProps()} className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${isDragActive ? 'border-primary bg-primary/5' : 'border-gray-300 hover:border-gray-400'}`}>
							<input {...getInputProps()} />
							<Upload className='mx-auto h-12 w-12 text-gray-400 mb-4' />
							{isDragActive ? (
								<p className='text-primary'>Drop the files here...</p>
							) : (
								<div>
									<p className='text-gray-600 mb-2'>Drag & drop files here, or click to select</p>
									<p className='text-sm text-gray-400'>Maximum file size: 10MB</p>
								</div>
							)}
						</div>

						{/* Selected Files */}
						{selectedFiles.length > 0 && (
							<div className='space-y-2'>
								<h4 className='font-medium'>Selected Files ({selectedFiles.length})</h4>
								<div className='space-y-2 max-h-40 overflow-y-auto'>
									{selectedFiles.map((file, index) => {
										const fileKey = `${file.name}-${file.size}`
										const progress = uploadProgress[fileKey] || 0

										return (
											<div key={index} className='flex items-center justify-between p-2 bg-gray-50 rounded'>
												<div className='flex-1'>
													<p className='text-sm font-medium'>{file.name}</p>
													<p className='text-xs text-gray-500'>{formatFileSize(file.size)}</p>
													{uploading && progress > 0 && (
														<div className='w-full bg-gray-200 rounded-full h-1 mt-1'>
															<div className='bg-primary h-1 rounded-full transition-all' style={{width: `${progress}%`}} />
														</div>
													)}
												</div>
												<Button variant='ghost' size='sm' onClick={() => removeSelectedFile(index)} disabled={uploading}>
													<X className='h-4 w-4' />
												</Button>
											</div>
										)
									})}
								</div>
								<Button onClick={handleUpload} disabled={uploading || selectedFiles.length === 0} className='w-full'>
									{uploading ? 'Uploading...' : `Upload ${selectedFiles.length} file(s)`}
								</Button>
							</div>
						)}
					</CardContent>
				</Card>

				{/* Search and Filter */}
				<Card>
					<CardContent className='pt-6'>
						<div className='flex gap-4'>
							<div className='flex-1 relative'>
								<Search className='absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400' />
								<Input placeholder='Search by filename or file type...' value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} className='pl-10' />
							</div>
							<Button variant='outline'>
								<Filter className='h-4 w-4 mr-2' />
								Filter
							</Button>
						</div>
					</CardContent>
				</Card>

				{/* Media Grid */}
				<Card>
					<CardHeader>
						<CardTitle>Media Library</CardTitle>
						<CardDescription>
							{filteredMedia.length} of {mediaItems.length} items
						</CardDescription>
					</CardHeader>
					<CardContent>
						{loading ? (
							<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'>
								{Array.from({length: 8}).map((_, i) => (
									<div key={i} className='space-y-2'>
										<Skeleton className='h-32 w-full' />
										<Skeleton className='h-4 w-3/4' />
										<Skeleton className='h-3 w-1/2' />
									</div>
								))}
							</div>
						) : filteredMedia.length === 0 ? (
							<div className='text-center py-12'>
								<Upload className='mx-auto h-12 w-12 text-gray-400 mb-4' />
								<h3 className='text-lg font-medium text-gray-900 mb-2'>No media files</h3>
								<p className='text-gray-500 mb-4'>{searchTerm ? 'No files match your search.' : 'Upload your first media file to get started.'}</p>
								{searchTerm && (
									<Button variant='outline' onClick={() => setSearchTerm('')}>
										Clear search
									</Button>
								)}
							</div>
						) : (
							<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'>
								{filteredMedia.map((item) => (
									<Card key={item.id} className='overflow-hidden'>
										<div className='aspect-square bg-gray-100 flex items-center justify-center'>
											{item.mime_type.startsWith('image/') ? (
												<img src={item.url} alt={item.filename} className='w-full h-full object-cover' />
											) : (
												<div className='text-center'>
													<div className='text-4xl mb-2'>{item.mime_type.startsWith('video/') ? '🎥' : item.mime_type.startsWith('audio/') ? '🎵' : item.mime_type.includes('pdf') ? '📄' : '📁'}</div>
													<p className='text-xs text-gray-500'>{item.mime_type}</p>
												</div>
											)}
										</div>
										<CardContent className='p-3'>
											<div className='space-y-2'>
												<h4 className='font-medium text-sm truncate' title={item.filename}>
													{item.filename}
												</h4>
												<div className='flex items-center justify-between'>
													<Badge className={getFileTypeBadge(item.mime_type)}>{item.mime_type.split('/')[0]}</Badge>
													<span className='text-xs text-gray-500'>{formatFileSize(item.size)}</span>
												</div>
												<p className='text-xs text-gray-400'>{new Date(item.created_at).toLocaleDateString()}</p>
												<div className='flex gap-1'>
													<Button variant='outline' size='sm' className='flex-1' onClick={() => window.open(item.url, '_blank')}>
														<Eye className='h-3 w-3 mr-1' />
														View
													</Button>
													<Button
														variant='outline'
														size='sm'
														onClick={() => {
															const link = document.createElement('a')
															link.href = item.url
															link.download = item.filename
															link.click()
														}}>
														<Download className='h-3 w-3' />
													</Button>
													<Button variant='outline' size='sm' onClick={() => handleDelete(item.id, item.filename)} className='text-red-600 hover:text-red-700'>
														<Trash2 className='h-3 w-3' />
													</Button>
												</div>
											</div>
										</CardContent>
									</Card>
								))}
							</div>
						)}
					</CardContent>
				</Card>
			</div>
		</PageContainer>
	)
}

export default MediaManagementPage
