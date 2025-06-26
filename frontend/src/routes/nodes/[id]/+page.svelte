<script>
  import { page } from "$app/stores";
  import { apiRequest, isAuthenticated, user } from "$lib/stores/auth.js";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    ArrowLeft,
    Calendar,
    User,
    Eye,
    ThumbsUp,
    ThumbsDown,
    Share2,
    Edit,
    Trash2,
    FileText,
    Video,
    Link,
    Code,
    Hash,
    ExternalLink,
    Play,
    Globe,
    Lock,
    Brain,
    AlertCircle,
  } from "lucide-svelte";

  let nodeId = $page.params.id;
  let node = null;
  let loading = true;
  let error = "";
  let voting = false;
  let userVote = null; // -1, 0, or 1

  // Node type configurations
  const nodeTypeConfig = {
    text: {
      icon: FileText,
      color: "text-neon-cyan",
      bg: "bg-neon-cyan/10",
      label: "Text Content",
    },
    video: {
      icon: Video,
      color: "text-neon-purple",
      bg: "bg-neon-purple/10",
      label: "Video Content",
    },
    link: {
      icon: Link,
      color: "text-neon-green",
      bg: "bg-neon-green/10",
      label: "External Link",
    },
    code: {
      icon: Code,
      color: "text-neon-yellow",
      bg: "bg-neon-yellow/10",
      label: "Code Example",
    },
  };

  // Fetch node data
  async function fetchNode() {
    loading = true;
    error = "";

    try {
      const response = await apiRequest(`/nodes/${nodeId}`);

      if (response.ok) {
        const data = await response.json();
        node = data.node;

        // Increment view count by calling the API
        incrementViewCount();
      } else if (response.status === 404) {
        error = "Node not found";
      } else if (response.status === 403) {
        error = "You don't have permission to view this node";
      } else {
        const errorData = await response.json();
        error = errorData.error || "Failed to load node";
      }
    } catch (err) {
      console.error("Fetch node error:", err);
      error = "Network error. Please try again.";
    } finally {
      loading = false;
    }
  }

  // Increment view count
  async function incrementViewCount() {
    // This would typically be handled by the GET request itself
    // or by a separate endpoint call
  }

  // Handle voting
  async function handleVote(voteType) {
    if (!$isAuthenticated) {
      goto("/login");
      return;
    }

    if (voting) return;

    voting = true;
    try {
      const response = await apiRequest(`/nodes/${nodeId}/vote`, {
        method: "POST",
        body: JSON.stringify({ vote_type: voteType }),
      });

      if (response.ok) {
        const data = await response.json();
        node = data.node;
        userVote = voteType;
      } else {
        const errorData = await response.json();
        console.error("Vote error:", errorData.error);
      }
    } catch (err) {
      console.error("Vote request error:", err);
    } finally {
      voting = false;
    }
  }

  // Handle delete
  async function handleDelete() {
    if (
      !confirm(
        "Are you sure you want to delete this node? This action cannot be undone."
      )
    ) {
      return;
    }

    try {
      const response = await apiRequest(`/nodes/${nodeId}`, {
        method: "DELETE",
      });

      if (response.ok) {
        goto("/nodes");
      } else {
        const errorData = await response.json();
        alert("Failed to delete node: " + (errorData.error || "Unknown error"));
      }
    } catch (err) {
      console.error("Delete error:", err);
      alert("Network error while deleting node");
    }
  }

  // Format date
  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  // Format content based on node type
  function formatContent(content, nodeType) {
    if (!content) return "";

    // Simple markdown-like formatting
    return content
      .replace(/\n/g, "<br>")
      .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
      .replace(/\*(.*?)\*/g, "<em>$1</em>")
      .replace(/`(.*?)`/g, "<code>$1</code>");
  }

  // Copy share link
  async function copyShareLink() {
    try {
      await navigator.clipboard.writeText(window.location.href);
      // You could add a toast notification here
    } catch (err) {
      console.error("Failed to copy link:", err);
    }
  }

  onMount(() => {
    if (!nodeId) {
      error = "Invalid node ID";
      loading = false;
      return;
    }
    fetchNode();
  });
</script>

<svelte:head>
  <title>{node ? node.title : "Loading..."} - Human Intelligence</title>
  <meta
    name="description"
    content={node
      ? node.description || "Learning node on Human Intelligence platform"
      : "Loading node..."}
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-10" />

  <div class="container-wide py-8">
    <!-- Navigation -->
    <div class="mb-6">
      <button
        on:click={() => goto("/nodes")}
        class="btn-ghost flex items-center space-x-2"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Back to Nodes</span>
      </button>
    </div>

    <!-- Content -->
    {#if loading}
      <div class="flex items-center justify-center py-16">
        <div class="text-center">
          <div
            class="w-16 h-16 border-4 border-neon-cyan/30 border-t-neon-cyan rounded-full animate-spin mx-auto mb-4"
          />
          <p class="text-gray-400">Loading node...</p>
        </div>
      </div>
    {:else if error}
      <div class="card border-red-500/30 bg-red-500/10">
        <div class="text-center py-12">
          <AlertCircle class="w-16 h-16 text-red-400 mx-auto mb-4" />
          <h2 class="text-2xl font-display font-bold text-red-400 mb-2">
            Error
          </h2>
          <p class="text-gray-300 mb-6">{error}</p>
          <div class="space-x-4">
            <button
              on:click={fetchNode}
              class="btn-secondary flex items-center"
            >
              Try Again
            </button>
            <button
              on:click={() => goto("/nodes")}
              class="btn-ghost flex items-center"
            >
              Go Back
            </button>
          </div>
        </div>
      </div>
    {:else if node}
      <!-- Main Content -->
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Node Content -->
        <div class="lg:col-span-3 space-y-6">
          <!-- Header -->
          <div class="card">
            <!-- Node Type Badge -->
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center space-x-3">
                <div
                  class="p-3 rounded-lg {nodeTypeConfig[node.node_type]?.bg ||
                    'bg-gray-700'}"
                >
                  <svelte:component
                    this={nodeTypeConfig[node.node_type]?.icon || FileText}
                    class="w-6 h-6 {nodeTypeConfig[node.node_type]?.color ||
                      'text-gray-400'}"
                  />
                </div>
                <div>
                  <span
                    class="text-sm font-medium text-gray-400 uppercase tracking-wide"
                  >
                    {nodeTypeConfig[node.node_type]?.label || node.node_type}
                  </span>
                  <div class="flex items-center space-x-2 mt-1">
                    {#if node.is_public}
                      <Globe class="w-4 h-4 text-neon-green" />
                      <span class="text-xs text-neon-green">Public</span>
                    {:else}
                      <Lock class="w-4 h-4 text-gray-500" />
                      <span class="text-xs text-gray-500">Private</span>
                    {/if}
                  </div>
                </div>
              </div>

              <!-- Actions -->
              {#if $isAuthenticated && node.created_by === $user?.id}
                <div class="flex items-center space-x-2">
                  <button class="btn-ghost p-2" title="Edit node">
                    <Edit class="w-4 h-4" />
                  </button>
                  <button
                    on:click={handleDelete}
                    class="btn-ghost p-2 text-red-400 hover:text-red-300"
                    title="Delete node"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              {/if}
            </div>

            <!-- Title -->
            <h1 class="text-4xl font-display font-bold text-white mb-4">
              {node.title}
            </h1>

            <!-- Description -->
            {#if node.description}
              <p class="text-xl text-gray-300 mb-6">
                {node.description}
              </p>
            {/if}

            <!-- Metadata for specific types -->
            {#if node.node_type === "video" && node.metadata?.video_url}
              <div class="mb-6">
                <a
                  href={node.metadata.video_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center space-x-2 btn-secondary"
                >
                  <Play class="w-4 h-4" />
                  <span>Watch Video</span>
                  <ExternalLink class="w-4 h-4" />
                </a>
              </div>
            {:else if node.node_type === "link" && node.metadata?.link_url}
              <div class="mb-6">
                <a
                  href={node.metadata.link_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center space-x-2 btn-primary"
                >
                  <ExternalLink class="w-4 h-4" />
                  <span>Visit Link</span>
                </a>
              </div>
            {/if}

            <!-- Tags -->
            {#if node.tags && node.tags.length > 0}
              <div class="flex flex-wrap gap-2 mb-6">
                {#each node.tags as tag}
                  <span
                    class="inline-flex items-center bg-neon-purple/20 text-neon-purple border border-neon-purple/30 rounded-full px-3 py-1 text-sm"
                  >
                    <Hash class="w-3 h-3 mr-1" />
                    {tag}
                  </span>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Content -->
          {#if node.content}
            <div class="card">
              <h2 class="text-xl font-display font-bold text-neon-cyan mb-4">
                Content
              </h2>
              <div class="prose prose-invert max-w-none">
                <div class="text-gray-300 leading-relaxed">
                  {@html formatContent(node.content, node.node_type)}
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Stats & Actions -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-white mb-4">
              Actions
            </h3>

            <!-- Voting -->
            <div class="flex items-center space-x-4 mb-4">
              <button
                on:click={() => handleVote(1)}
                disabled={voting || !$isAuthenticated}
                class="btn-ghost p-2 {userVote === 1
                  ? 'text-neon-green'
                  : 'text-gray-400 hover:text-neon-green'}"
                title="Upvote"
              >
                <ThumbsUp class="w-5 h-5" />
              </button>

              <span
                class="text-lg font-medium {node.vote_score > 0
                  ? 'text-neon-green'
                  : node.vote_score < 0
                  ? 'text-red-400'
                  : 'text-gray-400'}"
              >
                {node.vote_score || 0}
              </span>

              <button
                on:click={() => handleVote(-1)}
                disabled={voting || !$isAuthenticated}
                class="btn-ghost p-2 {userVote === -1
                  ? 'text-red-400'
                  : 'text-gray-400 hover:text-red-400'}"
                title="Downvote"
              >
                <ThumbsDown class="w-5 h-5" />
              </button>
            </div>

            <!-- Share -->
            <button
              on:click={copyShareLink}
              class="btn-secondary w-full justify-center space-x-2 mb-4 flex items-center"
            >
              <Share2 class="w-4 h-4" />
              <span>Copy Link</span>
            </button>

            {#if !$isAuthenticated}
              <p class="text-xs text-gray-500 text-center">
                <a href="/login" class="text-neon-cyan hover:text-neon-pink"
                  >Sign in</a
                > to vote and interact
              </p>
            {/if}
          </div>

          <!-- Metadata -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-white mb-4">
              Details
            </h3>

            <div class="space-y-3 text-sm">
              <div class="flex items-center justify-between">
                <span class="text-gray-400">Views</span>
                <div class="flex items-center space-x-1 text-gray-300">
                  <Eye class="w-4 h-4" />
                  <span>{node.view_count || 0}</span>
                </div>
              </div>

              <div class="flex items-center justify-between">
                <span class="text-gray-400">Created</span>
                <div class="flex items-center space-x-1 text-gray-300">
                  <Calendar class="w-4 h-4" />
                  <span>{formatDate(node.created_at)}</span>
                </div>
              </div>

              {#if node.updated_at !== node.created_at}
                <div class="flex items-center justify-between">
                  <span class="text-gray-400">Updated</span>
                  <div class="flex items-center space-x-1 text-gray-300">
                    <Calendar class="w-4 h-4" />
                    <span>{formatDate(node.updated_at)}</span>
                  </div>
                </div>
              {/if}

              <div class="flex items-center justify-between">
                <span class="text-gray-400">Author</span>
                <div class="flex items-center space-x-1 text-gray-300">
                  <User class="w-4 h-4" />
                  <span>
                    {#if node.created_by === $user?.id}
                      You
                    {:else}
                      User
                    {/if}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Related Actions -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-neon-green mb-4">
              Explore More
            </h3>
            <div class="space-y-2">
              <a
                href="/nodes"
                class="btn-ghost w-full justify-start flex items-center"
              >
                <Brain class="w-4 h-4 mr-2" />
                Browse All Nodes
              </a>
              <a
                href="/nodes/create"
                class="btn-secondary w-full justify-start flex items-center"
              >
                <FileText class="w-4 h-4 mr-2" />
                Create New Node
              </a>
            </div>
          </div>
        </div>
      </div>
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

  .prose code {
    @apply bg-dark-700 text-neon-cyan px-2 py-1 rounded text-sm;
  }
</style>
