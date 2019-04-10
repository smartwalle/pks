package pks

import (
	_ "github.com/micro/go-plugins/client/grpc"
	_ "github.com/micro/go-plugins/server/grpc"

	_ "github.com/smartwalle/pks/plugins/client/pks_grpc"
	_ "github.com/smartwalle/pks/plugins/server/pks_grpc"

	_ "github.com/smartwalle/pks/plugins/registry/pks_etcdv3"
)
