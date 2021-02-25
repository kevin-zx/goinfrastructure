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
func (b *BrandName) PureMatch(txt string) {

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
			if len(products) != 0 || len(properties) != 0 {
				m = matchPart(products, properties, txt, match, true)
				if m {
					return true
				}
			} else {
				return true
			}
		}
	}

	for _, match := range b.cnMatch {
		if strings.Contains(txt, match) {
			if len(products) != 0 || len(properties) != 0 {
				m = matchPart(products, properties, txt, match, true)
				if m {
					return true
				}
			} else {
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
	txt = strings.ReplaceAll(txt, " ", "")
	ts := strings.Split(txt, match)
	if len(ts) == 0 {
		return false
	}
	var pe bool
	for i, t := range ts {
		trs := []rune(t)
		if isEn {
			if i == 0 {
				//第一行 不做判断，第二行才能判断英文是完整的
				if len(t) == 0 {
					pe = false
					continue
				}
				pe = isLetter(trs[len(trs)-1])
				continue
			}
			ff := isLetter(trs[0])
			// 如果当前字符串的首字母不是英文 并且上一段尾部字母不是英文
			if !ff && !pe {
				// 如果是第二行需要判断两行的距离
				if i == 1 {
					if isDistanceOK(ts[0], product, 10, true, false) && isDistanceOK(ts[1], product, 10, false, false) {
						return true
					}
				}
				// 其它行正常判断即可
				if isDistanceOK(ts[i], product, 10, false, i == len(ts)-1) {
					return true
				}
			}
		} else {
			// 如果是中文 则直接判断距离
			ok := isDistanceOK(t, product, 10, i == 0, i == len(ts)-1)
			if ok {
				return true
			}
		}

	}
	return false
}

// 匹配品牌和产品 品牌是第一优先级，如果品牌匹配不成功，产品也会默认匹配不成功
// args
// txt: 匹配文本
// brand: 匹配的品牌
// product: 产品名
// productMinDistance: 产品和品牌的最小距离
// returns
// brandMatch: 是否匹配品牌
// productMatch: 产品是否匹配上
func matchBrandAndProduct(txt string, brand string, product string, productMinDistance int, isEn bool) (brandMatch bool, productMatch bool) {
	txt = strings.ReplaceAll(txt, " ", "")
	ts := strings.Split(txt, brand)
	// log.Printf("%s\n", strings.Join(ts, "-----------"))
	if len(ts) == 0 {
		return false, false
	}
	brandMatch = false
	if !isEn {
		brandMatch = true
	}
	var pe bool
	for i, t := range ts {
		trs := []rune(t)
		if isEn {
			if i == 0 {
				//第一行 不做判断，第二行才能判断英文是完整的
				if len(t) == 0 {
					pe = false
					continue
				}
				pe = isLetter(trs[len(trs)-1])
				continue
			}
			ff := isLetter(trs[0])
			// 如果当前字符串的首字母不是英文 并且上一段尾部字母不是英文
			if !ff && !pe {
				brandMatch = true
				// 如果是第二行需要判断两行的距离
				if i == 1 {
					if isDistanceOK(ts[0], product, 10, true, false) || isDistanceOK(ts[1], product, 10, false, false) {
						return true, true
					}
				} else if isDistanceOK(ts[i], product, 10, false, i == len(ts)-1) {
					// 其它行正常判断即可
					return true, true
				}
			}
		} else {
			// if i == 1 {
			// if isDistanceOK(ts[0], product, 10, true, false) && isDistanceOK(ts[1], product, 10, false, false) {
			// return true, true
			// }
			// }
			// 如果是中文 则直接判断距离
			ok := isDistanceOK(t, product, 10, i == 0, i == len(ts)-1)
			// if !ok {
			// log.Printf("%s product:%s %v %v \n", t, product, i == 0, i == len(ts)-1)
			// }
			if ok {
				return true, true
			}
		}

	}
	return brandMatch, false
}

func isCompleteENWord(word string, txt string) bool {

	txt = strings.ReplaceAll(txt, " ", "")
	ts := strings.Split(txt, word)
	if len(ts) == 0 {
		return false
	}
	for i, t := range ts {
		// 查看 前后是不是英文 避免识别错误
		// var f rune
		// var e rune
		trs := []rune(t)
		if len(trs) == 0 {
			continue
		}
		// 首字母
		f := trs[0]
		// 尾字母
		e := trs[len(trs)-1]
		// 分割的最后一条需要判断首字母
		if i == len(ts)-1 {
			if isLetter(f) {
				return false
			}
		}
		// 分割的第一条的只要判断尾部
		if i == 0 {
			if isLetter(e) {
				return false
			}
		}
		// 分割的中间需要首位都判断
		if isLetter(f) || isLetter(e) {
			return false
		}

	}
	return true
}
func isDistanceOK(t, p string, distance int, start bool, end bool) bool {
	if p == "" {
		return true
	}
	left, right := getDistance(t, p)
	// fmt.Printf("%d %d \n", left, right)

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
	tr := []rune(t)
	pr := []rune(p)
	return runeArrayIndex(tr, pr, true), runeArrayIndex(tr, pr, false)
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
func isLetter(r rune) bool {
	return (r >= 97 && r <= 122) || (r >= 65 && r <= 90)
}

func runeCountStr(str string) int {
	return len([]rune(str))
}

func strRuneIndex(str string, sub string, right bool) int {
	return runeArrayIndex([]rune(str), []rune(sub), right)
}

// func strToRuneArray(str string) []rune {
// rs := []rune{}
// for _, r := range str {
// rs = append(rs, r)
// }
// return rs
// }

func runeArrayIndex(mainRune []rune, subRune []rune, right bool) int {
	subLen := len(subRune)
	mainLen := len(mainRune)
	if subLen > mainLen || subLen == 0 {
		return -1
	}
	if right {

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
	} else {
		for i := mainLen; i >= subLen; i-- {
			equal := true
			for j, mr := range mainRune[i-subLen : i] {
				if mr != subRune[j] {
					equal = false
					break
				}
			}
			if equal {
				return mainLen - i
			}
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
