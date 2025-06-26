<script>
    import { goto } from '$app/navigation';
    import api from '$lib/api.js';

    let trees = [];
    let loading = false;
    let healthStatus = null;
    let seedStatus = null;
    let seedLoading = false;

    async function testAPI() {
        loading = true;
        try {
            const result = await api.checkHealth();
            if (result.success) {
                healthStatus = result.health;
                alert('✅ Backend is working!');
            } else {
                alert('❌ Backend responded but with errors: ' + result.error);
            }
        } catch (error) {
            alert('❌ Cannot connect to backend. Make sure it\'s running on port 8081');
        }
        loading = false;
    }

    async function checkSeedStatus() {
        try {
            const response = await fetch('http://localhost:8081/api/v1/admin/seed/status');
            const result = await response.json();
            if (result.success) {
                seedStatus = result.data;
            }
        } catch (error) {
            console.error('Error checking seed status:', error);
        }
    }

    async function seedDatabase() {
        seedLoading = true;
        try {
            const response = await fetch('http://localhost:8081/api/v1/admin/seed', {
                method: 'POST'
            });
            const result = await response.json();
            if (result.success) {
                alert('✅ Database seeded with demo data!');
                await checkSeedStatus();
            } else {
                alert('❌ Failed to seed database: ' + result.error?.message);
            }
        } catch (error) {
            alert('❌ Error seeding database. Check console for details.');
            console.error('Seed error:', error);
        }
        seedLoading = false;
    }

    async function clearSeedData() {
        if (!confirm('Are you sure you want to clear all demo data?')) return;

        seedLoading = true;
        try {
            const response = await fetch('http://localhost:8081/api/v1/admin/seed', {
                method: 'DELETE'
            });
            const result = await response.json();
            if (result.success) {
                alert('✅ Demo data cleared!');
                await checkSeedStatus();
            } else {
                alert('❌ Failed to clear demo data: ' + result.error?.message);
            }
        } catch (error) {
            alert('❌ Error clearing demo data. Check console for details.');
            console.error('Clear error:', error);
        }
        seedLoading = false;
    }

    // Check seed status on component mount
    import { onMount } from 'svelte';
    onMount(() => {
        checkSeedStatus();
    });

    function navigateToTrees() {
        goto('/trees');
    }
</script>

<svelte:head>
    <title>Home - AI Knowledge Trees</title>
</svelte:head>

<div class="container">
    <section class="hero">
        <h2>🌳 Welcome to Your Digital Garden of Wisdom</h2>
        <p>Grow beautiful knowledge trees enhanced by AI, with the timeless aesthetic of ancient Greek wisdom.</p>

        <div class="actions">
            <button class="btn btn-primary" on:click={navigateToTrees}>
                🌳 Explore Knowledge Trees
            </button>
            <button class="btn btn-secondary" on:click={testAPI} disabled={loading}>
                {loading ? '🔄 Testing...' : '🧪 Test Backend Connection'}
            </button>
        </div>

        <div class="seed-controls">
            <h4>🌱 Demo Data Controls</h4>
            {#if seedStatus}
                <div class="seed-status">
                    <span class="status-indicator" class:seeded={seedStatus.seeded}>
                        {seedStatus.seeded ? '✅ Demo data loaded' : '❌ No demo data'}
                    </span>
                </div>
                <div class="seed-actions">
                    {#if seedStatus.seeded}
                        <button
                            class="btn btn-secondary"
                            on:click={clearSeedData}
                            disabled={seedLoading}
                        >
                            {seedLoading ? '🔄 Clearing...' : '🧹 Clear Demo Data'}
                        </button>
                    {:else}
                        <button
                            class="btn btn-primary"
                            on:click={seedDatabase}
                            disabled={seedLoading}
                        >
                            {seedLoading ? '🔄 Loading...' : '🌱 Load Demo Data'}
                        </button>
                    {/if}
                    <button class="btn btn-secondary" on:click={checkSeedStatus}>
                        🔍 Check Status
                    </button>
                </div>
            {:else}
                <div class="seed-status">
                    <span class="status-indicator">🔄 Loading status...</span>
                </div>
            {/if}
        </div>
    </section>

    <section class="features">
        <div class="feature">
            <h3>🧠 AI-Enhanced</h3>
            <p>Offline AI processing for semantic understanding and analysis</p>
            <ul>
                <li>Automatic difficulty assessment</li>
                <li>Concept extraction</li>
                <li>Learning suggestions</li>
            </ul>
        </div>
        <div class="feature">
            <h3>🏛️ Classical Design</h3>
            <p>Timeless aesthetics inspired by ancient wisdom</p>
            <ul>
                <li>Marble and gold color palette</li>
                <li>Serif typography</li>
                <li>Elegant transitions</li>
            </ul>
        </div>
        <div class="feature">
            <h3>🌳 Hierarchical Trees</h3>
            <p>Organize knowledge in meaningful structures</p>
            <ul>
                <li>Parent-child relationships</li>
                <li>Drag & drop organization</li>
                <li>Visual depth indicators</li>
            </ul>
        </div>
    </section>

    {#if healthStatus}
        <section class="status-section">
            <h3>🏥 System Status</h3>
            <div class="status-grid">
                <div class="status-item">
                    <span class="status-label">Backend:</span>
                    <span class="status-value healthy">{healthStatus.status}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">Database:</span>
                    <span class="status-value healthy">{healthStatus.database}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">AI Service:</span>
                    <span class="status-value healthy">{healthStatus.localai}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">Version:</span>
                    <span class="status-value">{healthStatus.version}</span>
                </div>
            </div>
        </section>
    {/if}
</div>

<style>
    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 1rem;
    }

    .hero {
        text-align: center;
        padding: 3rem 0;
        background: white;
        border-radius: 0.5rem;
        margin-bottom: 2rem;
        border: 1px solid #d4af37;
    }

    .hero h2 {
        color: #212529;
        margin-bottom: 1rem;
    }

    .hero p {
        color: #6c757d;
        font-size: 1.1rem;
        margin-bottom: 2rem;
        max-width: 600px;
        margin-left: auto;
        margin-right: auto;
    }

    .actions {
        margin-top: 2rem;
    }

    .btn {
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 0.25rem;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-primary {
        background: #d4af37;
        color: white;
    }

    .btn-primary:hover:not(:disabled) {
        background: #b8941f;
        transform: translateY(-1px);
    }

    .btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .features {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 2rem;
        padding: 2rem 0;
    }

    .feature {
        background: white;
        padding: 2rem;
        border-radius: 0.5rem;
        text-align: center;
        border: 1px solid #e9ecef;
    }

    .feature h3 {
        color: #212529;
        margin-bottom: 1rem;
    }

    .feature p {
        color: #6c757d;
        margin-bottom: 1rem;
    }

    .feature ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .feature li {
        color: #6c757d;
        padding: 0.25rem 0;
        padding-left: 1.5rem;
        position: relative;
    }

    .feature li::before {
        content: "•";
        color: #d4af37;
        font-weight: bold;
        position: absolute;
        left: 0;
    }

    .seed-controls {
        margin-top: 2rem;
        padding: 1.5rem;
        background: #f8f9fa;
        border-radius: 0.5rem;
        border: 1px solid #d4af37;
        text-align: center;
    }

    .seed-controls h4 {
        margin: 0 0 1rem 0;
        color: #212529;
        font-size: 1.1rem;
    }

    .seed-status {
        margin-bottom: 1rem;
    }

    .status-indicator {
        display: inline-block;
        padding: 0.5rem 1rem;
        border-radius: 0.25rem;
        font-weight: 600;
        font-size: 0.9rem;
        background: #e9ecef;
        color: #6c757d;
    }

    .status-indicator.seeded {
        background: #d1edff;
        color: #0c5460;
    }

    .seed-actions {
        display: flex;
        gap: 0.5rem;
        justify-content: center;
        flex-wrap: wrap;
    }

    .status-section {
        background: white;
        border-radius: 0.5rem;
        padding: 2rem;
        margin-top: 2rem;
        border: 1px solid #d4af37;
    }

    .status-section h3 {
        text-align: center;
        color: #212529;
        margin-bottom: 1.5rem;
    }

    .status-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
    }

    .status-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 0;
        border-bottom: 1px solid #e9ecef;
    }

    .status-label {
        font-weight: 600;
        color: #6c757d;
    }

    .status-value {
        color: #212529;
        font-weight: 600;
    }

    .status-value.healthy {
        color: #28a745;
    }
</style>
