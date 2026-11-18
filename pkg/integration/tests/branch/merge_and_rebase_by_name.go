package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MergeAndRebaseByName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Merge and rebase onto a branch picked by name, while that branch is checked out in another worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("feature").
			EmptyCommit("feature commit").
			Checkout("master").
			AddWorktree("master", "../linked-worktree", "other").
			AddFileInWorktreeOrSubmodule("../linked-worktree", "other-file", "content").
			CommitInWorktreeOrSubmodule("../linked-worktree", "other commit").
			EmptyCommit("master commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Branches.CheckoutBranchByName)

		t.ExpectPopup().Prompt().Title(Equals("Branch name:")).Type("other").Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("other")).
			Select(Contains("Merge")).
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Regular merge (with merge commit)")).
			Confirm()

		t.Git().CurrentBranchName("master")
		t.Views().Commits().
			TopLines(
				Contains("Merge branch 'other'"),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature")).
			PressPrimaryAction()

		t.Views().Branches().
			Press(keys.Branches.CheckoutBranchByName)

		t.ExpectPopup().Prompt().Title(Equals("Branch name:")).Type("other").Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("other")).
			Select(Contains("Rebase")).
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'feature'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Git().CurrentBranchName("feature")
		t.Views().Commits().
			Lines(
				Contains("feature commit"),
				Contains("other commit"),
				Contains("base"),
			)
	},
})
