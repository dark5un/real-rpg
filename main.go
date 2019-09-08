package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	proto "real-rpg/proto"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
)

// Dice is
type Dice struct{}

// Roll is
func (d *Dice) Roll(ctx context.Context, req *proto.RollRequest, res *proto.RollResponse) error {
	re := regexp.MustCompile(`(?m)([dD]\d+)|((\d+)([dD]\d+))`)
	exp := req.Expression

	found := re.FindAllStringSubmatch(exp, -1)
	for _, matches := range found {
		if len(matches[2]) > 0 {
			count, err := strconv.Atoi(matches[3])
			if err != nil {
				return err
			}
			var newExpression string
			for index := 0; index < count; index++ {
				if index == 0 {
					newExpression = matches[4]
				} else {
					newExpression = newExpression + "+" + matches[4]
				}
			}
			exp = strings.Replace(exp, matches[0], newExpression, 1)
		}
	}
	evaluation := exp

	found = re.FindAllStringSubmatch(exp, -1)
	for _, matches := range found {
		di, err := strconv.Atoi(strings.ReplaceAll(strings.ToLower(matches[0]), "d", ""))
		if err != nil {
			return err
		}
		radomd := rand.Intn(di) + 1
		exp = strings.Replace(exp, matches[0], strconv.Itoa(radomd), 1)
		evaluation = strings.Replace(evaluation, matches[0], fmt.Sprint("{", radomd, "}"), 1)
	}

	parsedExp, err := parser.ParseExpr(exp)
	if err != nil {
		return err
	}

	res.Result = int64(Eval(parsedExp))
	res.Evaluation = evaluation
	return nil
}

// Eval is
func Eval(exp ast.Expr) int {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i
		}
	}

	return 0
}

// EvalBinaryExpr is
func EvalBinaryExpr(exp *ast.BinaryExpr) int {
	left := Eval(exp.X)
	right := Eval(exp.Y)

	switch exp.Op {
	case token.ADD:
		return left + right
	case token.SUB:
		return left - right
	case token.MUL:
		return left * right
	case token.QUO:
		return left / right
	}

	return 0
}

func runClientRoll(service micro.Service, eval string) {
	// Create new dice client
	dice := proto.NewDiceService("dice", service.Client())

	// Call the dice
	rsp, err := dice.Roll(context.TODO(), &proto.RollRequest{Expression: eval})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	fmt.Println(rsp.Result)
	fmt.Println(rsp.Evaluation)
}

func main() {
	var clientRoll string
	service := micro.NewService(
		micro.Name("dice"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "dnd",
		}),
		micro.Flags(cli.StringFlag{
			Name:        "roll",
			Destination: &clientRoll,
			Usage:       "Roll the dice",
		}),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			if c.IsSet("roll") {
				runClientRoll(service, clientRoll)
				os.Exit(0)
			}
		}),
	)

	proto.RegisterDiceHandler(service.Server(), new(Dice))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
