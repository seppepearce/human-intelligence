<script>
  import "../app.css";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { browser } from "$app/environment";
  import { writable } from "svelte/store";
  import { toast } from "svelte-french-toast";

  // Icons
  import {
    Search,
    User,
    Settings,
    Home,
    BookOpen,
    Zap,
    Users,
    Activity,
    Menu,
    X,
    Github,
    Terminal,
    Brain,
    Sparkles,
  } from "lucide-svelte";

  // Stores
  export const data = {};

  let isMenuOpen = false;
  let isSearchFocused = false;
  let searchQuery = "";
  let isOnline = true;
  let activityCount = 0;
  let currentUser = null;

  // Real-time connection status
  let wsConnected = false;
  let wsReconnectAttempts = 0;
  const maxReconnectAttempts = 5;

  // Navigation items
  const navItems = [
    { href: "/", icon: Home, label: "Forum", badge: null },
    { href: "/paths", icon: BookOpen, label: "Paths", badge: null },
    { href: "/nodes", icon: Zap, label: "Nodes", badge: null },
    { href: "/community", icon: Users, label: "Community", badge: null },
    { href: "/activity", icon: Activity, label: "Live", badge: activityCount },
  ];

  // WebSocket connection for real-time features
  let ws = null;

  onMount(() => {
    // Initialize theme
    initializeTheme();

    // Set up online/offline detection
    isOnline = navigator.onLine;
    window.addEventListener("online", () => (isOnline = true));
    window.addEventListener("offline", () => (isOnline = false));

    // Connect to WebSocket for real-time features
    connectWebSocket();

    // Check authentication status
    checkAuthStatus();

    // Setup keyboard shortcuts
    setupKeyboardShortcuts();

    return () => {
      if (ws) {
        ws.close();
      }
    };
  });

  function initializeTheme() {
    if (browser) {
      const stored = localStorage.getItem("theme");
      const prefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      const theme = stored || (prefersDark ? "dark" : "light");

      document.documentElement.classList.toggle("dark", theme === "dark");
      localStorage.setItem("theme", theme);
    }
  }

  function connectWebSocket() {
    if (!browser) return;

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    try {
      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        wsConnected = true;
        wsReconnectAttempts = 0;
        console.log("🔗 WebSocket connected");
        toast.success("Connected to live feed");
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
      };

      ws.onclose = () => {
        wsConnected = false;
        console.log("🔌 WebSocket disconnected");

        // Attempt to reconnect
        if (wsReconnectAttempts < maxReconnectAttempts) {
          wsReconnectAttempts++;
          setTimeout(connectWebSocket, 2000 * wsReconnectAttempts);
        } else {
          toast.error("Lost connection to live feed");
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    } catch (error) {
      console.error("Failed to connect WebSocket:", error);
    }
  }

  function handleWebSocketMessage(data) {
    switch (data.type) {
      case "activity":
        activityCount++;
        // Update activity feed
        break;
      case "heartbeat":
        // Keep connection alive
        break;
      case "notification":
        toast(data.message, {
          icon: data.icon || "🔔",
          style:
            "background: #1a1a1a; color: #ff006e; border: 1px solid #3a3a3a;",
        });
        break;
    }
  }

  async function checkAuthStatus() {
    try {
      const response = await fetch("/api/v1/users/me", {
        credentials: "include",
      });

      if (response.ok) {
        currentUser = await response.json();
      }
    } catch (error) {
      console.log("Not authenticated");
    }
  }

  function setupKeyboardShortcuts() {
    if (!browser) return;

    document.addEventListener("keydown", (e) => {
      // Ctrl/Cmd + K for search
      if ((e.ctrlKey || e.metaKey) && e.key === "k") {
        e.preventDefault();
        focusSearch();
      }

      // Escape to close menu/search
      if (e.key === "Escape") {
        isMenuOpen = false;
        isSearchFocused = false;
      }
    });
  }

  function focusSearch() {
    const searchInput = document.querySelector("#main-search");
    if (searchInput) {
      searchInput.focus();
      isSearchFocused = true;
    }
  }

  function handleSearch(e) {
    if (e.key === "Enter" && searchQuery.trim()) {
      goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  }

  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }

  function closeMenu() {
    isMenuOpen = false;
  }

  // Check if current route is active
  function isActiveRoute(href) {
    const pathname = $page.url.pathname;
    if (href === "/") {
      return pathname === "/";
    }
    return pathname.startsWith(href);
  }

  // Format activity count for display
  function formatActivityCount(count) {
    if (count === 0) return "";
    if (count > 99) return "99+";
    return count.toString();
  }

  // Compute page title reactively
  $: pageTitle = $page.data?.title
    ? `${$page.data.title} - Human Intelligence`
    : "Human Intelligence - HI!";
</script>

<svelte:head>
  <title>{pageTitle}</title>

  {#if $page.data?.description}
    <meta name="description" content={$page.data.description} />
  {/if}
</svelte:head>

<!-- Main App Container -->
<div class="min-h-screen bg-dark-900 text-white">
  <!-- Navigation Header -->
  <header
    class="sticky top-0 z-50 bg-dark-800 bg-opacity-95 backdrop-blur-sm border-b border-dark-600"
  >
    <div class="container-wide">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <div class="flex items-center space-x-4">
          <a href="/" class="flex items-center space-x-2 group">
            <div class="relative">
              <Brain class="w-8 h-8 text-neon-pink group-hover:animate-pulse" />
              <Sparkles
                class="w-4 h-4 text-neon-cyan absolute -top-1 -right-1 animate-pulse"
              />
            </div>
            <div
              class="font-display font-black text-2xl bg-gradient-to-r from-neon-pink to-neon-cyan bg-clip-text text-transparent"
            >
              HI!
            </div>
          </a>

          <!-- Online Status Indicator -->
          <div class="hidden md:flex items-center space-x-2">
            <div
              class="w-2 h-2 rounded-full {isOnline
                ? 'bg-neon-green animate-pulse'
                : 'bg-red-500'}"
            />
            <span class="text-xs text-gray-400 font-mono">
              {isOnline ? "ONLINE" : "OFFLINE"}
            </span>
            {#if wsConnected}
              <div
                class="w-2 h-2 rounded-full bg-neon-cyan animate-pulse"
                title="Live feed connected"
              />
            {/if}
          </div>
        </div>

        <!-- Desktop Navigation -->
        <nav class="hidden md:flex items-center space-x-1">
          {#each navItems as item}
            {@const isActive =
              item.href === "/"
                ? $page.url.pathname === "/"
                : $page.url.pathname.startsWith(item.href)}
            <a
              href={item.href}
              class="nav-link {isActive ? 'active' : ''}"
              on:click={closeMenu}
            >
              <div class="flex items-center space-x-2">
                <svelte:component this={item.icon} class="w-4 h-4" />
                <span>{item.label}</span>
                {#if item.badge && item.badge > 0}
                  <span
                    class="inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-dark-900 bg-neon-pink rounded-full min-w-[20px]"
                  >
                    {formatActivityCount(item.badge)}
                  </span>
                {/if}
              </div>
            </a>
          {/each}
        </nav>

        <!-- Search Bar -->
        <div class="hidden md:flex flex-1 max-w-md mx-8">
          <div class="relative w-full">
            <Search
              class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500"
            />
            <input
              id="main-search"
              type="text"
              placeholder="Search nodes, paths, users... (⌘K)"
              class="search-input pl-10 pr-4"
              bind:value={searchQuery}
              on:keydown={handleSearch}
              on:focus={() => (isSearchFocused = true)}
              on:blur={() => (isSearchFocused = false)}
            />
            {#if isSearchFocused && searchQuery}
              <div class="search-results">
                <div class="p-4 text-center text-gray-500 font-mono">
                  Press Enter to search for "{searchQuery}"
                </div>
              </div>
            {/if}
          </div>
        </div>

        <!-- User Menu / Auth -->
        <div class="flex items-center space-x-4">
          {#if currentUser}
            <div class="flex items-center space-x-2">
              <img
                src={currentUser.avatar || "/default-avatar.png"}
                alt={currentUser.username}
                class="user-avatar"
              />
              <span class="hidden md:block text-sm font-medium"
                >{currentUser.username}</span
              >
            </div>
          {:else}
            <div class="hidden md:flex items-center space-x-2">
              <a href="/login" class="btn-secondary text-sm">Login</a>
              <a href="/register" class="btn-primary text-sm">Join HI!</a>
            </div>
          {/if}

          <!-- Settings -->
          <button class="p-2 hover:bg-dark-700 rounded-md transition-colors">
            <Settings class="w-5 h-5" />
          </button>

          <!-- Mobile Menu Toggle -->
          <button
            class="md:hidden p-2 hover:bg-dark-700 rounded-md transition-colors"
            on:click={toggleMenu}
          >
            {#if isMenuOpen}
              <X class="w-5 h-5" />
            {:else}
              <Menu class="w-5 h-5" />
            {/if}
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Navigation -->
    {#if isMenuOpen}
      <div class="md:hidden bg-dark-800 border-t border-dark-600">
        <div class="px-4 py-2 space-y-1">
          <!-- Mobile Search -->
          <div class="py-2">
            <div class="relative">
              <Search
                class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500"
              />
              <input
                type="text"
                placeholder="Search..."
                class="search-input pl-10 pr-4 w-full"
                bind:value={searchQuery}
                on:keydown={handleSearch}
              />
            </div>
          </div>

          <!-- Mobile Nav Items -->
          {#each navItems as item}
            {@const isActive =
              item.href === "/"
                ? $page.url.pathname === "/"
                : $page.url.pathname.startsWith(item.href)}
            <a
              href={item.href}
              class="flex items-center space-x-3 px-3 py-2 rounded-md text-gray-300 hover:text-white hover:bg-dark-700 transition-colors {isActive
                ? 'bg-dark-700 text-neon-pink'
                : ''}"
              on:click={closeMenu}
            >
              <svelte:component this={item.icon} class="w-5 h-5" />
              <span>{item.label}</span>
              {#if item.badge && item.badge > 0}
                <span
                  class="ml-auto inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-dark-900 bg-neon-pink rounded-full"
                >
                  {formatActivityCount(item.badge)}
                </span>
              {/if}
            </a>
          {/each}

          <!-- Mobile Auth -->
          {#if !currentUser}
            <div class="pt-4 border-t border-dark-600 space-y-2">
              <a href="/login" class="block btn-secondary w-full text-center"
                >Login</a
              >
              <a href="/register" class="block btn-primary w-full text-center"
                >Join HI!</a
              >
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </header>

  <!-- Main Content -->
  <main class="min-h-screen">
    <slot />
  </main>

  <!-- Footer -->
  <footer class="bg-dark-800 border-t border-dark-600 mt-auto">
    <div class="container-wide py-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- Brand -->
        <div class="space-y-4">
          <div class="flex items-center space-x-2">
            <Brain class="w-6 h-6 text-neon-pink" />
            <span
              class="font-display font-bold text-xl bg-gradient-to-r from-neon-pink to-neon-cyan bg-clip-text text-transparent"
            >
              Human Intelligence
            </span>
          </div>
          <p class="text-gray-400 text-sm">
            Where AI augments HI. A factorial multiplier to human intelligence.
          </p>
          <div class="flex items-center space-x-4">
            <a
              href="https://github.com/your-org/human-intelligence"
              class="text-gray-400 hover:text-neon-cyan transition-colors"
            >
              <Github class="w-5 h-5" />
            </a>
            <Terminal class="w-5 h-5 text-terminal-400" />
          </div>
        </div>

        <!-- Quick Links -->
        <div class="space-y-4">
          <h3 class="font-semibold text-neon-cyan">Platform</h3>
          <div class="space-y-2 text-sm">
            <a
              href="/about"
              class="block text-gray-400 hover:text-white transition-colors"
              >About</a
            >
            <a
              href="/how-it-works"
              class="block text-gray-400 hover:text-white transition-colors"
              >How it Works</a
            >
            <a
              href="/pricing"
              class="block text-gray-400 hover:text-white transition-colors"
              >Pricing</a
            >
            <a
              href="/api"
              class="block text-gray-400 hover:text-white transition-colors"
              >API</a
            >
          </div>
        </div>

        <!-- Community -->
        <div class="space-y-4">
          <h3 class="font-semibold text-neon-purple">Community</h3>
          <div class="space-y-2 text-sm">
            <a
              href="/community"
              class="block text-gray-400 hover:text-white transition-colors"
              >Forum</a
            >
            <a
              href="/discord"
              class="block text-gray-400 hover:text-white transition-colors"
              >Discord</a
            >
            <a
              href="/contribute"
              class="block text-gray-400 hover:text-white transition-colors"
              >Contribute</a
            >
            <a
              href="/blog"
              class="block text-gray-400 hover:text-white transition-colors"
              >Blog</a
            >
          </div>
        </div>

        <!-- Legal -->
        <div class="space-y-4">
          <h3 class="font-semibold text-neon-green">Legal</h3>
          <div class="space-y-2 text-sm">
            <a
              href="/privacy"
              class="block text-gray-400 hover:text-white transition-colors"
              >Privacy</a
            >
            <a
              href="/terms"
              class="block text-gray-400 hover:text-white transition-colors"
              >Terms</a
            >
            <a
              href="/license"
              class="block text-gray-400 hover:text-white transition-colors"
              >License</a
            >
            <a
              href="/contact"
              class="block text-gray-400 hover:text-white transition-colors"
              >Contact</a
            >
          </div>
        </div>
      </div>

      <div
        class="border-t border-dark-600 mt-8 pt-8 flex flex-col md:flex-row justify-between items-center"
      >
        <div class="text-sm text-gray-500 font-mono">
          © 2024 Human Intelligence. Open source & collaborative.
        </div>
        <div class="text-sm text-gray-500 font-mono mt-2 md:mt-0">
          Made with 🧠 by humans, augmented by AI
        </div>
      </div>
    </div>
  </footer>
</div>

<!-- Toast Container -->
<div class="fixed top-4 right-4 z-50">
  <!-- Toasts will be rendered here by svelte-french-toast -->
</div>

<!-- Offline Indicator -->
{#if !isOnline}
  <div
    class="fixed bottom-4 left-4 bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg z-50 flex items-center space-x-2"
  >
    <div class="w-2 h-2 bg-white rounded-full animate-pulse" />
    <span class="text-sm font-medium">You're offline</span>
  </div>
{/if}

<style>
  /* Custom scrollbar for search results */
  .search-results::-webkit-scrollbar {
    width: 6px;
  }

  .search-results::-webkit-scrollbar-track {
    background: theme("colors.dark.700");
  }

  .search-results::-webkit-scrollbar-thumb {
    background: theme("colors.neon.pink");
    border-radius: 3px;
  }
</style>
