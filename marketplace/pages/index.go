package pages

import "html/template"

// Error holds the HTML template for the error page
var Error *template.Template

// Resources holds the HTML template for the resource page
var Resources *template.Template

func init() {
	var err error
	// Error
	Error = template.New("Error")
	Error, err = Error.Parse(pageError)
	if err != nil {
		panic(err)
	}
	// Resources
	Resources = template.New("Resources")
	Resources, err = Resources.Parse(pageResources)
	if err != nil {
		panic(err)
	}
}
