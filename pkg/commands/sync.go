package commands

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// Push pushes to a branch
type PushOpts struct {
	Force          bool
	UpstreamRemote string
	UpstreamBranch string
	SetUpstream    bool
}

func (c *GitCommand) PushCmdObj(opts PushOpts) (oscommands.ICmdObj, error) {
	cmdStr := "git push"

	if opts.Force {
		cmdStr += " --force-with-lease"
	}

	if opts.SetUpstream {
		cmdStr += " --set-upstream"
	}

	if opts.UpstreamRemote != "" {
		cmdStr += " " + c.OSCommand.Quote(opts.UpstreamRemote)
	}

	if opts.UpstreamBranch != "" {
		if opts.UpstreamRemote == "" {
			return nil, errors.New(c.Tr.MustSpecifyOriginError)
		}
		cmdStr += " " + c.OSCommand.Quote(opts.UpstreamBranch)
	}

	cmdObj := c.Cmd.New(cmdStr)
	return cmdObj, nil
}

func (c *GitCommand) Push(opts PushOpts, promptUserForCredential func(string) string) error {
	cmdObj, err := c.PushCmdObj(opts)
	if err != nil {
		return err
	}

	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}

type FetchOptions struct {
	PromptUserForCredential func(string) string
	RemoteName              string
	BranchName              string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	cmdStr := "git fetch"

	if opts.RemoteName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, c.OSCommand.Quote(opts.RemoteName))
	}
	if opts.BranchName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, c.OSCommand.Quote(opts.BranchName))
	}

	cmdObj := c.Cmd.New(cmdStr)
	return c.DetectUnamePass(cmdObj, func(question string) string {
		if opts.PromptUserForCredential != nil {
			return opts.PromptUserForCredential(question)
		}
		return "\n"
	})
}

type PullOptions struct {
	PromptUserForCredential func(string) string
	RemoteName              string
	BranchName              string
	FastForwardOnly         bool
}

func (c *GitCommand) Pull(opts PullOptions) error {
	if opts.PromptUserForCredential == nil {
		return errors.New("PromptUserForCredential is required")
	}

	cmdStr := "git pull --no-edit"

	if opts.FastForwardOnly {
		cmdStr += " --ff-only"
	}

	if opts.RemoteName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, c.OSCommand.Quote(opts.RemoteName))
	}
	if opts.BranchName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, c.OSCommand.Quote(opts.BranchName))
	}

	// setting GIT_SEQUENCE_EDITOR to ':' as a way of skipping it, in case the user
	// has 'pull.rebase = interactive' configured.
	cmdObj := c.Cmd.New(cmdStr).AddEnvVars("GIT_SEQUENCE_EDITOR=:")
	return c.DetectUnamePass(cmdObj, opts.PromptUserForCredential)
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(string) string) error {
	cmdStr := fmt.Sprintf("git fetch %s %s:%s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(remoteBranchName), c.OSCommand.Quote(branchName))
	cmdObj := c.Cmd.New(cmdStr)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}

func (c *GitCommand) FetchRemote(remoteName string, promptUserForCredential func(string) string) error {
	cmdStr := fmt.Sprintf("git fetch %s", c.OSCommand.Quote(remoteName))
	cmdObj := c.Cmd.New(cmdStr)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}
