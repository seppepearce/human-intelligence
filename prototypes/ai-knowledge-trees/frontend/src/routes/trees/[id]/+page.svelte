<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import api from '$lib/api.js';
    import FlatTreeView from '$lib/components/FlatTreeView.svelte';
    import TreeVisualization from '$lib/components/TreeVisualization.svelte';

    export let data;

    let { tree, treeId, error } = data;
    let loading = false;

    // Node management
    let showNodeEditor = false;
    let editingNode = null;
    let showDeleteConfirm = false;
    let nodeToDelete = null;

    // New node form
    let newNode = {
        title: '',
        content: '',
        parent_id: null,
        position: 0
    };

    // Form binding variables
    let formTitle = '';
    let formContent = '';
    let formParentId = null;
    let formPosition = 0;

    // AI Analysis
    let analyzing = false;
    let nodeWithAnalysis = null;

    // View mode toggle
    let viewMode = 'tree'; // 'tree' or 'flat'

    onMount(() => {
        if (!tree && !error) {
            loadTree();
        }
    });

    async function loadTree() {
        loading = true;
        const result = await api.loadTree(treeId);
        if (result.success) {
            tree = result.tree;
            error = null;
        } else {
            error = result.error;
        }
        loading = false;
    }

    async function createNode() {
        if (!formTitle.trim()) {
            alert('Please enter a node title');
            return;
        }

        const nodeData = {
            title: formTitle,
            content: formContent,
            parent_id: formParentId,
            position: formPosition,
            tree_id: treeId
        };

        const result = await api.saveNode(nodeData, false);
        if (result.success) {
            // Reset form
            newNode = { title: '', content: '', parent_id: null, position: 0 };
            formTitle = '';
            formContent = '';
            formParentId = null;
            formPosition = 0;
            showNodeEditor = false;
            editingNode = null;
            // Reload tree
            await loadTree();
            alert('✅ Node created successfully!');
        } else {
            alert('❌ Failed to create node: ' + result.error);
        }
    }

    async function updateNode() {
        if (!formTitle.trim()) {
            alert('Please enter a node title');
            return;
        }

        const nodeData = {
            id: editingNode.id,
            title: formTitle,
            content: formContent,
            parent_id: formParentId,
            position: formPosition
        };

        const result = await api.saveNode(nodeData, true);
        if (result.success) {
            showNodeEditor = false;
            editingNode = null;
            await loadTree();
            alert('✅ Node updated successfully!');
        } else {
            alert('❌ Failed to update node: ' + result.error);
        }
    }

    async function deleteNode(node) {
        const result = await api.removeNode(node.id);
        if (result.success) {
            await loadTree();
            alert('✅ Node deleted successfully!');
        } else {
            alert('❌ Failed to delete node: ' + result.error);
        }

        nodeToDelete = null;
        showDeleteConfirm = false;
    }

    async function analyzeNode(node) {
        analyzing = true;
        nodeWithAnalysis = node;

        const result = await api.analyzeNodeContent(node.id);
        if (result.success) {
            // Update the node in the tree with analysis
            const nodeIndex = tree.nodes.findIndex(n => n.id === node.id);
            if (nodeIndex !== -1) {
                tree.nodes[nodeIndex].ai_analysis = result.analysis;
                tree = { ...tree }; // Trigger reactivity
            }
        } else {
            alert('❌ AI analysis failed: ' + result.error);
        }

        analyzing = false;
        nodeWithAnalysis = null;
    }

    function startCreateNode(parentId = null) {
        newNode = { title: '', content: '', parent_id: parentId, position: 0 };
        editingNode = null;
        formTitle = '';
        formContent = '';
        formParentId = parentId;
        formPosition = 0;
        showNodeEditor = true;
    }

    function startEditNode(node) {
        editingNode = { ...node };
        newNode = { title: '', content: '', parent_id: null, position: 0 };
        formTitle = node.title || '';
        formContent = node.content || '';
        formParentId = node.parent_id;
        formPosition = node.position || 0;
        showNodeEditor = true;
    }

    function confirmDeleteNode(node) {
        nodeToDelete = node;
        showDeleteConfirm = true;
    }

    function closeEditor() {
        showNodeEditor = false;
        editingNode = null;
        newNode = { title: '', content: '', parent_id: null, position: 0 };
        formTitle = '';
        formContent = '';
        formParentId = null;
        formPosition = 0;
    }

    function closeDeleteModal() {
        showDeleteConfirm = false;
        nodeToDelete = null;
    }

    function formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function getNodeChildren(parentId) {
        if (!tree?.nodes) return [];
        return tree.nodes.filter(node => node.parent_id === parentId)
                        .sort((a, b) => (a.position || 0) - (b.position || 0));
    }

    function getRootNodes() {
        return getNodeChildren(null);
    }

    // Get available parents for node creation/editing
    function getAvailableParents(excludeId = null) {
        if (!tree?.nodes) return [];
        return tree.nodes.filter(node => node.id !== excludeId);
    }
</script>

<svelte:head>
    <title>{tree?.title || 'Tree'} - AI Knowledge Trees</title>
</svelte:head>

<div class="container">
    {#if error}
        <div class="error-state">
            <div class="error-icon">❌</div>
            <h2>Error Loading Tree</h2>
            <p>{error}</p>
            <div class="error-actions">
                <button class="btn btn-primary" on:click={loadTree}>🔄 Retry</button>
                <button class="btn btn-secondary" on:click={() => goto('/trees')}>← Back to Trees</button>
            </div>
        </div>
    {:else if loading}
        <div class="loading-state">
            <div class="loading-spinner">🌳</div>
            <p>Loading tree...</p>
        </div>
    {:else if tree}
        <!-- Tree Header -->
        <header class="tree-header">
            <div class="header-content">
                <div class="header-main">
                    <button class="btn-back" on:click={() => goto('/trees')}>
                        ← All Trees
                    </button>
                    <h1>🌳 {tree.title}</h1>
                    {#if tree.description}
                        <p class="tree-description">{tree.description}</p>
                    {/if}
                </div>
                <div class="header-actions">
                    <div class="view-toggle">
                        <button
                            class="btn-toggle"
                            class:active={viewMode === 'tree'}
                            on:click={() => viewMode = 'tree'}
                            title="Tree Visualization"
                        >
                            🌳 Tree View
                        </button>
                        <button
                            class="btn-toggle"
                            class:active={viewMode === 'flat'}
                            on:click={() => viewMode = 'flat'}
                            title="Flat Hierarchy"
                        >
                            📁 List View
                        </button>
                    </div>
                    <button class="btn btn-primary" on:click={() => startCreateNode()}>
                        ➕ Add Root Node
                    </button>
                </div>
            </div>

            <div class="tree-meta">
                <div class="meta-item">
                    <span class="meta-label">Created:</span>
                    <span class="meta-value">{formatDate(tree.created_at)}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Nodes:</span>
                    <span class="meta-value">{tree.nodes?.length || 0}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Visibility:</span>
                    <span class="meta-value">{tree.is_public ? '🌍 Public' : '🔒 Private'}</span>
                </div>
                {#if tree.updated_at}
                    <div class="meta-item">
                        <span class="meta-label">Updated:</span>
                        <span class="meta-value">{formatDate(tree.updated_at)}</span>
                    </div>
                {/if}
            </div>
        </header>

        <!-- Tree Content -->
        <main class="tree-content">
            {#if !tree.nodes || tree.nodes.length === 0}
                <div class="empty-tree">
                    <div class="empty-icon">🌱</div>
                    <h3>Empty Tree</h3>
                    <p>This tree has no nodes yet. Start growing your knowledge by adding the first node.</p>
                    <button class="btn btn-primary" on:click={() => startCreateNode()}>
                        🌱 Plant First Node
                    </button>
                </div>
            {:else}
                {#if viewMode === 'tree'}
                    <div class="tree-visualization-container">
                        <TreeVisualization
                            nodes={tree.nodes}
                            width={1200}
                            height={700}
                            {analyzing}
                            {nodeWithAnalysis}
                            on:addChild={(e) => startCreateNode(e.detail.node.id)}
                            on:edit={(e) => startEditNode(e.detail.node)}
                            on:delete={(e) => confirmDeleteNode(e.detail.node)}
                            on:analyze={(e) => analyzeNode(e.detail.node)}
                        />
                    </div>
                {:else}
                    <FlatTreeView
                        nodes={tree.nodes}
                        {analyzing}
                        {nodeWithAnalysis}
                        on:addChild={(e) => startCreateNode(e.detail.node.id)}
                        on:edit={(e) => startEditNode(e.detail.node)}
                        on:delete={(e) => confirmDeleteNode(e.detail.node)}
                        on:analyze={(e) => analyzeNode(e.detail.node)}
                    />
                {/if}
            {/if}
        </main>
    {/if}
</div>

<!-- Node Editor Modal -->
{#if showNodeEditor}
    <div class="modal-overlay" on:click={closeEditor}>
        <div class="modal modal-large" on:click|stopPropagation>
            <div class="modal-header">
                <h3>{editingNode ? '✏️ Edit Node' : '➕ Create New Node'}</h3>
                <button class="btn-close" on:click={closeEditor}>✕</button>
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label for="node-title">Node Title *</label>
                    <input
                        id="tree-title"
                        type="text"
                        bind:value={formTitle}
                        placeholder="e.g., Consciousness, Algorithms, Personal Growth"
                        class="form-input"
                        required
                    />
                </div>

                <div class="form-group">
                    <label for="node-parent">Parent Node</label>
                    <select
                        id="node-parent"
                        bind:value={formParentId}
                        class="form-select"
                    >
                        <option value={null}>Root Level (No Parent)</option>
                        {#each getAvailableParents(editingNode?.id) as parentOption}
                            <option value={parentOption.id}>{parentOption.title}</option>
                        {/each}
                    </select>
                </div>

                <div class="form-group">
                    <label for="node-content">Content</label>
                    <textarea
                        id="node-content"
                        bind:value={formContent}
                        placeholder="Describe this concept, idea, or knowledge area..."
                        class="form-textarea"
                        rows="8"
                    ></textarea>
                </div>

                <div class="form-group">
                    <label for="node-position">Position</label>
                    <input
                        id="node-position"
                        type="number"
                        bind:value={formPosition}
                        min="0"
                        class="form-input"
                        placeholder="0"
                    />
                    <small class="form-help">Lower numbers appear first among siblings</small>
                </div>
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" on:click={closeEditor}>
                    Cancel
                </button>
                <button
                    class="btn btn-primary"
                    on:click={editingNode ? updateNode : createNode}
                >
                    {editingNode ? '💾 Update Node' : '🌱 Create Node'}
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm && nodeToDelete}
    <div class="modal-overlay" on:click={closeDeleteModal}>
        <div class="modal modal-small" on:click|stopPropagation>
            <div class="modal-header">
                <h3>🗑️ Delete Node</h3>
                <button class="btn-close" on:click={closeDeleteModal}>✕</button>
            </div>

            <div class="modal-body">
                <p>Are you sure you want to delete <strong>"{nodeToDelete.title}"</strong>?</p>
                <p class="warning-text">⚠️ This action cannot be undone. All child nodes will also be deleted.</p>
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" on:click={closeDeleteModal}>
                    Cancel
                </button>
                <button class="btn btn-danger" on:click={() => deleteNode(nodeToDelete)}>
                    🗑️ Delete Node
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

    .error-state, .loading-state {
        text-align: center;
        padding: 3rem 0;
    }

    .error-icon, .loading-spinner {
        font-size: 3rem;
        margin-bottom: 1rem;
    }

    .loading-spinner {
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }

    .error-actions {
        margin-top: 2rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
    }

    .tree-header {
        background: white;
        border-radius: 0.5rem;
        border: 1px solid #d4af37;
        padding: 2rem;
        margin-bottom: 2rem;
    }

    .header-content {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.5rem;
        gap: 2rem;
    }

    .header-main h1 {
        color: #212529;
        margin: 0.5rem 0;
        font-size: 2rem;
    }

    .tree-description {
        color: #6c757d;
        font-size: 1.1rem;
        margin: 0;
    }

    .btn-back {
        background: none;
        border: none;
        color: #d4af37;
        font-size: 1rem;
        cursor: pointer;
        margin-bottom: 0.5rem;
        padding: 0.25rem 0;
    }

    .btn-back:hover {
        color: #b8941f;
        text-decoration: underline;
    }

    .tree-meta {
        display: flex;
        gap: 2rem;
        padding-top: 1.5rem;
        border-top: 1px solid #e9ecef;
    }

    .meta-item {
        display: flex;
        gap: 0.5rem;
    }

    .meta-label {
        color: #6c757d;
        font-weight: 600;
    }

    .meta-value {
        color: #212529;
    }

    .empty-tree {
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

    .tree-content {
        margin-top: 2rem;
    }

    .view-toggle {
        display: flex;
        gap: 0.25rem;
        margin-right: 1rem;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        overflow: hidden;
    }

    .btn-toggle {
        background: white;
        border: none;
        padding: 0.5rem 1rem;
        cursor: pointer;
        font-size: 0.9rem;
        font-family: Georgia, serif;
        transition: all 0.2s;
        color: #212529;
    }

    .btn-toggle:hover {
        background: #f8f9fa;
    }

    .btn-toggle.active {
        background: #d4af37;
        color: white;
    }

    .tree-visualization-container {
        background: white;
        border-radius: 0.5rem;
        border: 1px solid #e9ecef;
        overflow: hidden;
        height: 700px;
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

    .modal-large {
        max-width: 600px;
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

    .form-input, .form-textarea, .form-select {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        font-size: 1rem;
        font-family: Georgia, serif;
    }

    .form-input:focus, .form-textarea:focus, .form-select:focus {
        outline: none;
        border-color: #b8941f;
        box-shadow: 0 0 0 2px rgba(212, 175, 55, 0.2);
    }

    .form-help {
        display: block;
        margin-top: 0.25rem;
        color: #6c757d;
        font-size: 0.875rem;
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
        .header-content {
            flex-direction: column;
            gap: 1rem;
        }

        .tree-meta {
            flex-direction: column;
            gap: 0.5rem;
        }

        .meta-item {
            justify-content: space-between;
        }

        .node-header {
            flex-direction: column;
            gap: 1rem;
        }

        .node-actions {
            margin-left: 0;
            justify-content: flex-end;
        }

        .header-actions {
            flex-direction: column;
            gap: 1rem;
        }

        .view-toggle {
            margin-right: 0;
        }
    }
</style>
