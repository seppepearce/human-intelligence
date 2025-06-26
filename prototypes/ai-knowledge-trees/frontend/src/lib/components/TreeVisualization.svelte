<script>
    import { onMount, createEventDispatcher } from 'svelte';
    import * as d3 from 'd3';

    export let nodes = [];
    export let width = 1200;
    export let height = 800;
    export let showActions = true;
    export let analyzing = false;
    export let nodeWithAnalysis = null;

    const dispatch = createEventDispatcher();

    let svgElement;
    let containerElement;
    let tooltip;
    let svg, g, tree, root;
    let zoom;
    let tooltipPinned = false;
    let pinnedNode = null;

    // Classical color palette
    const colors = {
        marble: '#f8f9fa',
        marbleLight: '#ffffff',
        stone: '#6c757d',
        golden: '#d4af37',
        goldenLight: '#e6c968',
        goldenDark: '#b8941f',
        shadow: 'rgba(0,0,0,0.15)',
        text: '#212529',
        textLight: '#6c757d'
    };

    // Node size and spacing
    const nodeRadius = 25;
    const levelHeight = 150;
    const nodeSpacing = 200;

    onMount(() => {
        initializeVisualization();
        return () => cleanup();
    });

    $: if (svgElement && nodes.length > 0) {
        updateVisualization();
    }

    function initializeVisualization() {
        if (!svgElement) return;

        // Create SVG with zoom behavior
        svg = d3.select(svgElement)
            .attr('width', width)
            .attr('height', height)
            .style('background', `linear-gradient(135deg, ${colors.marble} 0%, ${colors.marbleLight} 100%)`)
            .style('border-radius', '0.5rem')
            .style('border', `2px solid ${colors.golden}`)
            .style('box-shadow', `0 8px 32px ${colors.shadow}`);

        // Create main group for zooming/panning
        g = svg.append('g').attr('class', 'tree-container');

        // Initialize zoom behavior
        zoom = d3.zoom()
            .scaleExtent([0.1, 3])
            .on('zoom', (event) => {
                g.attr('transform', event.transform);
            });

        svg.call(zoom);

        // Add click handler to close tooltip when clicking outside
        svg.on('click', (event) => {
            // Only close if clicking on the background (not on nodes)
            if (event.target === svgElement && tooltipPinned) {
                hideTooltip();
            }
        });

        // Create tree layout
        tree = d3.tree()
            .size([width - 200, height - 100])
            .separation((a, b) => (a.parent === b.parent ? 1 : 2) / a.depth);

        // Add classical decorative border
        svg.append('rect')
            .attr('x', 5)
            .attr('y', 5)
            .attr('width', width - 10)
            .attr('height', height - 10)
            .attr('fill', 'none')
            .attr('stroke', colors.golden)
            .attr('stroke-width', 1)
            .attr('stroke-dasharray', '5,5')
            .attr('opacity', 0.3);

        updateVisualization();
    }

    function updateVisualization() {
        if (!svg || !nodes.length) return;

        // Build hierarchy from flat node list
        const hierarchyData = buildHierarchy(nodes);
        root = d3.hierarchy(hierarchyData);

        // Calculate tree layout
        tree(root);

        // Center the tree
        const centerX = width / 2;
        const centerY = 50;
        root.descendants().forEach(d => {
            d.x += centerX - (root.x || 0);
            d.y = d.depth * levelHeight + centerY;
        });

        updateLinks();
        updateNodes();
    }

    function buildHierarchy(nodes) {
        const nodeMap = new Map();
        const rootChildren = [];

        // Create node map
        nodes.forEach(node => {
            nodeMap.set(node.id, { ...node, children: [] });
        });

        // Build parent-child relationships
        nodes.forEach(node => {
            if (node.parent_id && nodeMap.has(node.parent_id)) {
                nodeMap.get(node.parent_id).children.push(nodeMap.get(node.id));
            } else {
                rootChildren.push(nodeMap.get(node.id));
            }
        });

        return {
            id: 'virtual-root',
            title: 'Root',
            children: rootChildren
        };
    }

    function updateLinks() {
        const links = root.links().filter(d => d.source.data.id !== 'virtual-root');

        const linkSelection = g.selectAll('.tree-link')
            .data(links, d => `${d.source.data.id}-${d.target.data.id}`);

        // Remove old links
        linkSelection.exit()
            .transition()
            .duration(500)
            .style('opacity', 0)
            .remove();

        // Add new links
        const linkEnter = linkSelection.enter()
            .append('path')
            .attr('class', 'tree-link')
            .style('opacity', 0);

        // Update all links
        linkSelection.merge(linkEnter)
            .transition()
            .duration(500)
            .style('opacity', 1)
            .attr('d', d => {
                const sourceX = d.source.x;
                const sourceY = d.source.y;
                const targetX = d.target.x;
                const targetY = d.target.y;

                return `M ${sourceX},${sourceY}
                        C ${sourceX},${(sourceY + targetY) / 2}
                          ${targetX},${(sourceY + targetY) / 2}
                          ${targetX},${targetY}`;
            })
            .attr('stroke', colors.stone)
            .attr('stroke-width', 2)
            .attr('fill', 'none')
            .attr('stroke-linecap', 'round')
            .style('filter', 'drop-shadow(1px 1px 2px rgba(0,0,0,0.1))');
    }

    function updateNodes() {
        const nodeData = root.descendants().filter(d => d.data.id !== 'virtual-root');

        const nodeSelection = g.selectAll('.tree-node')
            .data(nodeData, d => d.data.id);

        // Remove old nodes
        nodeSelection.exit()
            .transition()
            .duration(500)
            .style('opacity', 0)
            .remove();

        // Add new nodes
        const nodeEnter = nodeSelection.enter()
            .append('g')
            .attr('class', 'tree-node')
            .style('opacity', 0)
            .style('cursor', 'pointer');

        // Add node circles (marble tablets)
        nodeEnter.append('circle')
            .attr('class', 'node-bg')
            .attr('r', nodeRadius);

        // Add difficulty indicator ring
        nodeEnter.append('circle')
            .attr('class', 'difficulty-ring')
            .attr('r', nodeRadius + 3)
            .attr('fill', 'none')
            .attr('stroke-width', 3);

        // Add node text
        nodeEnter.append('text')
            .attr('class', 'node-text')
            .attr('text-anchor', 'middle')
            .attr('dy', '.35em')
            .style('font-family', 'Georgia, serif')
            .style('font-weight', '600')
            .style('font-size', '11px')
            .style('fill', colors.text)
            .style('pointer-events', 'none');

        // Add concept indicators
        nodeEnter.append('circle')
            .attr('class', 'concept-indicator')
            .attr('r', 4)
            .attr('cx', nodeRadius - 8)
            .attr('cy', -nodeRadius + 8);

        // Update all nodes
        const nodeUpdate = nodeSelection.merge(nodeEnter);

        nodeUpdate
            .transition()
            .duration(500)
            .style('opacity', 1)
            .attr('transform', d => `translate(${d.x},${d.y})`);

        // Update node backgrounds
        nodeUpdate.select('.node-bg')
            .attr('fill', d => getNodeColor(d.data))
            .attr('stroke', d => getNodeStroke(d.data))
            .attr('stroke-width', d => isNodeAnalyzing(d.data) ? 3 : 2)
            .style('filter', d => getNodeFilter(d.data));

        // Update difficulty rings
        nodeUpdate.select('.difficulty-ring')
            .attr('stroke', d => getDifficultyColor(d.data))
            .style('opacity', d => d.data.ai_analysis?.difficulty ? 0.8 : 0);

        // Update text with intelligent wrapping
        nodeUpdate.select('.node-text')
            .text('')
            .each(function(d) {
                const text = d3.select(this);
                const words = d.data.title.split(/\s+/);
                const maxWidth = nodeRadius * 1.8;

                if (words.length === 1 && d.data.title.length <= 8) {
                    text.append('tspan').text(d.data.title);
                } else {
                    const line1 = words.slice(0, Math.ceil(words.length / 2)).join(' ');
                    const line2 = words.slice(Math.ceil(words.length / 2)).join(' ');

                    text.append('tspan')
                        .attr('x', 0)
                        .attr('dy', '-0.3em')
                        .text(truncateText(line1, 10));

                    if (line2) {
                        text.append('tspan')
                            .attr('x', 0)
                            .attr('dy', '1em')
                            .text(truncateText(line2, 10));
                    }
                }
            });

        // Update concept indicators
        nodeUpdate.select('.concept-indicator')
            .style('opacity', d => d.data.ai_analysis?.concepts?.length > 0 ? 1 : 0)
            .attr('fill', colors.golden);

        // Add event handlers
        nodeUpdate
            .on('click', handleNodeClick)
            .on('mouseenter', handleNodeMouseEnter)
            .on('mouseleave', handleNodeMouseLeave)
            .on('contextmenu', handleNodeRightClick);
    }

    function getNodeColor(node) {
        if (isNodeAnalyzing(node)) return colors.goldenLight;
        if (node.ai_analysis?.difficulty >= 7) return '#ffebee';
        if (node.ai_analysis?.difficulty >= 4) return '#fff8e1';
        if (node.ai_analysis?.difficulty >= 1) return '#e8f5e8';
        return colors.marbleLight;
    }

    function getNodeStroke(node) {
        if (isNodeAnalyzing(node)) return colors.golden;
        if (node.ai_analysis?.difficulty >= 7) return '#e57373';
        if (node.ai_analysis?.difficulty >= 4) return '#ffb74d';
        if (node.ai_analysis?.difficulty >= 1) return '#81c784';
        return colors.golden;
    }

    function getNodeFilter(node) {
        const baseFilter = 'drop-shadow(2px 2px 4px rgba(0,0,0,0.15))';
        if (isNodeAnalyzing(node)) {
            return `${baseFilter} drop-shadow(0 0 8px ${colors.golden})`;
        }
        return baseFilter;
    }

    function getDifficultyColor(node) {
        const difficulty = node.ai_analysis?.difficulty || 0;
        if (difficulty >= 8) return '#d32f2f';
        if (difficulty >= 6) return '#f57c00';
        if (difficulty >= 4) return '#fbc02d';
        if (difficulty >= 2) return '#689f38';
        return '#2e7d32';
    }

    function isNodeAnalyzing(node) {
        return analyzing && nodeWithAnalysis?.id === node.id;
    }

    function truncateText(text, maxLength) {
        return text.length > maxLength ? text.substring(0, maxLength - 1) + '…' : text;
    }

    function handleNodeClick(event, d) {
        event.stopPropagation();

        // If clicking the same node that's already pinned, close it
        if (tooltipPinned && pinnedNode?.id === d.data.id) {
            hideTooltip();
            return;
        }

        tooltipPinned = true;
        pinnedNode = d.data;
        showNodeTooltip(d.data, event, true);
    }

    function handleNodeMouseEnter(event, d) {
        d3.select(event.currentTarget)
            .select('.node-bg')
            .transition()
            .duration(200)
            .attr('r', nodeRadius + 3)
            .attr('stroke-width', 3);
    }

    function handleNodeMouseLeave(event, d) {
        d3.select(event.currentTarget)
            .select('.node-bg')
            .transition()
            .duration(200)
            .attr('r', nodeRadius)
            .attr('stroke-width', 2);

        // Only hide tooltip if it's not pinned
        if (!tooltipPinned) {
            hideTooltip();
        }
    }

    function handleNodeRightClick(event, d) {
        event.preventDefault();
        showContextMenu(d.data, event);
    }

    function showNodeTooltip(node, event, pinned = false) {
        const tooltipContent = `
            <div class="tooltip-content">
                <div class="tooltip-header">
                    <h4>${node.title}</h4>
                    ${pinned ? '<button class="tooltip-close" onclick="closeTooltip()">✕</button>' : ''}
                </div>
                ${node.content ? `<p>${node.content.substring(0, 150)}${node.content.length > 150 ? '...' : ''}</p>` : ''}
                ${node.ai_analysis ? `
                    <div class="ai-info">
                        ${node.ai_analysis.difficulty ? `<span class="difficulty">Difficulty: ${node.ai_analysis.difficulty}/10</span>` : ''}
                        ${node.ai_analysis.concepts?.length > 0 ? `
                            <div class="concepts">
                                <strong>Concepts:</strong> ${node.ai_analysis.concepts.slice(0, 3).join(', ')}
                                ${node.ai_analysis.concepts.length > 3 ? ` (+${node.ai_analysis.concepts.length - 3} more)` : ''}
                            </div>
                        ` : ''}
                    </div>
                ` : ''}
                <div class="actions">
                    <button onclick="editNode('${node.id}')">✏️ Edit</button>
                    <button onclick="addChild('${node.id}')">➕ Add Child</button>
                    <button onclick="analyzeNode('${node.id}')">🧠 Analyze</button>
                </div>
            </div>
        `;

        // Remove existing tooltip
        d3.select('body').selectAll('.tree-tooltip').remove();

        tooltip = d3.select('body')
            .append('div')
            .attr('class', 'tree-tooltip')
            .style('opacity', 0)
            .html(tooltipContent);

        // Add hover events to keep tooltip open when hovering over it
        tooltip.node().addEventListener('mouseenter', () => {
            // Keep tooltip open when hovering
        });

        tooltip.node().addEventListener('mouseleave', () => {
            if (!tooltipPinned) {
                hideTooltip();
            }
        });

        // Prevent tooltip clicks from bubbling up to SVG
        tooltip.node().addEventListener('click', (e) => {
            e.stopPropagation();
        });

        tooltip.transition()
            .duration(200)
            .style('opacity', 1)
            .style('left', (event.pageX + 10) + 'px')
            .style('top', (event.pageY - 10) + 'px');
    }

    function hideTooltip() {
        if (tooltip) {
            tooltip.transition()
                .duration(200)
                .style('opacity', 0)
                .remove();
        }
        tooltipPinned = false;
        pinnedNode = null;
    }

    function showContextMenu(node, event) {
        // Dispatch events for parent component to handle
        dispatch('edit', { node });
    }

    // Zoom controls
    function zoomIn() {
        svg.transition().call(zoom.scaleBy, 1.5);
    }

    function zoomOut() {
        svg.transition().call(zoom.scaleBy, 1 / 1.5);
    }

    function resetZoom() {
        svg.transition().call(zoom.transform, d3.zoomIdentity);
    }

    function fitToContent() {
        if (!root) return;

        const nodes = root.descendants();
        if (nodes.length === 0) return;

        const bounds = {
            x: d3.min(nodes, d => d.x) - nodeRadius - 20,
            y: d3.min(nodes, d => d.y) - nodeRadius - 20,
            width: d3.max(nodes, d => d.x) - d3.min(nodes, d => d.x) + 2 * nodeRadius + 40,
            height: d3.max(nodes, d => d.y) - d3.min(nodes, d => d.y) + 2 * nodeRadius + 40
        };

        const scale = Math.min(width / bounds.width, height / bounds.height, 1);
        const translateX = (width - bounds.width * scale) / 2 - bounds.x * scale;
        const translateY = (height - bounds.height * scale) / 2 - bounds.y * scale;

        svg.transition()
            .duration(750)
            .call(zoom.transform, d3.zoomIdentity.translate(translateX, translateY).scale(scale));
    }

    function cleanup() {
        hideTooltip();
        d3.select('body').selectAll('.tree-tooltip').remove();
    }

    // Make functions available globally for tooltip buttons
    if (typeof window !== 'undefined') {
        window.editNode = (nodeId) => {
            const node = nodes.find(n => n.id === nodeId);
            if (node) dispatch('edit', { node });
            hideTooltip();
        };

        window.addChild = (nodeId) => {
            const node = nodes.find(n => n.id === nodeId);
            if (node) dispatch('addChild', { node });
            hideTooltip();
        };

        window.analyzeNode = (nodeId) => {
            const node = nodes.find(n => n.id === nodeId);
            if (node) dispatch('analyze', { node });
            hideTooltip();
        };

        window.closeTooltip = () => {
            hideTooltip();
        };
    }
</script>

<div class="tree-visualization" bind:this={containerElement}>
    <!-- Zoom Controls -->
    <div class="zoom-controls">
        <button class="zoom-btn" on:click={zoomIn} title="Zoom In">🔍+</button>
        <button class="zoom-btn" on:click={zoomOut} title="Zoom Out">🔍−</button>
        <button class="zoom-btn" on:click={resetZoom} title="Reset Zoom">⌂</button>
        <button class="zoom-btn" on:click={fitToContent} title="Fit to Content">⊞</button>
    </div>

    <!-- Tree Legend -->
    <div class="tree-legend">
        <div class="legend-item">
            <div class="legend-node easy"></div>
            <span>Easy (1-3)</span>
        </div>
        <div class="legend-item">
            <div class="legend-node medium"></div>
            <span>Medium (4-6)</span>
        </div>
        <div class="legend-item">
            <div class="legend-node hard"></div>
            <span>Hard (7-10)</span>
        </div>
        <div class="legend-item">
            <div class="legend-indicator concepts"></div>
            <span>Has Concepts</span>
        </div>
    </div>

    <!-- SVG Container -->
    <svg bind:this={svgElement}></svg>

    <!-- Instructions -->
    <div class="instructions">
        <p>🖱️ <strong>Click</strong> nodes for details • <strong>Drag</strong> to pan • <strong>Scroll</strong> to zoom</p>
    </div>
</div>

<style>
    .tree-visualization {
        position: relative;
        width: 100%;
        height: 100%;
        background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
        border-radius: 0.5rem;
        overflow: hidden;
        font-family: Georgia, serif;
    }

    .zoom-controls {
        position: absolute;
        top: 1rem;
        right: 1rem;
        display: flex;
        gap: 0.5rem;
        z-index: 10;
    }

    .zoom-btn {
        background: rgba(255, 255, 255, 0.9);
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        padding: 0.5rem;
        cursor: pointer;
        font-size: 0.9rem;
        transition: all 0.2s;
        backdrop-filter: blur(4px);
    }

    .zoom-btn:hover {
        background: #d4af37;
        color: white;
        transform: translateY(-1px);
    }

    .tree-legend {
        position: absolute;
        top: 1rem;
        left: 1rem;
        background: rgba(255, 255, 255, 0.9);
        border: 1px solid #d4af37;
        border-radius: 0.5rem;
        padding: 1rem;
        z-index: 10;
        backdrop-filter: blur(4px);
    }

    .legend-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
        font-size: 0.8rem;
        color: #212529;
    }

    .legend-item:last-child {
        margin-bottom: 0;
    }

    .legend-node {
        width: 16px;
        height: 16px;
        border-radius: 50%;
        border: 2px solid;
    }

    .legend-node.easy {
        background: #e8f5e8;
        border-color: #81c784;
    }

    .legend-node.medium {
        background: #fff8e1;
        border-color: #ffb74d;
    }

    .legend-node.hard {
        background: #ffebee;
        border-color: #e57373;
    }

    .legend-indicator.concepts {
        width: 8px;
        height: 8px;
        background: #d4af37;
        border-radius: 50%;
        border: none;
    }

    .instructions {
        position: absolute;
        bottom: 1rem;
        left: 50%;
        transform: translateX(-50%);
        background: rgba(255, 255, 255, 0.9);
        padding: 0.75rem 1.5rem;
        border-radius: 2rem;
        border: 1px solid #d4af37;
        backdrop-filter: blur(4px);
        z-index: 10;
    }

    .instructions p {
        margin: 0;
        font-size: 0.8rem;
        color: #6c757d;
        text-align: center;
    }

    /* Global tooltip styles */
    :global(.tree-tooltip) {
        position: absolute;
        background: white;
        border: 2px solid #d4af37;
        border-radius: 0.5rem;
        padding: 1rem;
        box-shadow: 0 8px 32px rgba(0,0,0,0.15);
        max-width: 300px;
        z-index: 1000;
        font-family: Georgia, serif;
        backdrop-filter: blur(8px);
    }

    :global(.tree-tooltip .tooltip-header) {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 0.5rem;
    }

    :global(.tree-tooltip .tooltip-content h4) {
        margin: 0;
        color: #212529;
        font-size: 1.1rem;
        flex: 1;
    }

    :global(.tree-tooltip .tooltip-close) {
        background: none;
        border: none;
        font-size: 1.2rem;
        cursor: pointer;
        color: #6c757d;
        padding: 0;
        margin-left: 0.5rem;
        width: 1.5rem;
        height: 1.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        transition: all 0.2s;
    }

    :global(.tree-tooltip .tooltip-close:hover) {
        background: #f8f9fa;
        color: #212529;
    }

    :global(.tree-tooltip .tooltip-content p) {
        margin: 0 0 0.75rem 0;
        color: #6c757d;
        font-size: 0.9rem;
        line-height: 1.4;
    }

    :global(.tree-tooltip .ai-info) {
        background: #fffbf0;
        border: 1px solid #d4af37;
        border-radius: 0.25rem;
        padding: 0.5rem;
        margin-bottom: 0.75rem;
        font-size: 0.8rem;
    }

    :global(.tree-tooltip .difficulty) {
        display: inline-block;
        background: #d4af37;
        color: white;
        padding: 0.125rem 0.375rem;
        border-radius: 0.25rem;
        font-weight: 600;
        font-size: 0.7rem;
        margin-bottom: 0.25rem;
    }

    :global(.tree-tooltip .concepts) {
        color: #212529;
        line-height: 1.3;
    }

    :global(.tree-tooltip .actions) {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.75rem;
    }

    :global(.tree-tooltip .actions button) {
        flex: 1;
        background: #d4af37;
        color: white;
        border: none;
        border-radius: 0.25rem;
        padding: 0.375rem 0.5rem;
        font-size: 0.7rem;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    :global(.tree-tooltip .actions button:hover) {
        background: #b8941f;
    }
</style>
