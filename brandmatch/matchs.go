package brandmatch

import "strings"

type BrandMatch interface {
	MatchAll(txt string, products []string, properties []string) []MatchInfo
}
type MatchInfo struct {
	brand        string
	productMatch bool
}

type brandMatch struct {
	brands []*BrandName
}

func NewBrandMatch(brands []*BrandName) BrandMatch {
	for _, b := range brands {
		b.genMatch()
	}
	return &brandMatch{brands: brands}
}

// MatchAll 匹配所有的品牌并返回匹配信息
// txt: 需要匹配的文本
// products: 需要匹配的产品 每个元素都必须匹配 举例: []string{ "香水" ,"木质"}
// properties: 属性 只需要匹配一个即可.有时候需要对产品制定属性 如背包 指定 []string{"双肩","通勤"}
// MatchInfo: 匹配出来的信息
func (bm *brandMatch) MatchAll(txt string, products []string, properties []string) []MatchInfo {
	txt = strings.ToLower(txt)
	mis := []MatchInfo{}
	for _, b := range bm.brands {
		mb := false
		mp := false
		for _, en := range b.enMatch {
			tmb, tmp := matchOneBrand(products, properties, txt, en, true)
			mb = mb || tmb
			mp = mp || tmp
			if mb && mp {
				break
			}
		}

		if mb && mp {
			mis = append(mis, MatchInfo{
				brand:        b.MainName,
				productMatch: mp,
			})
			continue
		}
		for _, cn := range b.cnMatch {
			tmb, tmp := matchOneBrand(products, properties, txt, cn, false)
			mb = mb || tmb
			mp = mp || tmp
			if mb && mp {
				break
			}
		}
		if mb {
			mi := MatchInfo{brand: b.MainName}
			if mp {
				mi.productMatch = true
			}
			mis = append(mis, mi)
		}
	}
	var restulstMatchInfos []MatchInfo
	for _, mi := range mis {
		if mi.productMatch {
			restulstMatchInfos = append(restulstMatchInfos, mi)
		}
	}

	if len(restulstMatchInfos) == 0 && len(mis) == 1 {
		restulstMatchInfos = append(restulstMatchInfos, mis[0])
	}
	return restulstMatchInfos
}

func matchOneBrand(products, shoulds []string, txt, brandMatch string, isEN bool) (bool, bool) {
	matchBrand := false
	matchProduct := true
	for _, p := range products {
		mb, mp := matchBrandAndProduct(txt, brandMatch, p, 10, isEN)
		if !mb {
			break
		}
		matchBrand = true
		// 这里的匹配需要次次匹配上才算数
		matchProduct = matchProduct && mp
	}
	// matchProduct = matchProduct || len(products) == 0
	if !matchBrand && !matchProduct {
		return false, false
	}
	if !matchProduct {
		// 上面的判断确定品牌是肯定匹配上了,那么就直接为true
		// 如果这里为false 下面的should 其实不需要了
		return true, false
	}
	sf := false
	for _, p := range shoulds {
		mb, mp := matchBrandAndProduct(txt, p, brandMatch, 10, isEN)
		// 一次匹配即可
		matchBrand = matchBrand || mb
		sf = sf || mp
		if sf {
			break
		}
	}
	// 如果shoulds 没有,就算匹配上
	sf = sf || len(shoulds) == 0
	// 如果没有匹配到品牌名则产品默认未匹配上
	matchProduct = matchProduct && matchBrand

	return matchBrand, sf && matchProduct
}

// func (bm *brandMatch) getContainerBrands(txt string) []*BrandName {
// for _, b := range bm.brands {
//
// }
// }

// type BrandNames []BrandName
// func (bs *BrandNames) MatchAll(txt string, products []string, properties []string) []BrandName {
// return nil
// }
