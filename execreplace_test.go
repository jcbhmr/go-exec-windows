package execreplace_test

import (
	"context"
	"os"
	"testing"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func TestAll(t *testing.T) {
	ctx := context.Background()
	engine := &script.Engine{
		Conds: scripttest.DefaultConds(),
		Cmds:  scripttest.DefaultCmds(),
		Quiet: !testing.Verbose(),
	}
	env := os.Environ()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "PROJECT="+wd)
	scripttest.Test(t, ctx, engine, env, "testdata/*.txt")
}
