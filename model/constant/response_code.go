package constant

type ResCode string

const (
	APISuccess       ResCode = "%s100"
	APIUnauthorized  ResCode = "%s401"
	APIInternalError ResCode = "%s500"
)

const (
	ClientAlreadyExists     ResCode = "%s101"
	EmptyMandatoryParameter ResCode = "%s102"
	InvalidTypeParameter    ResCode = "%s103"
	InvalidRequestBody      ResCode = "%s104"
)
