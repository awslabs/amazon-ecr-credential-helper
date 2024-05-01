{{range . -}}
{{if eq .LicenseName "Apache-2.0" -}}
** {{.Name}}; version {{.Version}} - {{.LicenseURL}}
{{end -}}
{{end -}}
