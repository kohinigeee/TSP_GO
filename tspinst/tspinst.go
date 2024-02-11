package tspinst

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TspInst struct {
	ProblemName string
	Points      []*Point
	PointsDim   int
}

// ポイントのインデックス
type PointIndexT int

func (t *TspInst) String() string {
	ans := ""
	ans += "{ "
	ans += fmt.Sprintf("ProblemName: %v, Dimension: %v,\n", t.ProblemName, t.PointsDim)

	ans += "Points: [ "
	for _, p := range t.Points {
		ans += fmt.Sprintf("(%v, %v) ", p.X, p.Y)
	}
	ans += "] }"
	return ans
}

func (t *TspInst) Point(i PointIndexT) *Point {
	return t.Points[i]
}

// TSPのファイルからインスタンスを読み込む
func LoadTspInst(fpath string) (*TspInst, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tokenMap := make(map[string]string)

	// ファイルの先頭部分を読み込む
	for scanner.Scan() {
		s := scanner.Text()
		if s == "NODE_COORD_SECTION" {
			break
		}

		tokens := strings.Split(s, ":")
		if len(tokens) > 1 {
			tokens[0] = strings.TrimSpace(tokens[0])
			tokens[1] = strings.TrimSpace(tokens[1])
			tokenMap[tokens[0]] = tokens[1]
		}
	}

	tspInfos, err := NewTspFileInfos(tokenMap)
	if err != nil {
		return nil, fmt.Errorf("[LoadTspInst] tspInfos is incorrect : %v", err)
	}

	// ポイントの部分を読み込む
	points := make([]*Point, tspInfos.Dimension)
	for i := 0; i < tspInfos.Dimension; i++ {
		scanner.Scan()
		s := scanner.Text()
		tokens := strings.Fields(s)

		if len(tokens) != 3 {
			return nil, fmt.Errorf("[LoadTspInst] Point string's length is not 3 : %v [line %d]", tokens, i)
		}

		x, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Point[%v]: %v", i, tokens)
		}

		y, err := strconv.ParseFloat(tokens[2], 64)
		if err != nil {
			return nil, fmt.Errorf("Point[%v]: %v", i, tokens)
		}

		points[i] = NewPoint(int(x), int(y), i)
	}

	return &TspInst{
		ProblemName: tspInfos.ProblemName,
		Points:      points,
		PointsDim:   tspInfos.Dimension,
	}, nil
}

type TspFileInfos struct {
	ProblemName    string
	ProblemType    string
	Comment        string
	Dimension      int
	EdgeWeightType string
}

func NewTspFileInfos(tokens map[string]string) (*TspFileInfos, error) {
	Infos := new(TspFileInfos)

	problemName, ok := tokens["NAME"]
	if !ok {
		return nil, errors.New("Problem Name not found")
	}
	Infos.ProblemName = problemName

	Infos.ProblemType = tokens["TYPE"]
	Infos.Comment = tokens["COMMENT"]

	dimension, err := strconv.Atoi(tokens["DIMENSION"])
	if err != nil {
		return nil, fmt.Errorf("TspInfo[dimension=%v]: %v", tokens["DIMENSION"], err)
	}

	if dimension < 0 {
		return nil, fmt.Errorf("TspInfo[dimension=%v]: dimension must be greater than 0", tokens["DIMENSION"])
	}
	Infos.Dimension = dimension

	Infos.EdgeWeightType = tokens["EDGE_WEIGHT_TYPE"]

	return Infos, nil
}
