{{ range . -}}
{{ if ne .LicenseName "Apache-2.0" -}}
--------------------------------------------------------------------------------
** {{.Name}}; version {{.Version}} - {{.LicenseURL}}

{{ .LicenseText }}

{{end -}}
{{end -}}
