package objects

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Tree      string
	Parents   []string
	Author    string
	Committer string
	Message   string
	Timestamp time.Time
}

func NewCommit(tree, author, message string) *Commit {
	return &Commit{
		Tree:      tree,
		Parents:   make([]string, 0),
		Author:    author,
		Committer: author,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func (c *Commit) AddParent(hash string) {
	c.Parents = append(c.Parents, hash)
}

func (c *Commit) Type() ObjectType {
	return CommitObject
}

func (c *Commit) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("tree %s\n", c.Tree))

	for _, parent := range c.Parents {
		buf.WriteString(fmt.Sprintf("parent %s\n", parent))
	}

	timestamp := c.Timestamp.Unix()
	_, offset := c.Timestamp.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	timezone := fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes)

	buf.WriteString(fmt.Sprintf("author %s %d %s\n", c.Author, timestamp, timezone))
	buf.WriteString(fmt.Sprintf("committer %s %d %s\n", c.Committer, timestamp, timezone))

	buf.WriteString("\n")
	buf.WriteString(c.Message)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}

func (c *Commit) Deserialize(data []byte) error {
	lines := strings.Split(string(data), "\n")
	messageStart := -1
	for i, line := range lines {
		if line == "" {
			messageStart = i + 1
			break
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "tree":
			c.Tree = value
		case "parent":
			c.Parents = append(c.Parents, value)
		case "author":
			author, timestamp := parseAuthorLine(value)
			c.Author = author
			c.Timestamp = timestamp
		case "commiter":
			committer, _ := parseAuthorLine(value)
			c.Committer = committer
		}
	}

	if messageStart != -1 && messageStart < len(lines) {
		c.Message = strings.Join(lines[messageStart:], "\n")
		c.Message = strings.TrimSpace(c.Message)
	}

	return nil
}

func parseAuthorLine(line string) (string, time.Time) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return line, time.Now()
	}

	author := strings.Join(parts[:len(parts)-2], " ")
	timestampStr := parts[len(parts)-2]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return author, time.Now()
	}

	return author, time.Unix(timestamp, 0)
}
