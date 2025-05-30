package dto

type RequestShortenLink struct {
	Url string `json:"url"`
	ShortName string `json:"short_name"`
}

type ResponseShortenLink struct {
	Url string `json:"url"`
}

type RequestEditShortLink struct {
	NewUrl string `json:"new_url"`
}

type ResponseGetShortLink struct {
	Url string `json:"url"`
}

type ResponseEditShortLink struct {
	Url string `json:"new_url"`
	ShortName string `json:"short_name"`
}

