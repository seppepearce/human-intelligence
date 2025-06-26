<script>
  import { onMount, onDestroy } from "svelte";
  import { liveview, useLiveTopic } from "$lib/stores/liveview.js";
  import {
    Activity,
    Users,
    BookOpen,
    MessageSquare,
    Heart,
    GitBranch,
    Clock,
  } from "lucide-svelte";

  // Real-time activity data - Phoenix LiveView style
  const activityFeed = useLiveTopic("activity_feed", {
    activities: [],
    online_users: 0,
  });
  const onlineUsers = useLiveTopic("presence", { users: [], count: 0 });

  // Connection state
  $: connected = liveview.connected;
  $: connectionState = liveview.connectionState;

  let unsubscribe = [];

  onMount(() => {
    // Subscribe to live topics - Phoenix LiveView style
    unsubscribe.push(
      liveview.subscribe("activity_feed", handleActivityUpdate),
      liveview.subscribe("presence", handlePresenceUpdate),
      liveview.subscribe("node_created", handleNodeCreated),
      liveview.subscribe("path_completed", handlePathCompleted),
      liveview.subscribe("user_joined", handleUserJoined)
    );

    // Join presence to show we're online
    liveview.join("presence", {
      user_id: "current_user", // Replace with actual user ID
      timestamp: Date.now(),
    });
  });

  onDestroy(() => {
    // Clean up subscriptions
    unsubscribe.forEach((fn) => fn());
    liveview.leave("presence");
  });

  // Event handlers - LiveView style
  function handleActivityUpdate({ payload }) {
    console.log("Live activity update:", payload);
  }

  function handlePresenceUpdate({ payload }) {
    console.log("Presence update:", payload);
  }

  function handleNodeCreated({ payload }) {
    // Show real-time notification
    showNotification(`${payload.user} created "${payload.title}"`);
  }

  function handlePathCompleted({ payload }) {
    showNotification(`${payload.user} completed path "${payload.path}"`);
  }

  function handleUserJoined({ payload }) {
    showNotification(`${payload.user} joined the community!`);
  }

  function showNotification(message) {
    // Add visual notification - could integrate with toast system
    console.log("🔔", message);
  }

  // Activity type formatting
  function getActivityIcon(type) {
    switch (type) {
      case "node_created":
        return BookOpen;
      case "path_completed":
        return GitBranch;
      case "comment_added":
        return MessageSquare;
      case "node_liked":
        return Heart;
      case "user_registered":
        return Users;
      default:
        return Activity;
    }
  }

  function getActivityColor(type) {
    switch (type) {
      case "node_created":
        return "text-neon-cyan";
      case "path_completed":
        return "text-neon-green";
      case "comment_added":
        return "text-neon-purple";
      case "node_liked":
        return "text-neon-pink";
      case "user_registered":
        return "text-neon-yellow";
      default:
        return "text-gray-400";
    }
  }

  function formatTimeAgo(timestamp) {
    const now = Date.now();
    const diff = now - timestamp;

    if (diff < 60000) return "just now";
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    return `${Math.floor(diff / 86400000)}d ago`;
  }

  // Mock data for demonstration - in real app this comes from WebSocket
  $: activities = $activityFeed.activities || [
    {
      id: "1",
      type: "node_created",
      user: "alice_learns",
      action: "created node",
      target: "Introduction to Machine Learning",
      timestamp: Date.now() - 120000,
    },
    {
      id: "2",
      type: "path_completed",
      user: "bob_codes",
      action: "completed path",
      target: "Frontend Development Fundamentals",
      timestamp: Date.now() - 300000,
    },
    {
      id: "3",
      type: "comment_added",
      user: "charlie_dev",
      action: "commented on",
      target: "React Hooks Deep Dive",
      timestamp: Date.now() - 600000,
    },
  ];

  $: onlineCount = $onlineUsers.count || 42; // Mock data
</script>

<!-- Live Activity Feed - Phoenix LiveView Style -->
<div class="card space-y-6">
  <!-- Header with connection status -->
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-3">
      <Activity class="w-6 h-6 text-neon-pink" />
      <h2 class="text-xl font-display font-bold text-white">Live Activity</h2>
    </div>

    <!-- Connection indicator -->
    <div class="flex items-center space-x-2">
      <div class="flex items-center space-x-1">
        <Users class="w-4 h-4 text-neon-cyan" />
        <span class="text-sm text-neon-cyan">{onlineCount}</span>
      </div>

      <div class="flex items-center space-x-1">
        <div
          class="w-2 h-2 rounded-full {$connected
            ? 'bg-neon-green animate-pulse'
            : 'bg-red-500'}"
        />
        <span class="text-xs text-gray-400">
          {$connected ? "Live" : "Reconnecting..."}
        </span>
      </div>
    </div>
  </div>

  <!-- Activity Feed -->
  <div class="space-y-3 max-h-96 overflow-y-auto">
    {#each activities as activity (activity.id)}
      <div class="activity-item group">
        <div class="flex items-center space-x-3">
          <!-- Activity Icon -->
          <div class="relative">
            <svelte:component
              this={getActivityIcon(activity.type)}
              class="w-5 h-5 {getActivityColor(activity.type)}"
            />
            {#if activity.timestamp > Date.now() - 60000}
              <div
                class="absolute -top-1 -right-1 w-2 h-2 bg-neon-pink rounded-full animate-ping"
              />
            {/if}
          </div>

          <!-- Activity Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <span
                class="font-medium text-neon-cyan hover:text-neon-pink transition-colors cursor-pointer"
              >
                @{activity.user}
              </span>
              <span class="text-gray-400 text-sm">{activity.action}</span>
              <span class="text-white font-medium truncate">
                "{activity.target}"
              </span>
            </div>

            <div class="flex items-center space-x-2 mt-1">
              <Clock class="w-3 h-3 text-gray-500" />
              <span class="text-xs text-gray-500">
                {formatTimeAgo(activity.timestamp)}
              </span>
            </div>
          </div>

          <!-- Hover Actions -->
          <div class="opacity-0 group-hover:opacity-100 transition-opacity">
            <button class="btn-ghost p-1 text-xs"> View </button>
          </div>
        </div>
      </div>
    {/each}

    {#if activities.length === 0}
      <div class="text-center py-8">
        <Activity class="w-12 h-12 text-gray-600 mx-auto mb-3" />
        <p class="text-gray-500">No recent activity</p>
        <p class="text-sm text-gray-600">Start learning to see live updates!</p>
      </div>
    {/if}
  </div>

  <!-- LiveView-style presence indicator -->
  <div class="border-t border-dark-600 pt-4">
    <div class="flex items-center justify-between text-sm">
      <div class="flex items-center space-x-2 text-gray-400">
        <div class="w-2 h-2 bg-neon-green rounded-full animate-pulse" />
        <span>Real-time updates active</span>
      </div>

      <div class="flex items-center space-x-1 text-gray-500">
        <span>{onlineCount} learners online</span>
      </div>
    </div>
  </div>
</div>

<style>
  .activity-item {
    @apply relative p-3 rounded-lg border border-transparent;
    @apply hover:bg-dark-700 hover:border-dark-600 transition-all duration-200;
  }

  .activity-item::before {
    content: "";
    @apply absolute left-0 top-0 bottom-0 w-0.5 bg-gradient-to-b from-neon-pink to-neon-cyan;
    @apply opacity-0 transition-opacity duration-200;
  }

  .activity-item:hover::before {
    @apply opacity-100;
  }

  /* Custom scrollbar for activity feed */
  :global(.activity-feed::-webkit-scrollbar) {
    width: 4px;
  }

  :global(.activity-feed::-webkit-scrollbar-track) {
    @apply bg-dark-800;
  }

  :global(.activity-feed::-webkit-scrollbar-thumb) {
    @apply bg-neon-pink rounded-full;
  }

  :global(.activity-feed::-webkit-scrollbar-thumb:hover) {
    @apply bg-neon-cyan;
  }
</style>
