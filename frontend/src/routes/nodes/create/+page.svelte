<script>
  import { apiRequest, isAuthenticated } from "$lib/stores/auth.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    FileText,
    Video,
    Link,
    Code,
    Hash,
    Save,
    Eye,
    EyeOff,
    Plus,
    X,
    Brain,
    AlertCircle,
  } from "lucide-svelte";

  let title = "";
  let description = "";
  let content = "";
  let nodeType = "text";
  let tags = [];
  let newTag = "";
  let isPublic = true;
  let metadata = {};
  let errors = {};
  let serverError = "";
  let isSubmitting = false;
  let showPreview = false;

  // Node type options
  const nodeTypes = [
    {
      value: "text",
      label: "Text",
      icon: FileText,
      description: "Rich text content",
    },
    {
      value: "video",
      label: "Video",
      icon: Video,
      description: "Video content with URL",
    },
    {
      value: "link",
      label: "Link",
      icon: Link,
      description: "External resource link",
    },
    {
      value: "code",
      label: "Code",
      icon: Code,
      description: "Code snippets and examples",
    },
  ];

  // Reactive metadata fields based on node type
  $: metadataFields = getMetadataFields(nodeType);

  function getMetadataFields(type) {
    switch (type) {
      case "video":
        return [
          {
            key: "video_url",
            label: "Video URL",
            type: "url",
            placeholder: "https://youtube.com/watch?v=...",
          },
          {
            key: "video_duration",
            label: "Duration (seconds)",
            type: "number",
            placeholder: "3600",
          },
        ];
      case "link":
        return [
          {
            key: "link_url",
            label: "URL",
            type: "url",
            placeholder: "https://example.com",
          },
          {
            key: "link_title",
            label: "Link Title",
            type: "text",
            placeholder: "Automatically detected",
          },
        ];
      case "code":
        return [
          {
            key: "language",
            label: "Programming Language",
            type: "text",
            placeholder: "javascript",
          },
          {
            key: "repository",
            label: "Repository URL",
            type: "url",
            placeholder: "https://github.com/...",
          },
        ];
      default:
        return [];
    }
  }

  // Form validation
  function validateForm() {
    errors = {};

    if (!title.trim()) {
      errors.title = "Title is required";
    } else if (title.length > 255) {
      errors.title = "Title must be less than 255 characters";
    }

    if (
      nodeType === "video" &&
      metadata.video_url &&
      !isValidUrl(metadata.video_url)
    ) {
      errors.video_url = "Please enter a valid video URL";
    }

    if (
      nodeType === "link" &&
      metadata.link_url &&
      !isValidUrl(metadata.link_url)
    ) {
      errors.link_url = "Please enter a valid URL";
    }

    if (!content.trim() && nodeType !== "link") {
      errors.content = "Content is required";
    }

    return Object.keys(errors).length === 0;
  }

  function isValidUrl(string) {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  }

  // Handle form submission
  async function handleSubmit() {
    if (!validateForm()) return;

    isSubmitting = true;
    serverError = "";

    try {
      const response = await apiRequest("/nodes/", {
        method: "POST",
        body: JSON.stringify({
          title: title.trim(),
          description: description.trim(),
          content: content.trim(),
          node_type: nodeType,
          metadata: metadata,
          tags: tags,
          is_public: isPublic,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        // Redirect to the created node or nodes list
        goto(`/nodes/${result.node.id}`);
      } else {
        const error = await response.json();
        serverError = error.error || "Failed to create node";
      }
    } catch (error) {
      console.error("Node creation error:", error);
      serverError = "Network error. Please try again.";
    } finally {
      isSubmitting = false;
    }
  }

  // Tag management
  function addTag() {
    if (newTag.trim() && !tags.includes(newTag.trim())) {
      tags = [...tags, newTag.trim()];
      newTag = "";
    }
  }

  function removeTag(tagToRemove) {
    tags = tags.filter((tag) => tag !== tagToRemove);
  }

  function handleTagKeyPress(event) {
    if (event.key === "Enter" || event.key === ",") {
      event.preventDefault();
      addTag();
    }
  }

  // Metadata management
  function updateMetadata(key, value) {
    metadata = { ...metadata, [key]: value };
  }

  // Clear errors when user starts typing
  function clearError(field) {
    if (errors[field]) {
      errors = { ...errors };
      delete errors[field];
    }
    if (serverError) {
      serverError = "";
    }
  }

  // Redirect if not authenticated
  onMount(() => {
    if (!$isAuthenticated) {
      goto("/login");
      return;
    }

    // Focus on title field
    document.querySelector("#title")?.focus();
  });
</script>

<svelte:head>
  <title>Create Node - Human Intelligence</title>
  <meta
    name="description"
    content="Create a new learning node and share knowledge with the community."
  />
</svelte:head>

{#if !$isAuthenticated}
  <div class="min-h-screen bg-dark-900 flex items-center justify-center">
    <div class="text-center">
      <Brain class="w-16 h-16 text-neon-pink mx-auto mb-4 animate-pulse" />
      <p class="text-gray-400">Redirecting to login...</p>
    </div>
  </div>
{:else}
  <div class="min-h-screen bg-dark-900">
    <!-- Background Effects -->
    <div class="absolute inset-0 bg-grid opacity-10" />

    <div class="container-wide py-8">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center space-x-3 mb-4">
          <FileText class="w-8 h-8 text-neon-cyan" />
          <h1 class="text-4xl font-display font-bold text-white">
            Create New <span class="text-neon-pink">Node</span>
          </h1>
        </div>
        <p class="text-gray-400 text-lg">
          Share knowledge, create content, and contribute to the collective
          intelligence
        </p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Main Form -->
        <div class="lg:col-span-2">
          <form on:submit|preventDefault={handleSubmit} class="space-y-6">
            <!-- Server Error -->
            {#if serverError}
              <div
                class="bg-red-500/10 border border-red-500/30 rounded-lg p-4"
              >
                <div class="flex items-center space-x-2">
                  <AlertCircle class="w-5 h-5 text-red-400" />
                  <span class="text-red-400 text-sm">{serverError}</span>
                </div>
              </div>
            {/if}

            <!-- Title -->
            <div class="card space-y-4">
              <h3 class="text-xl font-display font-bold text-neon-cyan">
                Basic Information
              </h3>

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
                  bind:value={title}
                  on:input={() => clearError("title")}
                  class="input-neon {errors.title
                    ? 'border-red-500 focus:border-red-500'
                    : ''}"
                  placeholder="Enter a descriptive title for your node..."
                  disabled={isSubmitting}
                  maxlength="255"
                />
                {#if errors.title}
                  <p class="text-red-400 text-sm">{errors.title}</p>
                {:else}
                  <p class="text-gray-500 text-xs">
                    {title.length}/255 characters
                  </p>
                {/if}
              </div>

              <div class="space-y-2">
                <label
                  for="description"
                  class="block text-sm font-medium text-gray-300"
                >
                  Description
                </label>
                <textarea
                  id="description"
                  bind:value={description}
                  on:input={() => clearError("description")}
                  class="input-neon min-h-[80px] resize-none"
                  placeholder="Brief description of what this node covers..."
                  disabled={isSubmitting}
                />
                <p class="text-gray-500 text-xs">
                  Help others understand what they'll learn
                </p>
              </div>
            </div>

            <!-- Node Type Selection -->
            <div class="card space-y-4">
              <h3 class="text-xl font-display font-bold text-neon-cyan">
                Content Type
              </h3>

              <div class="grid grid-cols-2 gap-4">
                {#each nodeTypes as type}
                  <label class="cursor-pointer">
                    <input
                      type="radio"
                      bind:group={nodeType}
                      value={type.value}
                      class="sr-only"
                      disabled={isSubmitting}
                    />
                    <div
                      class="border-2 border-dark-600 rounded-lg p-4 transition-all duration-300 hover:border-neon-purple
                                {nodeType === type.value
                        ? 'border-neon-pink bg-neon-pink/10'
                        : ''}"
                    >
                      <div class="flex items-center space-x-3">
                        <svelte:component
                          this={type.icon}
                          class="w-6 h-6 text-neon-cyan"
                        />
                        <div>
                          <div class="font-medium text-white">{type.label}</div>
                          <div class="text-sm text-gray-400">
                            {type.description}
                          </div>
                        </div>
                      </div>
                    </div>
                  </label>
                {/each}
              </div>
            </div>

            <!-- Type-specific Metadata -->
            {#if metadataFields.length > 0}
              <div class="card space-y-4">
                <h3 class="text-xl font-display font-bold text-neon-cyan">
                  {nodeType.charAt(0).toUpperCase() + nodeType.slice(1)} Settings
                </h3>

                {#each metadataFields as field}
                  <div class="space-y-2">
                    <label
                      for={field.key}
                      class="block text-sm font-medium text-gray-300"
                    >
                      {field.label}
                    </label>
                    <input
                      id={field.key}
                      type={field.type}
                      value={metadata[field.key] || ""}
                      on:input={(e) =>
                        updateMetadata(field.key, e.target.value)}
                      on:input={() => clearError(field.key)}
                      class="input-neon {errors[field.key]
                        ? 'border-red-500 focus:border-red-500'
                        : ''}"
                      placeholder={field.placeholder}
                      disabled={isSubmitting}
                    />
                    {#if errors[field.key]}
                      <p class="text-red-400 text-sm">{errors[field.key]}</p>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}

            <!-- Content -->
            <div class="card space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="text-xl font-display font-bold text-neon-cyan">
                  Content
                </h3>
                <button
                  type="button"
                  on:click={() => (showPreview = !showPreview)}
                  class="btn-ghost text-sm"
                  disabled={isSubmitting}
                >
                  {#if showPreview}
                    <EyeOff class="w-4 h-4 mr-1" />
                    Edit
                  {:else}
                    <Eye class="w-4 h-4 mr-1" />
                    Preview
                  {/if}
                </button>
              </div>

              {#if showPreview}
                <div
                  class="bg-dark-800 border border-dark-600 rounded-lg p-4 min-h-[200px]"
                >
                  <div class="prose prose-invert max-w-none">
                    {#if content.trim()}
                      {@html content.replace(/\n/g, "<br>")}
                    {:else}
                      <p class="text-gray-500 italic">
                        Preview will appear here...
                      </p>
                    {/if}
                  </div>
                </div>
              {:else}
                <div class="space-y-2">
                  <textarea
                    bind:value={content}
                    on:input={() => clearError("content")}
                    class="input-neon min-h-[200px] resize-none font-mono {errors.content
                      ? 'border-red-500 focus:border-red-500'
                      : ''}"
                    placeholder={nodeType === "link"
                      ? "Optional: Add notes about this link..."
                      : "Enter your content here...\n\nUse markdown syntax for formatting."}
                    disabled={isSubmitting}
                  />
                  {#if errors.content}
                    <p class="text-red-400 text-sm">{errors.content}</p>
                  {:else}
                    <p class="text-gray-500 text-xs">
                      {nodeType === "text"
                        ? "Supports markdown formatting"
                        : nodeType === "code"
                        ? "Add code examples, explanations, and usage notes"
                        : nodeType === "video"
                        ? "Add notes, timestamps, and key takeaways"
                        : "Add context and additional information"}
                    </p>
                  {/if}
                </div>
              {/if}
            </div>

            <!-- Tags -->
            <div class="card space-y-4">
              <h3 class="text-xl font-display font-bold text-neon-cyan">
                Tags
              </h3>

              <div class="space-y-3">
                <div class="flex items-center space-x-2">
                  <Hash class="w-5 h-5 text-gray-500" />
                  <input
                    type="text"
                    bind:value={newTag}
                    on:keypress={handleTagKeyPress}
                    class="input-neon flex-1"
                    placeholder="Add tags (press Enter or comma to add)"
                    disabled={isSubmitting}
                  />
                  <button
                    type="button"
                    on:click={addTag}
                    class="btn-secondary text-sm px-3 py-2"
                    disabled={isSubmitting || !newTag.trim()}
                  >
                    <Plus class="w-4 h-4" />
                  </button>
                </div>

                {#if tags.length > 0}
                  <div class="flex flex-wrap gap-2">
                    {#each tags as tag}
                      <span
                        class="inline-flex items-center bg-neon-purple/20 text-neon-purple border border-neon-purple/30 rounded-full px-3 py-1 text-sm"
                      >
                        #{tag}
                        <button
                          type="button"
                          on:click={() => removeTag(tag)}
                          class="ml-2 text-neon-purple/70 hover:text-neon-purple"
                          disabled={isSubmitting}
                        >
                          <X class="w-3 h-3" />
                        </button>
                      </span>
                    {/each}
                  </div>
                {:else}
                  <p class="text-gray-500 text-sm">No tags added yet</p>
                {/if}
              </div>
            </div>

            <!-- Visibility -->
            <div class="card space-y-4">
              <h3 class="text-xl font-display font-bold text-neon-cyan">
                Visibility
              </h3>

              <div class="space-y-3">
                <label class="flex items-center space-x-3 cursor-pointer">
                  <input
                    type="radio"
                    bind:group={isPublic}
                    value={true}
                    class="w-4 h-4 text-neon-pink border-gray-300 focus:ring-neon-pink"
                    disabled={isSubmitting}
                  />
                  <div>
                    <div class="font-medium text-white">Public</div>
                    <div class="text-sm text-gray-400">
                      Anyone can discover and view this node
                    </div>
                  </div>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                  <input
                    type="radio"
                    bind:group={isPublic}
                    value={false}
                    class="w-4 h-4 text-neon-pink border-gray-300 focus:ring-neon-pink"
                    disabled={isSubmitting}
                  />
                  <div>
                    <div class="font-medium text-white">Private</div>
                    <div class="text-sm text-gray-400">
                      Only you can access this node
                    </div>
                  </div>
                </label>
              </div>
            </div>

            <!-- Submit Button -->
            <div class="flex justify-end space-x-4">
              <button
                type="button"
                on:click={() => goto("/")}
                class="btn-ghost"
                disabled={isSubmitting}
              >
                Cancel
              </button>

              <button
                type="submit"
                disabled={isSubmitting}
                class="btn-primary flex items-center space-x-2 {isSubmitting
                  ? 'opacity-50 cursor-not-allowed'
                  : ''}"
              >
                {#if isSubmitting}
                  <div
                    class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"
                  />
                  <span>Creating...</span>
                {:else}
                  <Save class="w-5 h-5" />
                  <span>Create Node</span>
                {/if}
              </button>
            </div>
          </form>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Quick Tips -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-neon-green mb-4">
              💡 Quick Tips
            </h3>
            <div class="space-y-3 text-sm text-gray-300">
              <div class="flex items-start space-x-2">
                <div class="w-2 h-2 bg-neon-cyan rounded-full mt-2" />
                <div>
                  <strong>Clear titles</strong> help others find your content
                </div>
              </div>
              <div class="flex items-start space-x-2">
                <div class="w-2 h-2 bg-neon-cyan rounded-full mt-2" />
                <div>
                  <strong>Good tags</strong> make content discoverable
                </div>
              </div>
              <div class="flex items-start space-x-2">
                <div class="w-2 h-2 bg-neon-cyan rounded-full mt-2" />
                <div>
                  <strong>Rich content</strong> provides more value to learners
                </div>
              </div>
            </div>
          </div>

          <!-- Content Guidelines -->
          <div class="card">
            <h3 class="text-lg font-display font-bold text-neon-yellow mb-4">
              📋 Guidelines
            </h3>
            <div class="space-y-2 text-sm text-gray-300">
              <p>• Create original, high-quality content</p>
              <p>• Respect intellectual property rights</p>
              <p>• Use appropriate tags and categories</p>
              <p>• Keep content focused and valuable</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .bg-grid {
    background-image: linear-gradient(
        rgba(255, 0, 110, 0.1) 1px,
        transparent 1px
      ),
      linear-gradient(90deg, rgba(0, 245, 255, 0.1) 1px, transparent 1px);
    background-size: 20px 20px;
  }
</style>
