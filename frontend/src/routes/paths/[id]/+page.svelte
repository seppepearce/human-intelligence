<script>
  import { onMount, onDestroy } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { isAuthenticated, user } from "$lib/stores/auth.js";
  import { liveview, useLiveTopic } from "$lib/stores/liveview.js";
  import {
    ArrowLeft,
    GitBranch,
    GitFork,
    Star,
    Play,
    Users,
    Clock,
    Target,
    BookOpen,
    Video,
    Link,
    Code,
    CheckCircle,
    Circle,
    Eye,
    MessageSquare,
    Share2,
    Download,
    Flag,
    MoreHorizontal,
    ChevronRight,
    TrendingUp,
    Calendar,
    Tag,
    Award,
    Zap,
  } from "lucide-svelte";

  // Path ID from URL
  $: pathId = $page.params.id;

  // Path data
  let pathData = null;
  let pathNodes = [];
  let userProgress = null;
  let loading = true;
  let error = null;

  // Real-time data
  const pathActivity = useLiveTopic(`path_${pathId}`, { activities: [] });

  // UI state
  let activeTab = "overview"; // 'overview', 'nodes', 'activity', 'reviews'
  let showShareModal = false;
  let showForkModal = false;

  onMount(() => {
    loadPathData();

    // Subscribe to real-time updates for this path
    if (pathId) {
      liveview.subscribe(`path_${pathId}`);
    }
  });

  onDestroy(() => {
    if (pathId) {
      liveview.unsubscribe(`path_${pathId}`);
    }
  });

  async function loadPathData() {
    loading = true;
    error = null;

    try {
      // Mock data - replace with actual API calls
      await new Promise((resolve) => setTimeout(resolve, 800));

      pathData = {
        id: pathId,
        title: "Full Stack Web Development",
        description:
          "Complete journey from frontend to backend development with modern JavaScript stack. Learn React, Node.js, databases, and deployment strategies.",
        author: "alice_dev",
        author_avatar: null,
        author_bio: "Full-stack developer with 8+ years experience",
        category: "programming",
        difficulty: "intermediate",
        estimated_hours: 40,
        stars: 89,
        forks: 23,
        completions: 156,
        views: 1247,
        is_featured: true,
        is_public: true,
        tags: ["javascript", "react", "node.js", "database", "fullstack"],
        created_at: "2024-01-15T10:00:00Z",
        updated_at: "2024-01-20T15:30:00Z",
        version: "1.2.0",
        is_starred: false,
        is_forked: false,
        fork_count: 23,
        completion_rate: 78,
        average_rating: 4.6,
        review_count: 34,
      };

      pathNodes = [
        {
          id: "1",
          title: "Introduction to Web Development",
          description:
            "Overview of web technologies and development environment setup",
          type: "text",
          estimated_minutes: 30,
          order: 0,
          is_required: true,
          is_completed: $isAuthenticated ? true : false,
          completion_rate: 92,
        },
        {
          id: "2",
          title: "HTML & CSS Fundamentals",
          description: "Learn the building blocks of web pages",
          type: "video",
          estimated_minutes: 90,
          order: 1,
          is_required: true,
          is_completed: $isAuthenticated ? true : false,
          completion_rate: 89,
        },
        {
          id: "3",
          title: "JavaScript Basics",
          description: "Programming fundamentals with JavaScript",
          type: "code",
          estimated_minutes: 120,
          order: 2,
          is_required: true,
          is_completed: $isAuthenticated ? false : false,
          completion_rate: 85,
        },
        {
          id: "4",
          title: "React Introduction",
          description: "Build your first React application",
          type: "link",
          estimated_minutes: 150,
          order: 3,
          is_required: true,
          is_completed: false,
          completion_rate: 76,
        },
        {
          id: "5",
          title: "Node.js Backend",
          description: "Create RESTful APIs with Node.js",
          type: "video",
          estimated_minutes: 180,
          order: 4,
          is_required: true,
          is_completed: false,
          completion_rate: 71,
        },
        {
          id: "6",
          title: "Database Integration",
          description: "Connect to databases and manage data",
          type: "text",
          estimated_minutes: 100,
          order: 5,
          is_required: false,
          is_completed: false,
          completion_rate: 68,
        },
      ];

      if ($isAuthenticated) {
        userProgress = {
          started_at: "2024-01-16T09:00:00Z",
          last_accessed: "2024-01-22T14:30:00Z",
          completed_nodes: 2,
          total_nodes: pathNodes.length,
          progress_percentage: 33,
          estimated_completion: "2024-02-15T00:00:00Z",
        };
      }
    } catch (err) {
      error = "Failed to load path data";
      console.error("Error loading path:", err);
    } finally {
      loading = false;
    }
  }

  // Actions
  function startPath() {
    if (!$isAuthenticated) {
      goto("/login");
      return;
    }

    const firstIncomplete = pathNodes.find((node) => !node.is_completed);
    if (firstIncomplete) {
      goto(`/nodes/${firstIncomplete.id}`);
    }
  }

  function forkPath() {
    if (!$isAuthenticated) {
      goto("/login");
      return;
    }

    goto(`/paths/${pathId}/fork`);
  }

  function toggleStar() {
    if (!$isAuthenticated) {
      goto("/login");
      return;
    }

    pathData.is_starred = !pathData.is_starred;
    pathData.stars += pathData.is_starred ? 1 : -1;

    // TODO: API call to update star status
  }

  function sharePath() {
    showShareModal = true;
  }

  function goToNode(nodeId) {
    goto(`/nodes/${nodeId}`);
  }

  // Utility functions
  function getNodeTypeIcon(type) {
    switch (type) {
      case "text":
        return BookOpen;
      case "video":
        return Video;
      case "link":
        return Link;
      case "code":
        return Code;
      default:
        return BookOpen;
    }
  }

  function getNodeTypeColor(type) {
    switch (type) {
      case "text":
        return "text-neon-cyan";
      case "video":
        return "text-neon-pink";
      case "link":
        return "text-neon-green";
      case "code":
        return "text-neon-purple";
      default:
        return "text-gray-400";
    }
  }

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

  function formatDuration(minutes) {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;

    if (hours > 0) {
      return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
    }
    return `${mins}m`;
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }

  $: completedNodes = pathNodes.filter((node) => node.is_completed).length;
  $: progressPercentage =
    userProgress?.progress_percentage ||
    (completedNodes / pathNodes.length) * 100;
</script>

<svelte:head>
  <title>{pathData?.title || "Loading..."} - Human Intelligence</title>
  <meta
    name="description"
    content={pathData?.description || "Learning path on Human Intelligence"}
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-20" />
  <div
    class="absolute inset-0 bg-gradient-to-br from-neon-purple/5 via-transparent to-neon-cyan/5"
  />

  <div class="relative container-wide py-8">
    {#if loading}
      <div class="flex justify-center items-center py-24">
        <div
          class="w-12 h-12 border-2 border-neon-purple border-t-transparent rounded-full animate-spin"
        />
      </div>
    {:else if error}
      <div class="text-center py-24">
        <div class="text-red-400 text-xl mb-4">{error}</div>
        <button on:click={() => goto("/paths")} class="btn-secondary">
          Back to Paths
        </button>
      </div>
    {:else if pathData}
      <!-- Header -->
      <div class="flex items-center space-x-4 mb-8">
        <button on:click={() => goto("/paths")} class="btn-ghost p-2">
          <ArrowLeft class="w-5 h-5" />
        </button>
        <div class="flex items-center space-x-2">
          <GitBranch class="w-6 h-6 text-neon-purple" />
          <span class="text-gray-400">Learning Path</span>
        </div>
      </div>

      <!-- Path Header -->
      <div class="card mb-8">
        <div class="flex flex-col lg:flex-row lg:items-start lg:space-x-8">
          <!-- Path Icon -->
          <div
            class="w-20 h-20 bg-gradient-to-br from-neon-purple to-neon-cyan rounded-xl flex items-center justify-center mb-6 lg:mb-0 flex-shrink-0"
          >
            <GitBranch class="w-10 h-10 text-white" />
          </div>

          <!-- Path Info -->
          <div class="flex-1 min-w-0">
            <div
              class="flex flex-col sm:flex-row sm:items-start sm:justify-between mb-4"
            >
              <div class="flex-1 min-w-0">
                <h1 class="text-3xl font-display font-bold text-white mb-2">
                  {pathData.title}
                </h1>
                <p class="text-gray-400 mb-4 leading-relaxed">
                  {pathData.description}
                </p>

                <!-- Author -->
                <div class="flex items-center space-x-3 mb-4">
                  <div
                    class="w-10 h-10 bg-gradient-to-br from-neon-pink to-neon-purple rounded-full flex items-center justify-center"
                  >
                    <span class="text-white font-bold text-sm">
                      {pathData.author[0].toUpperCase()}
                    </span>
                  </div>
                  <div>
                    <div class="text-neon-cyan font-medium">
                      @{pathData.author}
                    </div>
                    <div class="text-sm text-gray-500">
                      {pathData.author_bio}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="flex items-center space-x-3 mt-4 sm:mt-0">
                <button
                  on:click={toggleStar}
                  class="btn-ghost p-2 {pathData.is_starred
                    ? 'text-neon-yellow'
                    : 'text-gray-400'}"
                  title="Star this path"
                >
                  <Star
                    class="w-5 h-5 {pathData.is_starred ? 'fill-current' : ''}"
                  />
                </button>

                <button
                  on:click={forkPath}
                  class="btn-ghost p-2 text-gray-400 hover:text-neon-green"
                  title="Fork this path"
                >
                  <GitFork class="w-5 h-5" />
                </button>

                <button
                  on:click={sharePath}
                  class="btn-ghost p-2 text-gray-400 hover:text-neon-cyan"
                  title="Share this path"
                >
                  <Share2 class="w-5 h-5" />
                </button>

                <button
                  on:click={startPath}
                  class="btn-primary flex items-center space-x-2"
                >
                  <Play class="w-4 h-4" />
                  <span>{userProgress ? "Continue" : "Start Learning"}</span>
                </button>
              </div>
            </div>

            <!-- Stats -->
            <div
              class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-6 gap-4 mb-6"
            >
              <div class="text-center">
                <div class="text-xl font-bold text-neon-cyan">
                  {pathData.stars}
                </div>
                <div class="text-xs text-gray-400">Stars</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-neon-green">
                  {pathData.forks}
                </div>
                <div class="text-xs text-gray-400">Forks</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-neon-pink">
                  {pathData.completions}
                </div>
                <div class="text-xs text-gray-400">Completions</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-white">
                  {pathNodes.length}
                </div>
                <div class="text-xs text-gray-400">Nodes</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-white">
                  {pathData.estimated_hours}h
                </div>
                <div class="text-xs text-gray-400">Duration</div>
              </div>
              <div class="text-center">
                <div
                  class="text-xl font-bold {getDifficultyColor(
                    pathData.difficulty
                  )}"
                >
                  {pathData.difficulty}
                </div>
                <div class="text-xs text-gray-400">Level</div>
              </div>
            </div>

            <!-- Progress (if authenticated and started) -->
            {#if userProgress}
              <div class="bg-dark-700 rounded-lg p-4 border border-dark-600">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-300"
                    >Your Progress</span
                  >
                  <span class="text-sm text-neon-cyan"
                    >{Math.round(progressPercentage)}%</span
                  >
                </div>
                <div class="w-full bg-dark-800 rounded-full h-3 mb-2">
                  <div
                    class="bg-gradient-to-r from-neon-pink to-neon-cyan h-3 rounded-full transition-all duration-500"
                    style="width: {progressPercentage}%"
                  />
                </div>
                <div
                  class="flex items-center justify-between text-xs text-gray-500"
                >
                  <span
                    >{completedNodes} of {pathNodes.length} nodes completed</span
                  >
                  <span
                    >Last accessed {formatDate(
                      userProgress.last_accessed
                    )}</span
                  >
                </div>
              </div>
            {/if}

            <!-- Tags -->
            <div class="flex flex-wrap gap-2 mt-4">
              {#each pathData.tags as tag}
                <span
                  class="px-3 py-1 bg-dark-700 text-neon-cyan text-sm rounded-full border border-dark-600"
                >
                  #{tag}
                </span>
              {/each}
            </div>
          </div>
        </div>
      </div>

      <!-- Tab Navigation -->
      <div class="flex bg-dark-800 rounded-lg p-1 border border-dark-600 mb-8">
        <button
          on:click={() => (activeTab = "overview")}
          class="px-4 py-2 rounded text-sm transition-all {activeTab ===
          'overview'
            ? 'bg-neon-purple text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Overview
        </button>
        <button
          on:click={() => (activeTab = "nodes")}
          class="px-4 py-2 rounded text-sm transition-all {activeTab === 'nodes'
            ? 'bg-neon-purple text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Nodes ({pathNodes.length})
        </button>
        <button
          on:click={() => (activeTab = "activity")}
          class="px-4 py-2 rounded text-sm transition-all {activeTab ===
          'activity'
            ? 'bg-neon-purple text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Activity
        </button>
        <button
          on:click={() => (activeTab = "reviews")}
          class="px-4 py-2 rounded text-sm transition-all {activeTab ===
          'reviews'
            ? 'bg-neon-purple text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Reviews ({pathData.review_count})
        </button>
      </div>

      <!-- Tab Content -->
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Main Content -->
        <div class="lg:col-span-3">
          {#if activeTab === "overview"}
            <!-- Path Overview -->
            <div class="space-y-6">
              <!-- Learning Objectives -->
              <div class="card">
                <h2 class="text-xl font-display font-bold text-white mb-4">
                  What You'll Learn
                </h2>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300"
                      >Modern JavaScript ES6+ features</span
                    >
                  </div>
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300"
                      >React components and state management</span
                    >
                  </div>
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300">Node.js server development</span
                    >
                  </div>
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300"
                      >Database design and integration</span
                    >
                  </div>
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300">RESTful API development</span>
                  </div>
                  <div class="flex items-start space-x-2">
                    <CheckCircle
                      class="w-5 h-5 text-neon-green mt-0.5 flex-shrink-0"
                    />
                    <span class="text-gray-300"
                      >Deployment and DevOps basics</span
                    >
                  </div>
                </div>
              </div>

              <!-- Path Visualization -->
              <div class="card">
                <h2 class="text-xl font-display font-bold text-white mb-4">
                  Learning Journey
                </h2>
                <div class="space-y-4">
                  {#each pathNodes as node, index}
                    <div class="flex items-center space-x-4 group">
                      <!-- Progress Indicator -->
                      <div class="flex flex-col items-center">
                        <div
                          class="w-8 h-8 rounded-full border-2 {node.is_completed
                            ? 'bg-neon-green border-neon-green'
                            : 'border-gray-600'} flex items-center justify-center"
                        >
                          {#if node.is_completed}
                            <CheckCircle class="w-5 h-5 text-white" />
                          {:else}
                            <span
                              class="text-sm font-semibold {node.is_completed
                                ? 'text-white'
                                : 'text-gray-400'}"
                            >
                              {index + 1}
                            </span>
                          {/if}
                        </div>
                        {#if index < pathNodes.length - 1}
                          <div
                            class="w-0.5 h-8 {node.is_completed
                              ? 'bg-neon-green'
                              : 'bg-gray-600'} mt-2"
                          />
                        {/if}
                      </div>

                      <!-- Node Info -->
                      <div
                        class="flex-1 p-4 bg-dark-700 rounded-lg border border-dark-600 group-hover:border-neon-purple transition-colors cursor-pointer"
                        on:click={() => goToNode(node.id)}
                      >
                        <div class="flex items-center space-x-3 mb-2">
                          <svelte:component
                            this={getNodeTypeIcon(node.type)}
                            class="w-5 h-5 {getNodeTypeColor(node.type)}"
                          />
                          <h3
                            class="font-semibold text-white group-hover:text-neon-cyan transition-colors"
                          >
                            {node.title}
                          </h3>
                          {#if !node.is_required}
                            <span
                              class="text-xs bg-dark-600 text-gray-400 px-2 py-1 rounded"
                              >Optional</span
                            >
                          {/if}
                        </div>
                        <p class="text-gray-400 text-sm mb-2">
                          {node.description}
                        </p>
                        <div
                          class="flex items-center justify-between text-xs text-gray-500"
                        >
                          <span class="flex items-center space-x-1">
                            <Clock class="w-3 h-3" />
                            <span>{formatDuration(node.estimated_minutes)}</span
                            >
                          </span>
                          <span
                            >{node.completion_rate}% completed by learners</span
                          >
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            </div>
          {:else if activeTab === "nodes"}
            <!-- Detailed Nodes List -->
            <div class="space-y-4">
              {#each pathNodes as node, index}
                <div
                  class="card group cursor-pointer hover:border-neon-purple transition-colors"
                  on:click={() => goToNode(node.id)}
                >
                  <div class="flex items-center space-x-4">
                    <!-- Node Number -->
                    <div
                      class="w-10 h-10 rounded-full bg-neon-purple flex items-center justify-center text-white font-bold"
                    >
                      {index + 1}
                    </div>

                    <!-- Node Content -->
                    <div class="flex-1">
                      <div class="flex items-center space-x-3 mb-2">
                        <svelte:component
                          this={getNodeTypeIcon(node.type)}
                          class="w-5 h-5 {getNodeTypeColor(node.type)}"
                        />
                        <h3
                          class="text-lg font-semibold text-white group-hover:text-neon-cyan transition-colors"
                        >
                          {node.title}
                        </h3>
                        {#if node.is_completed}
                          <CheckCircle class="w-5 h-5 text-neon-green" />
                        {/if}
                        {#if !node.is_required}
                          <span
                            class="text-xs bg-dark-600 text-gray-400 px-2 py-1 rounded"
                            >Optional</span
                          >
                        {/if}
                      </div>

                      <p class="text-gray-400 mb-3">{node.description}</p>

                      <div
                        class="flex items-center space-x-6 text-sm text-gray-500"
                      >
                        <span class="flex items-center space-x-1">
                          <Clock class="w-4 h-4" />
                          <span>{formatDuration(node.estimated_minutes)}</span>
                        </span>
                        <span class="flex items-center space-x-1">
                          <Users class="w-4 h-4" />
                          <span>{node.completion_rate}% completion rate</span>
                        </span>
                      </div>
                    </div>

                    <!-- Arrow -->
                    <ChevronRight
                      class="w-5 h-5 text-gray-500 group-hover:text-neon-cyan transition-colors"
                    />
                  </div>
                </div>
              {/each}
            </div>
          {:else if activeTab === "activity"}
            <!-- Real-time Activity -->
            <div class="card">
              <h2 class="text-xl font-display font-bold text-white mb-4">
                Recent Activity
              </h2>
              <div class="space-y-3">
                <div
                  class="flex items-center space-x-3 p-3 bg-dark-700 rounded-lg"
                >
                  <div
                    class="w-8 h-8 bg-gradient-to-br from-neon-green to-neon-cyan rounded-full flex items-center justify-center"
                  >
                    <span class="text-white text-sm font-bold">B</span>
                  </div>
                  <div class="flex-1">
                    <p class="text-white text-sm">
                      <span class="text-neon-cyan">@bob_learns</span> completed
                      <span class="text-neon-green">"JavaScript Basics"</span>
                    </p>
                    <p class="text-xs text-gray-500">2 hours ago</p>
                  </div>
                </div>

                <div
                  class="flex items-center space-x-3 p-3 bg-dark-700 rounded-lg"
                >
                  <div
                    class="w-8 h-8 bg-gradient-to-br from-neon-pink to-neon-purple rounded-full flex items-center justify-center"
                  >
                    <span class="text-white text-sm font-bold">C</span>
                  </div>
                  <div class="flex-1">
                    <p class="text-white text-sm">
                      <span class="text-neon-cyan">@charlie_dev</span> forked this
                      path
                    </p>
                    <p class="text-xs text-gray-500">5 hours ago</p>
                  </div>
                </div>

                <div
                  class="flex items-center space-x-3 p-3 bg-dark-700 rounded-lg"
                >
                  <div
                    class="w-8 h-8 bg-gradient-to-br from-neon-yellow to-neon-orange rounded-full flex items-center justify-center"
                  >
                    <span class="text-white text-sm font-bold">D</span>
                  </div>
                  <div class="flex-1">
                    <p class="text-white text-sm">
                      <span class="text-neon-cyan">@diana_codes</span> started this
                      path
                    </p>
                    <p class="text-xs text-gray-500">1 day ago</p>
                  </div>
                </div>
              </div>
            </div>
          {:else if activeTab === "reviews"}
            <!-- Reviews and Ratings -->
            <div class="space-y-6">
              <!-- Rating Summary -->
              <div class="card">
                <div class="flex items-center space-x-6 mb-6">
                  <div class="text-center">
                    <div class="text-4xl font-bold text-neon-cyan">
                      {pathData.average_rating}
                    </div>
                    <div class="flex justify-center space-x-1 mb-1">
                      {#each Array(5) as _, i}
                        <Star
                          class="w-4 h-4 {i <
                          Math.floor(pathData.average_rating)
                            ? 'text-neon-yellow fill-current'
                            : 'text-gray-600'}"
                        />
                      {/each}
                    </div>
                    <div class="text-sm text-gray-400">
                      {pathData.review_count} reviews
                    </div>
                  </div>

                  <div class="flex-1">
                    {#each [5, 4, 3, 2, 1] as rating}
                      <div class="flex items-center space-x-2 mb-1">
                        <span class="text-sm text-gray-400 w-2">{rating}</span>
                        <div class="flex-1 bg-dark-700 rounded-full h-2">
                          <div
                            class="bg-neon-yellow h-2 rounded-full"
                            style="width: {rating === 5
                              ? 70
                              : rating === 4
                              ? 20
                              : rating === 3
                              ? 7
                              : rating === 2
                              ? 2
                              : 1}%"
                          />
                        </div>
                        <span class="text-sm text-gray-500 w-8">
                          {rating === 5
                            ? 24
                            : rating === 4
                            ? 7
                            : rating === 3
                            ? 2
                            : rating === 2
                            ? 1
                            : 0}
                        </span>
                      </div>
                    {/each}
                  </div>
                </div>
              </div>

              <!-- Individual Reviews -->
              <div class="space-y-4">
                <div class="card">
                  <div class="flex items-start space-x-4">
                    <div
                      class="w-10 h-10 bg-gradient-to-br from-neon-pink to-neon-purple rounded-full flex items-center justify-center"
                    >
                      <span class="text-white font-bold text-sm">J</span>
                    </div>
                    <div class="flex-1">
                      <div class="flex items-center space-x-2 mb-2">
                        <span class="font-medium text-white">@john_learner</span
                        >
                        <div class="flex space-x-1">
                          {#each Array(5) as _, i}
                            <Star
                              class="w-4 h-4 {i < 5
                                ? 'text-neon-yellow fill-current'
                                : 'text-gray-600'}"
                            />
                          {/each}
                        </div>
                        <span class="text-xs text-gray-500">2 days ago</span>
                      </div>
                      <p class="text-gray-300 text-sm leading-relaxed">
                        Excellent path for learning full-stack development! The
                        progression from basic concepts to advanced topics is
                        well-structured. The hands-on projects really helped
                        solidify my understanding.
                      </p>
                      <div
                        class="flex items-center space-x-4 mt-3 text-xs text-gray-500"
                      >
                        <button class="hover:text-neon-cyan transition-colors"
                          >👍 Helpful (12)</button
                        >
                        <button class="hover:text-neon-cyan transition-colors"
                          >Reply</button
                        >
                      </div>
                    </div>
                  </div>
                </div>

                <div class="card">
                  <div class="flex items-start space-x-4">
                    <div
                      class="w-10 h-10 bg-gradient-to-br from-neon-green to-neon-cyan rounded-full flex items-center justify-center"
                    >
                      <span class="text-white font-bold text-sm">S</span>
                    </div>
                    <div class="flex-1">
                      <div class="flex items-center space-x-2 mb-2">
                        <span class="font-medium text-white">@sarah_codes</span>
                        <div class="flex space-x-1">
                          {#each Array(5) as _, i}
                            <Star
                              class="w-4 h-4 {i < 4
                                ? 'text-neon-yellow fill-current'
                                : 'text-gray-600'}"
                            />
                          {/each}
                        </div>
                        <span class="text-xs text-gray-500">1 week ago</span>
                      </div>
                      <p class="text-gray-300 text-sm leading-relaxed">
                        Great content overall. The React section could use more
                        advanced patterns, but the fundamentals are solid.
                        Completed in about 35 hours.
                      </p>
                      <div
                        class="flex items-center space-x-4 mt-3 text-xs text-gray-500"
                      >
                        <button class="hover:text-neon-cyan transition-colors"
                          >👍 Helpful (8)</button
                        >
                        <button class="hover:text-neon-cyan transition-colors"
                          >Reply</button
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Quick Stats -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-white mb-4">
              Path Stats
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <span class="text-gray-400">Created</span>
                <span class="text-white">{formatDate(pathData.created_at)}</span
                >
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-400">Updated</span>
                <span class="text-white">{formatDate(pathData.updated_at)}</span
                >
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-400">Version</span>
                <span class="text-neon-cyan">{pathData.version}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-400">Views</span>
                <span class="text-white">{pathData.views.toLocaleString()}</span
                >
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-400">Success Rate</span>
                <span class="text-neon-green">{pathData.completion_rate}%</span>
              </div>
            </div>
          </div>

          <!-- Related Paths -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-white mb-4">
              Related Paths
            </h3>
            <div class="space-y-3">
              <div
                class="p-3 bg-dark-700 rounded-lg border border-dark-600 hover:border-neon-purple transition-colors cursor-pointer"
              >
                <h4 class="font-medium text-white text-sm mb-1">
                  Advanced React Patterns
                </h4>
                <p class="text-xs text-gray-400 mb-2">by @react_master</p>
                <div class="flex items-center space-x-2 text-xs text-gray-500">
                  <Star class="w-3 h-3 text-neon-yellow fill-current" />
                  <span>4.8</span>
                  <span>•</span>
                  <span>12h</span>
                </div>
              </div>

              <div
                class="p-3 bg-dark-700 rounded-lg border border-dark-600 hover:border-neon-purple transition-colors cursor-pointer"
              >
                <h4 class="font-medium text-white text-sm mb-1">
                  Node.js Microservices
                </h4>
                <p class="text-xs text-gray-400 mb-2">by @backend_guru</p>
                <div class="flex items-center space-x-2 text-xs text-gray-500">
                  <Star class="w-3 h-3 text-neon-yellow fill-current" />
                  <span>4.7</span>
                  <span>•</span>
                  <span>18h</span>
                </div>
              </div>

              <div
                class="p-3 bg-dark-700 rounded-lg border border-dark-600 hover:border-neon-purple transition-colors cursor-pointer"
              >
                <h4 class="font-medium text-white text-sm mb-1">
                  DevOps Fundamentals
                </h4>
                <p class="text-xs text-gray-400 mb-2">by @devops_alice</p>
                <div class="flex items-center space-x-2 text-xs text-gray-500">
                  <Star class="w-3 h-3 text-neon-yellow fill-current" />
                  <span>4.5</span>
                  <span>•</span>
                  <span>25h</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Author Info -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-white mb-4">
              About the Author
            </h3>
            <div class="flex items-center space-x-3 mb-3">
              <div
                class="w-12 h-12 bg-gradient-to-br from-neon-pink to-neon-purple rounded-full flex items-center justify-center"
              >
                <span class="text-white font-bold">
                  {pathData.author[0].toUpperCase()}
                </span>
              </div>
              <div>
                <div class="font-medium text-white">@{pathData.author}</div>
                <div class="text-sm text-gray-400">Full Stack Developer</div>
              </div>
            </div>
            <p class="text-sm text-gray-400 mb-3">{pathData.author_bio}</p>
            <div class="grid grid-cols-2 gap-3 text-center">
              <div>
                <div class="text-lg font-bold text-neon-cyan">12</div>
                <div class="text-xs text-gray-400">Paths</div>
              </div>
              <div>
                <div class="text-lg font-bold text-neon-pink">1.2k</div>
                <div class="text-xs text-gray-400">Followers</div>
              </div>
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

  /* Smooth progress bar animation */
  .progress-bar {
    transition: width 0.5s ease-in-out;
  }

  /* Hover effects for interactive elements */
  .card:hover {
    transform: translateY(-1px);
  }

  /* Tab transitions */
  .tab-content {
    animation: fadeIn 0.3s ease-in-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
