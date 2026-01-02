package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func confirmDelete(state *state, label, identifier string, forced bool) error {
	if forced {
		return nil
	}
	if state.NoInput || !isTTY(os.Stdin) {
		return errors.New("confirmation required (use --force)")
	}
	if _, err := fmt.Fprintf(state.Err, "Delete %s '%s'? [y/N]: ", label, identifier); err != nil {
		return err
	}
	var response string
	if _, err := fmt.Fscanln(os.Stdin, &response); err != nil {
		return err
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		return errors.New("aborted")
	}
	return nil
}
