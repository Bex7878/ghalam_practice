package domain

type DataExtractor interface {
	ExtractValues(htmlContent string) (map[string]string, error)
}
