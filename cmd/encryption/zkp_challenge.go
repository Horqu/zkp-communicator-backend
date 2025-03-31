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

	// Generowanie 1024-bitowego wyzwania e
	e := new(big.Int)
	for i := 0; i < 1024; i++ {
		bit, err := rand.Int(rand.Reader, big.NewInt(2)) // Losowy bit: 0 lub 1
		if err != nil {
			panic(err)
		}
		e.Lsh(e, 1) // Przesunięcie w lewo o 1 bit
		e.Or(e, bit)
	}
	log.Println("Wyzwanie e (1024 bity):", e.Text(2)) // Wyświetlenie w postaci binarnej

	return N, e
}

func GenerateFeigeFiatShamirProof(privateKey *big.Int, N *big.Int, e *big.Int) ([]*big.Int, []*big.Int, *big.Int) {
	v := new(big.Int).Exp(privateKey, big.NewInt(2), N) // v = privateKey^2 mod N
	var xList []*big.Int
	var yList []*big.Int

	eCopy := new(big.Int).Set(e) // Kopia e

	for i := 0; i < 1024; i++ {
		// Wybór losowej wartości r
		r, err := rand.Int(rand.Reader, N)
		if err != nil {
			panic(err)
		}

		// Obliczenie x = r^2 mod N
		x := new(big.Int).Exp(r, big.NewInt(2), N)
		xList = append(xList, x)

		// Pobranie i-tego bitu z e
		bit := new(big.Int).And(eCopy, big.NewInt(1)) // Pobranie najmłodszego bitu
		eCopy.Rsh(eCopy, 1)                           // Przesunięcie w prawo o 1 bit

		// Obliczenie odpowiedzi y
		var y *big.Int
		if bit.Cmp(big.NewInt(0)) == 0 {
			y = r // Jeśli bit = 0, odpowiedź to r
		} else {
			y = new(big.Int).Mul(r, privateKey) // Jeśli bit = 1, odpowiedź to r * privateKey mod N
			y.Mod(y, N)
		}
		yList = append(yList, y)
	}

	return xList, yList, v
}

func VerifyFeigeFiatShamir(xList []*big.Int, yList []*big.Int, v *big.Int, N *big.Int, e *big.Int) bool {
	eCopy := new(big.Int).Set(e) // Kopia e

	for i := 0; i < 1024; i++ {
		// Pobranie i-tego bitu z e
		bit := new(big.Int).And(eCopy, big.NewInt(1)) // Pobranie najmłodszego bitu
		eCopy.Rsh(eCopy, 1)                           // Przesunięcie w prawo o 1 bit

		// Obliczenie lewej strony równania: y^2 mod N
		lhs := new(big.Int).Exp(yList[i], big.NewInt(2), N)

		// Obliczenie prawej strony równania w zależności od bitu
		var rhs *big.Int
		if bit.Cmp(big.NewInt(0)) == 0 {
			// Jeśli bit = 0, sprawdzamy: y^2 ≡ x mod N
			rhs = xList[i]
		} else {
			// Jeśli bit = 1, sprawdzamy: y^2 ≡ x * v mod N
			rhs = new(big.Int).Mul(xList[i], v)
			rhs.Mod(rhs, N)
		}

		// Porównanie lewej i prawej strony
		if lhs.Cmp(rhs) != 0 {
			log.Printf("Weryfikacja nie powiodła się dla bitu %d", i)
			return false
		}
	}
	return true
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
