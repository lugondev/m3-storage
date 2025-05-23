import apiClient, {
	AuthResponse,
	LoginOutput, // <-- Add LoginOutput
	// Types specific to auth operations:
	LoginInput,
	// Passwordless Login
	VerifyLoginLinkInput, // Added for passwordless
	RequestLoginLinkInput, // Added for passwordless
	RegisterInput,
	ForgotPasswordInput,
	ResetPasswordInput,
	SocialTokenExchangeInput,
	EmailVerificationOutput,
	// Phone Verification
	VerifyPhoneInput,
	// 2FA Types
	Generate2FAResponse,
	Verify2FARequest,
	TwoFactorRecoveryCodesResponse,
	Disable2FARequest,
	AuthResult, // Explicitly import AuthResult
} from '@/lib/apiClient';
import { auth } from '@/lib/firebase';
import { signOut } from 'firebase/auth';

// Note: Token constants (ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY) and direct localStorage
// management functions (storeTokens, getTokens, clearTokens, dispatchStorageEvent)
// have been removed. Token adapters and retrieval are now expected to be handled
// by a React hook using `useLocalStorage` (e.g., within AuthContext).
// This service now focuses solely on API interactions related to authentication.

/**
 * Exchanges a social provider's token (e.g., Firebase ID token) for a custom JWT from the backend.
 * @param data The social token exchange data (SocialTokenExchangeInput).
 * @returns An AuthResponse containing tokens and potentially user/profile data.
 */
export const exchangeFirebaseToken = async (data: SocialTokenExchangeInput): Promise<AuthResponse> => {
	try {
		// Expect AuthResponse from the backend endpoint
		const response = await apiClient.post<AuthResponse>('/api/v1/auth/social-token-exchange', data); // Added /api/v1
		return response.data;
	} catch (error) {
		console.error('Error exchanging Firebase token:', error);
		throw error; // Re-throw the error to be handled by the caller
	}
};

/**
 * Signs in a user using email and password.
 * @param data Login credentials (LoginInput).
 * @returns A LoginOutput which might contain tokens or indicate 2FA is required.
 */
export const signInWithEmail = async (data: LoginInput): Promise<LoginOutput> => {
	try {
		// Expect LoginOutput from the backend /auth/login endpoint
		const response = await apiClient.post<LoginOutput>('/api/v1/auth/login', data); // Added /api/v1
		// AuthContext will handle storing tokens or prompting for 2FA
		return response.data;
	} catch (error) {
		console.error('Error signing in with email:', error);
		throw error; // Re-throw to be handled by the caller (AuthContext/LoginForm)
	}
};

/**
 * Verifies the 2FA code during login.
 * @param data The 2FA session token and code (Verify2FARequest).
 * @returns A LoginOutput containing the final tokens and user data upon success.
 */
export const verifyTwoFactorLogin = async (data: Verify2FARequest): Promise<LoginOutput> => {
	try {
		// Expect LoginOutput from the backend /auth/login/verify-2fa endpoint
		const response = await apiClient.post<LoginOutput>('/api/v1/auth/login/verify-2fa', data);
		// AuthContext will handle storing tokens on success
		return response.data;
	} catch (error) {
		console.error('Error verifying 2FA code during login:', error);
		throw error; // Re-throw to be handled by the caller (AuthContext)
	}
};

/**
 * Logs the user out by clearing tokens and signing out from Firebase.
 * Optionally calls the backend logout endpoint.
 */

/**
 * Registers a new user.
 * @param data The registration data (RegisterInput).
 * @returns An AuthResponse containing tokens and user data.
 */
export const register = async (data: RegisterInput): Promise<AuthResponse> => {
	try {
		const response = await apiClient.post<AuthResponse>('/api/v1/auth/register', data); // Added /api/v1
		// AuthContext will handle storing tokens and user state
		return response.data;
	} catch (error) {
		console.error('Error registering user:', error);
		throw error;
	}
};


/**
 * Sends a password reset email.
 * @param data The email address (ForgotPasswordInput).
 */
export const forgotPassword = async (data: ForgotPasswordInput): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/password/forgot', data); // Added /api/v1
		console.log('Forgot password request sent.');
	} catch (error) {
		console.error('Error sending forgot password request:', error);
		throw error;
	}
};

/**
 * Resets the user's password using a token.
 * @param data The reset token and new password (ResetPasswordInput).
 */
export const resetPassword = async (data: ResetPasswordInput): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/password/reset', data); // Added /api/v1
		console.log('Password reset successful.');
	} catch (error) {
		console.error('Error resetting password:', error);
		throw error;
	}
};

/**
 * Verifies the user's email using a token.
 * @param token The verification token from the URL/email link.
 * @returns EmailVerificationOutput (adjust based on backend response).
 */
export const verifyEmail = async (token: string): Promise<EmailVerificationOutput> => {
	try {
		// Backend expects GET request for email verification
		const response = await apiClient.get<EmailVerificationOutput>(`/api/v1/auth/email/verify/${token}`); // Added /api/v1
		console.log('Email verification successful.');
		return response.data;
	} catch (error) {
		console.error('Error verifying email:', error);
		throw error;
	}
};


/**
 * Refreshes the access token using the provided refresh token.
 * Uses the `/auth/refresh` endpoint. The caller is responsible for storing
 * the new tokens and handling failures (e.g., logging out).
 * @param currentRefreshToken The current refresh token.
 * @returns An AuthResponse containing the new tokens, or throws an error on failure.
 */
export const refreshToken = async (currentRefreshToken: string): Promise<AuthResponse> => {
	if (!currentRefreshToken) {
		console.error('refreshToken service: No refresh token provided.');
		// Throw an error or handle appropriately, depending on desired contract
		throw new Error('Refresh token is required.');
	}

	try {
		const response = await apiClient.post<AuthResponse>('/api/v1/auth/refresh', { // Added /api/v1
			refresh_token: currentRefreshToken,
		}, {
			headers: { '__skipAuthRefresh': 'true' } // Prevent interceptor loop
		});

		// Caller (e.g., AuthContext) is responsible for storing response.data.auth
		console.log('refreshToken service: Token refresh API call successful.');
		return response.data;
	} catch (error) {
		console.error('refreshToken service: Error during token refresh API call:', error);
		// Re-throw the error so the caller (e.g., AuthContext) can handle it,
		// potentially by logging the user out.
		throw error;
	}
};

/**
 * Logs the user out by notifying the backend and signing out from Firebase.
 * The caller (e.g., AuthContext) is responsible for clearing local tokens/state.
 */
export const logoutUser = async (): Promise<void> => {
	// Caller is responsible for clearing tokens (e.g., via useLocalStorage setter)
	console.log('logoutUser service: Initiating logout process.');

	try {
		// Notify backend about logout (best effort, don't block UI on failure)
		await apiClient.post('/api/v1/auth/logout', null, { // Added /api/v1
			headers: { '__skipAuthRefresh': 'true' } // Avoid potential issues if tokens were already cleared
		});
		console.log('logoutUser service: Backend logout notification sent.');
	} catch (error) {
		// Log non-critical error
		console.warn('logoutUser service: Error during backend logout notification:', error);
	}

	// Sign out from Firebase (best effort)
	try {
		await signOut(auth);
		console.log('logoutUser service: Firebase sign-out successful.');
	} catch (firebaseError) {
		console.error('logoutUser service: Error signing out from Firebase:', firebaseError);
		// Decide if this should be re-thrown or just logged
	}
	// No redirect here, handled by UI/AuthContext based on auth state change
};

// User profile management functions have been moved to userService.ts

// --- Phone Verification ---

/**
 * Requests a phone verification OTP to be sent to the authenticated user's phone.
 */
export const requestPhoneVerification = async (): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/verify-phone/request');
		console.log('Phone verification OTP request sent.');
	} catch (error) {
		console.error('Error requesting phone verification:', error);
		throw error;
	}
};

/**
 * Verifies the phone number using the provided OTP.
 * @param data The OTP code (VerifyPhoneInput).
 */
export const verifyPhone = async (data: VerifyPhoneInput): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/verify-phone/confirm', data);
		console.log('Phone verified successfully.');
	} catch (error) {
		console.error('Error verifying phone:', error);
		throw error;
	}
};

// --- Two-Factor Authentication (2FA) ---

/**
 * Generates a new 2FA secret and QR code URI for the authenticated user.
 * @returns Generate2FAResponse containing the secret and QR code data URI.
 */
export const generate2FASecret = async (): Promise<Generate2FAResponse> => {
	try {
		const response = await apiClient.get<Generate2FAResponse>('/api/v1/auth/2fa/generate');
		console.log('2FA secret generated.');
		console.log('generate2FASecret returning:', response.data); // Log what's being returned
		return response.data;
	} catch (error) {
		console.error('Error generating 2FA secret:', error);
		throw error;
	}
};

/**
 * Enables 2FA for the authenticated user by verifying the initial TOTP code.
 * @param data The verification code (Verify2FARequest).
 * @returns TwoFactorRecoveryCodesResponse containing recovery codes.
 */
export const enable2FA = async (data: Verify2FARequest): Promise<TwoFactorRecoveryCodesResponse> => {
	try {
		const response = await apiClient.post<TwoFactorRecoveryCodesResponse>('/api/v1/auth/2fa/enable', data);
		console.log('2FA enabled successfully.');
		return response.data;
	} catch (error) {
		console.error('Error enabling 2FA:', error);
		throw error;
	}
};

/**
 * Disables 2FA for the authenticated user. Requires password or current TOTP code.
 * @param data Verification data (Disable2FARequest).
 */
export const disable2FA = async (data: Disable2FARequest): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/2fa/disable', data);
		console.log('2FA disabled successfully.');
	} catch (error) {
		console.error('Error disabling 2FA:', error);
		throw error;
	}
};

/**
 * Verifies a login link token for passwordless authentication.
 * @param data The token from the login link (VerifyLoginLinkInput).
 * @returns A LoginOutput containing tokens and user data upon success.
 */
export const verifyLoginLink = async (data: VerifyLoginLinkInput): Promise<LoginOutput> => {
	try {
		const response = await apiClient.post<LoginOutput>('/api/v1/auth/login/verify-link', data);
		// AuthContext will handle storing tokens and user state
		return response.data;
	} catch (error) {
		console.error('Error verifying login link:', error);
		throw error;
	}
};

/**
 * Requests a login link to be sent to the user's email for passwordless authentication.
 * @param data The email and optional tenant slug (RequestLoginLinkInput).
 * @returns Promise<void>
 */
export const requestLoginLink = async (data: RequestLoginLinkInput): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/login/request-link', data);
		console.log('Request login link email sent.');
	} catch (error) {
		console.error('Error requesting login link:', error);
		throw error;
	}
};

/**
 * Resends the email verification link to the authenticated user's email.
 * @returns Promise<void>
 */
export const resendVerificationEmail = async (): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/email/resend-verify');
		console.log('Verification email resent successfully.');
	} catch (error) {
		console.error('Error resending verification email:', error);
		throw error;
	}
};

/**
 * Login the active tenant context for the authenticated user.
 * @param tenantId The ID of the tenant to switch to.
 * @returns An AuthResult containing new tokens for the selected tenant context.
 */
export const loginTenantContext = async (tenantId: string): Promise<AuthResult> => {
	try {
		const response = await apiClient.post<AuthResult>('/api/v1/auth/login-tenant', { tenant_id: tenantId });
		// The AuthContext will use this AuthResult to update tokens and user state.
		console.log('Tenant context switched successfully.');
		return response.data;
	} catch (error) {
		console.error('Error switching tenant context:', error);
		throw error; // Re-throw to be handled by the caller (AuthContext)
	}
};
