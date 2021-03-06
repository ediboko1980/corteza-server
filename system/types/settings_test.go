package types

import (
	"sort"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/pkg/settings"
)

// 	Hello! This file is auto-generated.

func Test_settingsExtAuthProvidersDecode(t *testing.T) {
	type (
		Dst struct {
			Providers ExternalAuthProviderSet
		}
	)

	var (
		aux = Dst{}
		kv  = settings.KV{
			"providers.foo.enabled":                types.JSONText(`true`),
			"providers.openid-connect.bar.enabled": types.JSONText(`true`),
			"providers.openid-connect.bar.key":     types.JSONText(`"K3Y"`),
			"providers.google.enabled":             types.JSONText(`true`),
			"providers.google.key":                 types.JSONText(`"g00gl3"`),
		}

		eq = Dst{
			Providers: ExternalAuthProviderSet{
				{Handle: "github"},
				{Handle: "facebook"},
				{Enabled: true, Key: "g00gl3", Handle: "google"},
				{Handle: "linkedin"},
				{Enabled: true, Key: "K3Y", Handle: "openid-connect.bar"},
			},
		}
	)

	sort.Sort(eq.Providers)

	require.NoError(t, settings.DecodeKV(kv, &aux))
	require.Len(t, aux.Providers, 5)

	require.Nil(t,
		aux.Providers.FindByHandle("foo"))

	require.Equal(t,
		aux.Providers.FindByHandle("openid-connect.bar"),
		&ExternalAuthProvider{Enabled: true, Key: "K3Y", Handle: "openid-connect.bar", Label: "Bar"})

	require.Equal(t,
		aux.Providers.FindByHandle("google"),
		&ExternalAuthProvider{Enabled: true, Key: "g00gl3", Handle: "google", Label: "Google"})

	require.Equal(t,
		aux.Providers.FindByHandle("linkedin"),
		&ExternalAuthProvider{Handle: "linkedin", Label: "LinkedIn"})

	require.Equal(t,
		aux.Providers.FindByHandle("github"),
		&ExternalAuthProvider{Handle: "github", Label: "GitHub"})

	require.Equal(t,
		aux.Providers.FindByHandle("facebook"),
		&ExternalAuthProvider{Handle: "facebook", Label: "Facebook"})

}
