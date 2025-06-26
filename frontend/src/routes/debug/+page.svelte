<script>
  let timestamp = new Date().toISOString();
  let testMessage = "Debug route is working correctly";

  function updateTime() {
    timestamp = new Date().toISOString();
  }
</script>

<svelte:head>
  <title>Debug - Human Intelligence</title>
</svelte:head>

<div class="min-h-screen bg-dark-900 text-white p-8">
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-neon-pink mb-6">🔧 Debug Route</h1>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Status Card -->
      <div class="bg-dark-800 p-6 rounded-lg border border-dark-600">
        <h2 class="text-xl font-semibold text-neon-green mb-4">✅ Route Status</h2>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-400">Route:</span>
            <span class="text-white">/debug</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Status:</span>
            <span class="text-neon-green">Working</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Message:</span>
            <span class="text-white">{testMessage}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Generated:</span>
            <span class="text-gray-300 font-mono text-xs">{timestamp}</span>
          </div>
        </div>
      </div>

      <!-- Interactive Test -->
      <div class="bg-dark-800 p-6 rounded-lg border border-dark-600">
        <h2 class="text-xl font-semibold text-neon-cyan mb-4">🧪 Interactive Test</h2>
        <div class="space-y-4">
          <button
            on:click={updateTime}
            class="w-full bg-neon-pink text-dark-900 px-4 py-2 rounded font-medium hover:bg-neon-cyan transition-colors"
          >
            Update Timestamp
          </button>

          <div class="text-sm text-gray-400">
            Click the button to test Svelte reactivity
          </div>
        </div>
      </div>
    </div>

    <!-- Test Results -->
    <div class="mt-8 bg-dark-800 p-6 rounded-lg border border-dark-600">
      <h2 class="text-xl font-semibold text-white mb-4">Test Results</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
        <div class="text-center">
          <div class="text-2xl mb-2">✅</div>
          <div class="text-neon-green font-medium">SvelteKit Routing</div>
          <div class="text-gray-400">Working</div>
        </div>
        <div class="text-center">
          <div class="text-2xl mb-2">✅</div>
          <div class="text-neon-green font-medium">Component Rendering</div>
          <div class="text-gray-400">Working</div>
        </div>
        <div class="text-center">
          <div class="text-2xl mb-2">✅</div>
          <div class="text-neon-green font-medium">CSS Styling</div>
          <div class="text-gray-400">Working</div>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <div class="mt-8 text-center">
      <div class="space-x-4">
        <a href="/" class="text-neon-cyan hover:text-neon-pink transition-colors">← Home</a>
        <a href="/activity" class="text-neon-cyan hover:text-neon-pink transition-colors">Activity Page</a>
        <a href="/nodes" class="text-neon-cyan hover:text-neon-pink transition-colors">Nodes</a>
      </div>
    </div>
  </div>
</div>
