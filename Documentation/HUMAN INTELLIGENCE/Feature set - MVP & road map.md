 
## Phase 1: MVP feature set

Summary:
- Create trees / nodes / branches with a batteries included markdown editor
- Architecture right from the start (Neo4j, Postgres, GraphQL, Go, Svelte)
- Basic search functionality

#### 1. **Tree Creation** 
- A **user** can create a new **Tree**, which serves as a top-level container for a structured set of information.
- Upon creation, a tree must include:
    - **Name** (string, required): The title of the tree.
    - **Description**
    - **Visibility** (enum, required): `Public` or `Private`.
    - **Created at** (date, required): 
#### 2. **Tree Structure**
- Each tree always starts with a single **Root Node**, which serves as:
    - An **index node** providing an overview of what the tree is meant to represent, pointers into the tree’s content.
- A **Tree** is composed of interconnected **Nodes** arranged in a hierarchical structure.
- Nodes can have zero or more child nodes, forming a branching structure.
#### 3. **Nodes**
- Each **Node** is a content unit within the tree.
- A node contains:
	- A title
	- A created date
	- A last modified date
    - A **Markdown document**: Supports rich text and images, allowing users to write content, embed visuals, and structure information freely.
- Nodes can be created or edited by users with appropriate permissions (typically the creator or collaborators, depending on privacy settings).
#### 4. **Branches**
- A **Branch** is a user-defined **linear path** through nodes, starting from the root or any node in the tree.
    
- Key properties of branches:
    - **User-specific**: Each user can create their own unique branches through a tree.
    - **Navigable**: Other users can explore branches created by different users.
    - **Semantic purpose**:
        - Represent a curated flow or chain of information.
        - Help break complex topics into digestible sequences.
        - Useful for roadmaps, learning paths, thematic guides, etc.
            
- **Branch Divergence**:
    - Multiple branches can originate from the same node.
    - Divergent branches (by the same or different users) allow alternative perspectives or paths within the same thematic scope.
      
#### 5. Basic search
- Search for trees by name
- Search trees for node by name

**Key Technical Decisions:**

- Use a graph database (Neo4j) or add graph capabilities to relational DB
- Real-time collaboration can wait - focus on async collaboration first
- Simple tree visualization (D3.js or similar) - fancy visuals come later
  
More: [[Technical architecture & design philosophy]]


## Phase 2: Collaboration

#### 1. Authentication, authorization
- User can sign up / sign in using social provider or email / password
- User can edit their account
- User can delete their account
- Users are authenticated using JWT tokens
- User information is stored in postgres + uuid in neo4j for tracking relations interactions, collaborations
#### 2. Collaboration
- Simple branching from existing nodes
- User profiles with contribution history
- Source citation system

#### 3. Interactions
- Nodes can be reacted to by users, and have a discussion sections, there are several different actions a user can take in response to a node:

Predefined reactions:

| Emoji | Label                | Meaning                           |
| ----- | -------------------- | --------------------------------- |
| 💡    | **Insightful**       | Strong idea or clever framing     |
| 🌱    | **Growing**          | Node shows promise                |
| ❤️    | **Appreciated**      | Social gratitude, support         |
| 🗝️   | **Key Contribution** | Essential / major input           |
| ⚠️    | **Needs Work**       | Confusing, incomplete, misleading |

Discussion types:

| Emoji | Label             | Purpose                                   |
| ----- | ----------------- | ----------------------------------------- |
| ❓     | **Question**      | Ask for clarification or details          |
| ⚔️    | **Challenge**     | Raise doubts or provide counterpoints     |
| ➕     | **Expand**        | Add new angles or information             |
| 📎    | **Contextualize** | Supply sources, history, or adjacent info |

## Phase 3: Leave for later
- Using tags for advanced search and linking between trees / nodes
- Real-time updates / discussions
- Advanced visualisations
- Recommendation engine
- AI features
- Gamification: levels, contribution points, boosts, daily streaks, ...
- Content moderation
- Citation / source system
- Controversy score
- Private collaboration / invites

## Idea: Content Moderation

**Node Classification:**

- **Fact**: Well-established, widely accepted information
- **Hypothesis**: Proposed explanations or theories under investigation
- **Myth**: Debunked or widely disputed claims
- **Opinion**: Personal viewpoints or interpretations
- **Incomplete**: Partial information that needs expansion
- **Fiction**: Prose / Poetry / Creative writing (Cant be reclassified or be built upon with facts or any other classification)

**Implementation:**
- Node creators set initial classification
- Community can suggest reclassification (requires voting threshold)
- Keep logs of reclassification open
- Citations required for keeping "Fact" classification
- Visual indicators (colored borders, icons) make classification immediately obvious
- Controversy score based on classification disagreements
- Citation formatting tools
- Community can rate source reliability (one aspect of peer review)



- Each node has a tldr system, people can create their own summary of the node to better understand it, and to help other people better understand that node. TLDR's can be voted for, and the top picks will show up more often. 
* Design
	* The website has a neutral academic like look, but really shines in the visual representation of the tree structures, these can be customised and are the eyecatcher of the app.
* Search:
	* Users can search for knowledge trees, trees and nodes can have tags which indicate what the topic is about
	* Lots of filter options to find the exact trees you wish to participate in
	* You can search for users which have created and collaborated on knowledge in a certain area.
	* If you find no tree to your liking, you will get the option to create your own tree. AI can help you make a base layout of what paths you could take
* Realtime:
	* On the home page, we can see trending trees which are being worked on at the moment to really sell the feeling that the platform is breathing with life. (counters going up, nodes popping up.
* AI
* User:
	* Your profile is built upon throughout your usage. It contains information about what your knowledge areas are, how much you've spent researching them.
	* It also contains the trees you have built, and other peoples trees you have traversed and spent time collaborating on.
	* You are incentiviced to collaborate and work on stuff through a gamified system, you can earn points/levels keeping up your daily streak, getting popular trees / nodes / tldrs, and just overall collaborating on public trees.




  
  