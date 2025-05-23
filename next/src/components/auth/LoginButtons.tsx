'use client'

import {useAuth} from '@/contexts/AuthContext'
import {Button} from '@/components/ui/button'
// Import icons from lucide-react
import {Chrome, Facebook, Apple, Twitter} from 'lucide-react' // Added Twitter

export function LoginButtons() {
	const {signInWithGoogle, signInWithFacebook, signInWithApple, signInWithTwitter, loading} = useAuth() // Added signInWithTwitter

	// TODO: Add error handling display (e.g., using sonner or similar)

	return (
		<div className='flex flex-col space-y-2'>
			<Button onClick={signInWithGoogle} disabled={loading} variant='outline' className='w-full'>
				<Chrome className='mr-2 h-4 w-4' /> {/* Google Icon */}
				{loading ? 'Signing in...' : 'Sign in with Google'}
			</Button>
			<Button onClick={signInWithFacebook} disabled={loading} variant='outline' className='w-full bg-blue-600 text-white hover:bg-blue-700'>
				<Facebook className='mr-2 h-4 w-4' /> {/* Facebook Icon */}
				{loading ? 'Signing in...' : 'Sign in with Facebook'}
			</Button>
			<Button onClick={signInWithApple} disabled={loading} variant='outline' className='w-full bg-black text-white hover:bg-gray-800'>
				<Apple className='mr-2 h-4 w-4' /> {/* Apple Icon */}
				{loading ? 'Signing in...' : 'Sign in with Apple'}
			</Button>
			<Button onClick={signInWithTwitter} disabled={loading} variant='outline' className='w-full bg-sky-500 text-white hover:bg-sky-600'>
				<Twitter className='mr-2 h-4 w-4' /> {/* Twitter Icon */}
				{loading ? 'Signing in...' : 'Sign in with Twitter'}
			</Button>
			{/* Add more providers as needed */}
		</div>
	)
}
