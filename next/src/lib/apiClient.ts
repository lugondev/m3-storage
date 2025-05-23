import axios, { AxiosError, InternalAxiosRequestConfig, AxiosResponse, AxiosHeaders } from 'axios';
// Import only refreshToken
// getTokens is removed, tokens read directly from localStorage here.
// serviceLogout is no longer called directly from here.
import { refreshToken as serviceRefreshToken } from '@/services/authService';

// Define localStorage keys here as they are no longer in authService
const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

// --- Start: Backend Type Definitions ---

// Based on internal/modules/account/domain/models.go
export type UserStatus = "active" | "pending" | "suspended" | "deleted";

// Based on internal/modules/account/domain/role.go
// Assuming Permission is a string for simplicity, adjust if it's more complex
export type Permission = string;

export interface RoleResponse {
	roles: string[]
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
	is_email_verified: boolean;
	is_phone_verified: boolean;
	profile?: UserProfile | null; // Added optional profile field
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

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

// Based on internal/modules/account/domain/auth_dto.go
export interface RegisterInput {
	email: string;
	password: string;
	first_name?: string; // Optional based on typical registration
	last_name?: string;  // Optional
}

export interface LoginInput {
	email: string;
	password: string;
	tenant_slug?: string; // Added tenant_slug as optional
}

export interface SocialTokenExchangeInput {
	token: string; // The ID token from the social provider (e.g., Firebase)
	provider: string; // Added provider field, e.g., "firebase", "google"
}

export interface VerifyLoginLinkInput {
	token: string;
}

export interface RequestLoginLinkInput {
	email: string;
	tenant_slug?: string;
}

export interface ForgotPasswordInput {
	email: string;
}

export interface ResetPasswordInput {
	token: string;
	new_password: string;
}

export interface EmailVerificationOutput {
	message: string; // Example: "Email verified successfully"
}

// --- Phone Verification DTOs ---
export interface VerifyPhoneInput {
	otp: string;
}

// --- 2FA DTOs ---
export interface Generate2FAResponse {
	secret: string;
	qr_code_uri: string; // Data URI for QR code image (Corrected to snake_case)
}

// Removed the empty interface definition above

// Added based on API spec for POST /api/v1/auth/login/verify-2fa
export interface Verify2FARequest {
	two_factor_session_token: string; // The token received from the initial login step
	code: string; // The TOTP code from the authenticator app
}

export interface TwoFactorRecoveryCodesResponse {
	recovery_codes: string[];
}

export interface Disable2FARequest {
	password?: string; // Optional: Current user password
	code?: string; // Optional: Current TOTP code
}

// Input DTOs for User operations (Keep existing)
export interface UpdateUserInput {
	first_name?: string;
	last_name?: string;
	phone?: string;
	status?: UserStatus; // Make optional for partial updates
}

// Added for user update requests
export interface UpdateUserRequest {
	first_name?: string;
	last_name?: string;
	phone?: string;
	status?: UserStatus; // Updated based on user feedback
	// Add other fields that can be updated via the API
}

export interface UpdateProfileInput {
	bio?: string;
	date_of_birth?: string | null; // Use string for date input, backend parses
	address?: string;
	interests?: string[];
	preferences?: Partial<UserPreferences>; // Allow partial updates for preferences
}

export interface UpdatePasswordInput {
	current_password: string;
	new_password: string;
	otp_code?: string; // Optional: Required if 2FA is enabled
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


// Paginated results for users and roles
export interface PaginatedUsers {
	users: UserOutput[];
	total: number; // Assuming int64 maps to number
	page: number;
	page_size: number;
	total_pages: number;
}

// Search query parameters
export interface UserSearchQuery {
	query?: string;
	status?: UserStatus;
	role_name?: string; // uuid.UUID
	offset?: number;
	limit?: number;
}

// Type for UpdatePasswordResponse in user_handler.go
export interface UpdatePasswordResponse {
	message: string;
}

// --- Start: Venue Module Type Definitions ---

// Based on internal/modules/venue/domain/venue.go
export type ApprovalStatus = "pending" | "active" | "inactive" | "rejected";

export interface VenuePhoto {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	url: string;
	caption: string;
	is_primary: boolean;
	created_at: string; // time.Time
}

export interface Venue {
	id: string; // uuid.UUID
	name: string;
	description: string;
	address: string;
	city: string;
	country: string;
	postal_code: string;
	latitude: number; // float64
	longitude: number; // float64
	phone_number: string;
	email: string;
	website: string;
	capacity: number; // int
	status: ApprovalStatus; // string -> ApprovalStatus
	category_id?: string | null; // *uuid.UUID
	owner_id: string; // uuid.UUID
	photos: VenuePhoto[];
	created_at: string; // time.Time
	updated_at: string; // time.Time
	deleted_at?: string | null; // *time.Time
}

export interface CreateVenueInput {
	name: string;
	description: string;
	address: string;
	city: string;
	country: string;
	postal_code: string;
	latitude: number;
	longitude: number;
	phone_number: string;
	email: string;
	website: string;
	capacity: number;
	category_id?: string | null; // Optional category ID
}

export interface UpdateVenueInput extends Partial<Omit<CreateVenueInput, 'owner_id'>> {
	// All fields are optional for PATCH
	status?: ApprovalStatus; // Allow updating status
}

export interface TransferVenueOwnershipInput {
	new_owner_id: string; // uuid.UUID
}

// Based on internal/modules/venue/domain/venue_staff.go
export type StaffRole = "owner" | "manager" | "staff" | "hostess" | "waiter" | "bartender";
export type StaffStatus = "active" | "inactive" | "pending";
// Assuming Venue Permission is a string for simplicity, map from Go constants
export type VenuePermission =
	| "manage_venue"
	| "manage_staff"
	| "manage_settings"
	| "manage_tables"
	| "manage_events"
	| "manage_products"
	| "manage_reservation"
	| "manage_orders"
	| "manage_promotions"
	| "view_reports";

export interface VenueStaff {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	user_id: string; // uuid.UUID
	role: StaffRole;
	permissions: VenuePermission[];
	status: StaffStatus;
	created_at: string; // time.Time
	updated_at: string; // time.Time
	deleted_at?: string | null; // *time.Time
}

export interface AddVenueStaffInput {
	user_id: string;
	role: StaffRole;
	permissions?: VenuePermission[]; // Optional, might be based on role
}

export interface UpdateVenueStaffInput {
	role?: StaffRole;
	permissions?: VenuePermission[];
	status?: StaffStatus;
}

// Based on internal/modules/venue/domain/venue_settings.go
export interface BusinessHour {
	day_of_week: number; // 0-6
	open_time: string; // Store as HH:mm string for simplicity? or full time.Time string
	close_time: string; // Store as HH:mm string for simplicity? or full time.Time string
	is_closed: boolean;
}

export interface BookingSettings {
	enable_booking: boolean;
	max_booking_per_time_slot: number;
	booking_lead_hours: number;
	booking_duration_minutes: number;
	require_approval: boolean;
}

export interface TierLevel {
	name: string;
	points_required: number;
	discount: number;
	perks: string;
}

export interface LoyaltySettings {
	enable_loyalty: boolean;
	points_per_purchase: number;
	points_redemption_rate: number;
	tier_levels: TierLevel[];
}

export interface AffiliateSettings {
	enable_affiliate: boolean;
	commission_rate: number;
	cookie_days: number;
	min_payout: number;
}

export interface VenueSettings {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	business_hours: BusinessHour[];
	time_zone: string;
	currency: string;
	booking_settings: BookingSettings;
	loyalty_settings: LoyaltySettings;
	affiliate_settings: AffiliateSettings;
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

// Input types for updating settings (likely partial updates)
export interface UpdateVenueSettingsInput {
	time_zone?: string;
	currency?: string;
	business_hours?: BusinessHour[]; // Sending the full array is often simpler for updates
	booking_settings?: Partial<BookingSettings>;
	loyalty_settings?: Partial<LoyaltySettings>; // Consider nested partials or full object replacement
	affiliate_settings?: Partial<AffiliateSettings>;
}

// Based on internal/modules/venue/domain/event.go
export type EventCategory = "music" | "sports" | "arts" | "food" | "business" | "conference" | "other";
export type EventStatus = "draft" | "published" | "cancelled" | "completed";

export interface EventPhoto {
	id: string; // uuid.UUID
	event_id: string; // uuid.UUID
	url: string;
	caption: string;
	is_primary: boolean;
	created_at: string; // time.Time
}

export interface Event {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	name: string;
	description: string;
	category: EventCategory;
	start_time: string; // time.Time
	end_time: string; // time.Time
	time_zone: string;
	is_recurring: boolean;
	recurrence_rule?: string; // iCal format
	max_capacity: number;
	ticket_price: number;
	is_featured: boolean;
	is_cancelled: boolean;
	photos: EventPhoto[];
	status: EventStatus;
	created_at: string; // time.Time
	updated_at: string; // time.Time
	deleted_at?: string | null; // *time.Time
	// Consider adding performers and tickets if they are part of the main Event response
}

export interface CreateEventInput {
	name: string;
	description: string;
	category: EventCategory;
	start_time: string; // ISO 8601 string
	end_time: string; // ISO 8601 string
	time_zone: string; // e.g., "Asia/Ho_Chi_Minh"
	max_capacity: number;
	ticket_price: number;
	is_recurring?: boolean;
	recurrence_rule?: string;
	is_featured?: boolean;
}

export interface UpdateEventInput extends Partial<CreateEventInput> {
	status?: EventStatus; // Allow updating status
	is_cancelled?: boolean;
}

// Based on internal/modules/venue/domain/product.go
export type ProductCategory = "appetizer" | "main" | "dessert" | "drink" | "alcohol" | "side" | "special";

export interface ProductPhoto {
	id: string; // uuid.UUID
	product_id: string; // uuid.UUID
	url: string;
	caption?: string;
	is_primary: boolean;
	created_at: string; // time.Time
}

export interface OptionChoice {
	id: string; // uuid.UUID
	option_id: string; // uuid.UUID
	name: string;
	description?: string;
	price_adjustment: number;
	is_default: boolean;
}

export interface ProductOption {
	id: string; // uuid.UUID
	product_id: string; // uuid.UUID
	name: string;
	description?: string;
	required: boolean;
	min_select: number;
	max_select: number;
	choices: OptionChoice[];
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

export interface NutritionalInfo {
	calories?: number;
	protein?: number;
	carbohydrates?: number;
	fat?: number;
	sodium?: number;
	sugar?: number;
	fiber?: number;
}

export interface Product {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	name: string;
	description: string;
	category: ProductCategory;
	price: number;
	discount_price?: number | null; // *float64
	currency: string;
	is_available: boolean;
	sku?: string;
	tags?: string[];
	photos: ProductPhoto[];
	options?: ProductOption[];
	ingredients?: string[];
	allergens?: string[];
	nutritional_info?: NutritionalInfo | null; // *NutritionalInfo
	created_at: string; // time.Time
	updated_at: string; // time.Time
	deleted_at?: string | null; // *time.Time
}

export interface CreateProductInput {
	name: string;
	description: string;
	category: ProductCategory;
	price: number;
	currency: string; // Should likely come from venue settings?
	discount_price?: number | null;
	is_available?: boolean; // Default to true?
	sku?: string;
	tags?: string[];
	ingredients?: string[];
	allergens?: string[];
	nutritional_info?: NutritionalInfo | null;
	// Options and Photos added separately?
}

export interface UpdateProductInput extends Partial<CreateProductInput> {
	is_available?: boolean;
	// Handle options/photos updates separately
}


// Based on internal/modules/venue/domain/table.go
export type TableStatus = "available" | "occupied" | "reserved" | "out_of_service";
export type TableType = "standard" | "bar" | "private" | "outdoor" | "lounge" | "vip";

export interface Table {
	id: string; // uuid.UUID
	venue_id: string; // uuid.UUID
	name: string;
	description: string;
	capacity: number;
	status: TableStatus;
	location: string;
	min_spend?: number;
	table_type: TableType;
	is_active: boolean;
	created_at: string; // time.Time
	updated_at: string; // time.Time
	deleted_at?: string | null; // *time.Time
}

export interface CreateTableInput {
	name: string;
	description: string;
	capacity: number;
	location: string;
	table_type: TableType;
	min_spend?: number;
	is_active?: boolean; // Default to true?
}

export interface UpdateTableInput extends Partial<CreateTableInput> {
	status?: TableStatus; // Status updated separately?
	is_active?: boolean;
}

// --- Start: Tenant Module Type Definitions ---

// Based on internal/modules/tenant/domain/dto.go

export interface CreateTenantRequest {
	name: string;
	slug: string;
	owner_email: string;
}

export interface TenantResponse {
	id: string; // uuid.UUID
	name: string;
	slug: string;
	owner_user_id: string; // uuid.UUID
	is_active: boolean;
	created_at: string; // time.Time
	updated_at: string; // time.Time
}

export interface UpdateTenantRequest {
	name?: string;
	is_active?: boolean;
}

export interface MinimalTenantInfo {
	id: string; // uuid.UUID
	name: string;
	slug: string;
}

export interface AddUserToTenantRequest {
	email: string;
	role_ids: string[]; // uuid.UUID array
}

export interface TenantUserResponse {
	user_id: string; // uuid.UUID
	email: string;
	first_name: string;
	last_name: string;
	avatar?: string;
	status_in_tenant: string; // e.g., "active", "invited"
	global_status: string;    // e.g., "active", "pending"
	roles: string[];          // Role names
	joined_at: string;        // time.Time
}

export interface UpdateTenantUserRequest {
	role_ids?: string[]; // uuid.UUID array, pointer in Go means optional
	status_in_tenant?: string; // e.g., "active", "invited", "suspended"
}

// Paginated responses for Tenant resources
export interface PaginatedTenants {
	tenants: TenantResponse[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface PaginatedTenantUsers {
	users: TenantUserResponse[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface UserTenantMembershipInfo {
	tenant_id: string;
	tenant_name: string;
	tenant_slug: string;
	tenant_is_active: boolean;
	user_roles: string[];
	user_status: string;
	joined_at: string;
}

// --- End: Tenant Module Type Definitions ---

// Paginated responses for Venue resources
export interface PaginatedVenues {
	venues: Venue[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface PaginatedVenueStaff {
	staff: VenueStaff[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface PaginatedEvents {
	events: Event[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface PaginatedProducts {
	products: Product[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

export interface PaginatedTables {
	tables: Table[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

// Generic search query parameters
export interface BaseSearchQuery {
	page?: number;
	page_size?: number;
	query?: string; // Generic search term
	// Add other common filters like sort_by, sort_order if applicable
}

export interface VenueSearchQuery extends BaseSearchQuery {
	status?: ApprovalStatus;
	category_id?: string;
	owner_id?: string;
	city?: string;
	country?: string;
}

export interface VenueStaffSearchQuery extends BaseSearchQuery {
	role?: StaffRole;
	status?: StaffStatus;
}

export interface EventSearchQuery extends BaseSearchQuery {
	category?: EventCategory;
	status?: EventStatus;
	start_date?: string; // ISO 8601 date string
	end_date?: string; // ISO 8601 date string
	is_featured?: boolean;
}

export interface ProductSearchQuery extends BaseSearchQuery {
	category?: ProductCategory;
	is_available?: boolean;
	min_price?: number;
	max_price?: number;
	tags?: string[]; // Comma-separated string or array? Check API
}

export interface TableSearchQuery extends BaseSearchQuery {
	status?: TableStatus;
	type?: TableType;
	location?: string;
	min_capacity?: number;
	is_active?: boolean;
}


function generateUUID() {
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
		const r = Math.random() * 16 | 0;
		const v = c === 'x' ? r : (r & 0x3 | 0x8);
		return v.toString(16);
	});
}


const apiClient = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1',
	headers: {
		'Content-Type': 'application/json',
		'X-Request-ID': generateUUID(),
	},
	withCredentials: true,
});

let isRefreshing = false;
type FailedQueueItem = {
	resolve: (value: unknown) => void;
	reject: (reason?: AxiosError | null) => void;
};
let failedQueue: FailedQueueItem[] = [];

const processQueue = (error: AxiosError | null, token: string | null = null) => {
	failedQueue.forEach(prom => {
		if (error) {
			prom.reject(error);
		} else {
			prom.resolve(token);
		}
	});
	failedQueue = [];
};

// Request Interceptor: Add JWT to headers
apiClient.interceptors.request.use(
	(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
		// Read access token directly from localStorage
		const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY)?.replaceAll('"', "");

		const skipAuth = config.headers?.['__skipAuthRefresh'] === 'true';
		if (accessToken && !skipAuth) {
			if (!config.headers) {
				config.headers = new AxiosHeaders();
			}
			config.headers.set('Authorization', `Bearer ${accessToken}`);
		}

		if (skipAuth) {
			delete config.headers['__skipAuthRefresh'];
		}

		return config;
	},
	(error: AxiosError) => Promise.reject(error)
);

// Response Interceptor: Handle 401 errors and token refresh
apiClient.interceptors.response.use(
	(response: AxiosResponse) => response,
	async (error: AxiosError) => {
		const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

		// Use status code from response if available, otherwise check error message potentially
		const statusCode = error.response?.status;

		// Check if it's a 401 error, not a refresh token failure itself, and not already retried
		// Backend refresh endpoint is /auth/refresh
		if (statusCode === 401 && originalRequest.url !== '/auth/refresh' && !originalRequest._retry) {

			if (isRefreshing) {
				return new Promise((resolve, reject) => {
					failedQueue.push({ resolve, reject });
				}).then(token => {
					if (originalRequest.headers) {
						originalRequest.headers['Authorization'] = `Bearer ${token}`;
					}
					// Ensure the retry uses the original method, data, etc.
					return apiClient(originalRequest);
				}).catch(err => {
					return Promise.reject(err); // Propagate refresh error
				});
			}

			originalRequest._retry = true;
			isRefreshing = true;

			// Read refresh token directly from localStorage
			const currentRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)?.replaceAll('"', "");

			if (!currentRefreshToken) {
				console.log('Interceptor: No refresh token found, redirecting to login.');
				processQueue(error, null); // Reject pending requests
				// Redirect to login page
				if (typeof window !== 'undefined') {
					window.location.href = '/login';
				}
				return Promise.reject(error); // Reject the original request after triggering redirect
			}

			try {
				// Call the service function, passing the refresh token
				const refreshResponse = await serviceRefreshToken(currentRefreshToken);

				// Check access_token within the nested 'auth' object
				const newAccessToken = refreshResponse?.auth.access_token;
				if (newAccessToken) {
					// Note: The hook useLocalStorage will automatically update localStorage
					// via the handleAuthenticationSuccess callback in AuthContext upon success.
					// We just need the new token to retry the original request.
					console.log('Interceptor: Token refreshed successfully. Retrying original request.');
					if (originalRequest.headers) {
						originalRequest.headers['Authorization'] = `Bearer ${newAccessToken}`;
					}
					processQueue(null, newAccessToken); // Resolve pending requests with the new token
					return apiClient(originalRequest); // Retry the original request
				} else {
					// This case should ideally not happen if serviceRefreshToken works correctly
					console.error('Interceptor: Refresh endpoint returned response without a new access token.');
					processQueue(error, null); // Reject pending requests
					// No need to call logoutUser here
					return Promise.reject(error); // Reject the original request
				}
			} catch (refreshError) {
				console.error('Interceptor: Failed to refresh token via service:', refreshError);
				processQueue(refreshError as AxiosError, null); // Reject pending requests with the refresh error
				// Let AuthContext handle the logout state/UI changes upon catching this error.
				// Redirect to login page on refresh failure
				if (typeof window !== 'undefined') {
					// Avoid redirect loop if the refresh endpoint itself returns 401 repeatedly
					// Although the main check already prevents retrying /auth/refresh
					if ((refreshError as AxiosError).response?.status === 401) {
						console.warn("Refresh token endpoint returned 401, potential issue with refresh token itself.");
					}
					window.location.href = '/login';
				}
				return Promise.reject(refreshError); // Reject the original request with the refresh error
			} finally {
				isRefreshing = false;
			}
		}

		// Log details for non-401 or already retried errors
		console.error('API call error:', error.response?.data || error.message);
		// Optionally handle specific error types (e.g., 403 Forbidden, 404 Not Found)
		// if (statusCode === 403) { ... }

		return Promise.reject(error);
	}
);

export default apiClient;


// NOTE: Removed redundant interface definitions at the end as they are now defined at the top.
// Interfaces like Venue, CreateVenueInput, etc., should be moved to a separate file
// (e.g., src/services/venueService.ts or src/types/venue.ts) if they grow complex.
