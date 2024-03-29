package search

import (
	"APIServerExercise/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// region AddToIndex

func TestSearcher_AddToIndex(t *testing.T) {
	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: false,
	}

	// Adding to empty index
	id1 := uuid.New()
	data1 := struct {
		Key string
	}{
		Key: "value",
	}
	searcher.AddToIndex(&data1, id1, "")

	assert.Len(t, searcher.Index["key"], 1)
	assert.Len(t, searcher.Index["key"]["value"], 1)
	assert.True(t, searcher.Index["key"]["value"][id1])

	// Adding with same Key, different value
	id2 := uuid.New()
	data2 := struct {
		Key string
	}{
		Key: "value2",
	}
	searcher.AddToIndex(&data2, id2, "")

	assert.Len(t, searcher.Index["key"], 2)
	assert.Len(t, searcher.Index["key"]["value2"], 1)
	assert.True(t, searcher.Index["key"]["value2"][id2])

	// Adding with different Key, different value
	id3 := uuid.New()
	data3 := struct {
		Key3 string
	}{
		Key3: "value3",
	}
	searcher.AddToIndex(&data3, id3, "")

	assert.Len(t, searcher.Index["key"], 2)
	assert.Len(t, searcher.Index["key3"], 1)
	assert.Len(t, searcher.Index["key"]["value"], 1)
	assert.Len(t, searcher.Index["key"]["value2"], 1)
	assert.Len(t, searcher.Index["key3"]["value3"], 1)
	assert.True(t, searcher.Index["key3"]["value3"][id3])
}

func TestSearcher_AddToIndex_WithSlice(t *testing.T) {
	type childData struct {
		Key1 string
		Key2 string
	}

	type parentData struct {
		ChildDatas []*childData
	}

	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: false,
	}

	id := uuid.New()
	data := parentData{
		ChildDatas: []*childData{
			{Key1: "value1", Key2: "value2"},
			{Key1: "value3", Key2: "value4"},
		},
	}
	searcher.AddToIndex(&data, id, "")

	assert.Len(t, searcher.Index["childdatas.key1"], 2)
	assert.Len(t, searcher.Index["childdatas.key1"]["value1"], 1)
	assert.True(t, searcher.Index["childdatas.key1"]["value1"][id])
	assert.Len(t, searcher.Index["childdatas.key1"]["value3"], 1)
	assert.True(t, searcher.Index["childdatas.key1"]["value3"][id])
	assert.Len(t, searcher.Index["childdatas.key2"], 2)
	assert.Len(t, searcher.Index["childdatas.key2"]["value2"], 1)
	assert.True(t, searcher.Index["childdatas.key2"]["value2"][id])
	assert.Len(t, searcher.Index["childdatas.key2"]["value4"], 1)
	assert.True(t, searcher.Index["childdatas.key2"]["value4"][id])
}

func TestSearcher_AddToIndex_WithDisableIndexWords(t *testing.T) {
	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: true,
	}

	id1 := uuid.New()
	data1 := struct {
		Key string
	}{
		Key: "value1 value2 value3",
	}
	searcher.AddToIndex(&data1, id1, "")

	assert.Len(t, searcher.Index["key"], 1)
	assert.Len(t, searcher.Index["key"]["value1 value2 value3"], 1)
	assert.True(t, searcher.Index["key"]["value1 value2 value3"][id1])
}

func TestSearcher_AddToIndex_WithoutDisableIndexWords(t *testing.T) {
	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: false,
	}

	id1 := uuid.New()
	data1 := struct {
		Key string
	}{
		Key: "value1 value2\nvalue3\tvalue4\n   value5",
	}
	searcher.AddToIndex(&data1, id1, "")

	assert.Len(t, searcher.Index["key"], 6)
	assert.Len(t, searcher.Index["key"]["value1 value2\nvalue3\tvalue4\n   value5"], 1)
	assert.True(t, searcher.Index["key"]["value1 value2\nvalue3\tvalue4\n   value5"][id1])
	assert.Len(t, searcher.Index["key"]["value1"], 1)
	assert.True(t, searcher.Index["key"]["value1"][id1])
	assert.Len(t, searcher.Index["key"]["value2"], 1)
	assert.True(t, searcher.Index["key"]["value2"][id1])
	assert.Len(t, searcher.Index["key"]["value3"], 1)
	assert.True(t, searcher.Index["key"]["value3"][id1])
	assert.Len(t, searcher.Index["key"]["value4"], 1)
	assert.True(t, searcher.Index["key"]["value4"][id1])
	assert.Len(t, searcher.Index["key"]["value5"], 1)
	assert.True(t, searcher.Index["key"]["value5"][id1])
}

// endregion

// region RemoveFromIndex

func TestSearcher_RemoveFromIndex(t *testing.T) {
	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: false,
	}

	id1 := uuid.New()
	data1 := struct {
		Key string
	}{
		Key: "value",
	}
	searcher.AddToIndex(&data1, id1, "")

	id2 := uuid.New()
	data2 := struct {
		Key string
	}{
		Key: "value",
	}
	searcher.AddToIndex(&data2, id2, "")

	searcher.RemoveFromIndex(id1)
	assert.Len(t, searcher.Index["key"], 1)
	assert.Len(t, searcher.Index["key"]["value"], 1)
	assert.True(t, searcher.Index["key"]["value"][id2])
}

func TestSearcher_RemoveFromIndex_RemoveValueIndex(t *testing.T) {
	searcher := Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: false,
	}

	id1 := uuid.New()
	data1 := struct {
		Key string
	}{
		Key: "value",
	}
	searcher.AddToIndex(&data1, id1, "")
	searcher.RemoveFromIndex(id1)
	assert.Empty(t, searcher.Index["key"])
}

// endregion

// region FilterMetadata

func TestSearcher_FilterMetadata(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	id4 := uuid.New()

	searcher := Searcher{
		Index: map[string]map[string]map[uuid.UUID]bool{
			"key": {
				"value1": {id1: true, id4: true, id3: true},
				"value2": {id2: true},
			},
		},
	}

	database := &core.Database{
		Metadatas: map[uuid.UUID]*core.Metadata{
			id1: {Id: id1},
			id2: {Id: id2},
			id4: {Id: id4},
			id3: {Id: id3},
		},
		Ordering: []uuid.UUID{
			id1,
			id2,
			id3,
			id4,
		},
	}

	query := map[string][]string{
		"key": {"value1"},
	}

	results, err := searcher.FilterMetadata(query, database)
	assert.Nil(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, id1, results[0].Id)
	assert.Equal(t, id3, results[1].Id)
	assert.Equal(t, id4, results[2].Id)
}

func TestSearcher_FilterMetadata_WithInvalidKey(t *testing.T) {
	id1 := uuid.New()

	searcher := Searcher{
		Index: map[string]map[string]map[uuid.UUID]bool{
			"key": {
				"value1": {id1: true},
			},
		},
	}

	database := &core.Database{
		Metadatas: map[uuid.UUID]*core.Metadata{
			id1: {Id: id1},
		},
		Ordering: []uuid.UUID{
			id1,
		},
	}

	query := map[string][]string{
		"invalidkey": {"value1"},
	}

	results, err := searcher.FilterMetadata(query, database)
	assert.Nil(t, results)
	assert.Error(t, err)
	assert.Equal(t, "no such field name invalidkey", err.Error())
}

// endregion

// region cleanWork

func TestSearcher_CleanWord(t *testing.T) {
	testString := "testString"
	result, include := cleanWord(testString)
	assert.True(t, include)
	assert.Equal(t, testString, result)
}

func TestSearcher_CleanWord_EmptyString(t *testing.T) {
	testString := ""
	_, include := cleanWord(testString)
	assert.False(t, include)
}

func TestSearcher_CleanWord_WithSpecialCharacters(t *testing.T) {
	testString := "**testString**"
	result, include := cleanWord(testString)
	assert.True(t, include)
	assert.Equal(t, "testString", result)
}

func TestSearcher_CleanWord_WithSpecialCharactersEmptyString(t *testing.T) {
	testString := "###"
	_, include := cleanWord(testString)
	assert.False(t, include)
}

// endregion
