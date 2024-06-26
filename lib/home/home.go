package home

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

type PathExpander struct {
	homePath string
}

func NewPathExpander(homePath ...string) *PathExpander {
	if len(homePath) == 0 {
		return &PathExpander{}
	}

	return &PathExpander{homePath: homePath[0]}
}

func (pe *PathExpander) makePath(path string) string {
	return path + string(os.PathSeparator)
}

func (pe *PathExpander) SetHomePath(path string) *PathExpander {
	pe.homePath = path
	return pe
}

func (pe *PathExpander) GetHomePath() string {
	return pe.homePath
}

func (pe *PathExpander) HomeDirectory() (string, error) {
	if len(pe.GetHomePath()) > 0 {
		return pe.GetHomePath(), nil
	}

	var userData *user.User
	var err error

	if userData, err = user.Current(); err != nil {
		return "", fmt.Errorf("find homedir: %w", err)
	}

	return userData.HomeDir, nil
}

func (pe *PathExpander) Expand(path string) (string, error) {
	if !strings.HasPrefix(path, pe.makePath("~")) {
		return path, nil
	}

	var err error
	var homeString string

	// first try using the home directory resolved by PathExpander
	homeString, err = pe.HomeDirectory()
	if err != nil || len(homeString) == 0 {
		return path, fmt.Errorf("expand tilde: %w", err)
	}

	return strings.Replace(path, pe.makePath("~"), pe.makePath(homeString), 1), nil
}

// Expand takes a path as input and replaces ~ at the start of the path with the user's
// home directory. Does nothing if the path does not start with '~'.
func Expand(path string, homeDirectoryPath ...string) (string, error) {
	var pe *PathExpander

	if len(homeDirectoryPath) == 0 {
		pe = NewPathExpander()
	} else {
		pe = NewPathExpander(homeDirectoryPath[0])
	}

	return pe.Expand(path)
}
