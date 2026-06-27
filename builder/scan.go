package builder

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var roleKeys = []string{
	"base", "surface", "overlay", "muted", "subtle",
	"text", "love", "gold", "rose", "pine", "foam",
	"iris", "highlightLow", "highlightMed", "highlightHigh",
	"accent", "onaccent",
}

var metaKeys = []string{
	"id", "name", "appearance", "type", "description", "accentname",
}

type ScannerOpts struct {
	Prefix rune
}

type Capture interface {
	captureSpan() span
}

type span struct {
	Start  uint
	Length uint
}

type variantArm struct {
	content  string
	captures []Capture
}

type (
	RoleCapture struct {
		span  span
		role  string
		alpha *float64
	}
	MetaCapture    struct{ span }
	VariantCapture struct {
		span             span
		main, moon, dawn variantArm
	}
	TextCapture struct{ span }
)

func (c RoleCapture) captureSpan() span    { return c.span }
func (c MetaCapture) captureSpan() span    { return c.span }
func (c VariantCapture) captureSpan() span { return c.span }
func (c TextCapture) captureSpan() span    { return c.span }

func Span(c Capture) (start, length uint) {
	s := c.captureSpan()
	return s.Start, s.Length
}

func isAsciiDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

type Scanner struct {
	Content     string
	Opts        ScannerOpts
	pos         uint
}

func (s *Scanner) curr() (byte, bool) {
	if int(s.pos) < len(s.Content) {
		return s.Content[s.pos], true
	}
	return 0, false
}

func (s *Scanner) peek() (byte, bool) {
	idx := int(s.pos + 1)
	if idx < len(s.Content) {
		return s.Content[idx], true
	}
	return 0, false
}

func (s *Scanner) remaining() string {
	return s.Content[s.pos:]
}

func (s *Scanner) advanceN(n int) {
	if int(s.pos)+n > len(s.Content) {
		return
	}
	s.pos += uint(n)
}

func (s *Scanner) advance() {
	if int(s.pos) >= len(s.Content) {
		return
	}
	s.pos++
}

func (s *Scanner) advanceToEOF() {
	s.pos = uint(len(s.Content))
}

func (s *Scanner) scanUntil(ch byte) bool {
	idx := strings.IndexByte(s.remaining(), ch)
	if idx != -1 {
		s.advanceN(idx)
		return true
	}
	return false
}

func (s *Scanner) advanceToPrefix() bool {
	return s.scanUntil(byte(s.Opts.Prefix))
}

func (s *Scanner) scanKey(keys []string) (string, bool) {
	src := s.remaining()
	for _, key := range keys {
		if strings.HasPrefix(src, key) {
			next := src[len(key):]
			if len(next) > 0 && (unicode.IsLetter(rune(next[0])) || unicode.IsDigit(rune(next[0]))) {
				continue
			}
			s.advanceN(len(key))
			return key, true
		}
	}
	return "", false
}

func (s *Scanner) scanVariantCapture(outerStartPos uint) (VariantCapture, error) {
	s.advance()
	start := s.pos
	s.scanUntil(')')
	content := s.Content[start:s.pos]
	parts := strings.SplitN(content, "|", 3)
	s.advance()

	mainCaptures, err := Scan(parts[0], s.Opts)
	if err != nil {
		return VariantCapture{}, err
	}

	moonCaptures, err := Scan(parts[1], s.Opts)
	if err != nil {
		return VariantCapture{}, err
	}

	dawnCaptures, err := Scan(parts[2], s.Opts)
	if err != nil {
		return VariantCapture{}, err
	}

	return VariantCapture{
		span: span{outerStartPos, s.pos - outerStartPos},
		main: variantArm{content: parts[0], captures: mainCaptures},
		moon: variantArm{content: parts[1], captures: moonCaptures},
		dawn: variantArm{content: parts[2], captures: dawnCaptures},
	}, nil
}

func (s *Scanner) scanRoleCapture(start uint, role string) (RoleCapture, error) {
	var alpha *float64
	if curr, _ := s.curr(); curr == '/' {
		if peeked, _ := s.peek(); isAsciiDigit(peeked) {
			s.advance()
			alphaStart := s.pos
			for {
				if curr, _ := s.curr(); isAsciiDigit(curr) {
					s.advance()
				} else {
					break
				}
			}
			parsed, err := strconv.ParseInt(s.Content[alphaStart:s.pos], 10, 32)
			if err == nil {
				value := float64(parsed) / 100
				alpha = &value
			}
		}
	}
	return RoleCapture{
		span:  span{start, s.pos - start},
		role:  role,
		alpha: alpha,
	}, nil
}

func (s *Scanner) scanCapture() (Capture, bool, error) {
	start := s.pos

	if curr, _ := s.curr(); curr == byte(s.Opts.Prefix) {
		s.advance()

		if curr, _ := s.curr(); curr == '(' {
			vc, err := s.scanVariantCapture(start)
			return vc, true, err
		}

		if _, found := s.scanKey(metaKeys); found {
			return MetaCapture{span{start, s.pos - start}}, true, nil
		}

		if role, found := s.scanKey(roleKeys); found {
			rc, err := s.scanRoleCapture(start, role)
			return rc, true, err
		}

		unknown := strings.IndexFunc(s.remaining(), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		var key string
		if unknown == -1 {
			key = s.remaining()
		} else {
			key = s.remaining()[:unknown]
		}
		return nil, false, fmt.Errorf("unknown variable %c%s", s.Opts.Prefix, key)

	} else {
		if !s.advanceToPrefix() {
			s.advanceToEOF()
		}
		if length := s.pos - start; length > 0 {
			return TextCapture{span{start, length}}, true, nil
		}
	}

	return nil, false, nil
}

func (s *Scanner) scanCaptures() ([]Capture, error) {
	prefixCount := strings.Count(s.Content, string(s.Opts.Prefix))
	captures := make([]Capture, 0, prefixCount)
	for {
		capture, found, err := s.scanCapture()
		if err != nil {
			return nil, err
		}
		if !found {
			break
		}
		captures = append(captures, capture)
	}
	return captures, nil
}

func Scan(content string, opts ScannerOpts) ([]Capture, error) {
	scanner := Scanner{
		Content: content,
		Opts:    opts,
	}
	return scanner.scanCaptures()
}
