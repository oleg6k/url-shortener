package types

type ShortingJSONBody struct {
	URL string `json:"url" validate:"required,url"`
}
