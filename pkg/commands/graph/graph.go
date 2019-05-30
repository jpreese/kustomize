// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package graph

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/ifc/transformer"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/pgmconfig"
	"sigs.k8s.io/kustomize/pkg/plugins"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/target"
)

// Options contain the options for creating a graph
type Options struct {
	kustomizationPath string
	outputPath        string
	loadRestrictor    loader.LoadRestrictorFunc
}

// NewOptions creates an Options object for the graph command
func NewOptions(p, o string) *Options {
	return &Options{
		kustomizationPath: p,
		outputPath:        o,
		loadRestrictor:    loader.RestrictionRootOnly,
	}
}

// NewCmdGraph creates a new graph command
func NewCmdGraph(
	out io.Writer, fs fs.FileSystem,
	v ifc.Validator, rf *resmap.Factory,
	ptf transformer.Factory) *cobra.Command {
	var o Options

	pluginConfig := plugins.DefaultPluginConfig()
	pl := plugins.NewLoader(pluginConfig, rf)

	cmd := &cobra.Command{
		Use:   "graph [path]",
		Short: "Create a visual graph of kustomize dependencies",
		Example: `
	# Using dot provided by GraphViz
	kustomize graph | dot -Tsvg > graph.svg
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}
			return o.RunGraph(out, v, fs, rf, ptf, pl)
		},
	}
	cmd.Flags().StringVarP(
		&o.outputPath,
		"output", "o", "",
		"If specified, create the graph in path.")

	loader.AddLoadRestrictionsFlag(cmd.Flags())
	plugins.AddEnablePluginsFlag(
		cmd.Flags(), &pluginConfig.Enabled)

	return cmd
}

// Validate validates the arguments of the graph command
func (o *Options) Validate(args []string) (err error) {
	if len(args) > 1 {
		return errors.New("specify one path to " + pgmconfig.KustomizationFileNames[0])
	}

	if len(args) == 0 {
		o.kustomizationPath = loader.CWD
	} else {
		o.kustomizationPath = args[0]
	}

	o.loadRestrictor, err = loader.ValidateLoadRestrictorFlag()
	return
}

// RunGraph runs the graph command
func (o *Options) RunGraph(
	out io.Writer, v ifc.Validator, fSys fs.FileSystem,
	rf *resmap.Factory, ptf transformer.Factory,
	pl *plugins.Loader) error {

	ldr, err := loader.NewLoader(
		o.loadRestrictor, v, o.kustomizationPath, fSys)
	if err != nil {
		return err
	}
	defer ldr.Cleanup()
	kt, err := target.NewKustTarget(ldr, rf, ptf, pl)
	if err != nil {
		return err
	}
	m, err := kt.MakeCustomizedResMap()
	if err != nil {
		return err
	}

	makeGraph(m)

	return nil
}

func makeGraph(r resmap.ResMap) {
}
