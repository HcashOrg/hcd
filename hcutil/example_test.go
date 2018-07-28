package hcutil_test

import (
	"fmt"
	"math"

	"github.com/HcashOrg/hcd/hcutil"
)

func ExampleAmount() {

	a := hcutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = hcutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = hcutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 HC
	// 100,000,000 Atoms: 1 HC
	// 100,000 Atoms: 0.001 HC
}

func ExampleNewAmount() {
	amountOne, err := hcutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := hcutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := hcutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := hcutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 HC
	// 0.01234567 HC
	// 0 HC
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := hcutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(hcutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(hcutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(hcutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(hcutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kHC
	// Atom to Coin: 444333.222111 HC
	// Atom to MilliCoin: 444333222.111 mHC
	// Atom to MicroCoin: 444333222111 μHC
	// Atom to Atom: 44433322211100 Atom
}
