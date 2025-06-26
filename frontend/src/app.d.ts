// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				id: string;
				username: string;
				email: string;
				avatar?: string;
				isActive: boolean;
				createdAt: string;
			};
			session?: {
				id: string;
				userId: string;
				expiresAt: string;
			};
		}
		interface PageData {
			user?: App.Locals['user'];
		}
		interface PageState {
			user?: App.Locals['user'];
		}
		// interface Platform {}
	}

	// Global type definitions
	interface Window {
		gtag?: (...args: any[]) => void;
		dataLayer?: any[];
	}

	// Environment variables
	declare namespace NodeJS {
		interface ProcessEnv {
			VITE_API_URL: string;
			VITE_WS_URL: string;
			VITE_ENABLE_AI_FEATURES: string;
			VITE_ENABLE_REAL_TIME: string;
		}
	}
}

// API Response Types
export interface APIResponse<T = any> {
	data?: T;
	error?: string;
	message?: string;
	success: boolean;
}

// User Types
export interface User {
	id: string;
	username: string;
	email: string;
	bio?: string;
	avatar?: string;
	createdAt: string;
	updatedAt: string;
	isActive: boolean;
}

export interface AuthResponse {
	user: User;
	accessToken: string;
	refreshToken: string;
	expiresAt: string;
}

// Node Types
export type NodeType = 'text' | 'video' | 'link' | 'quiz' | 'code' | 'latex' | 'jupyter' | 'wolfram' | 'custom';

export interface Node {
	id: string;
	title: string;
	description: string;
	content: string;
	nodeType: NodeType;
	metadata: Record<string, any>;
	tags: string[];
	createdBy: string;
	createdAt: string;
	updatedAt: string;
	isPublic: boolean;
	voteScore: number;
	viewCount: number;
}

export interface NodeMetadata {
	// Video nodes
	videoUrl?: string;
	videoDuration?: number;

	// Link nodes
	linkUrl?: string;
	linkTitle?: string;
	linkDomain?: string;

	// Quiz nodes
	questions?: QuizQuestion[];

	// Code nodes
	language?: string;
	repository?: string;

	// Custom plugin nodes
	pluginId?: string;
	pluginData?: Record<string, any>;
}

export interface QuizQuestion {
	id: string;
	question: string;
	options: string[];
	correct: number;
	explanation?: string;
}

// Learning Path Types
export type AccessLevel = 'public' | 'private' | 'classroom' | 'group';

export interface LearningPath {
	id: string;
	name: string;
	description: string;
	accessLevel: AccessLevel;
	createdBy: string;
	parentPathId?: string;
	createdAt: string;
	updatedAt: string;
	voteScore: number;
	forkCount: number;
	completionCount: number;
	tags: string[];
	nodes?: PathNode[];
	creator?: User;
	parentPath?: LearningPath;
}

export interface PathNode {
	id: string;
	pathId: string;
	nodeId: string;
	position: number;
	isRequired: boolean;
	prerequisites: string[];
	node?: Node;
}

// TLDR Types
export interface TLDR {
	id: string;
	pathId: string;
	userId: string;
	content: string;
	createdAt: string;
	voteScore: number;
	user?: User;
	learningPath?: LearningPath;
}

// Activity Types
export interface ActivityEvent {
	id: string;
	type: string;
	userId: string;
	targetId: string;
	pathId?: string;
	nodeId?: string;
	timestamp: string;
	user?: User;
	learningPath?: LearningPath;
	node?: Node;
}

// Search Types
export interface SearchQuery {
	query: string;
	tags?: string[];
	nodeTypes?: NodeType[];
	accessLevel?: AccessLevel;
	minRating?: number;
	dateFrom?: string;
	dateTo?: string;
	sortBy?: string;
	limit?: number;
	offset?: number;
}

export interface SearchResult {
	nodes: Node[];
	paths: LearningPath[];
	total: number;
	query: string;
	took: number;
}

// WebSocket Types
export interface WebSocketMessage {
	type: string;
	data?: any;
	timestamp?: string;
}

// Component Props Types
export interface ToastOptions {
	type?: 'success' | 'error' | 'warning' | 'info';
	duration?: number;
	dismissible?: boolean;
}

// Form Types
export interface LoginForm {
	email: string;
	password: string;
}

export interface RegisterForm {
	username: string;
	email: string;
	password: string;
	bio?: string;
}

export interface CreateNodeForm {
	title: string;
	description: string;
	content: string;
	nodeType: NodeType;
	metadata: Record<string, any>;
	tags: string[];
	isPublic: boolean;
}

export interface CreatePathForm {
	name: string;
	description: string;
	accessLevel: AccessLevel;
	tags: string[];
	nodeIds: string[];
}

// Plugin System Types
export interface Plugin {
	id: string;
	name: string;
	description: string;
	version: string;
	author: string;
	configSchema: Record<string, any>;
	isActive: boolean;
	isPremium: boolean;
	createdAt: string;
}

// Module declarations for static assets
declare module '*.svg' {
	const content: string;
	export default content;
}

declare module '*.png' {
	const content: string;
	export default content;
}

declare module '*.jpg' {
	const content: string;
	export default content;
}

declare module '*.jpeg' {
	const content: string;
	export default content;
}

declare module '*.gif' {
	const content: string;
	export default content;
}

declare module '*.webp' {
	const content: string;
	export default content;
}

declare module '*.ico' {
	const content: string;
	export default content;
}

export {};
