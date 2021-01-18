package brandmatch

import (
	"sort"
	"strings"
	"unicode"
)

// BrandName 品牌名称信息
type BrandName struct {
	MainName string   `json:"main_name"`
	CNNames  []string `json:"cn_names"`
	ENNames  []string `json:"en_names"`
	cnMatch  []string
	enMatch  []string
	gen      bool
}

// SetMainName 设置品牌展示名称
func (b *BrandName) SetMainName() {
	b.MainName = ""
	sort.Slice(b.CNNames, func(i, j int) bool {
		return len(b.CNNames[i]) > len(b.CNNames[j])
	})
	if len(b.CNNames) > 0 {
		b.MainName += b.CNNames[0]
	}
	if len(b.ENNames) > 0 {
		if len(b.MainName) == 0 {
			b.MainName = b.ENNames[0]
		} else {
			b.MainName += " (" + b.ENNames[0] + ")"
		}

	}
}

func (b *BrandName) genMatch() {
	if b.gen {
		return
	}
	for _, name := range b.CNNames {
		b.cnMatch = append(b.cnMatch, strings.ReplaceAll(name, " ", ""))
		clearName := handlerPuntc(name)
		if clearName != "" {
			b.cnMatch = append(b.cnMatch, clearName)
		}
	}
	sort.Slice(b.cnMatch, func(i, j int) bool {
		return len(b.cnMatch[i]) > len(b.cnMatch[j])
	})
	for _, name := range b.ENNames {
		b.enMatch = append(b.enMatch, strings.ToLower(strings.ReplaceAll(name, " ", "")))
		clearName := strings.ToLower(handlerPuntc(name))
		if clearName != "" {
			b.enMatch = append(b.enMatch, clearName)
		}
	}
	sort.Slice(b.enMatch, func(i, j int) bool {
		return len(b.enMatch[i]) > len(b.enMatch[j])
	})
	b.gen = true
}

func handlerPuntc(brandName string) string {
	withouPuntc := ""
	for _, r := range brandName {
		if !unicode.IsPunct(r) && !unicode.IsSpace(r) {
			withouPuntc += string(r)
		}
	}
	if withouPuntc != brandName {
		return withouPuntc
	}
	return ""
}

// Match 匹配txt是否包含品牌，
// products 品牌过滤
// properties 属性过滤
func (b *BrandName) Match(txt string, products []string, properties []string) bool {
	txt = strings.ToLower(txt)
	b.genMatch()
	m := false

	for _, match := range b.enMatch {
		if strings.Contains(txt, match) {
			m = matchPart(products, properties, txt, match, true)
			if m {
				return true
			}
		}

	}

	for _, match := range b.cnMatch {
		if strings.Contains(txt, match) {
			m = matchPart(products, properties, txt, match, true)
			if m {
				return true
			}
		}

	}
	return false
}

func matchPart(products, shoulds []string, txt, match string, isEN bool) bool {
	m := false
	for _, p := range products {
		m = matchSingleString(txt, p, match, 10, isEN) && true
		if !m {
			break
		}
	}
	var sflags []bool
	for _, p := range shoulds {
		s := matchSingleString(txt, p, match, 10, isEN) && true
		sflags = append(sflags, s)
	}
	sm := len(shoulds) == 0
	for _, s := range sflags {
		sm = sm || s
		if sm {
			break
		}
	}

	return m && sm
}

func matchSingleString(txt string, product string, match string, distance int, isEn bool) bool {
	ts := strings.Split(txt, match)
	if len(ts) == 0 {
		return false
	}
	for i, t := range ts {

		if isEn {
			// 查看 前后是不是英文 避免识别错误
			var f rune
			var e rune
			for _, r := range t {
				if f == 0 {
					f = r
				}
				e = r
			}
			if i == len(ts)-1 {
				if unicode.IsLower(rune(f)) || unicode.IsLower(f) {
					continue
				}
			}
			if i == 0 {
				if unicode.IsLower(rune(e)) || unicode.IsLower(e) {
					continue
				}
			}

			if unicode.IsLower(rune(e)) || unicode.IsLower(e) || unicode.IsLower(rune(f)) || unicode.IsLower(f) {
				continue
			}

		}
		ok := isDistanceOK(t, product, 10, i == 0, i == len(ts)-1)
		if ok {
			return true
		}

	}
	return false
}

func isDistanceOK(t, p string, distance int, start bool, end bool) bool {
	left, right := getDistance(t, p)
	if left == -1 || right == -1 {
		return false
	}
	if start {
		if right <= distance {
			return true
			// break
		}
		return false
	}

	if end {
		if left <= distance {
			return true
			// break
		}
		return false
	}

	if left <= distance || right <= distance {
		return true
		// break
	}
	return false
}

func getDistance(t string, p string) (int, int) {
	c := runeCountStr(t)
	rc := strRuneIndex(t, p)
	if c == 0 || rc == 0 {
		return -1, -1
	}
	return rc, c - rc
}

func matchSingleRune(txt string, product rune, match string, distance int, isEn bool) bool {
	ts := strings.Split(txt, match)
	if len(ts) == 0 {
		return false
	}
	for i, t := range ts {

		c := runeCountStr(t)
		if c == 0 {
			continue
		}
		if isEn {
			// ps := strings.Split(t, "")
			// f := ps[0]
			// e := ps[len(ps)-1]
			var f rune
			var e rune
			for _, r := range t {
				if f == 0 {
					f = r
				}
				e = r
			}
			if i == len(ts)-1 {
				if unicode.IsLower(rune(f)) || unicode.IsLower(f) {
					continue
				}
			}
			if i == 0 {
				if unicode.IsLower(rune(e)) || unicode.IsLower(e) {
					continue
				}
			}

			if unicode.IsLower(rune(e)) || unicode.IsLower(e) || unicode.IsLower(rune(f)) || unicode.IsLower(f) {
				continue
			}

		}
		rc := runeStrIndex(t, product)
		if rc < 0 {
			continue
		}

		if i == 0 {
			if c-rc <= distance {
				return true
				// break
			}
			continue
		}

		if i == len(ts)-1 {
			if rc <= distance {
				return true
				// break
			}
			continue
		}

		if rc <= distance || c-rc <= distance {
			return true
			// break
		}

	}
	return false
}

func runeCountStr(str string) int {
	c := 0
	for i := range str {
		c = i
	}
	return c + 1
}

func strRuneIndex(str string, sub string) int {
	return runeArrayIndex(strToRuneArray(str), strToRuneArray(sub))
}

func strToRuneArray(str string) []rune {
	rs := []rune{}
	for _, r := range str {
		rs = append(rs, r)
	}
	return rs
}

func runeArrayIndex(mainRune []rune, subRune []rune) int {
	subLen := len(subRune)
	mainLen := len(mainRune)
	if subLen > mainLen || subLen == 0 {
		return -1
	}

	for i := range mainRune {
		if i+subLen > mainLen {
			break
		}
		equal := true
		for j, mr := range mainRune[i : i+subLen] {
			if mr != subRune[j] {
				equal = false
				break
			}
		}
		if equal {
			return i
		}
	}

	return -1
}

func runeStrIndex(str string, r rune) int {
	for i, sr := range str {
		if sr == r {
			return i
		}
	}
	return -1

}
