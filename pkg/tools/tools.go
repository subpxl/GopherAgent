	package tools

	import "fmt"

	func calculate(input string) string {
		var a, b float64
		var op string

		_, err := fmt.Sscanf(input, "%f %s %f", &a, &op, &b)
		if err != nil {
			return "error: expected format 'a op b', got: " + input
		}

		switch op {
		case "+":
			return fmt.Sprintf("%g", a+b)
		case "-":
			return fmt.Sprintf("%g", a-b)
		case "*":
			return fmt.Sprintf("%g", a*b)
		case "/":
			if b == 0 {
				return "error: division by zero"
			}
			return fmt.Sprintf("%g", a/b)
		default:
			return "error: unknown operator: " + op
		}
	}

	func CalculatorTool() Tool {
		return Tool{
			Name:        "calculator",
			Description: "performs arithmetic on two numbers. input format: 'a op b' where op is +, -, *, /. example: '12.5 * 3'",
			Fn:          calculate,
		}
	}
