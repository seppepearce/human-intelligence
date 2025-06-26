<script>
    import './app.css';
</script>

<header class="header">
    <div class="container">
        <h1>🏛️ AI Knowledge Trees</h1>
        <p>Classical wisdom meets modern intelligence</p>
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
        text-align: center;
    }

    h1 {
        margin: 0;
        color: #212529;
        font-size: 2.5rem;
    }

    p {
        margin: 0.5rem 0 0 0;
        color: #6c757d;
        font-style: italic;
    }

    main {
        padding: 2rem;
    }
</style>
