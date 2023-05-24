package txscript

import (
	"bytes"
	"github.com/james-ray/hcd/wire"
	"testing"
)

var (
	// manyInputsBenchTx is a transaction that contains a lot of inputs which is
	// useful for benchmarking signature hash calculation.
	manyInputsBenchTx wire.MsgTx

	// A mock previous output script to use in the signing benchmark.
	prevOutScript = []parsedOpcode{{&(opcodeArray[1]), hexToBytes("a914f5916158e3e2c4551c1796708db8367207ed13bb87")}}
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

func BenchmarkCalcSigHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(manyInputsBenchTx.TxIn); j++ {
			_, err := CalcSignatureHash(prevOutScript, SigHashAll,
				&manyInputsBenchTx, j, nil)
			if err != nil {
				b.Fatalf("failed to calc signature hash: %v", err)
			}
		}
	}
}

// BenchmarkScriptParsing benchmarks how long it takes to parse a very large
// script.
func BenchmarkScriptParsing(b *testing.B) {
	script, err := genComplexScript()
	if err != nil {
		b.Fatalf("failed to create benchmark script: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pops, err := parseScript(script)
		if err != nil {
			b.Fatalf("failed to parse script: %v", err)
		}

		for _, pop := range pops {
			_ = pop.opcode
			_ = pop.data
		}
	}
}

// BenchmarkDisasmString benchmarks how long it takes to disassemble a very
// large script.
func BenchmarkDisasmString(b *testing.B) {
	script, err := genComplexScript()
	if err != nil {
		b.Fatalf("failed to create benchmark script: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DisasmString(script)
		if err != nil {
			b.Fatalf("failed to disasm script: %v", err)
		}
	}
}
