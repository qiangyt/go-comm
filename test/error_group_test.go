package test

import (
	"errors"
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_ErrorGroup_Add_happy(t *testing.T) {
	a := require.New(t)

	g := comm.NewErrorGroup(false)
	a.False(g.HasError())

	g.Add(errors.New("ERROR A"))
	g.Add(errors.New("ERROR B"))
	g.Add(errors.New("ERROR C"))

	a.True(g.HasError())
	a.Equal(`3 errors totally:
error #1 - ERROR A
error #2 - ERROR B
error #3 - ERROR C`, g.Error())
}

func Test_ErrorGroup_Merge_happy(t *testing.T) {
	a := require.New(t)

	g1 := comm.NewErrorGroup(false)
	g1.Add(errors.New("ERROR 1"))

	g2 := comm.NewErrorGroup(false)
	g2.Add(errors.New("ERROR 2"))

	g1.AddAll(g2)

	a.True(g1.HasError())
	a.Equal(`2 errors totally:
error #1 - ERROR 1
error #2 - ERROR 2`, g1.Error())
}
