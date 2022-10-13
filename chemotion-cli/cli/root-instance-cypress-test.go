package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

const REPO_URL = "https://github.com/mehmood86/chemotion-tests.git"
const PATH_TO_REPO = "./chemotion-test"

// getPortNumber returns the port number of running active chemotion ELN instance.
func getPortNumber(service string) (port_number uint16) {
	ctx, cli := setUpDockerCleint()
	filters := filters.NewArgs()
	filters.Add("name", getInternalName(currentInstance)+"-"+service)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filters})
	panicCheck(err)
	return containers[0].Ports[0].PublicPort
}

// ChangePortNumber changes the port number in the cypress.config.js using read and write to file methodology.
//
// When a cypress test is initiated, the port number from active instance of chemotion is extracted and then updates the cypress.config.js file accordingly.
func ChangePortNumber() {
	const filepath = "./chemotion-test/cypress.config.js"
	input, err := ioutil.ReadFile(filepath)
	panicCheck(err)
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "const port") {
			lines[i] = "const port=" + strconv.Itoa(int(getPortNumber("eln")))
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	panicCheck(err)
}

func handleSetupCypress() {
	fmt.Println("initializing cypress...")
	if _, err := os.Stat(PATH_TO_REPO); os.IsNotExist(err) {
		repo, err := git.PlainClone((PATH_TO_REPO), false, &git.CloneOptions{
			URL:      REPO_URL,
			Progress: os.Stdout,
		})
		panicCheck(err)
		repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
		})
		app := "npm"
		cmd := exec.Command(app, "install", "cypress", "--save-dev")
		cmd.Dir = "./chemotion-test"
		cmd.Stdout = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("path: %s", cmd)
		cmd_err := cmd.Run()
		if cmd_err != nil {
			log.Fatalf("failed to call cmd.Run(): %v", cmd_err)
		}
	} else {
		fmt.Println("The repo already exists!.")
	}
}

func handlePullChangesCypress() {
	r, err := git.PlainOpen(PATH_TO_REPO)
	panicCheck(err)

	// Get the working directory for the repository
	w, err := r.Worktree()
	panicCheck(err)

	// Pull the latest changes from the origin remote and merge into the current branch
	fmt.Println("git pull origin")
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err.Error() == "already up-to-date" {
		fmt.Println("already up-to-date")
	} else {
		// Print the latest commit that was just pulled
		ref, err := r.Head()
		panicCheck(err)
		commit, err := r.CommitObject(ref.Hash())
		panicCheck(err)
		fmt.Println(commit)
	}
}

func handleStartCypress() {
	fmt.Println("starting cypress...")
	ChangePortNumber()
	cmd := exec.Command("npx", "cypress", "open")
	cmd.Dir = "./chemotion-test"
	cmd.Stdout = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("path: %s", cmd)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
}

var setupCypressInstanceRootCmd = &cobra.Command{
	Use:     "setup",
	Aliases: []string{"i", "setup", "install"},
	Args:    cobra.NoArgs,
	Short:   "initialize cypress framework by downloading cypress and installing required node modules.",
	Run: func(cmd *cobra.Command, args []string) {
		if ownCall(cmd) {
			handleSetupCypress()
		} else {
			handleSetupCypress()
		}
	},
}

var pullChangesCypressInstanceRootCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "update"},
	Args:    cobra.NoArgs,
	Short:   "update exiting installation of cypress, if no existing installtion, a fresh installation will be made",
	Run: func(cmd *cobra.Command, args []string) {
		if ownCall(cmd) {
			handlePullChangesCypress()
		} else {
			handlePullChangesCypress()
		}
	},
}

var startCypressInstanceRootCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s", "start"},
	Args:    cobra.NoArgs,
	Short:   "start cypress framework",
	Run: func(cmd *cobra.Command, args []string) {
		if ownCall(cmd) {
			handleStartCypress()
		} else {
			handleStartCypress()
		}
	},
}

func init() {
	cypressCmd.AddCommand(setupCypressInstanceRootCmd)
	cypressCmd.AddCommand(pullChangesCypressInstanceRootCmd)
	cypressCmd.AddCommand(startCypressInstanceRootCmd)
}
