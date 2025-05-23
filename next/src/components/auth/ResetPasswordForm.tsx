'use client'

import {zodResolver} from '@hookform/resolvers/zod'
import {useForm} from 'react-hook-form'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form'
import {Input} from '@/components/ui/input'
import {resetPassword} from '@/services/authService' // Import directly from service
import {useState} from 'react'
import {toast} from 'sonner'
import {useRouter} from 'next/navigation'

// Schema for password reset
const formSchema = z
	.object({
		new_password: z.string().min(8, {message: 'Password must be at least 8 characters.'}),
		confirmPassword: z.string(),
	})
	.refine((data) => data.new_password === data.confirmPassword, {
		message: "Passwords don't match",
		path: ['confirmPassword'], // path of error
	})

interface ResetPasswordFormProps {
	token: string | null // Token received from URL param
}

export function ResetPasswordForm({token}: ResetPasswordFormProps) {
	const router = useRouter()
	const [loading, setLoading] = useState(false)
	const [message, setMessage] = useState<string | null>(null) // For success/error messages

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			new_password: '',
			confirmPassword: '',
		},
	})

	// Show error if token is missing early
	if (!token) {
		return <p className='text-red-500 text-center'>Invalid or missing reset token.</p>
	}

	async function onSubmit(values: z.infer<typeof formSchema>) {
		if (!token) return // Should not happen due to early return, but good practice

		setLoading(true)
		setMessage(null)
		try {
			await resetPassword({token: token, new_password: values.new_password})
			setMessage('Your password has been reset successfully.')
			toast.success('Password reset successful!')
			form.reset()
			// Redirect to login page after a short delay
			setTimeout(() => {
				router.push('/login')
			}, 2000)
		} catch (err) {
			console.error('Reset password error:', err)
			// Use const as errorMessage is not reassigned here
			const errorMessage = 'Failed to reset password. The link may be invalid or expired.'
			if (err instanceof Error) {
				// Avoid generic errors
			}
			setMessage(errorMessage)
			toast.error(errorMessage)
		} finally {
			setLoading(false)
		}
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
				<FormField
					control={form.control}
					name='new_password'
					render={({field}) => (
						<FormItem>
							<FormLabel>New Password</FormLabel>
							<FormControl>
								<Input placeholder='********' {...field} type='password' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name='confirmPassword'
					render={({field}) => (
						<FormItem>
							<FormLabel>Confirm New Password</FormLabel>
							<FormControl>
								<Input placeholder='********' {...field} type='password' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				{message && <p className={`text-sm ${message.includes('Failed') || message.includes('invalid') ? 'text-red-500' : 'text-green-600'}`}>{message}</p>}
				<Button type='submit' className='w-full' disabled={loading}>
					{loading ? 'Resetting...' : 'Reset Password'}
				</Button>
			</form>
		</Form>
	)
}
