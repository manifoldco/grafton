package pages

const pageResources = `
{{if .Message }}
  <h2>Done: {{.Message}} - {{.Code}}</h2>
{{end}}
<h1>Resources</h1>
{{range .Resources}}
<h3>{{.Name}}</h3>
{{else}}
<h2>No Resources</h2>
{{end}}
<form method="post" action="/resources">
  <input type="submit" value="Provision a Resource">
</form>
`
