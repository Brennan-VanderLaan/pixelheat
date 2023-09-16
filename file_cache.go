package main

import (
	"time"

	"github.com/go-git/go-git/v5"
)

type FileNode struct {
	Name   string
	Status string
	Active bool
}

// FileStatusCache represents the cached file status,
// including the time of the last check.
type FileStatusCache struct {
	Status    string
	LastCheck time.Time
}

type GitRepoCache struct {
	Repo      *git.Repository
	Status    *git.Status
	LastCheck time.Time
}

const CacheDuration = 5 * time.Second
