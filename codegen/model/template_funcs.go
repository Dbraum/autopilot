package model

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func (d ProjectData) Funcs() template.FuncMap {
	return template.FuncMap{
		// string utils
		"join":            strings.Join,
		"lower":           strings.ToLower,
		"lower_camel":     strcase.ToLowerCamel,
		"upper_camel":     strcase.ToCamel,
		"snake":           strcase.ToSnake,
		"split":           splitTrimEmpty,
		"string_contains": strings.Contains,

		// project data
		"has_outputs":          HasOutputs,
		"has_inputs":           HasInputs,
		"is_final":             isFinal,
		"is_metrics":           isMetrics,
		"worker_import_prefix": WorkerDirName,
		"worker_package":       d.workerPackage,
		"needs_metrics":        d.NeedsMetrics,
		"unique_outputs":       d.UniqueOutputs,
		"unique_params":        d.UniqueParams,
	}
}

func splitTrimEmpty(s, sep string) []string {
	return strings.Split(strings.TrimSpace(s), sep)
}

func HasOutputs(phase Phase) bool {
	return len(phase.Outputs) > 0
}

func HasInputs(phase Phase) bool {
	return len(phase.Inputs) > 0
}

func isFinal(phase Phase) bool {
	return !HasOutputs(phase) && !HasInputs(phase)
}

func isMetrics(param Parameter) bool {
	return param.Equals(Metrics)
}

func WorkerDirName(phase Phase) string {
	return strings.ToLower(phase.Name)
}

func (d ProjectData) workerPackage(phase Phase) string {
	return filepath.Join(phase.Project.ProjectPackage, "pkg", "workers", strings.ToLower(phase.Name))
}
