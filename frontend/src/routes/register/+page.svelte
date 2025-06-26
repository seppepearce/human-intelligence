<script>
  import { register, authLoading } from "$lib/stores/auth.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    Mail,
    Lock,
    User,
    FileText,
    UserPlus,
    Eye,
    EyeOff,
    Brain,
    CheckCircle,
  } from "lucide-svelte";

  let username = "";
  let email = "";
  let password = "";
  let confirmPassword = "";
  let bio = "";
  let showPassword = false;
  let showConfirmPassword = false;
  let errors = {};
  let serverError = "";
  let isSubmitting = false;

  // Password strength indicator
  $: passwordStrength = calculatePasswordStrength(password);

  function calculatePasswordStrength(pwd) {
    if (!pwd) return { score: 0, label: "", color: "" };

    let score = 0;
    if (pwd.length >= 8) score++;
    if (/[a-z]/.test(pwd)) score++;
    if (/[A-Z]/.test(pwd)) score++;
    if (/[0-9]/.test(pwd)) score++;
    if (/[^A-Za-z0-9]/.test(pwd)) score++;

    const levels = [
      { score: 0, label: "", color: "" },
      { score: 1, label: "Very Weak", color: "bg-red-500" },
      { score: 2, label: "Weak", color: "bg-orange-500" },
      { score: 3, label: "Fair", color: "bg-yellow-500" },
      { score: 4, label: "Good", color: "bg-blue-500" },
      { score: 5, label: "Strong", color: "bg-green-500" },
    ];

    return levels[score] || levels[0];
  }

  // Form validation
  function validateForm() {
    errors = {};

    if (!username) {
      errors.username = "Username is required";
    } else if (username.length < 3) {
      errors.username = "Username must be at least 3 characters";
    } else if (username.length > 50) {
      errors.username = "Username must be less than 50 characters";
    } else if (!/^[a-zA-Z0-9_-]+$/.test(username)) {
      errors.username =
        "Username can only contain letters, numbers, hyphens, and underscores";
    }

    if (!email) {
      errors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "Please enter a valid email address";
    }

    if (!password) {
      errors.password = "Password is required";
    } else if (password.length < 8) {
      errors.password = "Password must be at least 8 characters";
    } else if (passwordStrength.score < 3) {
      errors.password =
        "Password is too weak. Use a mix of letters, numbers, and symbols";
    }

    if (!confirmPassword) {
      errors.confirmPassword = "Please confirm your password";
    } else if (password !== confirmPassword) {
      errors.confirmPassword = "Passwords do not match";
    }

    if (bio && bio.length > 500) {
      errors.bio = "Bio must be less than 500 characters";
    }

    return Object.keys(errors).length === 0;
  }

  // Handle form submission
  async function handleSubmit() {
    if (!validateForm()) return;

    isSubmitting = true;
    serverError = "";

    const result = await register({
      username: username.trim(),
      email: email.trim().toLowerCase(),
      password: password,
      bio: bio.trim(),
    });

    if (result.success) {
      // Redirect to dashboard or welcome page
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

  function toggleConfirmPassword() {
    showConfirmPassword = !showConfirmPassword;
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
    // Focus on username field
    document.querySelector("#username")?.focus();
  });
</script>

<svelte:head>
  <title>Sign Up - Human Intelligence</title>
  <meta
    name="description"
    content="Join Human Intelligence and start your collaborative learning journey with AI augmentation."
  />
</svelte:head>

<div
  class="min-h-screen bg-dark-900 flex items-center justify-center px-4 py-8"
>
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
        Join the Community
      </h1>
      <p class="text-gray-400">Where AI augments Human Intelligence</p>
    </div>

    <!-- Registration Form -->
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

      <!-- Username Field -->
      <div class="space-y-2">
        <label for="username" class="block text-sm font-medium text-gray-300">
          Username
        </label>
        <div class="relative">
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <User class="w-5 h-5 text-gray-500" />
          </div>
          <input
            id="username"
            type="text"
            bind:value={username}
            on:input={() => clearError("username")}
            class="input-neon pl-10 {errors.username
              ? 'border-red-500 focus:border-red-500'
              : ''}"
            placeholder="your_username"
            disabled={isSubmitting}
            autocomplete="username"
          />
        </div>
        {#if errors.username}
          <p class="text-red-400 text-sm">{errors.username}</p>
        {:else}
          <p class="text-gray-500 text-xs">
            3-50 characters, letters, numbers, hyphens, and underscores only
          </p>
        {/if}
      </div>

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
            placeholder="Create a strong password"
            disabled={isSubmitting}
            autocomplete="new-password"
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

        <!-- Password Strength Indicator -->
        {#if password && passwordStrength.score > 0}
          <div class="flex items-center space-x-2">
            <div class="flex-1 bg-dark-700 rounded-full h-2">
              <div
                class="h-2 rounded-full transition-all duration-300 {passwordStrength.color}"
                style="width: {(passwordStrength.score / 5) * 100}%"
              />
            </div>
            <span class="text-xs text-gray-400">{passwordStrength.label}</span>
          </div>
        {/if}

        {#if errors.password}
          <p class="text-red-400 text-sm">{errors.password}</p>
        {:else}
          <p class="text-gray-500 text-xs">
            At least 8 characters with mix of letters, numbers, and symbols
          </p>
        {/if}
      </div>

      <!-- Confirm Password Field -->
      <div class="space-y-2">
        <label
          for="confirmPassword"
          class="block text-sm font-medium text-gray-300"
        >
          Confirm Password
        </label>
        <div class="relative">
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <Lock class="w-5 h-5 text-gray-500" />
          </div>
          <input
            id="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            value={confirmPassword}
            on:input={(e) => {
              confirmPassword = e.target.value;
              clearError("confirmPassword");
            }}
            class="input-neon pl-10 pr-10 {errors.confirmPassword
              ? 'border-red-500 focus:border-red-500'
              : ''}"
            placeholder="Confirm your password"
            disabled={isSubmitting}
            autocomplete="new-password"
          />
          <button
            type="button"
            on:click={toggleConfirmPassword}
            class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-300"
            disabled={isSubmitting}
          >
            {#if showConfirmPassword}
              <EyeOff class="w-5 h-5" />
            {:else}
              <Eye class="w-5 h-5" />
            {/if}
          </button>
        </div>
        {#if confirmPassword && password === confirmPassword}
          <div class="flex items-center space-x-1 text-green-400">
            <CheckCircle class="w-4 h-4" />
            <span class="text-xs">Passwords match</span>
          </div>
        {/if}
        {#if errors.confirmPassword}
          <p class="text-red-400 text-sm">{errors.confirmPassword}</p>
        {/if}
      </div>

      <!-- Bio Field -->
      <div class="space-y-2">
        <label for="bio" class="block text-sm font-medium text-gray-300">
          Bio <span class="text-gray-500">(optional)</span>
        </label>
        <div class="relative">
          <div
            class="absolute top-3 left-0 pl-3 flex items-start pointer-events-none"
          >
            <FileText class="w-5 h-5 text-gray-500" />
          </div>
          <textarea
            id="bio"
            bind:value={bio}
            on:input={() => clearError("bio")}
            class="input-neon pl-10 min-h-[80px] resize-none {errors.bio
              ? 'border-red-500 focus:border-red-500'
              : ''}"
            placeholder="Tell us a bit about yourself and your learning interests..."
            disabled={isSubmitting}
            maxlength="500"
          />
        </div>
        <div class="flex justify-between items-center">
          {#if errors.bio}
            <p class="text-red-400 text-sm">{errors.bio}</p>
          {:else}
            <p class="text-gray-500 text-xs">
              Share your interests, background, or learning goals
            </p>
          {/if}
          <span class="text-xs text-gray-500">{bio.length}/500</span>
        </div>
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
          <span>Creating Account...</span>
        {:else}
          <UserPlus class="w-5 h-5" />
          <span>Create Account</span>
        {/if}
      </button>

      <!-- Terms Note -->
      <div class="text-center">
        <p class="text-xs text-gray-500">
          By creating an account, you agree to our
          <a href="/terms" class="text-neon-cyan hover:text-neon-pink"
            >Terms of Service</a
          >
          and
          <a href="/privacy" class="text-neon-cyan hover:text-neon-pink"
            >Privacy Policy</a
          >
        </p>
      </div>
    </form>

    <!-- Sign In Link -->
    <div class="text-center mt-8">
      <p class="text-gray-400">
        Already have an account?
        <a
          href="/login"
          class="text-neon-cyan hover:text-neon-pink font-medium transition-colors"
        >
          Sign in here
        </a>
      </p>
    </div>

    <!-- Demo Note -->
    <div
      class="text-center mt-8 p-4 bg-dark-800 border border-dark-600 rounded-lg"
    >
      <p class="text-xs text-gray-500">
        🚀 <span class="text-neon-cyan">Demo Mode:</span> All data is for testing
        purposes only
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
