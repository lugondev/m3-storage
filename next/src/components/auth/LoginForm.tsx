'use client'

import {zodResolver} from '@hookform/resolvers/zod'
import {useForm} from 'react-hook-form'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form' // Added FormDescription
import {Input} from '@/components/ui/input'
import {InputOTP, InputOTPGroup, InputOTPSlot} from '@/components/ui/input-otp' // Import InputOTP
import {useAuth} from '@/contexts/AuthContext'
import {useState} from 'react'
import {Verify2FARequest} from '@/lib/apiClient' // Import Verify2FARequest type

const formSchema = z.object({
	email: z.string().email({message: 'Invalid email address.'}),
	password: z.string().min(6, {message: 'Password must be at least 6 characters.'}),
})

// Schema for the 2FA code
const twoFactorSchema = z.object({
	code: z.string().min(6, {message: 'Your one-time code must be 6 characters.'}),
})

export function LoginForm() {
	// Get necessary functions and state from AuthContext
	// Added twoFactorSessionToken from context
	const {signInWithEmail, verifyTwoFactorCode, isTwoFactorPending, twoFactorSessionToken} = useAuth()
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)
	// Removed local state for twoFactorSessionToken as it's now read from context

	// Form for email/password
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: '',
			password: '',
		},
	})

	async function onSubmit(values: z.infer<typeof formSchema>) {
		setLoading(true)
		setError(null)
		console.log('Attempting email/password sign in with:', values)

		const payload: {email: string; password: string} = {
			email: values.email,
			password: values.password,
		}

		try {
			// Call signInWithEmail and check the result
			const result = await signInWithEmail(payload)

			if (result.success && result.twoFactorRequired && result.sessionToken) {
				// 2FA is required. AuthContext now stores the token.
				// We no longer need to set local state here.
				console.log('2FA required. Token stored in AuthContext.')
				// AuthContext sets isTwoFactorPending, UI will switch to 2FA form
			} else if (result.success && !result.twoFactorRequired) {
				// Login successful without 2FA, AuthContext handles state and redirect
				console.log('Login successful (no 2FA).')
			} else {
				// Handle potential errors returned in the result object
				if (result.error instanceof Error) {
					setError(result.error.message)
				} else if (typeof result.error === 'string') {
					setError(result.error)
				} else {
					setError('Login failed. Please check your credentials.')
				}
			}
		} catch (err) {
			// Catch errors thrown by signInWithEmail (e.g., network errors)
			// AuthContext might show a toast, but set local error state too
			if (err instanceof Error) {
				setError(err.message)
			} else {
				setError('An unexpected error occurred during login.')
			}
			console.error('Email/Password Sign in error caught in form:', err)
		} finally {
			setLoading(false)
		}
	}

	// Form for 2FA code
	const twoFactorForm = useForm<z.infer<typeof twoFactorSchema>>({
		resolver: zodResolver(twoFactorSchema),
		defaultValues: {
			code: '',
		},
	})

	// Handler for 2FA form submission
	async function onTwoFactorSubmit(values: z.infer<typeof twoFactorSchema>) {
		setLoading(true)
		setError(null)
		console.log('[LoginForm] onTwoFactorSubmit called with values:', values) // Log 1: Function entry

		// Read the token directly from context here
		if (!twoFactorSessionToken) {
			console.error('[LoginForm] Error: twoFactorSessionToken from context is missing!') // Log 2: Token check failure
			setError('2FA session is invalid or expired. Please try logging in again.')
			setLoading(false)
			// Optionally force back to email/password step by clearing context state?
			// Or rely on user restarting the login flow.
			return
		}

		console.log('[LoginForm] Attempting 2FA verification with token:', twoFactorSessionToken, 'and code:', values.code) // Log 3: Before calling context function
		try {
			const verifyData: Verify2FARequest = {
				two_factor_session_token: twoFactorSessionToken,
				code: values.code,
			}
			// Call the context function
			const result = await verifyTwoFactorCode(verifyData)
			console.log('[LoginForm] verifyTwoFactorCode result:', result) // Log 4: After calling context function

			if (result.success) {
				console.log('[LoginForm] 2FA verification successful.')
			} else {
				console.error('[LoginForm] 2FA verification failed:', result.error) // Log 5: Verification failure
				// Handle verification failure
				if (result.error instanceof Error) {
					setError(result.error.message)
				} else if (typeof result.error === 'string') {
					setError(result.error)
				} else {
					setError('Invalid 2FA code or session expired.')
				}
				twoFactorForm.reset() // Clear the input field on error
			}
		} catch (err) {
			console.error('[LoginForm] Error caught during verifyTwoFactorCode call:', err) // Log 6: Catch block
			// Catch errors thrown by verifyTwoFactorCode
			if (err instanceof Error) {
				setError(err.message)
			} else {
				setError('An unexpected error occurred during 2FA verification.')
			}
			twoFactorForm.reset() // Clear the input field on error
		} finally {
			console.log('[LoginForm] onTwoFactorSubmit finally block.') // Log 7: Finally block
			setLoading(false)
		}
	}

	// Conditional Rendering based on isTwoFactorPending
	if (isTwoFactorPending) {
		return (
			<Form {...twoFactorForm}>
				<form onSubmit={twoFactorForm.handleSubmit(onTwoFactorSubmit)} className='space-y-6'>
					<FormField
						control={twoFactorForm.control}
						name='code'
						render={({field}) => (
							<FormItem>
								<FormLabel>One-Time Password</FormLabel>
								<FormControl>
									<InputOTP maxLength={6} {...field}>
										<InputOTPGroup>
											<InputOTPSlot index={0} />
											<InputOTPSlot index={1} />
											<InputOTPSlot index={2} />
											<InputOTPSlot index={3} />
											<InputOTPSlot index={4} />
											<InputOTPSlot index={5} />
										</InputOTPGroup>
									</InputOTP>
								</FormControl>
								<FormDescription>Please enter the 6-digit code from your authenticator app.</FormDescription>
								<FormMessage />
							</FormItem>
						)}
					/>
					{error && <p className='text-sm text-red-500'>{error}</p>}
					<Button type='submit' className='w-full' disabled={loading}>
						{loading ? 'Verifying...' : 'Verify Code'}
					</Button>
				</form>
			</Form>
		)
	}

	// Render email/password form if 2FA is not pending
	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
				<FormField
					control={form.control}
					name='email'
					render={({field}) => (
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input placeholder='your@email.com' {...field} type='email' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='password'
					render={({field}) => (
						<FormItem>
							<FormLabel>Password</FormLabel>
							<FormControl>
								<Input placeholder='********' {...field} type='password' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				{error && <p className='text-sm text-red-500'>{error}</p>}
				<Button type='submit' className='w-full' disabled={loading}>
					{loading ? 'Signing in...' : 'Sign in with Email'}
				</Button>
			</form>
		</Form>
	)
}
