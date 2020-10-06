package yaml

import (
	"bytes"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSteps(t *testing.T) {
	var (
		tc = []struct {
			name string
			yaml string
			test func(*testing.T, workflow.Node, workflow.Variables, error)
		}{
			{
				"one setter in map",
				`steps: { setter: { set: { foo: "bar" } }, (FINAL) }`,
				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
					scope, err = workflow.Workflow(context.Background(), start, scope)
					req.NoError(err)
					req.NotNil(scope)
					req.Contains(scope, "foo")
					req.Equal("bar", scope["foo"])
				},
			},
			{
				"one setter in sequence",
				`steps: [ { set: { foo: "bar" }, ref: setter }, (FINAL) ]`,
				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
					scope, err = workflow.Workflow(context.Background(), start, scope)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(scope)
					req.Contains(scope, "foo")
					req.Equal("bar", scope["foo"])
				},
			},
			{
				"simple gateway",
				`
steps:
  gw:
    exclusive-gateway:
    - if: this
      next: { setter: { set: { foo: "bar" } }, (FINAL) }
    - next: { (ERROR): "Poo" }
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
			{
				"two-stepper",
				`
steps:
  setter1: { set: { foo: 1 } }
  setter2: { set: { foo: 2 } }
  (FINAL)
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
			{
				"step-container-step",
				`
steps:
  setter1: { set: { foo: 1 } }
  container: { steps: [ { ref: setter2, set: { foo: 2 } } ] }
  (FINAL)
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
			{
				"nested steps under gw",
				`
steps:
  check:
    exclusive-gateway:
    - if: foo
      next: { setter1: { set: { foo: 1 } }, (FINAL) }
    - next: { setter2: { set: { foo: 2 } }, setter3: { set: { foo: 3 } }, (FINAL) }
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
			{
				"loop",
				`
steps:
  setter1: { set: { foo: foo + 1 } }
  while:
    exclusive-gateway:
    - if: foo < 5
      next: setter1
    - next: (FINAL)
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
			{
				"fork-join",
				`
steps:
  fork:
    fork-gateway:
    - { set: { s1: 1 }, ref: setter1, next: join }
    - { set: { s2: 2 }, ref: setter2, next: join }
    - { set: { s3: 3 }, ref: setter3, next: join }
  join:
    join-gateway:
    - setter1
    - setter2
    - setter3
  (FINAL):
`,

				func(t *testing.T, start workflow.Node, scope workflow.Variables, err error) {
					req := require.New(t)
					req.NoError(err)
					req.NotNil(scope)
					req.NotNil(start)
				},
			},
		}
	)

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			start, scope, err := Load(bytes.NewBufferString(c.yaml))
			c.test(t, start, scope, err)
		})
	}

}
