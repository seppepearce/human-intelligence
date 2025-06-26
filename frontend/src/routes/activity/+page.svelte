<script>
  import { onMount, onDestroy } from "svelte";
  import {
    Activity,
    Users,
    Wifi,
    WifiOff,
    Bell,
    BellOff,
    BookOpen,
    GitBranch,
    MessageCircle,
    Heart,
    Star,
    Clock,
    Eye,
    Zap,
    Code,
    Play,
    RefreshCw,
  } from "lucide-svelte";

  // Enhanced state
  let message = "Live activity feed operational";
  let connected = false;
  let notifications = true;
  let viewMode = "feed";
  let activities = [];
  let onlineUsers = 23;
  let lastUpdate = new Date();

  // Mock activity data
  const mockActivityTypes = [
    {
      type: "node_created",
      icon: BookOpen,
      color: "text-neon-cyan",
      verb: "created node",
    },
    {
      type: "path_completed",
      icon: GitBranch,
      color: "text-neon-green",
      verb: "completed path",
    },
    {
      type: "comment_added",
      icon: MessageCircle,
      color: "text-neon-purple",
      verb: "commented on",
    },
    { type: "node_liked", icon: Heart, color: "text-neon-pink", verb: "liked" },
    { type: "node_viewed", icon: Eye, color: "text-gray-400", verb: "viewed" },
    {
      type: "code_shared",
      icon: Code,
      color: "text-neon-yellow",
      verb: "shared code for",
    },
    {
      type: "path_starred",
      icon: Star,
      color: "text-neon-yellow",
      verb: "starred path",
    },
  ];

  const mockUsers = [
    "alice_learns",
    "bob_codes",
    "charlie_dev",
    "diana_data",
    "eve_ai",
    "frank_web",
    "grace_mobile",
    "henry_backend",
    "iris_design",
    "jack_ops",
  ];

  const mockTargets = [
    "React Hooks Deep Dive",
    "Python Data Science",
    "Node.js Fundamentals",
    "Machine Learning Basics",
    "CSS Grid Mastery",
    "JavaScript Async Patterns",
    "Docker Containerization",
    "API Design Principles",
    "Database Optimization",
    "Frontend Performance",
    "Web Security Essentials",
    "GraphQL Introduction",
  ];

  function generateMockActivity() {
    const activityType =
      mockActivityTypes[Math.floor(Math.random() * mockActivityTypes.length)];
    const user = mockUsers[Math.floor(Math.random() * mockUsers.length)];
    const target = mockTargets[Math.floor(Math.random() * mockTargets.length)];

    return {
      id: `activity_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type: activityType.type,
      icon: activityType.icon,
      color: activityType.color,
      user: user,
      verb: activityType.verb,
      target: target,
      timestamp: Date.now() - Math.floor(Math.random() * 3600000), // Random time within last hour
      isNew: Math.random() > 0.7, // 30% chance of being "new"
    };
  }

  function formatTimeAgo(timestamp) {
    const now = Date.now();
    const diff = now - timestamp;

    if (diff < 60000) return "just now";
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    return `${Math.floor(diff / 86400000)}d ago`;
  }

  function refreshActivities() {
    activities = Array.from({ length: 8 }, () => generateMockActivity()).sort(
      (a, b) => b.timestamp - a.timestamp
    );
    lastUpdate = new Date();
  }

  function addLiveActivity() {
    const newActivity = generateMockActivity();
    newActivity.timestamp = Date.now();
    newActivity.isNew = true;
    activities = [newActivity, ...activities.slice(0, 7)];

    // Remove "new" flag after 3 seconds
    setTimeout(() => {
      activities = activities.map((a) =>
        a.id === newActivity.id ? { ...a, isNew: false } : a
      );
    }, 3000);
  }

  onMount(() => {
    console.log("Enhanced activity page mounted");
    connected = true;
    refreshActivities();

    // Simulate live updates every 8-15 seconds
    const interval = setInterval(() => {
      if (connected && Math.random() > 0.3) {
        addLiveActivity();
      }
    }, 8000 + Math.random() * 7000);

    return () => clearInterval(interval);
  });
</script>

<svelte:head>
  <title>Live Activity - Human Intelligence</title>
  <meta
    name="description"
    content="Real-time learning activity across the Human Intelligence platform"
  />
</svelte:head>

<div class="min-h-screen bg-dark-900">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-20" />
  <div
    class="absolute inset-0 bg-gradient-to-br from-neon-pink/5 via-transparent to-neon-cyan/5"
  />

  <div class="relative container-wide py-8">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <Activity class="w-8 h-8 text-neon-pink animate-pulse" />
          <div>
            <h1 class="text-3xl font-display font-bold text-white">
              Live Activity
            </h1>
            <p class="text-gray-400">Real-time learning across the platform</p>
          </div>
        </div>

        <!-- Connection Status -->
        <div
          class="flex items-center space-x-2 px-3 py-1 rounded-full bg-dark-800 border border-dark-600"
        >
          {#if connected}
            <Wifi class="w-4 h-4 text-neon-green" />
            <span class="text-sm text-neon-green">Live</span>
          {:else}
            <WifiOff class="w-4 h-4 text-red-400" />
            <span class="text-sm text-red-400">Connecting...</span>
          {/if}
          <span class="text-xs text-gray-500 ml-2">23 online</span>
        </div>
      </div>

      <!-- View Toggle -->
      <div
        class="flex bg-dark-800 rounded-lg p-1 border border-dark-600 mt-4 w-fit"
      >
        <button
          on:click={() => (viewMode = "feed")}
          class="px-3 py-1 rounded text-sm transition-all {viewMode === 'feed'
            ? 'bg-neon-pink text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Feed
        </button>
        <button
          on:click={() => (viewMode = "stats")}
          class="px-3 py-1 rounded text-sm transition-all {viewMode === 'stats'
            ? 'bg-neon-pink text-white'
            : 'text-gray-400 hover:text-white'}"
        >
          Stats
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <div class="lg:col-span-3">
        {#if viewMode === "feed"}
          <div class="card">
            <div class="flex items-center justify-between mb-6">
              <h2
                class="text-xl font-display font-bold text-white flex items-center"
              >
                <Activity class="w-6 h-6 mr-3 text-neon-pink" />
                Live Activity Feed
              </h2>
              <div class="flex items-center space-x-3">
                <button
                  on:click={refreshActivities}
                  class="btn-ghost p-2 hover:text-neon-cyan transition-colors"
                  title="Refresh activities"
                >
                  <RefreshCw class="w-4 h-4" />
                </button>
                <div class="text-xs text-gray-500">
                  Updated: {lastUpdate.toLocaleTimeString()}
                </div>
              </div>
            </div>

            <div class="space-y-3 max-h-96 overflow-y-auto">
              {#each activities as activity (activity.id)}
                <div
                  class="activity-item group p-4 rounded-lg bg-dark-700/50 border border-dark-600 hover:border-neon-pink/30 transition-all duration-300 {activity.isNew
                    ? 'ring-2 ring-neon-pink/50'
                    : ''}"
                >
                  <div class="flex items-start space-x-3">
                    <div class="relative">
                      <svelte:component
                        this={activity.icon}
                        class="w-5 h-5 {activity.color}"
                      />
                      {#if activity.isNew}
                        <div
                          class="absolute -top-1 -right-1 w-2 h-2 bg-neon-pink rounded-full animate-ping"
                        />
                      {/if}
                    </div>

                    <div class="flex-1 min-w-0">
                      <div class="flex items-center space-x-2 mb-1">
                        <span
                          class="font-medium text-neon-cyan hover:text-neon-pink cursor-pointer transition-colors"
                        >
                          @{activity.user}
                        </span>
                        <span class="text-gray-400 text-sm"
                          >{activity.verb}</span
                        >
                        <span class="text-white font-medium truncate">
                          "{activity.target}"
                        </span>
                      </div>

                      <div
                        class="flex items-center space-x-2 text-xs text-gray-500"
                      >
                        <Clock class="w-3 h-3" />
                        <span>{formatTimeAgo(activity.timestamp)}</span>
                        {#if activity.isNew}
                          <span class="text-neon-pink font-medium">• NEW</span>
                        {/if}
                      </div>
                    </div>

                    <div
                      class="opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <button
                        class="btn-ghost p-1 text-xs hover:text-neon-cyan"
                      >
                        <Eye class="w-3 h-3" />
                      </button>
                    </div>
                  </div>
                </div>
              {/each}

              {#if activities.length === 0}
                <div class="text-center py-12">
                  <Activity
                    class="w-16 h-16 text-gray-600 mx-auto mb-4 animate-pulse"
                  />
                  <p class="text-gray-500 mb-2">No recent activity</p>
                  <button
                    on:click={refreshActivities}
                    class="btn-secondary text-sm flex items-center mx-auto"
                  >
                    <RefreshCw class="w-4 h-4 mr-2" />
                    Load Activities
                  </button>
                </div>
              {/if}
            </div>

            <!-- Activity Stats -->
            <div class="border-t border-dark-600 pt-4 mt-6">
              <div class="grid grid-cols-3 gap-4 text-center text-sm">
                <div>
                  <div class="text-neon-cyan font-semibold">
                    {activities.length}
                  </div>
                  <div class="text-gray-500">Recent</div>
                </div>
                <div>
                  <div class="text-neon-green font-semibold">{onlineUsers}</div>
                  <div class="text-gray-500">Online</div>
                </div>
                <div>
                  <div class="text-neon-purple font-semibold">
                    {activities.filter((a) => a.isNew).length}
                  </div>
                  <div class="text-gray-500">Live</div>
                </div>
              </div>
            </div>
          </div>
        {:else}
          <div class="card">
            <h2
              class="text-xl font-display font-bold text-white mb-6 flex items-center"
            >
              <Zap class="w-6 h-6 mr-3 text-neon-yellow" />
              Live Statistics
            </h2>
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-dark-700 p-4 rounded-lg border border-dark-600">
                <div class="flex items-center justify-between mb-2">
                  <BookOpen class="w-5 h-5 text-neon-cyan" />
                  <span class="text-2xl font-bold text-neon-cyan">42</span>
                </div>
                <div class="text-sm text-gray-400">Nodes Created Today</div>
              </div>
              <div class="bg-dark-700 p-4 rounded-lg border border-dark-600">
                <div class="flex items-center justify-between mb-2">
                  <GitBranch class="w-5 h-5 text-neon-green" />
                  <span class="text-2xl font-bold text-neon-green">15</span>
                </div>
                <div class="text-sm text-gray-400">Paths Completed</div>
              </div>
              <div class="bg-dark-700 p-4 rounded-lg border border-dark-600">
                <div class="flex items-center justify-between mb-2">
                  <Heart class="w-5 h-5 text-neon-pink" />
                  <span class="text-2xl font-bold text-neon-pink">128</span>
                </div>
                <div class="text-sm text-gray-400">Likes Given</div>
              </div>
              <div class="bg-dark-700 p-4 rounded-lg border border-dark-600">
                <div class="flex items-center justify-between mb-2">
                  <MessageCircle class="w-5 h-5 text-neon-purple" />
                  <span class="text-2xl font-bold text-neon-purple">67</span>
                </div>
                <div class="text-sm text-gray-400">Comments Posted</div>
              </div>
            </div>
          </div>
        {/if}
      </div>

      <div class="space-y-6">
        <div class="card">
          <h3 class="text-lg font-display font-bold text-white mb-4">
            Live Controls
          </h3>
          <div class="space-y-3">
            <button
              on:click={() => (notifications = !notifications)}
              class="btn-ghost w-full justify-start {notifications
                ? 'text-neon-cyan'
                : 'text-gray-500'} flex items-center"
            >
              {#if notifications}
                <Bell class="w-4 h-4 mr-2" />
                Notifications On
              {:else}
                <BellOff class="w-4 h-4 mr-2" />
                Notifications Off
              {/if}
            </button>

            <button
              on:click={addLiveActivity}
              class="btn-secondary w-full justify-start flex items-center"
            >
              <Play class="w-4 h-4 mr-2" />
              Simulate Activity
            </button>

            <button
              on:click={refreshActivities}
              class="btn-ghost w-full justify-start flex items-center"
            >
              <RefreshCw class="w-4 h-4 mr-2" />
              Refresh Feed
            </button>
          </div>
        </div>

        <div class="card">
          <h3 class="text-lg font-display font-bold text-neon-green mb-4">
            Activity Types
          </h3>
          <div class="space-y-2 text-sm">
            {#each mockActivityTypes.slice(0, 4) as type}
              <div class="flex items-center space-x-2">
                <svelte:component
                  this={type.icon}
                  class="w-4 h-4 {type.color}"
                />
                <span class="text-gray-300 capitalize"
                  >{type.type.replace("_", " ")}</span
                >
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
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

  .activity-item {
    position: relative;
  }

  .activity-item::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 2px;
    background: linear-gradient(
      to bottom,
      theme("colors.neon.pink"),
      theme("colors.neon.cyan")
    );
    opacity: 0;
    transition: opacity 0.3s;
  }

  .activity-item:hover::before {
    opacity: 1;
  }

  /* Custom scrollbar */
  .activity-item::-webkit-scrollbar {
    width: 4px;
  }

  .activity-item::-webkit-scrollbar-track {
    background: theme("colors.dark.800");
  }

  .activity-item::-webkit-scrollbar-thumb {
    background: theme("colors.neon.pink");
    border-radius: 2px;
  }
</style>
