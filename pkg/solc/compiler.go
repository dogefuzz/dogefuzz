package solc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common/compiler"
)

var ErrEmptySourceFile = errors.New("solc: empty source string")
var ErrSolidityBinaryCouldNotBeDownloaded = errors.New("the solidity binary could not be downloaded externally")

type Compiler struct {
	StorageFolder string
}

func NewCompiler(storageFolder string) *Compiler {
	return &Compiler{StorageFolder: storageFolder}
}

func (c *Compiler) CompileSource(source string) (map[string]*compiler.Contract, error) {
	if len(source) == 0 {
		return nil, ErrEmptySourceFile
	}

	solcVersion, err := getIdealSolcVersionBasedOnSource(source)
	if err != nil {
		return nil, err
	}

	solcBinaryLocation, err := c.downloadSolcBinaryBasedOnVersion(solcVersion)
	if err != nil {
		return nil, err
	}

	args := append(buildArgs(solcVersion), "--")
	cmd := exec.Command(solcBinaryLocation, append(args, "-")...)
	cmd.Stdin = strings.NewReader(source)
	return run(cmd, source, solcVersion)
}

func (c *Compiler) downloadSolcBinaryBasedOnVersion(version *semver.Version) (string, error) {
	solcDestinationFolder := path.Join(c.StorageFolder, "solc")
	if err := os.MkdirAll(solcDestinationFolder, os.ModePerm); err != nil {
		return "", err
	}

	solcBinaryName := fmt.Sprintf("solcV%s", getSimplifiedVersionString(version))
	solcBinaryAbsolutePath := path.Join(solcDestinationFolder, solcBinaryName)
	solcFile, err := os.Create(solcBinaryAbsolutePath)
	if err != nil {
		return "", err
	}
	defer solcFile.Close()

	solcBinaryDownloadURL := buildSolcBinaryForLinuxURLBasedOnVersion(version)
	resp, err := http.Get(solcBinaryDownloadURL)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrSolidityBinariesListCouldNotBeDownloaded
	}

	_, err = io.Copy(solcFile, resp.Body)
	if err != nil {
		return "", err
	}

	if err := os.Chmod(solcBinaryAbsolutePath, 0777); err != nil {
		return "", err
	}

	return solcBinaryAbsolutePath, nil
}

func buildArgs(version *semver.Version) []string {
	p := []string{
		"--combined-json", "bin,bin-runtime,srcmap,srcmap-runtime,abi,userdoc,devdoc",
		"--optimize",                  // code optimizer switched on
		"--allow-paths", "., ./, ../", // default to support relative path： ./  ../  .
	}
	version0_4_6, _ := semver.NewVersion("0.4.6")
	if version.GreaterThan(version0_4_6) {
		p[1] += ",metadata,hashes"
	}
	return p
}

func run(cmd *exec.Cmd, source string, maxVersion *semver.Version) (map[string]*compiler.Contract, error) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("solc: %v\n%s", err, stderr.Bytes())
	}

	return compiler.ParseCombinedJSON(stdout.Bytes(), source, maxVersion.String(), maxVersion.String(), strings.Join(buildArgs(maxVersion), " "))
}

func getIdealSolcVersionBasedOnSource(source string) (*semver.Version, error) {
	versions, err := getDescendingOrderedVersionsFromSolidyBinariesEndpoint()
	if err != nil {
		return nil, err
	}

	versionConstraint, err := extractVersionConstraintFromSource(source)
	if err != nil {
		return nil, err
	}

	maxVersion, err := getMaxVersionBasedOnContraint(versions, versionConstraint)
	if err != nil {
		return nil, err
	}
	return maxVersion, nil
}

func buildSolcBinaryForLinuxURLBasedOnVersion(version *semver.Version) string {
	const urlFormat = "https://github.com/ethereum/solidity/releases/download/v%s/solc-static-linux"
	return fmt.Sprintf(urlFormat, getSimplifiedVersionString(version))
}

func getSimplifiedVersionString(version *semver.Version) string {
	return fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch())
}
