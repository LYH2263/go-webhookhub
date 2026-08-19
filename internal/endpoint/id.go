package endpoint

import "github.com/LYH2263/go-webhookhub/internal/idgen"

func NewID() string { return idgen.EndpointID() }

func IsEndpointID(id string) bool { return idgen.IsPrefixed(id, "ep_") }
