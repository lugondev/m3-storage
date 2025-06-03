import axios, { AxiosError, InternalAxiosRequestConfig, AxiosResponse, AxiosHeaders } from 'axios';
import { refreshToken as serviceRefreshToken } from '@/services/authService';

// Define localStorage keys here as they are no longer in authService
const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

// --- Start: Backend Type Definitions ---

// Media-related types based on actual backend routes
export interface MediaItem {
	id: string;
	user_id: string;
	file_name: string;
	file_path: string;
	file_size: number;
	media_type: string;
	provider: string;
	public_url: string;
	uploaded_at: string;
	created_at: string;
	updated_at: string;
}

export interface UploadResponse {
	media: MediaItem;
	message?: string;
}

export interface MediaListResponse {
	data: MediaItem[];
	total?: number;
	page?: number;
	page_size?: number;
}

// Storage health check types
export interface StorageHealthResponse {
	status: string;
	message?: string;
	provider?: string;
}

export interface AllStorageHealthResponse {
	providers: Record<string, StorageHealthResponse>;
	overall_status: string;
}

// Auth-related types (minimal, based on auth service import)
export interface AuthResult {
	access_token: string;
	refresh_token?: string;
	expires_in?: number;
}

// --- End: Backend Type Definitions ---

function generateUUID() {
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
		const r = Math.random() * 16 | 0;
		const v = c === 'x' ? r : (r & 0x3 | 0x8);
		return v.toString(16);
	});
}

// Create separate clients for auth and app services
const authApiClient = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_AUTH_URL || 'http://localhost:8080/api/v1',
	headers: {
		'Content-Type': 'application/json',
		'X-Request-ID': generateUUID(),
	},
	withCredentials: true,
});

const appApiClient = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_APP_URL || 'http://localhost:8083/api/v1',
	headers: {
		'Content-Type': 'application/json',
		'X-Request-ID': generateUUID(),
	},
	withCredentials: true,
});

// Keep the original apiClient for backward compatibility, defaulting to auth URL
const apiClient = appApiClient;

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

// --- API Methods for Media Operations ---

// Media upload
export const uploadMedia = async (file: File): Promise<UploadResponse> => {
	const formData = new FormData();
	formData.append('file', file);

	const response = await apiClient.post<UploadResponse>('/media/upload', formData, {
		headers: {
			'Content-Type': 'multipart/form-data',
		},
	});
	return response.data;
};

// List media
export const listMedia = async (): Promise<MediaListResponse> => {
	const response = await apiClient.get<MediaListResponse>('/media');
	return response.data;
};

// Get single media item
export const getMedia = async (id: string): Promise<MediaItem> => {
	const response = await apiClient.get<MediaItem>(`/media/${id}`);
	return response.data;
};

// Delete media
export const deleteMedia = async (id: string): Promise<void> => {
	await apiClient.delete(`/media/${id}`);
};

// --- API Methods for Storage Operations ---

// Check storage health
export const checkStorageHealth = async (): Promise<StorageHealthResponse> => {
	const response = await apiClient.get<StorageHealthResponse>('/storage/health');
	return response.data;
};

// Check all storage providers health
export const checkAllStorageHealth = async (): Promise<AllStorageHealthResponse> => {
	const response = await apiClient.get<AllStorageHealthResponse>('/storage/health/all');
	return response.data;
};

export default apiClient;
export { authApiClient, appApiClient };
