package page

type TagMeta map[string]string

type TagA map[string]string

type Page struct {
	WebAddress		string
	TagsMeta		[]TagMeta
	TagsA			[]TagA

	Data 			string
}

