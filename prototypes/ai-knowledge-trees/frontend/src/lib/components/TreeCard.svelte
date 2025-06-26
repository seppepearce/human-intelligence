<script>
    import { createEventDispatcher } from 'svelte';

    export let tree;
    export let showActions = true;

    const dispatch = createEventDispatcher();

    function handleView() {
        dispatch('view', { tree });
    }

    function handleDelete() {
        dispatch('delete', { tree });
    }

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

<div class="tree-card">
    <div class="tree-header">
        <h3 class="tree-title">{tree.title}</h3>
        {#if showActions}
            <div class="tree-actions">
                <button class="btn-icon" on:click={handleView} title="View Tree">
                    👁️
                </button>
                <button class="btn-icon delete" on:click={handleDelete} title="Delete Tree">
                    🗑️
                </button>
            </div>
        {/if}
    </div>

    {#if tree.description}
        <p class="tree-description">{tree.description}</p>
    {/if}

    <div class="tree-meta">
        <div class="meta-item">
            <span class="meta-label">Created:</span>
            <span class="meta-value">{formatDate(tree.created_at)}</span>
        </div>
        {#if tree.node_count !== undefined}
            <div class="meta-item">
                <span class="meta-label">Nodes:</span>
                <span class="meta-value">{tree.node_count}</span>
            </div>
        {/if}
        <div class="meta-item">
            <span class="meta-label">Visibility:</span>
            <span class="meta-value">{tree.is_public ? '🌍 Public' : '🔒 Private'}</span>
        </div>
        {#if tree.updated_at && tree.updated_at !== tree.created_at}
            <div class="meta-item">
                <span class="meta-label">Updated:</span>
                <span class="meta-value">{formatDate(tree.updated_at)}</span>
            </div>
        {/if}
    </div>

    <div class="tree-footer">
        <button class="btn btn-primary" on:click={handleView}>
            🌳 Explore Tree
        </button>
    </div>
</div>

<style>
    .tree-card {
        background: white;
        border: 1px solid #e9ecef;
        border-radius: 0.5rem;
        padding: 1.5rem;
        transition: all 0.2s ease;
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    .tree-card:hover {
        border-color: #d4af37;
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0,0,0,0.1);
    }

    .tree-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;
    }

    .tree-title {
        color: #212529;
        margin: 0;
        flex: 1;
        font-size: 1.25rem;
        font-family: Georgia, serif;
    }

    .tree-actions {
        display: flex;
        gap: 0.5rem;
        margin-left: 1rem;
    }

    .btn-icon {
        background: none;
        border: none;
        font-size: 1.2rem;
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 0.25rem;
        transition: background-color 0.2s;
    }

    .btn-icon:hover {
        background: #f8f9fa;
    }

    .btn-icon.delete:hover {
        background: #fee;
        color: #dc3545;
    }

    .tree-description {
        color: #6c757d;
        margin-bottom: 1rem;
        line-height: 1.5;
        flex-grow: 1;
    }

    .tree-meta {
        margin-bottom: 1.5rem;
        border-top: 1px solid #e9ecef;
        padding-top: 1rem;
    }

    .meta-item {
        display: flex;
        justify-content: space-between;
        margin-bottom: 0.5rem;
    }

    .meta-label {
        color: #6c757d;
        font-weight: 600;
    }

    .meta-value {
        color: #212529;
    }

    .tree-footer {
        text-align: center;
        margin-top: auto;
    }

    .btn {
        display: inline-block;
        padding: 0.75rem 1.5rem;
        margin: 0.25rem;
        border: none;
        border-radius: 0.25rem;
        text-decoration: none;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s ease;
        font-family: Georgia, serif;
    }

    .btn-primary {
        background: #d4af37;
        color: white;
    }

    .btn-primary:hover {
        background: #b8941f;
        transform: translateY(-1px);
    }

    /* Responsive design */
    @media (max-width: 768px) {
        .tree-header {
            flex-direction: column;
            gap: 1rem;
        }

        .tree-actions {
            margin-left: 0;
            justify-content: flex-end;
        }
    }
</style>
