#!/bin/bash

# 🏛️ AI Knowledge Trees - Simple Setup Script
# Classical wisdom meets modern intelligence
#
# This script sets up the project without Docker dependencies

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Unicode symbols
CHECKMARK="✅"
CROSS="❌"
ARROW="➜"
GEAR="⚙️"
TREE="🌳"
PILLAR="🏛️"
WARNING="⚠️"

# Print functions
print_header() {
    echo -e "\n${PURPLE}${PILLAR}${WHITE} $1 ${PURPLE}${PILLAR}${NC}"
    echo -e "${CYAN}$(printf '=%.0s' {1..50})${NC}"
}

print_success() {
    echo -e "${GREEN}${CHECKMARK} $1${NC}"
}

print_error() {
    echo -e "${RED}${CROSS} $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}${WARNING} $1${NC}"
}

print_info() {
    echo -e "${BLUE}${ARROW} $1${NC}"
}

print_step() {
    echo -e "\n${CYAN}${GEAR} $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main setup function
main() {
    # ASCII Art Header
    echo -e "${PURPLE}"
    cat << "EOF"
    ╔══════════════════════════════════════════════╗
    ║                                              ║
    ║     🏛️  AI KNOWLEDGE TREES SETUP 🏛️         ║
    ║                                              ║
    ║     Simple setup without Docker hassles      ║
    ║                                              ║
    ╚══════════════════════════════════════════════╝
EOF
    echo -e "${NC}\n"

    # Check prerequisites
    print_header "Checking Prerequisites"

    local all_good=true

    # Check Node.js
    if command_exists node; then
        local node_version=$(node --version)
        print_success "Node.js found: $node_version"
    else
        print_error "Node.js not found. Please install Node.js 18+ from https://nodejs.org"
        all_good=false
    fi

    # Check npm
    if command_exists npm; then
        local npm_version=$(npm --version)
        print_success "npm found: v$npm_version"
    else
        print_error "npm not found. Please install npm"
        all_good=false
    fi

    # Check Go
    if command_exists go; then
        local go_version=$(go version | awk '{print $3}')
        print_success "Go found: $go_version"
    else
        print_error "Go not found. Please install Go 1.21+ from https://golang.org"
        all_good=false
    fi

    if [ "$all_good" = false ]; then
        print_error "Please install missing prerequisites and run this script again"
        exit 1
    fi

    # Check if we're in the right directory
    if [ ! -d "backend" ] || [ ! -d "frontend" ]; then
        print_error "Please run this script from the ai-knowledge-trees project root directory"
        exit 1
    fi

    # Initialize backend
    print_header "Setting Up Backend"

    print_step "Creating Go module"
    cd backend
    if [ ! -f "go.mod" ]; then
        go mod init github.com/human-intelligence/ai-knowledge-trees/backend
        print_success "Go module created"
    else
        print_info "Go module already exists"
    fi

    print_step "Adding Go dependencies"
    go get github.com/gin-gonic/gin@v1.9.1
    go get github.com/gin-contrib/cors@v1.5.0
    go get github.com/google/uuid@v1.6.0
    go get github.com/joho/godotenv@v1.5.1
    go get modernc.org/sqlite@v1.28.0
    print_success "Go dependencies added"

    print_step "Creating backend environment file"
    cat > .env << 'EOF'
# 🏛️ AI Knowledge Trees - Backend Configuration
PORT=8081
GIN_MODE=debug
SQLITE_PATH=./ai_knowledge_trees.db
LOCALAI_URL=http://localhost:8080
CORS_ORIGINS=http://localhost:5173,http://localhost:3000
EOF
    print_success "Backend .env created"
    cd ..

    # Initialize frontend
    print_header "Setting Up Frontend"

    cd frontend
    print_step "Creating SvelteKit project"

    if [ ! -f "package.json" ]; then
        # Create basic package.json
        cat > package.json << 'EOF'
{
  "name": "ai-knowledge-trees-frontend",
  "version": "1.0.0",
  "description": "AI Knowledge Trees Frontend",
  "type": "module",
  "scripts": {
    "dev": "vite dev",
    "build": "vite build",
    "preview": "vite preview"
  },
  "devDependencies": {
    "@sveltejs/adapter-auto": "^2.1.0",
    "@sveltejs/kit": "^1.30.0",
    "svelte": "^4.2.0",
    "vite": "^4.4.0",
    "@sveltejs/vite-plugin-svelte": "^2.4.0"
  },
  "dependencies": {
    "d3": "^7.8.5"
  }
}
EOF
        print_success "Package.json created"
    fi

    print_step "Installing frontend dependencies"
    npm install --legacy-peer-deps
    print_success "Frontend dependencies installed"

    # Create basic SvelteKit structure
    print_step "Creating SvelteKit structure"

    # Create directories
    mkdir -p src/routes src/lib static

    # Create app.html
    cat > src/app.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <link rel="icon" href="%sveltekit.assets%/favicon.png" />
    <meta name="viewport" content="width=device-width" />
    <title>🏛️ AI Knowledge Trees</title>
    %sveltekit.head%
</head>
<body data-sveltekit-preload-data="hover">
    <div style="display: contents">%sveltekit.body%</div>
</body>
</html>
EOF

    # Create basic layout
    cat > src/routes/+layout.svelte << 'EOF'
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
EOF

    # Create basic home page
    cat > src/routes/+page.svelte << 'EOF'
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
EOF

    # Create basic CSS
    cat > src/routes/app.css << 'EOF'
/* 🏛️ AI Knowledge Trees - Basic Classical Styles */

* {
    box-sizing: border-box;
}

body {
    margin: 0;
    padding: 0;
    font-family: 'Georgia', serif;
    line-height: 1.6;
    color: #212529;
    background: #f8f9fa;
}

h1, h2, h3, h4, h5, h6 {
    font-family: 'Georgia', serif;
    font-weight: 600;
    line-height: 1.2;
}

p {
    margin-bottom: 1rem;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
}

.btn {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    margin: 0.25rem;
    border: none;
    border-radius: 0.25rem;
    text-decoration: none;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
}

.btn-primary {
    background: #d4af37;
    color: white;
}

.btn-primary:hover {
    background: #b8941f;
    transform: translateY(-1px);
}

.btn-secondary {
    background: white;
    color: #212529;
    border: 1px solid #6c757d;
}

.btn-secondary:hover {
    background: #f8f9fa;
    border-color: #d4af37;
}
EOF

    # Create Vite config
    cat > vite.config.js << 'EOF'
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
    plugins: [sveltekit()],
    server: {
        port: 5173,
        host: true
    }
});
EOF

    # Create SvelteKit config
    cat > svelte.config.js << 'EOF'
import adapter from '@sveltejs/adapter-auto';

const config = {
    kit: {
        adapter: adapter()
    }
};

export default config;
EOF

    # Create environment file
    cat > .env << 'EOF'
VITE_API_URL=http://localhost:8081
VITE_ENABLE_AI_FEATURES=true
EOF

    print_success "SvelteKit structure created"
    cd ..

    # Create a simple test script
    print_header "Creating Test Scripts"

    cat > test-system.sh << 'EOF'
#!/bin/bash

echo "🧪 Testing AI Knowledge Trees System..."

# Test backend build
echo "Testing backend build..."
cd backend
if go build -o bin/server .; then
    echo "✅ Backend builds successfully"
else
    echo "❌ Backend build failed"
    exit 1
fi
cd ..

# Test frontend build
echo "Testing frontend build..."
cd frontend
if npm run build; then
    echo "✅ Frontend builds successfully"
else
    echo "❌ Frontend build failed"
    exit 1
fi
cd ..

echo "🎉 All tests passed!"
EOF
    chmod +x test-system.sh

    # Create development script
    cat > dev.sh << 'EOF'
#!/bin/bash

echo "🏛️ Starting AI Knowledge Trees Development Environment..."

# Check if backend is built
if [ ! -f "backend/bin/server" ]; then
    echo "Building backend..."
    cd backend && go build -o bin/server . && cd ..
fi

# Start backend in background
echo "Starting backend on port 8081..."
cd backend && ./bin/server &
BACKEND_PID=$!
cd ..

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "Starting frontend on port 5173..."
cd frontend && npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "🎉 Development servers started!"
echo "Frontend: http://localhost:5173"
echo "Backend:  http://localhost:8081"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for user interrupt
trap 'kill $BACKEND_PID $FRONTEND_PID 2>/dev/null' EXIT
wait
EOF
    chmod +x dev.sh

    # Final success message
    print_header "Setup Complete!"

    echo -e "${GREEN}${TREE} Your AI Knowledge Trees environment is ready! ${TREE}${NC}\n"

    echo -e "${WHITE}What's been set up:${NC}"
    echo -e "${GREEN}${CHECKMARK} Go backend with SQLite database${NC}"
    echo -e "${GREEN}${CHECKMARK} SvelteKit frontend with basic UI${NC}"
    echo -e "${GREEN}${CHECKMARK} Development and test scripts${NC}"
    echo ""

    echo -e "${WHITE}Next steps:${NC}"
    echo -e "${CYAN}${ARROW} Test the system: ${YELLOW}./test-system.sh${NC}"
    echo -e "${CYAN}${ARROW} Start development: ${YELLOW}./dev.sh${NC}"
    echo -e "${CYAN}${ARROW} Access frontend: ${YELLOW}http://localhost:5173${NC}"
    echo ""

    echo -e "${YELLOW}${WARNING} Note: AI features require Docker/LocalAI setup${NC}"
    echo -e "${CYAN}${ARROW} For now, focus on basic tree functionality${NC}"
    echo ""

    echo -e "${PURPLE}${PILLAR} \"The beginning of wisdom is the definition of terms.\" - Socrates ${PILLAR}${NC}"
    echo -e "${CYAN}Happy coding in your digital garden of wisdom!${NC}\n"
}

# Run main function
main "$@"
