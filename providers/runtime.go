// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package providers

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/llx"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/resources"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/upstream"
	"go.mondoo.com/cnquery/v9/types"
	"go.mondoo.com/cnquery/v9/utils/multierr"
	"google.golang.org/grpc/status"
)

const defaultShutdownTimeout = time.Duration(time.Second * 120)

// Runtimes are associated with one asset and carry all providers
// and open connections for that asset.
type Runtime struct {
	Provider       *ConnectedProvider
	UpstreamConfig *upstream.UpstreamConfig
	Recording      Recording
	AutoUpdate     UpdateProvidersConfig

	features []byte
	// coordinator is used to grab providers
	coordinator *coordinator
	// providers for with open connections
	providers map[string]*ConnectedProvider
	// schema aggregates all resources executable on this asset
	schema          extensibleSchema
	isClosed        bool
	isEphemeral     bool
	close           sync.Once
	shutdownTimeout time.Duration
}

type ConnectedProvider struct {
	Instance   *RunningProvider
	Connection *plugin.ConnectRes
}

func (c *coordinator) RuntimeWithShutdownTimeout(timeout time.Duration) *Runtime {
	runtime := c.NewRuntime()
	runtime.shutdownTimeout = timeout
	return runtime
}

type shutdownResult struct {
	Response *plugin.ShutdownRes
	Error    error
}

func (r *Runtime) tryShutdown() shutdownResult {
	// Ephemeral runtimes have their primary provider be ephemeral, i.e. non-shared.
	// All other providers are shared and will not be shut down from within the provider.
	if r.isEphemeral {
		err := r.coordinator.Stop(r.Provider.Instance, true)
		return shutdownResult{Error: err}
	}

	return shutdownResult{}
}

func (r *Runtime) Close() {
	r.isClosed = true
	r.close.Do(func() {
		if err := r.Recording.Save(); err != nil {
			log.Error().Err(err).Msg("failed to save recording")
		}

		response := make(chan shutdownResult, 1)
		go func() {
			response <- r.tryShutdown()
		}()
		select {
		case <-time.After(r.shutdownTimeout):
			log.Error().Str("provider", r.Provider.Instance.Name).Msg("timed out shutting down the provider")
		case result := <-response:
			if result.Error != nil {
				log.Error().Err(result.Error).Msg("failed to shutdown the provider")
			}
		}

		// TODO: ideally, we try to close the provider here but only if there are no more assets that need it
		// r.coordinator.Close(r.Provider.Instance)
		r.schema.Close()
	})
}

func (r *Runtime) DeactivateProviderDiscovery() {
	r.schema.allLoaded = true
}

func (r *Runtime) AssetMRN() string {
	if r.Provider != nil && r.Provider.Connection != nil && r.Provider.Connection.Asset != nil {
		return r.Provider.Connection.Asset.Mrn
	}
	return ""
}

// UseProvider sets the main provider for this runtime.
func (r *Runtime) UseProvider(id string) error {
	// We transfer isEphemeral here because:
	// 1. If the runtime is not ephemeral, it behaves as a shared provider by
	// default.
	// 2. If the runtime is ephemeral, we only want the main provider to be
	// ephemeral. All other providers by default are shared.
	// (Note: In the future we plan to have an isolated runtime mode,
	// where even other providers are ephemeral, ie not shared)
	res, err := r.addProvider(id, r.isEphemeral)
	if err != nil {
		return err
	}

	r.Provider = res
	return nil
}

func (r *Runtime) AddConnectedProvider(c *ConnectedProvider) {
	r.providers[c.Instance.ID] = c
	r.schema.Add(c.Instance.Name, c.Instance.Schema)
}

func (r *Runtime) addProvider(id string, isEphemeral bool) (*ConnectedProvider, error) {
	var running *RunningProvider
	var err error
	if isEphemeral {
		running, err = r.coordinator.Start(id, true, r.AutoUpdate)
		if err != nil {
			return nil, err
		}

	} else {
		// TODO: we need to detect only the shared running providers
		running = r.coordinator.RunningByID[id]
		if running == nil {
			var err error
			running, err = r.coordinator.Start(id, false, r.AutoUpdate)
			if err != nil {
				return nil, err
			}
		}
	}

	res := &ConnectedProvider{Instance: running}
	r.AddConnectedProvider(res)

	return res, nil
}

// DetectProvider will try to detect and start the right provider for this
// runtime. Generally recommended when you receive an asset to be scanned,
// but haven't initialized any provider. It will also try to install providers
// if necessary (and enabled)
func (r *Runtime) DetectProvider(asset *inventory.Asset) error {
	if asset == nil {
		return errors.New("please provide an asset to detect the provider")
	}
	if len(asset.Connections) == 0 {
		return errors.New("asset has no connections, can't detect provider")
	}

	var errs multierr.Errors
	for i := range asset.Connections {
		conn := asset.Connections[i]
		if conn.Type == "" {
			log.Warn().Msg("no connection `type` provided in inventory, falling back to deprecated `backend` field")
			conn.Type = inventory.ConnBackendToType(conn.Backend)
		}

		provider, err := EnsureProvider("", conn.Type, true, r.coordinator.Providers)
		if err != nil {
			errs.Add(err)
			continue
		}

		return r.UseProvider(provider.ID)
	}

	return multierr.Wrap(errs.Deduplicate(), "cannot find provider for this asset")
}

// Connect to an asset using the main provider
func (r *Runtime) Connect(req *plugin.ConnectReq) error {
	if r.Provider == nil {
		return errors.New("cannot connect, please select a provider first")
	}

	if req.Asset == nil {
		return errors.New("cannot connect, no asset info provided")
	}

	asset := req.Asset
	if len(asset.Connections) == 0 {
		return errors.New("cannot connect to asset, no connection info provided")
	}

	r.features = req.Features
	callbacks := providerCallbacks{
		runtime: r,
	}

	var err error
	r.Provider.Connection, err = r.Provider.Instance.Plugin.Connect(req, &callbacks)
	if err != nil {
		return err
	}
	r.Recording.EnsureAsset(r.Provider.Connection.Asset, r.Provider.Instance.ID, r.Provider.Connection.Id, asset.Connections[0])
	return nil
}

func (r *Runtime) CreateResource(name string, args map[string]*llx.Primitive) (llx.Resource, error) {
	provider, info, err := r.lookupResourceProvider(name)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("cannot create '" + name + "', no resource info found")
	}
	name = info.Id

	// Resources without providers are bridging resources only. They are static in nature.
	if provider == nil {
		return &llx.MockResource{Name: name}, nil
	}

	if provider.Connection == nil {
		return nil, errors.New("no connection to provider")
	}

	res, err := provider.Instance.Plugin.GetData(&plugin.DataReq{
		Connection: provider.Connection.Id,
		Resource:   name,
		Args:       args,
	})
	if err != nil {
		return nil, err
	}

	if _, ok := r.Recording.GetResource(provider.Connection.Id, name, string(res.Data.Value)); !ok {
		r.Recording.AddData(provider.Connection.Id, name, string(res.Data.Value), "", nil)
	}

	typ := types.Type(res.Data.Type)
	return &llx.MockResource{Name: typ.ResourceName(), ID: string(res.Data.Value)}, nil
}

func (r *Runtime) CloneResource(src llx.Resource, id string, fields []string, args map[string]*llx.Primitive) (llx.Resource, error) {
	name := src.MqlName()
	srcID := src.MqlID()

	provider, _, err := r.lookupResourceProvider(name)
	if err != nil {
		return nil, err
	}

	for i := range fields {
		field := fields[i]
		data, err := provider.Instance.Plugin.GetData(&plugin.DataReq{
			Connection: provider.Connection.Id,
			Resource:   name,
			ResourceId: srcID,
			Field:      field,
		})
		if err != nil {
			return nil, err
		}
		args[field] = data.Data
	}

	args["__id"] = llx.StringPrimitive(id)

	_, err = provider.Instance.Plugin.StoreData(&plugin.StoreReq{
		Connection: provider.Connection.Id,
		Resources: []*plugin.ResourceData{{
			Name:   name,
			Id:     id,
			Fields: PrimitiveArgsToResultArgs(args),
		}},
	})
	if err != nil {
		return nil, err
	}

	return &llx.MockResource{Name: name, ID: id}, nil
}

func (r *Runtime) Unregister(watcherUID string) error {
	// TODO: we don't unregister just yet...
	return nil
}

func fieldUID(resource string, id string, field string) string {
	return resource + "\x00" + id + "\x00" + field
}

// WatchAndUpdate a resource field and call the function if it changes with its current value
func (r *Runtime) WatchAndUpdate(resource llx.Resource, field string, watcherUID string, callback func(res interface{}, err error)) error {
	raw, err := r.watchAndUpdate(resource.MqlName(), resource.MqlID(), field, watcherUID)
	if raw != nil {
		callback(raw.Value, raw.Error)
	}
	return err
}

func (r *Runtime) watchAndUpdate(resource string, resourceID string, field string, watcherUID string) (*llx.RawData, error) {
	provider, info, fieldInfo, err := r.lookupFieldProvider(resource, field)
	if err != nil {
		return nil, err
	}
	if fieldInfo == nil {
		return nil, errors.New("cannot get field '" + field + "' for resource '" + resource + "'")
	}

	if info.Provider != fieldInfo.Provider {
		// technically we don't need to look up the resource provider, since
		// it had to have been called beforehand to get here
		_, err := provider.Instance.Plugin.StoreData(&plugin.StoreReq{
			Connection: provider.Connection.Id,
			Resources: []*plugin.ResourceData{{
				Name: resource,
				Id:   resourceID,
			}},
		})
		if err != nil {
			return nil, multierr.Wrap(err, "failed to create reference resource "+resource+" in provider "+provider.Instance.Name)
		}
	}

	if cached, ok := r.Recording.GetData(provider.Connection.Id, resource, resourceID, field); ok {
		return cached, nil
	}

	data, err := provider.Instance.Plugin.GetData(&plugin.DataReq{
		Connection: provider.Connection.Id,
		Resource:   resource,
		ResourceId: resourceID,
		Field:      field,
	})
	if err != nil {
		// Recoverable errors can continue with the exeuction,
		// they only store errors in the place of actual data.
		// Every other error is thrown up the chain.
		handled, err := r.handlePluginError(err, provider)
		if !handled {
			return nil, err
		}
		data = &plugin.DataRes{Error: err.Error()}
	}

	var raw *llx.RawData
	if data.Error != "" {
		raw = &llx.RawData{Error: errors.New(data.Error)}
	} else {
		raw = data.Data.RawData()
	}

	r.Recording.AddData(provider.Connection.Id, resource, resourceID, field, raw)
	return raw, nil
}

func (r *Runtime) handlePluginError(err error, provider *ConnectedProvider) (bool, error) {
	st, ok := status.FromError(err)
	if !ok {
		return false, err
	}

	switch st.Code() {
	case 14:
		// Error: Unavailable. Happens when the plugin crashes.
		// TODO: try to restart the plugin and reset its connections
		provider.Instance.isClosed = true
		provider.Instance.err = errors.New("the '" + provider.Instance.Name + "' provider crashed")
		return false, provider.Instance.err
	}
	return false, err
}

type providerCallbacks struct {
	recording *assetRecording
	runtime   *Runtime
}

func (p *providerCallbacks) GetRecording(req *plugin.DataReq) (*plugin.ResourceData, error) {
	resource, ok := p.recording.resources[req.Resource+"\x00"+req.ResourceId]
	if !ok {
		return nil, nil
	}

	res := plugin.ResourceData{
		Name:   req.Resource,
		Id:     req.ResourceId,
		Fields: make(map[string]*llx.Result, len(resource.Fields)),
	}
	for k, v := range resource.Fields {
		res.Fields[k] = v.Result()
	}

	return &res, nil
}

func (p *providerCallbacks) GetData(req *plugin.DataReq) (*plugin.DataRes, error) {
	if req.Field == "" {
		res, err := p.runtime.CreateResource(req.Resource, req.Args)
		if err != nil {
			return nil, err
		}

		return &plugin.DataRes{
			Data: &llx.Primitive{
				Type:  string(types.Resource(res.MqlName())),
				Value: []byte(res.MqlID()),
			},
		}, nil
	}

	raw, err := p.runtime.watchAndUpdate(req.Resource, req.ResourceId, req.Field, "")
	if raw == nil {
		return nil, err
	}
	res := raw.Result()
	return &plugin.DataRes{
		Data:  res.Data,
		Error: res.Error,
	}, err
}

func (p *providerCallbacks) Collect(req *plugin.DataRes) error {
	panic("NOT YET IMPLEMENTED")
	return nil
}

func (r *Runtime) SetRecording(recording Recording) error {
	r.Recording = recording
	if r.Provider == nil || r.Provider.Instance == nil {
		log.Warn().Msg("set recording while no provider is set on runtime")
		return nil
	}
	if r.Provider.Instance.ID != mockProvider.ID {
		return nil
	}

	service := r.Provider.Instance.Plugin.(*mockProviderService)
	// TODO: This is problematic, since we don't have multiple instances of
	// the service!!
	service.runtime = r

	return nil
}

func baseRecording(anyRecording Recording) *recording {
	var baseRecording *recording
	switch x := anyRecording.(type) {
	case *recording:
		baseRecording = x
	case *readOnlyRecording:
		baseRecording = x.recording
	}
	return baseRecording
}

// SetMockRecording is only used for test utilities. Please do not use it!
//
// Deprecated: This function may not be necessary anymore, consider removing.
func (r *Runtime) SetMockRecording(anyRecording Recording, providerID string, mockConnection bool) error {
	r.Recording = anyRecording

	baseRecording := baseRecording(anyRecording)
	if baseRecording == nil {
		return nil
	}

	provider, ok := r.providers[providerID]
	if !ok {
		return errors.New("cannot set recording, provider '" + providerID + "' not found")
	}

	assetRecording := &baseRecording.Assets[0]
	asset := assetRecording.Asset.ToInventory()

	if mockConnection {
		// Dom: we may need to retain the original asset ID, not sure yet...
		asset.Id = "mock-asset"
		asset.Connections = []*inventory.Config{{
			Type: "mock",
		}}

		callbacks := providerCallbacks{
			recording: assetRecording,
			runtime:   r,
		}

		res, err := provider.Instance.Plugin.Connect(&plugin.ConnectReq{
			Asset:        asset,
			HasRecording: true,
		}, &callbacks)
		if err != nil {
			return multierr.Wrap(err, "failed to set mock connection for recording")
		}
		provider.Connection = res
	}

	if provider.Connection == nil {
		// Dom: we may need to cancel the entire setup here, may need to be reconsidered...
		log.Warn().Msg("recording cannot determine asset, no connection was set up!")
	} else {
		baseRecording.assets[provider.Connection.Id] = assetRecording
	}

	return nil
}

func (r *Runtime) lookupResourceProvider(resource string) (*ConnectedProvider, *resources.ResourceInfo, error) {
	info := r.schema.Lookup(resource)
	if info == nil {
		return nil, nil, errors.New("cannot find resource '" + resource + "' in schema")
	}

	if info.Provider == "" {
		// This case happens when the resource is only bridging a resource chain,
		// i.e. it is extending in nature (which we only test for the warning).
		if !info.IsExtension {
			log.Warn().Msg("found a resource without a provider: '" + resource + "'")
		}
		return nil, info, nil
	}

	if provider := r.providers[info.Provider]; provider != nil {
		return provider, info, nil
	}

	providerConn := r.Provider.Instance.ID
	if info.Provider != providerConn && info.Provider != "go.mondoo.com/cnquery/v9/providers/core" && info.Provider != "go.mondoo.com/cnquery/v9/providers/network" {
		return nil, nil, errors.New("incorrect provider for asset, not adding")
	}

	res, err := r.addProvider(info.Provider, false)
	if err != nil {
		return nil, nil, multierr.Wrap(err, "failed to start provider '"+info.Provider+"'")
	}

	conn, err := res.Instance.Plugin.Connect(&plugin.ConnectReq{
		Features: r.features,
		Asset:    r.Provider.Connection.Asset,
	}, nil)
	if err != nil {
		return nil, nil, err
	}

	res.Connection = conn

	return res, info, nil
}

func (r *Runtime) lookupFieldProvider(resource string, field string) (*ConnectedProvider, *resources.ResourceInfo, *resources.Field, error) {
	resourceInfo, fieldInfo := r.schema.LookupField(resource, field)
	if resourceInfo == nil {
		return nil, nil, nil, errors.New("cannot find resource '" + resource + "' in schema")
	}
	if fieldInfo == nil {
		return nil, nil, nil, errors.New("cannot find field '" + field + "' in resource '" + resource + "'")
	}

	if provider := r.providers[fieldInfo.Provider]; provider != nil {
		return provider, resourceInfo, fieldInfo, nil
	}

	res, err := r.addProvider(fieldInfo.Provider, false)
	if err != nil {
		return nil, nil, nil, multierr.Wrap(err, "failed to start provider '"+fieldInfo.Provider+"'")
	}

	conn, err := res.Instance.Plugin.Connect(&plugin.ConnectReq{
		Features: r.features,
		Asset:    r.Provider.Connection.Asset,
	}, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	res.Connection = conn

	return res, resourceInfo, fieldInfo, nil
}

func (r *Runtime) Schema() llx.Schema {
	return &r.schema
}

func (r *Runtime) AddSchema(name string, schema *resources.Schema) {
	r.schema.Add(name, schema)
}

func (r *Runtime) asset() *inventory.Asset {
	if r.Provider == nil || r.Provider.Connection == nil {
		return nil
	}
	return r.Provider.Connection.Asset
}
