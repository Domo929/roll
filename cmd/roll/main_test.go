package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRollsOutput(t *testing.T) {
	var buf bytes.Buffer
	handleRollsTo(&buf, []string{"1d6"})
	output := buf.String()
	assert.Contains(t, output, "1d6")
}

func TestHandleAdvantageTo(t *testing.T) {
	var buf bytes.Buffer
	err := handleAdvantageTo(&buf, []string{"5"})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Advantage")
}

func TestHandleAdvantageToNoMod(t *testing.T) {
	var buf bytes.Buffer
	err := handleAdvantageTo(&buf, []string{})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Advantage")
}

func TestHandleAdvantageToInvalidMod(t *testing.T) {
	var buf bytes.Buffer
	err := handleAdvantageTo(&buf, []string{"abc"})
	assert.Error(t, err)
}

func TestHandleDisadvantageTo(t *testing.T) {
	var buf bytes.Buffer
	err := handleDisadvantageTo(&buf, []string{})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Disadvantage")
}

func TestHandleStatsTo(t *testing.T) {
	var buf bytes.Buffer
	handleStatsTo(&buf)
	output := buf.String()
	assert.Contains(t, output, "STR")
	assert.Contains(t, output, "CHA")
	assert.Contains(t, output, "Total")
}
