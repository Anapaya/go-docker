package container

import (
	"context"
	"strings"
	"testing"

	"github.com/cpuguy83/go-docker/errdefs"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestKill(t *testing.T) {
	ctx := context.Background()
	s := newTestService(t)

	err := s.Kill(ctx, "notexist")
	assert.Check(t, errdefs.IsNotFound(err), err)

	c, err := s.Create(ctx, "busybox:latest", WithCreateName(strings.ToLower(t.Name())), WithCreateTTY, WithCreateCmd("/bin/sh", "-c", "trap 'exit 0' SIGTERM; while true; do usleep 100000; done"))
	defer func() {
		assert.Check(t, s.Remove(ctx, c.ID(), WithRemoveForce))
	}()
	assert.NilError(t, err)
	assert.NilError(t, c.Start(ctx))

	err = c.Kill(ctx, WithKillSignal("FAKESIG"))
	assert.Check(t, errdefs.IsInvalid(err), err)

	err = c.Kill(ctx, WithKillSignal("SIGUSR1"))
	assert.NilError(t, err)

	inspect, err := c.Inspect(ctx)
	assert.NilError(t, err)
	assert.Check(t, cmp.Equal(inspect.State.Running, true))

	err = c.Kill(ctx, WithKillSignal("SIGKILL"))
	assert.NilError(t, err)
	inspect, err = c.Inspect(ctx)
	assert.NilError(t, err)
	assert.Check(t, cmp.Equal(inspect.State.Running, false))

	assert.NilError(t, c.Start(ctx))
	inspect, err = c.Inspect(ctx)
	assert.NilError(t, err)
	assert.Check(t, cmp.Equal(inspect.State.Running, true))

	err = c.Kill(ctx)
	assert.NilError(t, err)
	inspect, err = c.Inspect(ctx)
	assert.NilError(t, err)
	assert.Check(t, cmp.Equal(inspect.State.Running, false))
}
