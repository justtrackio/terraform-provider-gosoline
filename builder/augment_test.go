package builder_test

import (
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
	"github.com/stretchr/testify/require"
)

func provideAppId() builder.AppId {
	return builder.AppId{
		Project:     "prj",
		Environment: "env",
		Family:      "fam",
		Group:       "grp",
		Application: "app",
	}
}

func TestAugment(t *testing.T) {
	appId := provideAppId()

	hostnamePattern := "{scheme}://{project}-{env}-{family}-{group}-{app}-static.{metadata_domain}:{port}"
	additionalReplacements := map[string]string{
		"metadata_domain": "example.com",
		"scheme":          "https",
		"port":            "1337",
	}

	hostname := builder.Augment(hostnamePattern, appId, additionalReplacements)
	require.Equal(t, "https://prj-env-fam-grp-app-static.example.com:1337", hostname)
}
