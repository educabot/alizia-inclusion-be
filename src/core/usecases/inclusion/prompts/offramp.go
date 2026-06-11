package prompts

// Off-ramp (HU-6, T-6.3) — wording por defecto, editable, en la capa de prompts.
// El "paso al costado" es el ÚLTIMO recurso, no el primero: Alizia siempre intenta
// ayudar con lo que tiene (ver scopeRules). El off-ramp se usa solo en dos casos:
//
//  1. OffRampInvalidOutput: el guardrail por código rechazó la respuesta y no se
//     pudo reparar. El motor la sustituye por este texto (nunca muestra salida
//     inválida).
//  2. OffRampOutOfScope: el pedido queda fuera de alcance (clínico / crisis /
//     diagnóstico). El propio marco le pide al modelo usar este texto al derivar,
//     para que el wording sea consistente y no diagnostique.
const (
	OffRampInvalidOutput = "Perdón, no pude armar una recomendación válida con la valija en este momento. " +
		"¿Probamos de nuevo describiéndome la situación del alumno con otras palabras?"

	OffRampOutOfScope = "Esto se me va un poco del alcance pedagógico que puedo cubrir. No puedo hacer un " +
		"diagnóstico ni reemplazar una evaluación profesional. Si hay una situación clínica o de crisis, lo " +
		"mejor es derivar al equipo de orientación o a un profesional. ¿Querés que igual te ayude con alguna " +
		"adaptación de aula para acompañar mientras tanto?"
)
