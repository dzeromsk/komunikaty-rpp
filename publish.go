//go:build publish
// +build publish

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	dir    = flag.String("dir", "public/", "Path to directory to publish")
	branch = flag.String("branch", "gh-pages", "Branch to publish")
)

func main() {
	flag.Parse()

	tmpfs := memfs.New()
	filepath.Walk(*dir, func(path string, info fs.FileInfo, err error) error {
		name := strings.TrimPrefix(path, *dir)
		mode := info.Mode()
		perm := mode.Perm()
		switch {
		case mode.IsRegular():
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			return util.WriteFile(tmpfs, name, data, perm)
		case mode.IsDir():
			return tmpfs.MkdirAll(name, perm)
		}
		return nil
	})

	if err := util.WriteFile(tmpfs, ".nojekyll", []byte(""), 0644); err != nil {
		panic(err)
	}

	if err := publish(tmpfs); err != nil {
		panic(err)
	}
}

func publish(fs billy.Filesystem) error {
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return err
	}

	// set head to branch
	ref := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName(*branch))
	if err = r.Storer.SetReference(ref); err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return err
	}

	message := fmt.Sprintf("Publish: %s", os.Getenv("GITHUB_SHA"))
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "github-actions",
			Email: "github-actions@github.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "publish",
		URLs: []string{
			fmt.Sprintf("https://x-access-token:%s@github.com/%s.git",
				os.Getenv("GITHUB_TOKEN"),
				os.Getenv("GITHUB_REPOSITORY"),
			),
		},
	})

	err = r.Push(&git.PushOptions{
		RemoteName: "publish",
		Force:      true,
	})
	if err != nil {
		return err
	}

	return nil
}
