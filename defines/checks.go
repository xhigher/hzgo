package defines

var (
	pageLimits = []int32{10, 20, 50, 100}
)


func CheckPageLimit(limit int32) bool{
	for _, l := range pageLimits {
		if l == limit {
			return true
		}
	}
	return false
}
