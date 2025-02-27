// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package explorer

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9"
	"go.mondoo.com/cnquery/v9/checksums"
	llx "go.mondoo.com/cnquery/v9/llx"
	"go.mondoo.com/cnquery/v9/mqlc"
	"go.mondoo.com/cnquery/v9/mrn"
	"go.mondoo.com/cnquery/v9/types"
	"go.mondoo.com/cnquery/v9/utils/multierr"
	"google.golang.org/protobuf/proto"
)

// RefreshMRN computes a MRN from the UID or validates the existing MRN.
// Both of these need to fit the ownerMRN. It also removes the UID.
func (p *Property) RefreshMRN(ownerMRN string) error {
	nu, err := RefreshMRN(ownerMRN, p.Mrn, MRN_RESOURCE_QUERY, p.Uid)
	if err != nil {
		log.Error().Err(err).Str("owner", ownerMRN).Str("uid", p.Uid).Msg("failed to refresh mrn")
		return multierr.Wrap(err, "failed to refresh mrn for query "+p.Title)
	}

	p.Mrn = nu
	p.Uid = ""
	return nil
}

// Compile a given property and return the bundle.
func (p *Property) Compile(props map[string]*llx.Primitive, schema llx.Schema) (*llx.CodeBundle, error) {
	return mqlc.Compile(p.Mql, props, mqlc.NewConfig(schema, cnquery.DefaultFeatures))
}

// RefreshChecksumAndType by compiling the query and updating the Checksum field
func (p *Property) RefreshChecksumAndType(schema llx.Schema) (*llx.CodeBundle, error) {
	return p.refreshChecksumAndType(schema)
}

func (p *Property) refreshChecksumAndType(schema llx.Schema) (*llx.CodeBundle, error) {
	bundle, err := p.Compile(nil, schema)
	if err != nil {
		return bundle, multierr.Wrap(err, "failed to compile property '"+p.Mql+"'")
	}

	if bundle.GetCodeV2().GetId() == "" {
		return bundle, errors.New("failed to compile query: received empty result values")
	}

	// We think its ok to always use the new code id
	p.CodeId = bundle.CodeV2.Id

	// the compile step also dedents the code
	p.Mql = bundle.Source

	// TODO: record multiple entrypoints and types
	// TODO(jaym): is it possible that the 2 could produce different types
	if entrypoints := bundle.CodeV2.Entrypoints(); len(entrypoints) == 1 {
		ep := entrypoints[0]
		chunk := bundle.CodeV2.Chunk(ep)
		typ := chunk.Type()
		p.Type = string(typ)
	} else {
		p.Type = string(types.Any)
	}

	c := checksums.New.
		Add(p.Mql).
		Add(p.CodeId).
		Add(p.Mrn).
		Add(p.Type).
		Add(p.Context).
		Add(p.Title).Add("v2").
		Add(p.Desc)

	for i := range p.For {
		f := p.For[i]
		c = c.Add(f.Mrn)
	}

	p.Checksum = c.String()

	return bundle, nil
}

func (p *Property) Merge(base *Property) {
	if p.Mql == "" {
		p.Mql = base.Mql
	}
	if p.Type == "" {
		p.Type = base.Type
	}
	if p.Context == "" {
		p.Context = base.Context
	}
	if p.Title == "" {
		p.Title = base.Title
	}
	if p.Desc == "" {
		p.Desc = base.Desc
	}
	if len(p.For) == 0 {
		p.For = base.For
	}
}

type PropsCache struct {
	cache        map[string]*Property
	uidOnlyProps map[string]*Property
}

func NewPropsCache() PropsCache {
	return PropsCache{
		cache:        map[string]*Property{},
		uidOnlyProps: map[string]*Property{},
	}
}

// Add properties, NOT overwriting existing ones (instead we add them as base)
func (c PropsCache) Add(props ...*Property) {
	for i := range props {
		base := props[i]
		if base.Uid != "" && base.Mrn == "" {
			// keep track of properties that were specified by uid only.
			// we will merge them in later if we find a matching mrn
			c.uidOnlyProps[base.Uid] = base
			continue
		}
		// All properties at this point should have a mrn
		if base.Mrn != "" {
			if existingProp, ok := c.cache[base.Mrn]; ok {
				existingProp.Merge(base)
			} else {
				c.cache[base.Mrn] = base
			}
		}
	}
}

// try to Get the mrn, will also return uid-based
// properties if they exist first
func (c PropsCache) Get(ctx context.Context, propMrn string) (*Property, string, error) {
	if res, ok := c.cache[propMrn]; ok {
		name, err := mrn.GetResource(propMrn, MRN_RESOURCE_QUERY)
		if err != nil {
			return nil, "", errors.New("failed to get property name")
		}
		if uidProp, ok := c.uidOnlyProps[name]; ok {
			// We have a property that was specified by uid only. We need to merge it in
			// to get the full property.
			p := proto.Clone(uidProp).(*Property)
			p.Merge(res)
			return p, name, nil
		} else {
			// Everything was specified by mrn
			return res, name, nil
		}
	}

	// We currently don't grab properties from upstream. This requires further investigation.
	return nil, "", errors.New("property " + propMrn + " not found")
}
