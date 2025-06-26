// Test data service for AI Knowledge Trees frontend development
// Provides mock data that matches the expected API structure

export const mockTrees = [
    {
        id: "tree-1",
        title: "Philosophy Fundamentals",
        description: "Core philosophical concepts and ideas from ancient to modern times",
        owner_id: "default-user",
        is_public: true,
        created_at: "2025-06-26T10:00:00Z",
        updated_at: "2025-06-26T15:30:00Z",
        node_count: 8
    },
    {
        id: "tree-2",
        title: "Computer Science Essentials",
        description: "Fundamental concepts in computer science and programming",
        owner_id: "default-user",
        is_public: true,
        created_at: "2025-06-25T14:20:00Z",
        updated_at: "2025-06-26T09:15:00Z",
        node_count: 12
    },
    {
        id: "tree-3",
        title: "Personal Development",
        description: "Self-improvement and growth strategies",
        owner_id: "default-user",
        is_public: false,
        created_at: "2025-06-24T11:45:00Z",
        updated_at: "2025-06-24T11:45:00Z",
        node_count: 0
    },
    {
        id: "tree-4",
        title: "Mathematics Foundation",
        description: "Essential mathematical concepts and theories",
        owner_id: "default-user",
        is_public: true,
        created_at: "2025-06-23T16:10:00Z",
        updated_at: "2025-06-26T12:00:00Z",
        node_count: 15
    }
];

export const mockNodes = {
    "tree-1": [
        {
            id: "node-1-1",
            tree_id: "tree-1",
            parent_id: null,
            title: "Consciousness",
            content: "The subjective experience of awareness and mental states. A fundamental question in philosophy of mind.",
            depth: 0,
            position: 0,
            created_at: "2025-06-26T10:15:00Z",
            updated_at: "2025-06-26T14:30:00Z",
            ai_analysis: {
                difficulty: 9,
                concepts: ["phenomenology", "qualia", "awareness", "subjective experience", "mind-body problem"],
                suggestions: [
                    "Explore the hard problem of consciousness",
                    "Consider different theories of consciousness",
                    "Examine the relationship between consciousness and the brain",
                    "Study integrated information theory"
                ]
            }
        },
        {
            id: "node-1-2",
            tree_id: "tree-1",
            parent_id: "node-1-1",
            title: "Self-Awareness",
            content: "The ability to recognize oneself as an individual distinct from the environment and other individuals.",
            depth: 1,
            position: 0,
            created_at: "2025-06-26T10:30:00Z",
            updated_at: "2025-06-26T10:30:00Z",
            ai_analysis: {
                difficulty: 5,
                concepts: ["self-recognition", "identity", "introspection", "metacognition"],
                suggestions: [
                    "Study mirror self-recognition tests",
                    "Explore developmental aspects of self-awareness",
                    "Investigate theory of mind development"
                ]
            }
        },
        {
            id: "node-1-3",
            tree_id: "tree-1",
            parent_id: "node-1-1",
            title: "Stream of Consciousness",
            content: "William James's concept describing the continuous flow of thoughts and experiences.",
            depth: 1,
            position: 1,
            created_at: "2025-06-26T11:00:00Z",
            updated_at: "2025-06-26T11:00:00Z",
            ai_analysis: {
                difficulty: 3,
                concepts: ["William James", "consciousness flow", "psychological continuity"],
                suggestions: [
                    "Read James's 'Principles of Psychology'",
                    "Compare with modern theories of consciousness"
                ]
            }
        },
        {
            id: "node-1-4",
            tree_id: "tree-1",
            parent_id: null,
            title: "Ethics",
            content: "The branch of philosophy concerned with moral principles and values that govern behavior.",
            depth: 0,
            position: 1,
            created_at: "2025-06-26T11:30:00Z",
            updated_at: "2025-06-26T15:00:00Z",
            ai_analysis: {
                difficulty: 6,
                concepts: ["morality", "virtue", "duty", "consequentialism", "deontology"],
                suggestions: [
                    "Compare deontological and utilitarian approaches",
                    "Examine virtue ethics and character development",
                    "Study applied ethics in modern contexts",
                    "Explore meta-ethics and moral realism"
                ]
            }
        },
        {
            id: "node-1-5",
            tree_id: "tree-1",
            parent_id: "node-1-4",
            title: "Utilitarianism",
            content: "Ethical theory that judges actions by their consequences, aiming for the greatest good for the greatest number.",
            depth: 1,
            position: 0,
            created_at: "2025-06-26T12:00:00Z",
            updated_at: "2025-06-26T12:00:00Z",
            ai_analysis: {
                difficulty: 4,
                concepts: ["consequentialism", "greatest happiness", "Jeremy Bentham", "John Stuart Mill", "hedonistic calculus"],
                suggestions: [
                    "Study the trolley problem thought experiments",
                    "Compare act vs rule utilitarianism",
                    "Examine preference utilitarianism"
                ]
            }
        },
        {
            id: "node-1-6",
            tree_id: "tree-1",
            parent_id: "node-1-4",
            title: "Deontological Ethics",
            content: "Ethical theory focused on duties and rules, most famously developed by Immanuel Kant.",
            depth: 1,
            position: 1,
            created_at: "2025-06-26T12:30:00Z",
            updated_at: "2025-06-26T12:30:00Z",
            ai_analysis: {
                difficulty: 7,
                concepts: ["categorical imperative", "duty ethics", "Immanuel Kant", "moral law"],
                suggestions: [
                    "Study the categorical imperative formulations",
                    "Examine the concept of moral autonomy",
                    "Compare with hypothetical imperatives"
                ]
            }
        },
        {
            id: "node-1-7",
            tree_id: "tree-1",
            parent_id: null,
            title: "Epistemology",
            content: "The study of knowledge itself - what it is, how it's acquired, and what makes beliefs justified.",
            depth: 0,
            position: 2,
            created_at: "2025-06-26T13:00:00Z",
            updated_at: "2025-06-26T13:00:00Z",
            ai_analysis: {
                difficulty: 6,
                concepts: ["knowledge", "belief", "justification", "skepticism"],
                suggestions: [
                    "Explore the Gettier problem",
                    "Study different theories of justification",
                    "Examine the relationship between knowledge and truth"
                ]
            }
        },
        {
            id: "node-1-8",
            tree_id: "tree-1",
            parent_id: "node-1-7",
            title: "Empiricism vs Rationalism",
            content: "The fundamental debate about whether knowledge comes primarily from sensory experience or reason.",
            depth: 1,
            position: 0,
            created_at: "2025-06-26T13:30:00Z",
            updated_at: "2025-06-26T13:30:00Z",
            ai_analysis: {
                difficulty: 8,
                concepts: ["empiricism", "rationalism", "a priori knowledge", "sensory experience"],
                suggestions: [
                    "Study Locke's tabula rasa theory",
                    "Examine Descartes' method of doubt",
                    "Compare Hume and Kant on synthetic a priori"
                ]
            }
        }
    ],
    "tree-2": [
        {
            id: "node-2-1",
            tree_id: "tree-2",
            parent_id: null,
            title: "Algorithms",
            content: "Step-by-step procedures for solving computational problems efficiently.",
            depth: 0,
            position: 0,
            created_at: "2025-06-25T14:30:00Z",
            updated_at: "2025-06-26T09:00:00Z",
            ai_analysis: {
                difficulty: 3,
                concepts: ["complexity", "efficiency", "problem-solving", "optimization", "computational thinking"],
                suggestions: [
                    "Study Big O notation for complexity analysis",
                    "Practice implementing common algorithms",
                    "Learn about algorithm design paradigms",
                    "Explore recursive vs iterative approaches"
                ]
            }
        },
        {
            id: "node-2-2",
            tree_id: "tree-2",
            parent_id: "node-2-1",
            title: "Sorting Algorithms",
            content: "Algorithms that arrange elements in a specific order, fundamental to many computational tasks.",
            depth: 1,
            position: 0,
            created_at: "2025-06-25T15:00:00Z",
            updated_at: "2025-06-25T15:00:00Z",
            ai_analysis: {
                difficulty: 3,
                concepts: ["quicksort", "mergesort", "heapsort", "stability"],
                suggestions: [
                    "Compare time and space complexity of different sorting algorithms",
                    "Understand when to use each sorting algorithm"
                ]
            }
        },
        {
            id: "node-2-3",
            tree_id: "tree-2",
            parent_id: "node-2-1",
            title: "Search Algorithms",
            content: "Methods for finding specific elements within data structures efficiently.",
            depth: 1,
            position: 1,
            created_at: "2025-06-25T15:30:00Z",
            updated_at: "2025-06-25T15:30:00Z",
            ai_analysis: {
                difficulty: 2,
                concepts: ["binary search", "linear search", "hash tables", "indexing"],
                suggestions: [
                    "Practice binary search implementation",
                    "Understand search optimization techniques"
                ]
            }
        },
        {
            id: "node-2-4",
            tree_id: "tree-2",
            parent_id: null,
            title: "Data Structures",
            content: "Ways of organizing and storing data to enable efficient operations and access patterns.",
            depth: 0,
            position: 1,
            created_at: "2025-06-25T16:00:00Z",
            updated_at: "2025-06-25T16:00:00Z",
            ai_analysis: {
                difficulty: 5,
                concepts: ["arrays", "linked lists", "trees", "graphs", "hash tables"],
                suggestions: [
                    "Understand the trade-offs between different data structures",
                    "Practice implementing fundamental data structures",
                    "Learn about advanced data structures like B-trees and tries"
                ]
            }
        }
    ],
    "tree-3": [], // Empty tree for testing
    "tree-4": [
        {
            id: "node-4-1",
            tree_id: "tree-4",
            parent_id: null,
            title: "Number Theory",
            content: "The study of integers and their properties, relationships, and patterns.",
            depth: 0,
            position: 0,
            created_at: "2025-06-23T16:15:00Z",
            updated_at: "2025-06-23T16:15:00Z",
            ai_analysis: {
                difficulty: 7,
                concepts: ["prime numbers", "divisibility", "modular arithmetic", "Diophantine equations"],
                suggestions: [
                    "Study the fundamental theorem of arithmetic",
                    "Explore applications in cryptography",
                    "Learn about famous unsolved problems in number theory"
                ]
            }
        },
        {
            id: "node-4-2",
            tree_id: "tree-4",
            parent_id: null,
            title: "Calculus",
            content: "The mathematical study of continuous change, including derivatives and integrals.",
            depth: 0,
            position: 1,
            created_at: "2025-06-23T17:00:00Z",
            updated_at: "2025-06-26T12:00:00Z",
            ai_analysis: {
                difficulty: 7,
                concepts: ["limits", "derivatives", "integrals", "continuity", "infinite series"],
                suggestions: [
                    "Master the fundamental theorem of calculus",
                    "Practice applications in physics and engineering",
                    "Study multivariable calculus extensions",
                    "Explore real analysis foundations"
                ]
            }
        }
    ]
};

// Helper functions for generating realistic test data
export function generateRandomTree() {
    const titles = [
        "Artificial Intelligence", "Psychology", "History", "Literature",
        "Physics", "Chemistry", "Biology", "Art History", "Music Theory",
        "Economics", "Political Science", "Sociology", "Linguistics"
    ];

    const descriptions = [
        "Exploring fundamental concepts and principles",
        "A comprehensive study of key topics",
        "Essential knowledge and practical applications",
        "Core theories and modern developments",
        "Historical context and contemporary relevance"
    ];

    const randomTitle = titles[Math.floor(Math.random() * titles.length)];
    const randomDescription = descriptions[Math.floor(Math.random() * descriptions.length)];

    return {
        id: `tree-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        title: randomTitle,
        description: randomDescription,
        owner_id: "default-user",
        is_public: Math.random() > 0.3,
        created_at: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date().toISOString(),
        node_count: Math.floor(Math.random() * 20)
    };
}

export function generateRandomNode(treeId, parentId = null) {
    const topics = [
        "Introduction", "Basic Concepts", "Advanced Topics", "Applications",
        "Theory", "Practice", "Examples", "Case Studies", "Methodology",
        "Principles", "Frameworks", "Analysis", "Synthesis", "Evaluation"
    ];

    const randomTopic = topics[Math.floor(Math.random() * topics.length)];

    return {
        id: `node-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        tree_id: treeId,
        parent_id: parentId,
        title: randomTopic,
        content: `This is a detailed explanation of ${randomTopic.toLowerCase()} in the context of this knowledge tree. It provides essential information and context for understanding this concept.`,
        depth: parentId ? 1 : 0,
        position: Math.floor(Math.random() * 10),
        created_at: new Date(Date.now() - Math.random() * 24 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date().toISOString(),
        ai_analysis: Math.random() > 0.6 ? {
            difficulty: Math.floor(Math.random() * 10) + 1,
            concepts: ["concept1", "concept2", "concept3"].slice(0, Math.floor(Math.random() * 3) + 1),
            suggestions: [
                "Consider exploring related topics",
                "Review foundational concepts",
                "Practice with examples"
            ].slice(0, Math.floor(Math.random() * 3) + 1)
        } : undefined
    };
}

// Mock API responses
export const mockApiResponses = {
    trees: {
        success: true,
        data: mockTrees,
        meta: {
            total: mockTrees.length,
            limit: 20,
            offset: 0,
            page: 1,
            pages: 1
        }
    },

    getTree: (treeId) => {
        const tree = mockTrees.find(t => t.id === treeId);
        const nodes = mockNodes[treeId] || [];

        if (!tree) {
            return {
                success: false,
                error: "Tree not found"
            };
        }

        return {
            success: true,
            data: {
                ...tree,
                nodes: nodes
            }
        };
    },

    createTree: (treeData) => {
        const newTree = {
            id: `tree-${Date.now()}`,
            ...treeData,
            owner_id: "default-user",
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            node_count: 0
        };

        return {
            success: true,
            data: newTree
        };
    },

    createNode: (nodeData) => {
        const newNode = {
            id: `node-${Date.now()}`,
            ...nodeData,
            depth: nodeData.parent_id ? 1 : 0,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
        };

        return {
            success: true,
            data: newNode
        };
    },

    analyzeNode: () => {
        return {
            success: true,
            data: {
                difficulty: Math.floor(Math.random() * 10) + 1,
                concepts: [
                    "analysis", "understanding", "knowledge", "learning"
                ].slice(0, Math.floor(Math.random() * 4) + 1),
                suggestions: [
                    "Review related concepts for better understanding",
                    "Practice with concrete examples",
                    "Connect to real-world applications",
                    "Explore advanced topics in this area"
                ].slice(0, Math.floor(Math.random() * 4) + 1)
            }
        };
    },

    health: {
        success: true,
        data: {
            status: "healthy",
            database: "healthy",
            localai: "healthy",
            timestamp: new Date().toISOString(),
            version: "1.0.0"
        }
    }
};

// Test mode flag and utilities
export let useTestData = false;

export function enableTestMode() {
    useTestData = true;
    console.log("🧪 Test mode enabled - using mock data");
}

export function disableTestMode() {
    useTestData = false;
    console.log("🌐 Test mode disabled - using real API");
}

export function isTestMode() {
    return useTestData;
}

// Test data API wrapper
export const testApi = {
    async getTrees() {
        await new Promise(resolve => setTimeout(resolve, 300)); // Simulate network delay
        return mockApiResponses.trees;
    },

    async getTree(id) {
        await new Promise(resolve => setTimeout(resolve, 200));
        return mockApiResponses.getTree(id);
    },

    async createTree(treeData) {
        await new Promise(resolve => setTimeout(resolve, 400));
        return mockApiResponses.createTree(treeData);
    },

    async createNode(nodeData) {
        await new Promise(resolve => setTimeout(resolve, 300));
        return mockApiResponses.createNode(nodeData);
    },

    async updateNode(id, nodeData) {
        await new Promise(resolve => setTimeout(resolve, 300));
        return mockApiResponses.createNode({ ...nodeData, id });
    },

    async deleteTree(id) {
        await new Promise(resolve => setTimeout(resolve, 200));
        return { success: true };
    },

    async deleteNode(id) {
        await new Promise(resolve => setTimeout(resolve, 200));
        return { success: true };
    },

    async analyzeNode(id) {
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate AI processing time
        return mockApiResponses.analyzeNode();
    },

    async health() {
        await new Promise(resolve => setTimeout(resolve, 100));
        return mockApiResponses.health;
    }
};

export default {
    mockTrees,
    mockNodes,
    mockApiResponses,
    generateRandomTree,
    generateRandomNode,
    testApi,
    enableTestMode,
    disableTestMode,
    isTestMode
};
