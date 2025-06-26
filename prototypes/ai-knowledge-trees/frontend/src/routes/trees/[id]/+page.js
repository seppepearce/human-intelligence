// Load tree data with nodes for the tree detail page
const API_BASE_URL = 'http://localhost:8081';

export async function load({ params, fetch }) {
    const treeId = params.id;

    try {
        // Fetch tree data with nodes using SvelteKit's universal fetch
        const response = await fetch(`${API_BASE_URL}/api/v1/trees/${treeId}`);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const result = await response.json();

        if (!result.success) {
            throw new Error(result.error || 'API request failed');
        }

        return {
            tree: result.data,
            treeId: treeId
        };
    } catch (error) {
        console.error('Error loading tree:', error);
        return {
            tree: null,
            treeId: treeId,
            error: error.message || 'Failed to load tree data'
        };
    }
}
