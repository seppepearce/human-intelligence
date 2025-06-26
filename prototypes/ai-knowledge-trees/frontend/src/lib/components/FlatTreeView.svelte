<script>
    import { createEventDispatcher } from 'svelte';

    export let nodes = [];
    export let showActions = true;
    export let maxDepth = 10;
    export let analyzing = false;
    export let nodeWithAnalysis = null;

    const dispatch = createEventDispatcher();

    // State management for expanded/collapsed nodes
    let expandedNodes = new Set();
    let searchTerm = '';

    // Build node hierarchy map for efficient lookups
    $: nodeMap = new Map(nodes.map(node => [node.id, node]));
    $: childrenMap = buildChildrenMap(nodes);
    $: visibleNodes = getVisibleNodes(nodes, expandedNodes, searchTerm);

    function buildChildrenMap(nodes) {
        const map = new Map();
        nodes.forEach(node => {
            const parentId = node.parent_id || 'root';
            if (!map.has(parentId)) {
                map.set(parentId, []);
            }
            map.get(parentId).push(node);
        });

        // Sort children by position
        map.forEach(children => {
            children.sort((a, b) => (a.position || 0) - (b.position || 0));
        });

        return map;
    }

    function hasChildren(nodeId) {
        return childrenMap.has(nodeId) && childrenMap.get(nodeId).length > 0;
    }

    function isExpanded(nodeId) {
        return expandedNodes.has(nodeId);
    }

    function toggleExpanded(nodeId) {
        const newExpanded = new Set(expandedNodes);
        if (newExpanded.has(nodeId)) {
            newExpanded.delete(nodeId);
        } else {
            newExpanded.add(nodeId);
        }
        expandedNodes = newExpanded;
    }

    function getVisibleNodes(nodes, expanded, search) {
        if (!nodes.length) return [];

        const filtered = search
            ? nodes.filter(node =>
                node.title.toLowerCase().includes(search.toLowerCase()) ||
                (node.content && node.content.toLowerCase().includes(search.toLowerCase()))
              )
            : nodes;

        if (search) {
            // When searching, show all matching nodes
            return filtered.map(node => ({
                ...node,
                _visible: true,
                _hasChildren: hasChildren(node.id),
                _isExpanded: isExpanded(node.id)
            }));
        }

        // Build visible list based on expansion state
        const visible = [];
        const rootNodes = filtered.filter(node => !node.parent_id);

        function addNodeAndChildren(node, isVisible = true) {
            if (!isVisible) return;

            const nodeWithMeta = {
                ...node,
                _visible: true,
                _hasChildren: hasChildren(node.id),
                _isExpanded: isExpanded(node.id)
            };

            visible.push(nodeWithMeta);

            // Add children if expanded
            if (isExpanded(node.id) && hasChildren(node.id)) {
                const children = childrenMap.get(node.id) || [];
                children.forEach(child => addNodeAndChildren(child, true));
            }
        }

        rootNodes.forEach(node => addNodeAndChildren(node));
        return visible;
    }

    function expandAll() {
        expandedNodes = new Set(nodes.map(node => node.id));
    }

    function collapseAll() {
        expandedNodes = new Set();
    }

    function getIndentationLevel(node) {
        return Math.min(node.depth || 0, maxDepth);
    }

    function getDifficultyClass(difficulty) {
        if (!difficulty) return '';
        if (difficulty <= 2) return 'difficulty-easy';
        if (difficulty <= 5) return 'difficulty-medium';
        if (difficulty <= 7) return 'difficulty-hard';
        return 'difficulty-expert';
    }

    function formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // Event handlers
    function handleAddChild(node) {
        dispatch('addChild', { node });
    }

    function handleEdit(node) {
        dispatch('edit', { node });
    }

    function handleDelete(node) {
        dispatch('delete', { node });
    }

    function handleAnalyze(node) {
        dispatch('analyze', { node });
    }

    function clearSearch() {
        searchTerm = '';
    }
</script>

<div class="flat-tree-view">
    <!-- Search and Controls -->
    <div class="tree-controls">
        <div class="search-section">
            <div class="search-wrapper">
                <input
                    type="text"
                    placeholder="🔍 Search nodes..."
                    bind:value={searchTerm}
                    class="search-input"
                />
                {#if searchTerm}
                    <button class="clear-search" on:click={clearSearch}>✕</button>
                {/if}
            </div>
        </div>

        <div class="expand-controls">
            <button class="btn-control" on:click={expandAll} title="Expand All">
                📂 Expand All
            </button>
            <button class="btn-control" on:click={collapseAll} title="Collapse All">
                📁 Collapse All
            </button>
        </div>
    </div>

    <!-- Tree Display -->
    {#if visibleNodes.length === 0}
        <div class="empty-tree">
            {#if searchTerm}
                <div class="empty-icon">🔍</div>
                <h3>No matching nodes</h3>
                <p>Try a different search term or <button class="link-button" on:click={clearSearch}>clear search</button></p>
            {:else}
                <div class="empty-icon">🌱</div>
                <h3>No nodes yet</h3>
                <p>Start building your knowledge tree by adding the first node.</p>
            {/if}
        </div>
    {:else}
        <div class="tree-list">
            {#each visibleNodes as node (node.id)}
                <div
                    class="tree-node"
                    class:has-children={node._hasChildren}
                    class:expanded={node._isExpanded}
                    style="--indent: {getIndentationLevel(node)}"
                >
                    <!-- Indentation and Expansion Controls -->
                    <div class="node-indent">
                        {#each Array(getIndentationLevel(node)) as _, i}
                            <span class="indent-line" class:last={i === getIndentationLevel(node) - 1}></span>
                        {/each}

                        {#if node._hasChildren}
                            <button
                                class="expand-toggle"
                                on:click={() => toggleExpanded(node.id)}
                                title={node._isExpanded ? 'Collapse' : 'Expand'}
                            >
                                {node._isExpanded ? '📂' : '📁'}
                            </button>
                        {:else}
                            <span class="node-bullet">•</span>
                        {/if}
                    </div>

                    <!-- Node Content -->
                    <div class="node-content">
                        <div class="node-header">
                            <h4 class="node-title">{node.title}</h4>

                            <div class="node-meta-inline">
                                {#if node.ai_analysis?.difficulty}
                                    <span class="difficulty-badge {getDifficultyClass(node.ai_analysis.difficulty)}">
                                        {node.ai_analysis.difficulty}/10
                                    </span>
                                {/if}

                                {#if node._hasChildren}
                                    <span class="children-count">
                                        {childrenMap.get(node.id)?.length || 0} children
                                    </span>
                                {/if}
                            </div>

                            {#if showActions}
                                <div class="node-actions">
                                    <button
                                        class="btn-icon"
                                        on:click={() => handleAddChild(node)}
                                        title="Add Child Node"
                                    >
                                        ➕
                                    </button>
                                    <button
                                        class="btn-icon"
                                        on:click={() => handleEdit(node)}
                                        title="Edit Node"
                                    >
                                        ✏️
                                    </button>
                                    <button
                                        class="btn-icon ai-btn"
                                        on:click={() => handleAnalyze(node)}
                                        title="AI Analysis"
                                        disabled={analyzing && nodeWithAnalysis?.id === node.id}
                                    >
                                        {analyzing && nodeWithAnalysis?.id === node.id ? '🔄' : '🧠'}
                                    </button>
                                    <button
                                        class="btn-icon delete"
                                        on:click={() => handleDelete(node)}
                                        title="Delete Node"
                                    >
                                        🗑️
                                    </button>
                                </div>
                            {/if}
                        </div>

                        {#if node.content}
                            <p class="node-description">{node.content}</p>
                        {/if}

                        <!-- AI Analysis Summary -->
                        {#if node.ai_analysis}
                            <div class="ai-summary">
                                <span class="ai-label">🧠</span>
                                {#if node.ai_analysis.concepts && node.ai_analysis.concepts.length > 0}
                                    <div class="concepts-inline">
                                        {#each node.ai_analysis.concepts.slice(0, 3) as concept}
                                            <span class="concept-tag">{concept}</span>
                                        {/each}
                                        {#if node.ai_analysis.concepts.length > 3}
                                            <span class="more-concepts">+{node.ai_analysis.concepts.length - 3} more</span>
                                        {/if}
                                    </div>
                                {/if}
                            </div>
                        {/if}

                        <div class="node-footer">
                            <span class="node-timestamp">Created {formatDate(node.created_at)}</span>
                            <span class="node-depth">Depth: {node.depth || 0}</span>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <!-- Summary Stats -->
        <div class="tree-stats">
            <span>Showing {visibleNodes.length} of {nodes.length} nodes</span>
            {#if searchTerm}
                <span>• Filtered by: "{searchTerm}"</span>
            {/if}
        </div>
    {/if}
</div>

<style>
    .flat-tree-view {
        background: white;
        border-radius: 0.5rem;
        border: 1px solid #e9ecef;
        overflow: hidden;
    }

    .tree-controls {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        border-bottom: 1px solid #e9ecef;
        background: #f8f9fa;
        gap: 1rem;
        flex-wrap: wrap;
    }

    .search-section {
        flex: 1;
        min-width: 250px;
    }

    .search-wrapper {
        position: relative;
    }

    .search-input {
        width: 100%;
        padding: 0.5rem 2rem 0.5rem 0.75rem;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        font-size: 0.9rem;
    }

    .clear-search {
        position: absolute;
        right: 0.5rem;
        top: 50%;
        transform: translateY(-50%);
        background: none;
        border: none;
        cursor: pointer;
        color: #6c757d;
        font-size: 0.8rem;
    }

    .expand-controls {
        display: flex;
        gap: 0.5rem;
    }

    .btn-control {
        padding: 0.5rem 0.75rem;
        border: 1px solid #d4af37;
        background: white;
        border-radius: 0.25rem;
        cursor: pointer;
        font-size: 0.8rem;
        transition: all 0.2s;
    }

    .btn-control:hover {
        background: #d4af37;
        color: white;
    }

    .tree-list {
        max-height: 600px;
        overflow-y: auto;
    }

    .tree-node {
        display: flex;
        align-items: flex-start;
        padding: 0.75rem 0;
        border-bottom: 1px solid #f1f3f4;
        transition: background-color 0.2s;
    }

    .tree-node:hover {
        background: #f8f9fa;
    }

    .tree-node.has-children {
        border-left: 2px solid #d4af37;
    }

    .node-indent {
        display: flex;
        align-items: center;
        padding: 0 0.5rem;
        min-width: calc(var(--indent) * 1.5rem + 2rem);
    }

    .indent-line {
        width: 1.5rem;
        height: 1px;
        background: #e9ecef;
        margin-right: 0.25rem;
        position: relative;
    }

    .indent-line.last::after {
        content: '';
        position: absolute;
        right: -0.125rem;
        top: -0.5rem;
        width: 1px;
        height: 1rem;
        background: #e9ecef;
    }

    .expand-toggle {
        background: none;
        border: none;
        cursor: pointer;
        font-size: 1rem;
        padding: 0.25rem;
        border-radius: 0.25rem;
        transition: background-color 0.2s;
    }

    .expand-toggle:hover {
        background: #f1f3f4;
    }

    .node-bullet {
        color: #d4af37;
        font-size: 0.8rem;
        padding: 0.25rem;
    }

    .node-content {
        flex: 1;
        min-width: 0;
        padding-right: 1rem;
    }

    .node-header {
        display: flex;
        align-items: flex-start;
        gap: 1rem;
        margin-bottom: 0.5rem;
    }

    .node-title {
        flex: 1;
        margin: 0;
        font-size: 1rem;
        color: #212529;
        font-weight: 600;
    }

    .node-meta-inline {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.8rem;
    }

    .difficulty-badge {
        padding: 0.125rem 0.375rem;
        border-radius: 0.25rem;
        font-weight: 600;
        font-size: 0.7rem;
        color: white;
    }

    .difficulty-easy { background: #28a745; }
    .difficulty-medium { background: #ffc107; color: #212529; }
    .difficulty-hard { background: #fd7e14; }
    .difficulty-expert { background: #dc3545; }

    .children-count {
        color: #6c757d;
        font-size: 0.75rem;
    }

    .node-actions {
        display: flex;
        gap: 0.25rem;
    }

    .btn-icon {
        background: none;
        border: none;
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 0.25rem;
        font-size: 0.9rem;
        transition: all 0.2s;
        min-width: 1.75rem;
        height: 1.75rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .btn-icon:hover {
        background: #f1f3f4;
    }

    .btn-icon.ai-btn:hover {
        background: #fffbf0;
        color: #d4af37;
    }

    .btn-icon.delete:hover {
        background: #fee;
        color: #dc3545;
    }

    .node-description {
        margin: 0 0 0.5rem 0;
        color: #6c757d;
        font-size: 0.9rem;
        line-height: 1.4;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .ai-summary {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .ai-label {
        font-size: 0.8rem;
    }

    .concepts-inline {
        display: flex;
        gap: 0.25rem;
        flex-wrap: wrap;
    }

    .concept-tag {
        background: #fffbf0;
        border: 1px solid #d4af37;
        color: #212529;
        padding: 0.125rem 0.375rem;
        border-radius: 0.25rem;
        font-size: 0.7rem;
    }

    .more-concepts {
        color: #6c757d;
        font-size: 0.75rem;
        font-style: italic;
    }

    .node-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: #6c757d;
        font-size: 0.75rem;
        margin-top: 0.5rem;
    }

    .empty-tree {
        text-align: center;
        padding: 3rem 1rem;
        color: #6c757d;
    }

    .empty-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
    }

    .link-button {
        background: none;
        border: none;
        color: #d4af37;
        cursor: pointer;
        text-decoration: underline;
    }

    .tree-stats {
        padding: 0.75rem 1rem;
        background: #f8f9fa;
        border-top: 1px solid #e9ecef;
        font-size: 0.8rem;
        color: #6c757d;
        text-align: center;
    }

    /* Responsive design */
    @media (max-width: 768px) {
        .tree-controls {
            flex-direction: column;
            align-items: stretch;
        }

        .expand-controls {
            justify-content: center;
        }

        .node-header {
            flex-direction: column;
            gap: 0.5rem;
        }

        .node-actions {
            justify-content: flex-end;
        }

        .node-indent {
            min-width: calc(var(--indent) * 1rem + 1.5rem);
        }

        .indent-line {
            width: 1rem;
        }
    }
</style>
