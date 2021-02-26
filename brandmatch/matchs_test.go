package brandmatch

import (
	"reflect"
	"testing"
)

var deafultText = `
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

func Test_BrandMatch_MatchAll(t *testing.T) {
	brandRaws := []string{"埃森/essen", "OOTD", "角川", "一身"}
	type args struct {
		txt        string
		products   []string
		properties []string
		brandRaws  []string
	}
	tests := []struct {
		name string
		args args
		want []MatchInfo
	}{
		{
			name: "tm1",
			args: args{
				txt:        deafultText,
				products:   []string{"这"},
				properties: nil,
				brandRaws:  brandRaws,
			},
			want: []MatchInfo{
				{
					Brand:        "OOTD",
					ProductMatch: true,
				},
				{
					Brand:        "一身",
					ProductMatch: true,
				},
			},
		},
		{
			// name: "brand match but product not match",
			name: "single match",
			args: args{
				txt:        deafultText,
				products:   []string{"博物"},
				properties: nil,
				brandRaws:  brandRaws,
			},
			want: []MatchInfo{
				{
					Brand:        "角川",
					ProductMatch: true,
				},
				// {
				// Brand:        "一身",
				// ProductMatch: true,
				// },
			},
		},
		{
			name: "brand match but product not match",
			args: args{
				txt:        deafultText,
				products:   []string{"奇葩"},
				properties: nil,
				brandRaws:  []string{"unknown1", "unknown2", "角川"},
			},
			want: []MatchInfo{
				{
					Brand:        "角川",
					ProductMatch: false,
				},
				// {
				// Brand:        "一身",
				// ProductMatch: true,
				// },
			},
		},
		{
			name: "just must product",
			args: args{
				txt:        "一个奇葩的bug",
				products:   []string{"奇葩"},
				properties: nil,
				brandRaws:  []string{"unknown1", "unknown2", "角川"},
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bns, err := DecomposeBrands(test.args.brandRaws)
			if err != nil {
				t.Fatal(err)
			}
			var bbns []*BrandName
			for i := range bns {
				bbns = append(bbns, &bns[i])
			}
			bm := NewBrandMatch(bbns)
			matchInfo := bm.MatchAll(test.args.txt, test.args.products, test.args.properties)
			if !reflect.DeepEqual(matchInfo, test.want) {
				t.Errorf("testName:%s brandMatch.MatchAll() = %v, want %v", test.name, matchInfo, test.want)
			}
		})
	}
}

func Test_matchOneBrand(t *testing.T) {
	type args struct {
		products   []string
		properties []string
		txt        string
		matchBrand string
		isEN       bool
	}
	tests := []struct {
		name  string
		args  args
		want1 bool
		want2 bool
	}{
		{
			name: "en match",
			args: args{
				products:   []string{"日本", "女装"},
				properties: []string{"高端", "炫酷"},
				txt:        deafultText,
				matchBrand: "essen",
				isEN:       true,
			},
			want1: true,
			want2: true,
		},
		{
			name: "en product not match",
			args: args{
				products:   []string{"unknow", "日本", "女装"},
				properties: []string{"高端", "炫酷"},
				txt:        deafultText,
				matchBrand: "essen",
				isEN:       true,
			},
			want1: true,
			want2: false,
		}, {
			name: "cn product not match",
			args: args{
				products:   []string{"unknow", "日本", "女装"},
				properties: []string{"高端", "炫酷"},
				txt:        deafultText,
				matchBrand: "一身",
				isEN:       false,
			},
			want1: true,
			want2: false,
		}, {
			name: "cn match",
			args: args{
				products:   []string{"日本", "女装"},
				properties: []string{"高端", "炫酷"},
				txt:        deafultText,
				matchBrand: "一身",
				isEN:       false,
			},
			want1: true,
			want2: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got1, got2 := matchOneBrand(test.args.products, test.args.properties, test.args.txt, test.args.matchBrand, test.args.isEN)
			if got1 != test.want1 || got2 != test.want2 {
				t.Logf("testName:%s want1: %t, want2:%t, got1:%t, got2:%t\n", test.name, test.want1, test.want2, got1, got2)
				t.Fail()
			}
		})
	}
}

func Test_brandMatch_MatchAll(t *testing.T) {
	type fields struct {
		brands []*BrandName
	}
	type args struct {
		txt        string
		products   []string
		properties []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []MatchInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &brandMatch{
				brands: tt.fields.brands,
			}
			if got := bm.MatchAll(tt.args.txt, tt.args.products, tt.args.properties); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("brandMatch.MatchAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
