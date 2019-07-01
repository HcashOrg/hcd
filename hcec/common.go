package hcec

// SignatureType defines a specific cryptographic signature and curve pair for
// use in transaction scripts and addresses.

const (
	// STEcdsaSecp256k1 specifies that the signature is an ECDSA signature
	// over the secp256k1 elliptic curve.
	STEcdsaSecp256k1 int = 0

	// STEd25519 specifies that the signature is an ECDSA signature over the
	// edwards25519 twisted Edwards curve.
	STEd25519 = 1

	// STSchnorrSecp256k1 specifies that the signature is a Schnorr
	// signature over the secp256k1 elliptic curve.
	STSchnorrSecp256k1 = 2
)
