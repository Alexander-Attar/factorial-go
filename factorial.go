package factorial

import "math/big"

// GenerateFactorial calculates the factorial for the requested n.
func GenerateFactorial(n *big.Int) *big.Int {
	var (
		a   = big.NewInt(1)
		one = big.NewInt(1)
	)

	// Handle cases 0 and 1
	if n.Cmp(one) <= 0 {
		return a
	}

	// Handle all other cases
	for i := big.NewInt(2); i.Cmp(n) <= 0; i.Add(i, one) {
		a = a.Mul(a, i)
	}
	return a
}
