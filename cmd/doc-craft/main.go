package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lance52259/doc-draft/internal/ai"
	"github.com/Lance52259/doc-draft/internal/ai/provider"
	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/gitops"
	"github.com/Lance52259/doc-draft/internal/model"
	"github.com/Lance52259/doc-draft/internal/monitor"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "detect":
		os.Exit(runDetect(args))
	case "generate":
		os.Exit(runGenerate(args))
	case "run":
		os.Exit(runPipeline(args))
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `doc-craft — B examples → C docs PR (DeepSeek)

Usage:
  doc-craft detect [--out FILE] [--no-refresh]
  doc-craft generate [--practice ID] [--practices-file FILE] [--dry-run]
  doc-craft run [--practice ID] [--dry-run] [--no-refresh]

Env:
  MAX_PRACTICES   max new practices per run (0 = unlimited)
`)
}

func loadSettings() *config.Settings {
	s, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	return s
}

func runDetect(args []string) int {
	fs := flag.NewFlagSet("detect", flag.ExitOnError)
	out := fs.String("out", "", "write DetectionResult JSON")
	noRefresh := fs.Bool("no-refresh", false, "skip fetch/pull")
	_ = fs.Parse(args)

	s := loadSettings()
	if err := s.RequireRepos(); err != nil {
		log.Println(err)
		return 1
	}
	ctx, err := (&monitor.RepoWatcher{Settings: s}).PrepareRepos(!*noRefresh)
	if err != nil {
		log.Println(err)
		return 1
	}
	result, err := (&monitor.ChangeDetector{Settings: s}).Detect(ctx)
	if err != nil {
		log.Println(err)
		return 1
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	if *out != "" {
		if err := os.WriteFile(*out, append(data, '\n'), 0o644); err != nil {
			log.Println(err)
			return 1
		}
		fmt.Printf("Wrote %s (%d new)\n", *out, len(result.NewPractices))
		return 0
	}
	fmt.Println(string(data))
	return 0
}

func runGenerate(args []string) int {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	practice := fs.String("practice", "", "single practice_id")
	practicesFile := fs.String("practices-file", "", "JSON from detect")
	dryRun := fs.Bool("dry-run", true, "do not push/create PR")
	_ = fs.Parse(args)

	s := loadSettings()
	s.DryRun = *dryRun
	if err := s.RequireRepos(); err != nil {
		log.Println(err)
		return 1
	}
	if err := s.RequireAI(); err != nil {
		log.Println(err)
		return 1
	}

	repoCtx, err := (&monitor.RepoWatcher{Settings: s}).PrepareRepos(true)
	if err != nil {
		log.Println(err)
		return 1
	}
	detection, err := (&monitor.ChangeDetector{Settings: s}).Detect(repoCtx)
	if err != nil {
		log.Println(err)
		return 1
	}

	selected, err := selectPractices(s, repoCtx, detection, *practice, *practicesFile)
	if err != nil {
		log.Println(err)
		return 1
	}
	if len(selected) == 0 {
		fmt.Println("No practices to generate.")
		return 0
	}
	selected = limitPractices(selected, s.MaxPractices)

	p := provider.NewDeepSeek(s.AIAPIKey, s.AIBaseURL, s.AIModel, s.AITimeoutSeconds, s.AIMaxRetries)
	gen := ai.NewDocGenerator(s, p)
	for _, item := range selected {
		dir := filepath.Join(repoCtx.B.LocalPath, item.SourcePath)
		result, err := gen.Generate(context.Background(), item, dir)
		if err != nil {
			log.Printf("generate %s: %v", item.PracticeID, err)
			return 1
		}
		fmt.Printf("Generated %s: %d file(s)\n", item.PracticeID, len(result.Files))
		for _, f := range result.Files {
			fmt.Printf("  - %s %s\n", f.Action, f.Path)
		}
	}
	return 0
}

func runPipeline(args []string) int {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	practice := fs.String("practice", "", "only this practice_id")
	dryRunFlag := fs.String("dry-run", "", "true|false override")
	noRefresh := fs.Bool("no-refresh", false, "skip fetch/pull")
	_ = fs.Parse(args)

	s := loadSettings()
	if *dryRunFlag != "" {
		s.DryRun = strings.EqualFold(*dryRunFlag, "true") || *dryRunFlag == "1"
	}
	if err := s.RequireRepos(); err != nil {
		log.Println(err)
		return 1
	}

	repoCtx, err := (&monitor.RepoWatcher{Settings: s}).PrepareRepos(!*noRefresh)
	if err != nil {
		log.Println(err)
		return 1
	}
	detection, err := (&monitor.ChangeDetector{Settings: s}).Detect(repoCtx)
	if err != nil {
		log.Println(err)
		return 1
	}

	selected := detection.NewPractices
	if *practice != "" {
		var filtered []model.Practice
		for _, p := range selected {
			if p.PracticeID == *practice {
				filtered = append(filtered, p)
			}
		}
		selected = filtered
		if len(selected) == 0 {
			fmt.Printf("No new practice matching %s\n", *practice)
			return 0
		}
	}

	pipeline := model.PipelineResult{Detected: *detection, DryRun: s.DryRun}
	if len(selected) == 0 {
		fmt.Println("No new practices; exiting.")
		printJSON(pipeline)
		return 0
	}
	selected = limitPractices(selected, s.MaxPractices)
	if err := s.RequireAI(); err != nil {
		log.Println(err)
		return 1
	}

	p := provider.NewDeepSeek(s.AIAPIKey, s.AIBaseURL, s.AIModel, s.AITimeoutSeconds, s.AIMaxRetries)
	gen := ai.NewDocGenerator(s, p)
	op := &gitops.RepoOperator{Settings: s}
	prm, err := gitops.NewPRManager(s)
	if err != nil {
		log.Println(err)
		return 1
	}
	stateMgr := &monitor.StateManager{Path: s.AbsoluteStatePath()}
	state, _ := stateMgr.Load()

	for _, item := range selected {
		dir := filepath.Join(repoCtx.B.LocalPath, item.SourcePath)
		result, err := gen.Generate(context.Background(), item, dir)
		if err != nil {
			pipeline.Errors = append(pipeline.Errors, fmt.Sprintf("%s: %v", item.PracticeID, err))
			continue
		}
		pipeline.Generated = append(pipeline.Generated, *result)

		branch := branchName(item.PracticeID)
		title := fmt.Sprintf("docs: sync best practice %s", item.PracticeID)
		body := gitops.SummarizeChanges(result, s.BRepo, detection.BCommit)

		if _, err := op.ApplyAndPush(repoCtx.C.LocalPath, branch, s.CDefaultBranch, result, title, s.DryRun); err != nil {
			pipeline.Errors = append(pipeline.Errors, fmt.Sprintf("%s push: %v", item.PracticeID, err))
			continue
		}
		pr, err := prm.CreatePR(title, body, branch, s.CDefaultBranch, s.DryRun)
		if err != nil {
			pipeline.Errors = append(pipeline.Errors, fmt.Sprintf("%s pr: %v", item.PracticeID, err))
			continue
		}
		if pr != nil && pr.URL != "" {
			pipeline.PRURLs = append(pipeline.PRURLs, pr.URL)
			state.OpenPRs[item.PracticeID] = pr.URL
		}
		if !contains(state.ProcessedPractices, item.PracticeID) {
			state.ProcessedPractices = append(state.ProcessedPractices, item.PracticeID)
		}
	}

	state.BCommit = detection.BCommit
	state.CCommit = detection.CCommit
	if !s.DryRun {
		_ = stateMgr.Save(state)
	}
	printJSON(pipeline)
	if len(pipeline.Errors) > 0 {
		return 1
	}
	return 0
}

func limitPractices(practices []model.Practice, max int) []model.Practice {
	if max <= 0 || len(practices) <= max {
		return practices
	}
	fmt.Printf("Limiting practices: %d → %d (MAX_PRACTICES)\n", len(practices), max)
	return practices[:max]
}

func selectPractices(s *config.Settings, ctx *monitor.RepoContext, detection *model.DetectionResult, practice, practicesFile string) ([]model.Practice, error) {
	if practice != "" {
		for _, p := range detection.NewPractices {
			if p.PracticeID == practice {
				return []model.Practice{p}, nil
			}
		}
		all, err := monitor.EnumeratePractices(
			filepath.Join(ctx.B.LocalPath, s.BExamplesPath),
			s.BExamplesPath,
			s.IgnoreNames,
			s.Granularity,
		)
		if err != nil {
			return nil, err
		}
		for _, p := range all {
			if p.PracticeID == practice {
				return []model.Practice{p}, nil
			}
		}
		return nil, fmt.Errorf("practice not found: %s", practice)
	}
	if practicesFile != "" {
		data, err := os.ReadFile(practicesFile)
		if err != nil {
			return nil, err
		}
		var parsed model.DetectionResult
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		want := map[string]struct{}{}
		for _, p := range parsed.NewPractices {
			want[p.PracticeID] = struct{}{}
		}
		var selected []model.Practice
		for _, p := range detection.NewPractices {
			if _, ok := want[p.PracticeID]; ok {
				selected = append(selected, p)
			}
		}
		return selected, nil
	}
	return detection.NewPractices, nil
}

func branchName(practiceID string) string {
	slug := strings.ReplaceAll(practiceID, "/", "-")
	name := "doc-craft/" + slug
	if len(name) > 200 {
		return name[:200]
	}
	return name
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func printJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
