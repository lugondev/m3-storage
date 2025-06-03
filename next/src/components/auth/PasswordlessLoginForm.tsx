'use client'

import {useState} from 'react'
import {useForm} from 'react-hook-form'
import {zodResolver} from '@hookform/resolvers/zod'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {Label} from '@/components/ui/label'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card'
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert'
import {Loader2} from 'lucide-react'
import {requestLoginLink} from '@/services/authService'
import axios, {AxiosError} from 'axios' // Import axios and AxiosError

const passwordlessLoginSchema = z.object({
	email: z.string().email({message: 'Invalid email address.'}),
})

type PasswordlessLoginValues = z.infer<typeof passwordlessLoginSchema>

interface PasswordlessLoginFormProps {
	onLinkSent?: (email: string) => void // Callback after link is sent
}

export function PasswordlessLoginForm({onLinkSent}: PasswordlessLoginFormProps) {
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)
	const [successMessage, setSuccessMessage] = useState<string | null>(null)

	const form = useForm<PasswordlessLoginValues>({
		resolver: zodResolver(passwordlessLoginSchema),
		defaultValues: {
			email: '',
		},
	})

	const onSubmit = async (values: PasswordlessLoginValues) => {
		setIsLoading(true)
		setError(null)
		setSuccessMessage(null)

		try {
			await requestLoginLink({
				email: values.email,
			})
			setSuccessMessage(`A login link has been sent to ${values.email}. Please check your inbox.`)
			if (onLinkSent) {
				onLinkSent(values.email)
			}
			form.reset()
		} catch (err: unknown) {
			// Changed from any to unknown
			console.error('Request login link failed:', err)
			let errorMessage = 'Failed to send login link. Please try again.'
			if (axios.isAxiosError(err)) {
				const axiosError = err as AxiosError<{message?: string}>
				if (axiosError.response?.data?.message) {
					errorMessage = axiosError.response.data.message
				}
			} else if (err instanceof Error) {
				errorMessage = err.message
			}
			setError(errorMessage)
		} finally {
			setIsLoading(false)
		}
	}

	return (
		<Card>
			<CardHeader>
				<CardTitle>Passwordless Login</CardTitle>
				<CardDescription>Enter your email to receive a magic login link. You can also specify an organization slug if applicable.</CardDescription>
			</CardHeader>
			<CardContent>
				<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
					<div className='space-y-2'>
						<Label htmlFor='email'>Email</Label>
						<Input id='email' type='email' placeholder='name@example.com' {...form.register('email')} disabled={isLoading} />
						{form.formState.errors.email && <p className='text-sm text-red-600'>{form.formState.errors.email.message}</p>}
					</div>
					{error && (
						<Alert variant='destructive'>
							<AlertTitle>Error</AlertTitle>
							<AlertDescription>{error}</AlertDescription>
						</Alert>
					)}
					{successMessage && (
						<Alert variant='default'>
							{' '}
							{/* Changed to 'default' for success */}
							<AlertTitle>Success</AlertTitle>
							<AlertDescription>{successMessage}</AlertDescription>
						</Alert>
					)}
					<Button type='submit' className='w-full' disabled={isLoading}>
						{isLoading ? <Loader2 className='mr-2 h-4 w-4 animate-spin' /> : 'Send Login Link'}
					</Button>
				</form>
			</CardContent>
		</Card>
	)
}
