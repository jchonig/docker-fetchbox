package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type logger struct {
	verbose bool
	debug   bool
}

func (l *logger) infof(format string, args ...any) {
	if l.verbose || l.debug {
		log.Printf(format, args...)
	}
}

func (l *logger) debugf(format string, args ...any) {
	if l.debug {
		log.Printf("[debug] "+format, args...)
	}
}

func main() {
	var (
		daemon       = flag.Bool("daemon", false, "run continuously at configured interval")
		listFoldersF = flag.Bool("list-folders", false, "list IMAP folders and exit")
		verbose      = flag.Bool("v", false, "verbose logging")
		debug        = flag.Bool("d", false, "debug logging")
		noop         = flag.Bool("n", false, "dry run — fetch but do not save or mark seen")
		configPath   = flag.String("config", "/config/fetchbox.yml", "path to config file")
	)
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetchbox: load config: %v\n", err)
		os.Exit(1)
	}

	if *listFoldersF {
		for _, mb := range cfg.Mailboxes {
			fmt.Printf("=== %s ===\n", mb.Name)
			if err := listFolders(mb); err != nil {
				log.Printf("list folders %s: %v", mb.Name, err)
			}
		}
		return
	}

	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetchbox: invalid interval %q: %v\n", cfg.Interval, err)
		os.Exit(1)
	}

	p := &processor{
		cfg:    cfg,
		noop:   *noop,
		logger: &logger{verbose: *verbose, debug: *debug},
	}

	p.run()
	for *daemon {
		time.Sleep(interval)
		p.run()
	}
}
