// Package gqlconn provides utility types/functions for GraphQL Cursor
// Connections.
// <https://facebook.github.io/relay/graphql/connections.htm>
package gqlconn

import (
	"unicode"

	"github.com/graphql-go/graphql"
)

// PageInfo provides a struct for page information.
type PageInfo struct {
	HasNextPage bool   `json"hasNextPage"`
	HasPrevPage bool   `json"hasPrevPage"`
	StartCursor string `json"startCursor"`
	EndCursor   string `json"endCursor"`
}

// PageInfoType provides named type for page information.
var PageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"hasNextPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "whether more edges exist prior",
		},
		"hasPrevPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "whether more edges exist following",
		},
		"startCursor": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "cursor of the top edge",
		},
		"endCursor": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "cursor of the tail edge",
		},
	},
	Description: "pagination information",
})

// Edge provides a struct for each edge in connection.
type Edge struct {
	Node   interface{} `json:"node"`
	Cursor string      `json:"cursor"`
}

// Edges provides array type of Edge.
type Edges []Edge

// ResolveCountFn is a resolver for edge count.
type ResolveCountFn func(graphql.ResolveParams) (int, error)

// ResolveEdgesFn is a resolver for Edges ([]Edge).
type ResolveEdgesFn func(graphql.ResolveParams) (Edges, error)

// ResolvePageInfoFn is a resolver for PageInfo.
type ResolvePageInfoFn func(graphql.ResolveParams) (*PageInfo, error)

// CollectionConfig is a configuration for the collection type.
type CollectionConfig struct {
	NodeType graphql.Type

	ResolveCount    ResolveCountFn
	ResolveEdges    ResolveEdgesFn
	ResolvePageInfo ResolvePageInfoFn
}

func generateNames(base string) (typeName, smallName string) {
	runes := []rune(base)
	if len(runes) == 0 {
		return "", ""
	}
	if unicode.IsLower(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes), base
	}
	if unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
		return base, string(runes)
	}
	return base, base
}

// NewCollectionType creates a name type for collection.
func NewCollectionType(cfg CollectionConfig) *graphql.Object {
	typeName, smallName := generateNames(cfg.NodeType.Name())
	// create edge type.
	edgeType := graphql.NewObject(graphql.ObjectConfig{
		Name: typeName + "Edge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type:        cfg.NodeType,
				Description: "each " + smallName + " node",
			},
			"cursor": &graphql.Field{
				Type:        graphql.String,
				Description: "cursor for this node",
			},
		},
		Description: "edge type for " + smallName,
	})
	// create connection type.
	var (
		resolveCount    = cfg.ResolveCount
		resolveEdges    = cfg.ResolveEdges
		resolvePageInfo = cfg.ResolvePageInfo
	)
	connType := graphql.NewObject(graphql.ObjectConfig{
		Name: typeName + "sConnection",
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type:        graphql.Int,
				Description: "count of whole edges",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveCount(p)
				},
			},
			"edges": &graphql.Field{
				Type:        graphql.NewList(edgeType),
				Description: "count of whole edges",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveEdges(p)
				},
			},
			"pageInfo": &graphql.Field{
				Type:        PageInfoType,
				Description: "count of whole edges",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolvePageInfo(p)
				},
			},
		},
		Description: typeName + "s connection",
	})
	return connType
}
