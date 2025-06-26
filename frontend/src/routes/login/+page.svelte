<script>
  import { login, authLoading } from "$lib/stores/auth.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { Mail, Lock, LogIn, Eye, EyeOff, Brain } from "lucide-svelte";

  let email = "";
  let password = "";
  let showPassword = false;
  let errors = {};
  let serverError = "";
  let isSubmitting = false;

  // Form validation
  function validateForm() {
    errors = {};

    if (!email) {
      errors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "Please enter a valid email address";
    }

    if (!password) {
      errors.password = "Password is required";
    } else if (password.length < 8) {
      errors.password = "Password must be at least 8 characters";
    }

    return Object.keys(errors).length === 0;
  }

  // Handle form submission
  async function handleSubmit() {
    if (!validateForm()) return;

    isSubmitting = true;
    serverError = "";

    const result = await login({
      email: email.trim().toLowerCase(),
      password: password,
    });

    if (result.success) {
      // Redirect to dashboard or home
      goto("/");
    } else {
      serverError = result.error;
    }

    isSubmitting = false;
  }

  // Toggle password visibility
  function togglePassword() {
    showPassword = !showPassword;
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

  onMount(() => {
    // Focus on email field
    document.querySelector("#email")?.focus();
  });
</script>

<svelte:head>
  <title>Login - Human Intelligence</title>
  <meta
    name="description"
    content="Sign in to your Human Intelligence account and continue your learning journey."
  />
</svelte:head>

<div class="min-h-screen bg-dark-900 flex items-center justify-center px-4">
  <!-- Background Effects -->
  <div class="absolute inset-0 bg-grid opacity-20" />
  <div
    class="absolute inset-0 bg-gradient-to-br from-neon-pink/5 via-transparent to-neon-cyan/5"
  />

  <div class="relative w-full max-w-md">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="flex justify-center items-center space-x-3 mb-6">
        <Brain class="w-12 h-12 text-neon-pink animate-pulse" />
        <div class="text-3xl font-display font-black text-white">
          <span class="text-neon-pink">H</span><span class="text-neon-cyan"
            >I!</span
          >
        </div>
      </div>
      <h1 class="text-3xl font-display font-bold text-white mb-2">
        Welcome Back
      </h1>
      <p class="text-gray-400">Sign in to continue your learning journey</p>
    </div>

    <!-- Login Form -->
    <form on:submit|preventDefault={handleSubmit} class="card space-y-6">
      <!-- Server Error -->
      {#if serverError}
        <div class="bg-red-500/10 border border-red-500/30 rounded-lg p-4">
          <div class="flex items-center space-x-2">
            <div class="w-2 h-2 bg-red-500 rounded-full" />
            <span class="text-red-400 text-sm">{serverError}</span>
          </div>
        </div>
      {/if}

      <!-- Email Field -->
      <div class="space-y-2">
        <label for="email" class="block text-sm font-medium text-gray-300">
          Email Address
        </label>
        <div class="relative">
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <Mail class="w-5 h-5 text-gray-500" />
          </div>
          <input
            id="email"
            type="email"
            bind:value={email}
            on:input={() => clearError("email")}
            class="input-neon pl-10 {errors.email
              ? 'border-red-500 focus:border-red-500'
              : ''}"
            placeholder="your@email.com"
            disabled={isSubmitting}
            autocomplete="email"
          />
        </div>
        {#if errors.email}
          <p class="text-red-400 text-sm">{errors.email}</p>
        {/if}
      </div>

      <!-- Password Field -->
      <div class="space-y-2">
        <label for="password" class="block text-sm font-medium text-gray-300">
          Password
        </label>
        <div class="relative">
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <Lock class="w-5 h-5 text-gray-500" />
          </div>
          <input
            id="password"
            type={showPassword ? "text" : "password"}
            value={password}
            on:input={(e) => {
              password = e.target.value;
              clearError("password");
            }}
            class="input-neon pl-10 pr-10 {errors.password
              ? 'border-red-500 focus:border-red-500'
              : ''}"
            placeholder="Enter your password"
            disabled={isSubmitting}
            autocomplete="current-password"
          />
          <button
            type="button"
            on:click={togglePassword}
            class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-300"
            disabled={isSubmitting}
          >
            {#if showPassword}
              <EyeOff class="w-5 h-5" />
            {:else}
              <Eye class="w-5 h-5" />
            {/if}
          </button>
        </div>
        {#if errors.password}
          <p class="text-red-400 text-sm">{errors.password}</p>
        {/if}
      </div>

      <!-- Submit Button -->
      <button
        type="submit"
        disabled={isSubmitting || $authLoading}
        class="btn-primary w-full flex items-center justify-center space-x-2 {isSubmitting ||
        $authLoading
          ? 'opacity-50 cursor-not-allowed'
          : ''}"
      >
        {#if isSubmitting || $authLoading}
          <div
            class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"
          />
          <span>Signing In...</span>
        {:else}
          <LogIn class="w-5 h-5" />
          <span>Sign In</span>
        {/if}
      </button>

      <!-- Forgot Password Link -->
      <div class="text-center">
        <a
          href="/forgot-password"
          class="text-sm text-neon-cyan hover:text-neon-pink transition-colors"
        >
          Forgot your password?
        </a>
      </div>
    </form>

    <!-- Sign Up Link -->
    <div class="text-center mt-8">
      <p class="text-gray-400">
        Don't have an account?
        <a
          href="/register"
          class="text-neon-cyan hover:text-neon-pink font-medium transition-colors"
        >
          Sign up for free
        </a>
      </p>
    </div>

    <!-- Demo Note -->
    <div
      class="text-center mt-8 p-4 bg-dark-800 border border-dark-600 rounded-lg"
    >
      <p class="text-xs text-gray-500">
        🚀 <span class="text-neon-cyan">Demo Mode:</span> Create an account to test
        the platform!
      </p>
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
</style>
