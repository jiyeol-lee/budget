/**
 * API Client for Budget Tracker
 * Provides typed fetch wrapper for backend API calls
 */

/**
 * Dynamically determine the API base URL based on the current hostname.
 * This allows the app to work when:
 * - Accessing from localhost (development)
 * - Accessing from LAN IP (mobile testing)
 * - Accessing from a deployed domain (production)
 */
function getBaseUrl(): string {
	// Check if we're in a browser environment
	if (typeof window !== 'undefined') {
		const { protocol, hostname } = window.location;
		// Use the same hostname as the current page, but port 8080 for the API
		return `${protocol}//${hostname}:8080/api`;
	}
	// Fallback for SSR or non-browser environments
	return 'http://localhost:8080/api';
}

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
