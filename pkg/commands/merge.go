package commands

import "fmt"

type MergeOpts struct {
	FastForwardOnly bool
}

// Merge merge
func (c *GitCommand) Merge(branchName string, opts MergeOpts) error {
	mergeArg := ""
	if c.UserConfig.Git.Merging.Args != "" {
		mergeArg = " " + c.UserConfig.Git.Merging.Args
	}

	command := fmt.Sprintf("git merge --no-edit%s %s", mergeArg, c.OSCommand.Quote(branchName))
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return c.OSCommand.Cmd.New(command).Run()
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.Cmd.New("git merge --abort").Run()
}
