package main

import (
	"github.com/tidwall/redcon"
	"log"
	"red-trie/src"
	"strings"
	"sync"
)

var addr = ":6339"

var namespace = make(map[string]*src.Trie)

func main() {
	var mu sync.RWMutex
	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "ping":
				conn.WriteString("PONG")
			case "quit":
				conn.WriteString("OK")
				_ = conn.Close()
			case "add":
				// add ${namespace} ${word} (${entity})
				if len(cmd.Args) != 4 && len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				mu.Lock()
				ns := string(cmd.Args[1])
				trie, e := namespace[ns]
				if !e {
					trie = src.New()
					namespace[ns] = trie
				}
				word := string(cmd.Args[2])
				if len(cmd.Args) == 4 {
					entity := string(cmd.Args[3])
					trie.Add(word, entity)
				} else if len(cmd.Args) == 3 {
					trie.Add(word, "")
				}
				mu.Unlock()
				conn.WriteString("OK")
			case "contain":
				// contain ${namespace} ${word}
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				mu.RLock()
				ns := string(cmd.Args[1])
				word := string(cmd.Args[2])
				trie, e := namespace[ns]
				var contain = 0
				if e && trie.Contain(word) {
					contain = 1
				}
				mu.RUnlock()
				conn.WriteInt(contain)
			case "prefix":
				// prefix ${namespace} ${word}
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				mu.RLock()
				ns := string(cmd.Args[1])
				word := string(cmd.Args[2])
				trie, e := namespace[ns]
				var contain = 0
				if e && trie.Prefix(word) {
					contain = 1
				}
				mu.RUnlock()
				conn.WriteInt(contain)
			case "entities":
				// entities ${namespace} ${word}
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				mu.RLock()
				ns := string(cmd.Args[1])
				word := string(cmd.Args[2])
				trie, e := namespace[ns]
				if e {
					entities := trie.Entities(word)
					mu.RUnlock()
					if entities != nil && len(entities) != 0 {
						conn.WriteArray(len(entities) * 4)
						for _, entity := range entities {
							conn.WriteInt(0)
							conn.WriteInt(len([]rune(word)) - 1)
							conn.WriteString(entity)
							conn.WriteString(word)
						}
					} else {
						conn.WriteNull()
					}
				} else {
					mu.RUnlock()
					conn.WriteNull()
				}
			case "full":
				// full ${namespace} ${word} (${longest} default: 1)
				if len(cmd.Args) != 3 && len(cmd.Args) != 4 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				mu.RLock()
				ns := string(cmd.Args[1])
				word := string(cmd.Args[2])
				longest := true
				if len(cmd.Args) == 4 && string(cmd.Args[3]) == "0" {
					longest = false
				}
				trie, e := namespace[ns]
				if e {
					matches := trie.FullEntities(word, longest)
					mu.RUnlock()
					if matches != nil && len(matches) != 0 {
						conn.WriteArray(len(matches) * 4)
						for _, match := range matches {
							conn.WriteInt(match.Start)
							conn.WriteInt(match.End)
							conn.WriteString(match.Entity)
							conn.WriteString(match.Origin)
						}
					} else {
						conn.WriteNull()
					}
				} else {
					mu.RUnlock()
					conn.WriteNull()
				}
			}

		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
