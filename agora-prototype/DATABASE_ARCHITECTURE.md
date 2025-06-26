# 🗄️ Database Architecture - Human Intelligence Platform

> **Hybrid Database Strategy**: PostgreSQL + ChromaDB + Redis  
> **Goal**: Support both structured social data and semantic knowledge operations  
> **Timeline**: Phased migration with immediate ChromaDB integration

---

## 🎯 Architecture Overview

### **Current Challenge**
The Human Intelligence platform needs to handle:
- **Structured social data** (users, spaces, relationships)
- **Knowledge content** with semantic search capabilities
- **Real-time presence** and activity streams
- **AI-powered recommendations** and cross-pollination
- **Version control** for collaborative knowledge trees

### **Proposed Solution: Hybrid Database Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │    ChromaDB     │    │     Redis       │
│                 │    │                 │    │                 │
│ • User accounts │    │ • Knowledge     │    │ • Sessions      │
│ • Social graph  │    │   embeddings    │    │ • Real-time     │
│ • Permissions   │    │ • Semantic      │    │   presence      │
│ • Audit logs    │    │   search        │    │ • WebSocket     │
│ • Metadata      │    │ • AI features   │    │   state         │
│                 │    │ • Vector ops    │    │ • Cache layer   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 🧩 ChromaDB Integration Benefits

### **Why ChromaDB is Perfect for Human Intelligence**

#### **1. Semantic Search & Discovery**
```python
# Find knowledge related to user's current exploration
def find_related_knowledge(current_topic: str, user_context: dict):
    results = chroma_client.query(
        collection_name="knowledge_nodes",
        query_texts=[current_topic],
        n_results=10,
        where={"user_id": {"$ne": user_context["user_id"]}},  # Exclude user's own content
        include=["metadatas", "documents", "distances"]
    )
    return results
```

#### **2. Cross-Pollination Discovery**
```python
# Find unexpected connections between knowledge domains
def discover_cross_pollination(user_interests: list, current_space: str):
    # Get embeddings for user's usual topics
    user_embedding = chroma_client.query(
        collection_name="user_profiles",
        query_texts=user_interests,
        n_results=1
    )
    
    # Find semantically distant but potentially interesting content
    distant_topics = chroma_client.query(
        collection_name="knowledge_spaces",
        query_embeddings=user_embedding["embeddings"],
        n_results=5,
        where={"space_id": {"$ne": current_space}},
        # Use larger distance threshold for exploration
        distance_threshold=0.7
    )
    
    return distant_topics
```

#### **3. AI-Powered Recommendations**
```python
# Generate personalized learning paths
def generate_learning_path(user_id: str, target_concept: str):
    # Get user's knowledge state
    user_knowledge = chroma_client.get(
        collection_name="user_knowledge_state",
        where={"user_id": user_id}
    )
    
    # Find optimal learning sequence
    learning_path = chroma_client.query(
        collection_name="knowledge_prerequisites",
        query_texts=[target_concept],
        n_results=20,
        where={
            "difficulty": {"$lte": user_knowledge["level"]},
            "prerequisite_met": True
        }
    )
    
    return learning_path
```

---

## 🏗️ Data Model Design

### **PostgreSQL Schema (Structured Data)**

```sql
-- Core user and social data
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- User preferences
    interests JSONB DEFAULT '[]',
    learning_style JSONB DEFAULT '{}',
    privacy_settings JSONB DEFAULT '{}'
);

-- Discussion spaces (agora rooms)
CREATE TABLE spaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES users(id),
    capacity INTEGER DEFAULT 12,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Space configuration
    constraints JSONB DEFAULT '{"max_participants": 12, "visible_presence": true}',
    topics JSONB DEFAULT '[]',
    
    -- Analytics
    total_visits INTEGER DEFAULT 0,
    avg_session_duration INTERVAL,
    quality_score DECIMAL(3,2)
);

-- Social relationships
CREATE TABLE friendships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    friend_id UUID REFERENCES users(id),
    status VARCHAR(20) CHECK (status IN ('pending', 'accepted', 'blocked')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, friend_id)
);

-- Knowledge trees (metadata only, content in ChromaDB)
CREATE TABLE knowledge_trees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    owner_id UUID REFERENCES users(id),
    parent_tree_id UUID REFERENCES knowledge_trees(id), -- For forks
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Git-like versioning
    version VARCHAR(50) DEFAULT 'main',
    branch_name VARCHAR(100) DEFAULT 'main',
    
    -- Metadata
    tags JSONB DEFAULT '[]',
    language VARCHAR(10) DEFAULT 'en',
    license VARCHAR(50) DEFAULT 'CC-BY-SA',
    
    -- ChromaDB collection reference
    chroma_collection_id VARCHAR(255),
    
    -- Analytics
    view_count INTEGER DEFAULT 0,
    fork_count INTEGER DEFAULT 0,
    star_count INTEGER DEFAULT 0
);

-- Knowledge nodes (structure only, content in ChromaDB)
CREATE TABLE knowledge_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tree_id UUID REFERENCES knowledge_trees(id),
    parent_node_id UUID REFERENCES knowledge_nodes(id),
    title VARCHAR(500) NOT NULL,
    node_type VARCHAR(50) DEFAULT 'content', -- content, discussion, quiz, etc.
    order_index INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- ChromaDB document reference
    chroma_document_id VARCHAR(255),
    
    -- Node metadata
    estimated_read_time INTEGER, -- minutes
    difficulty_level INTEGER CHECK (difficulty_level BETWEEN 1 AND 10),
    prerequisites JSONB DEFAULT '[]',
    
    -- Versioning
    version INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true
);

-- User presence in spaces (real-time, also cached in Redis)
CREATE TABLE user_presence (
    user_id UUID REFERENCES users(id),
    space_id UUID REFERENCES spaces(id),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active', -- active, away, focused
    
    PRIMARY KEY (user_id, space_id)
);

-- Migration events (track cross-pollination)
CREATE TABLE migration_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    from_space_id UUID REFERENCES spaces(id),
    to_space_id UUID REFERENCES spaces(id),
    migration_reason VARCHAR(100), -- capacity_full, friend_suggestion, ai_recommendation
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Outcome tracking
    stayed_duration INTERVAL,
    engagement_level INTEGER CHECK (engagement_level BETWEEN 1 AND 5),
    led_to_new_interest BOOLEAN DEFAULT false
);

-- Activity feed
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    activity_type VARCHAR(50) NOT NULL, -- join_space, create_node, cross_pollinate, etc.
    target_type VARCHAR(50), -- space, tree, node, user
    target_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Activity metadata
    metadata JSONB DEFAULT '{}',
    visibility VARCHAR(20) DEFAULT 'friends' -- public, friends, private
);

-- Indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_spaces_public ON spaces(is_public) WHERE is_public = true;
CREATE INDEX idx_knowledge_trees_public ON knowledge_trees(is_public) WHERE is_public = true;
CREATE INDEX idx_knowledge_nodes_tree ON knowledge_nodes(tree_id);
CREATE INDEX idx_user_presence_active ON user_presence(space_id) WHERE status = 'active';
CREATE INDEX idx_activities_user_time ON activities(user_id, created_at DESC);
CREATE INDEX idx_migration_events_user ON migration_events(user_id, created_at DESC);
```

### **ChromaDB Collections Structure**

```python
# ChromaDB collections for semantic operations
CHROMA_COLLECTIONS = {
    # Knowledge content with embeddings
    "knowledge_nodes": {
        "embeddings": "text embeddings from content",
        "metadata": {
            "node_id": "UUID from PostgreSQL",
            "tree_id": "UUID from PostgreSQL", 
            "user_id": "content creator",
            "node_type": "content|discussion|quiz|code",
            "difficulty": "1-10 scale",
            "topics": ["AI", "Philosophy", "etc"],
            "language": "en|es|fr|etc",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "version": "version number",
            "word_count": "content length",
            "estimated_read_time": "minutes"
        },
        "documents": "actual content text (markdown, code, etc)"
    },
    
    # User interest profiles for recommendations
    "user_profiles": {
        "embeddings": "aggregated embeddings from user interactions",
        "metadata": {
            "user_id": "UUID from PostgreSQL",
            "interests": ["primary interest topics"],
            "expertise_level": {"topic": "level 1-10"},
            "learning_style": "visual|textual|interactive|mixed",
            "activity_patterns": "time preferences, session lengths",
            "cross_pollination_openness": "0.0-1.0 score",
            "last_updated": "timestamp"
        },
        "documents": "text representation of user interests and background"
    },
    
    # Discussion space embeddings for matching
    "discussion_spaces": {
        "embeddings": "embeddings from space descriptions and topics",
        "metadata": {
            "space_id": "UUID from PostgreSQL",
            "space_name": "human readable name",
            "primary_topics": ["main discussion topics"],
            "average_depth": "conversation quality metric",
            "typical_participants": "user types who join",
            "activity_level": "high|medium|low",
            "requires_expertise": "boolean",
            "language": "primary language",
            "created_at": "timestamp"
        },
        "documents": "space description and topic summaries"
    },
    
    # Cross-pollination opportunities
    "topic_relationships": {
        "embeddings": "embeddings representing topic connections",
        "metadata": {
            "topic_a": "first topic",
            "topic_b": "second topic", 
            "relationship_type": "prerequisite|related|opposing|complementary",
            "strength": "0.0-1.0 connection strength",
            "discovered_by": "user_generated|ai_detected|expert_verified",
            "evidence_count": "number of supporting interactions",
            "last_validated": "timestamp"
        },
        "documents": "description of how topics relate"
    },
    
    # Learning path optimization
    "learning_sequences": {
        "embeddings": "embeddings of successful learning progressions",
        "metadata": {
            "sequence_id": "unique identifier",
            "start_topic": "beginning concept",
            "end_topic": "target concept",
            "node_sequence": ["ordered list of node IDs"],
            "average_completion_time": "typical duration",
            "success_rate": "percentage who complete successfully",
            "difficulty_progression": "how difficulty changes",
            "prerequisite_knowledge": ["required background"],
            "created_by": "user_discovered|ai_generated|expert_designed"
        },
        "documents": "description of learning progression and rationale"
    }
}
```

### **Redis Data Structures**

```python
# Redis schema for real-time features
REDIS_KEYS = {
    # Real-time presence
    "presence:space:{space_id}": {
        "type": "hash",
        "fields": {
            "user:{user_id}": {
                "joined_at": "timestamp",
                "last_activity": "timestamp", 
                "status": "active|away|focused",
                "current_node": "node_id if viewing specific content"
            }
        },
        "ttl": "1 hour (auto-expire inactive users)"
    },
    
    # Active spaces (for quick lookup)
    "spaces:active": {
        "type": "sorted_set",
        "score": "participant_count",
        "members": "space_id values",
        "ttl": "5 minutes (refreshed regularly)"
    },
    
    # User session state
    "session:{user_id}": {
        "type": "hash", 
        "fields": {
            "current_space": "space_id",
            "current_tree": "tree_id",
            "current_node": "node_id",
            "session_start": "timestamp",
            "last_activity": "timestamp",
            "exploration_path": "ordered list of visited locations"
        },
        "ttl": "24 hours"
    },
    
    # Real-time activity feed
    "activity:global": {
        "type": "stream",
        "entries": {
            "user_id": "UUID",
            "action": "join_space|leave_space|create_node|etc",
            "target": "affected resource ID",
            "timestamp": "when it happened",
            "metadata": "additional context"
        },
        "maxlen": "10000 entries"
    },
    
    # Friend network activity
    "activity:friends:{user_id}": {
        "type": "stream",
        "entries": "filtered activity from user's friend network",
        "maxlen": "1000 entries",
        "ttl": "7 days"
    },
    
    # Cross-pollination suggestions cache
    "suggestions:{user_id}": {
        "type": "list",
        "values": "JSON objects with suggestion data",
        "ttl": "15 minutes (regenerated frequently)"
    },
    
    # WebSocket connection tracking
    "connections:space:{space_id}": {
        "type": "set",
        "members": "websocket_connection_id values",
        "ttl": "1 hour"
    }
}
```

---

## 🔄 Integration Patterns

### **1. Semantic Search Integration**

```python
async def semantic_search_knowledge(query: str, user_context: dict, limit: int = 10):
    """
    Search knowledge using both PostgreSQL metadata and ChromaDB embeddings
    """
    # Get semantic matches from ChromaDB
    semantic_results = await chroma_client.query(
        collection_name="knowledge_nodes",
        query_texts=[query],
        n_results=limit * 2,  # Get more than needed for filtering
        include=["metadatas", "documents", "distances"]
    )
    
    # Extract node IDs for PostgreSQL query
    node_ids = [result["node_id"] for result in semantic_results["metadatas"]]
    
    # Get detailed metadata from PostgreSQL
    detailed_results = await db.fetch_all("""
        SELECT n.*, t.title as tree_title, u.display_name as creator_name
        FROM knowledge_nodes n
        JOIN knowledge_trees t ON n.tree_id = t.id  
        JOIN users u ON t.owner_id = u.id
        WHERE n.id = ANY($1) AND n.is_active = true
        ORDER BY array_position($1, n.id)
    """, node_ids)
    
    # Combine ChromaDB scores with PostgreSQL metadata
    combined_results = []
    for i, node in enumerate(detailed_results):
        semantic_score = semantic_results["distances"][i]
        combined_results.append({
            **node,
            "content": semantic_results["documents"][i],
            "semantic_similarity": 1 - semantic_score,  # Convert distance to similarity
            "relevance_explanation": generate_relevance_explanation(
                query, semantic_results["documents"][i]
            )
        })
    
    return combined_results[:limit]
```

### **2. Cross-Pollination Discovery**

```python
async def discover_cross_pollination_opportunities(user_id: str):
    """
    Find unexpected learning opportunities using ChromaDB
    """
    # Get user's current interests and exploration history
    user_profile = await chroma_client.get(
        collection_name="user_profiles",
        where={"user_id": user_id}
    )
    
    # Find topics that are semantically distant but potentially interesting
    exploration_candidates = await chroma_client.query(
        collection_name="topic_relationships", 
        query_embeddings=user_profile["embeddings"],
        n_results=20,
        where={
            "relationship_type": {"$in": ["complementary", "related"]},
            "strength": {"$gte": 0.6}
        }
    )
    
    # Filter by current space availability and friend presence
    viable_opportunities = []
    for candidate in exploration_candidates["metadatas"]:
        # Check if there are active spaces for this topic
        active_spaces = await redis_client.zrangebyscore(
            "spaces:active", 1, 11  # Spaces with 1-11 participants (not full)
        )
        
        topic_spaces = await db.fetch_all("""
            SELECT s.*, COUNT(p.user_id) as current_participants
            FROM spaces s
            LEFT JOIN user_presence p ON s.id = p.space_id
            WHERE s.topics @> $1 AND s.id = ANY($2)
            GROUP BY s.id
            HAVING COUNT(p.user_id) < s.capacity
        """, [candidate["topic_b"]], active_spaces)
        
        if topic_spaces:
            viable_opportunities.append({
                "topic": candidate["topic_b"],
                "relationship": candidate["relationship_type"],
                "strength": candidate["strength"],
                "available_spaces": topic_spaces,
                "explanation": f"Explore {candidate['topic_b']} - it {candidate['relationship_type']} your interest in {candidate['topic_a']}"
            })
    
    return viable_opportunities
```

### **3. AI-Powered Learning Path Generation**

```python
async def generate_personalized_learning_path(user_id: str, target_concept: str):
    """
    Create optimal learning sequence using ChromaDB and user data
    """
    # Get user's current knowledge state
    user_knowledge = await chroma_client.get(
        collection_name="user_profiles",
        where={"user_id": user_id}
    )
    
    # Find existing successful learning sequences to target concept
    successful_paths = await chroma_client.query(
        collection_name="learning_sequences",
        query_texts=[target_concept],
        n_results=10,
        where={
            "end_topic": target_concept,
            "success_rate": {"$gte": 0.7}
        }
    )
    
    # Analyze user's prerequisite gaps
    for path in successful_paths["metadatas"]:
        prerequisites = path["prerequisite_knowledge"]
        user_expertise = user_knowledge["metadata"][0]["expertise_level"]
        
        # Calculate knowledge gaps
        gaps = []
        for prereq in prerequisites:
            if prereq not in user_expertise or user_expertise[prereq] < 6:
                gaps.append(prereq)
        
        # Generate filling sequence for gaps
        if gaps:
            gap_filling_nodes = await chroma_client.query(
                collection_name="knowledge_nodes",
                query_texts=gaps,
                n_results=20,
                where={
                    "difficulty": {"$lte": user_expertise.get("general", 5)},
                    "node_type": "content"
                }
            )
            
            # Optimize sequence order
            optimal_sequence = optimize_learning_sequence(
                gap_filling_nodes, path["node_sequence"], user_knowledge
            )
            
            return {
                "path_id": path["sequence_id"],
                "target": target_concept,
                "estimated_time": path["average_completion_time"],
                "success_probability": calculate_success_probability(user_knowledge, path),
                "prerequisite_nodes": optimal_sequence["prerequisites"],
                "main_sequence": path["node_sequence"],
                "total_nodes": len(optimal_sequence["prerequisites"]) + len(path["node_sequence"]),
                "personalization_notes": generate_personalization_notes(user_knowledge, gaps)
            }
    
    # If no good existing paths, generate new one with AI
    return await generate_novel_learning_path(user_id, target_concept)
```

---

## 🚀 Migration Strategy

### **Phase 1: ChromaDB Integration (Weeks 1-2)**

```python
# 1. Set up ChromaDB alongside existing PostgreSQL
async def setup_chromadb():
    # Initialize ChromaDB client
    chroma_client = chromadb.PersistentClient(path="./chroma_data")
    
    # Create collections
    for collection_name, config in CHROMA_COLLECTIONS.items():
        collection = chroma_client.create_collection(
            name=collection_name,
            embedding_function=OpenAIEmbeddingFunction(
                api_key=settings.OPENAI_API_KEY,
                model_name="text-embedding-3-small"
            )
        )
        
    return chroma_client

# 2. Migrate existing content to ChromaDB
async def migrate_existing_content():
    # Get all knowledge nodes from PostgreSQL
    nodes = await db.fetch_all("""
        SELECT n.*, t.title as tree_title 
        FROM knowledge_nodes n
        JOIN knowledge_trees t ON n.tree_id = t.id
        WHERE n.is_active = true
    """)
    
    # Process in batches
    batch_size = 100
    for i in range(0, len(nodes), batch_size):
        batch = nodes[i:i + batch_size]
        
        # Prepare data for ChromaDB
        documents = []
        metadatas = []
        ids = []
        
        for node in batch:
            documents.append(node["content"] or node["title"])
            metadatas.append({
                "node_id": str(node["id"]),
                "tree_id": str(node["tree_id"]),
                "node_type": node["node_type"],
                "difficulty": node["difficulty_level"],
                "created_at": node["created_at"].isoformat()
            })
            ids.append(str(node["id"]))
        
        # Add to ChromaDB
        chroma_client.get_collection("knowledge_nodes").add(
            documents=documents,
            metadatas=metadatas,
            ids=ids
        )
        
        # Update PostgreSQL with ChromaDB references
        await db.execute_many("""
            UPDATE knowledge_nodes 
            SET chroma_document_id = $1 
            WHERE id = $2
        """, [(node_id, node_id) for node_id in ids])
```

### **Phase 2: Semantic Features (Weeks 3-4)**

```python
# Add semantic search to existing endpoints
@app.get("/api/search/semantic")
async def semantic_search_endpoint(
    query: str,
    user_id: str = Depends(get_current_user),
    limit: int = 10
):
    results = await semantic_search_knowledge(query, {"user_id": user_id}, limit)
    return {"results": results, "query": query}

# Add cross-pollination suggestions
@app.get("/api/users/{user_id}/cross-pollination")
async def get_cross_pollination_suggestions(user_id: str):
    opportunities = await discover_cross_pollination_opportunities(user_id)
    return {"opportunities": opportunities}
```

### **Phase 3: AI-Powered Features (Weeks 5-6)**

```python
# Add learning path generation
@app.post("/api/learning-paths/generate")
async def generate_learning_path_endpoint(
    target_concept: str,
    user_id: str = Depends(get_current_user)
):
    path = await generate_personalized_learning_path(user_id, target_concept)
    return {"learning_path": path}

# Add intelligent recommendations
@app.get("/api/users/{user_id}/recommendations")
async def get_ai_recommendations(user_id: str):
    recommendations = await generate_ai_recommendations(user_id)
    return {"recommendations": recommendations}
```

---

## 📈 Performance Optimization

### **Caching Strategy**

```python
# Multi-layer caching for optimal performance
CACHE_LAYERS = {
    # Layer 1: Redis (fastest, real-time data)
    "real_time": {
        "user_presence": "5 minutes TTL",
        "space_activity": "1 minute TTL", 
        "friend_activity": "10 minutes TTL"
    },
    
    # Layer 2: Application cache (in-memory)
    "application": {
        "user_profiles": "1 hour TTL",
        "space_metadata": "30 minutes TTL",
        "popular_content": "6 hours TTL"
    },
    
    # Layer 3: ChromaDB query cache
    "semantic": {
        "similar_content": "24 hours TTL",
        "cross_pollination": "12 hours TTL",
        "learning_paths": "1 week TTL"
    }
}
```

### **Database Optimization**

```sql
-- PostgreSQL optimization
-- Partitioning for large tables
CREATE TABLE activities_2024 PARTITION OF activities
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- Materialized views for analytics
CREATE MATERIALIZED VIEW user_activity_summary AS
SELECT 
    user_id,
    COUNT(*) as total_activities,
    COUNT(DISTINCT target_id) as unique_interactions,
    AVG(CASE WHEN activity_type = 'cross_pollinate' THEN 1 ELSE 0 END) as cross_pollination_rate
FROM activities 
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY user_id;

-- Refresh periodically
CREATE OR REPLACE FUNCTION refresh_analytics()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_activity_summary;
END;
$$ LANGUAGE plpgsql;

-- Schedule refresh every hour
SELECT cron.schedule('refresh-analytics', '0 * * * *', 'SELECT refresh_analytics();');
```

---

## 🔧 Implementation Code Examples

### **ChromaDB Service Layer**

```python
class ChromaService:
    def __init__(self):
        self.client = chromadb.PersistentClient(path=settings.CHROMA_DB_PATH)
        self.embedding_function = OpenAIEmbeddingFunction(
            api_key=settings.OPENAI_API_KEY,
            model_name="text-embedding-3-small"
        )
    
    async def add_knowledge_node(self, node_data: dict):
        """Add new knowledge node to ChromaDB"""
        collection = self.client.get_collection("knowledge_nodes")
        
        collection.add(
            documents=[node_data["content"]],
            metadatas=[{
                "node_id": node_data["id"],
                "tree_id": node_data["tree_id"],
                "difficulty": node_data["difficulty"],
                "topics": node_data["topics"],
                "created_at": datetime.now().isoformat()
            }],
            ids=[str(node_data["id"])]
        )
    
    async def find_similar_content(self, content: str, exclude_ids: list = None):
        """Find semantically similar content"""
        collection = self.client.get_collection("knowledge_nodes")
        
        where_clause = {}
        if exclude_ids:
            where_clause["node_id"] = {"$nin": exclude_ids}
        
        results = collection.query(
            query_texts=[content],
            n_results=10,
            where=where_clause,
            include=["metadatas", "documents", "distances"]
        )
        
        return results
    
    async def update_user_profile(self, user_id: str, interactions: list):
        """Update user profile based on interactions"""
        collection = self.client.get_or_create_collection("user_profiles")
        
        # Generate profile embedding from interaction history
        interaction_texts = [interaction["content"] for interaction in interactions]
        profile_text = " ".join(interaction_texts[-50:])  # Last 50 interactions
        
        # Upsert user profile
        collection.upsert(
            documents=[profile_text],
            metadatas=[{
                "user_id": user_id,
                "interests": extract_interests(interactions),
                "expertise_level": calculate_expertise_levels(interactions),
                "last_updated": datetime.now().isoformat()
            }],
            ids=[user_id]
        )
```

### **Hybrid Query Service**

```python
class HybridQueryService:
    def __init__(self, db_pool, chroma_service, redis_client):
        self.db = db_pool
        self.chroma = chroma_service
        self.redis = redis_client
    
    async def search_knowledge(self, query: str, user_id: str, filters: dict = None):
        """Hybrid search combining PostgreSQL filters and ChromaDB similarity"""
        
        # 1. Get semantic matches from ChromaDB
        semantic_results = await self.chroma.find_similar_content(query)
        node_ids = [meta["node_id"] for meta in semantic_results["metadatas"]]
        
        # 2. Apply PostgreSQL filters and get metadata
        sql_filters = []
        params = [node_ids]
        
        if filters:
            if filters.get("difficulty_max"):
                    sql_filters.append("n.difficulty_level <= $%d" % (len(params) + 1))
                    params.append(filters["difficulty_max"])
            
                if filters.get("language"):
                    sql_filters.append("t.language = $%d" % (len(params) + 1))
                    params.append(filters["language"])
        
            where_clause = "WHERE n.id = ANY($1) AND n.is_active = true"
            if sql_filters:
                where_clause += " AND " + " AND ".join(sql_filters)
        
            detailed_results = await self.db.fetch_all(f"""
                SELECT n.*, t.title as tree_title, t.language, u.display_name as creator_name,
                       array_position($1, n.id) as relevance_order
                FROM knowledge_nodes n
                JOIN knowledge_trees t ON n.tree_id = t.id  
                JOIN users u ON t.owner_id = u.id
                {where_clause}
                ORDER BY array_position($1, n.id)
            """, *params)
        
            # 3. Combine results with semantic scores
            results = []
            for i, node in enumerate(detailed_results):
                semantic_score = 1 - semantic_results["distances"][i]  # Convert to similarity
            
                results.append({
                    **dict(node),
                    "content_preview": semantic_results["documents"][i][:200] + "...",
                    "semantic_similarity": semantic_score,
                    "hybrid_score": self._calculate_hybrid_score(semantic_score, node, user_id)
                })
        
            return sorted(results, key=lambda x: x["hybrid_score"], reverse=True)
    
        def _calculate_hybrid_score(self, semantic_score: float, node: dict, user_id: str) -> float:
            """Calculate weighted score combining semantic similarity with user preferences"""
            base_score = semantic_score * 0.7  # 70% weight on semantic similarity
        
            # Add user preference bonus (30% weight)
            preference_bonus = 0.0
        
            # Bonus for appropriate difficulty level (get user level from cache/db)
            # This would be implemented with actual user data
        
            # Bonus for content from followed creators
            # This would check if user follows the creator
        
            # Bonus for recently updated content
            if node.get("updated_at"):
                days_since_update = (datetime.now() - node["updated_at"]).days
                if days_since_update < 30:
                    preference_bonus += 0.1
        
            return min(1.0, base_score + preference_bonus)
    ```

    ---

    ## 🔒 Security & Privacy

    ### **Data Protection Strategy**

    ```python
    # Implement data anonymization for ChromaDB
    class PrivacyProtectedChromaService:
        def __init__(self):
            self.encryption_key = Fernet.generate_key()
            self.fernet = Fernet(self.encryption_key)
    
        async def add_private_content(self, content: str, user_id: str, metadata: dict):
            """Add content with privacy protection"""
            # Anonymize user identifiers
            anonymous_id = self._generate_anonymous_id(user_id)
        
            # Encrypt sensitive metadata
            protected_metadata = {
                **metadata,
                "user_id": anonymous_id,
                "original_user_hash": hashlib.sha256(user_id.encode()).hexdigest()[:16]
            }
        
            # Remove personally identifiable information from content
            cleaned_content = self._remove_pii(content)
        
            collection = self.client.get_collection("knowledge_nodes")
            collection.add(
                documents=[cleaned_content],
                metadatas=[protected_metadata],
                ids=[str(uuid.uuid4())]
            )
    
        def _remove_pii(self, content: str) -> str:
            """Remove personally identifiable information from content"""
            # Implement PII detection and removal
            # This would use libraries like spaCy or regex patterns
            return content  # Simplified for example
    ```

    ### **GDPR Compliance**

    ```sql
    -- Data retention policies
    CREATE TABLE data_retention_policies (
        data_type VARCHAR(100) PRIMARY KEY,
        retention_period INTERVAL NOT NULL,
        deletion_strategy VARCHAR(50) NOT NULL
    );

    INSERT INTO data_retention_policies VALUES
    ('user_activities', '2 years', 'hard_delete'),
    ('user_presence', '90 days', 'hard_delete'),
    ('knowledge_content', '10 years', 'anonymize'),
    ('user_profiles', 'until_account_deletion', 'hard_delete');

    -- User data export functionality
    CREATE OR REPLACE FUNCTION export_user_data(target_user_id UUID)
    RETURNS JSON AS $$
    DECLARE
        user_data JSON;
    BEGIN
        SELECT json_build_object(
            'user_profile', (SELECT row_to_json(u) FROM users u WHERE id = target_user_id),
            'knowledge_trees', (SELECT json_agg(t) FROM knowledge_trees t WHERE owner_id = target_user_id),
            'activities', (SELECT json_agg(a) FROM activities a WHERE user_id = target_user_id),
            'friendships', (SELECT json_agg(f) FROM friendships f WHERE user_id = target_user_id OR friend_id = target_user_id)
        ) INTO user_data;
    
        RETURN user_data;
    END;
    $$ LANGUAGE plpgsql;
    ```

    ---

    ## 📊 Monitoring & Analytics

    ### **Performance Monitoring**

    ```python
    # Database performance monitoring
    class DatabaseMonitor:
        def __init__(self, db_pool, chroma_client, redis_client):
            self.db = db_pool
            self.chroma = chroma_client
            self.redis = redis_client
        
        async def collect_metrics(self):
            """Collect performance metrics from all database systems"""
            metrics = {
                "postgresql": await self._collect_postgres_metrics(),
                "chromadb": await self._collect_chroma_metrics(),
                "redis": await self._collect_redis_metrics(),
                "hybrid_queries": await self._collect_hybrid_metrics()
            }
        
            # Send to monitoring system (Prometheus, etc.)
            await self._send_to_monitoring(metrics)
        
            return metrics
    
        async def _collect_postgres_metrics(self):
            """Collect PostgreSQL performance metrics"""
            stats = await self.db.fetch_one("""
                SELECT 
                    (SELECT count(*) FROM users) as total_users,
                    (SELECT count(*) FROM knowledge_trees WHERE is_public = true) as public_trees,
                    (SELECT count(*) FROM knowledge_nodes WHERE is_active = true) as active_nodes,
                    (SELECT avg(extract(epoch from (now() - last_activity))) FROM user_presence) as avg_session_duration
            """)
        
            return dict(stats)
    
        async def _collect_chroma_metrics(self):
            """Collect ChromaDB performance metrics"""
            collections = self.chroma.list_collections()
            metrics = {}
        
            for collection in collections:
                coll = self.chroma.get_collection(collection.name)
                metrics[collection.name] = {
                    "document_count": coll.count(),
                    "last_updated": datetime.now().isoformat()
                }
        
            return metrics
    
        async def _collect_redis_metrics(self):
            """Collect Redis performance metrics"""
            info = self.redis.info()
            return {
                "connected_clients": info.get("connected_clients", 0),
                "used_memory": info.get("used_memory", 0),
                "keyspace_hits": info.get("keyspace_hits", 0),
                "keyspace_misses": info.get("keyspace_misses", 0)
            }
    ```

    ---

    ## 🔄 Backup & Recovery

    ### **Backup Strategy**

    ```bash
    #!/bin/bash
    # Comprehensive backup script for all database systems

    # PostgreSQL backup
    pg_dump $DATABASE_URL | gzip > "backups/postgres_$(date +%Y%m%d_%H%M%S).sql.gz"

    # ChromaDB backup
    tar -czf "backups/chromadb_$(date +%Y%m%d_%H%M%S).tar.gz" ./chroma_data/

    # Redis backup
    redis-cli --rdb "backups/redis_$(date +%Y%m%d_%H%M%S).rdb"

    # Upload to cloud storage
    aws s3 sync backups/ s3://human-intelligence-backups/$(date +%Y/%m/%d)/

    # Clean up old local backups (keep last 7 days)
    find backups/ -name "*.gz" -mtime +7 -delete
    find backups/ -name "*.rdb" -mtime +7 -delete
    ```

    ### **Disaster Recovery Plan**

    ```yaml
    Recovery_Procedures:
      Total_System_Loss:
        1. "Restore PostgreSQL from latest backup"
        2. "Restore ChromaDB data directory"
        3. "Restore Redis data"
        4. "Rebuild ChromaDB indexes"
        5. "Verify data consistency"
        6. "Resume application services"
    
      Partial_Data_Loss:
        ChromaDB_Loss:
          - "Restore from backup"
          - "Re-generate embeddings from PostgreSQL content"
          - "Rebuild semantic indexes"
    
        PostgreSQL_Loss:
          - "Restore from backup"
          - "Rebuild foreign key relationships"
          - "Re-sync with ChromaDB metadata"
    
        Redis_Loss:
          - "Accept data loss (cache layer)"
          - "Rebuild from PostgreSQL data"
          - "Resume real-time features"

      Recovery_Time_Objectives:
        PostgreSQL: "4 hours maximum"
        ChromaDB: "8 hours maximum (includes re-indexing)"
        Redis: "30 minutes maximum"
        Full_System: "12 hours maximum"
    ```

    ---

    ## 🚀 Future Enhancements

    ### **Planned Database Features**

    1. **Multi-Modal Embeddings**
       - Support for image, audio, and video content in ChromaDB
       - Cross-modal semantic search capabilities

    2. **Distributed ChromaDB**
       - Horizontal scaling for large knowledge bases
       - Regional deployment for global latency optimization

    3. **Advanced Analytics**
       - Knowledge flow analysis through graph databases
       - Predictive modeling for learning outcomes

    4. **Blockchain Integration**
       - Decentralized knowledge ownership
       - Immutable contribution tracking

    ### **Scaling Roadmap**

    ```
    Current (MVP): 1K users, 10K nodes
    └── PostgreSQL + ChromaDB + Redis (single instance)

    Phase 1 (10K users): 100K nodes  
    └── Read replicas + ChromaDB clustering

    Phase 2 (100K users): 1M nodes
    └── Database sharding + distributed ChromaDB

    Phase 3 (1M+ users): 10M+ nodes
    └── Multi-region deployment + edge caching
    ```

    ---

    ## 💡 Key Takeaways

    ### **Why This Architecture Works for Human Intelligence**

    1. **Hybrid Approach**: Combines the strengths of relational and vector databases
    2. **Semantic Capabilities**: ChromaDB enables AI-powered features from day one
    3. **Scalability**: Each component can scale independently based on load patterns
    4. **Real-time Performance**: Redis ensures responsive user interactions
    5. **Future-Proof**: Architecture supports planned AI and semantic features

    ### **Implementation Benefits**

    - **Immediate Semantic Search**: ChromaDB integration enables powerful content discovery
    - **Cross-Pollination AI**: Vector similarities drive intelligent content recommendations  
    - **Performance Optimization**: Multi-layer caching ensures responsive user experience
    - **Data Consistency**: PostgreSQL maintains relational integrity for core platform data
    - **Developer Experience**: Clear separation of concerns between structured and semantic data

    **This hybrid database architecture provides the perfect foundation for transforming Human Intelligence from a social platform into an AI-augmented collaborative knowledge system.**

    ---

    *Database Architecture completed June 26, 2025*  
    *PostgreSQL + ChromaDB + Redis hybrid approach*  
    *Ready for immediate ChromaDB integration*