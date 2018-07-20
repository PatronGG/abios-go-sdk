package structs

// OrganisationStructPaginated holds a list of OrganisationStruct as well as information about pages.
type OrganisationStructPaginated struct {
	LastPage    int64                `json:"last_page,omitempty"`
	CurrentPage int64                `json:"current_page,omitempty"`
	Data        []OrganisationStruct `json:"data,omitempty"`
}

// OrganisationStruct represents a groupings of Teams across different Games.
type OrganisationStruct struct {
	Id    int64        `json:"id,omitempty"`
	Name  string       `json:"name,omitempty"`
	Teams []TeamStruct `json:"teams"`
}
