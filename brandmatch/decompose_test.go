package brandmatch

import (
	"fmt"
	"testing"
)

func TestDecomposeBrands(t *testing.T) {
	rawBrands := []string{"YULYNA/虞琳娜", "Franic/法兰琳卡", "亚细亚", "柏蕊诗", "美康粉黛", "奈呈", "粉后（美妆）", "PCXC", "Yalget/雅丽洁", "植颖", "巧彩格", "JudydoLL/橘朵", "Addiction", "MAQuillAGE/心机", "Laura Mercier", "rcma", "MODiSSA/梦迪莎", "ROLANJONA/露兰姬娜", "惠居家"}
	bs, err := DecomposeBrands(rawBrands)
	if err != nil {
		panic(err)
	}
	for _, b := range bs {
		fmt.Printf("%+v\n", b)
	}
}
