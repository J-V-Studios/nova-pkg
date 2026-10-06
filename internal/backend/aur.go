package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// aurRPCURL is the AUR RPC endpoint, version 5. Variable (not const)
// so tests can point it at httptest servers.
var aurRPCURL = "https://aur.archlinux.org/rpc/?v=5"

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
	// Cap the body at 10MB: a hostile or broken proxy must not OOM us.
	if err := json.NewDecoder(http.MaxBytesReader(nil, resp.Body, 10<<20)).Decode(&ar); err != nil {
		return nil, fmt.Errorf("aur search: %w", err)
	}
	return parseAURSearch(ar), nil
}

func (AUR) Install(pkg string) error {
	if err := ValidateName(pkg); err != nil {
		return err
	}
	// makepkg refuses to run as root.
	if os.Geteuid() == 0 {
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

// Remove uninstalls an AUR package through pacman: makepkg-installed
// packages are ordinary pacman packages once built.
func (AUR) Remove(pkg string) error {
	return Pacman{}.Remove(pkg)
}

// Update checks foreign (AUR) packages against the RPC and reports
// newer versions. It does not rebuild automatically; the user reruns
// `nova-pkg install <name>` for the ones listed. Version comparison is
// a plain string inequality, so occasional false positives are possible.
func (AUR) Update() error {
	foreign, err := AUR{}.List()
	if err != nil {
		return fmt.Errorf("aur update: %w", err)
	}
	if len(foreign) == 0 {
		fmt.Println("aur: no foreign packages installed")
		return nil
	}

	q := url.Values{}
	for _, p := range foreign {
		q.Add("arg[]", p.Name)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(aurRPCURL + "&type=info&" + q.Encode())
	if err != nil {
		return fmt.Errorf("aur update: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("aur update: unexpected status %s", resp.Status)
	}
	var ar aurResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return fmt.Errorf("aur update: %w", err)
	}

	latest := make(map[string]string, len(ar.Results))
	for _, r := range ar.Results {
		latest[r.Name] = r.Version
	}
	outdated := outdatedAUR(foreign, latest)
	for _, o := range outdated {
		fmt.Printf("aur: %s %s -> %s (run: nova-pkg install %s)\n", o.Name, o.Version, latest[o.Name], o.Name)
	}
	if len(outdated) == 0 {
		fmt.Println("aur: all foreign packages up to date")
	}
	return nil
}

// outdatedAUR returns foreign packages whose installed version differs
// from the latest AUR version. Pure string inequality: occasional false
// positives are possible (epoch changes, packaging suffixes).
func outdatedAUR(foreign []Package, latest map[string]string) []Package {
	var out []Package
	for _, p := range foreign {
		if v, ok := latest[p.Name]; ok && v != p.Version {
			out = append(out, p)
		}
	}
	return out
}

// List returns foreign packages (`pacman -Qm`), which on this system
// means AUR-built packages.
func (AUR) List() ([]Package, error) {
	out, err := exec.Command("pacman", "-Qm").Output()
	if err != nil {
		// pacman exits 1 when no foreign packages exist.
		if len(out) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("aur list: %w", err)
	}
	return parsePacmanQ(string(out), "aur"), nil
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
