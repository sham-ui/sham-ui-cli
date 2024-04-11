package slugify

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"context"
	"testing"
)

func TestService_CategorySlug(t *testing.T) {
	srv := New()

	testCases := []struct {
		in  string
		out model.CategorySlug
	}{
		{"DOBROSLAWZYBORT", "dobroslawzybort"},
		{"Dobroslaw Zybort", "dobroslaw-zybort"},
		{"  Dobroslaw     Zybort  ?", "dobroslaw-zybort"},
		{"Dobrosław Żybort", "dobroslaw-zybort"},
		{"Ala ma 6 kotów.", "ala-ma-6-kotow"},
		{"áÁàÀãÃâÂäÄąĄą̊Ą̊", "aaaaaaaaaaaaaa"},
		{"ćĆĉĈçÇčČ", "cccccccc"},
		{"éÉèÈẽẼêÊëËęĘěĚ", "eeeeeeeeeeeeee"},
		{"íÍìÌĩĨîÎïÏįĮ", "iiiiiiiiiiii"},
		{"łŁ", "ll"},
		{"ńŃ", "nn"},
		{"óÓòÒõÕôÔöÖǫǪǭǬø", "ooooooooooooooo"},
		{"śŚšŠ", "ssss"},
		{"řŘ", "rr"},
		{"ťŤ", "tt"},
		{"úÚùÙũŨûÛüÜųŲůŮ", "uuuuuuuuuuuuuu"},
		{"y̨Y̨ýÝ", "yyyy"},
		{"źŹżŹžŽ", "zzzzzz"},
		{"·/,:;`˜'\"", ""},
		{"2000–2013", "2000-2013"},
		{"style—not", "style-not"},
		{"test_slug", "test_slug"},
		{"_test_slug_", "test_slug"},
		{"-test-slug-", "test-slug"},
		{"Æ", "ae"},
		{"Ich heiße", "ich-heisse"},
		{"% 5 @ 4 $ 3 / 2 & 1 & 2 # 3 @ 4 _ 5", "5-at-4-3-2-and-1-and-2-3-at-4-_-5"},
		{"This & that", "this-and-that"},
		{"fácil €", "facil-eu"},
		{"smile ☺", "smile"},
		{"Hellö Wörld хелло ворлд", "hello-world-khello-vorld"},
		{"\"C'est déjà l’été.\"", "cest-deja-lete"},
		{"jaja---lol-méméméoo--a", "jaja-lol-mememeoo-a"},
		{"影師", "ying-shi"},
		{"Đanković & Kožušček", "dankovic-and-kozuscek"},
		{"ĂăÂâÎîȘșȚț", "aaaaiisstt"},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			res := srv.SlugifyCategory(context.Background(), test.in)
			asserts.Equals(t, test.out, res)
		})
	}
}

func TestService_TagSlug(t *testing.T) {
	srv := New()

	testCases := []struct {
		in  string
		out model.TagSlug
	}{
		{"DOBROSLAWZYBORT", "dobroslawzybort"},
		{"Dobroslaw Zybort", "dobroslaw-zybort"},
		{"  Dobroslaw     Zybort  ?", "dobroslaw-zybort"},
		{"Dobrosław Żybort", "dobroslaw-zybort"},
		{"Ala ma 6 kotów.", "ala-ma-6-kotow"},
		{"áÁàÀãÃâÂäÄąĄą̊Ą̊", "aaaaaaaaaaaaaa"},
		{"ćĆĉĈçÇčČ", "cccccccc"},
		{"éÉèÈẽẼêÊëËęĘěĚ", "eeeeeeeeeeeeee"},
		{"íÍìÌĩĨîÎïÏįĮ", "iiiiiiiiiiii"},
		{"łŁ", "ll"},
		{"ńŃ", "nn"},
		{"óÓòÒõÕôÔöÖǫǪǭǬø", "ooooooooooooooo"},
		{"śŚšŠ", "ssss"},
		{"řŘ", "rr"},
		{"ťŤ", "tt"},
		{"úÚùÙũŨûÛüÜųŲůŮ", "uuuuuuuuuuuuuu"},
		{"y̨Y̨ýÝ", "yyyy"},
		{"źŹżŹžŽ", "zzzzzz"},
		{"·/,:;`˜'\"", ""},
		{"2000–2013", "2000-2013"},
		{"style—not", "style-not"},
		{"test_slug", "test_slug"},
		{"_test_slug_", "test_slug"},
		{"-test-slug-", "test-slug"},
		{"Æ", "ae"},
		{"Ich heiße", "ich-heisse"},
		{"% 5 @ 4 $ 3 / 2 & 1 & 2 # 3 @ 4 _ 5", "5-at-4-3-2-and-1-and-2-3-at-4-_-5"},
		{"This & that", "this-and-that"},
		{"fácil €", "facil-eu"},
		{"smile ☺", "smile"},
		{"Hellö Wörld хелло ворлд", "hello-world-khello-vorld"},
		{"\"C'est déjà l’été.\"", "cest-deja-lete"},
		{"jaja---lol-méméméoo--a", "jaja-lol-mememeoo-a"},
		{"影師", "ying-shi"},
		{"Đanković & Kožušček", "dankovic-and-kozuscek"},
		{"ĂăÂâÎîȘșȚț", "aaaaiisstt"},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			res := srv.SlugifyTag(context.Background(), test.in)
			asserts.Equals(t, test.out, res)
		})
	}
}

func TestService_ArticleSlug(t *testing.T) {
	srv := New()

	testCases := []struct {
		in  string
		out model.ArticleSlug
	}{
		{"DOBROSLAWZYBORT", "dobroslawzybort"},
		{"Dobroslaw Zybort", "dobroslaw-zybort"},
		{"  Dobroslaw     Zybort  ?", "dobroslaw-zybort"},
		{"Dobrosław Żybort", "dobroslaw-zybort"},
		{"Ala ma 6 kotów.", "ala-ma-6-kotow"},
		{"áÁàÀãÃâÂäÄąĄą̊Ą̊", "aaaaaaaaaaaaaa"},
		{"ćĆĉĈçÇčČ", "cccccccc"},
		{"éÉèÈẽẼêÊëËęĘěĚ", "eeeeeeeeeeeeee"},
		{"íÍìÌĩĨîÎïÏįĮ", "iiiiiiiiiiii"},
		{"łŁ", "ll"},
		{"ńŃ", "nn"},
		{"óÓòÒõÕôÔöÖǫǪǭǬø", "ooooooooooooooo"},
		{"śŚšŠ", "ssss"},
		{"řŘ", "rr"},
		{"ťŤ", "tt"},
		{"úÚùÙũŨûÛüÜųŲůŮ", "uuuuuuuuuuuuuu"},
		{"y̨Y̨ýÝ", "yyyy"},
		{"źŹżŹžŽ", "zzzzzz"},
		{"·/,:;`˜'\"", ""},
		{"2000–2013", "2000-2013"},
		{"style—not", "style-not"},
		{"test_slug", "test_slug"},
		{"_test_slug_", "test_slug"},
		{"-test-slug-", "test-slug"},
		{"Æ", "ae"},
		{"Ich heiße", "ich-heisse"},
		{"% 5 @ 4 $ 3 / 2 & 1 & 2 # 3 @ 4 _ 5", "5-at-4-3-2-and-1-and-2-3-at-4-_-5"},
		{"This & that", "this-and-that"},
		{"fácil €", "facil-eu"},
		{"smile ☺", "smile"},
		{"Hellö Wörld хелло ворлд", "hello-world-khello-vorld"},
		{"\"C'est déjà l’été.\"", "cest-deja-lete"},
		{"jaja---lol-méméméoo--a", "jaja-lol-mememeoo-a"},
		{"影師", "ying-shi"},
		{"Đanković & Kožušček", "dankovic-and-kozuscek"},
		{"ĂăÂâÎîȘșȚț", "aaaaiisstt"},
	}

	for _, test := range testCases {
		t.Run(test.in, func(t *testing.T) {
			res := srv.SlugifyArticle(context.Background(), test.in)
			asserts.Equals(t, test.out, res)
		})
	}
}
