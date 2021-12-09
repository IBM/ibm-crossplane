/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
//
// Copyright 2021 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package core

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"time"

	"github.com/alecthomas/kong"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/crossplane/crossplane/internal/controller/apiextensions"
	"github.com/crossplane/crossplane/internal/controller/pkg"
	"github.com/crossplane/crossplane/internal/feature"
	"github.com/crossplane/crossplane/internal/xpkg"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
)

// Command runs the core crossplane controllers
type Command struct {
	Start startCommand `cmd:"" help:"Start Crossplane controllers."`
	Init  initCommand  `cmd:"" help:"Make cluster ready for Crossplane controllers."`
}

// KongVars represent the kong variables associated with the CLI parser
// required for the Registry default variable interpolation.
var KongVars = kong.Vars{
	"default_registry": name.DefaultRegistry,
}

// Run is the no-op method required for kong call tree
// Kong requires each node in the calling path to have associated
// Run method.
func (c *Command) Run() error {
	return nil
}

type startCommand struct {
	// IBM Patch: change POD_NAMESPACE to WATCH_NAMESPACE
	Namespace      string        `short:"n" help:"Namespace used to unpack and run packages." default:"crossplane-system" env:"WATCH_NAMESPACE"`
	CacheDir       string        `short:"c" help:"Directory used for caching package images." default:"/cache" env:"CACHE_DIR"`
	LeaderElection bool          `short:"l" help:"Use leader election for the controller manager." default:"false" env:"LEADER_ELECTION"`
	Registry       string        `short:"r" help:"Default registry used to fetch packages when not specified in tag." default:"${default_registry}" env:"REGISTRY"`
	Sync           time.Duration `short:"s" help:"Controller manager sync period duration such as 300ms, 1.5h or 2h45m" default:"1h"`

	EnableCompositionRevisions bool `group:"Alpha Features:" help:"Enable support for CompositionRevisions."`
}

// Run core Crossplane controllers.
func (c *startCommand) Run(s *runtime.Scheme, log logging.Logger) error { //nolint:gocyclo
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return errors.Wrap(err, "Cannot get config")
	}
	log.Debug("Starting", "sync-period", c.Sync.String())

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           s,
		LeaderElection:   c.LeaderElection,
		LeaderElectionID: "crossplane-leader-election-core",
		SyncPeriod:       &c.Sync,
	})
	if err != nil {
		return errors.Wrap(err, "Cannot create manager")
	}

	f := &feature.Flags{}
	if c.EnableCompositionRevisions {
		f.Enable(feature.FlagEnableAlphaCompositionRevisions)
		log.Info("Alpha feature enabled", "flag", feature.FlagEnableAlphaCompositionRevisions.String())
	}

	if err := apiextensions.Setup(mgr, log, f); err != nil {
		return errors.Wrap(err, "Cannot setup API extension controllers")
	}

	pkgCache := xpkg.NewImageCache(c.CacheDir, afero.NewOsFs())

	if err := pkg.Setup(mgr, log, pkgCache, c.Namespace, c.Registry); err != nil {
		return errors.Wrap(err, "Cannot add packages controllers to manager")
	}

	// IBM Patch: Migration to use Provider.
	// This migration is needed in order to support the user having previous working Kafka or Postgres.
	// New logic could fail for them and in the worst scenario it could recreate the services
	cli := mgr.GetClient()
	crdNames := []string{"postgrescomposites.shim.bedrock.ibm.com", "kafkacomposites.shim.bedrock.ibm.com"}
	crds := &unstructured.UnstructuredList{}
	crds.SetGroupVersionKind(schema.GroupVersionKind{Version: "apiextensions.k8s.io/v1", Kind: "CustomResourceDefinition"})
	if err := cli.List(context.Background(), crds); err != nil {
		return errors.Wrap(err, "Cannot list CRDs for migration")
	}

	// run migration only when required CRDs exist (would fail on integration tests)
out:
	for _, crd := range crds.Items {
		for _, ck := range crdNames {
			if crd.GetName() == ck {
				if err := migrateToProviderLogic(cli); err != nil {
					return errors.Wrap(err, "Migration failed")
				}
				break out
			}
		}
	}

	return errors.Wrap(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

// IBM Patch: Migration to use Provider.
// migrateToProviderLogic: detects if Composites need to be updated
// so that Crossplane will adapt them to Provider logic.
// This function deletes "compositionRef" and "resourceRefs" if needed.
func migrateToProviderLogic(cli client.Client) error {
	ctx := context.Background()
	// List of compositions that were updated from Crossplane to Provider logic.
	ctu := map[string]bool{
		"kafka-iaf.odlm.bedrock.ibm.com":           true,
		"kafka-iaf-skip-user.odlm.bedrock.ibm.com": true,
		"postgres.odlm.bedrock.ibm.com":            true,
	}
	// List of Composite resources
	cl := &unstructured.UnstructuredList{}
	ckinds := []string{"PostgresComposite", "KafkaComposite"}
	for _, ck := range ckinds {
		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(schema.GroupVersionKind{Version: "shim.bedrock.ibm.com/v1alpha1", Kind: ck})
		if err := cli.List(ctx, l); err != nil {
			return errors.Wrap(err, "Cannot list Composites for migration")
		}
		cl.Items = append(cl.Items, l.Items...)
	}
	// Range over all composite resources and check if any of them needs to be updated
OUTER:
	for _, cpt := range cl.Items {
		var cref string
		cref, _ = fieldpath.Pave(cpt.Object).GetString("spec.compositionRef.name")
		// Handle only if specified Composition is indicated
		if ctu[cref] {
			csite := &unstructured.Unstructured{}
			csite.SetGroupVersionKind(schema.GroupVersionKind{Version: "shim.bedrock.ibm.com/v1alpha1", Kind: cpt.GetKind()})
			if err := cli.Get(ctx, types.NamespacedName{Namespace: "", Name: cpt.GetName()}, csite); err != nil {
				return errors.Wrap(err, "Cannot get Composite for migration")
			}
			spec := csite.Object["spec"].(map[string]interface{})
			// Skip if there are no resourceRefs
			if spec["resourceRefs"] == nil {
				continue
			}
			rrs := spec["resourceRefs"].([]interface{})
			for _, rr := range rrs {
				if rr.(map[string]interface{})["kind"] == "Object" {
					// Do not apply migration if there is at least one `Object` kind
					// Then we assume that Provider logic has already been applied here (future cases)
					continue OUTER
				}
			}
			// Clear following sections. Crossplane will again match Composition
			delete(spec, "compositionRef")
			delete(spec, "resourceRefs")
			if err := cli.Update(ctx, csite); err != nil {
				return errors.Wrap(err, "Cannot update Composite for migration")
			}
			fmt.Println("Successfully updated existing Composite for migration: " + csite.GetName())
		}
	}
	return nil
}
