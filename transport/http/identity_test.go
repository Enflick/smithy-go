package http

import (
	"context"
	smithy "github.com/Enflick/smithy-go"
	"github.com/Enflick/smithy-go/auth"
	"testing"
)

func TestIdentity(t *testing.T) {
	var expected auth.Identity = &auth.AnonymousIdentity{}

	resolver := auth.AnonymousIdentityResolver{}
	actual, _ := resolver.GetIdentity(context.TODO(), smithy.Properties{})
	if expected != actual {
		t.Errorf("Anonymous identity resolver does not produce correct identity")
	}
}
