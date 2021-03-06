package event

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackFrame(t *testing.T) {
	pc, file, line, ok := runtime.Caller(0) //foo
	require.True(t, ok)
	f := runtime.FuncForPC(pc)
	require.NotNil(t, f)

	frame := newStackFrame(pc)
	require.NotNil(t, frame)

	assert.Equal(t, pc, frame.PC)
	assert.Equal(t, file, frame.File)
	assert.Equal(t, line, frame.Line)
	assert.Equal(t, f.Name(), frame.Function)
	assert.Equal(t, "event", frame.Package)
	assert.Equal(t, "TestStackFrame", frame.Func)

	assert.Equal(t, "	pc, file, line, ok := runtime.Caller(0) //foo", string(frame.Source()))

	scLinesBefore, scLine, scLinesAfter := frame.SourceContext(2, 2)
	assert.Equal(t, [][]byte{[]byte{}, []byte("func TestStackFrame(t *testing.T) {")}, scLinesBefore)
	assert.Equal(t, []byte("	pc, file, line, ok := runtime.Caller(0) //foo"), scLine)
	assert.Equal(t, [][]byte{[]byte("	require.True(t, ok)"), []byte("	f := runtime.FuncForPC(pc)")}, scLinesAfter)
}

func TestLevelName(t *testing.T) {
	assert.Equal(t, "debug", Debug.String())
	assert.Equal(t, "info", Info.String())
	assert.Equal(t, "notice", Notice.String())
	assert.Equal(t, "warning", Warning.String())
	assert.Equal(t, "error", Error.String())
	assert.Equal(t, "critical", Critical.String())
	assert.Equal(t, "emergency", Emergency.String())
}

func TestLevelInt(t *testing.T) {
	assert.IsType(t, int(0), Info.Int())
	assert.Equal(t, int(Info), Info.Int())
}

func TestNew(t *testing.T) {
	e := New(
		123,
		Notice,
		"test New",
		map[string]interface{}{"foo": map[string]string{"bar": "baz"}},
		false,
	)

	assert.Equal(t, uint64(123), e.Id)
	assert.Equal(t, Notice, e.Level)
	assert.Equal(t, "test New", e.Message)
	assert.Equal(t, map[string]interface{}{"foo.bar": "baz"}, e.FlatFields)
	assert.Equal(t, map[interface{}]interface{}{"foo": map[interface{}]interface{}{"bar": "baz"}}, e.Fields)
	assert.Empty(t, e.Stack)
	assert.WithinDuration(t, time.Now(), e.Time, time.Second)
}

func TestNew_Stack(t *testing.T) {
	// Temporarily override repoPath as New is going to look for the first
	// file outside the package, which won't work here...
	repoPathBkup := RepoPath
	RepoPath = FilePath
	defer func() { RepoPath = repoPathBkup }()

	e := New(
		123,
		Notice,
		"test New Stack",
		nil,
		true,
	)

	_, file, _, _ := runtime.Caller(0)

	frame1 := e.Stack[0]
	require.NotNil(t, frame1)
	assert.Equal(t, file, frame1.File)
}
