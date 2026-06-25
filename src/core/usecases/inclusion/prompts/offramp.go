package prompts

// Off-ramp default wording, editable at the prompts layer.
// Handoff is the LAST resort, not the first: Alizia always attempts to help
// with available resources (see scopeRules). The off-ramp is used in two cases only:
//
//  1. OffRampInvalidOutput: the code-level guardrail rejected the response and
//     repair was not possible. The engine substitutes this text (invalid output
//     is never shown to the user).
//  2. OffRampOutOfScope: the request falls outside pedagogical scope (clinical /
//     crisis / diagnosis). The prompt framework instructs the model to use this
//     text when handing off, ensuring consistent wording and no diagnosis.
const (
	OffRampInvalidOutput = "Perdón, no pude armar una recomendación válida con la valija en este momento. " +
		"¿Probamos de nuevo describiéndome la situación del alumno con otras palabras?"

	OffRampOutOfScope = "Esto se me va un poco del alcance pedagógico que puedo cubrir. No puedo hacer un " +
		"diagnóstico ni reemplazar una evaluación profesional. Si hay una situación clínica o de crisis, lo " +
		"mejor es derivar al equipo de orientación o a un profesional. ¿Querés que igual te ayude con alguna " +
		"adaptación de aula para acompañar mientras tanto?"
)
