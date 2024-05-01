{{ range . -}}
{{ if ne .LicenseName "Apache-2.0" -}}
--------------------------------------------------------------------------------
** {{.Name}}; version {{.Version}} - https://{{.Name}}

{{ .LicenseText }}

{{end -}}
{{end -}}
