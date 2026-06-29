package inclusion

import (
	"regexp"
	"strings"
)

// offRampMessage reemplaza la respuesta del modelo cuando el guardrail detecta que
// cruzó el límite clínico (afirmar un diagnóstico o dar una indicación clínica).
// Deriva al equipo de orientación pero deja la conversación abierta en lo
// pedagógico, alineado al §5 del flujo de comportamiento (alizia-comportamiento-flujo-v1.md).
const offRampMessage = "Eso ya entra en terreno clínico, que lo ve mejor el equipo de orientación o un profesional de la salud; te recomiendo derivarlo ahí. Lo que sí puedo acompañar es lo del aula: contame qué situación concreta observás (qué le cuesta y en qué actividad) y vemos juntos cómo rediseñar la propuesta para que pueda participar."

// guardrailAccents normaliza las vocales acentuadas para que las expresiones de más
// abajo se escriban sin tildes y matcheen igual (la entrada ya viene en minúsculas).
var guardrailAccents = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u",
)

// guardrailSentenceSplit parte el texto en oraciones para acotar la coincidencia de
// patrones: un marcador de diagnóstico y una condición clínica deben caer en la
// MISMA oración para disparar (evita falsos positivos a distancia).
var guardrailSentenceSplit = regexp.MustCompile(`[.!?;\n]+`)

// clinicalConditionRe son condiciones clínicas nombradas. Mencionarlas NO basta para
// disparar (Alizia las nombra a propósito al derivar); hace falta también un acto
// diagnóstico explícito en la misma oración.
var clinicalConditionRe = regexp.MustCompile(`\b(tdah|tda|tea|autismo|autista|asperger|dislexia|dislexico|discalculia|disgrafia|disortografia|sindrome de down|trastorno del espectro|discapacidad intelectual|retraso madurativo|retraso mental|esquizofrenia|trastorno bipolar|bipolaridad)\b`)

// diagnosticActRe son marcas de que el modelo está AFIRMANDO un diagnóstico con
// certeza (no mencionándolo ni difiriéndolo). Conservador a propósito.
var diagnosticActRe = regexp.MustCompile(`(el diagnostico es|tiene claramente|claramente tiene|definitivamente (tiene|es)|se trata de un caso de|es un caso de|presenta un cuadro de|tiene un cuadro de|yo diria que (tiene|es)|estoy segur[ao] de que (tiene|es)|seguro tiene|padece|sufre de|diagnostico de)`)

// clinicalPrescriptionRe son indicaciones clínicas (medicación / tratamiento) que
// Alizia nunca daría legítimamente: disparan por sí solas, sin condición asociada.
var clinicalPrescriptionRe = regexp.MustCompile(`\b(medic\w*|dosis|miligramos|mg|recet\w*|deberia tomar|tiene que tomar|el tratamiento (es|seria|debe|deberia))\b`)

// guardrailExemptRe marca oraciones que NO deben disparar aunque mencionen una
// condición: son derivaciones, negativas o condicionales (justo lo que queremos que
// Alizia haga). Mantiene los falsos positivos casi en cero.
var guardrailExemptRe = regexp.MustCompile(`(no puedo|no me corresponde|no soy quien|no podria|no estoy en condiciones|deriv|equipo de orientacion|profesional|especialista|medico|consulta|consultar|habria que evaluar|seria importante evaluar|si (tiene|tuviera) un diagnostico|podria|quizas|tal vez)`)

// normalizeForGuardrail deja el texto en minúsculas y sin tildes para que los
// patrones (escritos sin acentos) matcheen de forma robusta.
func normalizeForGuardrail(s string) string {
	return guardrailAccents.Replace(strings.ToLower(s))
}

// crossedClinicalLine reporta si la respuesta del modelo cruzó el límite clínico
// (afirmar un diagnóstico o dar una indicación clínica) y por qué motivo. Es un
// guardrail duro, conservador: prioriza no pisar respuestas legítimas por sobre
// atrapar toda forma posible de desvío.
func crossedClinicalLine(text string) (tripped bool, reason string) {
	for _, raw := range guardrailSentenceSplit.Split(normalizeForGuardrail(text), -1) {
		sent := strings.TrimSpace(raw)
		if sent == "" || guardrailExemptRe.MatchString(sent) {
			continue
		}
		if clinicalPrescriptionRe.MatchString(sent) {
			return true, "prescription"
		}
		if diagnosticActRe.MatchString(sent) && clinicalConditionRe.MatchString(sent) {
			return true, "diagnosis"
		}
	}
	return false, ""
}
