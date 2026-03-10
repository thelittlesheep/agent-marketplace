package main

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	bm25K1 = 1.5
	bm25B  = 0.75
)

var nonWordRe = regexp.MustCompile(`[^\p{L}\p{N}_\s]`)

type BM25 struct {
	docTokens [][]string
	docLens   []float64
	avgDL     float64
	idf       map[string]float64
	n         int
}

type ScoredDoc struct {
	Index int
	Score float64
}

func Tokenize(text string) []string {
	cleaned := nonWordRe.ReplaceAllString(strings.ToLower(text), " ")
	words := strings.Fields(cleaned)
	var tokens []string
	for _, w := range words {
		if len([]rune(w)) > 1 {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

func NewBM25(documents []string) *BM25 {
	n := len(documents)
	if n == 0 {
		return &BM25{}
	}

	docTokens := make([][]string, n)
	docLens := make([]float64, n)
	var totalLen float64

	// Tokenize all documents
	for i, doc := range documents {
		tokens := Tokenize(doc)
		docTokens[i] = tokens
		docLens[i] = float64(len(tokens))
		totalLen += docLens[i]
	}
	avgDL := totalLen / float64(n)

	// Calculate document frequency for each term
	docFreqs := make(map[string]int)
	for _, tokens := range docTokens {
		seen := make(map[string]bool)
		for _, t := range tokens {
			if !seen[t] {
				docFreqs[t]++
				seen[t] = true
			}
		}
	}

	// Calculate IDF (Robertson variant)
	idf := make(map[string]float64)
	for term, freq := range docFreqs {
		idf[term] = math.Log((float64(n)-float64(freq)+0.5)/(float64(freq)+0.5) + 1.0)
	}

	return &BM25{
		docTokens: docTokens,
		docLens:   docLens,
		avgDL:     avgDL,
		idf:       idf,
		n:         n,
	}
}

func (b *BM25) Score(query string) []ScoredDoc {
	if b.n == 0 {
		return nil
	}

	queryTokens := Tokenize(query)
	results := make([]ScoredDoc, 0, b.n)

	for i := 0; i < b.n; i++ {
		// Count term frequencies in this document
		tf := make(map[string]float64)
		for _, t := range b.docTokens[i] {
			tf[t]++
		}

		var score float64
		dl := b.docLens[i]

		for _, qt := range queryTokens {
			idfVal, ok := b.idf[qt]
			if !ok {
				continue
			}
			termFreq := tf[qt]
			num := termFreq * (bm25K1 + 1)
			denom := termFreq + bm25K1*(1-bm25B+bm25B*dl/b.avgDL)
			score += idfVal * num / denom
		}

		if score > 0 {
			results = append(results, ScoredDoc{Index: i, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}
