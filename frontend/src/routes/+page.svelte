<script>
  import { onMount, onDestroy } from "svelte";
  import { writable } from "svelte/store";
  import { fade, fly, scale } from "svelte/transition";
  import { browser } from "$app/environment";

  // Icons
  import {
    Activity,
    Users,
    BookOpen,
    Zap,
    TrendingUp,
    Clock,
    Star,
    MessageCircle,
    Plus,
    ArrowRight,
    Brain,
    Sparkles,
    Target,
    Award,
    GitBranch,
    Eye,
  } from "lucide-svelte";

  // Stores for real-time data
  let liveActivity = writable([]);
  let popularPaths = writable([]);
  let trendingNodes = writable([]);
  let communityStats = writable({
    totalUsers: 0,
    activePaths: 0,
    completedNodes: 0,
    onlineUsers: 0,
  });

  // Component state
  let ws = null;
  let mounted = false;
  let selectedCategory = "all";
  let viewMode = "grid"; // 'grid' or 'tree'

  // Mock data for initial render
  let mockActivity = [
    {
      id: 1,
      type: "node_completed",
      user: { username: "alice_learns", avatar: "/avatars/alice.png" },
      node: { title: "React Hooks Deep Dive", type: "code" },
      path: { name: "Modern React Development" },
      timestamp: new Date(Date.now() - 120000),
      intensity: "high",
    },
    {
      id: 2,
      type: "path_started",
      user: { username: "bob_codes", avatar: "/avatars/bob.png" },
      path: { name: "Machine Learning Fundamentals", nodes: 12 },
      timestamp: new Date(Date.now() - 300000),
      intensity: "medium",
    },
    {
      id: 3,
      type: "tldr_created",
      user: { username: "charlie_ai", avatar: "/avatars/charlie.png" },
      tldr: "Neural networks are just fancy function approximators that learn from data!",
      path: { name: "AI for Beginners" },
      timestamp: new Date(Date.now() - 480000),
      intensity: "low",
    },
    {
      id: 4,
      type: "path_forked",
      user: { username: "diana_dev", avatar: "/avatars/diana.png" },
      path: { name: "Full Stack JavaScript", parent: "Web Development Basics" },
      timestamp: new Date(Date.now() - 600000),
      intensity: "medium",
    },
    {
      id: 5,
      type: "collaboration",
      users: [
        { username: "eve_teacher", avatar: "/avatars/eve.png" },
        { username: "frank_student", avatar: "/avatars/frank.png" },
      ],
      path: { name: "Calculus Made Simple" },
      timestamp: new Date(Date.now() - 720000),
      intensity: "high",
    },
  ];

  let mockPaths = [
    {
      id: 1,
      name: "Modern Web Development",
      description: "Learn React, Node.js, and modern web technologies",
      nodes: 24,
      completions: 1337,
      forks: 89,
      tags: ["javascript", "react", "nodejs"],
      difficulty: "intermediate",
      lastActivity: new Date(Date.now() - 60000),
      creator: { username: "webmaster_pro", avatar: "/avatars/webmaster.png" },
      isLive: true,
    },
    {
      id: 2,
      name: "AI & Machine Learning",
      description: "From basics to advanced neural networks",
      nodes: 42,
      completions: 2156,
      forks: 156,
      tags: ["ai", "python", "tensorflow"],
      difficulty: "advanced",
      lastActivity: new Date(Date.now() - 180000),
      creator: { username: "ai_researcher", avatar: "/avatars/researcher.png" },
      isLive: true,
    },
    {
      id: 3,
      name: "Blockchain Fundamentals",
      description: "Understand cryptocurrency and DeFi",
      nodes: 18,
      completions: 892,
      forks: 67,
      tags: ["blockchain", "solidity", "defi"],
      difficulty: "beginner",
      lastActivity: new Date(Date.now() - 300000),
      creator: { username: "crypto_guru", avatar: "/avatars/crypto.png" },
      isLive: false,
    },
  ];

  let mockStats = {
    totalUsers: 12547,
    activePaths: 1893,
    completedNodes: 45621,
    onlineUsers: 234,
  };

  onMount(() => {
    mounted = true;

    // Initialize stores with mock data
    liveActivity.set(mockActivity);
    popularPaths.set(mockPaths);
    communityStats.set(mockStats);

    // Connect to WebSocket for real-time updates
    connectWebSocket();

    // Set up periodic updates
    const interval = setInterval(updateMockData, 3000);

    return () => {
      clearInterval(interval);
      if (ws) ws.close();
    };
  });

  function connectWebSocket() {
    if (!browser) return;

    try {
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
      const wsUrl = `${protocol}//${window.location.host}/ws`;

      ws = new WebSocket(wsUrl);

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleRealtimeUpdate(data);
      };
    } catch (error) {
      console.log("WebSocket connection failed, using mock data");
    }
  }

  function handleRealtimeUpdate(data) {
    switch (data.type) {
      case "activity":
        liveActivity.update((activities) => [
          data.activity,
          ...activities.slice(0, 19),
        ]);
        break;
      case "stats":
        communityStats.set(data.stats);
        break;
    }
  }

  function updateMockData() {
    // Simulate new activity
    const newActivity = {
      id: Date.now(),
      type: ["node_completed", "path_started", "tldr_created", "path_forked"][
        Math.floor(Math.random() * 4)
      ],
      user: {
        username: [
          "learner" + Math.floor(Math.random() * 1000),
          "student" + Math.floor(Math.random() * 1000),
        ][Math.floor(Math.random() * 2)],
        avatar: "/avatars/default.png",
      },
      timestamp: new Date(),
      intensity: ["low", "medium", "high"][Math.floor(Math.random() * 3)],
    };

    if (Math.random() > 0.7) {
      liveActivity.update((activities) => [
        newActivity,
        ...activities.slice(0, 19),
      ]);
    }

    // Update stats
    communityStats.update((stats) => ({
      ...stats,
      onlineUsers: stats.onlineUsers + Math.floor(Math.random() * 3) - 1,
    }));
  }

  function formatTimeAgo(date) {
    const now = new Date();
    const diff = now - date;
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return "just now";
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    return `${days}d ago`;
  }

  function getActivityIcon(type) {
    switch (type) {
      case "node_completed":
        return Zap;
      case "path_started":
        return BookOpen;
      case "tldr_created":
        return MessageCircle;
      case "path_forked":
        return GitBranch;
      case "collaboration":
        return Users;
      default:
        return Activity;
    }
  }

  function getActivityColor(type, intensity = "medium") {
    const colors = {
      node_completed:
        intensity === "high" ? "text-neon-green" : "text-terminal-400",
      path_started: intensity === "high" ? "text-neon-blue" : "text-blue-400",
      tldr_created:
        intensity === "high" ? "text-neon-purple" : "text-purple-400",
      path_forked: intensity === "high" ? "text-neon-cyan" : "text-cyan-400",
      collaboration: intensity === "high" ? "text-neon-pink" : "text-pink-400",
    };
    return colors[type] || "text-gray-400";
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
</script>

<svelte:head>
  <title>Human Intelligence - Where AI Augments HI</title>
  <meta
    name="description"
    content="Join the collaborative learning revolution. Share knowledge, build learning paths, and grow together with AI-augmented human intelligence."
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Hero Section -->
  <section
    class="relative overflow-hidden bg-gradient-to-br from-dark-900 via-dark-800 to-dark-900"
  >
    <div class="absolute inset-0 bg-grid opacity-30" />
    <div
      class="absolute inset-0 bg-gradient-to-r from-neon-pink/10 via-transparent to-neon-cyan/10"
    />

    <div class="relative container-wide py-20">
      <div class="text-center space-y-8">
        <div class="flex justify-center items-center space-x-4 mb-8">
          <Brain class="w-16 h-16 text-neon-pink animate-pulse" />
          <div
            class="text-8xl font-display font-black bg-gradient-to-r from-neon-pink via-neon-purple to-neon-cyan bg-clip-text text-transparent"
          >
            HI!
          </div>
          <Sparkles class="w-12 h-12 text-neon-cyan animate-float" />
        </div>

        <h1 class="text-6xl md:text-8xl font-display font-black">
          <span
            class="bg-gradient-to-r from-neon-pink to-neon-cyan bg-clip-text text-transparent"
          >
            Human Intelligence
          </span>
        </h1>

        <p class="text-2xl md:text-3xl text-gray-300 max-w-4xl mx-auto">
          Where <span class="text-neon-cyan font-bold">AI augments HI</span>. A
          factorial multiplier to human intelligence.
        </p>

        <p class="text-lg text-gray-400 max-w-2xl mx-auto">
          Join the collaborative learning revolution. Share knowledge, build
          learning paths, and grow together in the most vibrant learning
          community on the web.
        </p>

        <div
          class="flex flex-col sm:flex-row gap-4 justify-center items-center pt-8"
        >
          <a
            href="/register"
            class="btn-primary text-lg px-8 py-4 flex items-center"
          >
            <Plus class="w-5 h-5 mr-2" />
            Join the Revolution
          </a>
          <a
            href="/paths"
            class="btn-secondary text-lg px-8 py-4 flex items-center"
          >
            <BookOpen class="w-5 h-5 mr-2" />
            Explore Paths
          </a>
        </div>
      </div>
    </div>
  </section>

  <!-- Stats Bar -->
  <section class="bg-dark-800 border-y border-dark-600">
    <div class="container-wide py-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8">
        {#each Object.entries($communityStats) as [key, value]}
          <div class="text-center">
            <div class="text-3xl font-bold text-neon-cyan font-mono">
              {value.toLocaleString()}
            </div>
            <div class="text-sm text-gray-400 uppercase tracking-wider">
              {key.replace(/([A-Z])/g, " $1").trim()}
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Main Content -->
  <div class="container-wide py-12">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Live Activity Feed -->
      <div class="lg:col-span-2 space-y-6">
        <div class="flex items-center justify-between">
          <h2
            class="text-3xl font-display font-bold text-neon-cyan flex items-center"
          >
            <Activity class="w-8 h-8 mr-3 animate-pulse" />
            Live Forum Activity
          </h2>
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-neon-green rounded-full animate-pulse" />
            <span class="text-sm text-gray-400 font-mono">LIVE</span>
          </div>
        </div>

        <div class="space-y-4">
          {#each $liveActivity as activity (activity.id)}
            <div
              class="activity-item group hover:bg-dark-700 rounded-lg p-4 border border-dark-600 hover:border-{getActivityColor(
                activity.type,
                activity.intensity
              ).split('-')[1]}-500 transition-all duration-300"
              in:fly={{ y: -20, duration: 300 }}
            >
              <div class="flex items-start space-x-4">
                <div class="flex-shrink-0">
                  <div
                    class="w-3 h-3 rounded-full bg-{getActivityColor(
                      activity.type,
                      activity.intensity
                    ).split('-')[1]}-500 animate-pulse"
                  />
                </div>

                <div class="flex-1 min-w-0">
                  <div class="flex items-center space-x-2 mb-2">
                    <svelte:component
                      this={getActivityIcon(activity.type)}
                      class="w-4 h-4 {getActivityColor(
                        activity.type,
                        activity.intensity
                      )}"
                    />
                    <span class="font-medium text-white">
                      {activity.user?.username || activity.users?.[0]?.username}
                    </span>
                    {#if activity.users && activity.users.length > 1}
                      <span class="text-gray-400"
                        >+{activity.users.length - 1} others</span
                      >
                    {/if}
                    <span class="text-gray-500 text-sm">
                      {formatTimeAgo(activity.timestamp)}
                    </span>
                  </div>

                  <div class="text-gray-300">
                    {#if activity.type === "node_completed"}
                      completed <span class="text-neon-green font-medium"
                        >{activity.node?.title || "Unknown Node"}</span
                      >
                      {#if activity.path?.name}
                        in <span class="text-neon-cyan"
                          >{activity.path.name}</span
                        >
                      {/if}
                    {:else if activity.type === "path_started"}
                      started learning <span class="text-neon-blue font-medium"
                        >{activity.path?.name || "Unknown Path"}</span
                      >
                      {#if activity.path?.nodes}
                        <span class="text-gray-500"
                          >({activity.path.nodes} nodes)</span
                        >
                      {/if}
                    {:else if activity.type === "tldr_created"}
                      shared a TLDR: <span class="text-neon-purple italic"
                        >"{activity.tldr || "No content"}"</span
                      >
                    {:else if activity.type === "path_forked"}
                      forked <span class="text-neon-cyan font-medium"
                        >{activity.path?.name || "Unknown Path"}</span
                      >
                      {#if activity.path?.parent}
                        from <span class="text-gray-400"
                          >{activity.path.parent}</span
                        >
                      {/if}
                    {:else if activity.type === "collaboration"}
                      are collaborating on <span
                        class="text-neon-pink font-medium"
                        >{activity.path?.name || "Unknown Path"}</span
                      >
                    {/if}
                  </div>
                </div>

                <div class="flex-shrink-0">
                  <ArrowRight
                    class="w-4 h-4 text-gray-600 group-hover:text-neon-cyan transition-colors"
                  />
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-8">
        <!-- Popular Learning Paths -->
        <div class="card">
          <h3
            class="text-xl font-display font-bold text-neon-purple mb-6 flex items-center"
          >
            <TrendingUp class="w-5 h-5 mr-2" />
            Trending Paths
          </h3>

          <div class="space-y-4">
            {#each $popularPaths as path}
              <div class="group cursor-pointer">
                <div
                  class="flex items-start space-x-3 p-3 rounded-lg hover:bg-dark-700 transition-colors"
                >
                  <div class="flex-shrink-0 mt-1">
                    {#if path.isLive}
                      <div
                        class="w-3 h-3 bg-neon-green rounded-full animate-pulse"
                      />
                    {:else}
                      <div class="w-3 h-3 bg-gray-600 rounded-full" />
                    {/if}
                  </div>

                  <div class="flex-1 min-w-0">
                    <h4
                      class="font-semibold text-white group-hover:text-neon-cyan transition-colors"
                    >
                      {path.name}
                    </h4>
                    <p class="text-sm text-gray-400 mt-1">
                      {path.description}
                    </p>

                    <div
                      class="flex items-center space-x-4 mt-2 text-xs text-gray-500"
                    >
                      <span class="flex items-center">
                        <BookOpen class="w-3 h-3 mr-1" />
                        {path.nodes} nodes
                      </span>
                      <span class="flex items-center">
                        <Users class="w-3 h-3 mr-1" />
                        {path.completions}
                      </span>
                      <span class="flex items-center">
                        <GitBranch class="w-3 h-3 mr-1" />
                        {path.forks}
                      </span>
                    </div>

                    <div class="flex items-center justify-between mt-2">
                      <div class="flex flex-wrap gap-1">
                        {#each path.tags.slice(0, 2) as tag}
                          <span
                            class="px-2 py-1 bg-dark-700 text-xs rounded text-gray-300"
                          >
                            {tag}
                          </span>
                        {/each}
                      </div>
                      <span
                        class="text-xs {getDifficultyColor(path.difficulty)}"
                      >
                        {path.difficulty}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            {/each}
          </div>

          <div class="pt-4 border-t border-dark-600">
            <a href="/paths" class="btn-ghost w-full justify-center">
              View All Paths
              <ArrowRight class="w-4 h-4 ml-2" />
            </a>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="card">
          <h3 class="text-xl font-display font-bold text-neon-green mb-6">
            Quick Start
          </h3>

          <div class="space-y-3">
            <a
              href="/nodes/create"
              class="btn-terminal w-full justify-start flex items-center"
            >
              <Plus class="w-4 h-4 mr-2" />
              Create Node
            </a>
            <a
              href="/paths/create"
              class="btn-secondary w-full justify-start flex items-center"
            >
              <BookOpen class="w-4 h-4 mr-2" />
              Start Path
            </a>
            <a
              href="/search"
              class="btn-ghost w-full justify-start flex items-center"
            >
              <Target class="w-4 h-4 mr-2" />
              Find Topic
            </a>
          </div>
        </div>

        <!-- Community Highlights -->
        <div class="card">
          <h3
            class="text-xl font-display font-bold text-neon-yellow mb-6 flex items-center"
          >
            <Award class="w-5 h-5 mr-2" />
            This Week
          </h3>

          <div class="space-y-4 text-sm">
            <div class="flex items-center justify-between">
              <span class="text-gray-300">Most Active Learner</span>
              <span class="text-neon-cyan font-medium">@coding_ninja</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-300">Top Contributor</span>
              <span class="text-neon-purple font-medium">@teach_master</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-300">Rising Path</span>
              <span class="text-neon-green font-medium">Web3 Basics</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-300">Nodes Created</span>
              <span class="text-neon-pink font-medium">847</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Call to Action -->
  <section
    class="bg-gradient-to-r from-neon-pink/10 to-neon-cyan/10 border-t border-dark-600"
  >
    <div class="container-wide py-16 text-center">
      <h2 class="text-4xl font-display font-bold mb-4">
        Ready to <span class="text-neon-pink">multiply</span> your intelligence?
      </h2>
      <p class="text-lg text-gray-300 mb-8 max-w-2xl mx-auto">
        Join thousands of learners already sharing knowledge and building the
        future of education.
      </p>
      <div class="flex flex-col sm:flex-row gap-4 justify-center">
        <a
          href="/register"
          class="btn-primary text-lg px-8 py-4 flex items-center"
        >
          Get Started Free
        </a>
        <a href="/about" class="btn-ghost text-lg px-8 py-4 flex items-center">
          Learn More
        </a>
      </div>
    </div>
  </section>
</div>

<style>
  .animate-float {
    animation: float 3s ease-in-out infinite;
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }
</style>
