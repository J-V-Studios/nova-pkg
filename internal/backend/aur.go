package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

// aurRPCURL is the AUR RPC endpoint, version 5.
const aurRPCURL = "https://aur.archlinux.org/rpc/?v=5"

// AUR is the backend for the Arch User Repository.
type AUR struct{}

func (AUR) Name() string { return "aur" }

// Available requires git and makepkg: AUR packages are built from source.
func (AUR) Available() bool {
	_, gitErr := exec.LookPath("git")
	_, makeErr := exec.LookPath("makepkg")
	return gitErr == nil && makeErr == nil
}

// aurResponse is the relevant subset of the RPC v5 JSON reply.
type aurResponse struct {
	Results []struct {
		Name        string
		Version     string
		Description string
	}
}

func (AUR) Search(query string) ([]Package, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	reqURL := aurRPCURL + "&type=search&arg=" + url.QueryEscape(query)
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("aur search: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aur search: unexpected status %s", resp.Status)
	}

	var ar aurResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, fmt.Errorf("aur search: %w", err)
	}
	return parseAURSearch(ar), nil
}

func (AUR) Install(pkg string) error {
	// makepkg refuses to run as root.
	if u, err := user.Current(); err == nil && u.Uid == "0" {
		return fmt.Errorf("aur install %s: refusing to run makepkg as root", pkg)
	}

	dir, err := os.MkdirTemp("", "nova-pkg-aur-")
	if err != nil {
		return fmt.Errorf("aur install %s: %w", pkg, err)
	}
	// Do not remove the temp dir: makepkg output stays for debugging.
	fmt.Println("aur: building in", dir)

	clone := exec.Command("git", "clone", "https://aur.archlinux.org/"+pkg+".git", filepath.Join(dir, pkg))
	if out, err := clone.CombinedOutput(); err != nil {
		return fmt.Errorf("aur install %s: git clone: %w\n%s", pkg, err, out)
	}

	build := exec.Command("makepkg", "-si")
	build.Dir = filepath.Join(dir, pkg)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	build.Stdin = os.Stdin
	if err := build.Run(); err != nil {
		return fmt.Errorf("aur install %s: makepkg: %w", pkg, err)
	}
	return nil
}

// parseAURSearch maps decoded RPC results to Packages. Pure function
// so it can be unit tested without network access.
func parseAURSearch(ar aurResponse) []Package {
	pkgs := make([]Package, 0, len(ar.Results))
	for _, r := range ar.Results {
		pkgs = append(pkgs, Package{
			Name:        r.Name,
			Version:     r.Version,
			Description: r.Description,
			Source:      "aur",
		})
	}
	return pkgs
}
