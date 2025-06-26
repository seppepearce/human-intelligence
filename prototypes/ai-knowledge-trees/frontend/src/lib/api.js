// API client utility for AI Knowledge Trees backend
import { testApi, isTestMode } from './testData.js';

const API_BASE_URL = 'http://localhost:8081';

class APIError extends Error {
    constructor(message, status, data) {
        super(message);
        this.name = 'APIError';
        this.status = status;
        this.data = data;
    }
}

class APIClient {
    constructor(baseUrl = API_BASE_URL) {
        this.baseUrl = baseUrl;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseUrl}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new APIError(
                    data.error || `HTTP ${response.status}`,
                    response.status,
                    data
                );
            }

            if (!data.success) {
                throw new APIError(
                    data.error || 'API request failed',
                    response.status,
                    data
                );
            }

            return data;
        } catch (error) {
            if (error instanceof APIError) {
                throw error;
            }

            // Network or other errors
            throw new APIError(
                error.message || 'Network error',
                0,
                null
            );
        }
    }

    // Health check
    async health() {
        if (isTestMode()) {
            return testApi.health();
        }
        return this.request('/api/v1/health');
    }

    // Tree operations
    async getTrees(page = 1, limit = 20) {
        if (isTestMode()) {
            return testApi.getTrees();
        }
        const params = new URLSearchParams({ page: page.toString(), limit: limit.toString() });
        return this.request(`/api/v1/trees?${params}`);
    }

    async getTree(id) {
        if (isTestMode()) {
            return testApi.getTree(id);
        }
        return this.request(`/api/v1/trees/${id}`);
    }

    async createTree(treeData) {
        if (isTestMode()) {
            return testApi.createTree(treeData);
        }
        return this.request('/api/v1/trees', {
            method: 'POST',
            body: JSON.stringify(treeData)
        });
    }

    async updateTree(id, treeData) {
        return this.request(`/api/v1/trees/${id}`, {
            method: 'PUT',
            body: JSON.stringify(treeData)
        });
    }

    async deleteTree(id) {
        if (isTestMode()) {
            return testApi.deleteTree(id);
        }
        return this.request(`/api/v1/trees/${id}`, {
            method: 'DELETE'
        });
    }

    // Node operations
    async getNode(id) {
        return this.request(`/api/v1/nodes/${id}`);
    }

    async createNode(nodeData) {
        if (isTestMode()) {
            return testApi.createNode(nodeData);
        }
        return this.request('/api/v1/nodes', {
            method: 'POST',
            body: JSON.stringify(nodeData)
        });
    }

    async updateNode(id, nodeData) {
        if (isTestMode()) {
            return testApi.updateNode(id, nodeData);
        }
        return this.request(`/api/v1/nodes/${id}`, {
            method: 'PUT',
            body: JSON.stringify(nodeData)
        });
    }

    async deleteNode(id) {
        if (isTestMode()) {
            return testApi.deleteNode(id);
        }
        return this.request(`/api/v1/nodes/${id}`, {
            method: 'DELETE'
        });
    }

    async analyzeNode(id) {
        if (isTestMode()) {
            return testApi.analyzeNode(id);
        }
        return this.request(`/api/v1/nodes/${id}/analyze`, {
            method: 'POST'
        });
    }

    // AI operations
    async analyzeContent(title, content) {
        return this.request('/api/v1/ai/analyze', {
            method: 'POST',
            body: JSON.stringify({ title, content })
        });
    }

    async suggestConnections(nodeId, limit = 5) {
        return this.request('/api/v1/ai/suggest-connections', {
            method: 'POST',
            body: JSON.stringify({ node_id: nodeId, limit })
        });
    }

    async getAIHealth() {
        return this.request('/api/v1/ai/health');
    }

    // Batch operations
    async getTreesWithStats() {
        return this.request('/api/v1/trees/stats');
    }

    async searchTrees(query) {
        const params = new URLSearchParams({ q: query });
        return this.request(`/api/v1/trees/search?${params}`);
    }

    async searchNodes(query, treeId = null) {
        const params = new URLSearchParams({ q: query });
        if (treeId) {
            params.append('tree_id', treeId);
        }
        return this.request(`/api/v1/nodes/search?${params}`);
    }
}

// Create singleton instance
export const apiClient = new APIClient();

// Export the class for custom instances
export { APIClient, APIError };

// Helper functions for common operations
export const api = {
    // Tree helpers
    async loadTrees() {
        try {
            const result = await apiClient.getTrees();
            return { success: true, trees: result.data || [], meta: result.meta };
        } catch (error) {
            console.error('Error loading trees:', error);
            return { success: false, error: error.message, trees: [] };
        }
    },

    async loadTree(id) {
        try {
            const result = await apiClient.getTree(id);
            return { success: true, tree: result.data };
        } catch (error) {
            console.error('Error loading tree:', error);
            return { success: false, error: error.message, tree: null };
        }
    },

    async saveTree(treeData, isEditing = false) {
        try {
            const result = isEditing
                ? await apiClient.updateTree(treeData.id, treeData)
                : await apiClient.createTree(treeData);
            return { success: true, tree: result.data };
        } catch (error) {
            console.error('Error saving tree:', error);
            return { success: false, error: error.message };
        }
    },

    async removeTree(id) {
        try {
            await apiClient.deleteTree(id);
            return { success: true };
        } catch (error) {
            console.error('Error deleting tree:', error);
            return { success: false, error: error.message };
        }
    },

    // Node helpers
    async loadNode(id) {
        try {
            const result = await apiClient.getNode(id);
            return { success: true, node: result.data };
        } catch (error) {
            console.error('Error loading node:', error);
            return { success: false, error: error.message, node: null };
        }
    },

    async saveNode(nodeData, isEditing = false) {
        try {
            const result = isEditing
                ? await apiClient.updateNode(nodeData.id, nodeData)
                : await apiClient.createNode(nodeData);
            return { success: true, node: result.data };
        } catch (error) {
            console.error('Error saving node:', error);
            return { success: false, error: error.message };
        }
    },

    async removeNode(id) {
        try {
            await apiClient.deleteNode(id);
            return { success: true };
        } catch (error) {
            console.error('Error deleting node:', error);
            return { success: false, error: error.message };
        }
    },

    async analyzeNodeContent(id) {
        try {
            const result = await apiClient.analyzeNode(id);
            return { success: true, analysis: result.data };
        } catch (error) {
            console.error('Error analyzing node:', error);
            return { success: false, error: error.message, analysis: null };
        }
    },

    // AI helpers
    async analyzeText(title, content) {
        try {
            const result = await apiClient.analyzeContent(title, content);
            return { success: true, analysis: result.data };
        } catch (error) {
            console.error('Error analyzing content:', error);
            return { success: false, error: error.message, analysis: null };
        }
    },

    async getConnectionSuggestions(nodeId, limit = 5) {
        try {
            const result = await apiClient.suggestConnections(nodeId, limit);
            return { success: true, suggestions: result.data };
        } catch (error) {
            console.error('Error getting connection suggestions:', error);
            return { success: false, error: error.message, suggestions: [] };
        }
    },

    // System helpers
    async checkHealth() {
        try {
            const result = await apiClient.health();
            return { success: true, health: result.data };
        } catch (error) {
            console.error('Error checking health:', error);
            return { success: false, error: error.message, health: null };
        }
    },

    async checkAIHealth() {
        try {
            const result = await apiClient.getAIHealth();
            return { success: true, aiHealth: result.data };
        } catch (error) {
            console.error('Error checking AI health:', error);
            return { success: false, error: error.message, aiHealth: null };
        }
    }
};

// Export test mode controls
export { enableTestMode, disableTestMode, isTestMode } from './testData.js';

// Export default for convenience
export default api;
