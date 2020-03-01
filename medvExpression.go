package segmentation

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

type MedvExpression struct {
}

func (e MedvExpression) compile(expression string) (interface{}, error) {
	return expr.Compile(expression)
}

func (e MedvExpression) execute(program interface{}, env interface{}) (bool, error) {
	result, err := expr.Run(program.(*vm.Program), env)

	if err != nil {
		return false, err
	}

	return result.(bool), err
}
