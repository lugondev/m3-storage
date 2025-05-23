'use client'

import {zodResolver} from '@hookform/resolvers/zod'
import {useForm} from 'react-hook-form'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form'
import {Input} from '@/components/ui/input'
import {forgotPassword} from '@/services/authService' // Import directly from service
import {useState} from 'react'
import {toast} from 'sonner'

const formSchema = z.object({
	email: z.string().email({message: 'Invalid email address.'}),
})

export function ForgotPasswordForm() {
	const [loading, setLoading] = useState(false)
	const [message, setMessage] = useState<string | null>(null) // For success/error messages

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: '',
		},
	})

	async function onSubmit(values: z.infer<typeof formSchema>) {
		setLoading(true)
		setMessage(null) // Clear previous messages
		try {
			await forgotPassword({email: values.email})
			setMessage('If an account exists for this email, a password reset link has been sent.')
			toast.success('Password reset instructions sent!')
			form.reset() // Clear the form on success
		} catch (err) {
			console.error('Forgot password error:', err)
			// Use const as errorMessage is not reassigned here
			const errorMessage = 'Failed to send reset instructions. Please try again.'
			if (err instanceof Error) {
				// Avoid showing generic backend errors directly if possible
				// errorMessage = err.message;
			}
			setMessage(errorMessage) // Show error message
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
					name='email'
					render={({field}) => (
						<FormItem>
							<FormLabel>Email Address</FormLabel>
							<FormControl>
								<Input placeholder='your@email.com' {...field} type='email' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				{message && <p className={`text-sm ${message.includes('Failed') ? 'text-red-500' : 'text-green-600'}`}>{message}</p>}
				<Button type='submit' className='w-full' disabled={loading}>
					{loading ? 'Sending...' : 'Send Reset Instructions'}
				</Button>
			</form>
		</Form>
	)
}
