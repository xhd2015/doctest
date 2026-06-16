package build

import (
	"fmt"
	"io"
	"os"

	"github.com/xhd2015/doctest/libdoc/core"
)

const (
	ansiReset = "\x1b[0m"
	ansiRed   = "\x1b[31m"
	ansiGreen = "\x1b[32m"
	ansiGray  = "\x1b[90m"
)

func colorEnabled(mode core.ColorMode, w io.Writer) bool {
	switch mode {
	case core.ColorAlways:
		return true
	case core.ColorNever:
		return false
	case core.ColorAuto:
		f, ok := w.(*os.File)
		if !ok {
			return false
		}
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	default:
		return false
	}
}

type colorStyle struct {
	enabled bool
}

func newColorStyle(mode core.ColorMode, w io.Writer) colorStyle {
	return colorStyle{enabled: colorEnabled(mode, w)}
}

func (c colorStyle) red(s string) string {
	if !c.enabled {
		return s
	}
	return ansiRed + s + ansiReset
}

func (c colorStyle) green(s string) string {
	if !c.enabled {
		return s
	}
	return ansiGreen + s + ansiReset
}

func (c colorStyle) gray(s string) string {
	if !c.enabled {
		return s
	}
	return ansiGray + s + ansiReset
}

func formatSummary(style colorStyle, runCount, passCount, failCount, cachedCount int) string {
	runSeg := fmt.Sprintf("%d Run", runCount)
	passSeg := fmt.Sprintf("%d Pass", passCount)
	failSeg := fmt.Sprintf("%d Fail", failCount)
	cachedSeg := fmt.Sprintf("%d Cached", cachedCount)

	if style.enabled {
		if passCount > 0 {
			passSeg = style.green(passSeg)
		}
		if failCount > 0 {
			failSeg = style.red(failSeg)
		} else {
			failSeg = style.gray(failSeg)
		}
		cachedSeg = style.gray(cachedSeg)
	}

	return fmt.Sprintf("  (%s, %s, %s, %s)", runSeg, passSeg, failSeg, cachedSeg)
}