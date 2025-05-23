'use client'

import {useEffect, useState, useRef, ChangeEvent} from 'react'
import {useAuth} from '@/contexts/AuthContext'
import {getCurrentUser, updateCurrentUser, updateCurrentUserProfile, updateCurrentUserPassword, updateUserAvatar} from '@/services/userService' // Changed to updateUserAvatar
import {UserOutput, UserProfile, UpdateUserInput, UpdateProfileInput, UpdatePasswordInput} from '@/lib/apiClient' // Ensure these types match your generated client
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import {Skeleton} from '@/components/ui/skeleton'
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar'
import {Tabs, TabsContent, TabsList, TabsTrigger} from '@/components/ui/tabs'
import {Button} from '@/components/ui/button'
import {UploadCloudIcon} from 'lucide-react'
import {Input} from '@/components/ui/input'
import {Textarea} from '@/components/ui/textarea'
import {Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form'
import {useForm} from 'react-hook-form'
import {zodResolver} from '@hookform/resolvers/zod'
import * as z from 'zod'
import {toast} from 'sonner'
import {Popover, PopoverContent, PopoverTrigger} from '@/components/ui/popover'
import {CalendarIcon} from 'lucide-react'
import {Calendar} from '@/components/ui/calendar'
import {cn} from '@/lib/utils'
import {Checkbox} from '@/components/ui/checkbox'

// --- Import Security Tab Components ---
import EmailVerificationStatus from '@/features/profile/components/security/EmailVerificationStatus'
import PhoneVerification from '@/features/profile/components/security/PhoneVerification'
import TwoFactorAuthManagement from '@/features/profile/components/security/TwoFactorAuthManagement'

// Helper function to get initials
const getInitials = (firstName?: string, lastName?: string, email?: string): string => {
	if (firstName && lastName) {
		return `${firstName[0]}${lastName[0]}`.toUpperCase()
	}
	if (firstName) {
		return firstName.substring(0, 2).toUpperCase()
	}
	if (email) {
		return email.substring(0, 2).toUpperCase()
	}
	return 'U'
}

// Zod Schemas for Forms
const generalInfoSchema = z.object({
	first_name: z.string().min(1, 'First name is required').optional().or(z.literal('')),
	last_name: z.string().min(1, 'Last name is required').optional().or(z.literal('')),
	phone: z.string().optional().or(z.literal('')),
	// Status update might need specific logic/permissions, omitted for now
})

// Define the base password schema including otp_code and password match refinement
const passwordSchema = z
	.object({
		current_password: z.string().min(1, 'Current password is required'),
		new_password: z.string().min(8, 'New password must be at least 8 characters long'),
		confirm_new_password: z.string(),
		otp_code: z.string().optional(), // Optional by default
	})
	.refine((data) => data.new_password === data.confirm_new_password, {
		message: 'New passwords do not match',
		path: ['confirm_new_password'],
	})

const profileDetailsSchema = z.object({
	address: z.string().optional().or(z.literal('')),
	bio: z.string().optional().or(z.literal('')),
	date_of_birth: z.date().optional(),
	interests: z.array(z.string()).optional(), // Needs a specific input component (e.g., tags input)
	preferences: z
		.object({
			email_notifications: z.boolean().optional(),
			language: z.string().optional().or(z.literal('')),
			push_notifications: z.boolean().optional(),
			theme: z.string().optional().or(z.literal('')), // Needs a select component
		})
		.optional(),
})

// Infer type from the base password schema
type PasswordFormData = z.infer<typeof passwordSchema>

type GeneralInfoFormData = z.infer<typeof generalInfoSchema>
type ProfileDetailsFormData = z.infer<typeof profileDetailsSchema>

// --- Form Components ---

// --- Add Props for PasswordForm ---
interface PasswordFormProps {
	userData: UserOutput | null
}

interface GeneralInfoFormProps {
	userData: UserOutput | null
	onSubmitSuccess: (updatedUser: UserOutput) => void
}

function GeneralInfoForm({userData, onSubmitSuccess}: GeneralInfoFormProps) {
	const form = useForm<GeneralInfoFormData>({
		resolver: zodResolver(generalInfoSchema),
		defaultValues: {
			first_name: userData?.first_name || '',
			last_name: userData?.last_name || '',
			phone: userData?.phone || '',
		},
	})

	const [isSubmitting, setIsSubmitting] = useState(false)

	useEffect(() => {
		// Reset form when userData changes
		form.reset({
			first_name: userData?.first_name || '',
			last_name: userData?.last_name || '',
			phone: userData?.phone || '',
		})
	}, [userData, form])

	async function onSubmit(values: GeneralInfoFormData) {
		setIsSubmitting(true)
		const updateData: UpdateUserInput = {}
		if (values.first_name !== undefined) updateData.first_name = values.first_name
		if (values.last_name !== undefined) updateData.last_name = values.last_name
		if (values.phone !== undefined) updateData.phone = values.phone
		// Status is intentionally omitted here, could be added based on role/logic

		try {
			const updatedUser = await updateCurrentUser(updateData)
			toast.success('General information updated successfully.')
			onSubmitSuccess(updatedUser) // Callback to update parent state
		} catch (error: unknown) {
			console.error('Error updating general info:', error)
			const message = error instanceof Error ? error.message : 'Failed to update information.'
			toast.error(`Update failed: ${message}`)
		} finally {
			setIsSubmitting(false)
		}
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
				<FormField
					control={form.control}
					name='first_name'
					render={({field}) => (
						<FormItem>
							<FormLabel>First Name</FormLabel>
							<FormControl>
								<Input placeholder='Enter first name' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='last_name'
					render={({field}) => (
						<FormItem>
							<FormLabel>Last Name</FormLabel>
							<FormControl>
								<Input placeholder='Enter last name' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='phone'
					render={({field}) => (
						<FormItem>
							<FormLabel>Phone</FormLabel>
							<FormControl>
								<Input placeholder='Enter phone number' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<Button type='submit' disabled={isSubmitting || !form.formState.isDirty}>
					{isSubmitting ? 'Saving...' : 'Save Changes'}
				</Button>
			</form>
		</Form>
	)
}

// Update PasswordForm to accept userData
function PasswordForm({userData}: PasswordFormProps) {
	const is2FAEnabled = userData?.is_two_factor_enabled ?? false

	// Conditionally add the OTP refinement if 2FA is enabled
	const finalPasswordSchema = is2FAEnabled
		? passwordSchema.refine((data) => !!data.otp_code && data.otp_code.trim().length > 0, {
				message: 'OTP code is required when 2FA is enabled',
				path: ['otp_code'],
		  })
		: passwordSchema // Use the base schema if 2FA is not enabled

	// Use the final schema for the resolver
	const form = useForm<PasswordFormData>({
		resolver: zodResolver(finalPasswordSchema), // Use the potentially refined schema
		defaultValues: {
			current_password: '',
			new_password: '',
			confirm_new_password: '',
			otp_code: '', // Add default for otp_code
		},
	})

	const [isSubmitting, setIsSubmitting] = useState(false)

	// Watch 2FA status in case it changes (though unlikely within this form's lifecycle)
	useEffect(() => {
		// Re-validate or reset based on schema change if needed,
		// but zodResolver should handle this based on the schema passed during initialization.
		// If userData could change *while* the form is mounted, more complex handling might be needed.
	}, [is2FAEnabled, form])

	async function onSubmit(values: PasswordFormData) {
		// Double-check OTP requirement based on current userData state
		if (is2FAEnabled && (!values.otp_code || values.otp_code.trim() === '')) {
			form.setError('otp_code', {
				type: 'manual',
				message: 'OTP code is required when 2FA is enabled.',
			})
			return // Prevent submission
		}

		setIsSubmitting(true)
		const updateData: UpdatePasswordInput = {
			current_password: values.current_password,
			new_password: values.new_password,
			// Conditionally include otp_code
			...(is2FAEnabled && values.otp_code && {otp_code: values.otp_code}),
		}

		try {
			await updateCurrentUserPassword(updateData)
			toast.success('Password updated successfully.')
			form.reset() // Clear form on success
		} catch (error: unknown) {
			console.error('Error updating password:', error)
			// Check if the error response has a specific message
			let message = 'Failed to update password.'
			// Use unknown and type guard for better type safety
			if (error && typeof error === 'object' && 'response' in error && error.response) {
				const responseError = error.response as {data?: {error?: string} | string}
				if (responseError.data && typeof responseError.data === 'object' && responseError.data.error) {
					message = responseError.data.error
				} else if (typeof responseError.data === 'string') {
					message = responseError.data
				}
			} else if (error instanceof Error) {
				message = error.message
			}
			toast.error(`Update failed: ${message}`)
		} finally {
			setIsSubmitting(false)
		}
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
				<FormField
					control={form.control}
					name='current_password'
					render={({field}) => (
						<FormItem>
							<FormLabel>Current Password</FormLabel>
							<FormControl>
								<Input type='password' placeholder='Enter current password' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='new_password'
					render={({field}) => (
						<FormItem>
							<FormLabel>New Password</FormLabel>
							<FormControl>
								<Input type='password' placeholder='Enter new password' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='confirm_new_password'
					render={({field}) => (
						<FormItem>
							<FormLabel>Confirm New Password</FormLabel>
							<FormControl>
								<Input type='password' placeholder='Confirm new password' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				{/* Conditionally render OTP input */}
				{is2FAEnabled && (
					<FormField
						control={form.control}
						name='otp_code'
						render={({field}) => (
							<FormItem>
								<FormLabel>Two-Factor Authentication Code</FormLabel>
								<FormControl>
									<Input type='text' placeholder='Enter your 6-digit code' {...field} maxLength={6} />
								</FormControl>
								<FormDescription>Enter the code from your authenticator app.</FormDescription>
								<FormMessage />
							</FormItem>
						)}
					/>
				)}
				<Button type='submit' disabled={isSubmitting || !form.formState.isDirty}>
					{isSubmitting ? 'Updating...' : 'Update Password'}
				</Button>
			</form>
		</Form>
	)
}

interface ProfileDetailsFormProps {
	profileData: UserProfile | null
	onSubmitSuccess: (updatedProfile: UserProfile) => void
}

function ProfileDetailsForm({profileData, onSubmitSuccess}: ProfileDetailsFormProps) {
	const form = useForm<ProfileDetailsFormData>({
		resolver: zodResolver(profileDetailsSchema),
		defaultValues: {
			address: profileData?.address || '',
			bio: profileData?.bio || '',
			date_of_birth: profileData?.date_of_birth ? new Date(profileData.date_of_birth) : undefined,
			interests: profileData?.interests || [],
			preferences: {
				email_notifications: profileData?.preferences?.email_notifications ?? true, // Default to true if undefined
				language: profileData?.preferences?.language || '',
				push_notifications: profileData?.preferences?.push_notifications ?? true, // Default to true if undefined
				theme: profileData?.preferences?.theme || '',
			},
		},
	})

	const [isSubmitting, setIsSubmitting] = useState(false)

	useEffect(() => {
		form.reset({
			address: profileData?.address || '',
			bio: profileData?.bio || '',
			date_of_birth: profileData?.date_of_birth ? new Date(profileData.date_of_birth) : undefined,
			interests: profileData?.interests || [],
			preferences: {
				email_notifications: profileData?.preferences?.email_notifications ?? true,
				language: profileData?.preferences?.language || '',
				push_notifications: profileData?.preferences?.push_notifications ?? true,
				theme: profileData?.preferences?.theme || '',
			},
		})
	}, [profileData, form])

	async function onSubmit(values: ProfileDetailsFormData) {
		setIsSubmitting(true)

		// Construct the update payload carefully, handling optional fields
		const updateData: UpdateProfileInput = {
			// Only include fields if they have a value or are explicitly meant to be cleared
			...(values.address !== undefined && {address: values.address}),
			...(values.bio !== undefined && {bio: values.bio}),
			...(values.date_of_birth && {date_of_birth: values.date_of_birth.toISOString()}), // Send as ISO string
			...(values.interests && {interests: values.interests}), // Assuming backend accepts array directly
			...(values.preferences && {
				preferences: {
					...(values.preferences.email_notifications !== undefined && {email_notifications: values.preferences.email_notifications}),
					...(values.preferences.language !== undefined && {language: values.preferences.language}),
					...(values.preferences.push_notifications !== undefined && {push_notifications: values.preferences.push_notifications}),
					...(values.preferences.theme !== undefined && {theme: values.preferences.theme}),
				},
			}),
		}

		try {
			const updatedProfile = await updateCurrentUserProfile(updateData)
			toast.success('Profile details updated successfully.')
			onSubmitSuccess(updatedProfile)
		} catch (error: unknown) {
			console.error('Error updating profile details:', error)
			const message = error instanceof Error ? error.message : 'Failed to update profile details.'
			toast.error(`Update failed: ${message}`)
		} finally {
			setIsSubmitting(false)
		}
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
				<FormField
					control={form.control}
					name='bio'
					render={({field}) => (
						<FormItem>
							<FormLabel>Bio</FormLabel>
							<FormControl>
								<Textarea placeholder='Tell us a little bit about yourself' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='address'
					render={({field}) => (
						<FormItem>
							<FormLabel>Address</FormLabel>
							<FormControl>
								<Input placeholder='Enter your address' {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='date_of_birth'
					render={({field}) => (
						<FormItem className='flex flex-col'>
							<FormLabel>Date of Birth</FormLabel>
							<Popover>
								<PopoverTrigger asChild>
									<FormControl>
										<Button variant={'outline'} className={cn('w-[240px] pl-3 text-left font-normal', !field.value && 'text-muted-foreground')}>
											{field.value ? (
												field.value.toLocaleDateString('en-US', {
													year: 'numeric',
													month: 'long',
													day: 'numeric',
												})
											) : (
												<span>Pick a date</span>
											)}
											<CalendarIcon className='ml-auto h-4 w-4 opacity-50' />
										</Button>
									</FormControl>
								</PopoverTrigger>
								<PopoverContent className='w-auto p-0' align='start'>
									<Calendar mode='single' selected={field.value} onSelect={field.onChange} disabled={(date) => date > new Date() || date < new Date('1900-01-01')} initialFocus />
								</PopoverContent>
							</Popover>
							<FormMessage />
						</FormItem>
					)}
				/>
				{/* TODO: Add a better input for interests (e.g., TagInput) */}
				{/* <FormField ... name="interests" ... /> */}

				<h3 className='text-lg font-medium pt-4'>Preferences</h3>
				<FormField
					control={form.control}
					name='preferences.email_notifications'
					render={({field}) => (
						<FormItem className='flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4'>
							<FormControl>
								<Checkbox checked={field.value} onCheckedChange={field.onChange} />
							</FormControl>
							<div className='space-y-1 leading-none'>
								<FormLabel>Email Notifications</FormLabel>
								<FormDescription>Receive notifications via email.</FormDescription>
							</div>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='preferences.push_notifications'
					render={({field}) => (
						<FormItem className='flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4'>
							<FormControl>
								<Checkbox checked={field.value} onCheckedChange={field.onChange} />
							</FormControl>
							<div className='space-y-1 leading-none'>
								<FormLabel>Push Notifications</FormLabel>
								<FormDescription>Receive push notifications on your devices.</FormDescription>
							</div>
							<FormMessage />
						</FormItem>
					)}
				/>
				{/* TODO: Add Select component for Theme and Language */}
				{/* <FormField ... name="preferences.theme" ... /> */}
				{/* <FormField ... name="preferences.language" ... /> */}

				<Button type='submit' disabled={isSubmitting || !form.formState.isDirty}>
					{isSubmitting ? 'Saving...' : 'Save Profile Details'}
				</Button>
			</form>
		</Form>
	)
}

// --- Main Page Component ---

export default function ProfilePage() {
	const {user: authUser, loading: authLoading} = useAuth()
	const [userData, setUserData] = useState<UserOutput | null>(null)
	const [profileData, setProfileData] = useState<UserProfile | null>(null)
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)
	const [isAvatarHovered, setIsAvatarHovered] = useState(false)
	const fileInputRef = useRef<HTMLInputElement>(null)
	const [isUploadingAvatar, setIsUploadingAvatar] = useState(false)

	// Fetch data on mount
	useEffect(() => {
		const fetchData = async () => {
			if (!authLoading && authUser) {
				try {
					setLoading(true)
					setError(null)
					// Fetch only the current user, profile is included
					const fetchedUser = await getCurrentUser()
					setUserData(fetchedUser)
					// Extract profile from user data, handle potential null profile
					setProfileData(fetchedUser.profile || null)
				} catch (err: unknown) {
					console.error('Error fetching user data:', err) // Updated error message context
					let message = 'Failed to load user information.' // Updated error message context
					if (err instanceof Error) {
						message = err.message
					}
					setError(message)
					toast.error(message) // Show toast on fetch error
				} finally {
					setLoading(false)
				}
			} else if (!authLoading && !authUser) {
				setError('User not authenticated.')
				setLoading(false)
			}
		}
		fetchData()
	}, [authUser, authLoading])

	const isLoading = loading || authLoading

	// Handlers to update state after successful form submissions
	const handleGeneralInfoUpdate = (updatedUser: UserOutput) => {
		setUserData(updatedUser)
	}

	const handleProfileDetailsUpdate = (updatedProfile: UserProfile) => {
		setProfileData(updatedProfile)
	}

	const handleAvatarFileSelect = async (event: ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0]
		if (file && userData) {
			setIsUploadingAvatar(true)
			// formData is not needed here as updateUserAvatar expects a File object directly

			try {
				// Assuming updateUserAvatar returns the updated user or at least the new avatar URL
				const updatedUser = await updateUserAvatar(file) // Changed to pass file directly
				setUserData((prevData) => ({
					...prevData!,
					avatar: updatedUser.avatar, // Adjust based on actual response structure
				}))
				toast.success('Avatar updated successfully!')
			} catch (err: unknown) {
				console.error('Error uploading avatar:', err)
				const message = err instanceof Error ? err.message : 'Failed to upload avatar.'
				toast.error(`Upload failed: ${message}`)
			} finally {
				setIsUploadingAvatar(false)
				// Reset file input value to allow selecting the same file again if needed
				if (fileInputRef.current) {
					fileInputRef.current.value = ''
				}
			}
		}
	}

	return (
		<div className='container mx-auto p-4 space-y-6'>
			<h1 className='text-3xl font-bold'>User Profile</h1>
			<input type='file' ref={fileInputRef} onChange={handleAvatarFileSelect} accept='image/*' style={{display: 'none'}} />

			{isLoading ? (
				<Card>
					<CardHeader className='flex flex-row items-center space-x-4'>
						<Skeleton className='h-16 w-16 rounded-full' />
						<div className='space-y-2'>
							<Skeleton className='h-6 w-[250px]' />
							<Skeleton className='h-4 w-[200px]' />
						</div>
					</CardHeader>
					<CardContent className='space-y-4 pt-6'>
						<Skeleton className='h-4 w-full' />
						<Skeleton className='h-4 w-3/4' />
						<Skeleton className='h-4 w-1/2' />
					</CardContent>
				</Card>
			) : error ? (
				<Card className='border-destructive'>
					<CardHeader>
						<CardTitle className='text-destructive'>Error Loading Profile</CardTitle>
					</CardHeader>
					<CardContent>
						<p className='text-destructive-foreground'>{error}</p>
						<p className='text-sm text-muted-foreground'>Please try refreshing the page.</p>
					</CardContent>
				</Card>
			) : userData ? (
				<>
					{/* User Header Section */}
					<Card>
						<CardHeader className='flex flex-row items-center space-x-4'>
							<div className='relative h-16 w-16' onMouseEnter={() => setIsAvatarHovered(true)} onMouseLeave={() => setIsAvatarHovered(false)}>
								<Avatar className='h-16 w-16'>
									<AvatarImage src={userData.avatar || undefined} alt='User Avatar' />
									<AvatarFallback>{getInitials(userData.first_name, userData.last_name, userData.email)}</AvatarFallback>
								</Avatar>
								{isAvatarHovered && (
									<div className='absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 rounded-full cursor-pointer'>
										<Button variant='ghost' size='icon' className='text-white hover:text-gray-200' onClick={() => fileInputRef.current?.click()} disabled={isUploadingAvatar}>
											{isUploadingAvatar ? <div className='h-5 w-5 animate-spin rounded-full border-b-2 border-white'></div> : <UploadCloudIcon className='h-6 w-6' />}
										</Button>
									</div>
								)}
							</div>
							<div className='flex-grow'>
								<CardTitle className='text-2xl'>{`${userData.first_name || ''} ${userData.last_name || ''}`.trim()}</CardTitle>
								<CardDescription>{userData.email}</CardDescription>
								<div className='text-sm text-muted-foreground mt-1'>
									{!!userData.roles?.length ? <span>Role: {userData.roles?.join(', ')}</span> : <span>Role: User</span>}
									{' | '}
									<span>
										Status: <span className={`capitalize ${userData.status === 'active' ? 'text-green-600' : 'text-yellow-600'}`}>{userData.status}</span>
									</span>{' '}
									| <span>Joined: {new Date(userData.created_at).toLocaleDateString()}</span>
								</div>
							</div>
							{/* Optional: Add Avatar Upload Button Here */}
						</CardHeader>
					</Card>

					{/* Tabs for Editing Different Sections */}
					<Tabs defaultValue='general' className='w-full'>
						{/* Updated grid-cols to 4 */}
						<TabsList className='grid w-full grid-cols-4'>
							<TabsTrigger value='general'>General Info</TabsTrigger>
							<TabsTrigger value='profile'>Profile Details</TabsTrigger>
							<TabsTrigger value='password'>Password</TabsTrigger>
							{/* Added Security Tab Trigger */}
							<TabsTrigger value='security'>Security</TabsTrigger>
						</TabsList>

						{/* General Info Tab */}
						<TabsContent value='general'>
							<Card>
								<CardHeader>
									<CardTitle>General Information</CardTitle>
									<CardDescription>Update your first name, last name, and phone number.</CardDescription>
								</CardHeader>
								<CardContent>
									<GeneralInfoForm userData={userData} onSubmitSuccess={handleGeneralInfoUpdate} />
								</CardContent>
							</Card>
						</TabsContent>

						{/* Profile Details Tab */}
						<TabsContent value='profile'>
							<Card>
								<CardHeader>
									<CardTitle>Profile Details</CardTitle>
									<CardDescription>Update your bio, address, date of birth, and preferences.</CardDescription>
								</CardHeader>
								<CardContent>
									{/* Pass profileData and handler */}
									<ProfileDetailsForm profileData={profileData} onSubmitSuccess={handleProfileDetailsUpdate} />
								</CardContent>
							</Card>
						</TabsContent>

						{/* Password Tab */}
						<TabsContent value='password'>
							<Card>
								<CardHeader>
									<CardTitle>Change Password</CardTitle>
									<CardDescription>Update your account password. Make sure it&#39;s strong!</CardDescription>
								</CardHeader>
								<CardContent>
									{/* Pass userData to PasswordForm */}
									<PasswordForm userData={userData} />
								</CardContent>
							</Card>
						</TabsContent>

						{/* Security Tab */}
						<TabsContent value='security'>
							<Card>
								<CardHeader>
									<CardTitle>Security Settings</CardTitle>
									<CardDescription>Manage your email/phone verification and two-factor authentication.</CardDescription>
								</CardHeader>
								<CardContent className='space-y-6'>
									{/* Integrate Security Components */}
									<EmailVerificationStatus userData={userData} />
									<PhoneVerification userData={userData} onUpdate={handleGeneralInfoUpdate} />
									<TwoFactorAuthManagement userData={userData} onUpdate={handleGeneralInfoUpdate} />
								</CardContent>
							</Card>
						</TabsContent>
					</Tabs>
				</>
			) : (
				<Card>
					<CardHeader>
						<CardTitle>No User Data</CardTitle>
					</CardHeader>
					<CardContent>
						<p>Could not load user information.</p>
					</CardContent>
				</Card>
			)}
		</div>
	)
}
