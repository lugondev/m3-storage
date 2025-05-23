'use client'

import {zodResolver} from '@hookform/resolvers/zod'
import {useForm} from 'react-hook-form'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form'
import {Input} from '@/components/ui/input'
import {useAuth} from '@/contexts/AuthContext'
import {useState} from 'react'
import {useRouter} from 'next/navigation'
// Removed unused toast import

// Schema includes optional first_name and last_name
const formSchema = z
	.object({
		email: z.string().email({message: 'Invalid email address.'}),
		password: z.string().min(8, {message: 'Password must be at least 8 characters.'}),
		confirmPassword: z.string(),
		first_name: z.string().optional(),
		last_name: z.string().optional(),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match",
		path: ['confirmPassword'], // path of error
	})

export function RegisterForm() {
	const {register} = useAuth() // Get the register function from context
	const router = useRouter() // Get router instance
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: '',
			password: '',
			confirmPassword: '',
			first_name: '',
			last_name: '',
		},
	})

	async function onSubmit(values: z.infer<typeof formSchema>) {
		setLoading(true)
		setError(null)
		console.log('Attempting registration with:', values)
		try {
			// Destructure confirmPassword (it's validated but not sent), pass the rest
			// eslint-disable-next-line @typescript-eslint/no-unused-vars
			const {confirmPassword, ...registerData} = values
			await register(registerData)
			// AuthContext handles success toast, now navigate explicitly
			router.push('/login') // Navigate to login page
		} catch (err) {
			// AuthContext throws the error, catch it here
			// AuthContext already shows an error toast
			if (err instanceof Error) {
				setError(err.message) // Display error below the form
			} else {
				setError('An unexpected error occurred during registration.')
			}
			console.error('Registration error caught in form:', err)
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
					name='first_name'
					render={({field}) => (
						<FormItem>
							<FormLabel>First Name (Optional)</FormLabel>
							<FormControl>
								<Input placeholder='John' {...field} />
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
							<FormLabel>Last Name (Optional)</FormLabel>
							<FormControl>
								<Input placeholder='Doe' {...field} />
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
				<FormField
					control={form.control}
					name='confirmPassword'
					render={({field}) => (
						<FormItem>
							<FormLabel>Confirm Password</FormLabel>
							<FormControl>
								<Input placeholder='********' {...field} type='password' />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				{error && <p className='text-sm text-red-500'>{error}</p>}
				<Button type='submit' className='w-full' disabled={loading}>
					{loading ? 'Registering...' : 'Register'}
				</Button>
			</form>
		</Form>
	)
}
