package hashref

type HashType int

const (
	Hash      = 0
	Text      = 1
	File      = 2
	Publisher = 3
)

var Lookup = map[HashType]string{Hash: "hash", Text: "text", File: "file", Publisher: "publisher"}
