package byteinterval

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	tr := New[int]()
	t.Run("insert", func(t *testing.T) {
		tr.Insert([]byte(`100`), []byte(`200`), 1)
		tr.Insert([]byte(`200`), []byte(`300`), 2)
		tr.Insert([]byte(`300`), []byte(`400`), 3)
		tr.Insert([]byte(`400`), []byte(`500`), 4)
		less := tr.Insert([]byte(`300`), []byte(`200`), 400)
		assert.Nil(t, less)
	})
	c := tr.Insert([]byte(`350`), []byte(`600`), 11)
	d := tr.Insert([]byte(`010`), []byte(`100`), 11)
	e := tr.Insert([]byte(`020`), []byte(`100`), 11)
	f := tr.Insert([]byte(`700`), []byte(`800`), 7)
	g := tr.Insert([]byte(`701`), []byte(`799`), 8)
	h := tr.Insert([]byte(`702`), []byte(`798`), 9)
	t.Run("find", func(t *testing.T) {
		vals := tr.Find([]byte(`001`))
		require.Len(t, vals, 0)
		vals = tr.Find([]byte(`100`))
		require.Len(t, vals, 1)
		assert.Equal(t, vals[0], 1)
		vals = tr.Find([]byte(`101`))
		require.Len(t, vals, 1)
		assert.Equal(t, vals[0], 1)
		vals = tr.Find([]byte(`200`))
		require.Len(t, vals, 1)
		assert.Equal(t, vals[0], 2)
		vals = tr.Find([]byte(`350`))
		require.Len(t, vals, 2)
		for _, v := range vals {
			assert.True(t, v == 3 || v == 11)
		}
		vals = tr.Find([]byte(`500`))
		require.Len(t, vals, 1)
		assert.Equal(t, vals[0], 11)
		vals = tr.Find([]byte(`600`))
		require.Len(t, vals, 0)
	})
	t.Run("remove", func(t *testing.T) {
		t.Run("once", func(t *testing.T) {
			vals := tr.Find([]byte(`700`))
			require.Len(t, vals, 1)
			vals = tr.Find([]byte(`701`))
			require.Len(t, vals, 2)
			vals = tr.Find([]byte(`702`))
			require.Len(t, vals, 3)
			h.Remove()
			vals = tr.Find([]byte(`702`))
			require.Len(t, vals, 2)
			g.Remove()
			vals = tr.Find([]byte(`702`))
			require.Len(t, vals, 1)
			f.Remove()
			vals = tr.Find([]byte(`702`))
			require.Len(t, vals, 0)
			c.Remove()
			vals = tr.Find([]byte(`350`))
			require.Len(t, vals, 1)
			assert.Equal(t, vals[0], 3)
			vals = tr.Find([]byte(`450`))
			require.Len(t, vals, 1)
			assert.Equal(t, vals[0], 4)
			vals = tr.Find([]byte(`550`))
			require.Len(t, vals, 0)
			vals = tr.Find([]byte(`650`))
			require.Len(t, vals, 0)
			vals = tr.Find([]byte(`000`))
			require.Len(t, vals, 0)
			d.Remove()
			vals = tr.Find([]byte(`010`))
			require.Len(t, vals, 0)
			vals = tr.Find([]byte(`100`))
			require.Len(t, vals, 1)
			vals = tr.Find([]byte(`020`))
			require.Len(t, vals, 1)
			e.Remove()
			vals = tr.Find([]byte(`020`))
			require.Len(t, vals, 0)
			vals = tr.Find([]byte(`100`))
			require.Len(t, vals, 1)
		})
		t.Run("twice", func(t *testing.T) {
			e.Remove()
			d.rem = false
			d.Remove()
		})
		t.Run("duplicate", func(t *testing.T) {
			i := tr.Insert([]byte(`801`), []byte(`802`), 1)
			tr.Insert([]byte(`801`), []byte(`802`), 2)
			i.Remove()
			vals := tr.Find([]byte(`801`))
			require.Len(t, vals, 1)
		})

	})
	t.Run("findAny", func(t *testing.T) {
		vals := tr.FindAny([]byte(`101`), []byte(`102`), []byte(`103`))
		require.Len(t, vals, 1)
		assert.Equal(t, vals[0], 1)
		vals = tr.FindAny([]byte(`100`), []byte(`200`), []byte(`300`))
		require.Len(t, vals, 3)
		vals = tr.FindAny([]byte(`000`), []byte(`100`))
		require.Len(t, vals, 1)
		vals = tr.FindAny([]byte(`000`), []byte(`001`), []byte(`002`))
		require.Len(t, vals, 0)
		vals = tr.FindAny([]byte(`999`))
		require.Len(t, vals, 0)
	})
}

func TestReadme(t *testing.T) {
	tree := New[int]()
	tree.Insert([]byte(`alpha`), []byte(`bravo`), 100)
	tree.Insert([]byte(`bravo`), []byte(`charlie`), 200)
	c := tree.Insert([]byte(`bravo`), []byte(`delta`), 300)

	items := tree.Find([]byte(`bravo`))
	require.Equal(t, items, []int{200, 300})

	items = tree.Find([]byte(`charlie`))
	require.Equal(t, items, []int{300})

	items = tree.FindAny([]byte(`alpha`), []byte(`bravo`), []byte(`charlie`))
	require.Equal(t, items, []int{100, 200, 300})

	c.Remove()

	items = tree.Find([]byte(`bravo`))
	require.Equal(t, items, []int{200})

	items = tree.Find([]byte(`charlie`))
	require.Equal(t, items, []int(nil))
}
