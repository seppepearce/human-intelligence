<script>
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated, user } from "$lib/stores/auth.js";
  import {
    Search,
    Filter,
    Plus,
    GitBranch,
    GitFork,
    Star,
    Clock,
    Users,
    BookOpen,
    Play,
    CheckCircle,
    Eye,
    ArrowUpDown,
    Zap,
    Target,
    TrendingUp,
  } from "lucide-svelte";

  // Search and filter state
  let searchQuery = "";
  let filterCategory = "all";
  let sortBy = "updated";
  let sortOrder = "desc";
  let loading = false;

  // Paths data
  let paths = [];
  let filteredPaths = [];
  let featuredPaths = [];

  // Categories for filtering
  const categories = [
    { value: "all", label: "All Categories" },
    { value: "programming", label: "Programming" },
    { value: "data-science", label: "Data Science" },
    { value: "design", label: "Design" },
    { value: "business", label: "Business" },
    { value: "languages", label: "Languages" },
    { value: "science", label: "Science" },
  ];

  const sortOptions = [
    { value: "updated", label: "Recently Updated" },
    { value: "created", label: "Recently Created" },
    { value: "popular", label: "Most Popular" },
    { value: "stars", label: "Most Starred" },
    { value: "forks", label: "Most Forked" },
    { value: "nodes", label: "Most Nodes" },
  ];

  onMount(() => {
    loadPaths();
  });

  async function loadPaths() {
    loading = true;

    try {
      // Mock data for now - replace with actual API call
      paths = [
        {
          id: "1",
          title: "Full Stack Web Development",
          description: "Complete journey from frontend to backend development",
          author: "alice_dev",
          author_avatar: null,
          category: "programming",
          nodes_count: 12,
          estimated_hours: 40,
          difficulty: "intermediate",
          stars: 89,
          forks: 23,
          completions: 156,
          is_featured: true,
          tags: ["javascript", "react", "node.js", "database"],
          created_at: "2024-01-15T10:00:00Z",
          updated_at: "2024-01-20T15:30:00Z",
          progress: $isAuthenticated ? 0.3 : 0,
          is_forked: false,
          is_starred: false,
          thumbnail: null,
        },
        {
          id: "2",
          title: "Machine Learning Fundamentals",
          description: "Learn the basics of ML from theory to practice",
          author: "bob_ml",
          author_avatar: null,
          category: "data-science",
          nodes_count: 8,
          estimated_hours: 25,
          difficulty: "beginner",
          stars: 67,
          forks: 15,
          completions: 89,
          is_featured: true,
          tags: ["python", "machine-learning", "statistics", "numpy"],
          created_at: "2024-01-10T09:00:00Z",
          updated_at: "2024-01-18T12:00:00Z",
          progress: $isAuthenticated ? 0.7 : 0,
          is_forked: true,
          is_starred: true,
          thumbnail: null,
        },
        {
          id: "3",
          title: "UI/UX Design Mastery",
          description: "From wireframes to prototypes and user testing",
          author: "charlie_design",
          author_avatar: null,
          category: "design",
          nodes_count: 10,
          estimated_hours: 30,
          difficulty: "intermediate",
          stars: 45,
          forks: 8,
          completions: 67,
          is_featured: false,
          tags: ["design", "figma", "user-research", "prototyping"],
          created_at: "2024-01-05T14:00:00Z",
          updated_at: "2024-01-16T11:00:00Z",
          progress: $isAuthenticated ? 0.1 : 0,
          is_forked: false,
          is_starred: false,
          thumbnail: null,
        },
        {
          id: "4",
          title: "Blockchain Development",
          description: "Build decentralized applications from scratch",
          author: "diana_blockchain",
          author_avatar: null,
          category: "programming",
          nodes_count: 15,
          estimated_hours: 50,
          difficulty: "advanced",
          stars: 78,
          forks: 19,
          completions: 34,
          is_featured: false,
          tags: ["blockchain", "ethereum", "solidity", "web3"],
          created_at: "2024-01-01T08:00:00Z",
          updated_at: "2024-01-19T16:00:00Z",
          progress: $isAuthenticated ? 0 : 0,
          is_forked: false,
          is_starred: true,
          thumbnail: null,
        },
      ];

      featuredPaths = paths.filter((path) => path.is_featured);
      filteredPaths = [...paths];
    } catch (error) {
      console.error("Failed to load paths:", error);
    } finally {
      loading = false;
    }
  }

  // Search and filter functions
  function handleSearch() {
    applyFilters();
  }

  function applyFilters() {
    filteredPaths = paths.filter((path) => {
      const matchesSearch =
        path.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        path.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
        path.tags.some((tag) =>
          tag.toLowerCase().includes(searchQuery.toLowerCase())
        );

      const matchesCategory =
        filterCategory === "all" || path.category === filterCategory;

      return matchesSearch && matchesCategory;
    });

    // Apply sorting
    filteredPaths.sort((a, b) => {
      let aVal, bVal;

      switch (sortBy) {
        case "updated":
          aVal = new Date(a.updated_at);
          bVal = new Date(b.updated_at);
          break;
        case "created":
          aVal = new Date(a.created_at);
          bVal = new Date(b.created_at);
          break;
        case "popular":
          aVal = a.completions;
          bVal = b.completions;
          break;
        case "stars":
          aVal = a.stars;
          bVal = b.stars;
          break;
        case "forks":
          aVal = a.forks;
          bVal = b.forks;
          break;
        case "nodes":
          aVal = a.nodes_count;
          bVal = b.nodes_count;
          break;
        default:
          return 0;
      }

      return sortOrder === "desc" ? bVal - aVal : aVal - bVal;
    });
  }

  function toggleSortOrder() {
    sortOrder = sortOrder === "desc" ? "asc" : "desc";
    applyFilters();
  }

  function handleSortChange() {
    applyFilters();
  }

  // Path actions
  function startPath(pathId) {
    goto(`/paths/${pathId}`);
  }

  function forkPath(pathId) {
    goto(`/paths/${pathId}/fork`);
  }

  function starPath(pathId) {
    // TODO: Implement starring
    console.log("Star path:", pathId);
  }

  // Utility functions
  function getDifficultyColor(difficulty) {
    switch (difficulty) {
      case "beginner":
        return "text-neon-green";
      case "intermediate":
        return "text-neon-yellow";
      case "advanced":
        return "text-neon-pink";
      default:
        return "text-gray-400";
    }
  }

  function formatTimeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;

    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    if (minutes > 0) return `${minutes}m ago`;
    return "just now";
  }

  function formatDuration(hours) {
    if (hours < 1) return `${Math.round(hours * 60)}min`;
    return `${hours}h`;
  }

  // Reactive statements
  $: {
    if (searchQuery !== undefined || filterCategory !== undefined) {
      applyFilters();
    }
  }
</script>

<svelte:head>
  <title>Learning Paths - Human Intelligence</title>
  <meta
    name="description"
    content="Discover and create learning paths on Human Intelligence platform"
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-20" />
  <div
    class="absolute inset-0 bg-gradient-to-br from-neon-purple/5 via-transparent to-neon-cyan/5"
  />

  <div class="relative container-wide py-8">
    <!-- Header -->
    <div class="mb-8">
      <div
        class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0"
      >
        <div class="flex items-center space-x-4">
          <GitBranch class="w-8 h-8 text-neon-purple" />
          <div>
            <h1 class="text-3xl font-display font-bold text-white">
              Learning Paths
            </h1>
            <p class="text-gray-400">Git-like sequences of learning nodes</p>
          </div>
        </div>

        {#if $isAuthenticated}
          <a
            href="/paths/create"
            class="btn-primary flex items-center space-x-2"
          >
            <Plus class="w-5 h-5" />
            <span>Create Path</span>
          </a>
        {/if}
      </div>
    </div>

    <!-- Featured Paths -->
    {#if featuredPaths.length > 0}
      <div class="mb-8">
        <h2
          class="text-xl font-display font-bold text-white mb-4 flex items-center space-x-2"
        >
          <Star class="w-5 h-5 text-neon-yellow" />
          <span>Featured Paths</span>
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {#each featuredPaths as path}
            <div
              class="card-neon group cursor-pointer"
              on:click={() => startPath(path.id)}
            >
              <div class="flex items-start space-x-4">
                <div
                  class="w-16 h-16 bg-gradient-to-br from-neon-purple to-neon-pink rounded-lg flex items-center justify-center"
                >
                  <GitBranch class="w-8 h-8 text-white" />
                </div>

                <div class="flex-1 min-w-0">
                  <h3
                    class="text-lg font-semibold text-white group-hover:text-neon-cyan transition-colors"
                  >
                    {path.title}
                  </h3>
                  <p class="text-gray-400 text-sm mt-1 line-clamp-2">
                    {path.description}
                  </p>

                  <div class="flex items-center space-x-4 mt-3 text-sm">
                    <span class="flex items-center space-x-1 text-gray-400">
                      <BookOpen class="w-4 h-4" />
                      <span>{path.nodes_count} nodes</span>
                    </span>
                    <span class="flex items-center space-x-1 text-gray-400">
                      <Clock class="w-4 h-4" />
                      <span>{formatDuration(path.estimated_hours)}</span>
                    </span>
                    <span
                      class="flex items-center space-x-1 {getDifficultyColor(
                        path.difficulty
                      )}"
                    >
                      <Target class="w-4 h-4" />
                      <span>{path.difficulty}</span>
                    </span>
                  </div>

                  {#if $isAuthenticated && path.progress > 0}
                    <div class="mt-3">
                      <div
                        class="flex justify-between text-xs text-gray-400 mb-1"
                      >
                        <span>Progress</span>
                        <span>{Math.round(path.progress * 100)}%</span>
                      </div>
                      <div class="w-full bg-dark-700 rounded-full h-2">
                        <div
                          class="bg-gradient-to-r from-neon-pink to-neon-cyan h-2 rounded-full transition-all duration-300"
                          style="width: {path.progress * 100}%"
                        />
                      </div>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Search and Filters -->
    <div class="card mb-8">
      <div
        class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0 lg:space-x-6"
      >
        <!-- Search -->
        <div class="flex-1 max-w-md">
          <div class="relative">
            <div
              class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
            >
              <Search class="w-5 h-5 text-gray-500" />
            </div>
            <input
              type="text"
              bind:value={searchQuery}
              on:keypress={(e) => e.key === "Enter" && handleSearch()}
              class="input-neon pl-10 w-full"
              placeholder="Search paths..."
            />
          </div>
        </div>

        <!-- Filters -->
        <div class="flex flex-wrap items-center gap-4">
          <div class="flex items-center space-x-2">
            <Filter class="w-4 h-4 text-gray-500" />
            <select bind:value={filterCategory} class="input-neon text-sm">
              {#each categories as category}
                <option value={category.value}>{category.label}</option>
              {/each}
            </select>
          </div>

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
              <ArrowUpDown class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Paths Grid -->
    {#if loading}
      <div class="flex justify-center items-center py-12">
        <div
          class="w-8 h-8 border-2 border-neon-pink border-t-transparent rounded-full animate-spin"
        />
      </div>
    {:else if filteredPaths.length === 0}
      <div class="text-center py-12">
        <GitBranch class="w-16 h-16 text-gray-600 mx-auto mb-4" />
        <h3 class="text-xl font-semibold text-gray-400 mb-2">No paths found</h3>
        <p class="text-gray-500 mb-6">
          Try adjusting your search or create a new path
        </p>
        {#if $isAuthenticated}
          <a href="/paths/create" class="btn-primary flex items-center">
            <Plus class="w-5 h-5 mr-2" />
            Create First Path
          </a>
        {/if}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each filteredPaths as path}
          <div
            class="card group hover:border-neon-purple transition-all duration-300"
          >
            <!-- Path Header -->
            <div class="flex items-start justify-between mb-4">
              <div class="flex items-center space-x-3">
                <div
                  class="w-10 h-10 bg-gradient-to-br from-neon-purple to-neon-cyan rounded-full flex items-center justify-center"
                >
                  <span class="text-white font-bold text-sm">
                    {path.author[0].toUpperCase()}
                  </span>
                </div>
                <div>
                  <div class="text-sm text-neon-cyan">@{path.author}</div>
                  <div class="text-xs text-gray-500">
                    {formatTimeAgo(path.updated_at)}
                  </div>
                </div>
              </div>

              <div class="flex items-center space-x-1">
                {#if path.is_starred}
                  <Star class="w-4 h-4 text-neon-yellow fill-current" />
                {/if}
                {#if path.is_forked}
                  <GitFork class="w-4 h-4 text-neon-green" />
                {/if}
              </div>
            </div>

            <!-- Path Content -->
            <div class="mb-4">
              <h3
                class="text-lg font-semibold text-white mb-2 group-hover:text-neon-cyan transition-colors cursor-pointer"
                on:click={() => startPath(path.id)}
              >
                {path.title}
              </h3>
              <p class="text-gray-400 text-sm line-clamp-2 mb-3">
                {path.description}
              </p>

              <!-- Tags -->
              <div class="flex flex-wrap gap-1 mb-3">
                {#each path.tags.slice(0, 3) as tag}
                  <span
                    class="px-2 py-1 bg-dark-700 text-neon-cyan text-xs rounded-full border border-dark-600"
                  >
                    {tag}
                  </span>
                {/each}
                {#if path.tags.length > 3}
                  <span
                    class="px-2 py-1 bg-dark-700 text-gray-400 text-xs rounded-full border border-dark-600"
                  >
                    +{path.tags.length - 3}
                  </span>
                {/if}
              </div>
            </div>

            <!-- Path Stats -->
            <div class="grid grid-cols-3 gap-4 mb-4 text-sm">
              <div class="text-center">
                <div class="text-gray-400">Nodes</div>
                <div class="font-semibold text-white">{path.nodes_count}</div>
              </div>
              <div class="text-center">
                <div class="text-gray-400">Duration</div>
                <div class="font-semibold text-white">
                  {formatDuration(path.estimated_hours)}
                </div>
              </div>
              <div class="text-center">
                <div class="text-gray-400">Level</div>
                <div
                  class="font-semibold {getDifficultyColor(path.difficulty)}"
                >
                  {path.difficulty}
                </div>
              </div>
            </div>

            <!-- Progress Bar -->
            {#if $isAuthenticated && path.progress > 0}
              <div class="mb-4">
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                  <span>Your Progress</span>
                  <span>{Math.round(path.progress * 100)}%</span>
                </div>
                <div class="w-full bg-dark-700 rounded-full h-2">
                  <div
                    class="bg-gradient-to-r from-neon-pink to-neon-cyan h-2 rounded-full transition-all duration-300"
                    style="width: {path.progress * 100}%"
                  />
                </div>
              </div>
            {/if}

            <!-- Actions -->
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4 text-sm text-gray-400">
                <span class="flex items-center space-x-1">
                  <Star class="w-4 h-4" />
                  <span>{path.stars}</span>
                </span>
                <span class="flex items-center space-x-1">
                  <GitFork class="w-4 h-4" />
                  <span>{path.forks}</span>
                </span>
                <span class="flex items-center space-x-1">
                  <CheckCircle class="w-4 h-4" />
                  <span>{path.completions}</span>
                </span>
              </div>

              {#if $isAuthenticated}
                <div class="flex items-center space-x-2">
                  <button
                    on:click={() => starPath(path.id)}
                    class="btn-ghost p-2 {path.is_starred
                      ? 'text-neon-yellow'
                      : 'text-gray-400'}"
                    title="Star path"
                  >
                    <Star
                      class="w-4 h-4 {path.is_starred ? 'fill-current' : ''}"
                    />
                  </button>
                  <button
                    on:click={() => forkPath(path.id)}
                    class="btn-ghost p-2 text-gray-400 hover:text-neon-green"
                    title="Fork path"
                  >
                    <GitFork class="w-4 h-4" />
                  </button>
                  <button
                    on:click={() => startPath(path.id)}
                    class="btn-secondary px-3 py-1 text-sm flex items-center"
                  >
                    <Play class="w-4 h-4 mr-1" />
                    Start
                  </button>
                </div>
              {:else}
                <button
                  on:click={() => startPath(path.id)}
                  class="btn-secondary px-3 py-1 text-sm flex items-center"
                >
                  <Eye class="w-4 h-4 mr-1" />
                  View
                </button>
              {/if}
            </div>
          </div>
        {/each}
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

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* Hover effects for path cards */
  .card:hover .progress-bar {
    box-shadow: 0 0 10px rgba(255, 0, 110, 0.5);
  }

  /* Custom scrollbar */
  ::-webkit-scrollbar {
    width: 8px;
  }

  ::-webkit-scrollbar-track {
    background: var(--dark-800);
    border-radius: 4px;
  }

  ::-webkit-scrollbar-thumb {
    background: linear-gradient(45deg, var(--neon-pink), var(--neon-purple));
    border-radius: 4px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(45deg, var(--neon-purple), var(--neon-cyan));
  }
</style>
