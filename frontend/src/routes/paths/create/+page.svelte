<script>
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated, user } from "$lib/stores/auth.js";
  import {
    ArrowLeft,
    Plus,
    Search,
    BookOpen,
    Video,
    Link,
    Code,
    GripVertical,
    X,
    Eye,
    Save,
    GitBranch,
    Target,
    Clock,
    Users,
    Tag,
    ChevronRight,
    ChevronDown,
    Play,
    AlertCircle,
    CheckCircle,
  } from "lucide-svelte";

  // Redirect if not authenticated
  $: if (!$isAuthenticated) {
    goto("/login");
  }

  // Path metadata
  let pathData = {
    title: "",
    description: "",
    category: "programming",
    difficulty: "beginner",
    tags: [],
    is_public: true,
    estimated_hours: 0,
  };

  // Node selection and ordering
  let selectedNodes = [];
  let availableNodes = [];
  let searchQuery = "";
  let searchResults = [];
  let isSearching = false;

  // UI state
  let activeTab = "metadata"; // 'metadata', 'nodes', 'preview'
  let showNodeSearch = false;
  let draggedIndex = null;
  let newTag = "";
  let errors = {};
  let isSubmitting = false;
  let previewMode = false;

  // Categories and difficulties
  const categories = [
    { value: "programming", label: "Programming" },
    { value: "data-science", label: "Data Science" },
    { value: "design", label: "Design" },
    { value: "business", label: "Business" },
    { value: "languages", label: "Languages" },
    { value: "science", label: "Science" },
  ];

  const difficulties = [
    { value: "beginner", label: "Beginner", color: "text-neon-green" },
    { value: "intermediate", label: "Intermediate", color: "text-neon-yellow" },
    { value: "advanced", label: "Advanced", color: "text-neon-pink" },
  ];

  onMount(() => {
    loadAvailableNodes();
  });

  // Load available nodes for selection
  async function loadAvailableNodes() {
    try {
      // Mock data - replace with actual API call
      availableNodes = [
        {
          id: "1",
          title: "Introduction to JavaScript",
          description: "Learn the basics of JavaScript programming",
          type: "text",
          estimated_minutes: 45,
          difficulty: "beginner",
          tags: ["javascript", "basics", "programming"],
          author: "alice_dev",
        },
        {
          id: "2",
          title: "JavaScript Functions Deep Dive",
          description: "Master functions, closures, and scope",
          type: "video",
          estimated_minutes: 60,
          difficulty: "intermediate",
          tags: ["javascript", "functions", "closures"],
          author: "bob_teaches",
        },
        {
          id: "3",
          title: "Building Your First React App",
          description: "Step-by-step React application tutorial",
          type: "code",
          estimated_minutes: 90,
          difficulty: "beginner",
          tags: ["react", "javascript", "tutorial"],
          author: "charlie_react",
        },
        {
          id: "4",
          title: "React Hooks Explained",
          description: "useState, useEffect, and custom hooks",
          type: "text",
          estimated_minutes: 40,
          difficulty: "intermediate",
          tags: ["react", "hooks", "modern"],
          author: "diana_frontend",
        },
        {
          id: "5",
          title: "API Integration with Fetch",
          description: "Learn to connect to REST APIs",
          type: "link",
          estimated_minutes: 30,
          difficulty: "beginner",
          tags: ["api", "fetch", "javascript"],
          author: "eve_api",
        },
      ];
    } catch (error) {
      console.error("Failed to load nodes:", error);
    }
  }

  // Search nodes
  function searchNodes() {
    if (!searchQuery.trim()) {
      searchResults = [];
      return;
    }

    isSearching = true;

    // Simulate search delay
    setTimeout(() => {
      searchResults = availableNodes
        .filter(
          (node) =>
            node.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            node.description
              .toLowerCase()
              .includes(searchQuery.toLowerCase()) ||
            node.tags.some((tag) =>
              tag.toLowerCase().includes(searchQuery.toLowerCase())
            )
        )
        .filter(
          (node) => !selectedNodes.find((selected) => selected.id === node.id)
        );

      isSearching = false;
    }, 300);
  }

  // Add node to path
  function addNode(node) {
    if (!selectedNodes.find((n) => n.id === node.id)) {
      selectedNodes = [
        ...selectedNodes,
        { ...node, order: selectedNodes.length },
      ];
      updateEstimatedHours();
      searchQuery = "";
      searchResults = [];
      showNodeSearch = false;
    }
  }

  // Remove node from path
  function removeNode(index) {
    selectedNodes = selectedNodes.filter((_, i) => i !== index);
    updateEstimatedHours();
  }

  // Drag and drop handling
  function handleDragStart(event, index) {
    draggedIndex = index;
    event.dataTransfer.effectAllowed = "move";
  }

  function handleDragOver(event) {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }

  function handleDrop(event, dropIndex) {
    event.preventDefault();

    if (draggedIndex !== null && draggedIndex !== dropIndex) {
      const draggedNode = selectedNodes[draggedIndex];
      const newNodes = [...selectedNodes];

      // Remove dragged item
      newNodes.splice(draggedIndex, 1);

      // Insert at new position
      if (dropIndex > draggedIndex) {
        newNodes.splice(dropIndex - 1, 0, draggedNode);
      } else {
        newNodes.splice(dropIndex, 0, draggedNode);
      }

      selectedNodes = newNodes;
    }

    draggedIndex = null;
  }

  // Tag management
  function addTag() {
    if (newTag.trim() && !pathData.tags.includes(newTag.trim())) {
      pathData.tags = [...pathData.tags, newTag.trim()];
      newTag = "";
    }
  }

  function removeTag(tag) {
    pathData.tags = pathData.tags.filter((t) => t !== tag);
  }

  function handleTagKeyPress(event) {
    if (event.key === "Enter" || event.key === ",") {
      event.preventDefault();
      addTag();
    }
  }

  // Update estimated hours based on selected nodes
  function updateEstimatedHours() {
    const totalMinutes = selectedNodes.reduce(
      (sum, node) => sum + node.estimated_minutes,
      0
    );
    pathData.estimated_hours = Math.round((totalMinutes / 60) * 10) / 10;
  }

  // Validation
  function validatePath() {
    errors = {};

    if (!pathData.title.trim()) {
      errors.title = "Title is required";
    }

    if (!pathData.description.trim()) {
      errors.description = "Description is required";
    }

    if (selectedNodes.length === 0) {
      errors.nodes = "At least one node is required";
    }

    return Object.keys(errors).length === 0;
  }

  // Save path
  async function savePath() {
    if (!validatePath()) return;

    isSubmitting = true;

    try {
      const pathPayload = {
        ...pathData,
        nodes: selectedNodes.map((node, index) => ({
          node_id: node.id,
          order: index,
          is_required: true,
        })),
        created_by: $user.id,
      };

      // Mock API call - replace with actual implementation
      console.log("Creating path:", pathPayload);

      // Simulate API delay
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Redirect to the new path
      goto("/paths");
    } catch (error) {
      console.error("Failed to create path:", error);
      errors.submit = "Failed to create path. Please try again.";
    } finally {
      isSubmitting = false;
    }
  }

  // Node type icons
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

  // Reactive search
  $: if (searchQuery) {
    searchNodes();
  }
</script>

<svelte:head>
  <title>Create Learning Path - Human Intelligence</title>
  <meta
    name="description"
    content="Create a new learning path by composing nodes"
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
    <div class="flex items-center justify-between mb-8">
      <div class="flex items-center space-x-4">
        <button on:click={() => goto("/paths")} class="btn-ghost p-2">
          <ArrowLeft class="w-5 h-5" />
        </button>
        <div>
          <h1 class="text-3xl font-display font-bold text-white">
            Create Learning Path
          </h1>
          <p class="text-gray-400">Compose a sequence of learning nodes</p>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <button
          on:click={() => (previewMode = !previewMode)}
          class="btn-ghost space-x-2"
        >
          <Eye class="w-4 h-4" />
          <span>Preview</span>
        </button>

        <button
          on:click={savePath}
          disabled={isSubmitting}
          class="btn-primary flex items-center space-x-2 {isSubmitting
            ? 'opacity-50 cursor-not-allowed'
            : ''}"
        >
          {#if isSubmitting}
            <div
              class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
            />
            <span>Creating...</span>
          {:else}
            <Save class="w-4 h-4" />
            <span>Create Path</span>
          {/if}
        </button>
      </div>
    </div>

    <!-- Tab Navigation -->
    <div class="flex bg-dark-800 rounded-lg p-1 border border-dark-600 mb-8">
      <button
        on:click={() => (activeTab = "metadata")}
        class="px-4 py-2 rounded text-sm transition-all {activeTab ===
        'metadata'
          ? 'bg-neon-purple text-white'
          : 'text-gray-400 hover:text-white'}"
      >
        Path Details
      </button>
      <button
        on:click={() => (activeTab = "nodes")}
        class="px-4 py-2 rounded text-sm transition-all {activeTab === 'nodes'
          ? 'bg-neon-purple text-white'
          : 'text-gray-400 hover:text-white'}"
      >
        Add Nodes ({selectedNodes.length})
      </button>
      <button
        on:click={() => (activeTab = "preview")}
        class="px-4 py-2 rounded text-sm transition-all {activeTab === 'preview'
          ? 'bg-neon-purple text-white'
          : 'text-gray-400 hover:text-white'}"
      >
        Preview
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <!-- Main Content -->
      <div class="lg:col-span-3">
        {#if activeTab === "metadata"}
          <!-- Path Metadata Form -->
          <div class="card space-y-6">
            <h2 class="text-xl font-display font-bold text-white">
              Path Information
            </h2>

            <!-- Title -->
            <div class="space-y-2">
              <label
                for="title"
                class="block text-sm font-medium text-gray-300"
              >
                Title *
              </label>
              <input
                id="title"
                type="text"
                bind:value={pathData.title}
                class="input-neon {errors.title
                  ? 'border-red-500 focus:border-red-500'
                  : ''}"
                placeholder="Enter a descriptive title for your learning path..."
                maxlength="255"
              />
              {#if errors.title}
                <p class="text-red-400 text-sm">{errors.title}</p>
              {/if}
            </div>

            <!-- Description -->
            <div class="space-y-2">
              <label
                for="description"
                class="block text-sm font-medium text-gray-300"
              >
                Description *
              </label>
              <textarea
                id="description"
                bind:value={pathData.description}
                class="input-neon min-h-[100px] resize-none {errors.description
                  ? 'border-red-500 focus:border-red-500'
                  : ''}"
                placeholder="Describe what learners will achieve by completing this path..."
                maxlength="1000"
              />
              {#if errors.description}
                <p class="text-red-400 text-sm">{errors.description}</p>
              {/if}
            </div>

            <!-- Category and Difficulty -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label
                  for="category"
                  class="block text-sm font-medium text-gray-300"
                >
                  Category
                </label>
                <select
                  id="category"
                  bind:value={pathData.category}
                  class="input-neon"
                >
                  {#each categories as category}
                    <option value={category.value}>{category.label}</option>
                  {/each}
                </select>
              </div>

              <div class="space-y-2">
                <label
                  for="difficulty"
                  class="block text-sm font-medium text-gray-300"
                >
                  Difficulty Level
                </label>
                <select
                  id="difficulty"
                  bind:value={pathData.difficulty}
                  class="input-neon"
                >
                  {#each difficulties as difficulty}
                    <option value={difficulty.value}>{difficulty.label}</option>
                  {/each}
                </select>
              </div>
            </div>

            <!-- Tags -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-300">
                Tags
              </label>
              <div class="flex flex-wrap gap-2 mb-3">
                {#each pathData.tags as tag}
                  <span
                    class="inline-flex items-center px-3 py-1 bg-dark-700 text-neon-cyan text-sm rounded-full border border-dark-600"
                  >
                    {tag}
                    <button
                      on:click={() => removeTag(tag)}
                      class="ml-2 text-gray-400 hover:text-red-400"
                    >
                      <X class="w-3 h-3" />
                    </button>
                  </span>
                {/each}
              </div>
              <input
                type="text"
                bind:value={newTag}
                on:keypress={handleTagKeyPress}
                class="input-neon"
                placeholder="Add tags (press Enter or comma to add)"
              />
            </div>

            <!-- Visibility -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-300">
                Visibility
              </label>
              <div class="flex items-center space-x-3">
                <label class="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="radio"
                    bind:group={pathData.is_public}
                    value={true}
                    class="text-neon-pink focus:ring-neon-pink"
                  />
                  <span class="text-gray-300"
                    >Public - Anyone can view and fork</span
                  >
                </label>
                <label class="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="radio"
                    bind:group={pathData.is_public}
                    value={false}
                    class="text-neon-pink focus:ring-neon-pink"
                  />
                  <span class="text-gray-300">Private - Only you can view</span>
                </label>
              </div>
            </div>
          </div>
        {:else if activeTab === "nodes"}
          <!-- Node Selection and Ordering -->
          <div class="space-y-6">
            <!-- Add Nodes Section -->
            <div class="card">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-display font-bold text-white">
                  Add Nodes
                </h2>
                <button
                  on:click={() => (showNodeSearch = !showNodeSearch)}
                  class="btn-secondary space-x-2"
                >
                  <Plus class="w-4 h-4" />
                  <span>Add Node</span>
                </button>
              </div>

              {#if showNodeSearch}
                <div class="border-t border-dark-600 pt-4">
                  <!-- Search -->
                  <div class="relative mb-4">
                    <div
                      class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                    >
                      <Search class="w-5 h-5 text-gray-500" />
                    </div>
                    <input
                      type="text"
                      bind:value={searchQuery}
                      class="input-neon pl-10"
                      placeholder="Search nodes to add..."
                    />
                  </div>

                  <!-- Search Results -->
                  {#if isSearching}
                    <div class="flex justify-center py-4">
                      <div
                        class="w-6 h-6 border-2 border-neon-cyan border-t-transparent rounded-full animate-spin"
                      />
                    </div>
                  {:else if searchResults.length > 0}
                    <div class="space-y-2">
                      {#each searchResults as node}
                        <div
                          class="flex items-center justify-between p-3 bg-dark-700 rounded-lg border border-dark-600 hover:border-neon-purple transition-colors"
                        >
                          <div class="flex items-center space-x-3">
                            <svelte:component
                              this={getNodeTypeIcon(node.type)}
                              class="w-5 h-5 {getNodeTypeColor(node.type)}"
                            />
                            <div>
                              <h3 class="font-medium text-white">
                                {node.title}
                              </h3>
                              <p class="text-sm text-gray-400 line-clamp-1">
                                {node.description}
                              </p>
                              <div
                                class="flex items-center space-x-2 mt-1 text-xs text-gray-500"
                              >
                                <span>{node.estimated_minutes}min</span>
                                <span>•</span>
                                <span>{node.difficulty}</span>
                                <span>•</span>
                                <span>by @{node.author}</span>
                              </div>
                            </div>
                          </div>
                          <button
                            on:click={() => addNode(node)}
                            class="btn-primary px-3 py-1 text-sm"
                          >
                            Add
                          </button>
                        </div>
                      {/each}
                    </div>
                  {:else if searchQuery}
                    <div class="text-center py-8 text-gray-500">
                      <Search class="w-12 h-12 mx-auto mb-2" />
                      <p>No nodes found matching "{searchQuery}"</p>
                    </div>
                  {/if}
                </div>
              {/if}
            </div>

            <!-- Selected Nodes (Path Sequence) -->
            <div class="card">
              <h2 class="text-xl font-display font-bold text-white mb-4">
                Path Sequence ({selectedNodes.length} nodes)
              </h2>

              {#if errors.nodes}
                <div
                  class="bg-red-500/10 border border-red-500/30 rounded-lg p-3 mb-4"
                >
                  <div class="flex items-center space-x-2">
                    <AlertCircle class="w-4 h-4 text-red-400" />
                    <span class="text-red-400 text-sm">{errors.nodes}</span>
                  </div>
                </div>
              {/if}

              {#if selectedNodes.length === 0}
                <div class="text-center py-8 text-gray-500">
                  <GitBranch class="w-12 h-12 mx-auto mb-2" />
                  <p>No nodes added yet</p>
                  <p class="text-sm">Add nodes to create your learning path</p>
                </div>
              {:else}
                <div class="space-y-3">
                  {#each selectedNodes as node, index}
                    <div
                      class="flex items-center space-x-4 p-4 bg-dark-700 rounded-lg border border-dark-600 group"
                      draggable="true"
                      on:dragstart={(e) => handleDragStart(e, index)}
                      on:dragover={handleDragOver}
                      on:drop={(e) => handleDrop(e, index)}
                    >
                      <!-- Drag Handle -->
                      <div
                        class="cursor-move text-gray-500 group-hover:text-gray-300"
                      >
                        <GripVertical class="w-5 h-5" />
                      </div>

                      <!-- Step Number -->
                      <div
                        class="w-8 h-8 bg-neon-purple rounded-full flex items-center justify-center text-white font-bold text-sm"
                      >
                        {index + 1}
                      </div>

                      <!-- Node Info -->
                      <div class="flex-1">
                        <div class="flex items-center space-x-2 mb-1">
                          <svelte:component
                            this={getNodeTypeIcon(node.type)}
                            class="w-4 h-4 {getNodeTypeColor(node.type)}"
                          />
                          <h3 class="font-medium text-white">{node.title}</h3>
                        </div>
                        <p class="text-sm text-gray-400 line-clamp-1">
                          {node.description}
                        </p>
                        <div
                          class="flex items-center space-x-2 mt-1 text-xs text-gray-500"
                        >
                          <Clock class="w-3 h-3" />
                          <span>{node.estimated_minutes}min</span>
                          <span>•</span>
                          <span>{node.difficulty}</span>
                        </div>
                      </div>

                      <!-- Remove Button -->
                      <button
                        on:click={() => removeNode(index)}
                        class="opacity-0 group-hover:opacity-100 btn-ghost p-2 text-red-400 hover:text-red-300 transition-all"
                      >
                        <X class="w-4 h-4" />
                      </button>
                    </div>

                    {#if index < selectedNodes.length - 1}
                      <div class="flex justify-center">
                        <ChevronDown class="w-5 h-5 text-gray-600" />
                      </div>
                    {/if}
                  {/each}
                </div>
              {/if}
            </div>
          </div>
        {:else if activeTab === "preview"}
          <!-- Path Preview -->
          <div class="card">
            <h2 class="text-xl font-display font-bold text-white mb-6">
              Path Preview
            </h2>

            <!-- Path Header -->
            <div class="border-b border-dark-600 pb-6 mb-6">
              <div class="flex items-start space-x-4">
                <div
                  class="w-16 h-16 bg-gradient-to-br from-neon-purple to-neon-cyan rounded-lg flex items-center justify-center"
                >
                  <GitBranch class="w-8 h-8 text-white" />
                </div>
                <div class="flex-1">
                  <h1 class="text-2xl font-bold text-white mb-2">
                    {pathData.title || "Untitled Path"}
                  </h1>
                  <p class="text-gray-400 mb-3">
                    {pathData.description || "No description provided"}
                  </p>

                  <div class="flex items-center space-x-4 text-sm">
                    <span class="flex items-center space-x-1 text-gray-400">
                      <BookOpen class="w-4 h-4" />
                      <span>{selectedNodes.length} nodes</span>
                    </span>
                    <span class="flex items-center space-x-1 text-gray-400">
                      <Clock class="w-4 h-4" />
                      <span>{pathData.estimated_hours}h</span>
                    </span>
                    <span class="flex items-center space-x-1">
                      <Target class="w-4 h-4" />
                      <span
                        class="capitalize {difficulties.find(
                          (d) => d.value === pathData.difficulty
                        )?.color}"
                      >
                        {pathData.difficulty}
                      </span>
                    </span>
                  </div>

                  {#if pathData.tags.length > 0}
                    <div class="flex flex-wrap gap-1 mt-3">
                      {#each pathData.tags as tag}
                        <span
                          class="px-2 py-1 bg-dark-700 text-neon-cyan text-xs rounded-full border border-dark-600"
                        >
                          {tag}
                        </span>
                      {/each}
                    </div>
                  {/if}
                </div>
              </div>
            </div>

            <!-- Path Content -->
            {#if selectedNodes.length > 0}
              <div class="space-y-4">
                <h3 class="text-lg font-semibold text-white">
                  Learning Sequence
                </h3>

                {#each selectedNodes as node, index}
                  <div
                    class="flex items-center space-x-4 p-4 bg-dark-700 rounded-lg border border-dark-600"
                  >
                    <div
                      class="w-8 h-8 bg-neon-purple rounded-full flex items-center justify-center text-white font-bold text-sm"
                    >
                      {index + 1}
                    </div>

                    <div class="flex-1">
                      <div class="flex items-center space-x-2 mb-1">
                        <svelte:component
                          this={getNodeTypeIcon(node.type)}
                          class="w-4 h-4 {getNodeTypeColor(node.type)}"
                        />
                        <h4 class="font-medium text-white">{node.title}</h4>
                      </div>
                      <p class="text-sm text-gray-400">{node.description}</p>
                      <div
                        class="flex items-center space-x-2 mt-2 text-xs text-gray-500"
                      >
                        <Clock class="w-3 h-3" />
                        <span>{node.estimated_minutes} minutes</span>
                        <span>•</span>
                        <Target class="w-3 h-3" />
                        <span>{node.difficulty}</span>
                      </div>
                    </div>

                    <ChevronRight class="w-5 h-5 text-gray-500" />
                  </div>
                {/each}
              </div>
            {:else}
              <div class="text-center py-8 text-gray-500">
                <GitBranch class="w-12 h-12 mx-auto mb-2" />
                <p>Add some nodes to see the path preview</p>
              </div>
            {/if}
          </div>
        {/if}
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Path Stats -->
        <div class="card">
          <h3 class="text-lg font-display font-bold text-white mb-4">
            Path Statistics
          </h3>

          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-gray-400">Total Nodes</span>
              <span class="text-neon-cyan font-semibold"
                >{selectedNodes.length}</span
              >
            </div>

            <div class="flex justify-between items-center">
              <span class="text-gray-400">Estimated Time</span>
              <span class="text-neon-cyan font-semibold"
                >{pathData.estimated_hours}h</span
              >
            </div>

            <div class="flex justify-between items-center">
              <span class="text-gray-400">Difficulty</span>
              <span
                class="capitalize {difficulties.find(
                  (d) => d.value === pathData.difficulty
                )?.color} font-semibold"
              >
                {pathData.difficulty}
              </span>
            </div>

            <div class="flex justify-between items-center">
              <span class="text-gray-400">Visibility</span>
              <span
                class="text-{pathData.is_public
                  ? 'neon-green'
                  : 'neon-yellow'} font-semibold"
              >
                {pathData.is_public ? "Public" : "Private"}
              </span>
            </div>

            <div class="flex justify-between items-center">
              <span class="text-gray-400">Category</span>
              <span class="text-white font-semibold capitalize"
                >{pathData.category}</span
              >
            </div>
          </div>
        </div>

        <!-- Validation Status -->
        <div class="card">
          <h3 class="text-lg font-display font-bold text-white mb-4">
            Validation
          </h3>

          <div class="space-y-2">
            <div class="flex items-center space-x-2">
              {#if pathData.title.trim()}
                <CheckCircle class="w-4 h-4 text-neon-green" />
              {:else}
                <AlertCircle class="w-4 h-4 text-red-400" />
              {/if}
              <span
                class="text-sm {pathData.title.trim()
                  ? 'text-gray-300'
                  : 'text-red-400'}"
              >
                Title provided
              </span>
            </div>

            <div class="flex items-center space-x-2">
              {#if pathData.description.trim()}
                <CheckCircle class="w-4 h-4 text-neon-green" />
              {:else}
                <AlertCircle class="w-4 h-4 text-red-400" />
              {/if}
              <span
                class="text-sm {pathData.description.trim()
                  ? 'text-gray-300'
                  : 'text-red-400'}"
              >
                Description provided
              </span>
            </div>

            <div class="flex items-center space-x-2">
              {#if selectedNodes.length > 0}
                <CheckCircle class="w-4 h-4 text-neon-green" />
              {:else}
                <AlertCircle class="w-4 h-4 text-red-400" />
              {/if}
              <span
                class="text-sm {selectedNodes.length > 0
                  ? 'text-gray-300'
                  : 'text-red-400'}"
              >
                At least one node added
              </span>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="card">
          <h3 class="text-lg font-display font-bold text-white mb-4">
            Quick Actions
          </h3>

          <div class="space-y-2">
            <button
              on:click={() => (activeTab = "metadata")}
              class="w-full text-left p-2 rounded hover:bg-dark-700 transition-colors text-sm text-gray-300"
            >
              Edit Path Details
            </button>
            <button
              on:click={() => (activeTab = "nodes")}
              class="w-full text-left p-2 rounded hover:bg-dark-700 transition-colors text-sm text-gray-300"
            >
              Manage Nodes
            </button>
            <button
              on:click={() => (activeTab = "preview")}
              class="w-full text-left p-2 rounded hover:bg-dark-700 transition-colors text-sm text-gray-300"
            >
              Preview Path
            </button>
          </div>
        </div>

        <!-- Tips -->
        <div class="card">
          <h3 class="text-lg font-display font-bold text-white mb-4">Tips</h3>

          <div class="space-y-3 text-sm text-gray-400">
            <div class="flex items-start space-x-2">
              <div
                class="w-1.5 h-1.5 bg-neon-cyan rounded-full mt-2 flex-shrink-0"
              />
              <p>
                Start with foundational concepts and build complexity gradually
              </p>
            </div>
            <div class="flex items-start space-x-2">
              <div
                class="w-1.5 h-1.5 bg-neon-cyan rounded-full mt-2 flex-shrink-0"
              />
              <p>Mix different node types to keep learners engaged</p>
            </div>
            <div class="flex items-start space-x-2">
              <div
                class="w-1.5 h-1.5 bg-neon-cyan rounded-full mt-2 flex-shrink-0"
              />
              <p>Drag and drop to reorder nodes in your sequence</p>
            </div>
            <div class="flex items-start space-x-2">
              <div
                class="w-1.5 h-1.5 bg-neon-cyan rounded-full mt-2 flex-shrink-0"
              />
              <p>Use descriptive tags to help others discover your path</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Messages -->
    {#if errors.submit}
      <div
        class="fixed bottom-4 right-4 bg-red-500/10 border border-red-500/30 rounded-lg p-4 max-w-md"
      >
        <div class="flex items-center space-x-2">
          <AlertCircle class="w-5 h-5 text-red-400" />
          <span class="text-red-400">{errors.submit}</span>
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

  .line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* Drag and drop visual feedback */
  .dragging {
    opacity: 0.5;
    transform: rotate(5deg);
  }

  /* Smooth transitions for node reordering */
  .space-y-3 > * {
    transition: transform 0.2s ease;
  }

  /* Custom radio button styling */
  input[type="radio"] {
    accent-color: var(--neon-pink);
  }

  /* Tab active state animation */
  .flex.bg-dark-800 button {
    position: relative;
    transition: all 0.3s ease;
  }

  .flex.bg-dark-800 button.bg-neon-purple::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(45deg, var(--neon-purple), var(--neon-pink));
    border-radius: 0.375rem;
    opacity: 0.1;
    z-index: -1;
  }
</style>
