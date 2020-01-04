package cymidb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_UpdateDataBuf(t *testing.T) {
	dt := NewDataType("blue.gasser/test")
	n := Node{
		datas: []Data{{Type: dt, Data: []byte("testing")}},
	}
	require.Nil(t, n.DataBuf)
	require.NoError(t, n.updateDataBuf())
	require.NotNil(t, n.DataBuf)

	n2 := Node{
		DataBuf: n.DataBuf,
	}
	d, err := n2.GetData(dt)
	require.NoError(t, err)
	require.Equal(t, d, n.datas[0].Data)
}
