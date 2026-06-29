package inclusion

import (
	"regexp"
	"strconv"
	"strings"
)

// gradeOrdinalWords mapea ordinales en español (con y sin tilde) al número de grado.
var gradeOrdinalWords = map[string]int{
	"primero": 1, "primer": 1, "primera": 1,
	"segundo": 2, "segunda": 2,
	"tercero": 3, "tercer": 3, "tercera": 3,
	"cuarto": 4, "cuarta": 4,
	"quinto": 5, "quinta": 5,
	"sexto": 6, "sexta": 6,
	"septimo": 7, "séptimo": 7, "septima": 7, "séptima": 7,
}

// gradeCanonical es la forma canónica de cada grado (1..7).
var gradeCanonical = map[int]string{
	1: "1ro", 2: "2do", 3: "3ro", 4: "4to", 5: "5to", 6: "6to", 7: "7mo",
}

var (
	gradeDigitRegex  = regexp.MustCompile(`\d+`)
	gradeSuffixRegex = regexp.MustCompile(`^(ro|do|to|mo|er|ª|º|°)`)
)

// normalizeGrade interpreta cómo el docente nombra un aula ("3ro A", "tercero B",
// "3°", "3 A", "primer grado") y devuelve su forma canónica: name ("3ro A"),
// grade ("3ro") y section ("A", mayúscula; vacío si no la dio). Si no logra
// reconocer un grado válido (1..7), devuelve los tres vacíos.
func normalizeGrade(input string) (name, grade, section string) {
	s := strings.ToLower(strings.TrimSpace(input))
	if s == "" {
		return "", "", ""
	}

	num := 0
	rest := s
	if loc := gradeDigitRegex.FindStringIndex(s); loc != nil {
		num, _ = strconv.Atoi(s[loc[0]:loc[1]])
		// Sacamos el dígito y su sufijo ordinal inmediato ("3ro" -> resto " A").
		after := gradeSuffixRegex.ReplaceAllString(s[loc[1]:], "")
		rest = s[:loc[0]] + " " + after
	} else {
		// Match por token exacto: evita que "primero" matchee el substring "primer"
		// y deje una letra suelta que se confunda con la sección.
		tokens := strings.Fields(s)
		for i, tok := range tokens {
			if n, ok := gradeOrdinalWords[tok]; ok {
				num = n
				rest = strings.Join(append(tokens[:i:i], tokens[i+1:]...), " ")
				break
			}
		}
	}

	canonical, ok := gradeCanonical[num]
	if !ok {
		return "", "", ""
	}
	grade = canonical

	// La sección es la primera letra suelta que quede (ignorando "grado", "°", etc.).
	rest = strings.ReplaceAll(rest, "grado", " ")
	for _, tok := range strings.Fields(rest) {
		if len(tok) == 1 && tok >= "a" && tok <= "z" {
			section = strings.ToUpper(tok)
			break
		}
	}

	if section != "" {
		return grade + " " + section, grade, section
	}
	return grade, grade, ""
}

var matchNameSpaceRegex = regexp.MustCompile(`\s+`)

// normalizeName canoniza un nombre para comparar sin distinguir mayúsculas,
// acentos ni espacios extra. Se usa para reconocer alumnos ya existentes.
func normalizeName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = matchNameSpaceRegex.ReplaceAllString(s, " ")
	return stripAccents(s)
}

var accentReplacer = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
)

func stripAccents(s string) string {
	return accentReplacer.Replace(s)
}
