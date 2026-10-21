package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brottum-brass/pdfmanager/cmd"
	"github.com/brottum-brass/pdfmanager/pkg"
)

func main() {
	if len(os.Args) < 2 {
		cmd.Usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "-h", "--help", "help":
		cmd.Usage()
		return
	case "info", "split", "merge":
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		cmd.Usage()
		os.Exit(2)
	}

	ctx := context.Background()

	manager, err := pkg.NewDefaultManager()
	if err == nil {
		switch os.Args[1] {
		case "info":
			err = cmd.RunInfo(ctx, manager, os.Args[2:])
		case "split":
			err = cmd.RunSplit(ctx, manager, os.Args[2:])
		case "merge":
			err = cmd.RunMerge(ctx, manager, os.Args[2:])
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
