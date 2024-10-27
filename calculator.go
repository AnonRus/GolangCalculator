package main

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
)

func Calc(expression string) (float64, error) {
  expression = strings.ReplaceAll(expression, " ", "")

  result, _, err := parseExpression(expression, 0)
  if err != nil {
    return 0, err
  }
  return result, nil
}

func parseExpression(expr string, pos int) (float64, int, error) {
  value, pos, err := parseTerm(expr, pos)
  if err != nil {
    return 0, pos, err
  }

  for pos < len(expr) {
    op := expr[pos]
    if op != '+' && op != '-' {
      break
    }
    pos++

    nextValue, nextPos, err := parseTerm(expr, pos)
    if err != nil {
      return 0, nextPos, err
    }

    if op == '+' {
      value += nextValue
    } else {
      value -= nextValue
    }
    pos = nextPos
  }

  return value, pos, nil
}

func parseTerm(expr string, pos int) (float64, int, error) {
  value, pos, err := parseFactor(expr, pos)
  if err != nil {
    return 0, pos, err
  }

  for pos < len(expr) {
    op := expr[pos]
    if op != '*' && op != '/' {
      break
    }
    pos++

    nextValue, nextPos, err := parseFactor(expr, pos)
    if err != nil {
      return 0, nextPos, err
    }

    if op == '*' {
      value *= nextValue
    } else {
      if nextValue == 0 {
        return 0, nextPos, errors.New("деление на ноль")
      }
      value /= nextValue
    }
    pos = nextPos
  }

  return value, pos, nil
}

func parseFactor(expr string, pos int) (float64, int, error) {
  if pos < len(expr) && expr[pos] == '(' {
    pos++
    value, nextPos, err := parseExpression(expr, pos)
    if err != nil {
      return 0, nextPos, err
    }
    if nextPos >= len(expr) || expr[nextPos] != ')' {
      return 0, nextPos, errors.New("нехватка закрывающей скобки")
    }
    return value, nextPos + 1, nil
  }

  return parseNumber(expr, pos)
}

func parseNumber(expr string, pos int) (float64, int, error) {
  startPos := pos
  for pos < len(expr) && (expr[pos] >= '0' && expr[pos] <= '9' || expr[pos] == '.') {
    pos++
  }
  if startPos == pos {
    return 0, pos, errors.New("ожидалось число")
  }

  value, err := strconv.ParseFloat(expr[startPos:pos], 64)
  if err != nil {
    return 0, pos, errors.New("неправильный формат числа")
  }
  return value, pos, nil
}

func main() {
  expression := "3 + 5 * (2 - 8)"
  result, err := Calc(expression)
  if err != nil {
    fmt.Println("Ошибка:", err)
  } else {
    fmt.Printf("Результат: %f\n", result)
  }
}
