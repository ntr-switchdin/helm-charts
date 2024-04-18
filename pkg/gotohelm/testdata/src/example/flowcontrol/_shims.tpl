{{- /* Generated from "" */ -}}

{{- define "_shims.typetest" -}}
{{- $type := (index .a 0) -}}
{{- $value := (index .a 1) -}}
{{- dict "r" (list $value (typeIs $type $value)) | toJson -}}
{{- end -}}

{{- define "_shims.dicttest" -}}
{{- $dict := (index .a 0) -}}
{{- $key := (index .a 1) -}}
{{- if (hasKey $dict $key) -}}
{{- (dict "r" (list (index $dict $key) true)) | toJson -}}
{{- else -}}
{{- (dict "r" (list "" false)) | toJson -}}
{{- end -}}
{{- end -}}

{{- define "_shims.typeassertion" -}}
{{- $type := (index .a 0) -}}
{{- $value := (index .a 1) -}}
{{- if (not (typeIs $type $value)) -}}
{{- (fail "TODO MAKE THIS A NICE MESSAGE") -}}
{{- end -}}
{{- (dict "r" $value) | toJson -}}
{{- end -}}

{{- define "_shims.compact" -}}
{{- $out := (dict) -}}
{{- range $i, $e := (index .a 0) }}
{{- $_ := (set $out (printf "T%d" (add1 $i)) $e) -}}
{{- end -}}
{{- (dict "r" $out) | toJson -}}
{{- end -}}

{{- define "_shims.sitobytes" -}}
  {{/*
  This template converts the incoming SI value to whole number bytes.
  Input can be: b | B | k | K | m | M | g | G | Ki | Mi | Gi
  Or number without suffix
  */}}
  {{- $si := . -}}
  {{- if not (typeIs "string" . ) -}}
    {{- $si = int64 $si | toString -}}
  {{- end -}}
  {{- $bytes := 0 -}}
  {{- if or (hasSuffix "B" $si) (hasSuffix "b" $si) -}}
    {{- $bytes = $si | trimSuffix "B" | trimSuffix "b" | float64 | floor -}}
  {{- else if or (hasSuffix "K" $si) (hasSuffix "k" $si) -}}
    {{- $raw := $si | trimSuffix "K" | trimSuffix "k" | float64 -}}
    {{- $bytes = mulf $raw (mul 1000) | floor -}}
  {{- else if or (hasSuffix "M" $si) (hasSuffix "m" $si) -}}
    {{- $raw := $si | trimSuffix "M" | trimSuffix "m" | float64 -}}
    {{- $bytes = mulf $raw (mul 1000 1000) | floor -}}
  {{- else if or (hasSuffix "G" $si) (hasSuffix "g" $si) -}}
    {{- $raw := $si | trimSuffix "G" | trimSuffix "g" | float64 -}}
    {{- $bytes = mulf $raw (mul 1000 1000 1000) | floor -}}
  {{- else if hasSuffix "Ki" $si -}}
    {{- $raw := $si | trimSuffix "Ki" | float64 -}}
    {{- $bytes = mulf $raw (mul 1024) | floor -}}
  {{- else if hasSuffix "Mi" $si -}}
    {{- $raw := $si | trimSuffix "Mi" | float64 -}}
    {{- $bytes = mulf $raw (mul 1024 1024) | floor -}}
  {{- else if hasSuffix "Gi" $si -}}
    {{- $raw := $si | trimSuffix "Gi" | float64 -}}
    {{- $bytes = mulf $raw (mul 1024 1024 1024) | floor -}}
  {{- else if (mustRegexMatch "^[0-9]+$" $si) -}}
    {{- $bytes = $si -}}
  {{- else -}}
    {{- printf "\n%s is invalid SI quantity\nSuffixes can be: b | B | k | K | m | M | g | G | Ki | Mi | Gi or without any Suffixes" $si | fail -}}
  {{- end -}}
  {{- $bytes | int64 -}}
{{- end -}}
