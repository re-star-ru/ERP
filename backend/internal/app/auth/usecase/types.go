package authUsecase

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
)

const (
	// EffectAllow allow effect
	EffectAllow = Effect("allow")
	// EffectDeny deny effect
	EffectDeny = Effect("deny")
)

// Resource alias Resource
type Resource string

// Action alias type Action
type Action string

func (act Action) String() string {
	return string(act)
}

// Effect the type of effect
type Effect string

func (eff Effect) String() string {
	return string(eff)
}

// Policy the type of policy
type Policy struct {
	Resource
	Action
	Effect
}

// GetEffect returns effect of resource, default is allow
func (p *Policy) GetEffect() string {
	eft := p.Effect
	if eft == "" {
		eft = EffectAllow
	}

	return eft.String()
}

func (p *Policy) String() string {
	return p.Resource.String() + ":" + p.Action.String() + ":" + p.GetEffect()
}

func (res Resource) String() string {
	return string(res)
}

// RelativeTo returns relative resource to other resource
func (res Resource) RelativeTo(other Resource) (Resource, error) {
	prefix := other.String()
	str := res.String()

	if !strings.HasPrefix(str, prefix) {
		return Resource(""), errors.New("value error")
	}

	relative := strings.TrimPrefix(strings.TrimPrefix(str, prefix), "/")
	if relative == "" {
		relative = "."
	}

	return Resource(relative), nil
}

// Subresource returns subresource
func (res Resource) Subresource(resources ...Resource) Resource {
	elements := []string{res.String()}

	for _, resource := range resources {
		elements = append(elements, resource.String())
	}

	return Resource(path.Join(elements...))
}

// GetNamespace returns namespace from resource
func (res Resource) GetNamespace() (Namespace, error) {
	return nil, fmt.Errorf("no namespace found for %s", res)
}

var (
	parsesMu sync.RWMutex
	parses   = map[string]NamespaceParse{}
)

// NamespaceParse parse namespace from the resource
type NamespaceParse func(Resource) (Namespace, bool)

// Namespace the namespace interface
type Namespace interface {
	// Kind returns the kind of namespace
	Kind() string
	// Resource returns new resource for subresources with the namespace
	Resource(subresources ...Resource) Resource
	// Identity returns identity attached with namespace
	Identity() interface{}
	// GetPolicies returns all policies of the namespace
	GetPolicies() []*Policy
}

// RegistryNamespaceParse ...
func RegistryNamespaceParse(name string, parse NamespaceParse) {
	parsesMu.Lock()
	defer parsesMu.Unlock()
	if parse == nil {
		panic("permission: Register namespace parse is nil")
	}
	if _, dup := parses[name]; dup {
		panic("permission: Register called twice for namespace parse " + name)
	}

	parses[name] = parse
}

// NamespaceFromResource returns namespace from resource
func NamespaceFromResource(resource Resource) (Namespace, bool) {
	parsesMu.RLock()
	defer parsesMu.RUnlock()

	for _, parse := range parses {
		if ns, ok := parse(resource); ok {
			return ns, true
		}
	}

	return nil, false
}

// ResourceAllowedInNamespace returns true when resource's namespace equal the ns
func ResourceAllowedInNamespace(resource Resource, ns Namespace) bool {
	n, ok := NamespaceFromResource(resource)
	if ok {
		return n.Kind() == ns.Kind() && n.Identity() == ns.Identity()
	}

	return false
}
