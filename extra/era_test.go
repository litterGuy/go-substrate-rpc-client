package extra

import (
	"fmt"
	"testing"
)

func TestGetEra(t *testing.T) {
	//104 26
	//4137894- 4138406

	//88 2
	//4138853- 4138917
	era := GetEra(4138853)
	fmt.Printf("%+v", era)
}
