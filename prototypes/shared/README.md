# 🎨 Shared Assets & Requirements for Human Intelligence Prototypes

> **Common Foundation**: Ensure all prototypes implement the same core features for fair comparison

## 🎯 Purpose

This directory contains shared assets, requirements, and specifications that **all prototypes must implement** to ensure fair comparison. Think of this as the "spec sheet" that every prototype needs to meet.

## 📋 Core Feature Requirements

### MVP Feature Set (All Prototypes Must Implement)

#### 1. **User Authentication** 🔐
```yaml
Features:
  - User registration with email/password
  - User login with session management
  - Password reset functionality
  - Basic user profile (username, avatar)
  - Logout functionality

Technical Requirements:
  - JWT or session-based authentication
  - Password hashing (bcrypt or similar)
  - Email validation
  - Session persistence across browser sessions
  - Secure password reset flow

Success Criteria:
  - Register → Login → Use app → Logout → Login again
  - Password reset via email (can be simulated)
  - Profile information persists
```

#### 2. **Learning Nodes** 📝
```yaml
Features:
  - Create text-based learning nodes
  - Edit existing nodes (if owner)
  - Delete nodes (if owner)
  - View public nodes
  - Rich text formatting (bold, italic, links)
  - Basic file upload (images)

Data Model:
  - id: unique identifier
  - title: string, required
  - content: rich text/JSON
  - type: text|video|quiz|code
  - author_id: foreign key to user
  - created_at: timestamp
  - updated_at: timestamp
  - visibility: public|private

Success Criteria:
  - Create node with title and content
  - Edit and save changes
  - Upload and display images
  - View other users' public nodes
```

#### 3. **Learning Paths** 🛤️
```yaml
Features:
  - Create learning paths (sequences of nodes)
  - Add/remove nodes from paths
  - Reorder nodes in paths
  - View and follow learning paths
  - Mark nodes as completed
  - Fork existing paths (create copy)

Data Model:
  - id: unique identifier
  - title: string, required
  - description: text, optional
  - node_ids: array of node IDs (ordered)
  - author_id: foreign key to user
  - created_at: timestamp
  - forked_from: optional foreign key

Success Criteria:
  - Create path with 3+ nodes
  - Reorder nodes via drag-and-drop or similar
  - Follow path, mark nodes complete
  - Fork someone else's path
```

#### 4. **Real-time Activity Feed** ⚡
```yaml
Features:
  - Show recent platform activity
  - Real-time updates (WebSocket, SSE, or polling)
  - Activity types: node created, path completed, user joined
  - User avatars and timestamps
  - Live user count/presence

Activity Types:
  - node_created: User created a new node
  - path_created: User created a new learning path
  - path_completed: User completed a learning path
  - node_liked: User liked a node (if implemented)

Success Criteria:
  - Activity appears in real-time across browser tabs
  - Show last 20 activities on page load
  - Updates without page refresh
  - Proper timestamps and user attribution
```

#### 5. **AI Integration** 🤖
```yaml
Features:
  - Generate TLDR summaries for nodes (140 chars max)
  - AI-powered content suggestions
  - Basic semantic search (nice to have)

Requirements:
  - Use OpenAI API (gpt-3.5-turbo minimum)
  - Handle API failures gracefully
  - Show loading states during generation
  - Cache generated content to avoid re-calling API

Success Criteria:
  - Generate TLDR for a node with 200+ words
  - Handle API errors without breaking UI
  - TLDR appears in under 5 seconds
  - Generated content is relevant and useful
```

#### 6. **Vaporwave UI** 🌈
```yaml
Design Requirements:
  - Implement provided color palette
  - Smooth animations and hover effects
  - Mobile-responsive design (works on phones)
  - Dark theme with neon accents
  - Retro/cyberpunk aesthetic

Success Criteria:
  - Matches provided design mockups
  - Smooth 60fps animations
  - Usable on mobile devices
  - Consistent visual hierarchy
  - Accessible (keyboard navigation, screen readers)
```

## 🎨 Design System

### Color Palette
```css
/* Primary Colors */
--neon-pink: #FF006E;      /* Primary actions, highlights */
--neon-cyan: #00F5FF;      /* Secondary actions, links */
--neon-purple: #8B5CF6;    /* Accents, special elements */
--neon-green: #39FF14;     /* Success states, online indicators */
--neon-yellow: #FFFF00;    /* Warnings, notifications */

/* Background Colors */
--dark-950: #050505;       /* Deepest background */
--dark-900: #0a0a0a;       /* Main background */
--dark-800: #1a1a1a;       /* Card backgrounds */
--dark-700: #2a2a2a;       /* Elevated surfaces */
--dark-600: #3a3a3a;       /* Borders */

/* Gradients */
--gradient-primary: linear-gradient(135deg, #FF006E 0%, #8B5CF6 100%);
--gradient-secondary: linear-gradient(135deg, #00F5FF 0%, #39FF14 100%);
--gradient-background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
```

### Typography
```css
/* Font Stack */
--font-primary: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
--font-mono: 'JetBrains Mono', 'Monaco', 'Consolas', monospace;

/* Font Sizes */
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
--text-4xl: 2.25rem;   /* 36px */
```

### Component Styles
```css
/* Button Variants */
.btn-primary {
  background: var(--gradient-primary);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(255, 0, 110, 0.4);
}

.btn-secondary {
  background: transparent;
  color: var(--neon-cyan);
  border: 2px solid var(--neon-cyan);
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  transition: all 0.3s ease;
}

.btn-secondary:hover {
  background: var(--neon-cyan);
  color: var(--dark-900);
  box-shadow: 0 10px 25px rgba(0, 245, 255, 0.4);
}

/* Cards */
.card {
  background: var(--dark-800);
  border: 1px solid var(--dark-600);
  border-radius: 0.75rem;
  padding: 1.5rem;
  transition: all 0.3s ease;
}

.card:hover {
  border-color: var(--neon-cyan);
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

/* Animations */
@keyframes neon-pulse {
  0% { text-shadow: 0 0 5px currentColor; }
  100% { text-shadow: 0 0 20px currentColor, 0 0 30px currentColor; }
}

.neon-glow {
  animation: neon-pulse 2s ease-in-out infinite alternate;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

.float {
  animation: float 6s ease-in-out infinite;
}
```

## 📊 Database Schema

### Common Tables (Implement in your chosen database)

#### Users Table
```sql
-- Users (Authentication)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP DEFAULT NOW()
);
```

#### Nodes Table
```sql
-- Learning Nodes
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    content JSONB NOT NULL, -- Store rich content
    type VARCHAR(20) DEFAULT 'text', -- text|video|quiz|code
    author_id UUID REFERENCES users(id) ON DELETE CASCADE,
    visibility VARCHAR(20) DEFAULT 'public', -- public|private
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_nodes_author ON nodes(author_id);
CREATE INDEX idx_nodes_created ON nodes(created_at);
CREATE INDEX idx_nodes_visibility ON nodes(visibility);
```

#### Learning Paths Table
```sql
-- Learning Paths
CREATE TABLE learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    node_ids UUID[] DEFAULT '{}', -- Ordered array of node IDs
    author_id UUID REFERENCES users(id) ON DELETE CASCADE,
    forked_from UUID REFERENCES learning_paths(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_paths_author ON learning_paths(author_id);
CREATE INDEX idx_paths_created ON learning_paths(created_at);
```

#### Activities Table
```sql
-- Activity Feed
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL, -- node_created, path_completed, etc.
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(50), -- node|path|user
    target_id UUID, -- ID of the target object
    target_title VARCHAR(200), -- Cached title for display
    metadata JSONB DEFAULT '{}', -- Additional activity data
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_activities_created ON activities(created_at);
CREATE INDEX idx_activities_user ON activities(user_id);
CREATE INDEX idx_activities_type ON activities(type);
```

#### User Progress Table
```sql
-- Track user progress through learning paths
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    path_id UUID REFERENCES learning_paths(id) ON DELETE CASCADE,
    completed_nodes UUID[] DEFAULT '{}', -- Array of completed node IDs
    current_node_id UUID, -- Currently active node
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    UNIQUE(user_id, path_id)
);

CREATE INDEX idx_progress_user ON user_progress(user_id);
CREATE INDEX idx_progress_path ON user_progress(path_id);
```

## 🎭 Sample Data

### Test Users
```json
{
  "users": [
    {
      "username": "alice_learns",
      "email": "alice@example.com",
      "password": "password123"
    },
    {
      "username": "bob_codes",
      "email": "bob@example.com", 
      "password": "password123"
    },
    {
      "username": "charlie_ai",
      "email": "charlie@example.com",
      "password": "password123"
    }
  ]
}
```

### Sample Learning Nodes
```json
{
  "nodes": [
    {
      "title": "Introduction to React Hooks",
      "content": {
        "text": "React Hooks are functions that let you 'hook into' React state and lifecycle features from function components. They were introduced in React 16.8 as a way to use state and other React features without writing a class component.\n\nThe most commonly used hooks are:\n- useState: For managing local state\n- useEffect: For side effects and lifecycle events\n- useContext: For consuming React context\n- useReducer: For complex state management\n\nHooks follow two main rules:\n1. Only call hooks at the top level of your function\n2. Only call hooks from React functions or custom hooks"
      },
      "type": "text",
      "author": "alice_learns"
    },
    {
      "title": "Building Your First API",
      "content": {
        "text": "Creating a REST API is a fundamental skill for backend developers. In this node, we'll explore the key concepts:\n\n**What is an API?**\nAn Application Programming Interface (API) is a set of rules and protocols that allows different software applications to communicate with each other.\n\n**REST Principles:**\n- Stateless: Each request contains all necessary information\n- Resource-based: URLs represent resources\n- HTTP methods: GET, POST, PUT, DELETE for different operations\n- JSON format: Standard data exchange format\n\n**Best Practices:**\n- Use meaningful HTTP status codes\n- Version your API (v1, v2, etc.)\n- Implement proper error handling\n- Add authentication and rate limiting\n- Document your endpoints thoroughly"
      },
      "type": "text",
      "author": "bob_codes"
    },
    {
      "title": "Understanding Machine Learning Basics",
      "content": {
        "text": "Machine Learning (ML) is a subset of artificial intelligence that enables computers to learn and make decisions from data without being explicitly programmed for every scenario.\n\n**Types of Machine Learning:**\n\n1. **Supervised Learning**: Learning with labeled training data\n   - Examples: Email spam detection, image classification\n   - Algorithms: Decision trees, neural networks, support vector machines\n\n2. **Unsupervised Learning**: Finding patterns in data without labels\n   - Examples: Customer segmentation, anomaly detection\n   - Algorithms: K-means clustering, principal component analysis\n\n3. **Reinforcement Learning**: Learning through rewards and penalties\n   - Examples: Game playing, robotics, recommendation systems\n   - Key concept: Agent learns optimal actions through trial and error\n\n**Getting Started:**\n- Learn Python and basic statistics\n- Understand data preprocessing\n- Practice with datasets from Kaggle\n- Start with scikit-learn library for beginners"
      },
      "type": "text",
      "author": "charlie_ai"
    }
  ]
}
```

### Sample Learning Paths
```json
{
  "learning_paths": [
    {
      "title": "Frontend Development Fundamentals",
      "description": "Master the basics of modern frontend development with React, TypeScript, and best practices.",
      "nodes": ["intro-to-react-hooks", "typescript-basics", "state-management"],
      "author": "alice_learns"
    },
    {
      "title": "Backend API Development",
      "description": "Learn to build robust, scalable APIs from scratch using modern backend technologies.",
      "nodes": ["building-first-api", "database-design", "authentication-security"],
      "author": "bob_codes"
    },
    {
      "title": "AI/ML for Beginners",
      "description": "Start your journey into artificial intelligence and machine learning with practical examples.",
      "nodes": ["ml-basics", "python-for-data-science", "first-ml-model"],
      "author": "charlie_ai"
    }
  ]
}
```

## 🔧 API Requirements

### RESTful Endpoints (All Prototypes Should Support)

#### Authentication
```
POST /api/auth/register
POST /api/auth/login
POST /api/auth/logout
POST /api/auth/refresh
GET  /api/auth/me
```

#### Users
```
GET    /api/users              # List public users
GET    /api/users/:id          # Get user profile
PUT    /api/users/:id          # Update own profile
POST   /api/users/:id/avatar   # Upload avatar
```

#### Nodes
```
GET    /api/nodes              # List public nodes
POST   /api/nodes              # Create node
GET    /api/nodes/:id          # Get node
PUT    /api/nodes/:id          # Update node (if owner)
DELETE /api/nodes/:id          # Delete node (if owner)
POST   /api/nodes/:id/tldr     # Generate AI TLDR
```

#### Learning Paths
```
GET    /api/paths              # List public paths
POST   /api/paths              # Create path
GET    /api/paths/:id          # Get path
PUT    /api/paths/:id          # Update path (if owner)
DELETE /api/paths/:id          # Delete path (if owner)
POST   /api/paths/:id/fork     # Fork path
POST   /api/paths/:id/progress # Update progress
```

#### Activities
```
GET    /api/activities         # Get activity feed
POST   /api/activities         # Create activity (system use)
```

#### Real-time Events
```
WebSocket: /ws/activities      # Real-time activity updates
WebSocket: /ws/presence        # User presence tracking
```

### Response Formats
```json
// Success Response
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "1.0"
  }
}

// Error Response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Username is required",
    "details": { ... }
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "1.0"
  }
}
```

## 📱 UI/UX Requirements

### Page Structure (All Prototypes Should Have)
```
/                     # Home page with hero, features, CTA
/login               # User login form
/register            # User registration form
/dashboard           # User dashboard with activity feed
/nodes               # Browse public nodes
/nodes/create        # Create new node
/nodes/:id           # View/edit specific node
/paths               # Browse learning paths
/paths/create        # Create new learning path
/paths/:id           # View/follow learning path
/profile/:username   # User profile page
/settings            # User settings (if logged in)
```

### Mobile Responsiveness
- Works on screens 320px+ wide
- Touch-friendly buttons (44px+ tap targets)
- Readable text (16px+ on mobile)
- Horizontal scrolling for wide content
- Collapsible navigation menu

### Accessibility
- Semantic HTML elements
- Alt text for images
- Keyboard navigation support
- Focus indicators
- Screen reader friendly
- Color contrast ratio 4.5:1+

## ⚡ Performance Requirements

### Loading Times
- **First page load**: < 3 seconds (3G network)
- **Subsequent pages**: < 1 second
- **API responses**: < 500ms average
- **Real-time updates**: < 200ms latency

### Concurrent Users
- **Minimum**: 10 concurrent users without performance degradation
- **Target**: 50 concurrent users
- **Stretch**: 100+ concurrent users

### Resource Usage
- **Client bundle size**: < 500KB gzipped
- **Memory usage**: < 100MB per user session
- **Server memory**: < 512MB for 50 concurrent users

## 🧪 Testing Requirements

### Manual Testing Checklist
```
User Flow Testing:
[ ] Register new account
[ ] Login with existing account
[ ] Create a learning node
[ ] Edit your own node
[ ] Create a learning path with 3+ nodes
[ ] Follow someone else's learning path
[ ] Fork a learning path
[ ] Generate AI TLDR for a node
[ ] View activity feed updates in real-time
[ ] Test on mobile device
[ ] Test keyboard navigation
[ ] Test with screen reader (if possible)

Cross-browser Testing:
[ ] Chrome/Chromium
[ ] Firefox
[ ] Safari (if Mac available)
[ ] Mobile browsers

Performance Testing:
[ ] Lighthouse audit (score 80+)
[ ] Real-time latency test
[ ] Concurrent user simulation
```

### Load Testing
```bash
# Example Artillery test
# Save as load-test.yml
config:
  target: 'http://localhost:3000'
  phases:
    - duration: 60
      arrivalRate: 5
scenarios:
  - name: "Browse and create content"
    flow:
      - get:
          url: "/"
      - get:
          url: "/nodes"
      - post:
          url: "/api/nodes"
          json:
            title: "Test Node"
            content: { text: "Test content" }
```

## 📋 Evaluation Criteria

Each prototype will be evaluated on:

1. **Feature Completeness** (40%)
   - All MVP features implemented and working
   - Quality of implementation (polish, error handling)
   - Mobile responsiveness

2. **Development Experience** (25%)
   - Time to implement features
   - Debugging and troubleshooting ease
   - Code maintainability

3. **Performance** (20%)
   - Page load times
   - Real-time latency
   - Concurrent user handling

4. **Technical Excellence** (15%)
   - Code quality and architecture
   - Security implementation
   - Error handling and edge cases

## 🚀 Getting Started

1. **Choose your prototype** from the available options
2. **Copy this requirements document** to your prototype directory
3. **Implement features incrementally**, checking them off as you go
4. **Test frequently** with the provided sample data
5. **Document challenges and solutions** for evaluation
6. **Deploy to a public URL** for team review

## 📞 Questions & Clarifications

If any requirements are unclear or you need to make assumptions:

1. **Document your assumptions** in your prototype README
2. **Explain your reasoning** for technical choices
3. **Note any limitations** or trade-offs made
4. **Suggest improvements** for the requirements

Remember: The goal is fair comparison, not perfect implementation. Focus on demonstrating the strengths and weaknesses of your chosen tech stack.

---

*Requirements Version: 1.0*  
*Last Updated: 2025-06-26*  
*Status: Ready for implementation*