package base

const ( 
	MaxHP int = 1000
	Death int = -618
)

type Life struct {
	Vitality float64
	Rate int
	// barrier
}


// NEW
func InitLife() *Life {
	var buffer Life
	buffer.Rate = 618
	return &buffer
}

// READ 
func (life *Life) Full() bool { return (*life).Rate > MaxHP }
func (life *Life) Wounded() bool { return (*life).Rate <= 0 }
func (life *Life) Dead() bool { return (*life).Rate < Death }

// MOD
func (life *Life) HealDamage(amount int) {
	if life.Wounded() {
		if amount > 0 { (*life).Rate += 1 } else if amount < 0 { (*life).Rate += -1 }
	} else { (*life).Rate += amount }
	if life.Dead() { (*life).Rate = Death ; return }
	if life.Full() { (*life).Rate = MaxHP ; return }
}

