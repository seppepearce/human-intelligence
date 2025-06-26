<script>
    let trees = [];
    let loading = false;
    let apiUrl = 'http://localhost:8081';

    async function testAPI() {
        loading = true;
        try {
            const response = await fetch(`${apiUrl}/api/v1/health`);
            const data = await response.json();
            if (data.success) {
                alert('✅ Backend is working!');
            } else {
                alert('❌ Backend responded but with errors');
            }
        } catch (error) {
            alert('❌ Cannot connect to backend. Make sure it\'s running on port 8081');
        }
        loading = false;
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
            <button class="btn btn-primary" on:click={testAPI} disabled={loading}>
                {loading ? '🔄 Testing...' : '🧪 Test Backend Connection'}
            </button>
        </div>
    </section>

    <section class="features">
        <div class="feature">
            <h3>🧠 AI-Enhanced</h3>
            <p>Offline AI processing for semantic understanding</p>
        </div>
        <div class="feature">
            <h3>🏛️ Classical Design</h3>
            <p>Timeless aesthetics inspired by ancient wisdom</p>
        </div>
        <div class="feature">
            <h3>🌳 Beautiful Trees</h3>
            <p>Interactive visualizations of your knowledge</p>
        </div>
    </section>
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
    }
</style>
