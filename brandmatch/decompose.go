package brandmatch

import (
	"strings"
	"unicode"

	"github.com/kevin-zx/goinfrastructure/stringer"
)

// DecomposeBrands 将品牌字符串分解为品牌
func DecomposeBrands(rawBrands []string) ([]BrandName, error) {
	fBrands, err := handRawBrands(rawBrands)
	if err != nil {
		return nil, err
	}
	fBrands = combination(fBrands)
	for i := range fBrands {
		fBrands[i].SetMainName()
	}
	return fBrands, nil
}

func handRawBrands(rawBrands []string) ([]BrandName, error) {
	bs := []BrandName{}
	// remove duplicate
	rawBrands = stringer.RemoveDuplicateAndBlankString(rawBrands)
	// split RawBrand
	for _, rb := range rawBrands {
		parts := split(rb)
		b := BrandName{}
		for _, p := range parts {
			e, c := getName(p)
			if len(e) != 0 {
				b.ENNames = append(b.ENNames, e)
			}
			if len(c) != 0 {
				b.CNNames = append(b.CNNames, c)
			}

		}
		b.SetMainName()
		bs = append(bs, b)
	}
	return bs, nil
}

func split(rb string) []string {
	rbs := []string{}
	if strings.Contains(rb, "/") {
		rbs = strings.Split(rb, "/")
	} else if strings.Contains(rb, "(") {
		rbt := strings.ReplaceAll(rb, ")", "")
		rbs = strings.Split(rbt, "(")
	} else if strings.Contains(rb, "（") {

		rbt := strings.ReplaceAll(rb, "）", "")
		rbs = strings.Split(rbt, "（")
	}
	if len(rbs) > 0 {
		if len(rbs) == 2 {
			isHan1 := isPureHan(rbs[0])
			isHan2 := isPureHan(rbs[1])
			if (isHan1 && isHan2) || (!isHan1 && !isHan2) {
				if strings.Contains(rb, "（") {
					return []string{strings.Split(rb, "（")[0]}
				}
				if strings.Contains(rb, "(") {
					return []string{strings.Split(rb, "(")[0]}
				}
				return []string{rb}
			}

		}
		return rbs
	}

	return []string{rb}
}

func combination(bs []BrandName) []BrandName {
	nBrands := []BrandName{}
	for i := len(bs) - 1; i >= 0; i-- {
		if bs[i].MainName == "" {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if needBeCombination(bs[i], bs[j]) {
				bs[i] = combinationTwo(bs[i], bs[j])
				bs[j] = BrandName{}
			}
		}
	}
	for _, b := range bs {
		if b.MainName != "" {
			nBrands = append(nBrands, b)
		}
	}
	return nBrands
}

func needBeCombination(b1 BrandName, b2 BrandName) bool {
	if b1.MainName == b2.MainName {
		return true
	}
	for i1 := range b1.CNNames {
		for i2 := range b2.CNNames {
			if strings.ReplaceAll(b1.CNNames[i1], " ", "") == strings.ReplaceAll(b2.CNNames[i2], " ", "") {
				return true
			}
		}
	}

	for i1 := range b1.ENNames {
		for i2 := range b2.ENNames {
			if minimizeENName(b1.ENNames[i1]) == minimizeENName(b2.ENNames[i2]) {
				return true
			}
		}
	}
	return false
}

func minimizeENName(name string) string {
	rn := ""
	for _, n := range name {
		if unicode.IsLetter(n) {
			rn += string(n)
		}
	}

	return strings.ToLower(rn)
}

func combinationTwo(b1 BrandName, b2 BrandName) BrandName {
	b := BrandName{}

	b.MainName = getMainNameFromTwo(b1, b2)
	b.CNNames = stringer.RemoveDuplicateAndBlankString(append(b1.CNNames, b2.CNNames...))
	b.ENNames = stringer.RemoveDuplicateAndBlankString(append(b1.ENNames, b2.ENNames...))
	return b
}
func getMainNameFromTwo(b1 BrandName, b2 BrandName) string {
	if b1.MainName == b2.MainName {
		return b1.MainName
	}
	b1MainNameIsCN := false

	for _, cn := range b1.CNNames {
		if b1.MainName == cn {
			b1MainNameIsCN = true
			break
		}
	}
	b2MainNameIsCN := false

	for _, cn := range b2.CNNames {
		if b2.MainName == cn {
			b2MainNameIsCN = true
			break
		}
	}

	//  第一个品牌主名称 是中文，第二个不是 返回第一个名称
	if b1MainNameIsCN && !b2MainNameIsCN {
		return b1.MainName
	}
	//  第2个品牌主名称 是中文，第1个不是 返回第2个名称
	if b2MainNameIsCN && !b1MainNameIsCN {
		return b2.MainName
	}
	// 剩下可能性是 两个都是英文或者两个都是中文，谁短选谁
	if len(b1.MainName) <= len(b2.MainName) {
		return b1.MainName
	}
	return b2.MainName
}

func getName(brandP string) (e string, c string) {

	e = ""

	c = ""

	if isPureHan(brandP) {
		c = strings.TrimSpace(brandP)
		return
	}
	e = strings.TrimSpace(brandP)

	return

}

func isPureHan(str string) bool {
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}
