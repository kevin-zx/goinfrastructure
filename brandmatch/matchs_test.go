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

var brandRaws = []string{"埃森/essen", "OOTD", "角川", "一身"}

func Test_BrandMatch_MatchAll(t *testing.T) {
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
		// {
		//   name: "tm1",
		//   args: args{
		//     txt:        deafultText,
		//     products:   []string{"这"},
		//     properties: nil,
		//     brandRaws:  brandRaws,
		//   },
		//   want: []MatchInfo{
		//     {
		//       Brand:        "OOTD",
		//       ProductMatch: true,
		//     },
		//     {
		//       Brand:        "一身",
		//       ProductMatch: true,
		//     },
		//   },
		// },
		// {
		//   // name: "brand match but product not match",
		//   name: "single match",
		//   args: args{
		//     txt:        deafultText,
		//     products:   []string{"博物"},
		//     properties: nil,
		//     brandRaws:  brandRaws,
		//   },
		//   want: []MatchInfo{
		//     {
		//       Brand:        "角川",
		//       ProductMatch: true,
		//     },
		//     // {
		//     // Brand:        "一身",
		//     // ProductMatch: true,
		//     // },
		//   },
		// },
		// {
		//   name: "brand match but product not match",
		//   args: args{
		//     txt:        deafultText,
		//     products:   []string{"奇葩"},
		//     properties: nil,
		//     brandRaws:  []string{"unknown1", "unknown2", "角川"},
		//   },
		//   want: []MatchInfo{
		//     {
		//       Brand:        "角川",
		//       ProductMatch: false,
		//     },
		//     // {
		//     // Brand:        "一身",
		//     // ProductMatch: true,
		//     // },
		//   },
		// },
		// {
		//   name: "just must product",
		//   args: args{
		//     txt:        "一个奇葩的bug",
		//     products:   []string{"奇葩"},
		//     properties: nil,
		//     brandRaws:  []string{"unknown1", "unknown2", "角川"},
		//   },
		//   want: nil,
		// },
		{
			name: "distance test",
			args: args{
				txt: `看我发现了什么！南京居然有Been Trill的门店！
-
Been Trill @BEEN TRILL 是一个在美国非常🔥的街头潮牌，设计团队由设计师、艺术家和DJ等人物组成。跟街头潮牌Stussy和HBA成功合作之后，Been Trill开始活跃在美国潮牌圈。年轻、时尚、潮流、活力、运动感都是品牌设计所追求的[赞R]
-
Been Trill的设计极具创造力和视觉艺术感，它有两个标志性的特征，一是像油漆往下流淌的字体 logo ，二是标签符号“#”，不管是在门店装饰里还是在服装设计的小细节里都处处出现这些元素💯
Been Trill 和HBA、KTZ、Undefeated、Stussy、SSUR等众多高街品牌都曾推出过联名的合作系列，每次的联名款都大受好评。
另外，它也一直都受到时尚圈和国内外明星爱豆的喜爱，陈冠希、吴亦凡、王源和Bobby等韩国一众爱豆都很喜欢这个品牌的衣服[得意R][得意R][得意R]
-
除了服装之外，帽子、背包、口罩、袜子等配饰也很优秀。尤其是Been Trill的棒球帽，既潮流又有运动感，帽型舒服好戴，也是明星同款“火爆区”。
质量超好，设计又时尚吸睛，即使是平常不怎么穿潮牌的人，路过门店也会被这种设计吸引吧~如果是本来就喜欢街头潮牌的人，那就更幸运啦，直接就能在南京的门店买到Been Trill，完全省去找代购的麻烦[派对R][派对R][派对R]
——————————
店铺：Been Trill
地址：南京市建邺区应天大街888号金鹰世界4楼（近集庆门大街地铁站）
价位：T恤500左右，外套价格在1000-2000区间左右（实体门店会有折扣活动，南京金鹰门店现在两件会打折）
——————————
PS：虽然🍑上有店，但是在实地探店后发现南京金鹰的门店里还是有很多🍑上没有更新的款式，而且都更好看，所以条件允许的话还是建议去实体门店逛一逛~而且门店还有优惠活动！
@南京薯  @穿搭薯 #南京吃喝玩乐 #南京购物 #南京探店 #周末探店 #跟我来探店 #南京潮牌 #南京网红店打卡`,
				products:   []string{"箱", "包"},
				properties: nil,
				brandRaws:  []string{"GASTON LUGA", "迪士尼 (Disney)", "MYD", "路易威登 (lv)", "耐克 (NIKE)", "匡威 (Converse)", "苹果 (Apple)", "途加", "路过 (Luguo)", "alloy", "JANSPORT", "小米 (MI)", "迪奥 (DIOR)", "爱马仕 (Hermes)", "古驰 (gucci)", "ACE", "日默瓦 (RIMOWA)", "CILOCALA", "拉科斯特 (LACOSTE)", "漫游 (ROAMING)", "皇冠 (CROWN)", "倒计时 (COUNT DOWN)", "北极狐 (FJALLRAVEN)", "蔻驰 (COACH)", "不莱玫 (bromen bags)", "MUJI", "顽皮 (Wanpy)", "内涵", "90分", "熊猫 (PANDA)"},
			},
			want: []MatchInfo{{
				Brand:        "路过 (Luguo)",
				ProductMatch: false,
			}},
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
