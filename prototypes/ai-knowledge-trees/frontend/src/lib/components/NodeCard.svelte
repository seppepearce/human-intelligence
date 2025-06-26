<script>
    import { createEventDispatcher } from 'svelte';

    export let node;
    export let showActions = true;
    export let showChildren = true;
    export let children = [];
    export let analyzing = false;
    export let isChild = false;

    const dispatch = createEventDispatcher();

    function handleEdit() {
        dispatch('edit', { node });
    }

    function handleDelete() {
        dispatch('delete', { node });
    }

    function handleAddChild() {
        dispatch('addChild', { node });
    }

    function handleAnalyze() {
        dispatch('analyze', { node });
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

    function getDifficultyClass(difficulty) {
        if (difficulty <= 2) return 'level-easy';
        if (difficulty <= 5) return 'level-medium';
        if (difficulty <= 7) return 'level-hard';
        return 'level-expert';
    }
</script>

<div class="node-card" class:child-node={isChild} class:has-children={children.length > 0}>
    <div class="node-header">
        <h3 class="node-title" class:child-title={isChild}>{node.title}</h3>
        {#if showActions}
            <div class="node-actions">
                <button
                    class="btn-icon"
                    on:click={handleAddChild}
                    title="Add Child Node"
                >
                    ➕
                </button>
                <button
                    class="btn-icon"
                    on:click={handleEdit}
                    title="Edit Node"
                >
                    ✏️
                </button>
                <button
                    class="btn-icon ai-btn"
                    on:click={handleAnalyze}
                    title="AI Analysis"
                    disabled={analyzing}
                >
                    {analyzing ? '🔄' : '🧠'}
                </button>
                <button
                    class="btn-icon delete"
                    on:click={handleDelete}
                    title="Delete Node"
                >
                    🗑️
                </button>
            </div>
        {/if}
    </div>

    {#if node.content}
        <div class="node-content">
            <p>{node.content}</p>
        </div>
    {/if}

    <!-- AI Analysis Display -->
    {#if node.ai_analysis}
        <div class="ai-insights">
            <h5>🧠 AI Analysis</h5>

            {#if node.ai_analysis.difficulty !== undefined}
                <div class="difficulty">
                    <span class="difficulty-label">Difficulty: </span>
                    <span class="difficulty-badge {getDifficultyClass(node.ai_analysis.difficulty)}">
                        {node.ai_analysis.difficulty}/10
                    </span>
                </div>
            {/if}

            {#if node.ai_analysis.concepts && node.ai_analysis.concepts.length > 0}
                <div class="concepts">
                    <h6>Key Concepts:</h6>
                    <div class="concept-tags">
                        {#each node.ai_analysis.concepts as concept}
                            <span class="concept-tag">{concept}</span>
                        {/each}
                    </div>
                </div>
            {/if}

            {#if node.ai_analysis.suggestions && node.ai_analysis.suggestions.length > 0}
                <div class="suggestions">
                    <h6>AI Suggestions:</h6>
                    <ul>
                        {#each node.ai_analysis.suggestions as suggestion}
                            <li>{suggestion}</li>
                        {/each}
                    </ul>
                </div>
            {/if}
        </div>
    {/if}

    <div class="node-meta">
        <small>Created: {formatDate(node.created_at)}</small>
        {#if node.updated_at && node.updated_at !== node.created_at}
            <small>• Updated: {formatDate(node.updated_at)}</small>
        {/if}
        {#if node.depth !== undefined}
            <small>• Depth: {node.depth}</small>
        {/if}
    </div>

    <!-- Child Nodes -->
    {#if showChildren && children.length > 0}
        <div class="child-nodes">
            <div class="child-header">
                <h6>Child Nodes ({children.length})</h6>
            </div>
            {#each children as childNode (childNode.id)}
                <svelte:self
                    node={childNode}
                    {showActions}
                    showChildren={false}
                    isChild={true}
                    {analyzing}
                    on:edit
                    on:delete
                    on:addChild
                    on:analyze
                />
            {/each}
        </div>
    {/if}
</div>

<style>
    .node-card {
        background: white;
        border: 1px solid #e9ecef;
        border-radius: 0.5rem;
        padding: 1.5rem;
        margin-bottom: 1rem;
        transition: all 0.2s ease;
    }

    .node-card:hover {
        border-color: #d4af37;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .node-card.has-children {
        border-left: 4px solid #d4af37;
    }

    .node-card.child-node {
        margin-left: 2rem;
        border-left: 2px solid #e9ecef;
        background: #f8f9fa;
        padding: 1rem;
    }

    .node-card.child-node:hover {
        background: #f1f3f4;
        border-left-color: #d4af37;
    }

    .node-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;
    }

    .node-title {
        color: #212529;
        margin: 0;
        flex: 1;
        font-size: 1.25rem;
        font-family: Georgia, serif;
    }

    .child-title {
        font-size: 1.1rem;
    }

    .node-actions {
        display: flex;
        gap: 0.5rem;
        margin-left: 1rem;
    }

    .btn-icon {
        background: none;
        border: none;
        font-size: 1.1rem;
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 0.25rem;
        transition: background-color 0.2s;
        min-width: 2rem;
        height: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .btn-icon:hover {
        background: #f8f9fa;
    }

    .btn-icon.delete:hover {
        background: #fee;
        color: #dc3545;
    }

    .btn-icon.ai-btn:hover {
        background: #fffbf0;
        color: #d4af37;
    }

    .btn-icon:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .node-content {
        margin-bottom: 1rem;
        color: #212529;
        line-height: 1.6;
    }

    .node-content p {
        margin: 0;
    }

    .ai-insights {
        background: #fffbf0;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        padding: 1rem;
        margin: 1rem 0;
    }

    .ai-insights h5, .ai-insights h6 {
        color: #212529;
        margin: 0 0 0.5rem 0;
        font-size: 1rem;
    }

    .ai-insights h6 {
        font-size: 0.9rem;
        margin-top: 0.75rem;
    }

    .difficulty {
        margin-bottom: 0.75rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .difficulty-label {
        font-weight: 600;
        color: #212529;
    }

    .difficulty-badge {
        padding: 0.2rem 0.5rem;
        border-radius: 0.25rem;
        font-size: 0.8rem;
        font-weight: 600;
        color: white;
    }

    .difficulty-badge.level-easy {
        background: #28a745;
    }

    .difficulty-badge.level-medium {
        background: #ffc107;
        color: #212529;
    }

    .difficulty-badge.level-hard {
        background: #fd7e14;
    }

    .difficulty-badge.level-expert {
        background: #dc3545;
    }

    .concepts {
        margin-bottom: 0.75rem;
    }

    .concept-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.25rem;
        margin-top: 0.25rem;
    }

    .concept-tag {
        background: white;
        border: 1px solid #d4af37;
        color: #212529;
        padding: 0.2rem 0.5rem;
        border-radius: 0.25rem;
        font-size: 0.8rem;
    }

    .suggestions ul {
        margin: 0.5rem 0 0 0;
        padding-left: 1.5rem;
    }

    .suggestions li {
        margin-bottom: 0.25rem;
        color: #212529;
        font-size: 0.9rem;
    }

    .node-meta {
        color: #6c757d;
        font-size: 0.9rem;
        border-top: 1px solid #e9ecef;
        padding-top: 0.5rem;
        margin-top: 1rem;
    }

    .child-nodes {
        margin-top: 1.5rem;
        border-top: 1px solid #e9ecef;
        padding-top: 1rem;
    }

    .child-header {
        margin-bottom: 1rem;
    }

    .child-header h6 {
        color: #6c757d;
        margin: 0;
        font-size: 0.9rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    /* Responsive design */
    @media (max-width: 768px) {
        .node-header {
            flex-direction: column;
            gap: 1rem;
        }

        .node-actions {
            margin-left: 0;
            justify-content: flex-end;
        }

        .node-card.child-node {
            margin-left: 1rem;
        }

        .concept-tags {
            flex-direction: column;
            align-items: flex-start;
        }

        .difficulty {
            flex-direction: column;
            align-items: flex-start;
            gap: 0.25rem;
        }
    }
</style>
