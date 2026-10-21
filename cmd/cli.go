package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/brottum-brass/pdfmanager/pkg"
)

func Usage() {
	fmt.Fprint(os.Stderr, `pdfmanager - offline PDF split/merge tool
 
Commands:
  info  <file.pdf>
        Print the page count.
 
  split -in <file.pdf> -out <folder-name> [mode]
        Output goes to <home>/Downloads/<folder-name> (created if missing).
        Modes (choose one):
          -ranges "1-3,5,8-"   one output file per range ("8-" = page 8 to end)
          -ranges "1-3,5" -combine
                               all listed ranges combined into ONE output file
          -every N             split into chunks of N pages
          -each                one file per page (same as -every 1)
 
  merge -out <file-name> [-divider] <a.pdf> <b.pdf> ...
        Merge files in the order given. The result is saved as
        <home>/Downloads/<file-name>.pdf (.pdf is added if missing).
        -divider inserts a blank page between each merged file.
`)
}

func RunInfo(ctx context.Context, manager *pkg.Manager, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: pdfmanager info <file.pdf>")
	}
	n, err := manager.PageCount(ctx, args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s: %d pages\n", args[0], n)
	return nil
}

func RunSplit(ctx context.Context, manager *pkg.Manager, args []string) error {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	in := fs.String("in", "", "input PDF")
	out := fs.String("out", "", "name of the output folder, created inside your Downloads folder")
	ranges := fs.String("ranges", "", `page ranges, e.g. "1-3,5,8-"`)
	every := fs.Int("every", 0, "split into chunks of N pages")
	each := fs.Bool("each", false, "one file per page")
	combine := fs.Bool("combine", false, "with -ranges: put all ranges in a single file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *in == "" || *out == "" {
		return errors.New("split requires -in and -out")
	}
	modes := 0
	if *ranges != "" {
		modes++
	}
	if *every > 0 {
		modes++
	}
	if *each {
		modes++
	}
	if modes != 1 {
		return errors.New("choose exactly one of -ranges, -every, -each")
	}
	if *combine && *ranges == "" {
		return errors.New("-combine only applies to -ranges")
	}

	var (
		written []string
		err     error
	)
	switch {
	case *each:
		written, err = manager.SplitEachPage(ctx, *in, *out)
	case *every > 0:
		written, err = manager.SplitEveryNPages(ctx, *in, *out, *every)
	default:
		written, err = manager.SplitRanges(ctx, *in, *out, *ranges, *combine)
	}
	for _, f := range written {
		fmt.Println("wrote", f)
	}
	return err
}

func RunMerge(ctx context.Context, manager *pkg.Manager, args []string) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	out := fs.String("out", "", "name of the merged PDF, saved in your Downloads folder")
	divider := fs.Bool("divider", false, "insert a blank page between merged files")
	if err := fs.Parse(args); err != nil {
		return err
	}
	inFiles := fs.Args()
	if *out == "" || len(inFiles) < 2 {
		return errors.New("usage: pdfmanager merge -out <file-name> <a.pdf> <b.pdf> ...")
	}

	outFile, err := manager.Merge(ctx, inFiles, *out, *divider)
	if err != nil {
		return err
	}
	fmt.Printf("merged %d files into %s\n", len(inFiles), outFile)
	return nil
}
