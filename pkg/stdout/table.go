package stdout

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"
)

type TablePrinter struct {
	title   string
	headers []string
	rows    [][]string
	added   bool
}

func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}
func (tp *TablePrinter) SetTitle(title string) {
	tp.title = title
}

// Add 添加一行数据：第一行为表头，其余为数据行
func (tp *TablePrinter) Add(row []string) {
	if !tp.added {
		tp.headers = row
		tp.added = true
	} else {
		if len(row) != len(tp.headers) {
			return
		}
		tp.rows = append(tp.rows, row)
	}
}
func (tp *TablePrinter) Reset() {
	tp.headers = nil
	tp.rows = nil
	tp.added = false
	tp.title = ""
}

// Print 打印表格
func (tp *TablePrinter) Print() error {
	if len(tp.headers) == 0 {
		return fmt.Errorf("no headers")
	}
	if tp.title != "" {
		fmt.Println(tp.title)
	}
	widths := tp.calcColumnWidths()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	border := tp.buildBorderLine(widths)
	fmt.Fprintln(w, border)
	// 打印表头
	var headerCells []string
	for i, h := range tp.headers {
		headerCells = append(headerCells, "| "+tp.padToWidth(h, widths[i])+" ")
	}
	fmt.Fprintln(w, strings.Join(headerCells, "")+"|")
	fmt.Fprintln(w, border)

	// 打印每行数据
	for _, row := range tp.rows {
		var cells []string
		for i, val := range row {
			cells = append(cells, "| "+tp.padToWidth(val, widths[i])+" ")
		}
		fmt.Fprintln(w, strings.Join(cells, "")+"|")
	}
	fmt.Fprintln(w, border)
	w.Flush()
	return nil
}

// ==================== 私有方法 ====================

func (tp *TablePrinter) isWideRune(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.In(r, unicode.Hangul, unicode.Hiragana, unicode.Katakana) ||
		(r >= 0xFF01 && r <= 0xFF60)
}

func (tp *TablePrinter) displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if tp.isWideRune(r) {
			width += 2
		} else {
			width++
		}
	}
	return width
}

func (tp *TablePrinter) padToWidth(s string, target int) string {
	cur := tp.displayWidth(s)
	if cur >= target {
		return s
	}
	return s + strings.Repeat(" ", target-cur)
}

func (tp *TablePrinter) buildBorderLine(widths []int) string {
	var b strings.Builder
	b.WriteString("+")
	for _, w := range widths {
		b.WriteString(strings.Repeat("-", w+2))
		b.WriteString("+")
	}
	return b.String()
}

func (tp *TablePrinter) calcColumnWidths() []int {
	widths := make([]int, len(tp.headers))
	for i, h := range tp.headers {
		widths[i] = tp.displayWidth(h)
	}
	for _, row := range tp.rows {
		for i, val := range row {
			w := tp.displayWidth(val)
			if w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}
