package comm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewErrorGroup(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	a.NotNil(eg)
	a.False(eg.HasError())
	a.Equal(0, eg.AmountOfErrors())
}

func TestErrorGroup_Add(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)

	// Add first error
	err1 := errors.New("error 1")
	eg.Add(err1)
	a.True(eg.HasError())
	a.Equal(1, eg.AmountOfErrors())

	// Add second error
	err2 := errors.New("error 2")
	eg.Add(err2)
	a.Equal(2, eg.AmountOfErrors())
}

func TestErrorGroup_AddNil(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)

	// Adding nil should not increase count
	eg.Add(nil)
	a.False(eg.HasError())
	a.Equal(0, eg.AmountOfErrors())
}

func TestErrorGroup_HasError(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	a.False(eg.HasError())

	eg.Add(errors.New("test error"))
	a.True(eg.HasError())
}

func TestErrorGroup_AmountOfErrors(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	a.Equal(0, eg.AmountOfErrors())

	eg.Add(errors.New("error 1"))
	a.Equal(1, eg.AmountOfErrors())

	eg.Add(errors.New("error 2"))
	a.Equal(2, eg.AmountOfErrors())

	eg.Add(errors.New("error 3"))
	a.Equal(3, eg.AmountOfErrors())
}

func TestErrorGroup_Error(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	eg.Add(errors.New("first error"))
	eg.Add(errors.New("second error"))

	errorMsg := eg.Error()
	a.Contains(errorMsg, "2 errors totally")
	a.Contains(errorMsg, "error #1")
	a.Contains(errorMsg, "first error")
	a.Contains(errorMsg, "error #2")
	a.Contains(errorMsg, "second error")
}

func TestErrorGroup_Error_withDumpStack(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(true)
	eg.Add(errors.New("error 1"))
	eg.Add(errors.New("error 2"))

	errorMsg := eg.Error()
	a.Contains(errorMsg, "2 errors totally")
}

func TestErrorGroup_AddAll(t *testing.T) {
	a := require.New(t)

	eg1 := NewErrorGroup(false)
	eg1.Add(errors.New("error 1"))
	eg1.Add(errors.New("error 2"))

	eg2 := NewErrorGroup(false)
	eg2.Add(errors.New("error 3"))

	eg1.AddAll(eg2)
	a.Equal(3, eg1.AmountOfErrors())
}

func TestErrorGroup_AddAll_nil(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	eg.Add(errors.New("error 1"))

	eg.AddAll(nil)
	a.Equal(1, eg.AmountOfErrors())
}

func TestErrorGroup_AddAll_empty(t *testing.T) {
	a := require.New(t)

	eg1 := NewErrorGroup(false)
	eg1.Add(errors.New("error 1"))

	eg2 := NewErrorGroup(false)

	eg1.AddAll(eg2)
	a.Equal(1, eg1.AmountOfErrors())
}

func TestErrorGroup_MayError_withErrors(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)
	eg.Add(errors.New("test error"))

	err := eg.MayError()
	a.NotNil(err)
	a.Error(err)
}

func TestErrorGroup_MayError_noErrors(t *testing.T) {
	a := require.New(t)

	eg := NewErrorGroup(false)

	err := eg.MayError()
	a.Nil(err)
}
