package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	db "github.com/Horqu/zkp-communicator-backend/cmd/database"
	"github.com/Horqu/zkp-communicator-backend/cmd/database/queries"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func GenerateZkSnarkProof(username string) bn254.G1Affine {

	publicKey, err := queries.GetPublicKeyByUsername(db.GetDBInstance(), username)
	if err != nil {
		log.Println(err)
	}

	challenge, err := StringToG1Affine(publicKey)
	if err != nil {
		log.Println(err)
	}

	return challenge
}

func VerifyZkSnarkProof(check bn254.G1Affine, y bn254.G1Affine) bool {

	if check.Equal(&y) {
		return true
	} else {
		return false
	}

}

func GenerateSchnorrChallenge(publicKey string) (*big.Int, *big.Int) {

	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	y, err := StringToG1Affine(publicKey)
	if err != nil {
		log.Println(err)
	}

	r, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	var R bn254.G1Affine
	R.ScalarMultiplication(&g, r) // R needs to be stored

	h := sha256.New()
	h.Write(R.Marshal())
	h.Write(y.Marshal())
	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, bn254.ID.ScalarField()) // e needs to be stored

	return e, r
}

func GenerateSchnorrProof(publicKey string, privateKey *big.Int, e *big.Int, r *big.Int) *big.Int {

	s := new(big.Int).Mul(e, privateKey)
	s.Add(s, r)
	s.Mod(s, bn254.ID.ScalarField())
	return s
}

func VerifySchnorrProof(R bn254.G1Affine, e *big.Int, s *big.Int, publicKey string) bool {

	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	y, err := StringToG1Affine(publicKey)
	if err != nil {
		log.Println(err)
	}

	var gS, yE, RplusYE bn254.G1Affine
	gS.ScalarMultiplication(&g, s)
	yE.ScalarMultiplication(&y, e)
	RplusYE.Add(&R, &yE)

	return gS.Equal(&RplusYE)
}

func GenerateFeigeFiatShamirChallenge() (*big.Int, *big.Int) {

	p, err := rand.Prime(rand.Reader, 1024) // Generowanie liczby pierwszej p
	if err != nil {
		panic(err)
	}
	q, err := rand.Prime(rand.Reader, 1024) // Generowanie liczby pierwszej q
	if err != nil {
		panic(err)
	}
	N := new(big.Int).Mul(p, q) // N = p * q (modulo)
	log.Println("N:", N.String())

	e, err := rand.Int(rand.Reader, big.NewInt(2)) // Losowe wyzwanie: 0 lub 1
	if err != nil {
		panic(err)
	}
	log.Println("e:", e.String())

	return N, e
}

func GenerateFeigeFiatShamirProof(privateKey *big.Int, N *big.Int, e *big.Int) (*big.Int, *big.Int, *big.Int) {

	v := new(big.Int).Exp(privateKey, big.NewInt(2), N)

	x := new(big.Int).Exp(v, big.NewInt(2), N) // x = v^2 mod N

	log.Println("x:", x.String())
	// Obliczenie odpowiedzi y
	var y *big.Int
	if e == big.NewInt(0) {
		y = v // Jeśli e = 0, odpowiedź to v
	} else {
		y = new(big.Int).Mul(v, privateKey) // Jeśli e = 1, odpowiedź to v * s mod N
		y.Mod(y, N)
	}

	return x, y, v
}

func VerifyFeigeFiatShamir(x *big.Int, y *big.Int, v *big.Int, N *big.Int, e *big.Int) bool {
	// Weryfikacja odpowiedzi
	var lhs, rhs *big.Int
	lhs = new(big.Int).Exp(y, big.NewInt(2), N) // y^2 mod N

	if e == big.NewInt(0) {
		// Jeśli e = 0, sprawdzamy: y^2 ≡ x mod N
		rhs = x
	} else {
		// Jeśli e = 1, sprawdzamy: y^2 ≡ x * v mod N
		rhs = new(big.Int).Mul(x, v)
		rhs.Mod(rhs, N)
	}

	return lhs.Cmp(rhs) == 0
}

func Schnorr_proof() {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego x
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	log.Println("Klucz prywatny x:", x)
	if err != nil {
		panic(err)
	}

	// Obliczenie klucza publicznego y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)
	log.Println("Klucz publiczny y:", y.String())

	// Wygenerowanie losowej wartości r
	r, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	// Obliczenie zobowiązania R = g^r
	var R bn254.G1Affine
	R.ScalarMultiplication(&g, r)

	// Obliczenie wyzwania e = H(R || y)
	h := sha256.New()
	h.Write(R.Marshal())
	h.Write(y.Marshal())
	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, bn254.ID.ScalarField())

	// Obliczenie odpowiedzi s = r + e * x
	// Klient potrzebuje e, r (wszystko od serwera)
	s := new(big.Int).Mul(e, x)
	s.Add(s, r)
	s.Mod(s, bn254.ID.ScalarField())

	// Weryfikacja dowodu: sprawdzamy, czy g^s = R + y^e
	// Serwer potrzebuje g, y, R, e, s (s od klienta)
	var gS, yE, RplusYE bn254.G1Affine
	gS.ScalarMultiplication(&g, s)
	yE.ScalarMultiplication(&y, e)
	RplusYE.Add(&R, &yE)

	if gS.Equal(&RplusYE) {
		log.Println("Dowód Schnorra poprawny: klient zna klucz prywatny.")
	} else {
		log.Println("Dowód Schnorra niepoprawny: klient nie zna klucza.")
	}
}

func Gnark_crypto_main() {
	// Generator punktu na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1") // Wartość generatora X
	g.Y.SetString("2") // Wartość generatora Y

	if !g.IsOnCurve() {
		panic("Punkt nie należy do krzywej")
	}
	// Wygenerowanie sekretnego klucza x
	// x, err := rand.Int(rand.Reader, fr.Modulus())
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	log.Println("Sekretny klucz x:", x)

	// Obliczenie y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	log.Println("Publiczny klucz y:", y.String())

	// Weryfikacja dowodu (sprawdzenie, czy x rzeczywiście generuje y)
	var check bn254.G1Affine
	check.ScalarMultiplication(&g, x) // to musi zwrocic klient do serwera

	if check.Equal(&y) {
		log.Println("Dowód poprawny: posiadacz zna x")
	} else {
		log.Println("Dowód błędny: x nie pasuje")
	}
}
