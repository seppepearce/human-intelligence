Keywords: #knowledge #opensource #collaboration #social-wikipedia 

A social information-collaboration platform called human intelligence, where people create trees / paths of knowledge.

Features: 
- Trees
	* Users can plant a seed (root node) 
	* Each node contains information in markdown supported format 
	* You can build on a node by either building on the existing path, or branching off from it. 
	* Trees can be public or private
	* When a tree is public, users can collaborate on the existing tree by branching of from a particular node. When they are happy with their changes, they can request the original nodes owner to make the changes definitive and visible for everyone else. 
	* Within each node, comments and questions can be made, nodes and trees can be liked. 
	* We can open the tree and start exploring it dungeon style, when traversing the tree, you can choose the knowledge path you walk, or you can pave your own
	* Each node has a tldr system, people can create their own summary of the node to better understand it, and to help other people better understand that node. TLDR's can be voted for, and the top picks will show up more often. 
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



Possible use cases could be: 
- Personal projects / second brain. You could also just build out your own private knowledge if you dont wish to share it with others.
- Private use in a classroom, for students to collaborate on subject matter. 
- News stories, A storyline could develop and nodes could be fact checked / voted on 
* Just all around having fun developing and sharing knowledge together with other people about esoteric things 
* Research, tree structures and nodes could be combined to create new insights
* Potential employers looking for knowledgeable people in certain fields
* Educators could create their own roadmaps to follow in an interactive way
  
  
  MVP:
## MVP Feature Set

**Core MVP (Phase 1):**

- Basic tree creation and node editing (markdown support)
- Simple branching from existing nodes
- Public/private tree toggle
- Basic search by title/tags
- User profiles with contribution history
- Content classification system (fact/hypothesis/myth)
- Source citation system

**Leave for Later:**

- AI features, gamification, real-time updates, advanced visualizations, merge requests

## Content Moderation System

Your fact/hypothesis/myth classification is smart - it's transparent and educational rather than just censorial. Here's how to structure it:

**Content Classification:**

- **Fact**: Well-established, widely accepted information
- **Hypothesis**: Proposed explanations or theories under investigation
- **Myth**: Debunked or widely disputed claims
- **Opinion**: Personal viewpoints or interpretations
- **Incomplete**: Partial information that needs expansion

**Implementation:**

- Node creators set initial classification
- Community can suggest reclassification (requires voting threshold)
- Citations required for "Fact" classification
- Visual indicators (colored borders, icons) make classification immediately obvious
- "Controversy score" based on classification disagreements

**Source System:**

- Mandatory citation fields for facts
- Source quality indicators (academic journals, primary sources, news, etc.)
- Community can rate source reliability
- Easy citation formatting tools

## Monetization Strategy

**Free Tier:**

- Create unlimited public trees
- Join public collaborations
- Basic search and discovery
- Up to 10 private trees

**Creator Plan ($9/month):**

- Unlimited private trees
- Early access to AI features (tree suggestions, auto-citations)
- Analytics on tree engagement
- Monetization tools (paid access to premium trees)
- Priority support

**Institution Plan ($49/month + per-user):**

- Private organization workspaces
- Advanced collaboration tools
- Admin controls and user management
- Integration APIs
- Custom branding options
- Usage analytics dashboard

## MVP Development Roadmap

**Week 1-4: Core Infrastructure**

- User authentication
- Basic tree/node CRUD operations
- Simple tree visualization
- Markdown editor

**Week 5-8: Collaboration Features**

- Public/private trees
- Basic branching
- User profiles
- Content classification system

**Week 9-12: Discovery & Polish**

- Search functionality
- Source citation system
- Basic moderation tools
- Mobile responsiveness

## Educator Acquisition Strategy

**Content Creator Program:**

- Invite educators to migrate existing course materials
- Revenue sharing on premium educational trees
- Featured educator spotlights
- "Verified Educator" badges for credibility

**Initial Outreach:**

- Start with independent educators, tutors, online course creators
- Offer migration assistance for existing content
- Create case studies showing engagement improvements
- Target educators frustrated with static platforms

## Technical Considerations for MVP

**Data Structure:**

```
Tree {
  id, title, description, tags[], isPrivate, ownerId, createdAt
}

Node {
  id, treeId, parentNodeId, title, content, classification,
  sources[], authorId, createdAt, likes, views
}

Source {
  url, title, author, type, reliability_score
}
```

**Key Technical Decisions:**

- Use a graph database (Neo4j) or add graph capabilities to relational DB
- Real-time collaboration can wait - focus on async collaboration first
- Simple tree visualization (D3.js or similar) - fancy visuals come later

## Early Validation Steps

**Before Building:**

1. Create mockups and test with 10-15 educators
2. Survey target users about current pain points
3. Validate pricing with potential institutional customers

**After MVP Launch:**

1. Onboard 5-10 educator partners manually
2. Track key metrics: trees created, nodes per tree, collaboration events
3. Focus on one subject area initially (maybe computer science or biology)

## Risk Mitigation

**Content Quality:**

- Start with invited educators only
- Manual review of first 100 trees
- Clear community guidelines from day one
- Easy reporting mechanisms

**User Adoption:**

- Provide migration tools from existing platforms
- Create compelling demo content
- Focus on solving real educator pain points