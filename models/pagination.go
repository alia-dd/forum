package models

import "html/template"

type Pagination struct {
	Page     int
	NextPage int
	PrevPage int
	HasNext  bool
	HasPrev  bool
	Base     template.URL
}
