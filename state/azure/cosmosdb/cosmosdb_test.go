// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package cosmosdb

import (
	"encoding/json"
	"testing"

	"github.com/dapr/components-contrib/state"
	"github.com/stretchr/testify/assert"
)

type widget struct {
	Color string `json:"color"`
}

func TestCreateCosmosItem(t *testing.T) {
	value := widget{Color: "red"}
	partitionKey := "/partitionKey"
	t.Run("create item for golang struct", func(t *testing.T) {
		req := state.SetRequest{
			Key:   "testKey",
			Value: value,
		}

		item := createUpsertItem(req, partitionKey)
		assert.Equal(t, partitionKey, item.PartitionKey)
		assert.Equal(t, "testKey", item.ID)
		assert.Equal(t, value, item.Value)

		// items need to be marshallable to JSON with encoding/json
		b, err := json.Marshal(item)
		assert.NoError(t, err)

		j := map[string]interface{}{}
		err = json.Unmarshal(b, &j)
		assert.NoError(t, err)

		value := j["value"]
		m, ok := value.(map[string]interface{})
		assert.Truef(t, ok, "value should be a map")

		assert.Equal(t, "red", m["color"])
	})

	t.Run("create item for JSON bytes", func(t *testing.T) {
		// We have special handling for bytes, we assume that it's UTF-8 text containing
		// JSON and special case it to avoid the serializer turning it into a base64 encoded string
		bytes, err := json.Marshal(value)
		assert.NoError(t, err)

		req := state.SetRequest{
			Key:   "testKey",
			Value: bytes,
		}

		item := createUpsertItem(req, partitionKey)
		assert.Equal(t, partitionKey, item.PartitionKey)
		assert.Equal(t, "testKey", item.ID)

		// items need to be marshallable to JSON with encoding/json
		b, err := json.Marshal(item)
		assert.NoError(t, err)

		j := map[string]interface{}{}
		err = json.Unmarshal(b, &j)
		assert.NoError(t, err)

		value := j["value"]
		m, ok := value.(map[string]interface{})
		assert.Truef(t, ok, "value should be a map")

		assert.Equal(t, "red", m["color"])
	})
}
