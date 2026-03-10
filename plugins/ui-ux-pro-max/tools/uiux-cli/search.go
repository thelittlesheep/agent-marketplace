package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

//go:embed data/*.csv data/stacks/*.csv
var dataFS embed.FS

type SearchResult struct {
	Domain  string              `json:"domain"`
	Stack   string              `json:"stack,omitempty"`
	Query   string              `json:"query"`
	File    string              `json:"file"`
	Count   int                 `json:"count"`
	Results []map[string]string `json:"results"`
}

func Search(query string, domain string, maxResults int) (SearchResult, error) {
	if domain == "" {
		domain = DetectDomain(query)
	}

	cfg, ok := DomainConfigs[domain]
	if !ok {
		cfg = DomainConfigs["style"]
		domain = "style"
	}

	results, err := searchCSV(cfg.File, cfg.SearchCols, cfg.OutputCols, query, maxResults)
	if err != nil {
		return SearchResult{}, err
	}

	return SearchResult{
		Domain:  domain,
		Query:   query,
		File:    cfg.File,
		Count:   len(results),
		Results: results,
	}, nil
}

func SearchStack(query string, stack string, maxResults int) (SearchResult, error) {
	file, ok := StackFiles[stack]
	if !ok {
		return SearchResult{}, fmt.Errorf("unknown stack: %s", stack)
	}

	results, err := searchCSV(file, StackSearchCols, StackOutputCols, query, maxResults)
	if err != nil {
		return SearchResult{}, err
	}

	return SearchResult{
		Domain:  "stack",
		Stack:   stack,
		Query:   query,
		File:    file,
		Count:   len(results),
		Results: results,
	}, nil
}

func DetectDomain(query string) string {
	queryLower := strings.ToLower(query)

	bestDomain := "style"
	bestScore := 0

	for domain, keywords := range DomainKeywords {
		score := 0
		for _, kw := range keywords {
			if strings.Contains(queryLower, kw) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			bestDomain = domain
		}
	}

	return bestDomain
}

func searchCSV(filepath string, searchCols, outputCols []string, query string, maxResults int) ([]map[string]string, error) {
	rows, headers, err := loadCSV(filepath)
	if err != nil {
		return nil, err
	}

	// Build header index
	headerIdx := make(map[string]int)
	for i, h := range headers {
		headerIdx[strings.TrimSpace(h)] = i
	}

	// Build document strings from search columns
	documents := make([]string, len(rows))
	for i, row := range rows {
		var parts []string
		for _, col := range searchCols {
			idx, ok := headerIdx[col]
			if ok && idx < len(row) {
				parts = append(parts, row[idx])
			}
		}
		documents[i] = strings.Join(parts, " ")
	}

	// BM25 search
	bm25 := NewBM25(documents)
	scored := bm25.Score(query)

	// Take top results
	var results []map[string]string
	for _, sd := range scored {
		if len(results) >= maxResults {
			break
		}
		row := rows[sd.Index]
		entry := make(map[string]string)
		for _, col := range outputCols {
			idx, ok := headerIdx[col]
			if ok && idx < len(row) {
				entry[col] = strings.TrimSpace(row[idx])
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

func loadCSV(filepath string) ([][]string, []string, error) {
	f, err := dataFS.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", filepath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f.(io.Reader))
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable field counts

	allRows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("parse %s: %w", filepath, err)
	}

	if len(allRows) < 2 {
		return nil, nil, fmt.Errorf("%s: no data rows", filepath)
	}

	return allRows[1:], allRows[0], nil
}
