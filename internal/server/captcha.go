package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

// captchaStore holds captcha codes in memory with expiry.
type captchaEntry struct {
	code    string
	expires time.Time
}

var (
	captchaMu    sync.Mutex
	captchaStore  = make(map[string]captchaEntry)
	captchaRunes  = []rune("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789")
	captchaColors = []string{"#6366f1", "#8b5cf6", "#ec4899", "#f59e0b", "#10b981", "#ef4444", "#3b82f6"}
)

// generateCaptcha creates a new captcha, stores it, and returns (id, svgImage).
// The SVG is a data URI string ready to use as <img src>.
func generateCaptcha() (string, string) {
	code := randomCaptchaCode(4)
	id := randomToken(16)

	captchaMu.Lock()
	captchaStore[id] = captchaEntry{code: code, expires: time.Now().Add(5 * time.Minute)}
	// Clean expired entries
	now := time.Now()
	for k, v := range captchaStore {
		if now.After(v.expires) {
			delete(captchaStore, k)
		}
	}
	captchaMu.Unlock()

	svg := renderCaptchaSVG(code)
	return id, svg
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
	return strings.EqualFold(entry.code, code)
}

func randomCaptchaCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(captchaRunes))))
		b[i] = captchaRunes[idx.Int64()]
	}
	return string(b)
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

// renderCaptchaSVG generates an SVG captcha image as a data URI.
func renderCaptchaSVG(code string) string {
	var b strings.Builder
	b.WriteString("data:image/svg+xml;utf8,")
	b.WriteString("<svg xmlns='http://www.w3.org/2000/svg' width='150' height='50'>")

	// Background
	b.WriteString("<rect width='150' height='50' fill='#f1f5f9' rx='8'/>")

	// Noise lines
	for i := 0; i < 6; i++ {
		x1 := randIntRange(150)
		y1 := randIntRange(50)
		x2 := randIntRange(150)
		y2 := randIntRange(50)
		color := captchaColors[randIntRange(len(captchaColors))]
		b.WriteString(fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' stroke='%s' stroke-width='1' opacity='0.3'/>", x1, y1, x2, y2, color))
	}

	// Noise dots
	for i := 0; i < 15; i++ {
		x := randIntRange(150)
		y := randIntRange(50)
		b.WriteString(fmt.Sprintf("<circle cx='%d' cy='%d' r='1' fill='#94a3b8' opacity='0.5'/>", x, y))
	}

	// Characters
	charWidth := 30
	startX := 15
	for i, ch := range code {
		x := startX + i*charWidth
		y := 30 + randIntRange(10) - 5
		rotation := randIntRange(30) - 15
		color := captchaColors[randIntRange(len(captchaColors))]
		fontSize := 22 + randIntRange(6)
		b.WriteString(fmt.Sprintf("<text x='%d' y='%d' font-size='%d' font-family='monospace' fill='%s' transform='rotate(%d %d %d)' font-weight='bold'>%c</text>",
			x, y, fontSize, color, rotation, x, y, ch))
	}

	b.WriteString("</svg>")
	return b.String()
}
