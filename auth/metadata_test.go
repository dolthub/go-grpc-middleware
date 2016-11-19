package grpc_auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestAuthFromMD(t *testing.T) {
	for _, run := range []struct {
		md      metadata.MD
		value   string
		errCode codes.Code
		msg     string
	}{
		{
			md:    metadata.Pairs(":authorization", "bearer some_token"),
			value: "some_token",
			msg:   "must extract simple bearer tokens without case checking",
		},
		{
			md:    metadata.Pairs(":authorization", "Bearer some_token"),
			value: "some_token",
			msg:   "must extract simple bearer tokens with case checking",
		},
		{
			md:    metadata.Pairs(":authorization", "Bearer some multi string bearer"),
			value: "some multi string bearer",
			msg:   "must handle string based bearers",
		},
		{
			md:      metadata.Pairs(":authorization", "Basic login:passwd"),
			value:   "",
			errCode: codes.Unauthenticated,
			msg:     "must check authentication type",
		},
		{
			md:      metadata.Pairs(":authorization", "Basic login:passwd", ":authorization", "bearer some_token"),
			value:   "",
			errCode: codes.Unauthenticated,
			msg:     "must not allow multiple authentication methods",
		},
	} {
		ctx := metadata.NewContext(context.TODO(), run.md)
		out, err := AuthFromMD(ctx, "bearer")
		if run.errCode != codes.OK {
			assert.Equal(t, run.errCode, grpc.Code(err), run.msg)
		} else {
			assert.NoError(t, err, run.msg)
		}
		assert.Equal(t, run.value, out, run.msg)
	}

}
