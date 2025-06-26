<script>
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Users, MessageCircle, ArrowRight, Clock, Heart, Brain, Zap } from 'lucide-svelte';

	// Mock user data - in real app this would come from auth
	const currentUser = writable({
		id: 'user_123',
		name: 'Socrates',
		avatar: '🧠',
		interests: ['Philosophy', 'AI Ethics', 'Ancient Wisdom']
	});

	// Mock friends data - in real app this would come from social graph
	const friends = writable([
		{ id: 'friend_1', name: 'Plato', avatar: '📚', currentSpace: 'Philosophy of Mind', status: 'active' },
		{ id: 'friend_2', name: 'Aristotle', avatar: '🔬', currentSpace: 'Logic & Reasoning', status: 'active' },
		{ id: 'friend_3', name: 'Hypatia', avatar: '🌟', currentSpace: 'Mathematics', status: 'active' }
	]);

	// Discussion spaces with natural constraints
	const spaces = writable([
		{
			id: 'philosophy',
			name: 'Philosophy of Mind',
			description: 'Exploring consciousness, qualia, and the hard problem',
			participants: 11,
			capacity: 12,
			activity: 'high',
			topics: ['Consciousness', 'Free Will', 'Qualia'],
			friends: ['Plato'],
			avgDepth: 8.5,
			lastActive: '2 min ago'
		},
		{
			id: 'ai-ethics',
			name: 'AI Ethics',
			description: 'Moral implications of artificial intelligence',
			participants: 12,
			capacity: 12,
			activity: 'very-high',
			topics: ['Alignment', 'Bias', 'Privacy'],
			friends: [],
			avgDepth: 7.2,
			lastActive: 'now',
			status: 'full'
		},
		{
			id: 'cognitive-science',
			name: 'Cognitive Science',
			description: 'How minds work across species and machines',
			participants: 4,
			capacity: 12,
			activity: 'medium',
			topics: ['Neuroscience', 'Psychology', 'AI'],
			friends: [],
			avgDepth: 6.8,
			lastActive: '5 min ago'
		},
		{
			id: 'ancient-wisdom',
			name: 'Ancient Wisdom',
			description: 'Timeless insights from classical philosophers',
			participants: 7,
			capacity: 12,
			activity: 'medium',
			topics: ['Stoicism', 'Virtue Ethics', 'Dialectics'],
			friends: [],
			avgDepth: 9.1,
			lastActive: '3 min ago'
		},
		{
			id: 'logic-reasoning',
			name: 'Logic & Reasoning',
			description: 'Formal and informal reasoning patterns',
			participants: 8,
			capacity: 12,
			activity: 'medium',
			topics: ['Formal Logic', 'Fallacies', 'Arguments'],
			friends: ['Aristotle'],
			avgDepth: 7.9,
			lastActive: '1 min ago'
		},
		{
			id: 'mathematics',
			name: 'Mathematics',
			description: 'The language of the universe',
			participants: 6,
			capacity: 12,
			activity: 'low',
			topics: ['Geometry', 'Number Theory', 'Infinity'],
			friends: ['Hypatia'],
			avgDepth: 8.7,
			lastActive: '4 min ago'
		},
		{
			id: 'quantum-physics',
			name: 'Quantum Physics',
			description: 'The strange world of quantum mechanics',
			participants: 3,
			capacity: 12,
			activity: 'low',
			topics: ['Superposition', 'Entanglement', 'Measurement'],
			friends: [],
			avgDepth: 5.4,
			lastActive: '8 min ago'
		}
	]);

	// Reactive suggestions based on crowding and interests
	$: suggestions = $spaces
		.filter(space => space.status !== 'full')
		.filter(space => {
			const userInterests = $currentUser.interests;
			return space.topics.some(topic =>
				userInterests.some(interest =>
					topic.toLowerCase().includes(interest.toLowerCase()) ||
					interest.toLowerCase().includes(topic.toLowerCase())
				)
			);
		})
		.sort((a, b) => {
			// Prioritize spaces with friends
			const aHasFriends = a.friends.length > 0 ? 1 : 0;
			const bHasFriends = b.friends.length > 0 ? 1 : 0;
			if (aHasFriends !== bHasFriends) return bHasFriends - aHasFriends;

			// Then by activity level
			const activityScore = { 'very-high': 4, 'high': 3, 'medium': 2, 'low': 1 };
			return activityScore[b.activity] - activityScore[a.activity];
		});

	// Cross-pollination opportunities
	$: crossPollination = $spaces
		.filter(space => !$currentUser.interests.some(interest =>
			space.topics.some(topic =>
				topic.toLowerCase().includes(interest.toLowerCase()) ||
				interest.toLowerCase().includes(topic.toLowerCase())
			)
		))
		.filter(space => space.participants < space.capacity * 0.7) // Not too crowded
		.slice(0, 3);

	// Simulate real-time updates
	onMount(() => {
		const interval = setInterval(() => {
			spaces.update(currentSpaces => {
				return currentSpaces.map(space => {
					// Random small changes to simulate activity
					const change = Math.random() > 0.7 ? (Math.random() > 0.5 ? 1 : -1) : 0;
					const newParticipants = Math.max(1, Math.min(space.capacity, space.participants + change));

					return {
						...space,
						participants: newParticipants,
						status: newParticipants >= space.capacity ? 'full' : undefined,
						lastActive: change !== 0 ? 'now' : space.lastActive
					};
				});
			});
		}, 3000);

		return () => clearInterval(interval);
	});

	function getActivityColor(activity, isFull = false) {
		if (isFull) return 'border-red-500 bg-red-500/10';
		switch (activity) {
			case 'very-high': return 'border-green-400 bg-green-400/10';
			case 'high': return 'border-yellow-400 bg-yellow-400/10';
			case 'medium': return 'border-blue-400 bg-blue-400/10';
			case 'low': return 'border-gray-400 bg-gray-400/10';
			default: return 'border-gray-500 bg-gray-500/10';
		}
	}

	function getCapacityBarColor(percentage) {
		if (percentage >= 100) return 'bg-red-500';
		if (percentage >= 80) return 'bg-yellow-500';
		if (percentage >= 60) return 'bg-blue-500';
		return 'bg-green-500';
	}

	function joinSpace(spaceId) {
		// In real app, this would connect to WebSocket and join the space
		console.log('Joining space:', spaceId);
		// For demo, just show an alert
		alert(`Joining ${$spaces.find(s => s.id === spaceId)?.name}...`);
	}

	function exploreCrossPollination(spaceId) {
		// In real app, this would track cross-pollination events
		console.log('Cross-pollination opportunity:', spaceId);
		joinSpace(spaceId);
	}
</script>

<svelte:head>
	<title>Digital Agora - The Plaza</title>
	<meta name="description" content="Welcome to the Digital Agora - where meaningful discourse flows like ancient streams of wisdom." />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-agora-night via-agora-deep to-agora-night">
	<!-- Agora Header -->
	<header class="border-b border-agora-gold/20 bg-agora-night/50 backdrop-blur-sm">
		<div class="container mx-auto px-6 py-4">
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-3">
					<div class="text-2xl">🏛️</div>
					<div>
						<h1 class="text-xl font-bold text-agora-gold">Digital Agora</h1>
						<p class="text-sm text-agora-stone/70">The Plaza of Ideas</p>
					</div>
				</div>

				<div class="flex items-center space-x-4">
					<div class="text-sm text-agora-stone/70">
						Welcome, <span class="text-agora-gold font-medium">{$currentUser.name}</span>
					</div>
					<div class="w-8 h-8 rounded-full bg-agora-gold/20 flex items-center justify-center text-lg">
						{$currentUser.avatar}
					</div>
				</div>
			</div>
		</div>
	</header>

	<div class="container mx-auto px-6 py-8">
		<!-- Current Activity Overview -->
		<div class="mb-8">
			<h2 class="text-2xl font-bold text-agora-gold mb-4 flex items-center">
				<Zap class="w-6 h-6 mr-2" />
				The Agora Today
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
				<div class="bg-agora-deep/50 rounded-lg p-4 border border-agora-gold/20">
					<div class="text-2xl font-bold text-agora-gold">{$spaces.reduce((sum, s) => sum + s.participants, 0)}</div>
					<div class="text-sm text-agora-stone/70">Active Participants</div>
				</div>
				<div class="bg-agora-deep/50 rounded-lg p-4 border border-agora-gold/20">
					<div class="text-2xl font-bold text-agora-gold">{$spaces.length}</div>
					<div class="text-sm text-agora-stone/70">Discussion Spaces</div>
				</div>
				<div class="bg-agora-deep/50 rounded-lg p-4 border border-agora-gold/20">
					<div class="text-2xl font-bold text-agora-gold">{$friends.filter(f => f.status === 'active').length}</div>
					<div class="text-sm text-agora-stone/70">Friends Online</div>
				</div>
				<div class="bg-agora-deep/50 rounded-lg p-4 border border-agora-gold/20">
					<div class="text-2xl font-bold text-agora-gold">{$spaces.filter(s => s.status === 'full').length}</div>
					<div class="text-sm text-agora-stone/70">Crowded Spaces</div>
				</div>
			</div>
		</div>

		<!-- Friend Network Activity -->
		{#if $friends.some(f => f.status === 'active')}
		<div class="mb-8">
			<h3 class="text-lg font-semibold text-agora-gold mb-4 flex items-center">
				<Heart class="w-5 h-5 mr-2" />
				Your Circle's Wanderings
			</h3>
			<div class="flex flex-wrap gap-3">
				{#each $friends.filter(f => f.status === 'active') as friend}
				<div class="bg-agora-deep/30 rounded-full px-4 py-2 border border-agora-gold/20 flex items-center space-x-2">
					<span class="text-lg">{friend.avatar}</span>
					<span class="text-sm text-agora-stone">{friend.name} is exploring</span>
					<span class="text-sm font-medium text-agora-gold">{friend.currentSpace}</span>
					<button
						class="text-xs text-agora-neon-cyan hover:text-agora-neon-purple transition-colors"
						on:click={() => joinSpace(friend.currentSpace.toLowerCase().replace(/ /g, '-'))}
					>
						Join →
					</button>
				</div>
				{/each}
			</div>
		</div>
		{/if}

		<!-- Discussion Spaces Grid -->
		<div class="mb-8">
			<h3 class="text-lg font-semibold text-agora-gold mb-4 flex items-center">
				<MessageCircle class="w-5 h-5 mr-2" />
				Active Discussion Spaces
			</h3>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each $spaces as space}
				<div class="bg-agora-deep/50 rounded-lg p-6 border {getActivityColor(space.activity, space.status === 'full')} hover:border-agora-gold/50 transition-all duration-300 hover:shadow-lg hover:shadow-agora-gold/10">
					<!-- Space Header -->
					<div class="flex items-start justify-between mb-4">
						<div class="flex-1">
							<h4 class="font-semibold text-agora-marble mb-1">{space.name}</h4>
							<p class="text-sm text-agora-stone/70 leading-relaxed">{space.description}</p>
						</div>
						<div class="ml-4 text-right">
							<div class="text-sm font-medium text-agora-gold">
								{space.participants}/{space.capacity}
							</div>
							<div class="text-xs text-agora-stone/60">participants</div>
						</div>
					</div>

					<!-- Capacity Bar -->
					<div class="w-full bg-agora-night/50 rounded-full h-2 mb-4">
						<div
							class="h-2 rounded-full transition-all duration-300 {getCapacityBarColor((space.participants / space.capacity) * 100)}"
							style="width: {Math.min(100, (space.participants / space.capacity) * 100)}%"
						></div>
					</div>

					<!-- Space Stats -->
					<div class="flex items-center justify-between text-xs text-agora-stone/60 mb-4">
						<div class="flex items-center space-x-4">
							<span class="flex items-center">
								<Clock class="w-3 h-3 mr-1" />
								{space.lastActive}
							</span>
							<span class="flex items-center">
								<Brain class="w-3 h-3 mr-1" />
								{space.avgDepth}/10 depth
							</span>
						</div>
						{#if space.friends.length > 0}
						<div class="text-agora-neon-cyan">
							Friends: {space.friends.join(', ')}
						</div>
						{/if}
					</div>

					<!-- Topics -->
					<div class="flex flex-wrap gap-2 mb-4">
						{#each space.topics as topic}
						<span class="bg-agora-night/50 text-agora-stone/80 px-2 py-1 rounded-full text-xs">
							{topic}
						</span>
						{/each}
					</div>

					<!-- Action Button -->
					<button
						class="w-full py-2 px-4 rounded-lg font-medium transition-all duration-200 {
							space.status === 'full'
								? 'bg-red-500/20 text-red-300 border border-red-500/30 cursor-not-allowed'
								: 'bg-agora-gold/20 text-agora-gold border border-agora-gold/30 hover:bg-agora-gold/30 hover:shadow-md'
						}"
						disabled={space.status === 'full'}
						on:click={() => joinSpace(space.id)}
					>
						{space.status === 'full' ? 'Space Full' : 'Join Discussion'}
					</button>

					{#if space.status === 'full'}
					<p class="text-xs text-agora-stone/60 mt-2 text-center">
						This space is at capacity. Try a related space below!
					</p>
					{/if}
				</div>
				{/each}
			</div>
		</div>

		<!-- Suggestions Section -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Recommended Spaces -->
			<div>
				<h3 class="text-lg font-semibold text-agora-gold mb-4 flex items-center">
					<ArrowRight class="w-5 h-5 mr-2" />
					Recommended for You
				</h3>
				<div class="space-y-3">
					{#each suggestions.slice(0, 3) as space}
					<div class="bg-agora-deep/30 rounded-lg p-4 border border-agora-gold/20 hover:border-agora-gold/40 transition-all">
						<div class="flex items-center justify-between">
							<div class="flex-1">
								<h5 class="font-medium text-agora-marble">{space.name}</h5>
								<p class="text-sm text-agora-stone/70">{space.participants}/{space.capacity} participants</p>
								{#if space.friends.length > 0}
								<p class="text-xs text-agora-neon-cyan mt-1">Friends: {space.friends.join(', ')}</p>
								{/if}
							</div>
							<button
								class="px-3 py-1 bg-agora-gold/20 text-agora-gold rounded-md text-sm hover:bg-agora-gold/30 transition-colors"
								on:click={() => joinSpace(space.id)}
							>
								Join
							</button>
						</div>
					</div>
					{/each}
				</div>
			</div>

			<!-- Cross-Pollination Opportunities -->
			<div>
				<h3 class="text-lg font-semibold text-agora-gold mb-4 flex items-center">
					<Zap class="w-5 h-5 mr-2" />
					Explore New Territory
				</h3>
				<p class="text-sm text-agora-stone/70 mb-4">
					Venture beyond your usual domains. Great insights often come from unexpected connections.
				</p>
				<div class="space-y-3">
					{#each crossPollination as space}
					<div class="bg-agora-deep/30 rounded-lg p-4 border border-agora-neon-purple/20 hover:border-agora-neon-purple/40 transition-all">
						<div class="flex items-center justify-between">
							<div class="flex-1">
								<h5 class="font-medium text-agora-marble">{space.name}</h5>
								<p class="text-sm text-agora-stone/70">{space.participants}/{space.capacity} participants • {space.avgDepth}/10 depth</p>
								<div class="flex flex-wrap gap-1 mt-2">
									{#each space.topics.slice(0, 2) as topic}
									<span class="bg-agora-neon-purple/20 text-agora-neon-purple px-2 py-0.5 rounded-full text-xs">
										{topic}
									</span>
									{/each}
								</div>
							</div>
							<button
								class="px-3 py-1 bg-agora-neon-purple/20 text-agora-neon-purple rounded-md text-sm hover:bg-agora-neon-purple/30 transition-colors"
								on:click={() => exploreCrossPollination(space.id)}
							>
								Explore
							</button>
						</div>
					</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	:global(.agora-app) {
		--agora-marble: #f8f9fa;
		--agora-stone: #e9ecef;
		--agora-shadow: #495057;
		--agora-night: #1a1a2e;
		--agora-deep: #16213e;
		--agora-gold: #d4af37;
		--agora-bronze: #cd7f32;
		--agora-crimson: #dc143c;
		--agora-olive: #808000;
		--agora-neon-cyan: #00f5ff;
		--agora-neon-purple: #8b5cf6;
		--agora-neon-pink: #ff006e;
	}

	.container {
		max-width: 1200px;
	}
</style>
