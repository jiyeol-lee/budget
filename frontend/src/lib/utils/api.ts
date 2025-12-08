/**
 * API Client for Budget Tracker
 * Provides typed fetch wrapper for backend API calls
 */

/**
 * Dynamically determine the API base URL based on environment and context.
 * This allows the app to work when:
 * - Development mode: Uses dynamic hostname with port 8080 (supports mobile access)
 * - Production mode: Uses same origin (nginx proxies /api to backend)
 * - SSR: Uses VITE_API_URL or falls back to localhost
 */
const getBaseUrl = (): string => {
	// Server-side rendering
	if (typeof window === 'undefined') {
		return import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
	}

	// Client-side: Check for explicit API URL override
	const envUrl = import.meta.env.VITE_API_URL;

	// Development mode: use same hostname with port 8080
	// This allows mobile access (192.168.x.x:5173 -> 192.168.x.x:8080/api)
	if (import.meta.env.DEV) {
		if (envUrl) {
			// If envUrl is 'auto', use dynamic hostname with port 8080
			if (envUrl === 'auto') {
				const { protocol, hostname } = window.location;
				return `${protocol}//${hostname}:8080/api`;
			}
			return envUrl;
		}
		// Default dev behavior: dynamic hostname with port 8080
		const { protocol, hostname } = window.location;
		return `${protocol}//${hostname}:8080/api`;
	}

	// Production: use same origin (nginx proxy)
	return `${window.location.origin}/api`;
};

const BASE_URL = getBaseUrl();

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
	constructor(
		public status: number,
		public statusText: string,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

/**
 * Extracts a human-readable error message from an API response text.
 * Handles JSON responses with an 'error' field (e.g., {"error": "message"})
 * and falls back to the original text for non-JSON responses.
 */
export function extractErrorMessage(text: string): string {
	if (!text) {
		return 'Unknown error';
	}

	try {
		const json = JSON.parse(text);
		if (typeof json.error === 'string') {
			return json.error;
		}
		// If JSON but no error field, return original text
		return text;
	} catch {
		// Not valid JSON, return as-is (plain text error)
		return text;
	}
}

/**
 * Generic response handler
 */
async function handleResponse<T>(response: Response): Promise<T> {
	if (!response.ok) {
		const rawText = await response.text().catch(() => 'Unknown error');
		const errorMessage = extractErrorMessage(rawText);
		throw new ApiError(response.status, response.statusText, errorMessage);
	}

	// Handle empty responses
	const text = await response.text();
	if (!text) {
		return {} as T;
	}

	try {
		return JSON.parse(text) as T;
	} catch {
		throw new ApiError(response.status, 'Parse Error', 'Failed to parse JSON response');
	}
}

/**
 * GET request
 */
export async function get<T>(endpoint: string): Promise<T> {
	const response = await fetch(`${BASE_URL}${endpoint}`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	});
	return handleResponse<T>(response);
}

/**
 * POST request
 */
export async function post<T, D = unknown>(endpoint: string, data?: D): Promise<T> {
	const response = await fetch(`${BASE_URL}${endpoint}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: data ? JSON.stringify(data) : undefined
	});
	return handleResponse<T>(response);
}

/**
 * PUT request
 */
export async function put<T, D = unknown>(endpoint: string, data?: D): Promise<T> {
	const response = await fetch(`${BASE_URL}${endpoint}`, {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json'
		},
		body: data ? JSON.stringify(data) : undefined
	});
	return handleResponse<T>(response);
}

/**
 * DELETE request
 */
export async function del<T>(endpoint: string): Promise<T> {
	const response = await fetch(`${BASE_URL}${endpoint}`, {
		method: 'DELETE',
		headers: {
			'Content-Type': 'application/json'
		}
	});
	return handleResponse<T>(response);
}

/**
 * POST request with FormData (for file uploads)
 */
export async function postFormData<T>(endpoint: string, formData: FormData): Promise<T> {
	const response = await fetch(`${BASE_URL}${endpoint}`, {
		method: 'POST',
		body: formData
	});
	return handleResponse<T>(response);
}

// Export all methods as named exports
export const api = {
	get,
	post,
	put,
	delete: del,
	postFormData
};

export default api;
