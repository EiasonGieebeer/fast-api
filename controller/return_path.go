package controller

import (
	"strings"

	"github.com/QuantumNous/fast-api/common"
	"github.com/QuantumNous/fast-api/setting/system_setting"
)

func paymentReturnPath(suffix string) string {
	base := strings.TrimRight(system_setting.ServerAddress, "/")
	return base + common.ThemeAwarePath(suffix)
}
