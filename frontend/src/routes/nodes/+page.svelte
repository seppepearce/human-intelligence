<script>
  import { apiRequest, isAuthenticated, user } from "$lib/stores/auth.js";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    Plus,
    FileText,
    Video,
    Link,
    Code,
    Search,
    Filter,
    SortAsc,
    SortDesc,
    Eye,
    ThumbsUp,
    Calendar,
    User,
    Hash,
    Brain,
  } from "lucide-svelte";

  let nodes = [];
  let loading = true;
  let error = "";
  let searchQuery = "";
  let sortBy = "created_at";
  let sortOrder = "desc";
  let filterType = "all";
  let showMyNodes = false;
  let currentPage = 1;
  let totalPages = 1;
  let limit = 12;

  // Node type configurations
  const nodeTypeConfig = {
    text: { icon: FileText, color: "text-neon-cyan", bg: "bg-neon-cyan/10" },
    video: { icon: Video, color: "text-neon-purple", bg: "bg-neon-purple/10" },
    link: { icon: Link, color: "text-neon-green", bg: "bg-neon-green/10" },
    code: { icon: Code, color: "text-neon-yellow", bg: "bg-neon-yellow/10" },
  };

  const sortOptions = [
    { value: "created_at", label: "Date Created" },
    { value: "updated_at", label: "Last Updated" },
    { value: "vote_score", label: "Popular" },
    { value: "view_count", label: "Most Viewed" },
    { value: "title", label: "Title A-Z" },
  ];

  // Fetch nodes from API
  async function fetchNodes() {
    loading = true;
    error = "";

    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        offset: ((currentPage - 1) * limit).toString(),
        sort: sortOrder === "desc" ? sortBy : sortBy,
      });

      if (searchQuery.trim()) {
        params.append("q", searchQuery.trim());
      }

      if (filterType !== "all") {
        params.append("type", filterType);
      }

      const endpoint =
        showMyNodes && $isAuthenticated
          ? `/users/me/nodes?${params}`
          : `/public/nodes?${params}`;

      const response = await apiRequest(endpoint);

      if (response.ok) {
        const data = await response.json();
        nodes = data.nodes || [];
        totalPages = Math.ceil((data.total || nodes.length) / limit);
      } else {
        const errorData = await response.json();
        error = errorData.error || "Failed to fetch nodes";
      }
    } catch (err) {
      console.error("Fetch nodes error:", err);
      error = "Network error. Please try again.";
    } finally {
      loading = false;
    }
  }

  // Handle search
  function handleSearch() {
    currentPage = 1;
    fetchNodes();
  }

  // Handle filter change
  function handleFilterChange() {
    currentPage = 1;
    fetchNodes();
  }

  // Handle sort change
  function handleSortChange() {
    currentPage = 1;
    fetchNodes();
  }

  // Toggle sort order
  function toggleSortOrder() {
    sortOrder = sortOrder === "asc" ? "desc" : "asc";
    handleSortChange();
  }

  // Format date
  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }

  // Navigate to node
  function viewNode(nodeId) {
    goto(`/nodes/${nodeId}`);
  }

  // Load data on mount and when filters change
  onMount(() => {
    fetchNodes();
  });

  $: if (
    searchQuery !== undefined ||
    filterType !== undefined ||
    showMyNodes !== undefined
  ) {
    handleFilterChange();
  }
</script>

<svelte:head>
  <title>Learning Nodes - Human Intelligence</title>
  <meta
    name="description"
    content="Explore learning nodes created by the Human Intelligence community."
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-10" />

  <div class="container-wide py-8">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <div class="flex items-center space-x-3 mb-2">
          <Brain class="w-8 h-8 text-neon-pink animate-pulse" />
          <h1 class="text-4xl font-display font-bold text-white">
            Learning <span class="text-neon-cyan">Nodes</span>
          </h1>
        </div>
        <p class="text-gray-400 text-lg">
          Discover and explore knowledge shared by the community
        </p>
      </div>

      {#if $isAuthenticated}
        <a href="/nodes/create" class="btn-primary flex items-center space-x-2">
          <Plus class="w-5 h-5" />
          <span>Create Node</span>
        </a>
      {/if}
    </div>

    <!-- Filters and Search -->
    <div class="card mb-8">
      <div
        class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0 lg:space-x-6"
      >
        <!-- Search -->
        <div class="flex-1 max-w-md">
          <div class="relative">
            <Search
              class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-500"
            />
            <input
              type="text"
              bind:value={searchQuery}
              on:keypress={(e) => e.key === "Enter" && handleSearch()}
              class="input-neon pl-10 w-full"
              placeholder="Search nodes..."
            />
          </div>
        </div>

        <!-- Filters -->
        <div class="flex flex-wrap items-center gap-4">
          <!-- Node Type Filter -->
          <div class="flex items-center space-x-2">
            <Filter class="w-4 h-4 text-gray-500" />
            <select bind:value={filterType} class="input-neon text-sm">
              <option value="all">All Types</option>
              <option value="text">Text</option>
              <option value="video">Video</option>
              <option value="link">Link</option>
              <option value="code">Code</option>
            </select>
          </div>

          <!-- Sort -->
          <div class="flex items-center space-x-2">
            <select
              bind:value={sortBy}
              on:change={handleSortChange}
              class="input-neon text-sm"
            >
              {#each sortOptions as option}
                <option value={option.value}>{option.label}</option>
              {/each}
            </select>
            <button
              on:click={toggleSortOrder}
              class="btn-ghost p-2"
              title="Toggle sort order"
            >
              {#if sortOrder === "desc"}
                <SortDesc class="w-4 h-4" />
              {:else}
                <SortAsc class="w-4 h-4" />
              {/if}
            </button>
          </div>

          <!-- My Nodes Toggle -->
          {#if $isAuthenticated}
            <label class="flex items-center space-x-2 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={showMyNodes}
                class="w-4 h-4 text-neon-pink border-gray-600 rounded focus:ring-neon-pink bg-dark-700"
              />
              <span class="text-sm text-gray-300">My Nodes</span>
            </label>
          {/if}
        </div>
      </div>
    </div>

    <!-- Content -->
    {#if loading}
      <div class="flex items-center justify-center py-16">
        <div class="text-center">
          <div
            class="w-16 h-16 border-4 border-neon-cyan/30 border-t-neon-cyan rounded-full animate-spin mx-auto mb-4"
          />
          <p class="text-gray-400">Loading nodes...</p>
        </div>
      </div>
    {:else if error}
      <div class="card border-red-500/30 bg-red-500/10">
        <div class="text-center py-8">
          <div
            class="w-16 h-16 bg-red-500/20 rounded-full flex items-center justify-center mx-auto mb-4"
          >
            <Brain class="w-8 h-8 text-red-400" />
          </div>
          <h3 class="text-lg font-semibold text-red-400 mb-2">
            Error Loading Nodes
          </h3>
          <p class="text-gray-400 mb-4">{error}</p>
          <button on:click={fetchNodes} class="btn-secondary">
            Try Again
          </button>
        </div>
      </div>
    {:else if nodes.length === 0}
      <div class="card">
        <div class="text-center py-16">
          <FileText class="w-16 h-16 text-gray-600 mx-auto mb-4" />
          <h3 class="text-xl font-semibold text-gray-400 mb-2">
            {showMyNodes ? "No nodes created yet" : "No nodes found"}
          </h3>
          <p class="text-gray-500 mb-6">
            {showMyNodes
              ? "Start creating and sharing your knowledge!"
              : searchQuery
              ? "Try adjusting your search or filters"
              : "Be the first to contribute to the community"}
          </p>
          {#if $isAuthenticated}
            <a
              href="/nodes/create"
              class="btn-primary flex items-center space-x-2"
            >
              <Plus class="w-5 h-5" />
              <span>Create Your First Node</span>
            </a>
          {/if}
        </div>
      </div>
    {:else}
      <!-- Nodes Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        {#each nodes as node (node.id)}
          <div
            class="card hover:border-neon-purple cursor-pointer transition-all duration-300 group"
            on:click={() => viewNode(node.id)}
            on:keypress={(e) => e.key === "Enter" && viewNode(node.id)}
            role="button"
            tabindex="0"
          >
            <!-- Node Type Badge -->
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center space-x-2">
                <div
                  class="p-2 rounded-lg {nodeTypeConfig[node.node_type]?.bg ||
                    'bg-gray-700'}"
                >
                  <svelte:component
                    this={nodeTypeConfig[node.node_type]?.icon || FileText}
                    class="w-4 h-4 {nodeTypeConfig[node.node_type]?.color ||
                      'text-gray-400'}"
                  />
                </div>
                <span
                  class="text-xs font-medium text-gray-400 uppercase tracking-wide"
                >
                  {node.node_type}
                </span>
              </div>
              <div class="text-xs text-gray-500">
                {formatDate(node.created_at)}
              </div>
            </div>

            <!-- Title and Description -->
            <h3
              class="text-lg font-semibold text-white mb-2 group-hover:text-neon-cyan transition-colors"
            >
              {node.title}
            </h3>

            {#if node.description}
              <p class="text-gray-400 text-sm mb-4 line-clamp-2">
                {node.description}
              </p>
            {/if}

            <!-- Tags -->
            {#if node.tags && node.tags.length > 0}
              <div class="flex flex-wrap gap-1 mb-4">
                {#each node.tags.slice(0, 3) as tag}
                  <span
                    class="inline-flex items-center bg-dark-700 text-gray-300 text-xs px-2 py-1 rounded-full"
                  >
                    <Hash class="w-3 h-3 mr-1" />
                    {tag}
                  </span>
                {/each}
                {#if node.tags.length > 3}
                  <span class="text-xs text-gray-500"
                    >+{node.tags.length - 3} more</span
                  >
                {/if}
              </div>
            {/if}

            <!-- Stats -->
            <div
              class="flex items-center justify-between text-sm text-gray-500 border-t border-dark-600 pt-3"
            >
              <div class="flex items-center space-x-4">
                <div class="flex items-center space-x-1">
                  <Eye class="w-4 h-4" />
                  <span>{node.view_count || 0}</span>
                </div>
                <div class="flex items-center space-x-1">
                  <ThumbsUp class="w-4 h-4" />
                  <span>{node.vote_score || 0}</span>
                </div>
              </div>

              {#if node.created_by === $user?.id}
                <span class="text-neon-cyan text-xs font-medium">Your Node</span
                >
              {/if}
            </div>
          </div>
        {/each}
      </div>

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="flex items-center justify-center space-x-2">
          <button
            on:click={() => {
              currentPage = Math.max(1, currentPage - 1);
              fetchNodes();
            }}
            disabled={currentPage === 1}
            class="btn-ghost px-3 py-2 {currentPage === 1
              ? 'opacity-50 cursor-not-allowed'
              : ''}"
          >
            Previous
          </button>

          {#each Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
            const startPage = Math.max(1, currentPage - 2);
            return startPage + i;
          }) as page}
            {#if page <= totalPages}
              <button
                on:click={() => {
                  currentPage = page;
                  fetchNodes();
                }}
                class="px-3 py-2 rounded transition-colors {currentPage === page
                  ? 'bg-neon-pink text-white'
                  : 'text-gray-400 hover:text-white hover:bg-dark-700'}"
              >
                {page}
              </button>
            {/if}
          {/each}

          <button
            on:click={() => {
              currentPage = Math.min(totalPages, currentPage + 1);
              fetchNodes();
            }}
            disabled={currentPage === totalPages}
            class="btn-ghost px-3 py-2 {currentPage === totalPages
              ? 'opacity-50 cursor-not-allowed'
              : ''}"
          >
            Next
          </button>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .bg-grid {
    background-image: linear-gradient(
        rgba(255, 0, 110, 0.1) 1px,
        transparent 1px
      ),
      linear-gradient(90deg, rgba(0, 245, 255, 0.1) 1px, transparent 1px);
    background-size: 20px 20px;
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
