package platform

import (
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/slices"
)

func NewGraph(platformModule *tfconfig.Module) Graph {
	g := Graph{}
	g.Parse(platformModule)
	return g
}

func (g *Graph) Parse(srcModule *tfconfig.Module) {
	toTraverse := map[BlockID]struct{}{}

	for k := range srcModule.ModuleCalls {
		if !strings.HasPrefix(k, ComponentPrefix) {
			continue
		}
		bID := NewBlockID(BlockType_ModuleCall, k)
		toTraverse[bID] = struct{}{}
	}

	for k := range srcModule.Outputs {
		bID := NewBlockID(BlockType_Output, k)
		toTraverse[bID] = struct{}{}
	}

	for len(toTraverse) > 0 {
		for bID := range toTraverse {
			if g.GetByID(bID) != nil {
				delete(toTraverse, bID)
				continue
			}

			blockRequirements := bID.FindRequirements(srcModule)
			g.Append(bID, blockRequirements)
			delete(toTraverse, bID)

			for _, reqBId := range blockRequirements {
				toTraverse[reqBId] = struct{}{}
			}
		}
	}

	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].ID < (*g)[j].ID
	})
}

func (g Graph) GetByID(id BlockID) *GraphNode {
	for i, v := range g {
		if v.ID == id {
			return &g[i]
		}
	}

	return nil
}

func (g *Graph) Append(bID BlockID, requirements []BlockID) *GraphNode {
	(*g) = append((*g), GraphNode{ID: bID, Requirements: requirements})
	return &(*g)[len(*g)-1]
}

type GraphWalkerCB func(blockId BlockID) error

func (g *Graph) Walk(roots []BlockID, fu GraphWalkerCB) error {
	roots = slices.Compact(roots)
	traverser := make([]BlockID, len(roots)) // nodes before `i` are visited and after `i` are queued
	copy(traverser, roots)

	err := g.traverseRootBlocks(&traverser, fu)
	if err != nil {
		return err
	}

	err = g.traverseOutputBlocks(&traverser, fu)
	if err != nil {
		return err
	}

	return nil
}

func (g *Graph) traverseRootBlocks(traverser *[]BlockID, fu GraphWalkerCB) error {
	for i := 0; i < len(*traverser); i++ {
		node := g.GetByID((*traverser)[i])
		if node == nil {
			continue
		}

		err := fu(node.ID)
		if err != nil {
			return err
		}

		g.appendRequirements(traverser, node.Requirements)
	}

	return nil
}

func (g *Graph) appendRequirements(traverser *[]BlockID, requirements []BlockID) {
	for _, bID := range requirements {
		if !slices.Contains(*traverser, bID) {
			*traverser = append(*traverser, bID)
		}
	}
}

func (g *Graph) traverseOutputBlocks(traverser *[]BlockID, fu GraphWalkerCB) error {
	for _, node := range *g {
		bt, _ := node.ID.Parse()
		if bt != BlockType_Output {
			continue
		}

		if !g.allDependenciesTraversed(traverser, node.Requirements) {
			continue
		}

		err := fu(node.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Graph) allDependenciesTraversed(traverser *[]BlockID, requirements []BlockID) bool {
	for _, bId := range requirements {
		if !slices.Contains(*traverser, bId) {
			return false
		}
	}

	return true
}
