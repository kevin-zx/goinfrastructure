package brandmatch

import (
	"fmt"
	"strings"
	"testing"
)

func Test_matchBrandAndProduct(t *testing.T) {
	type args struct {
		txt                string
		brand              string
		product            string
		productMinDistance int
		isEn               bool
	}
	deafultText := `
OOTD |爱上这条墨绿灰的直筒裤｜碾压Cos
今日周六，穿上这一身慵懒舒适的穿搭出门溜达了一圈～

早上去了 “角川武蔵野博物館”，傍晚去了吉祥寺，

看了很多古着和精品女装店后，

还是最喜欢自己身上，

这条松紧带的墨绿灰直筒裤！😁（无修图）

超级显瘦，又不勒肚子，可舒服自在了～😎

这一身都来自日本的高端女装essen家。

不过老粉丝们都别问了啊,

这一身的衣服都已售罄，我再给你们推荐其它的哈！

	`
	tests := []struct {
		name  string
		args  args
		want1 bool
		want2 bool
	}{
		{
			name: "en brand",
			args: args{
				txt:                deafultText,
				brand:              "essen",
				product:            "",
				productMinDistance: 10,
				isEn:               true,
			},
			want1: true,
			want2: true,
		},
		{
			name: "cn brand",
			args: args{
				txt:                deafultText,
				brand:              "吉祥寺",
				product:            "",
				productMinDistance: 10,
				isEn:               false,
			},
			want1: true,
			want2: true,
		},
		{
			name: "cn brand with product",
			args: args{
				txt:                deafultText,
				brand:              "吉祥寺",
				product:            "博物",
				productMinDistance: 10,
				isEn:               false,
			},
			want1: true,
			want2: true,
		},
		{
			name: "en brand with product",
			args: args{
				txt:                deafultText,
				brand:              "essen",
				product:            "高",
				productMinDistance: 10,
				isEn:               true,
			},
			want1: true,
			want2: true,
		},
		{
			name: "en brand with product and not matchp",
			args: args{
				txt:                deafultText,
				brand:              "essen",
				product:            "他",
				productMinDistance: 10,
				isEn:               true,
			},
			want1: true,
			want2: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got1, got2 := matchBrandAndProduct(test.args.txt, test.args.brand, test.args.product, test.args.productMinDistance, test.args.isEn)
			if got1 != test.want1 || got2 != test.want2 {
				t.Logf("testName:%s want1: %t, want2:%t, got1:%t, got2:%t\n", test.name, test.want1, test.want2, got1, got2)
				t.Fail()
			}
		})
	}
}

func Test_isDistanceOK(t *testing.T) {
	type args struct {
		txt      string
		product  string
		distance int
		start    bool
		end      bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test name",
			args: args{
				txt: `OOTD |爱上这条墨绿灰的直筒裤｜碾压Cos
		今日周六，穿上这一身慵懒舒适的穿搭出门溜达了一圈～

		早上去了 “角川武蔵野博物館”，傍晚去了`,
				product:  "博物",
				distance: 15,
				start:    true,
				end:      false,
			},
			want: true,
		}, {
			name: "test name",
			args: args{
				txt:      `这一身都来自日本的高端女装`,
				product:  "女装",
				distance: 15,
				start:    true,
				end:      false,
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isDistanceOK(test.args.txt, test.args.product, test.args.distance, test.args.start, test.args.end)
			if got != test.want {
				t.Logf("testName:%s, want:%v got:%v\n", test.name, test.want, got)
				t.Fail()
			}
		})
	}
}

func Test_getDistance(t *testing.T) {
	type args struct {
		text    string
		product string
	}
	tests := []struct {
		name      string
		args      args
		wantLeft  int
		wantRight int
	}{
		{
			name: "normal",
			args: args{
				text:    `早上去了 “角川武蔵野博物館”，傍晚去了`,
				product: "博物",
			},
			wantLeft:  11,
			wantRight: 7,
		},
		{
			name: "test end",
			args: args{
				text:    `这一身都来自日本的高端女装`,
				product: "女装",
			},
			wantLeft:  11,
			wantRight: 0,
		},
		{
			name: "test start",
			args: args{
				text:    `这一身都来自日本的高端女装`,
				product: "这一",
			},
			wantLeft:  0,
			wantRight: 11,
		},
		{
			name: "test two direct",
			args: args{
				text:    `这一身女都来自日本的高端女装`,
				product: "女",
			},
			wantLeft:  3,
			wantRight: 1,
		},
	}
	for _, test := range tests {
		fmt.Println("txt len", len(strings.Split(test.args.text, "")))
		t.Run(test.name, func(t *testing.T) {
			gotLeft, gotRight := getDistance(test.args.text, test.args.product)
			if gotLeft != test.wantLeft || gotRight != test.wantRight {
				t.Logf("testName:%s, wantLeft:%v gotLeft:%v wantRight:%v gotRight:%v\n", test.name, test.wantLeft, gotLeft, test.wantRight, gotRight)
				t.Fail()
			}
		})
	}
}

func Test_strRuneIndex(t *testing.T) {
	type args struct {
		txt   string
		sub   string
		right bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test name",
			args: args{
				txt:   "",
				sub:   "",
				right: false,
			},
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
