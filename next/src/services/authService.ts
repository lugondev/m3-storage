import apiClient from '@/lib/apiClient';
import { signOut } from 'firebase/auth';
import { auth } from '@/lib/firebase';

export interface Verify2FARequest {
	two_factor_session_token: string; // The token received from the initial login step
	code: string; // The TOTP code from the authenticator app
}

export interface SocialTokenExchangeInput {
	token: string; // The ID token from the social provider (e.g., Firebase)
	provider: string; // Added provider field, e.g., "firebase", "google"
}

export interface LoginInput {
	email: string;
	password: string;
}

// Based on internal/modules/account/domain/auth_dto.go
export interface RegisterInput {
	email: string;
	password: string;
	first_name?: string; // Optional based on typical registration
	last_name?: string;  // Optional
}

// Based on internal/modules/account/domain/models.go
export type UserStatus = "active" | "pending" | "blocked" | "inactive";

export interface UserPreferences {
	email_notifications: boolean;
	push_notifications: boolean;
	language: string;
	theme: string;
}

export interface UserProfile {
	id: string; // uuid.UUID
	user_id: string; // uuid.UUID
	bio?: string;
	date_of_birth?: string | null; // time.Time can be null/zero
	address?: string;
	interests?: string[];
	preferences?: UserPreferences; // Can be optional or have defaults
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

// Based on internal/modules/account/domain/user_dto.go and models.go
export interface UserOutput {
	id: string; // uuid.UUID
	email: string;
	first_name: string;
	last_name: string;
	phone?: string;
	avatar?: string;
	roles?: string[];
	status: UserStatus;
	email_verified_at?: string | null; // Reverted back to timestamp
	phone_verified_at?: string | null; // Reverted back to timestamp
	is_two_factor_enabled: boolean; // Added for 2FA status
	profile?: UserProfile | null; // Added optional profile field
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

// Based on internal/modules/account/domain/auth_dto.go
export interface AuthResult {
	access_token: string;
	refresh_token?: string;
	expires_in?: number; // Optional: Expiry time if provided
}

// Combined Login/Register/Refresh/Exchange Output
// Combining AuthResult and potential user/profile data returned on login/exchange
export interface AuthResponse {
	auth: AuthResult;
	user?: UserOutput; // User might be returned on login/register/refresh
	profile?: UserProfile; // Profile might be returned on login/register/refresh
}

// Added based on API spec for POST /api/v1/auth/login and POST /api/v1/auth/login/verify-2fa
export interface LoginOutput {
	user?: UserOutput; // User details might be returned
	auth?: AuthResult | null; // Auth tokens (null if 2FA is required initially)
	two_factor_required: boolean; // Flag indicating if 2FA step is needed
	two_factor_session_token?: string; // Token for the 2FA verification step (only if two_factor_required is true)
}

export const exchangeFirebaseToken = async (data: SocialTokenExchangeInput): Promise<AuthResponse> => {
	try {
		// Expect AuthResponse from the backend endpoint
		const response = await apiClient.post<AuthResponse>('/auth/social-token-exchange', data);
		return response.data;
	} catch (error) {
		console.error('Error exchanging Firebase token:', error);
		throw error; // Re-throw the error to be handled by the caller
	}
};

export const signInWithEmail = async (data: LoginInput): Promise<LoginOutput> => {
	try {
		// Expect LoginOutput from the backend /auth/login endpoint
		const response = await apiClient.post<LoginOutput>('/auth/login', data);
		// AuthContext will handle storing tokens or prompting for 2FA
		return response.data;
	} catch (error) {
		console.error('Error signing in with email:', error);
		throw error; // Re-throw to be handled by the caller (AuthContext/LoginForm)
	}
};

export const refreshToken = async (currentRefreshToken: string): Promise<AuthResponse> => {
	if (!currentRefreshToken) {
		console.error('refreshToken service: No refresh token provided.');
		// Throw an error or handle appropriately, depending on desired contract
		throw new Error('Refresh token is required.');
	}

	try {
		const response = await apiClient.post<AuthResponse>('/auth/refresh', {
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

export const register = async (data: RegisterInput): Promise<AuthResponse> => {
	try {
		const response = await apiClient.post<AuthResponse>('/auth/register', data);
		return response.data;
	} catch (error) {
		console.error('Error during registration:', error);
		throw error; // Re-throw to be handled by the caller
	}
};

export const verifyTwoFactorLogin = async (data: Verify2FARequest): Promise<LoginOutput> => {
	try {
		const response = await apiClient.post<LoginOutput>('/auth/login/verify-2fa', data);
		return response.data;
	} catch (error) {
		console.error('Error during 2FA verification:', error);
		throw error; // Re-throw to be handled by the caller
	}
};

export const logoutUser = async (): Promise<void> => {
	// Caller is responsible for clearing tokens (e.g., via useLocalStorage setter)
	console.log('logoutUser service: Initiating logout process.');

	try {
		// Notify backend about logout (best effort, don't block UI on failure)
		await apiClient.post('/auth/logout', null, {
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


export interface RequestLoginLinkInput {
	email: string;
}

export const requestLoginLink = async (data: RequestLoginLinkInput): Promise<void> => {
	try {
		await apiClient.post('/api/v1/auth/login/request-link', data);
		console.log('Request login link email sent.');
	} catch (error) {
		console.error('Error requesting login link:', error);
		throw error;
	}
};