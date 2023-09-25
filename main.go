package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/mod/modfile"
)

const MainRepoURL = "https://github.com/gnolang/gno.git"

var rootFS = flag.NewFlagSet("gno-sync-examples", flag.ExitOnError)
var (
	skipFetchFlag    = rootFS.Bool("skip-fetch", false, "Skip the fetch operation")
	mainRepoCacheDir = rootFS.String("cache-dir", "~/gno/cache/gno-main-repo", "Main repository cache directory")
)

var rootCommand = &ffcli.Command{
	ShortUsage: "gno-sync-examples [flags] <subcommand>",
	FlagSet:    rootFS,
	Subcommands: []*ffcli.Command{
		{Name: "info", Exec: showRepoInfo},
		{Name: "push", Exec: pushToMainRepo},
		{Name: "pull", Exec: pullFromMainRepo},
		{Name: "clean", Exec: cleanMainRepo},
	},
	Exec: func(ctx context.Context, args []string) error {
		return flag.ErrHelp
	},
}

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func showRepoInfo(ctx context.Context, args []string) error {
	projectDir := "."
	if len(args) > 0 {
		projectDir = args[0]
	}
	fmt.Printf("projectDir: %q\n", projectDir)

	modulePath := getModulePath(projectDir)
	fmt.Printf("modulePath: %q\n", modulePath)

	if !*skipFetchFlag {
		pullOrCloneMainRepo()
	}

	// TODO: add more info

	return nil
}

func pullOrCloneMainRepo() {
	expandedDir := expandHome(*mainRepoCacheDir)

	if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
		if err := runCommand("git", "clone", MainRepoURL, expandedDir); err != nil {
			log.Fatalf("Failed to clone repo: %s", err)
		}
	} else {
		if err := runCommand("git", "-C", expandedDir, "pull"); err != nil {
			log.Fatalf("Failed to pull repo: %s", err)
		}
	}
}

func pullFromMainRepo(ctx context.Context, args []string) error {
	expandedDir := expandHome(*mainRepoCacheDir)
	_ = expandedDir

	if !*skipFetchFlag {
		pullOrCloneMainRepo()
	}

	panic("not implemented")
	return nil
}

func pushToMainRepo(ctx context.Context, args []string) error {
	panic("not implemented")
	return nil
}

func cleanMainRepo(ctx context.Context, args []string) error {
	dir := expandHome(*mainRepoCacheDir)
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("Failed to clean main repo dir: %w.", err)
	}
	return nil
}

func expandHome(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %w.", err)
	}
	return strings.Replace(path, "~", usr.HomeDir, 1)
}

func getModulePath(workingDir string) string {
	data, err := os.ReadFile(filepath.Join(workingDir, "gno.mod"))
	if err != nil {
		log.Fatalf("Failed to read gno.mod: %w.", err)
	}

	modFile, err := modfile.Parse("gno.mod", data, nil)
	if err != nil {
		log.Fatalf("Failed to parse gno.mod: %w.", err)
	}

	return modFile.Module.Mod.Path
}

func main() {
	if err := rootCommand.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
