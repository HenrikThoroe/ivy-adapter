package uci

import (
	"errors"
	"strconv"
	"strings"
)

func readUntil(parts []string, start int, key ...string) []string {
	res := make([]string, 0)
	n := len(parts)

	for idx := start; idx < n; idx++ {
		part := parts[idx]

		for _, k := range key {
			if part == k {
				return res
			}
		}

		res = append(res, part)
	}

	return res
}

func parseOptionConfigStr(info string) (OptionConfig, error) {
	parts := strings.Split(info, " ")
	res := OptionConfig{Var: make([]string, 0)}
	idx := 1
	keys := []string{"type", "min", "max", "default", "name", "var", "value"}

	if len(parts) < 3 {
		return res, errors.New("invalid option. option is too short")
	}

	for idx < len(parts)-1 {
		part := parts[idx]

		if part == "type" {
			next := parts[idx+1]
			res.Type = OptionType(next)
			idx += 2

			if res.Type != Check && res.Type != Spin && res.Type != Combo && res.Type != Button && res.Type != String {
				return res, errors.New("invalid option type")
			}
		} else if part == "min" {
			next := parts[idx+1]
			idx += 2

			if i, e := strconv.Atoi(next); e == nil {
				res.Min = i
			} else {
				return res, errors.New("invalid option min value")
			}
		} else if part == "max" {
			next := parts[idx+1]
			idx += 2

			if i, e := strconv.Atoi(next); e == nil {
				res.Max = i
			} else {
				return res, errors.New("invalid option max value")
			}
		} else if part == "default" {
			buf := readUntil(parts, idx+1, keys...)
			res.Def = strings.Join(buf, " ")
			idx += len(buf) + 1
		} else if part == "name" {
			buf := readUntil(parts, idx+1, keys...)
			res.Name = strings.Join(buf, " ")
			idx += len(buf) + 1
		} else if part == "var" {
			buf := readUntil(parts, idx+1, keys...)
			res.Var = append(res.Var, strings.Join(buf, " "))
			idx += len(buf) + 1
		} else {
			return res, errors.New("invalid option config string. unknown key '" + part + "' in '" + info + "'")
		}
	}

	if res.Name == "" || res.Type == "" {
		return res, errors.New("invalid option config string. name or type is missing")
	}

	return res, nil
}

func parseInfoStr(info string) *MoveInfo {
	parts := strings.Split(info, " ")
	res := MoveInfo{}
	var buf *[]string

	for idx := 0; idx < len(parts); idx++ {
		var hit int8
		part := parts[idx]

		if idx < len(parts)-1 {
			hit |= setInt("depth", part, parts[idx+1], &res.Depth)
			hit |= setInt("seldepth", part, parts[idx+1], &res.SelDepth)
			hit |= setInt("time", part, parts[idx+1], &res.Time)
			hit |= setInt("nodes", part, parts[idx+1], &res.Nodes)
			hit |= setInt("multipv", part, parts[idx+1], &res.MultiPv)
			hit |= setInt("currmovenumber", part, parts[idx+1], &res.CurrentMoveNumber)
			hit |= setInt("hashfull", part, parts[idx+1], &res.HashFull)
			hit |= setInt("nps", part, parts[idx+1], &res.Nps)
			hit |= setInt("tbhits", part, parts[idx+1], &res.TbHits)
			hit |= setInt("sbhits", part, parts[idx+1], &res.Sbhits)
			hit |= setInt("cpuload", part, parts[idx+1], &res.CpuLoad)

			if part == "currmove" {
				res.CurrentMove = parts[idx+1]
				hit = 1
			}

			if part == "string" {
				res.String = strings.Join(parts[idx+1:], " ")
				break
			}
		}

		if part == "score" {
			n := parseScore(parts, idx, &res.Score)
			idx += n - 1
		} else if part == "refutation" {
			buf = &res.Refutation
		} else if part == "currline" {
			buf = &res.Currline
		} else if part == "pv" {
			buf = &res.Pv
		} else if hit == 0 && buf != nil {
			*buf = append(*buf, part)
		} else {
			buf = nil
		}

	}

	return &res
}

func parseScore(parts []string, start int, target *Score) int {
	if start >= len(parts) {
		return 0
	}

	if parts[start] != "score" {
		return 0
	}

	idx := start + 1
	n := len(parts)

	for idx < n {
		part := parts[idx]

		if idx+1 < n && setInt("cp", part, parts[idx+1], &target.Value) != 0 {
			target.Type = CP
			idx += 2
		} else if idx+1 < n && setInt("mate", part, parts[idx+1], &target.Value) != 0 {
			target.Type = Mate
			idx += 2
		} else if part == "lowerbound" {
			target.Lowerbound = true
			idx += 1
		} else if part == "upperbound" {
			target.Upperbound = true
			idx += 1
		} else {
			return idx - start
		}
	}

	return idx - start
}

func setInt(match string, key string, val string, target *int) int8 {
	if key == match {
		i, e := strconv.Atoi(val)

		if e == nil {
			*target = i
		}

		return 1
	}

	return 0
}
