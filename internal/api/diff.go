package api

import "strings"

type DiffLine struct {
	Type string `json:"type"` // "equal" | "add" | "remove"
	Text string `json:"text"`
}

// diffLines は行単位のLCSベースdiff。プロンプト程度のサイズ（数百行）を想定した
// シンプルなO(n*m)実装で十分（大規模ソースコード向けの高速diffは不要）。
func diffLines(a, b string) []DiffLine {
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")

	n, m := len(linesA), len(linesB)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if linesA[i] == linesB[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	var result []DiffLine
	i, j := 0, 0
	for i < n && j < m {
		switch {
		case linesA[i] == linesB[j]:
			result = append(result, DiffLine{Type: "equal", Text: linesA[i]})
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			result = append(result, DiffLine{Type: "remove", Text: linesA[i]})
			i++
		default:
			result = append(result, DiffLine{Type: "add", Text: linesB[j]})
			j++
		}
	}
	for ; i < n; i++ {
		result = append(result, DiffLine{Type: "remove", Text: linesA[i]})
	}
	for ; j < m; j++ {
		result = append(result, DiffLine{Type: "add", Text: linesB[j]})
	}
	return result
}
