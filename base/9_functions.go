package base

import (
	"math"
	"time"
	"math/rand"
	"crypto/sha512"
  "encoding/binary"
)

const MinEntropy = 0.0132437
const StartEpoch = int(time.Now().UnixNano())


// TIME
func Epoch() int { return (EpochNS()-StartEpoch)/1000000 }
func EpochNS() int { return int(time.Now().UnixNano())-StartEpoch }
func Wait(ms float64) { time.Sleep( time.Millisecond * time.Duration( ms )) }


// HP, STATS
func Vector(args ...float64) float64 { sum := 0.0 ; for _,each := range args { sum += each*each } ; return math.Sqrt(sum) }
func Round(a float64) int { return int(math.Round(a)) }
func CeilRound(a float64) int { return int(math.Ceil(a)) } 
func FloorRound(a float64) int { return int(math.Floor(a)) }
func ChancedRound(a float64) int {
  b,l:=math.Ceil(a),math.Floor(a)
  c:=math.Abs(math.Abs(a)-math.Abs(math.Min(b, l)))
  if a<0 {c = 1-c}
  if Rand() < c {return int(b)} else {return int(l)}
  return 0
}


// RANDOMIZER
// func Stability(max, vector float64) float64 {
//   return math.Log2( math.Abs(vector) / math.Abs(max) )/math.Log2(math.Sqrt(3))
//   // can be used to compare heat: vec = threshold, max = current [ x>1 ? overheat : x<-1 ? dumb : normal]
// }
func Near(c int) float64 { sum := 0.0 ; for x:=0; x<c; x++ { sum += Rand()+Rand() } ; return sum } 
func Rand() float64 {
  x := (time.Now().UnixNano())
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(x))
  hsum := sha512.Sum512(in_bytes)
  sum  := binary.BigEndian.Uint64(hsum[:])
  return rand.New(rand.NewSource( int64(sum) )).Float64()
}
func Ntrp(a float64) float64 { 
  randy := (Epoch() % 1000) / 300
  entropy := math.Log10( math.Abs(a)+1 )/25 
  if randy == 2 { a = a*(1+MinEntropy+entropy) }
  if randy == 0 { a = a/(1+MinEntropy+entropy) }
  return math.Round( a*1000 ) / 1000
}
