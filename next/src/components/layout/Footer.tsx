export function Footer() {
	const currentYear = new Date().getFullYear()

	return (
		<footer className='border-t bg-background p-4 text-center text-sm text-muted-foreground'>
			<p>&copy; {currentYear} Your App Name. All rights reserved.</p>
			{/* Add any other footer content like links here */}
		</footer>
	)
}
