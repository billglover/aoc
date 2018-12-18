package device

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExample(t *testing.T) {
	possibleIdent := map[string]bool{}

	before := [4]int{3, 2, 1, 1}
	after := [4]int{3, 2, 2, 1}
	ins := [4]int{9, 2, 1, 2}
	want := map[string]bool{"addi": true, "mulr": true, "seti": true}

	insSet := LoadInstructions()

	for _, oc := range insSet {
		Reg = before
		oc.Do(ins[1], ins[2], ins[3])

		if Reg == after {
			possibleIdent[oc.Name] = true
		}
	}

	if reflect.DeepEqual(possibleIdent, want) == false {
		t.Errorf("unexpected opcodes matched:\ngot:  %v\nwant: %v", possibleIdent, want)
	}
}

func TestPartOne(t *testing.T) {
	samples, err := ReadSamples("input_p1.txt")
	if err != nil {
		t.Fatalf("unable to read samples: %v", err)
	}

	if len(samples) != 776 {
		t.Errorf("unexpected number of samples in input: got %d, want %d", len(samples), 776)
	}

	count := 0

	for _, s := range samples {
		matchedOpcodes := 0

		insSet := LoadInstructions()

		for _, oc := range insSet {
			Reg = s.Before
			oc.Do(s.Instruction[1], s.Instruction[2], s.Instruction[3])

			if Reg == s.After {
				matchedOpcodes++
			}
		}

		if matchedOpcodes >= 3 {
			count++
		}
	}

	fmt.Println("Samples matching 3 or more opcodes:", count)
}
