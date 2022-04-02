package src

type Trie struct {
	Root *TrieNode
}

type TrieNode struct {
	target   bool
	ChildMap map[rune]*TrieNode
	// 防止重复数据
	EntityList map[string]bool
}

func New() *Trie {
	trie := &Trie{}
	trie.Root = &TrieNode{false, make(map[rune]*TrieNode), make(map[string]bool)}
	return trie
}

func (t *Trie) Add(word string, entity string) {
	chars := []rune(word)
	cursor := t.Root
	for _, char := range chars {
		child, e := cursor.ChildMap[char]
		if !e {
			child = &TrieNode{false, make(map[rune]*TrieNode), make(map[string]bool)}
			cursor.ChildMap[char] = child
		}
		cursor = child
	}
	cursor.target = true
	if entity != "" {
		cursor.EntityList[entity] = true
	}
}

func (t *Trie) Contain(word string) bool {
	chars := []rune(word)
	cursor := t.Root
	for _, char := range chars {
		child, e := cursor.ChildMap[char]
		if !e {
			return false
		}
		cursor = child
	}
	return cursor.target
}

func (t *Trie) Prefix(word string) bool {
	chars := []rune(word)
	cursor := t.Root
	for _, char := range chars {
		child, e := cursor.ChildMap[char]
		if !e {
			return false
		}
		cursor = child
	}
	return true
}

func (t *Trie) Entities(word string) []string {
	chars := []rune(word)
	cursor := t.Root
	for _, char := range chars {
		child, e := cursor.ChildMap[char]
		if !e {
			return []string{}
		}
		cursor = child
	}
	if cursor.target {
		return map2strings(cursor.EntityList)
	} else {
		return []string{}
	}
}

type Match struct {
	//[start, end]
	Start  int
	End    int
	Entity string
	Origin string
}

func (t *Trie) FullEntities(word string, longest bool) []Match {
	var entites []Match
	chars := []rune(word)
	for i := 0; i < len(chars); i++ {
		cursor := t.Root
		var sub []Match
		for j := i; j < len(chars); j++ {
			child, e := cursor.ChildMap[chars[j]]
			if !e {
				break
			}
			if child.EntityList != nil {
				for k := range child.EntityList {
					sub = append(sub, Match{i, j, k, string(chars[i : j+1])})
				}
			}
			cursor = child
		}
		if sub != nil {
			if longest {
				entites = append(entites, sub[len(sub)-1])
			} else {
				for _, match := range sub {
					entites = append(entites, match)
				}
			}
		}
	}
	return entites
}

func map2strings(entities map[string]bool) (res []string) {
	if entities == nil {
		return
	}
	for k := range entities {
		res = append(res, k)
	}
	return
}
