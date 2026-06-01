package relay

import (
	relaycommon "github.com/QuantumNous/fast-api/relay/common"
	"github.com/QuantumNous/fast-api/types"
)

func newAPIErrorFromParamOverride(err error) *types.FastAPIError {
	if fixedErr, ok := relaycommon.AsParamOverrideReturnError(err); ok {
		return relaycommon.FastAPIErrorFromParamOverride(fixedErr)
	}
	return types.NewError(err, types.ErrorCodeChannelParamOverrideInvalid, types.ErrOptionWithSkipRetry())
}
