package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
)

// StashDo modify stash
func (c *GitCommand) StashDo(index int, method string) error {
	return c.Cmd.New(fmt.Sprintf("git stash %s stash@{%d}", method, index)).Run()
}

// StashSave save stash
// TODO: before calling this, check if there is anything to save
func (c *GitCommand) StashSave(message string) error {
	return c.Cmd.New("git stash save " + c.OSCommand.Quote(message)).Run()
}

// GetStashEntryDiff stash diff
func (c *GitCommand) ShowStashEntryCmdStr(index int) string {
	return fmt.Sprintf("git stash show -p --stat --color=%s --unified=%d stash@{%d}", c.colorArg(), c.UserConfig.Git.DiffContextSize, index)
}

// StashSaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (c *GitCommand) StashSaveStagedChanges(message string) error {
	// wrap in 'writing', which uses a mutex
	if err := c.Cmd.New("git stash --keep-index").Run(); err != nil {
		return err
	}

	if err := c.StashSave(message); err != nil {
		return err
	}

	if err := c.Cmd.New("git stash apply stash@{1}").Run(); err != nil {
		return err
	}

	if err := c.OSCommand.PipeCommands("git stash show -p", "git apply -R"); err != nil {
		return err
	}

	if err := c.Cmd.New("git stash drop stash@{1}").Run(); err != nil {
		return err
	}

	// if you had staged an untracked file, that will now appear as 'AD' in git status
	// meaning it's deleted in your working tree but added in your index. Given that it's
	// now safely stashed, we need to remove it.
	files := loaders.
		NewFileLoader(c.Common, c.Cmd, c.GitConfig).
		GetStatusFiles(loaders.GetStatusFileOptions{})

	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := c.UnStageFile(file.Names(), false); err != nil {
				return err
			}
		}
	}

	return nil
}
