package operation

import (
	"fmt"
)

const (
	servicesLayer = "services"
)

func buildOperation(layer, service, method string) string {
	return fmt.Sprintf("%s.%s.%s", layer, service, method)
}

func ServicesOperation(service, method string) string {
	return buildOperation(servicesLayer, service, method)
}
