package utils

import "strings"

func IsSeriesType(programmeType string) bool {
	return programmeType == "Season" || strings.Contains(strings.ToLower(programmeType), strings.ToLower("Series"))
}

func IsMovieType(programmeType string) bool {
	return !IsSeriesType(programmeType)
}
