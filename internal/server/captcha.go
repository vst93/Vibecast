package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

// captchaStore holds captcha answers in memory with expiry.
type captchaEntry struct {
	answer  string
	expires time.Time
}

var (
	captchaMu   sync.Mutex
	captchaStore = make(map[string]captchaEntry)
)

// generateCaptcha creates a new arithmetic captcha, stores it, and returns (id, svgDataURI).
// The SVG image renders a math problem with visual noise to prevent simple OCR.
func generateCaptcha() (string, string) {
	a := randIntRange(9) + 1 // 1-9
	b := randIntRange(9) + 1 // 1-9
	var question, answer string

	op := randIntRange(2)
	if op == 0 {
		// Addition
		question = fmt.Sprintf("%d + %d = ?", a, b)
		answer = fmt.Sprintf("%d", a+b)
	} else {
		// Subtraction (ensure non-negative)
		if a < b {
			a, b = b, a
		}
		question = fmt.Sprintf("%d − %d = ?", a, b)
		answer = fmt.Sprintf("%d", a-b)
	}

	svgDataURI := renderCaptchaSVG(question)

	id := randomToken(16)

	captchaMu.Lock()
	captchaStore[id] = captchaEntry{answer: answer, expires: time.Now().Add(5 * time.Minute)}
	// Clean expired entries
	now := time.Now()
	for k, v := range captchaStore {
		if now.After(v.expires) {
			delete(captchaStore, k)
		}
	}
	captchaMu.Unlock()

	return id, svgDataURI
}

// renderCaptchaSVG generates an SVG image with the captcha text, noise lines, and distortion.
// Returns a base64 data URI suitable for use in <img src="...">.
func renderCaptchaSVG(text string) string {
	const (
		width  = 150
		height = 40
	)

	// Remove spaces — they waste horizontal space and add nothing
	text = strings.ReplaceAll(text, " ", "")

	var b strings.Builder

	// Generate random noise lines (3-5 lines)
	numLines := 3 + randIntRange(3)
	var lines strings.Builder
	for i := 0; i < numLines; i++ {
		x1 := randIntRange(width)
		y1 := randIntRange(height)
		x2 := randIntRange(width)
		y2 := randIntRange(height)
		color := randomColor()
		opacity := 0.2 + float64(randIntRange(30))/100.0
		lines.WriteString(fmt.Sprintf(
			`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1" opacity="%.2f"/>`,
			x1, y1, x2, y2, color, opacity,
		))
	}

	// Generate noise dots (15-25)
	numDots := 15 + randIntRange(11)
	var dots strings.Builder
	for i := 0; i < numDots; i++ {
		x := randIntRange(width)
		y := randIntRange(height)
		r := 1 + randIntRange(2)
		color := randomColor()
		dots.WriteString(fmt.Sprintf(
			`<circle cx="%d" cy="%d" r="%d" fill="%s" opacity="0.3"/>`,
			x, y, r, color,
		))
	}

	// Render each character with slight rotation and vertical offset
	chars := strings.Split(text, "")
	charW := 22
	totalW := charW * len(chars)
	startX := (width - totalW) / 2 + 6
	var charSVG strings.Builder
	for i, ch := range chars {
		x := startX + i*charW
		y := 27 + randIntRange(5) - 2 // vertical jitter ±2px
		rotation := randIntRange(24) - 12 // ±12 degrees
		color := randomDarkColor()
		charSVG.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-family="'JetBrains Mono','Fira Code',monospace" font-size="20" font-weight="700" fill="%s" transform="rotate(%d %d %d)">%s</text>`,
			x, y, color, rotation, x, y, escapeXML(ch),
		))
	}

	// Build full SVG — use opaque background for cross-environment compatibility
	b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+
			`<rect width="%d" height="%d" fill="#f6f8fa"/>`+
			`%s`+ // noise dots
			`%s`+ // noise lines
			`%s`+ // characters
			`</svg>`,
		width, height, width, height,
		width, height,
		dots.String(),
		lines.String(),
		charSVG.String(),
	))

	svgBytes := []byte(b.String())
	dataURI := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(svgBytes)
	return dataURI
}

// randomColor returns a random hex color string.
func randomColor() string {
	colors := []string{"#e74c3c", "#3498db", "#2ecc71", "#f39c12", "#9b59b6", "#1abc9c", "#e67e22", "#e84393"}
	return colors[randIntRange(len(colors))]
}

// randomDarkColor returns a random dark color for text (ensures readability).
func randomDarkColor() string {
	colors := []string{"#2c3e50", "#c0392b", "#16a085", "#8e44ad", "#2980b9", "#d35400", "#27ae60"}
	return colors[randIntRange(len(colors))]
}

// escapeXML escapes special XML characters.
func escapeXML(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	).Replace(s)
}

// verifyCaptcha checks and removes the captcha. Returns true if matched.
func verifyCaptcha(id, code string) bool {
	captchaMu.Lock()
	defer captchaMu.Unlock()

	entry, ok := captchaStore[id]
	if !ok {
		return false
	}
	delete(captchaStore, id) // one-time use

	if time.Now().After(entry.expires) {
		return false
	}
	// Normalize: trim spaces
	code = strings.TrimSpace(code)
	return strings.EqualFold(entry.answer, code)
}

func randomToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func randIntRange(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}
