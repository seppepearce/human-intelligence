package seed

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/db"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/models"
)

// SeedData populates the database with demo trees and nodes
func SeedData(database *db.DB) error {
	log.Println("🌱 Seeding database with demo data...")

	// Create Philosophy tree
	philosophyTree, err := createPhilosophyTree(database)
	if err != nil {
		return fmt.Errorf("failed to create philosophy tree: %w", err)
	}

	// Create Computer Science tree
	csTree, err := createComputerScienceTree(database)
	if err != nil {
		return fmt.Errorf("failed to create computer science tree: %w", err)
	}

	// Create Mathematics tree
	mathTree, err := createMathematicsTree(database)
	if err != nil {
		return fmt.Errorf("failed to create mathematics tree: %w", err)
	}

	log.Printf("✅ Successfully seeded database with %d demo trees", 3)
	log.Printf("   - %s (Philosophy)", philosophyTree.ID)
	log.Printf("   - %s (Computer Science)", csTree.ID)
	log.Printf("   - %s (Mathematics)", mathTree.ID)

	return nil
}

func createPhilosophyTree(database *db.DB) (*models.Tree, error) {
	// Create the tree
	tree := models.NewTree(
		"Philosophy Fundamentals",
		"Core philosophical concepts and ideas from ancient to modern times, exploring consciousness, ethics, and epistemology",
		"default-user",
		true,
	)

	if err := database.CreateTree(tree); err != nil {
		return nil, err
	}

	// Create root nodes
	consciousness := models.NewNode(tree.ID, nil, "Consciousness",
		"The subjective experience of awareness and mental states. A fundamental question in philosophy of mind that explores what it means to be aware and have subjective experiences.",
		0, 0)
	consciousness.Concepts = marshalConcepts([]string{"phenomenology", "qualia", "awareness", "subjective experience", "mind-body problem"})
	consciousness.Difficulty = intPtr(9)
	if err := database.CreateNode(consciousness); err != nil {
		return nil, err
	}

	ethics := models.NewNode(tree.ID, nil, "Ethics",
		"The branch of philosophy concerned with moral principles and values that govern behavior. Studies what is right and wrong, good and bad.",
		0, 1)
	ethics.Concepts = marshalConcepts([]string{"morality", "virtue", "duty", "consequentialism", "deontology"})
	ethics.Difficulty = intPtr(6)
	if err := database.CreateNode(ethics); err != nil {
		return nil, err
	}

	epistemology := models.NewNode(tree.ID, nil, "Epistemology",
		"The study of knowledge itself - what it is, how it's acquired, and what makes beliefs justified. The foundation of all learning.",
		0, 2)
	epistemology.Concepts = marshalConcepts([]string{"knowledge", "belief", "justification", "skepticism", "truth"})
	epistemology.Difficulty = intPtr(7)
	if err := database.CreateNode(epistemology); err != nil {
		return nil, err
	}

	// Create child nodes for consciousness
	selfAwareness := models.NewNode(tree.ID, &consciousness.ID, "Self-Awareness",
		"The ability to recognize oneself as an individual distinct from the environment and other individuals. The foundation of personal identity.",
		1, 0)
	selfAwareness.Concepts = marshalConcepts([]string{"self-recognition", "identity", "introspection", "metacognition"})
	selfAwareness.Difficulty = intPtr(5)
	if err := database.CreateNode(selfAwareness); err != nil {
		return nil, err
	}

	streamOfConsciousness := models.NewNode(tree.ID, &consciousness.ID, "Stream of Consciousness",
		"William James's concept describing the continuous flow of thoughts and experiences that make up our mental life.",
		1, 1)
	streamOfConsciousness.Concepts = marshalConcepts([]string{"William James", "consciousness flow", "psychological continuity", "mental states"})
	streamOfConsciousness.Difficulty = intPtr(3)
	if err := database.CreateNode(streamOfConsciousness); err != nil {
		return nil, err
	}

	// Create child nodes for ethics
	utilitarianism := models.NewNode(tree.ID, &ethics.ID, "Utilitarianism",
		"Ethical theory that judges actions by their consequences, aiming for the greatest good for the greatest number of people.",
		1, 0)
	utilitarianism.Concepts = marshalConcepts([]string{"consequentialism", "greatest happiness", "Jeremy Bentham", "John Stuart Mill", "hedonistic calculus"})
	utilitarianism.Difficulty = intPtr(4)
	if err := database.CreateNode(utilitarianism); err != nil {
		return nil, err
	}

	deontology := models.NewNode(tree.ID, &ethics.ID, "Deontological Ethics",
		"Ethical theory focused on duties and rules, most famously developed by Immanuel Kant. Actions are right or wrong regardless of consequences.",
		1, 1)
	deontology.Concepts = marshalConcepts([]string{"categorical imperative", "duty ethics", "Immanuel Kant", "moral law", "autonomy"})
	deontology.Difficulty = intPtr(7)
	if err := database.CreateNode(deontology); err != nil {
		return nil, err
	}

	// Create child nodes for epistemology
	empiricismVsRationalism := models.NewNode(tree.ID, &epistemology.ID, "Empiricism vs Rationalism",
		"The fundamental debate about whether knowledge comes primarily from sensory experience (empiricism) or reason (rationalism).",
		1, 0)
	empiricismVsRationalism.Concepts = marshalConcepts([]string{"empiricism", "rationalism", "a priori knowledge", "sensory experience", "innate ideas"})
	empiricismVsRationalism.Difficulty = intPtr(8)
	if err := database.CreateNode(empiricismVsRationalism); err != nil {
		return nil, err
	}

	return tree, nil
}

func createComputerScienceTree(database *db.DB) (*models.Tree, error) {
	// Create the tree
	tree := models.NewTree(
		"Computer Science Essentials",
		"Fundamental concepts in computer science and programming, covering algorithms, data structures, and computational thinking",
		"default-user",
		true,
	)

	if err := database.CreateTree(tree); err != nil {
		return nil, err
	}

	// Create root nodes
	algorithms := models.NewNode(tree.ID, nil, "Algorithms",
		"Step-by-step procedures for solving computational problems efficiently. The heart of computer science.",
		0, 0)
	algorithms.Concepts = marshalConcepts([]string{"complexity", "efficiency", "problem-solving", "optimization", "computational thinking"})
	algorithms.Difficulty = intPtr(3)
	if err := database.CreateNode(algorithms); err != nil {
		return nil, err
	}

	dataStructures := models.NewNode(tree.ID, nil, "Data Structures",
		"Ways of organizing and storing data to enable efficient operations and access patterns. Foundation of all programs.",
		0, 1)
	dataStructures.Concepts = marshalConcepts([]string{"arrays", "linked lists", "trees", "graphs", "hash tables", "memory management"})
	dataStructures.Difficulty = intPtr(5)
	if err := database.CreateNode(dataStructures); err != nil {
		return nil, err
	}

	// Create child nodes for algorithms
	sorting := models.NewNode(tree.ID, &algorithms.ID, "Sorting Algorithms",
		"Algorithms that arrange elements in a specific order. Fundamental to many computational tasks and a great introduction to algorithm analysis.",
		1, 0)
	sorting.Concepts = marshalConcepts([]string{"quicksort", "mergesort", "heapsort", "stability", "time complexity"})
	sorting.Difficulty = intPtr(2)
	if err := database.CreateNode(sorting); err != nil {
		return nil, err
	}

	searching := models.NewNode(tree.ID, &algorithms.ID, "Search Algorithms",
		"Methods for finding specific elements within data structures efficiently. Essential for information retrieval.",
		1, 1)
	searching.Concepts = marshalConcepts([]string{"binary search", "linear search", "hash tables", "indexing", "search optimization"})
	searching.Difficulty = intPtr(2)
	if err := database.CreateNode(searching); err != nil {
		return nil, err
	}

	// Create child nodes for data structures
	arrays := models.NewNode(tree.ID, &dataStructures.ID, "Arrays & Lists",
		"Sequential collections of elements. The most basic and widely used data structures in programming.",
		1, 0)
	arrays.Concepts = marshalConcepts([]string{"contiguous memory", "indexing", "dynamic arrays", "linked lists", "iteration"})
	arrays.Difficulty = intPtr(1)
	if err := database.CreateNode(arrays); err != nil {
		return nil, err
	}

	trees := models.NewNode(tree.ID, &dataStructures.ID, "Trees & Graphs",
		"Hierarchical and networked data structures that model complex relationships between data elements.",
		1, 1)
	trees.Concepts = marshalConcepts([]string{"binary trees", "graph theory", "traversal", "nodes", "edges"})
	trees.Difficulty = intPtr(6)
	if err := database.CreateNode(trees); err != nil {
		return nil, err
	}

	return tree, nil
}

func createMathematicsTree(database *db.DB) (*models.Tree, error) {
	// Create the tree
	tree := models.NewTree(
		"Mathematics Foundation",
		"Essential mathematical concepts and theories that form the basis of scientific and logical thinking",
		"default-user",
		true,
	)

	if err := database.CreateTree(tree); err != nil {
		return nil, err
	}

	// Create root nodes
	numberTheory := models.NewNode(tree.ID, nil, "Number Theory",
		"The study of integers and their properties, relationships, and patterns. Often called the 'queen of mathematics'.",
		0, 0)
	numberTheory.Concepts = marshalConcepts([]string{"prime numbers", "divisibility", "modular arithmetic", "Diophantine equations", "cryptography"})
	numberTheory.Difficulty = intPtr(8)
	if err := database.CreateNode(numberTheory); err != nil {
		return nil, err
	}

	calculus := models.NewNode(tree.ID, nil, "Calculus",
		"The mathematical study of continuous change, including derivatives and integrals. Essential for understanding rates of change.",
		0, 1)
	calculus.Concepts = marshalConcepts([]string{"limits", "derivatives", "integrals", "continuity", "infinite series"})
	calculus.Difficulty = intPtr(7)
	if err := database.CreateNode(calculus); err != nil {
		return nil, err
	}

	algebra := models.NewNode(tree.ID, nil, "Abstract Algebra",
		"The study of algebraic structures such as groups, rings, and fields. Generalizes arithmetic operations.",
		0, 2)
	algebra.Concepts = marshalConcepts([]string{"groups", "rings", "fields", "symmetry", "abstract operations"})
	algebra.Difficulty = intPtr(9)
	if err := database.CreateNode(algebra); err != nil {
		return nil, err
	}

	// Create child nodes
	primes := models.NewNode(tree.ID, &numberTheory.ID, "Prime Numbers",
		"Natural numbers greater than 1 that have no positive divisors other than 1 and themselves. Building blocks of all integers.",
		1, 0)
	primes.Concepts = marshalConcepts([]string{"fundamental theorem", "sieve methods", "prime distribution", "cryptographic applications"})
	primes.Difficulty = intPtr(6)
	if err := database.CreateNode(primes); err != nil {
		return nil, err
	}

	derivatives := models.NewNode(tree.ID, &calculus.ID, "Derivatives",
		"Mathematical objects that measure how a function changes as its input changes. The instantaneous rate of change.",
		1, 0)
	derivatives.Concepts = marshalConcepts([]string{"rate of change", "limits", "differentiation rules", "chain rule", "optimization"})
	derivatives.Difficulty = intPtr(5)
	if err := database.CreateNode(derivatives); err != nil {
		return nil, err
	}

	return tree, nil
}

// Helper functions
func marshalConcepts(concepts []string) json.RawMessage {
	data, _ := json.Marshal(concepts)
	return data
}

func intPtr(i int) *int {
	return &i
}

// CheckIfSeeded checks if the database already has seed data
func CheckIfSeeded(database *db.DB) (bool, error) {
	// Check if we have any trees with specific seed titles
	trees, _, err := database.SearchTrees(&models.TreeSearchRequest{
		Query:     "Philosophy Fundamentals",
		Limit:     1,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "desc",
	})
	if err != nil {
		return false, err
	}

	return len(trees) > 0, nil
}

// ClearSeedData removes all seed data from the database
func ClearSeedData(database *db.DB) error {
	log.Println("🧹 Clearing seed data...")

	// Get all trees and delete seed ones
	trees, _, err := database.SearchTrees(&models.TreeSearchRequest{
		Limit:     100,
		SortBy:    "created_at",
		SortOrder: "desc",
	})
	if err != nil {
		return err
	}

	seedTitles := map[string]bool{
		"Philosophy Fundamentals":     true,
		"Computer Science Essentials": true,
		"Mathematics Foundation":      true,
	}

	deletedCount := 0
	for _, tree := range trees {
		if seedTitles[tree.Title] {
			if err := database.DeleteTree(tree.ID); err != nil {
				log.Printf("Warning: failed to delete seed tree %s: %v", tree.Title, err)
			} else {
				deletedCount++
			}
		}
	}

	log.Printf("✅ Cleared %d seed trees", deletedCount)
	return nil
}
