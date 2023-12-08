package trade

import (
	"context"
	"hft-assistant-clicker/setting"
	"testing"
	"time"
)

func TestMouseControl(t *testing.T) {
	s := setting.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := New(ctx, s)
	c.Setting.IsMoveObject = true

	c.MouseControl()

	t.Log("success done")
}
