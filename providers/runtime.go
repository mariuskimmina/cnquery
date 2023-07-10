package providers

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/llx"
	"go.mondoo.com/cnquery/providers/proto"
	"go.mondoo.com/cnquery/resources"
	"go.mondoo.com/cnquery/types"
	"go.mondoo.com/ranger-rpc"
)

type Runtime struct {
	coordinator    *coordinator
	Provider       *RunningProvider
	Connection     *proto.Connection
	SchemaData     *resources.Schema
	UpstreamConfig *UpstreamConfig
	Recording      Recording

	isClosed bool
}

// mondoo platform config so that resource scan talk upstream
// TODO: this configuration struct does not belong into the MQL package
// nevertheless the MQL runtime needs to have something that allows users
// to store additional credentials so that resource can use those for
// their resources.
type UpstreamConfig struct {
	AssetMrn    string
	SpaceMrn    string
	ApiEndpoint string
	Plugins     []ranger.ClientPlugin
	Incognito   bool
	HttpClient  *http.Client
}

func (c *coordinator) NewRuntime() *Runtime {
	return &Runtime{
		coordinator: c,
	}
}

func (r *Runtime) Close() {
	if r.isClosed {
		return
	}
	r.isClosed = true

	if err := r.Recording.Save(); err != nil {
		log.Error().Err(err).Msg("failed to save recording")
	}

	r.coordinator.Close(r.Provider)
	r.SchemaData = nil
}

// UseProvider sets the main provider for this runtime.
func (r *Runtime) UseProvider(name string) error {
	var running *RunningProvider
	for _, p := range r.coordinator.Running {
		if p.Name == name {
			running = p
			break
		}
	}

	if running == nil {
		var err error
		running, err = r.coordinator.Start(name)
		if err != nil {
			return err
		}
	}

	r.Provider = running
	r.SchemaData = running.Schema

	return nil
}

// Connect to an asset using the main provider
func (r *Runtime) Connect(req *proto.ConnectReq) error {
	if r.Provider == nil {
		return errors.New("cannot connect, please select a provider first")
	}

	if req.Asset == nil || req.Asset.Spec == nil || len(req.Asset.Spec.Assets) == 0 {
		return errors.New("cannot connect, no asset info provided")
	}

	asset := req.Asset.Spec.Assets[0]
	if len(asset.Connections) == 0 {
		return errors.New("cannot connect to asset, no connection info provided")
	}

	conn := asset.Connections[0]
	if conn.Id != 0 {
		r.Connection = &proto.Connection{Id: asset.Connections[0].Id}
		r.Recording.EnsureAsset(asset, r.Provider.Name, conn)
		return nil
	}

	var err error
	r.Connection, err = r.Provider.Plugin.Connect(req)
	if err != nil {
		return err
	}

	r.Recording.EnsureAsset(asset, r.Provider.Name, conn)
	return nil
}

func (r *Runtime) CreateResource(name string, args map[string]*llx.Primitive) (llx.Resource, error) {
	res, err := r.Provider.Plugin.GetData(&proto.DataReq{
		Connection: r.Connection.Id,
		Resource:   name,
		Args:       args,
	}, nil)
	if err != nil {
		return nil, err
	}

	if cached, ok := r.Recording.GetResource(r.Connection.Id, name, string(res.Data.Value)); ok {
		fields, err := RawDataArgsToPrimitiveArgs(cached)
		if err != nil {
			log.Error().Str("resource", name).Str("id", string(res.Data.Value)).Msg("failed to load resource from recording")
		} else {
			r.Provider.Plugin.StoreData(&proto.StoreReq{
				Connection: r.Connection.Id,
				Resources: []*proto.Resource{{
					Name:   name,
					Id:     string(res.Data.Value),
					Fields: fields,
				}},
			})
		}
	} else {
		r.Recording.AddData(r.Connection.Id, name, string(res.Data.Value), "", nil)
	}

	typ := types.Type(res.Data.Type)
	return &llx.MockResource{Name: typ.ResourceName(), ID: string(res.Data.Value)}, nil
}

func (r *Runtime) CreateResourceWithID(name string, id string, args map[string]*llx.Primitive) (llx.Resource, error) {
	panic("NOT YET")
}

func (r *Runtime) Resource(name string) (*resources.ResourceInfo, bool) {
	x, ok := r.SchemaData.Resources[name]
	return x, ok
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
	name := resource.MqlName()
	id := resource.MqlID()
	info, ok := r.SchemaData.Resources[name]
	if !ok {
		return errors.New("cannot get resource info on " + name)
	}
	if _, ok := info.Fields[field]; !ok {
		return errors.New("cannot get field '" + field + "' for resource '" + name + "'")
	}

	if cached, ok := r.Recording.GetData(r.Connection.Id, name, id, field); ok {
		callback(cached.Value, cached.Error)
		return nil
	}

	data, err := r.Provider.Plugin.GetData(&proto.DataReq{
		Connection: r.Connection.Id,
		Resource:   name,
		ResourceId: id,
		Field:      field,
	}, nil)
	if err != nil {
		return err
	}

	if data.Error != "" {
		err = errors.New(data.Error)
	}
	raw := data.Data.RawData()

	r.Recording.AddData(r.Connection.Id, name, id, field, raw)

	callback(raw.Value, err)
	return nil
}

func (r *Runtime) Schema() *resources.Schema {
	return r.SchemaData
}
