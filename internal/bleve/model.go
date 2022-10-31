package bleve

type Book struct {
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Chapters []Chapter `json:"chapters"`
	Text     string    `json:"text"`
}

type Chapter struct {
	ChapterTitle string      `json:"chapter_title"`
	Paragraphs   []Paragraph `json:"paragraphs"`
}

type Paragraph struct {
	Text string `json:"text"`
}
