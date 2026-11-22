package sync

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func parseJSON(cuerpo string) (map[string]any, error) {
	if cuerpo == "" {
		return nil, nil
	}
	var raiz map[string]any
	if err := json.Unmarshal([]byte(cuerpo), &raiz); err != nil {
		return nil, err
	}
	return raiz, nil
}

func texto(n map[string]any, campo string) string {
	if n == nil {
		return ""
	}
	v, ok := n[campo]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	default:
		return fmt.Sprint(v)
	}
}

func entero(n map[string]any, campo string) *int {
	if n == nil {
		return nil
	}
	v, ok := n[campo]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case float64:
		i := int(t)
		return &i
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil
		}
		return &i
	default:
		return nil
	}
}

func arrayDatos(raiz map[string]any) []map[string]any {
	if raiz == nil {
		return nil
	}
	raw, ok := raiz["data"]
	if !ok {
		return nil
	}
	slice, ok := raw.([]any)
	if !ok {
		return nil
	}
	var out []map[string]any
	for _, item := range slice {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func escribirJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func mapaValores(pares ...any) map[string]any {
	m := make(map[string]any, len(pares)/2)
	for i := 0; i+1 < len(pares); i += 2 {
		k, _ := pares[i].(string)
		m[k] = pares[i+1]
	}
	return m
}