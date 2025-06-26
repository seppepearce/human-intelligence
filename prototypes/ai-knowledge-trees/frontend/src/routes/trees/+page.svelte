<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import api, { isTestMode } from '$lib/api.js';
    import TreeCard from '$lib/components/TreeCard.svelte';

    let trees = [];
    let loading = false;
    let searchTerm = '';
    let showCreateModal = false;

    // Create tree form data
    let newTree = {
        title: '',
        description: '',
        is_public: true
    };

    // Delete confirmation
    let treeToDelete = null;
    let showDeleteConfirm = false;

    onMount(() => {
        loadTrees();
    });

    async function loadTrees() {
        loading = true;
        const result = await api.loadTrees();
        if (result.success) {
            trees = result.trees;
        } else {
            console.error('Failed to load trees:', result.error);
            alert('❌ Cannot connect to backend. Make sure it\'s running on port 8081');
            trees = [];
        }
        loading = false;
    }

    async function createTree() {
        if (!newTree.title.trim()) {
            alert('Please enter a tree title');
            return;
        }

        const result = await api.saveTree(newTree, false);
        if (result.success) {
            // Reset form
            newTree = { title: '', description: '', is_public: true };
            showCreateModal = false;
            // Reload trees
            await loadTrees();
            alert('✅ Tree created successfully!');
        } else {
            alert('❌ Failed to create tree: ' + result.error);
        }
    }

    async function deleteTree(tree) {
        const result = await api.removeTree(tree.id);
        if (result.success) {
            await loadTrees();
            alert('✅ Tree deleted successfully!');
        } else {
            alert('❌ Failed to delete tree: ' + result.error);
        }

        treeToDelete = null;
        showDeleteConfirm = false;
    }

    function confirmDelete(tree) {
        treeToDelete = tree;
        showDeleteConfirm = true;
    }

    function handleTreeView(event) {
        const tree = event.detail.tree;
        goto(`/trees/${tree.id}`);
    }

    function handleTreeDelete(event) {
        const tree = event.detail.tree;
        confirmDelete(tree);
    }

    function closeCreateModal() {
        showCreateModal = false;
        newTree = { title: '', description: '', is_public: true };
    }

    function closeDeleteModal() {
        showDeleteConfirm = false;
        treeToDelete = null;
    }

    // Filter trees based on search term
    $: filteredTrees = trees.filter(tree =>
        tree.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (tree.description && tree.description.toLowerCase().includes(searchTerm.toLowerCase()))
    );

    function formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }
</script>

<svelte:head>
    <title>Knowledge Trees - AI Knowledge Trees</title>
</svelte:head>

<div class="container">
    <header class="page-header">
        <h1>🌳 Knowledge Trees</h1>
        <p>Grow your digital garden of wisdom</p>
        {#if isTestMode()}
            <div class="test-mode-badge">
                🧪 Test Mode - Using Mock Data
            </div>
        {/if}
    </header>

    <div class="actions-bar">
        <div class="search-section">
            <input
                type="text"
                placeholder="🔍 Search trees..."
                bind:value={searchTerm}
                class="search-input"
            />
        </div>
        <div class="button-section">
            <button class="btn btn-primary" on:click={() => showCreateModal = true}>
                ➕ Create New Tree
            </button>
            <button class="btn btn-secondary" on:click={loadTrees} disabled={loading}>
                {loading ? '🔄 Loading...' : '🔄 Refresh'}
            </button>
        </div>
    </div>

    {#if loading}
        <div class="loading-state">
            <div class="loading-spinner">🌳</div>
            <p>Loading your knowledge trees...</p>
        </div>
    {:else if filteredTrees.length === 0}
        <div class="empty-state">
            {#if trees.length === 0}
                <div class="empty-icon">🌱</div>
                <h3>No Trees Yet</h3>
                <p>Start your journey of knowledge by creating your first tree.</p>
                <button class="btn btn-primary" on:click={() => showCreateModal = true}>
                    🌳 Plant Your First Tree
                </button>
            {:else}
                <div class="empty-icon">🔍</div>
                <h3>No Trees Match Your Search</h3>
                <p>Try a different search term or create a new tree.</p>
            {/if}
        </div>
    {:else}
        <div class="trees-grid">
            {#each filteredTrees as tree (tree.id)}
                <TreeCard
                    {tree}
                    on:view={handleTreeView}
                    on:delete={handleTreeDelete}
                />
            {/each}
        </div>
    {/if}
</div>

<!-- Create Tree Modal -->
{#if showCreateModal}
    <div class="modal-overlay" on:click={closeCreateModal}>
        <div class="modal" on:click|stopPropagation>
            <div class="modal-header">
                <h3>🌱 Create New Knowledge Tree</h3>
                <button class="btn-close" on:click={closeCreateModal}>✕</button>
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label for="tree-title">Tree Title *</label>
                    <input
                        id="tree-title"
                        type="text"
                        bind:value={newTree.title}
                        placeholder="e.g., Philosophy, Computer Science, Personal Development"
                        class="form-input"
                        required
                    />
                </div>

                <div class="form-group">
                    <label for="tree-description">Description</label>
                    <textarea
                        id="tree-description"
                        bind:value={newTree.description}
                        placeholder="Describe what this knowledge tree will contain..."
                        class="form-textarea"
                        rows="3"
                    ></textarea>
                </div>

                <div class="form-group">
                    <label class="checkbox-label">
                        <input
                            type="checkbox"
                            bind:checked={newTree.is_public}
                        />
                        <span class="checkmark">✓</span>
                        Make this tree public
                    </label>
                </div>
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" on:click={closeCreateModal}>
                    Cancel
                </button>
                <button class="btn btn-primary" on:click={createTree}>
                    🌳 Create Tree
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm && treeToDelete}
    <div class="modal-overlay" on:click={closeDeleteModal}>
        <div class="modal modal-small" on:click|stopPropagation>
            <div class="modal-header">
                <h3>🗑️ Delete Tree</h3>
                <button class="btn-close" on:click={closeDeleteModal}>✕</button>
            </div>

            <div class="modal-body">
                <p>Are you sure you want to delete <strong>"{treeToDelete.title}"</strong>?</p>
                <p class="warning-text">⚠️ This action cannot be undone. All nodes in this tree will be permanently deleted.</p>
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" on:click={closeDeleteModal}>
                    Cancel
                </button>
                <button class="btn btn-danger" on:click={() => deleteTree(treeToDelete)}>
                    🗑️ Delete Tree
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 1rem;
    }

    .page-header {
        text-align: center;
        padding: 2rem 0;
        margin-bottom: 2rem;
    }

    .page-header h1 {
        color: #212529;
        margin-bottom: 0.5rem;
        font-size: 2.5rem;
    }

    .page-header p {
        color: #6c757d;
        font-size: 1.1rem;
    }

    .test-mode-badge {
        display: inline-block;
        background: #fff3cd;
        color: #856404;
        border: 1px solid #ffeaa7;
        border-radius: 0.25rem;
        padding: 0.5rem 1rem;
        font-size: 0.9rem;
        font-weight: 600;
        margin-top: 1rem;
    }

    .actions-bar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
        gap: 1rem;
        flex-wrap: wrap;
    }

    .search-section {
        flex: 1;
        min-width: 250px;
    }

    .search-input {
        width: 100%;
        padding: 0.75rem 1rem;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        font-size: 1rem;
        font-family: Georgia, serif;
    }

    .search-input:focus {
        outline: none;
        border-color: #b8941f;
        box-shadow: 0 0 0 2px rgba(212, 175, 55, 0.2);
    }

    .button-section {
        display: flex;
        gap: 0.5rem;
    }

    .loading-state {
        text-align: center;
        padding: 3rem 0;
    }

    .loading-spinner {
        font-size: 3rem;
        margin-bottom: 1rem;
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }

    .empty-state {
        text-align: center;
        padding: 3rem 0;
        background: white;
        border-radius: 0.5rem;
        border: 1px solid #e9ecef;
    }

    .empty-icon {
        font-size: 4rem;
        margin-bottom: 1rem;
    }

    .empty-state h3 {
        color: #212529;
        margin-bottom: 1rem;
    }

    .empty-state p {
        color: #6c757d;
        margin-bottom: 2rem;
    }

    .trees-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    /* Modal Styles */
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .modal {
        background: white;
        border-radius: 0.5rem;
        width: 100%;
        max-width: 500px;
        max-height: 90vh;
        overflow-y: auto;
        margin: 1rem;
        border: 1px solid #d4af37;
    }

    .modal-small {
        max-width: 400px;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 1.5rem 0;
        border-bottom: 1px solid #e9ecef;
        margin-bottom: 1.5rem;
    }

    .modal-header h3 {
        margin: 0;
        color: #212529;
    }

    .btn-close {
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        color: #6c757d;
        padding: 0;
        width: 2rem;
        height: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 0.25rem;
    }

    .btn-close:hover {
        background: #f8f9fa;
        color: #212529;
    }

    .modal-body {
        padding: 0 1.5rem;
    }

    .modal-footer {
        padding: 1.5rem;
        border-top: 1px solid #e9ecef;
        margin-top: 1.5rem;
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        color: #212529;
        font-weight: 600;
    }

    .form-input, .form-textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        font-size: 1rem;
        font-family: Georgia, serif;
    }

    .form-input:focus, .form-textarea:focus {
        outline: none;
        border-color: #b8941f;
        box-shadow: 0 0 0 2px rgba(212, 175, 55, 0.2);
    }

    .checkbox-label {
        display: flex !important;
        align-items: center;
        cursor: pointer;
    }

    .checkbox-label input[type="checkbox"] {
        margin-right: 0.5rem;
    }

    .warning-text {
        color: #dc3545;
        font-weight: 600;
    }

    .btn-danger {
        background: #dc3545;
        color: white;
    }

    .btn-danger:hover {
        background: #c82333;
    }

    /* Responsive design */
    @media (max-width: 768px) {
        .actions-bar {
            flex-direction: column;
            align-items: stretch;
        }

        .button-section {
            justify-content: center;
        }

        .trees-grid {
            grid-template-columns: 1fr;
        }


    }
</style>
