package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderNotifyContentReplacesValuePlaceholdersInOrder(t *testing.T) {
	content := "{{value}}，当前剩余额度为 {{value}}，充值链接：<a href='{{value}}'>{{value}}</a>"

	rendered := renderNotifyContent(content, []interface{}{
		"您的额度即将用尽",
		"$228.99",
		"https://api.pureapi.net/console/topup",
		"https://api.pureapi.net/console/topup",
	})

	require.Equal(t, "您的额度即将用尽，当前剩余额度为 $228.99，充值链接：<a href='https://api.pureapi.net/console/topup'>https://api.pureapi.net/console/topup</a>", rendered)
	require.NotContains(t, rendered, "%!(EXTRA")
	require.NotContains(t, rendered, "{{value}}")
}

func TestRenderNotifyContentLeavesExtraValuesOutOfContent(t *testing.T) {
	rendered := renderNotifyContent("余额：{{value}}", []interface{}{"$1.00", "extra"})

	require.Equal(t, "余额：$1.00", rendered)
}
