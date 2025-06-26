package seeds

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Seeder handles database seeding operations
type Seeder struct {
	db *sql.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedAll runs all seeding operations
func (s *Seeder) SeedAll() error {
	log.Println("🌱 Starting database seeding...")

	if err := s.ClearData(); err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	if err := s.SeedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := s.SeedNodes(); err != nil {
		return fmt.Errorf("failed to seed nodes: %w", err)
	}

	if err := s.SeedLearningPaths(); err != nil {
		return fmt.Errorf("failed to seed learning paths: %w", err)
	}

	if err := s.SeedVotes(); err != nil {
		return fmt.Errorf("failed to seed votes: %w", err)
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

// ClearData removes all seeded data (for development)
func (s *Seeder) ClearData() error {
	log.Println("🧹 Clearing existing data...")

	tables := []string{
		"path_nodes",
		"votes",
		"learning_paths",
		"nodes",
		"users",
	}

	for _, table := range tables {
		_, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}

	return nil
}

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// SeedUsers creates sample users
func (s *Seeder) SeedUsers() error {
	log.Println("👥 Seeding users...")

	users := []struct {
		username string
		email    string
		bio      string
		password string
	}{
		{
			username: "alice_dev",
			email:    "alice@example.com",
			bio:      "Full-stack developer with 8+ years experience. Passionate about teaching and clean code.",
			password: "password123",
		},
		{
			username: "bob_ml",
			email:    "bob@example.com",
			bio:      "Machine learning engineer and data scientist. Love making AI accessible to everyone.",
			password: "password123",
		},
		{
			username: "charlie_design",
			email:    "charlie@example.com",
			bio:      "UX/UI designer focused on creating beautiful and functional user experiences.",
			password: "password123",
		},
		{
			username: "diana_blockchain",
			email:    "diana@example.com",
			bio:      "Blockchain developer and smart contract specialist. Building the decentralized future.",
			password: "password123",
		},
		{
			username: "eve_api",
			email:    "eve@example.com",
			bio:      "Backend engineer specializing in API design and microservices architecture.",
			password: "password123",
		},
		{
			username: "frank_frontend",
			email:    "frank@example.com",
			bio:      "Frontend developer with expertise in React, Vue, and modern JavaScript.",
			password: "password123",
		},
	}

	for _, user := range users {
		userID := uuid.New().String()

		// Hash the password
		passwordHash, err := hashPassword(user.password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", user.username, err)
		}

		_, err = s.db.Exec(`
			INSERT INTO users (id, username, email, password_hash, bio, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, userID, user.username, user.email, passwordHash, user.bio, time.Now(), time.Now())

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.username, err)
		}
	}

	return nil
}

// SeedNodes creates sample learning nodes
func (s *Seeder) SeedNodes() error {
	log.Println("📚 Seeding nodes...")

	// Get user IDs for attribution
	users, err := s.getUserMap()
	if err != nil {
		return err
	}

	nodes := []struct {
		title       string
		description string
		nodeType    string
		content     string
		tags        []string
		author      string
		metadata    map[string]interface{}
		isPublic    bool
	}{
		{
			title:       "Introduction to JavaScript",
			description: "Learn the basics of JavaScript programming language",
			nodeType:    "text",
			content: `# Introduction to JavaScript

JavaScript is a versatile programming language that powers the web. Originally created for client-side scripting, it has evolved into a full-stack language.

## Key Features

- **Dynamic typing**: Variables can hold different types of values
- **First-class functions**: Functions are values that can be passed around
- **Prototype-based**: Objects can inherit directly from other objects
- **Event-driven**: Perfect for handling user interactions

## Basic Syntax

Variables and functions:
- let name = "Alice";
- const age = 25;
- function greet(name) { return "Hello, " + name + "!"; }

JavaScript is the foundation of modern web development!`,
			tags:     []string{"javascript", "programming", "basics", "web-development"},
			author:   "alice_dev",
			metadata: map[string]interface{}{},
			isPublic: true,
		},
		{
			title:       "React Hooks Deep Dive",
			description: "Master React Hooks with practical examples and best practices",
			nodeType:    "video",
			content: `# React Hooks Deep Dive

This comprehensive video covers all essential React Hooks and how to use them effectively.

## Topics Covered

1. **useState** - Managing component state
2. **useEffect** - Side effects and lifecycle
3. **useContext** - Sharing state across components
4. **useReducer** - Complex state management
5. **Custom Hooks** - Creating reusable logic

## Key Takeaways

- Hooks let you use state and other React features without writing a class
- Always call hooks at the top level of your React function
- Custom hooks are a powerful way to share stateful logic

Perfect for intermediate React developers looking to level up!`,
			tags:   []string{"react", "hooks", "javascript", "frontend"},
			author: "frank_frontend",
			metadata: map[string]interface{}{
				"video_url":      "https://youtube.com/watch?v=example123",
				"video_duration": 3600,
			},
			isPublic: true,
		},
		{
			title:       "Machine Learning Fundamentals",
			description: "Complete guide to getting started with machine learning",
			nodeType:    "text",
			content: `# Machine Learning Fundamentals

Machine Learning is a subset of artificial intelligence that enables computers to learn and make decisions from data.

## Types of Machine Learning

### 1. Supervised Learning
- Uses labeled training data
- Examples: Classification, Regression
- Algorithms: Linear Regression, Decision Trees, Neural Networks

### 2. Unsupervised Learning
- Works with unlabeled data
- Examples: Clustering, Dimensionality Reduction
- Algorithms: K-Means, PCA, Autoencoders

### 3. Reinforcement Learning
- Learning through interaction and feedback
- Examples: Game AI, Robotics
- Algorithms: Q-Learning, Policy Gradients

## Getting Started

1. **Learn the Math**: Statistics, Linear Algebra, Calculus
2. **Choose a Language**: Python, R, or Julia
3. **Practice with Data**: Kaggle, UCI ML Repository
4. **Build Projects**: Start simple and gradually increase complexity

The key is consistent practice and understanding the fundamentals!`,
			tags:     []string{"machine-learning", "ai", "python", "data-science"},
			author:   "bob_ml",
			metadata: map[string]interface{}{},
			isPublic: true,
		},
		{
			title:       "Building RESTful APIs with Node.js",
			description: "Learn to create scalable REST APIs using Node.js and Express",
			nodeType:    "code",
			content: `# Building RESTful APIs with Node.js

Learn how to build robust, scalable REST APIs using Node.js and Express.

## Project Setup

Install dependencies:
- npm init -y
- npm install express cors helmet morgan
- npm install -D nodemon

## Basic Server Setup

Create a basic Express server with middleware:
- const express = require('express');
- const app = express();
- app.use(express.json());
- app.listen(3000);

## Best Practices

1. **Use middleware** for common functionality
2. **Validate input** to prevent security issues
3. **Handle errors** gracefully
4. **Use environment variables** for configuration
5. **Implement rate limiting** to prevent abuse

This foundation will serve you well for any API project!`,
			tags:   []string{"nodejs", "api", "express", "backend", "javascript"},
			author: "eve_api",
			metadata: map[string]interface{}{
				"language":   "javascript",
				"repository": "https://github.com/example/nodejs-api",
			},
			isPublic: true,
		},
		{
			title:       "UI/UX Design Principles",
			description: "Essential design principles for creating user-friendly interfaces",
			nodeType:    "link",
			content: `# UI/UX Design Principles

This comprehensive guide covers the fundamental principles every designer should know.

## Key Principles Covered

1. **User-Centered Design**: Always prioritize user needs
2. **Consistency**: Maintain visual and functional consistency
3. **Hierarchy**: Guide users through clear information architecture
4. **Accessibility**: Design for all users, including those with disabilities
5. **Feedback**: Provide clear feedback for user actions

## Tools and Resources

- **Figma**: For design and prototyping
- **Sketch**: Mac-based design tool
- **Adobe XD**: Complete design system
- **InVision**: Prototyping and collaboration

## Additional Reading

The linked resource provides interactive examples and case studies that demonstrate these principles in action.

Perfect for both beginners and experienced designers looking to refresh their fundamentals!`,
			tags:   []string{"design", "ux", "ui", "figma", "user-experience"},
			author: "charlie_design",
			metadata: map[string]interface{}{
				"link_url":   "https://lawsofux.com",
				"link_title": "Laws of UX - Complete Design Principles Guide",
			},
			isPublic: true,
		},
		{
			title:       "Blockchain Development with Solidity",
			description: "Build smart contracts and DApps on the Ethereum blockchain",
			nodeType:    "text",
			content: `# Blockchain Development with Solidity

Learn to build decentralized applications (DApps) and smart contracts using Solidity.

## What is Solidity?

Solidity is a high-level programming language designed for implementing smart contracts on the Ethereum Virtual Machine (EVM).

## Basic Smart Contract

A simple storage contract:
- pragma solidity ^0.8.0;
- contract SimpleStorage {
-   uint256 private storedData;
-   function set(uint256 value) public { storedData = value; }
-   function get() public view returns (uint256) { return storedData; }
- }

## Key Concepts

1. **Gas**: Computational cost for operations
2. **Events**: Logging mechanism for the blockchain
3. **Modifiers**: Reusable code to change function behavior
4. **Inheritance**: Contracts can inherit from other contracts

## Development Tools

- **Remix**: Browser-based IDE
- **Hardhat**: Development environment
- **Truffle**: Framework for Ethereum development
- **MetaMask**: Browser wallet for testing

## Security Best Practices

- Always validate inputs
- Be aware of reentrancy attacks
- Use established patterns and libraries
- Test thoroughly on testnets

The blockchain space is evolving rapidly - stay curious and keep learning!`,
			tags:     []string{"blockchain", "ethereum", "solidity", "smart-contracts", "web3"},
			author:   "diana_blockchain",
			metadata: map[string]interface{}{},
			isPublic: true,
		},
	}

	for _, node := range nodes {
		nodeID := uuid.New().String()
		userID := users[node.author]

		// Convert metadata to JSON
		var metadataJSON []byte
		var err error
		if len(node.metadata) > 0 {
			metadataJSON, err = json.Marshal(node.metadata)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata for node %s: %w", node.title, err)
			}
		} else {
			metadataJSON = []byte("{}")
		}

		_, err = s.db.Exec(`
			INSERT INTO nodes (
				id, title, description, node_type, content, tags, metadata,
				created_by, is_public, created_at, updated_at, vote_score, view_count
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, nodeID, node.title, node.description, node.nodeType, node.content,
			pq.Array(node.tags), metadataJSON, userID, node.isPublic,
			time.Now(), time.Now(), 0, 0)

		if err != nil {
			return fmt.Errorf("failed to insert node %s: %w", node.title, err)
		}
	}

	return nil
}

// SeedLearningPaths creates sample learning paths
func (s *Seeder) SeedLearningPaths() error {
	log.Println("🛤️  Seeding learning paths...")

	users, err := s.getUserMap()
	if err != nil {
		return err
	}

	nodes, err := s.getNodeMap()
	if err != nil {
		return err
	}

	paths := []struct {
		title          string
		description    string
		category       string
		difficulty     string
		author         string
		tags           []string
		isPublic       bool
		estimatedHours int
		nodeSequence   []string // Node titles in order
	}{
		{
			title:          "Full Stack Web Development",
			description:    "Complete journey from frontend to backend development with modern JavaScript stack",
			category:       "programming",
			difficulty:     "intermediate",
			author:         "alice_dev",
			tags:           []string{"javascript", "react", "nodejs", "fullstack"},
			isPublic:       true,
			estimatedHours: 40,
			nodeSequence: []string{
				"Introduction to JavaScript",
				"React Hooks Deep Dive",
				"Building RESTful APIs with Node.js",
			},
		},
		{
			title:          "Machine Learning Fundamentals",
			description:    "Learn the basics of ML from theory to practice",
			category:       "data-science",
			difficulty:     "beginner",
			author:         "bob_ml",
			tags:           []string{"machine-learning", "python", "ai", "data-science"},
			isPublic:       true,
			estimatedHours: 25,
			nodeSequence: []string{
				"Machine Learning Fundamentals",
			},
		},
		{
			title:          "UI/UX Design Mastery",
			description:    "From wireframes to prototypes and user testing",
			category:       "design",
			difficulty:     "intermediate",
			author:         "charlie_design",
			tags:           []string{"design", "ux", "ui", "figma"},
			isPublic:       true,
			estimatedHours: 30,
			nodeSequence: []string{
				"UI/UX Design Principles",
			},
		},
		{
			title:          "Blockchain Development",
			description:    "Build decentralized applications from scratch",
			category:       "programming",
			difficulty:     "advanced",
			author:         "diana_blockchain",
			tags:           []string{"blockchain", "ethereum", "solidity", "web3"},
			isPublic:       true,
			estimatedHours: 50,
			nodeSequence: []string{
				"Blockchain Development with Solidity",
			},
		},
	}

	for _, path := range paths {
		pathID := uuid.New().String()
		userID := users[path.author]

		// Insert learning path
		_, err := s.db.Exec(`
			INSERT INTO learning_paths (
				id, name, description, access_level, created_by, tags,
				created_at, updated_at, vote_score, fork_count, completion_count
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`, pathID, path.title, path.description, func() string {
			if path.isPublic {
				return "public"
			}
			return "private"
		}(), userID, pq.Array(path.tags), time.Now(), time.Now(), 0, 0, 0)

		if err != nil {
			return fmt.Errorf("failed to insert learning path %s: %w", path.title, err)
		}

		// Insert path nodes
		for order, nodeTitle := range path.nodeSequence {
			if nodeID, exists := nodes[nodeTitle]; exists {
				_, err := s.db.Exec(`
					INSERT INTO path_nodes (path_id, node_id, position, is_required)
					VALUES ($1, $2, $3, $4)
				`, pathID, nodeID, order, true)

				if err != nil {
					return fmt.Errorf("failed to insert path node: %w", err)
				}
			}
		}
	}

	return nil
}

// SeedVotes creates sample votes and interactions
func (s *Seeder) SeedVotes() error {
	log.Println("👍 Seeding votes and interactions...")

	users, err := s.getUserMap()
	if err != nil {
		return err
	}

	nodes, err := s.getNodeMap()
	if err != nil {
		return err
	}

	// Add some votes to nodes
	votePatterns := []struct {
		nodeTitle string
		votes     int
		views     int
	}{
		{"Introduction to JavaScript", 45, 1200},
		{"React Hooks Deep Dive", 38, 890},
		{"Machine Learning Fundamentals", 52, 1100},
		{"Building RESTful APIs with Node.js", 41, 756},
		{"UI/UX Design Principles", 29, 634},
		{"Blockchain Development with Solidity", 33, 445},
	}

	for _, pattern := range votePatterns {
		if nodeID, exists := nodes[pattern.nodeTitle]; exists {
			// Update vote and view counts
			_, err := s.db.Exec(`
				UPDATE nodes SET vote_score = $1, view_count = $2 WHERE id = $3
			`, pattern.votes, pattern.views, nodeID)

			if err != nil {
				return fmt.Errorf("failed to update node stats: %w", err)
			}

			// Create some individual vote records
			userList := []string{"alice_dev", "bob_ml", "charlie_design", "diana_blockchain"}
			for i, username := range userList {
				if i >= pattern.votes/10 { // Only create votes for some users
					break
				}

				if userID, exists := users[username]; exists {
					voteID := uuid.New().String()
					_, err := s.db.Exec(`
						INSERT INTO votes (id, user_id, target_id, target_type, vote_type, created_at)
						VALUES ($1, $2, $3, $4, $5, $6)
					`, voteID, userID, nodeID, "node", 1, time.Now())

					if err != nil {
						return fmt.Errorf("failed to insert vote: %w", err)
					}
				}
			}
		}
	}

	return nil
}

// Helper functions

func (s *Seeder) getUserMap() (map[string]string, error) {
	users := make(map[string]string)

	rows, err := s.db.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, username string
		if err := rows.Scan(&id, &username); err != nil {
			return nil, err
		}
		users[username] = id
	}

	return users, nil
}

func (s *Seeder) getNodeMap() (map[string]string, error) {
	nodes := make(map[string]string)

	rows, err := s.db.Query("SELECT id, title FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}
		nodes[title] = id
	}

	return nodes, nil
}
