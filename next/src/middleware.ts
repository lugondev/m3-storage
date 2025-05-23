import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export async function middleware(request: NextRequest) {
	return NextResponse.next();
}

export const config = {
	matcher: [
		/*
		 * Match all request paths except for the ones starting with:
		 * - api (API routes)
		 * - _next/static (static files)
		 * - _next/image (image optimization files)
		 * - favicon.ico (favicon file)
		 * - login, register, forgot-password, etc. (public auth pages)
		 */
		'/((?!api|_next/static|_next/image|favicon.ico|login|register|forgot-password|reset-password|verify-email).*)',
	],
};
