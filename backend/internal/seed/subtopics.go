package seed

import "puyrg/backend/internal/models"

// subtopicNodes returns all subtopic/pattern nodes parented under topics.
func subtopicNodes(nodeMap map[string]uint) []models.KnowledgeNode {
	p := func(id uint) *uint { return &id }
	e := models.DifficultyEasy
	m := models.DifficultyMedium
	h := models.DifficultyHard
	x := models.DifficultyExpert

	var nodes []models.KnowledgeNode

	add := func(parentSlug, domainSlug, name, slug string, diff models.DifficultyLevel, mins, rev, order int) {
		parentID := nodeMap[parentSlug]
		domainID := nodeMap[domainSlug]
		if parentID == 0 {
			return
		}
		nodes = append(nodes, models.KnowledgeNode{
			ParentID: p(parentID), DomainID: p(domainID),
			Type: models.NodeTypePattern, Name: name, Slug: slug,
			Difficulty: diff, EstimatedMinutes: mins,
			RevisionIntervalDays: rev, SortOrder: order,
		})
	}

	// ── Arrays ───────────────────────────────────────────────────────────────
	add("arrays","dsa","Kadane's Algorithm","arrays-kadane",m,60,7,1)
	add("arrays","dsa","Dutch National Flag","arrays-dnf",m,60,7,2)
	add("arrays","dsa","Difference Array","arrays-diff",m,60,7,3)
	add("arrays","dsa","Circular Array","arrays-circular",m,60,10,4)
	add("arrays","dsa","Next Greater Element","arrays-nge",m,60,7,5)
	add("arrays","dsa","Product Array","arrays-product",m,60,7,6)
	add("arrays","dsa","Merge Intervals","arrays-merge-intervals",m,90,7,7)
	add("arrays","dsa","Rotate Array","arrays-rotate",m,60,7,8)
	add("arrays","dsa","Subarray Problems","arrays-subarray",m,90,7,9)
	add("arrays","dsa","Majority Element","arrays-majority",e,45,7,10)

	// ── Strings ──────────────────────────────────────────────────────────────
	add("strings","dsa","Palindrome Detection","strings-palindrome",e,60,7,1)
	add("strings","dsa","Anagram Check","strings-anagram",e,45,7,2)
	add("strings","dsa","String Hashing","strings-hashing",m,90,7,3)
	add("strings","dsa","KMP Algorithm","strings-kmp",h,120,10,4)
	add("strings","dsa","Z Algorithm","strings-z",h,120,10,5)
	add("strings","dsa","Rabin-Karp","strings-rabin-karp",h,120,10,6)
	add("strings","dsa","Longest Common Prefix","strings-lcp",m,60,10,7)
	add("strings","dsa","String Compression","strings-compression",m,60,7,8)
	add("strings","dsa","Wildcard Matching","strings-wildcard",h,90,10,9)
	add("strings","dsa","Parentheses Matching","strings-parens",m,60,7,10)

	// ── Binary Search ─────────────────────────────────────────────────────────
	add("binary-search","dsa","Classic Binary Search","bs-classic",e,45,7,1)
	add("binary-search","dsa","Lower Bound / Upper Bound","bs-bounds",m,60,7,2)
	add("binary-search","dsa","First / Last Occurrence","bs-first-last",m,60,7,3)
	add("binary-search","dsa","Answer Binary Search","bs-answer",m,90,7,4)
	add("binary-search","dsa","Rotated Array","bs-rotated",m,90,7,5)
	add("binary-search","dsa","Peak Element","bs-peak",m,60,7,6)
	add("binary-search","dsa","Binary Search on Matrix","bs-matrix",m,90,7,7)
	add("binary-search","dsa","Floating Point Binary Search","bs-float",m,60,10,8)
	add("binary-search","dsa","Binary Search + Greedy","bs-greedy",h,120,7,9)
	add("binary-search","dsa","Binary Search + DP","bs-dp",h,120,7,10)
	add("binary-search","dsa","Parallel Binary Search","bs-parallel",x,180,14,11)

	// ── Two Pointers ──────────────────────────────────────────────────────────
	add("two-pointers","dsa","Opposite Ends","tp-opposite",e,60,7,1)
	add("two-pointers","dsa","Fast & Slow Pointer","tp-fast-slow",m,90,7,2)
	add("two-pointers","dsa","3Sum / kSum","tp-ksum",m,90,7,3)
	add("two-pointers","dsa","Remove Duplicates","tp-dedup",e,45,7,4)
	add("two-pointers","dsa","Container With Most Water","tp-water",m,60,7,5)
	add("two-pointers","dsa","Trapping Rain Water","tp-rain",h,90,7,6)

	// ── Sliding Window ────────────────────────────────────────────────────────
	add("sliding-window","dsa","Fixed Size Window","sw-fixed",e,60,7,1)
	add("sliding-window","dsa","Variable Size Window","sw-variable",m,90,7,2)
	add("sliding-window","dsa","Window + Hashing","sw-hash",m,90,7,3)
	add("sliding-window","dsa","Minimum Window Substring","sw-min-window",h,120,7,4)
	add("sliding-window","dsa","Longest Substring K Distinct","sw-k-distinct",m,90,7,5)

	// ── Hashing ───────────────────────────────────────────────────────────────
	add("hashing","dsa","HashMap Basics","hash-basics",e,60,7,1)
	add("hashing","dsa","Frequency Count","hash-freq",e,45,7,2)
	add("hashing","dsa","Two Sum Pattern","hash-twosum",e,45,7,3)
	add("hashing","dsa","Group Anagrams","hash-anagram",m,60,7,4)
	add("hashing","dsa","LRU Cache","hash-lru",h,120,7,5)
	add("hashing","dsa","Rolling Hash","hash-rolling",h,120,10,6)

	// ── Prefix Sum ────────────────────────────────────────────────────────────
	add("prefix-sum","dsa","1D Prefix Sum","ps-1d",e,45,7,1)
	add("prefix-sum","dsa","2D Prefix Sum","ps-2d",m,90,7,2)
	add("prefix-sum","dsa","Difference Array","ps-diff",m,90,7,3)
	add("prefix-sum","dsa","Prefix XOR","ps-xor",m,60,7,4)
	add("prefix-sum","dsa","Prefix Sum + HashMap","ps-hashmap",m,90,7,5)

	// ── Sorting ───────────────────────────────────────────────────────────────
	add("sorting","dsa","Merge Sort","sort-merge",m,90,10,1)
	add("sorting","dsa","Quick Sort","sort-quick",m,90,10,2)
	add("sorting","dsa","Counting Sort","sort-counting",e,60,14,3)
	add("sorting","dsa","Custom Comparator","sort-comparator",m,60,7,4)
	add("sorting","dsa","Sort + Greedy","sort-greedy",m,90,7,5)
	add("sorting","dsa","External Sort","sort-external",h,120,21,6)

	// ── Greedy ────────────────────────────────────────────────────────────────
	add("greedy","dsa","Activity Selection","greedy-activity",m,90,7,1)
	add("greedy","dsa","Interval Scheduling","greedy-interval",m,90,7,2)
	add("greedy","dsa","Huffman Coding","greedy-huffman",h,120,14,3)
	add("greedy","dsa","Jump Game","greedy-jump",m,60,7,4)
	add("greedy","dsa","Gas Station","greedy-gas",m,90,7,5)
	add("greedy","dsa","Task Scheduling","greedy-task",m,90,7,6)
	add("greedy","dsa","Fractional Knapsack","greedy-knapsack",m,60,14,7)

	// ── Recursion & Backtracking ──────────────────────────────────────────────
	add("recursion-backtracking","dsa","Basic Recursion","rec-basic",e,60,7,1)
	add("recursion-backtracking","dsa","Subsets","rec-subsets",m,90,7,2)
	add("recursion-backtracking","dsa","Permutations","rec-perms",m,90,7,3)
	add("recursion-backtracking","dsa","Combinations","rec-combos",m,90,7,4)
	add("recursion-backtracking","dsa","N-Queens","rec-nqueens",h,120,7,5)
	add("recursion-backtracking","dsa","Sudoku Solver","rec-sudoku",h,120,7,6)
	add("recursion-backtracking","dsa","Word Search","rec-word-search",h,120,7,7)
	add("recursion-backtracking","dsa","Pruning Techniques","rec-pruning",h,120,7,8)

	// ── Bit Manipulation ──────────────────────────────────────────────────────
	add("bit-manipulation","dsa","Basic Bit Operations","bit-basic",e,60,10,1)
	add("bit-manipulation","dsa","Power of Two","bit-pow2",e,45,10,2)
	add("bit-manipulation","dsa","XOR Tricks","bit-xor",m,90,7,3)
	add("bit-manipulation","dsa","Bitmask Enumeration","bit-mask",m,90,7,4)
	add("bit-manipulation","dsa","Count Set Bits","bit-count",m,60,10,5)
	add("bit-manipulation","dsa","Bit DP","bit-dp",h,180,10,6)

	// ── Linked List ───────────────────────────────────────────────────────────
	add("linked-list","dsa","Reversal","ll-reverse",m,60,7,1)
	add("linked-list","dsa","Cycle Detection (Floyd)","ll-cycle",m,90,7,2)
	add("linked-list","dsa","Merge Two Lists","ll-merge",e,60,7,3)
	add("linked-list","dsa","Find Middle","ll-middle",e,45,7,4)
	add("linked-list","dsa","LRU Cache (LL+Hash)","ll-lru",h,120,7,5)
	add("linked-list","dsa","Reverse K Group","ll-rev-k",h,120,7,6)
	add("linked-list","dsa","Intersection of Lists","ll-intersect",m,60,7,7)

	// ── Stack ─────────────────────────────────────────────────────────────────
	add("stack","dsa","Balanced Parentheses","stack-parens",e,45,7,1)
	add("stack","dsa","Next Greater Element","stack-nge",m,90,7,2)
	add("stack","dsa","Min Stack","stack-min",m,60,7,3)
	add("stack","dsa","Evaluate Expression","stack-eval",m,90,7,4)
	add("stack","dsa","Largest Rectangle Histogram","stack-histogram",h,120,7,5)
	add("stack","dsa","Monotonic Stack Pattern","stack-mono",m,90,7,6)

	// ── Heap / Priority Queue ─────────────────────────────────────────────────
	add("heap","dsa","Kth Largest / Smallest","heap-kth",m,90,7,1)
	add("heap","dsa","Merge K Sorted Lists","heap-merge-k",h,120,7,2)
	add("heap","dsa","Top K Frequent","heap-topk",m,90,7,3)
	add("heap","dsa","Median of Data Stream","heap-median",h,120,7,4)
	add("heap","dsa","Task Scheduler","heap-scheduler",m,90,7,5)
	add("heap","dsa","Heap Sort","heap-sort",m,90,14,6)

	// ── Queue & Deque ─────────────────────────────────────────────────────────
	add("queue-deque","dsa","Sliding Window Maximum","dq-sw-max",h,90,7,1)
	add("queue-deque","dsa","Monotonic Deque","dq-mono",m,90,7,2)
	add("queue-deque","dsa","BFS Queue Pattern","dq-bfs",m,60,7,3)
	add("queue-deque","dsa","Circular Queue","dq-circular",m,60,10,4)

	// ── Binary Tree ───────────────────────────────────────────────────────────
	add("binary-tree","dsa","Inorder / Preorder / Postorder","bt-traversal",e,60,7,1)
	add("binary-tree","dsa","Level Order BFS","bt-level",m,90,7,2)
	add("binary-tree","dsa","Height & Depth","bt-height",e,45,7,3)
	add("binary-tree","dsa","Path Sum","bt-pathsum",m,90,7,4)
	add("binary-tree","dsa","Lowest Common Ancestor","bt-lca",m,120,7,5)
	add("binary-tree","dsa","Diameter of Tree","bt-diameter",m,90,7,6)
	add("binary-tree","dsa","Serialize / Deserialize","bt-serialize",h,120,7,7)
	add("binary-tree","dsa","Binary Tree DP","bt-dp",h,180,7,8)
	add("binary-tree","dsa","Right Side View","bt-rside",m,60,7,9)
	add("binary-tree","dsa","Flatten Binary Tree","bt-flatten",m,90,7,10)

	// ── BST ───────────────────────────────────────────────────────────────────
	add("bst","dsa","Insert / Delete / Search","bst-basic",e,60,7,1)
	add("bst","dsa","Kth Smallest / Largest","bst-kth",m,60,7,2)
	add("bst","dsa","Validate BST","bst-validate",m,60,7,3)
	add("bst","dsa","BST to Sorted Array","bst-inorder",e,45,7,4)
	add("bst","dsa","Balanced BST","bst-balanced",m,90,7,5)
	add("bst","dsa","BST + LCA","bst-lca",m,90,7,6)

	// ── Trie ──────────────────────────────────────────────────────────────────
	add("trie","dsa","Basic Trie Insert/Search","trie-basic",m,120,7,1)
	add("trie","dsa","Prefix Search","trie-prefix",m,90,7,2)
	add("trie","dsa","Word Dictionary (Wildcard)","trie-wildcard",h,120,7,3)
	add("trie","dsa","Trie + DFS","trie-dfs",h,120,7,4)
	add("trie","dsa","XOR Trie","trie-xor",h,150,10,5)
	add("trie","dsa","Suffix Trie","trie-suffix",h,180,14,6)

	// ── Tree Algorithms ───────────────────────────────────────────────────────
	add("tree-algorithms","dsa","Binary Lifting","tree-binary-lifting",h,180,10,1)
	add("tree-algorithms","dsa","LCA with Binary Lifting","tree-lca-lifting",h,180,10,2)
	add("tree-algorithms","dsa","Euler Tour","tree-euler-tour",h,150,10,3)
	add("tree-algorithms","dsa","Tree DP","tree-dp",h,240,7,4)
	add("tree-algorithms","dsa","Rerooting DP","tree-rerooting",h,240,10,5)
	add("tree-algorithms","dsa","Centroid Decomposition","tree-centroid",x,300,14,6)
	add("tree-algorithms","dsa","Heavy Light Decomposition","tree-hld",x,360,14,7)
	add("tree-algorithms","dsa","DSU on Tree","tree-dsu",x,300,14,8)
	add("tree-algorithms","dsa","Tree Diameter","tree-diameter",m,120,10,9)

	// ── Graph Traversal ───────────────────────────────────────────────────────
	add("graph-traversal","dsa","Basic DFS","graph-dfs",m,90,7,1)
	add("graph-traversal","dsa","Basic BFS","graph-bfs",m,90,7,2)
	add("graph-traversal","dsa","Connected Components","graph-components",m,90,7,3)
	add("graph-traversal","dsa","Multi-Source BFS","graph-ms-bfs",m,90,7,4)
	add("graph-traversal","dsa","0-1 BFS","graph-01bfs",h,120,7,5)
	add("graph-traversal","dsa","Flood Fill","graph-flood",e,60,7,6)
	add("graph-traversal","dsa","Cycle Detection","graph-cycle",m,120,7,7)
	add("graph-traversal","dsa","Bipartite Check","graph-bipartite",m,90,7,8)
	add("graph-traversal","dsa","State BFS","graph-state-bfs",h,150,7,9)
	add("graph-traversal","dsa","Grid Graph BFS","graph-grid",m,90,7,10)

	// ── Shortest Paths ────────────────────────────────────────────────────────
	add("shortest-paths","dsa","Dijkstra","sp-dijkstra",m,120,7,1)
	add("shortest-paths","dsa","Bellman Ford","sp-bellman",m,120,10,2)
	add("shortest-paths","dsa","Floyd Warshall","sp-floyd",m,90,10,3)
	add("shortest-paths","dsa","SPFA","sp-spfa",h,120,14,4)
	add("shortest-paths","dsa","A*","sp-astar",h,180,14,5)
	add("shortest-paths","dsa","Multi-Source Dijkstra","sp-ms-dijkstra",h,120,7,6)
	add("shortest-paths","dsa","Path Reconstruction","sp-reconstruct",m,90,7,7)
	add("shortest-paths","dsa","Negative Cycle Detection","sp-neg-cycle",h,120,14,8)

	// ── Topological Sort ──────────────────────────────────────────────────────
	add("topological-sort","dsa","Kahn's Algorithm (BFS)","topo-kahn",m,120,7,1)
	add("topological-sort","dsa","DFS-based Topo Sort","topo-dfs",m,90,7,2)
	add("topological-sort","dsa","Course Schedule Pattern","topo-course",m,90,7,3)
	add("topological-sort","dsa","DAG DP","topo-dag-dp",h,150,7,4)
	add("topological-sort","dsa","Longest Path in DAG","topo-longest",h,120,10,5)
	add("topological-sort","dsa","Alien Dictionary","topo-alien",h,150,10,6)

	// ── DSU ───────────────────────────────────────────────────────────────────
	add("dsu","dsa","Path Compression","dsu-path",m,120,7,1)
	add("dsu","dsa","Union by Rank","dsu-rank",m,90,7,2)
	add("dsu","dsa","Union by Size","dsu-size",m,90,7,3)
	add("dsu","dsa","Cycle Detection with DSU","dsu-cycle",m,90,7,4)
	add("dsu","dsa","DSU + Kruskal","dsu-kruskal",m,120,7,5)
	add("dsu","dsa","Rollback DSU","dsu-rollback",x,240,14,6)
	add("dsu","dsa","Offline Queries","dsu-offline",x,240,14,7)

	// ── MST ───────────────────────────────────────────────────────────────────
	add("mst","dsa","Kruskal's Algorithm","mst-kruskal",m,120,7,1)
	add("mst","dsa","Prim's Algorithm","mst-prim",m,120,7,2)
	add("mst","dsa","Second MST","mst-second",h,180,14,3)
	add("mst","dsa","Minimax / Maximin Path","mst-minimax",h,150,14,4)

	// ── SCC & Bridges ─────────────────────────────────────────────────────────
	add("scc-bridges","dsa","Kosaraju's SCC","scc-kosaraju",h,180,14,1)
	add("scc-bridges","dsa","Tarjan's SCC","scc-tarjan",h,180,14,2)
	add("scc-bridges","dsa","Condensation Graph","scc-condensation",h,150,14,3)
	add("scc-bridges","dsa","Bridges (Tarjan)","bridges-tarjan",h,180,14,4)
	add("scc-bridges","dsa","Articulation Points","bridges-ap",h,180,14,5)

	// ── Network Flow ──────────────────────────────────────────────────────────
	add("network-flow","dsa","Max Flow (Dinic)","flow-dinic",x,300,14,1)
	add("network-flow","dsa","Max Flow (Edmonds-Karp)","flow-ek",h,240,14,2)
	add("network-flow","dsa","Min Cut","flow-mincut",h,180,14,3)
	add("network-flow","dsa","Bipartite Matching","flow-bipartite",h,180,10,4)
	add("network-flow","dsa","Min Cost Max Flow","flow-mcmf",x,360,21,5)

	// ── Graph Matching ────────────────────────────────────────────────────────
	add("graph-matching","dsa","Bipartite Matching (Hopcroft-Karp)","match-hk",h,240,14,1)
	add("graph-matching","dsa","Hungarian Algorithm","match-hungarian",x,300,21,2)
	add("graph-matching","dsa","Stable Marriage","match-stable",h,180,14,3)

	// ── Segment Tree ──────────────────────────────────────────────────────────
	add("segment-tree","dsa","Point Update Range Query","st-purq",m,180,7,1)
	add("segment-tree","dsa","Range Update Point Query","st-rupq",m,180,7,2)
	add("segment-tree","dsa","Lazy Propagation","st-lazy",h,240,7,3)
	add("segment-tree","dsa","Segment Tree + Merge Sort","st-merge-sort",h,240,10,4)
	add("segment-tree","dsa","2D Segment Tree","st-2d",x,360,14,5)
	add("segment-tree","dsa","Persistent Segment Tree","st-persistent",x,360,14,6)
	add("segment-tree","dsa","Segment Tree Beats","st-beats",x,360,21,7)
	add("segment-tree","dsa","Implicit Segment Tree","st-implicit",x,300,14,8)

	// ── Fenwick Tree ──────────────────────────────────────────────────────────
	add("fenwick-tree","dsa","1D BIT","bit-1d",m,120,7,1)
	add("fenwick-tree","dsa","2D BIT","bit-2d",h,180,10,2)
	add("fenwick-tree","dsa","BIT with Order Statistics","bit-order",h,180,10,3)
	add("fenwick-tree","dsa","BIT Range Update","bit-range",m,120,7,4)

	// ── Sparse Table ──────────────────────────────────────────────────────────
	add("sparse-table","dsa","RMQ Sparse Table","spt-rmq",m,120,10,1)
	add("sparse-table","dsa","LCA with Sparse Table","spt-lca",h,180,10,2)

	// ── Mo's Algorithm ────────────────────────────────────────────────────────
	add("mos-algorithm","dsa","Basic Mo's","mo-basic",h,240,14,1)
	add("mos-algorithm","dsa","Mo's on Trees","mo-trees",x,360,21,2)
	add("mos-algorithm","dsa","Mo's with Updates","mo-updates",x,300,21,3)

	// ── Dynamic Programming ───────────────────────────────────────────────────
	add("dp","dsa","1D DP","dp-1d",m,120,7,1)
	add("dp","dsa","2D DP","dp-2d",m,150,7,2)
	add("dp","dsa","Knapsack 0/1","dp-knapsack01",m,180,7,3)
	add("dp","dsa","Unbounded Knapsack","dp-knapsack-unbound",m,120,7,4)
	add("dp","dsa","LIS","dp-lis",m,150,7,5)
	add("dp","dsa","LCS","dp-lcs",m,150,7,6)
	add("dp","dsa","Grid DP","dp-grid",m,150,7,7)
	add("dp","dsa","Tree DP","dp-tree",h,240,7,8)
	add("dp","dsa","Bitmask DP","dp-bitmask",h,240,7,9)
	add("dp","dsa","Digit DP","dp-digit",h,240,10,10)
	add("dp","dsa","Interval DP","dp-interval",h,240,10,11)
	add("dp","dsa","Profile DP","dp-profile",x,360,14,12)
	add("dp","dsa","Rerooting DP","dp-rerooting",h,240,10,13)
	add("dp","dsa","Probability DP","dp-probability",h,240,14,14)
	add("dp","dsa","Game DP","dp-game",h,180,10,15)
	add("dp","dsa","SOS DP","dp-sos",x,300,14,16)
	add("dp","dsa","Divide & Conquer DP","dp-dc",x,300,14,17)
	add("dp","dsa","Convex Hull Trick","dp-cht",x,360,14,18)
	add("dp","dsa","Matrix Exponentiation DP","dp-matrix-exp",x,360,21,19)
	add("dp","dsa","State Compression DP","dp-state-compress",h,240,10,20)

	// ── Number Theory ─────────────────────────────────────────────────────────
	add("number-theory","mathematics","GCD / LCM","nt-gcd",e,60,14,1)
	add("number-theory","mathematics","Prime Sieve","nt-sieve",m,90,14,2)
	add("number-theory","mathematics","Prime Factorization","nt-factorize",m,90,14,3)
	add("number-theory","mathematics","Modular Arithmetic","nt-modular",m,120,7,4)
	add("number-theory","mathematics","Modular Inverse","nt-modinv",m,120,7,5)
	add("number-theory","mathematics","Euler's Totient","nt-totient",h,120,14,6)
	add("number-theory","mathematics","CRT","nt-crt",h,180,14,7)
	add("number-theory","mathematics","Miller-Rabin Primality","nt-miller",x,240,21,8)
	add("number-theory","mathematics","Pollard Rho","nt-pollard",x,300,21,9)

	// ── Combinatorics ─────────────────────────────────────────────────────────
	add("combinatorics","mathematics","nCr / nPr","comb-ncr",m,90,7,1)
	add("combinatorics","mathematics","Pascal's Triangle","comb-pascal",e,60,14,2)
	add("combinatorics","mathematics","Inclusion-Exclusion","comb-pie",h,180,10,3)
	add("combinatorics","mathematics","Stirling Numbers","comb-stirling",x,240,21,4)
	add("combinatorics","mathematics","Catalan Numbers","comb-catalan",m,120,14,5)
	add("combinatorics","mathematics","Stars & Bars","comb-stars",m,90,14,6)

	// ── Probability ───────────────────────────────────────────────────────────
	add("probability","mathematics","Basic Probability","prob-basic",m,90,14,1)
	add("probability","mathematics","Expected Value","prob-ev",m,120,10,2)
	add("probability","mathematics","Conditional Probability","prob-cond",m,120,14,3)
	add("probability","mathematics","Random Walk","prob-rw",h,180,14,4)

	// ── Advanced Strings (CP) ─────────────────────────────────────────────────
	add("advanced-strings","cp","Suffix Array","adv-str-sa",x,360,14,1)
	add("advanced-strings","cp","Suffix Automaton","adv-str-sam",x,360,14,2)
	add("advanced-strings","cp","Aho-Corasick","adv-str-ac",x,300,14,3)
	add("advanced-strings","cp","Manacher's Algorithm","adv-str-manacher",h,180,14,4)
	add("advanced-strings","cp","Palindromic Tree (Eertree)","adv-str-eertree",x,360,21,5)

	// ── FFT / NTT ─────────────────────────────────────────────────────────────
	add("fft-ntt","cp","FFT Basics","fft-basic",x,300,21,1)
	add("fft-ntt","cp","NTT","fft-ntt-impl",x,300,21,2)
	add("fft-ntt","cp","Polynomial Multiplication","fft-poly",x,300,21,3)

	// ── DP Optimizations (CP) ─────────────────────────────────────────────────
	add("dp-optimizations","cp","Convex Hull Trick","dp-opt-cht",x,360,14,1)
	add("dp-optimizations","cp","Divide & Conquer DP","dp-opt-dc",x,300,14,2)
	add("dp-optimizations","cp","Knuth's Optimization","dp-opt-knuth",x,300,21,3)
	add("dp-optimizations","cp","SMAWK Algorithm","dp-opt-smawk",x,360,21,4)

	// ── Advanced Data Structures (CP) ─────────────────────────────────────────
	add("advanced-ds","cp","Treap","ads-treap",x,360,14,1)
	add("advanced-ds","cp","Splay Tree","ads-splay",x,360,21,2)
	add("advanced-ds","cp","Link-Cut Tree","ads-lct",x,480,21,3)
	add("advanced-ds","cp","Li Chao Tree","ads-lichao",x,360,14,4)
	add("advanced-ds","cp","Wavelet Tree","ads-wavelet",x,480,21,5)
	add("advanced-ds","cp","PBDS (Policy-Based)","ads-pbds",h,180,14,6)

	// ── Constructive ──────────────────────────────────────────────────────────
	add("constructive","cp","Constructive Basics","const-basic",h,120,7,1)
	add("constructive","cp","Invariant Construction","const-invariant",h,180,10,2)
	add("constructive","cp","Greedy Construction","const-greedy",h,150,10,3)

	// ── Meet in the Middle ────────────────────────────────────────────────────
	add("meet-in-middle","cp","Subset Split","mitm-subset",h,180,14,1)
	add("meet-in-middle","cp","Baby-step Giant-step","mitm-bsgs",x,300,21,2)

	// ── OS ────────────────────────────────────────────────────────────────────
	add("operating-systems","core-cs","Processes & Threads","os-processes",m,120,7,1)
	add("operating-systems","core-cs","Scheduling Algorithms","os-scheduling",m,120,7,2)
	add("operating-systems","core-cs","Memory Management","os-memory",m,120,7,3)
	add("operating-systems","core-cs","Virtual Memory & Paging","os-paging",m,120,7,4)
	add("operating-systems","core-cs","File Systems","os-fs",m,90,10,5)
	add("operating-systems","core-cs","Deadlock","os-deadlock",m,90,7,6)
	add("operating-systems","core-cs","Synchronization (Mutex/Semaphore)","os-sync",m,120,7,7)
	add("operating-systems","core-cs","IPC","os-ipc",m,90,10,8)

	// ── DBMS ──────────────────────────────────────────────────────────────────
	add("dbms","core-cs","SQL Basics","dbms-sql",e,120,7,1)
	add("dbms","core-cs","Joins & Aggregations","dbms-joins",m,120,7,2)
	add("dbms","core-cs","Indexing","dbms-indexing",m,120,7,3)
	add("dbms","core-cs","Transactions & ACID","dbms-txn",m,120,7,4)
	add("dbms","core-cs","Normalization","dbms-norm",m,90,10,5)
	add("dbms","core-cs","Query Optimization","dbms-query-opt",h,150,10,6)
	add("dbms","core-cs","NoSQL vs SQL","dbms-nosql",m,90,14,7)

	// ── Computer Networks ─────────────────────────────────────────────────────
	add("computer-networks","core-cs","OSI Model","cn-osi",e,90,14,1)
	add("computer-networks","core-cs","TCP vs UDP","cn-tcp-udp",m,90,7,2)
	add("computer-networks","core-cs","HTTP/HTTPS","cn-http",m,90,7,3)
	add("computer-networks","core-cs","DNS","cn-dns",m,60,14,4)
	add("computer-networks","core-cs","Sockets","cn-sockets",m,120,10,5)
	add("computer-networks","core-cs","Load Balancing Basics","cn-lb",m,90,10,6)

	// ── Concurrency ───────────────────────────────────────────────────────────
	add("concurrency","core-cs","Goroutines & Channels","conc-goroutines",m,180,7,1)
	add("concurrency","core-cs","Mutex & RWMutex","conc-mutex",m,120,7,2)
	add("concurrency","core-cs","WaitGroup & Context","conc-waitgroup",m,120,7,3)
	add("concurrency","core-cs","Worker Pool","conc-worker-pool",h,180,7,4)
	add("concurrency","core-cs","Race Conditions","conc-race",m,120,7,5)
	add("concurrency","core-cs","Deadlock Prevention","conc-deadlock",m,120,10,6)

	// ── Go Language ───────────────────────────────────────────────────────────
	add("go-language","backend","Go Basics & Syntax","go-basics",e,120,14,1)
	add("go-language","backend","Structs & Interfaces","go-structs",m,120,7,2)
	add("go-language","backend","Error Handling","go-errors",m,90,7,3)
	add("go-language","backend","Goroutines","go-goroutines",m,180,7,4)
	add("go-language","backend","Channels","go-channels",m,180,7,5)
	add("go-language","backend","Context Package","go-context",m,120,7,6)
	add("go-language","backend","Generics","go-generics",m,120,14,7)
	add("go-language","backend","Testing in Go","go-testing",m,120,10,8)

	// ── PostgreSQL ────────────────────────────────────────────────────────────
	add("postgresql","backend","CRUD & DDL","pg-crud",e,120,7,1)
	add("postgresql","backend","Indexes & EXPLAIN","pg-indexes",m,180,7,2)
	add("postgresql","backend","Transactions & Isolation","pg-txn",m,180,7,3)
	add("postgresql","backend","Window Functions","pg-window",h,180,10,4)
	add("postgresql","backend","CTEs & Recursive Queries","pg-cte",h,150,10,5)
	add("postgresql","backend","JSONB","pg-jsonb",m,120,14,6)
	add("postgresql","backend","Connection Pooling (pgx)","pg-pooling",m,120,10,7)

	// ── Redis ─────────────────────────────────────────────────────────────────
	add("redis","backend","Strings & TTL","redis-strings",e,90,14,1)
	add("redis","backend","Lists & Sets","redis-collections",m,90,14,2)
	add("redis","backend","Sorted Sets","redis-sorted",m,90,10,3)
	add("redis","backend","Pub/Sub","redis-pubsub",m,120,14,4)
	add("redis","backend","Caching Patterns","redis-cache",m,120,7,5)
	add("redis","backend","Distributed Lock","redis-lock",h,180,10,6)

	// ── LLD ───────────────────────────────────────────────────────────────────
	add("lld","system-design","SOLID Principles","lld-solid",m,180,7,1)
	add("lld","system-design","Design Patterns (Creational)","lld-dp-creational",m,180,10,2)
	add("lld","system-design","Design Patterns (Structural)","lld-dp-structural",m,180,10,3)
	add("lld","system-design","Design Patterns (Behavioral)","lld-dp-behavioral",m,180,10,4)
	add("lld","system-design","Parking Lot Design","lld-parking",h,180,7,5)
	add("lld","system-design","Library Management System","lld-library",h,180,7,6)
	add("lld","system-design","Chess / Snake & Ladder","lld-chess",h,180,7,7)

	// ── HLD ───────────────────────────────────────────────────────────────────
	add("hld","system-design","URL Shortener","hld-url",m,180,7,1)
	add("hld","system-design","Rate Limiter","hld-rate-limiter",h,180,7,2)
	add("hld","system-design","Cache Design (CDN + Redis)","hld-cache",h,240,7,3)
	add("hld","system-design","Database Sharding","hld-sharding",h,240,10,4)
	add("hld","system-design","Consistent Hashing","hld-ch",h,180,7,5)
	add("hld","system-design","Message Queue Design","hld-mq",h,240,10,6)
	add("hld","system-design","Distributed File System","hld-dfs",h,240,14,7)

	return nodes
}
