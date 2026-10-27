package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		renderHelp(c.OutOrStdout(), c)
	})
	rootCmd.SetUsageFunc(func(c *cobra.Command) error {
		renderUsage(c.OutOrStdout(), c)
		return nil
	})
}

// renderHelp draws the branded help screen. It is installed on the root
// command so every subcommand inherits it (cobra walks up the parent chain).
func renderHelp(w io.Writer, c *cobra.Command) {
	if !c.HasParent() {
		// H1: compact wordmark with the galaxy gradient + plain version line
		// (kept in a single render so the literal "autodev v<X>" is stable).
		fmt.Fprintf(w, "  %s\n", Galaxy("AUTODEV"))
		fmt.Fprintf(w, "  %s\n\n",
			styleDim.Render(fmt.Sprintf("autodev v%s — the autonomous engineering hands", c.Version)))
	}

	desc := c.Long
	if desc == "" {
		desc = c.Short
	}
	if desc != "" {
		fmt.Fprintln(w, desc)
		fmt.Fprintln(w)
	}

	writeUsage(w, c)

	if c.HasExample() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render("Examples:"))
		for _, line := range strings.Split(c.Example, "\n") {
			fmt.Fprintln(w, styleDim.Render(line))
		}
	}

	if c.HasAvailableSubCommands() {
		writeCommands(w, c)
	}

	if c.HasAvailableLocalFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render("Flags:"))
		fmt.Fprint(w, strings.TrimRight(c.LocalFlags().FlagUsages(), "\n"))
		fmt.Fprintln(w)
	}

	if c.HasAvailableInheritedFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render("Global Flags:"))
		fmt.Fprint(w, strings.TrimRight(c.InheritedFlags().FlagUsages(), "\n"))
		fmt.Fprintln(w)
	}

	if c.HasAvailableSubCommands() {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s\n", styleDim.Render(
			fmt.Sprintf("Run '%s <command> --help' for more information about a command.", c.CommandPath())))
	}

	if !c.HasParent() {
		fmt.Fprintln(w, styleDim.Render("Shell completion: autodev completion bash|zsh|fish|powershell"))
		fmt.Fprintln(w, styleDim.Render("Docs: https://github.com/HEETMEHTA18/autodev"))
	}
}

// renderUsage prints only the usage block (used for on-demand usage output).
func renderUsage(w io.Writer, c *cobra.Command) {
	writeUsage(w, c)
}

func writeUsage(w io.Writer, c *cobra.Command) {
	fmt.Fprintln(w, styleHeading.Render("Usage:"))
	if c.Runnable() {
		fmt.Fprintf(w, "  %s\n", c.UseLine())
	}
	if c.HasAvailableSubCommands() {
		fmt.Fprintf(w, "  %s [command]\n", c.CommandPath())
	}
	if len(c.Aliases) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render("Aliases:"))
		fmt.Fprintf(w, "  %s, %s\n", c.Name(), strings.Join(c.Aliases, ", "))
	}
}

func writeCommands(w io.Writer, c *cobra.Command) {
	subs := visibleCommands(c)
	if len(subs) == 0 {
		return
	}

	width := 12
	for _, sub := range subs {
		if n := len(sub.Name()) + 2; n > width {
			width = n
		}
	}

	line := func(sub *cobra.Command) {
		name := rpad(sub.Name(), width)
		short := sub.Short
		if sub.Deprecated != "" {
			short = "(deprecated) " + short
		}
		fmt.Fprintf(w, "  %s %s\n", styleBold.Render(name), short)
	}

	printed := false
	for _, group := range c.Groups() {
		var members []*cobra.Command
		for _, sub := range subs {
			if sub.GroupID == group.ID {
				members = append(members, sub)
			}
		}
		if len(members) == 0 {
			continue
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render(group.Title))
		for _, sub := range members {
			line(sub)
		}
		printed = true
	}

	var rest []*cobra.Command
	for _, sub := range subs {
		if sub.GroupID == "" {
			rest = append(rest, sub)
		}
	}
	if len(rest) > 0 {
		title := "Additional Commands:"
		if !printed {
			title = "Available Commands:"
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, styleHeading.Render(title))
		for _, sub := range rest {
			line(sub)
		}
	}
}

// visibleCommands returns the same set cobra's default template lists: all
// available children plus the auto-generated help command.
func visibleCommands(c *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, sub := range c.Commands() {
		if sub.IsAvailableCommand() || sub.Name() == "help" {
			out = append(out, sub)
		}
	}
	return out
}

func rpad(s string, padding int) string {
	if len(s) >= padding {
		return s
	}
	return s + strings.Repeat(" ", padding-len(s))
}
