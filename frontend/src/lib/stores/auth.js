import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// API base URL - can be configured via environment variables
const API_BASE = 'http://localhost:8081/api/v1';

// User store
export const user = writable(null);
export const isAuthenticated = writable(false);
export const authLoading = writable(false);

// Token management
let accessToken = null;
let refreshToken = null;

// Initialize auth state from localStorage
if (browser) {
	const storedUser = localStorage.getItem('hi_user');
	const storedAccessToken = localStorage.getItem('hi_access_token');
	const storedRefreshToken = localStorage.getItem('hi_refresh_token');

	if (storedUser && storedAccessToken && storedRefreshToken) {
		try {
			const userData = JSON.parse(storedUser);
			user.set(userData);
			accessToken = storedAccessToken;
			refreshToken = storedRefreshToken;
			isAuthenticated.set(true);
		} catch (error) {
			console.error('Failed to parse stored user data:', error);
			clearAuth();
		}
	}
}

// Clear authentication data
function clearAuth() {
	user.set(null);
	isAuthenticated.set(false);
	accessToken = null;
	refreshToken = null;

	if (browser) {
		localStorage.removeItem('hi_user');
		localStorage.removeItem('hi_access_token');
		localStorage.removeItem('hi_refresh_token');
	}
}

// Store authentication data
function storeAuth(authResponse) {
	const { user: userData, access_token, refresh_token } = authResponse;

	user.set(userData);
	isAuthenticated.set(true);
	accessToken = access_token;
	refreshToken = refresh_token;

	if (browser) {
		localStorage.setItem('hi_user', JSON.stringify(userData));
		localStorage.setItem('hi_access_token', access_token);
		localStorage.setItem('hi_refresh_token', refresh_token);
	}
}

// Make authenticated API requests
export async function apiRequest(endpoint, options = {}) {
	const url = `${API_BASE}${endpoint}`;
	const headers = {
		'Content-Type': 'application/json',
		...options.headers
	};

	// Add auth header if we have a token
	if (accessToken) {
		headers.Authorization = `Bearer ${accessToken}`;
	}

	try {
		const response = await fetch(url, {
			...options,
			headers
		});

		// If we get a 401, try to refresh the token
		if (response.status === 401 && refreshToken && endpoint !== '/auth/refresh') {
			const refreshed = await refreshAccessToken();
			if (refreshed) {
				// Retry the original request with new token
				headers.Authorization = `Bearer ${accessToken}`;
				return fetch(url, { ...options, headers });
			} else {
				// Refresh failed, clear auth and redirect to login
				clearAuth();
				if (browser) {
					window.location.href = '/login';
				}
				throw new Error('Authentication expired');
			}
		}

		return response;
	} catch (error) {
		console.error('API request failed:', error);
		throw error;
	}
}

// Refresh access token
async function refreshAccessToken() {
	if (!refreshToken) return false;

	try {
		const response = await fetch(`${API_BASE}/auth/refresh`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				refresh_token: refreshToken
			})
		});

		if (response.ok) {
			const authResponse = await response.json();
			storeAuth(authResponse);
			return true;
		} else {
			clearAuth();
			return false;
		}
	} catch (error) {
		console.error('Token refresh failed:', error);
		clearAuth();
		return false;
	}
}

// Register new user
export async function register(userData) {
	authLoading.set(true);

	try {
		const response = await fetch(`${API_BASE}/auth/register`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(userData)
		});

		const data = await response.json();

		if (response.ok) {
			storeAuth(data);
			return { success: true, data };
		} else {
			return {
				success: false,
				error: data.error || 'Registration failed'
			};
		}
	} catch (error) {
		console.error('Registration error:', error);
		return {
			success: false,
			error: 'Network error. Please try again.'
		};
	} finally {
		authLoading.set(false);
	}
}

// Login user
export async function login(credentials) {
	authLoading.set(true);

	try {
		const response = await fetch(`${API_BASE}/auth/login`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(credentials)
		});

		const data = await response.json();

		if (response.ok) {
			storeAuth(data);
			return { success: true, data };
		} else {
			return {
				success: false,
				error: data.error || 'Login failed'
			};
		}
	} catch (error) {
		console.error('Login error:', error);
		return {
			success: false,
			error: 'Network error. Please try again.'
		};
	} finally {
		authLoading.set(false);
	}
}

// Logout user
export function logout() {
	clearAuth();
	if (browser) {
		window.location.href = '/';
	}
}

// Get current auth token (for manual API calls)
export function getAccessToken() {
	return accessToken;
}

// Check if user is authenticated
export function checkAuth() {
	return accessToken !== null && refreshToken !== null;
}
