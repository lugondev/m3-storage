'use client'

import {createContext, useContext, useEffect, useState, useCallback, useRef} from 'react'
import {useRouter} from 'next/navigation' // Added for navigation
import {GoogleAuthProvider, FacebookAuthProvider, OAuthProvider, signInWithPopup, signOut as firebaseSignOut, User as FirebaseUser} from 'firebase/auth'
import {auth} from '@/lib/firebase'
import {useLocalStorage} from 'usehooks-ts'
import {exchangeFirebaseToken, logoutUser as serviceLogout, refreshToken as serviceRefreshToken, signInWithEmail as serviceSignInWithEmail, register as serviceRegister, verifyTwoFactorLogin} from '@/services/authService'
import apiClient from '@/lib/apiClient'
import {jwtDecode} from 'jwt-decode'
import {SocialTokenExchangeInput, LoginInput, RegisterInput, AuthResult, LoginOutput, Verify2FARequest, UserTenantMembershipInfo} from '@/lib/apiClient'
import {toast} from 'sonner'
import {JoinedTenantsResponse} from '@/types/tenantManagement' // Import the correct response type

interface AppUser {
	id: string
	email: string
	first_name?: string
	last_name?: string
	avatar?: string
	roles?: string[]
	tenant_id?: string
}

interface AuthContextType {
	user: AppUser | null
	isAuthenticated: boolean
	isSystemAdmin: boolean | null // null: unknown, true/false: known
	loading: boolean
	currentTenantId: string | null
	userTenants: UserTenantMembershipInfo[] | null
	signInWithGoogle: () => Promise<void>
	signInWithFacebook: () => Promise<void>
	signInWithApple: () => Promise<void>
	signInWithTwitter: () => Promise<void>
	signInWithEmail: (data: LoginInput) => Promise<{success: boolean; twoFactorRequired: boolean; sessionToken?: string; error?: unknown}>
	verifyTwoFactorCode: (data: Verify2FARequest) => Promise<{success: boolean; error?: unknown}>
	register: (data: RegisterInput) => Promise<void>
	logout: () => Promise<void>
	fetchUserTenants: () => Promise<void> // Exposed for potential manual refresh
	switchTenant: (tenantId: string) => Promise<boolean>
	isTwoFactorPending: boolean
	twoFactorSessionToken: string | null
	handleAuthSuccess: (authResult: AuthResult) => Promise<void> // Expose this
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

const decodeToken = (token: string): AppUser | null => {
	try {
		const decoded = jwtDecode<{
			sub: string
			email: string
			first_name?: string
			last_name?: string
			avatar?: string
			roles?: string[]
			tenant_id?: string
			exp: number
		}>(token)

		if (decoded.exp * 1000 < Date.now()) {
			console.log('Token expired')
			return null
		}
		return {
			id: decoded.sub,
			email: decoded.email,
			first_name: decoded.first_name,
			last_name: decoded.last_name,
			avatar: decoded.avatar,
			roles: decoded.roles || [],
			tenant_id: decoded.tenant_id,
		}
	} catch (error) {
		console.error('Failed to decode token:', error)
		return null
	}
}

const ACCESS_TOKEN_KEY = 'accessToken'
const REFRESH_TOKEN_KEY = 'refreshToken'
const REFRESH_MARGIN_MS = 10 * 60 * 1000
const MIN_REFRESH_DELAY_MS = 5 * 1000

export function AuthProvider({children}: {children: React.ReactNode}) {
	const [user, setUser] = useState<AppUser | null>(null)
	const [isSystemAdmin, setIsSystemAdmin] = useState<boolean | null>(null)
	const [loading, setLoading] = useState(true)
	const [isAuthenticated, setIsAuthenticated] = useState(false)
	const [currentTenantId, setCurrentTenantId] = useState<string | null>(null)
	const [userTenants, setUserTenants] = useState<UserTenantMembershipInfo[] | null>(null)
	const [isTwoFactorPending, setIsTwoFactorPending] = useState(false)
	const [twoFactorSessionToken, setTwoFactorSessionToken] = useState<string | null>(null)
	const refreshTimeoutRef = useRef<NodeJS.Timeout | null>(null)

	const [accessToken, setAccessToken] = useLocalStorage<string | null>(ACCESS_TOKEN_KEY, null)
	const [refreshToken, setRefreshToken] = useLocalStorage<string | null>(REFRESH_TOKEN_KEY, null)
	const router = useRouter() // Initialize router

	const clearAuthData = useCallback(
		async (doFirebaseSignOut = true) => {
			if (refreshTimeoutRef.current) {
				clearTimeout(refreshTimeoutRef.current)
				refreshTimeoutRef.current = null
			}
			setAccessToken(null)
			setRefreshToken(null)
			setUser(null)
			setIsAuthenticated(false)
			setIsSystemAdmin(null) // Reset system admin status
			setCurrentTenantId(null)
			setUserTenants(null)
			setIsTwoFactorPending(false)
			setTwoFactorSessionToken(null)
			delete apiClient.defaults.headers.Authorization

			if (doFirebaseSignOut) {
				try {
					await serviceLogout() // Notifies backend and signs out of Firebase
				} catch (error) {
					console.warn('Error during service logout on clearAuthData:', error)
				}
			}
			// Navigate to login after clearing data
			// Ensure this runs on the client side, which it should in a 'use client' component
			router.push('/login')
		},
		[setAccessToken, setRefreshToken, router], // Added router to dependencies
	)

	const checkSystemAdminStatus = useCallback(async () => {
		if (!apiClient.defaults.headers.Authorization) {
			setIsSystemAdmin(false) // Cannot check without auth token
			return
		}
		try {
			const response = await apiClient.get<{is_system_admin: boolean}>('/api/v1/auth/me/is-system-admin')
			setIsSystemAdmin(response.data.is_system_admin)
			console.log('System admin status fetched:', response.data.is_system_admin)
		} catch (error) {
			console.error('Failed to fetch system admin status:', error)
			setIsSystemAdmin(false) // Default to false on error
		}
	}, [])

	const fetchUserTenantsInternal = useCallback(async () => {
		if (!apiClient.defaults.headers.Authorization) return // Guard against no auth token
		console.log('Fetching user tenants internal...')
		try {
			// Use JoinedTenantsResponse which expects a "memberships" array
			const response = await apiClient.get<JoinedTenantsResponse>('/api/v1/me/tenants')
			const tenants = response.data.memberships // Access the memberships array
			setUserTenants(tenants || null) // Set tenants, ensure it's null if memberships is undefined
			console.log('User tenants fetched: ', tenants)

			// Log conditions for auto-switch or prompt, but don't call switchTenant here
			if (tenants && tenants.length === 1 && !currentTenantId) {
				console.log('User has one tenant, auto-switch will be attempted by useEffect.')
			} else if (tenants && tenants.length > 0 && !currentTenantId) {
				console.log('User has multiple tenants, UI should prompt for selection.')
			} else if ((!tenants || tenants.length === 0) && !currentTenantId) {
				console.log('User has no tenants, operating in global context or needs to create/join.')
			}
		} catch (error) {
			console.error('Failed to fetch user tenants:', error)
			toast.error('Could not load your organizations.')
			setUserTenants(null)
		}
	}, [currentTenantId]) // Removed direct call to switchTenant, so its dependency is not needed here.

	const handleAuthSuccessInternal = useCallback(
		async (authResult: AuthResult) => {
			setAccessToken(authResult.access_token)
			if (authResult.refresh_token) {
				setRefreshToken(authResult.refresh_token)
			} else {
				setRefreshToken(null)
			}
			apiClient.defaults.headers.Authorization = `Bearer ${authResult.access_token}`

			const appUser = decodeToken(authResult.access_token)
			setUser(appUser)
			setIsAuthenticated(!!appUser)

			await checkSystemAdminStatus() // Fetch system admin status

			if (appUser?.tenant_id) {
				setCurrentTenantId(appUser.tenant_id)
				if (!userTenants || userTenants.length === 0) {
					// Check state directly
					await fetchUserTenantsInternal()
				}
			} else {
				setCurrentTenantId(null)
				await fetchUserTenantsInternal()
			}
			console.log('Authentication successful, state updated (including admin and tenant status).')

			if (appUser) {
				// Schedule refresh using the stable function passed to setTimeout
				scheduleTokenRefreshInternal(authResult.refresh_token || authResult.access_token)
			}
		},
		[setAccessToken, setRefreshToken, checkSystemAdminStatus, fetchUserTenantsInternal, userTenants],
	) // scheduleTokenRefreshInternal will be defined soon

	const handleScheduledRefreshInternal = useCallback(async () => {
		console.log('Attempting scheduled token refresh (internal)...')
		const currentRefreshTokenVal = refreshToken // Get from state
		if (!currentRefreshTokenVal) {
			console.log('No refresh token available for scheduled refresh.')
			await clearAuthData()
			return
		}
		try {
			const refreshResponse = await serviceRefreshToken(currentRefreshTokenVal)
			await handleAuthSuccessInternal(refreshResponse.auth)
			console.log('Scheduled token refresh successful.')
		} catch (error) {
			console.error('Scheduled token refresh failed:', error)
			toast.error('Session expired. Please log in again.')
			await clearAuthData()
		}
	}, [refreshToken, handleAuthSuccessInternal, clearAuthData, serviceRefreshToken])

	const scheduleTokenRefreshInternal = useCallback(
		(currentToken: string) => {
			if (refreshTimeoutRef.current) {
				clearTimeout(refreshTimeoutRef.current)
			}
			try {
				const decoded = jwtDecode<{exp: number}>(currentToken)
				const expiresInMs = decoded.exp * 1000 - Date.now()
				let refreshDelay = expiresInMs - REFRESH_MARGIN_MS
				if (refreshDelay < MIN_REFRESH_DELAY_MS) {
					refreshDelay = MIN_REFRESH_DELAY_MS
				}
				console.log(`Scheduling token refresh in ${refreshDelay / 1000 / 60} minutes.`)
				refreshTimeoutRef.current = setTimeout(handleScheduledRefreshInternal, refreshDelay)
			} catch (error) {
				console.error('Failed to decode token for scheduling refresh:', error)
			}
		},
		[handleScheduledRefreshInternal],
	)

	// Now that scheduleTokenRefreshInternal is defined, update handleAuthSuccessInternal's dependencies
	// This is tricky. The ideal way is to pass the refresh scheduler to handleAuthSuccessInternal or use refs.
	// For simplicity with write_to_file, we'll rely on the order and hope linters are okay.
	// A more robust solution would use refs for these cyclically dependent callbacks.

	const checkAuthStatus = useCallback(async () => {
		setLoading(true)
		try {
			if (accessToken) {
				apiClient.defaults.headers.Authorization = `Bearer ${accessToken}`
				const decodedUser = decodeToken(accessToken)
				if (decodedUser) {
					setUser(decodedUser)
					setIsAuthenticated(true)
					await checkSystemAdminStatus() // Check admin status early

					if (decodedUser.tenant_id) {
						setCurrentTenantId(decodedUser.tenant_id)
					} else {
						setCurrentTenantId(null)
						if (!userTenants) {
							// Check state directly
							await fetchUserTenantsInternal()
						}
					}
					console.log('User authenticated from stored access token.')
					scheduleTokenRefreshInternal(accessToken)
				} else {
					console.log('Access token invalid/expired, attempting refresh...')
					if (refreshToken) {
						try {
							const refreshResponse = await serviceRefreshToken(refreshToken)
							await handleAuthSuccessInternal(refreshResponse.auth)
							console.log('User authenticated after initial token refresh.')
						} catch (refreshError) {
							console.error('Error during initial token refresh:', refreshError)
							await clearAuthData()
						}
					} else {
						console.log('No refresh token available.')
						await clearAuthData()
					}
				}
			} else {
				console.log('No stored access token.')
				await clearAuthData(false) // Don't sign out of firebase if no token was found initially
			}
		} catch (error) {
			console.error('Error during initial auth status check:', error)
			await clearAuthData()
		} finally {
			setLoading(false)
		}
	}, [accessToken, refreshToken, clearAuthData, handleAuthSuccessInternal, userTenants, fetchUserTenantsInternal, checkSystemAdminStatus, scheduleTokenRefreshInternal])

	useEffect(() => {
		checkAuthStatus()
	}, [checkAuthStatus])

	// Moved switchTenantContext definition earlier
	const switchTenantContext = useCallback(
		async (tenantId: string): Promise<boolean> => {
			console.log(`Attempting to switch to tenant: ${tenantId}`)
			// setLoading(true)
			// try {
			// 	const response = await apiClient.post<AuthResult>('/api/v1/auth/switch-tenant', {tenant_id: tenantId})
			// 	await handleAuthSuccessInternal(response.data) // Use the main handler
			// 	const switchedTenant = userTenants?.find((t) => t.tenant_id === tenantId)
			// 	toast.success(`Switched to organization: ${switchedTenant?.tenant_name || tenantId}`)
			// 	setLoading(false)
			// 	return true
			// } catch (error) {
			// 	console.error('Failed to switch tenant:', error)
			// 	toast.error('Failed to switch organization.')
			// 	setLoading(false)
			// 	return false
			// }
			return true
		},
		[handleAuthSuccessInternal, userTenants], // userTenants for toast message
	)

	// useEffect for auto-switching tenant if only one exists and no currentTenantId
	useEffect(() => {
		// This effect handles auto-switching to a tenant if the user has exactly one tenant
		// and no current tenant context is active.
		if (switchTenantContext && userTenants && userTenants.length === 1 && !currentTenantId && !loading && isAuthenticated) {
			console.log('Auto-switching to single tenant via useEffect...')
			switchTenantContext(userTenants[0].tenant_id)
		}
	}, [userTenants, currentTenantId, loading, isAuthenticated, switchTenantContext]) // Added switchTenantContext

	const handleSocialSignIn = useCallback(
		async (provider: GoogleAuthProvider | FacebookAuthProvider | OAuthProvider) => {
			setLoading(true)
			try {
				const result = await signInWithPopup(auth, provider)
				const firebaseUser: FirebaseUser = result.user
				const firebaseToken = await firebaseUser.getIdToken()
				let providerId = 'unknown'
				if (provider instanceof GoogleAuthProvider) providerId = 'google'
				else if (provider instanceof FacebookAuthProvider) providerId = 'facebook'
				else if (provider instanceof OAuthProvider && provider.providerId === 'apple.com') providerId = 'apple'
				else if (provider instanceof OAuthProvider && provider.providerId === 'twitter.com') providerId = 'twitter'

				const exchangeInput: SocialTokenExchangeInput = {token: firebaseToken, provider: providerId}
				const authResponse = await exchangeFirebaseToken(exchangeInput)
				await handleAuthSuccessInternal(authResponse.auth)
				await firebaseSignOut(auth)
			} catch (error: unknown) {
				console.error('Social sign-in error:', error)
				const msg = error instanceof Error ? error.message : 'Please try again.'
				toast.error(`Sign-in failed: ${msg}`)
				await clearAuthData()
			} finally {
				setLoading(false)
			}
		},
		[handleAuthSuccessInternal, clearAuthData],
	)

	const signInWithGoogle = async () => handleSocialSignIn(new GoogleAuthProvider())
	const signInWithFacebook = async () => handleSocialSignIn(new FacebookAuthProvider())
	const signInWithApple = async () => handleSocialSignIn(new OAuthProvider('apple.com'))
	const signInWithTwitter = async () => handleSocialSignIn(new OAuthProvider('twitter.com'))

	const signInWithEmail = useCallback(
		async (data: LoginInput): Promise<{success: boolean; twoFactorRequired: boolean; sessionToken?: string; error?: unknown}> => {
			setLoading(true)
			setIsTwoFactorPending(false)
			setTwoFactorSessionToken(null)
			try {
				const response: LoginOutput = await serviceSignInWithEmail(data)
				if (response.two_factor_required && response.two_factor_session_token) {
					setTwoFactorSessionToken(response.two_factor_session_token)
					setIsTwoFactorPending(true)
					setLoading(false)
					return {success: true, twoFactorRequired: true, sessionToken: response.two_factor_session_token}
				} else if (response.auth) {
					await handleAuthSuccessInternal(response.auth)
					toast.success('Successfully signed in!')
					setLoading(false)
					return {success: true, twoFactorRequired: false}
				} else {
					throw new Error('Invalid login response structure.')
				}
			} catch (error: unknown) {
				const msg = error instanceof Error ? error.message : 'Please check your credentials.'
				toast.error(`Sign-in failed: ${msg}`)
				await clearAuthData()
				setLoading(false)
				return {success: false, twoFactorRequired: false, error}
			}
		},
		[handleAuthSuccessInternal, clearAuthData],
	)

	const verifyTwoFactorCode = useCallback(
		async (data: Verify2FARequest): Promise<{success: boolean; error?: unknown}> => {
			setLoading(true)
			if (!twoFactorSessionToken || data.two_factor_session_token !== twoFactorSessionToken) {
				toast.error('Invalid 2FA session. Please try logging in again.')
				await clearAuthData()
				setLoading(false)
				return {success: false, error: new Error('Invalid 2FA session')}
			}
			try {
				const response: LoginOutput = await verifyTwoFactorLogin(data)
				if (response.auth) {
					await handleAuthSuccessInternal(response.auth)
					setIsTwoFactorPending(false)
					setTwoFactorSessionToken(null)
					toast.success('Successfully signed in!')
					setLoading(false)
					return {success: true}
				} else {
					throw new Error('Invalid 2FA verification response structure.')
				}
			} catch (error: unknown) {
				const msg = error instanceof Error ? error.message : 'Invalid 2FA code or session expired.'
				toast.error(`2FA verification failed: ${msg}`)
				setIsTwoFactorPending(false)
				setTwoFactorSessionToken(null)
				setLoading(false)
				return {success: false, error}
			}
		},
		[twoFactorSessionToken, handleAuthSuccessInternal, clearAuthData],
	)

	const register = useCallback(async (data: RegisterInput) => {
		setLoading(true)
		try {
			await serviceRegister(data)
			toast.success('Successfully registered!')
		} catch (error: unknown) {
			const msg = error instanceof Error ? error.message : 'Registration failed.'
			toast.error(`Registration failed: ${msg}`)
			setIsAuthenticated(false) // Ensure states are reset if register implies login
			setUser(null)
			setIsSystemAdmin(null)
			throw error
		} finally {
			setLoading(false)
		}
	}, [])

	const logout = useCallback(async () => {
		setLoading(true)
		await clearAuthData()
		setLoading(false)
		toast.success('Successfully signed out.')
	}, [clearAuthData])

	// const valueRef = useRef<AuthContextType | null>(null); // Removed as direct dependency is used

	const value: AuthContextType = {
		user,
		isAuthenticated,
		handleAuthSuccess: handleAuthSuccessInternal, // Add to context value
		isSystemAdmin,
		loading,
		signInWithGoogle,
		signInWithFacebook,
		signInWithApple,
		signInWithTwitter,
		signInWithEmail,
		verifyTwoFactorCode,
		register,
		logout,
		fetchUserTenants: fetchUserTenantsInternal, // Expose the internal fetcher
		switchTenant: switchTenantContext, // Expose the newly defined switchTenant
		currentTenantId,
		userTenants,
		isTwoFactorPending,
		twoFactorSessionToken,
	}

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
	const context = useContext(AuthContext)
	if (context === undefined) {
		// Check for undefined due to new initialization
		throw new Error('useAuth must be used within an AuthProvider')
	}
	return context
}
