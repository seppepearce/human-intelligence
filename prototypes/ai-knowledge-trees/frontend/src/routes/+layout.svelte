<script>
    import './app.css';
</script>

<header class="header">
    <div class="container">
        <div class="header-content">
            <div class="brand">
                <h1>🏛️ AI Knowledge Trees</h1>
                <p>Classical wisdom meets modern intelligence</p>
            </div>
            <nav class="nav">
                <a href="/" class="nav-link">🏠 Home</a>
                <a href="/trees" class="nav-link">🌳 Trees</a>
            </nav>
        </div>
    </div>
</header>

<main>
    <slot />
</main>

<style>
    :global(body) {
        margin: 0;
        font-family: 'Arial', sans-serif;
        background: #f8f9fa;
    }

    .header {
        background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
        border-bottom: 1px solid #d4af37;
        padding: 2rem 0;
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 1rem;
    }

    .header-content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        text-align: left;
    }

    .brand h1 {
        margin: 0;
        color: #212529;
        font-size: 2.5rem;
    }

    .brand p {
        margin: 0.5rem 0 0 0;
        color: #6c757d;
        font-style: italic;
    }

    .nav {
        display: flex;
        gap: 1rem;
    }

    .nav-link {
        text-decoration: none;
        color: #212529;
        font-weight: 600;
        padding: 0.5rem 1rem;
        border-radius: 0.25rem;
        transition: all 0.2s ease;
        border: 1px solid transparent;
    }

    .nav-link:hover {
        background: #d4af37;
        color: white;
        border-color: #b8941f;
    }

    @media (max-width: 768px) {
        .header-content {
            flex-direction: column;
            text-align: center;
            gap: 1rem;
        }

        .brand h1 {
            font-size: 2rem;
        }
    }

    main {
        padding: 2rem;
    }
</style>
