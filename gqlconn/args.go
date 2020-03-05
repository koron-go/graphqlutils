package gqlconn

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

// PaginationParams is a parameter set for pagination.
type PaginationParams struct {
	// Backword is `true` for backward pagination, `false` for forward.
	Backword bool

	// Pivot is cursor of the exclusive pivot for pagination.
	Pivot string

	// Size is number of edges in a page.
	Size int
}

// PaginationType is directions of pagination.
type PaginationType int

const (
	// Both means bidirectional pagination.
	Both PaginationType = 0
	// Forward means forward pagination.
	Forward = 1
	// Backword means backword pagination.
	Backword = 2
)

// PaginationConfig is configuration for pagination args.
type PaginationConfig struct {
	Type         PaginationType
	DefaultFirst int
	DefaultLast  int
}

func (pc PaginationConfig) isForward() bool {
	return pc.Type == Both || pc.Type == Forward
}

func (pc PaginationConfig) isBackword() bool {
	return pc.Type == Both || pc.Type == Forward
}

// FieldConfigArgument creates FieldConfigArgument for pagination.
func (pc PaginationConfig) FieldConfigArgument() graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}
	if pc.isForward() {
		first := &graphql.ArgumentConfig{
			Type:        graphql.Int,
			Description: "number of edges to get",
		}
		if pc.DefaultFirst > 0 {
			first.DefaultValue = pc.DefaultFirst
		}
		args["first"] = first
		args["after"] = &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "cursor for edge. after the edge to get (exclusively)",
		}
	}
	if pc.isBackword() {
		last := &graphql.ArgumentConfig{
			Type:        graphql.Int,
			Description: "number of edges to get",
		}
		if pc.DefaultLast > 0 {
			last.DefaultValue = pc.DefaultLast
		}
		args["last"] = last
		args["before"] = &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "cursor for edge. before the edge to get (exclusively)",
		}
	}
	return args
}

// Parse parses args as PaginationParams.
func (pc PaginationConfig) Parse(args map[string]interface{}) (*PaginationParams, error) {

	if pc.isForward() {
		if v, ok := args["first"]; ok {
			size, ok := v.(int)
			if !ok || size < 0 {
				return nil, fmt.Errorf("invalid value for \"first\": %v", v)
			}
			var pivot string
			if w, ok := args["after"]; ok {
				pivot, ok = w.(string)
				if !ok {
					return nil, fmt.Errorf("unexpected \"after\" type: want=string got=%T", w)
				}
			}
			return &PaginationParams{
				Backword: false,
				Pivot:    pivot,
				Size:     size,
			}, nil
		}
	}

	if pc.isBackword() {
		if v, ok := args["last"]; ok {
			size, ok := v.(int)
			if !ok || size < 0 {
				return nil, fmt.Errorf("invalid value for \"last\": %v", v)
			}
			var pivot string
			if w, ok := args["before"]; ok {
				pivot, ok = w.(string)
				if !ok {
					return nil, fmt.Errorf("unexpected \"before\" type: want=string got=%T", w)
				}
			}
			return &PaginationParams{
				Backword: true,
				Pivot:    pivot,
				Size:     size,
			}, nil
		}
	}

	if pc.isForward() {
		if v, ok := args["after"]; ok {
			pivot, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected \"after\" type: want=string got=%T", v)
			}
			if pc.DefaultFirst < 0 {
				return nil, fmt.Errorf("negative default size: %d", pc.DefaultFirst)
			}
			return &PaginationParams{
				Backword: false,
				Pivot:    pivot,
				Size:     pc.DefaultFirst,
			}, nil
		}
	}

	if pc.isBackword() {
		if v, ok := args["before"]; ok {
			pivot, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected \"before\" type: want=string got=%T", v)
			}
			if pc.DefaultLast < 0 {
				return nil, fmt.Errorf("negative default size: %d", pc.DefaultLast)
			}
			return &PaginationParams{
				Backword: true,
				Pivot:    pivot,
				Size:     pc.DefaultLast,
			}, nil
		}
	}

	if pc.DefaultFirst < 0 {
		return nil, fmt.Errorf("negative default size: %d", pc.DefaultFirst)
	}
	return &PaginationParams{
		Backword: pc.isBackword(),
		Pivot:    "",
		Size:     pc.DefaultFirst,
	}, nil
}
