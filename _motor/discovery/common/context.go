// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package common

import "context"

type ContextInitializer interface {
	InitCtx(ctx context.Context) context.Context
}
