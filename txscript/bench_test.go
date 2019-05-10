package txscript

import (
	"bytes"
	"testing"
)


func genComplexScript() ([]byte, error) {
	var scriptLen int
	builder := NewScriptBuilder()
	for i := 0; i < MaxOpsPerScript/2; i++ {
		builder.AddOp(OP_TRUE)
		scriptLen++
	}
	maxData := bytes.Repeat([]byte{0x02}, MaxScriptElementSize)
	for i := 0; i < (maxScriptSize-scriptLen)/MaxScriptElementSize; i++ {
		builder.AddData(maxData)
	}
	return builder.Script()
}

func BenchmarkIsStakeGenerationScript(b *testing.B) {
	script, err := genComplexScript()
	if err != nil {
		b.Fatalf("failed to create benchmark script: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pops, _ := parseScript(script)
		_ = isStakeGen(pops)
	}
}
