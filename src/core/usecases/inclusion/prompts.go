package inclusion

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type GeneratedAdaptation struct {
	Title       string                    `json:"title"`
	Type        string                    `json:"type"`
	Strategy    string                    `json:"strategy"`
	DeviceIDs   []int64                   `json:"device_ids"`
	DeviceNames []string                  `json:"device_names"`
	RampID      *int64                    `json:"ramp_id,omitempty"`
	Steps       []entities.AdaptationStep `json:"steps,omitempty"`
}

// aliziaPersona es la identidad ÚNICA (capa 1, cacheable) que encabeza todo system
// prompt. Una sola voz: misma identidad, tono y límites recomiende, asista o guíe.
// Lo que cambia por momento (cadencia, repregunta, RAG) es comportamiento, no esto.
// Ver alizia-persona-base-v2.md.
const aliziaPersona = `Sos Alizia, la asistente de inclusión educativa de Educabot. Acompañás a docentes de aula y a maestras y maestros integradores a remover barreras de aprendizaje y a diseñar la clase para que todos puedan participar. Partís siempre de la situación observable del aula, no del diagnóstico. Trabajás desde el Diseño Universal para el Aprendizaje (DUA): ofrecés distintas formas de representar el contenido, de participar y de expresar lo aprendido, con ajustes proporcionados a cada alumno.

VOZ Y TONO:
- Cálida pero medida, profesional. Español rioplatense, tratás de "vos". Sin jerga clínica.
- Concreta y accionable: el docente suele leerte en plena clase, así que vas al grano. Una idea por vez.

TU LUGAR:
- Aportás ideas y acompañás la decisión del docente; la última palabra es suya.
- Tu terreno es lo pedagógico; lo clínico lo conducen los profesionales de salud.
- Hablás con un especialista: no expliques lo obvio ni describas para qué sirve un material que el docente ya conoce. Sumá criterio pedagógico, no repitas catálogo.
- Tu primer reflejo es ayudar con lo pedagógico que tengas; derivar es el último recurso, no la salida por defecto. Solo ante algo clínico, una crisis o un pedido de diagnóstico, nombralo con cuidado y derivá al equipo de orientación o a un profesional, sin cerrar la conversación: seguís disponible para lo del aula.

CÓMO RESPONDÉS:
- Primero la estrategia pedagógica (DUA). Un dispositivo de la valija es UNA opción posible, no el objetivo: muchas adaptaciones no necesitan material físico.
- Proponés ajustes proporcionados, partiendo de lo observable.
- Recomendás apoyos o dispositivos solo si existen en el catálogo, nombrándolos por lo que son.

HONESTIDAD (no negociable):
- Nunca afirmes haber consultado bibliografía, fuentes, papers, guías o "material" si no lo hiciste en este turno con una herramienta de búsqueda. No inventes ni des a entender una búsqueda que no ocurrió.
- Lo que sale de tu criterio decilo como tal ("desde el enfoque DUA", "por lo general en el aula"), sin atribuirlo a una fuente que no abriste.
- Si el docente te pide en qué te basás y no tenés material a mano, sé honesta: ofrecé el fundamento pedagógico que sí tenés y aclaralo, en vez de simular respaldo bibliográfico.
`

// repreguntaGate es el gate de repregunta (pedido central de pedagogía): antes de
// proponer, si falta contexto clave, una sola pregunta. Ver
// alizia-comportamiento-flujo-v1.md §2.
const repreguntaGate = `ANTES DE PROPONER:
- Si falta contexto clave (la barrera observable concreta, para quién y en qué actividad), hacé UNA sola pregunta clara y esperá. No respondas genérico.
- Ej.: "le cuesta escribir" puede ser el agarre/motricidad, sostener la atención, organizar las ideas o copiar del pizarrón: cada uno lleva a otra adaptación. Si dice "le tiembla la mano al escribir", preguntá lo que afina la propuesta (¿siempre o en ciertos momentos?, ¿una mano o las dos?, ¿al empezar o tras un rato?) antes de recomendar un soporte concreto.
- Si el docente ya dio el dato, no lo vuelvas a pedir. Si pide algo rápido o el dato no es imprescindible, proponé con un supuesto explícito ("Asumo X; si es otra cosa, decime y ajusto").
`

// fundamentosRAG instruye el uso del RAG agéntico. SOLO se inyecta cuando el modo
// agéntico está activo (AI_AGENTIC_ENABLED=true): si no, las tools search_content/
// get_content no existen y no hay que instruir su uso. Ver flujo §4.
const fundamentosRAG = `FUNDAMENTOS (material pedagógico real):
- Ante cualquier pregunta sobre una discapacidad, barrera, estrategia pedagógica, marco o normativa, DEBÉS llamar search_content_hibrido ANTES de responder. No la uses para charla trivial.
- Si el docente pide bibliografía, fuentes, evidencia, referencias o "en qué te basás", DEBÉS llamar search_content_hibrido y responder con lo que devuelva. Nunca contestes que buscaste si no llamaste la tool en este turno.
- Pasá la pregunta del docente completa en semantic_question y las palabras clave en terms (temas + nombre de la discapacidad/barrera).
- Fundamentá tu respuesta con lo que devuelve, integrándolo de forma natural (no hace falta citar el título del documento). Si el preview es pertinente y necesitás más, usá get_content.
- Si la búsqueda vuelve vacía, no inventes: respondé con los lineamientos base aclarando que no hay material cargado sobre ese punto.
- Los materiales de la valija ya están en el catálogo de este prompt; no los busques por search_content_hibrido.
- Si te apoyás en material del corpus, citá la fuente con [CONTENT_ID:X], usando el id del recurso (resource_id) que devolvió la búsqueda. El sistema lo convierte en un chip; no menciones el id de otra forma.
`

// alumnoFlow guía el Caso 2 ("tengo un alumno con tal barrera/diagnóstico"):
// reconocer primero, traer contexto si ya lo conoce, y ofrecer crearlo sin forzar
// si no. La creación real la hace la tool create_student, SOLO tras confirmación.
// Depende de las tools agénticas: se inyecta únicamente cuando agentic=true.
const alumnoFlow = `CUANDO EL DOCENTE TE TRAE UN ALUMNO:
- Primero fijate si YA lo conocés: mirá la lista de "alumnos que conocés" (más abajo). Si no aparece pero el docente da un nombre, buscalo con find_student_by_name antes de asumir que es nuevo. Si lo encontrás, traé su contexto con get_student (y get_student_history / get_past_adaptations) ANTES de proponer, para construir sobre lo que ya se sabe y no repetir lo probado.
- Si no lo conocés y falta la barrera observable concreta, hacé UNA pregunta para entenderla (no pidas el diagnóstico).
- No exijas el nombre para ayudar: podés proponer igual. Pero ofrecelo sin presionar, ej.: "si me decís el nombre y el aula lo creamos y queda guardado para la próxima". Podés sugerirlo también más tarde.
- SOLO cuando el docente confirme que sí, dalo de alta:
  1. Resolvé el aula: pedila en formato "3ro A" / "tercero B" si no la sabés. Buscala con list_classrooms; si no existe, creala con create_classroom (pasá el grado tal como lo dijo el docente) y usá el id que devuelve.
  2. Llamá create_student con name + classroom_id (+ la barrera observable como difficulties). Es idempotente: si ya existía, te devuelve el alumno sin duplicar.
- Nunca llames create_student sin confirmación explícita, ni para un alumno que ya reconociste. Con el id que devuelve, usá [STUDENT_ID:X] y enlazá la adaptación a ese alumno.
`

// writeDeviceCatalog imprime el catálogo de dispositivos en formato estable.
func writeDeviceCatalog(b *strings.Builder, devices []entities.Device, withDetail bool) {
	for i := range devices {
		d := &devices[i]
		fmt.Fprintf(b, "- [ID:%d] %s", d.ID, d.Name)
		if d.NeedsDescription != nil {
			fmt.Fprintf(b, " — %s", *d.NeedsDescription)
		}
		b.WriteString("\n")
		if withDetail && d.Rationale != nil {
			fmt.Fprintf(b, "  Fundamento: %s\n", *d.Rationale)
		}
		if withDetail && d.HowToUse != nil {
			fmt.Fprintf(b, "  Uso: %s\n", *d.HowToUse)
		}
	}
}

func buildRecommendSystemPrompt(devices []entities.Device) string {
	var b strings.Builder

	b.WriteString(aliziaPersona)
	b.WriteString("\n")

	b.WriteString("LINEAMIENTOS:\n")
	b.WriteString("- Entrada pedagógica, no clínica: partís de situaciones de aula, no de diagnósticos.\n")
	b.WriteString("- Diferenciación (DUA): proponé variaciones de la actividad en al menos tres niveles.\n")
	b.WriteString("- Respuestas accionables: concretas, breves, aplicables enseguida.\n")
	b.WriteString("- Coherencia: ofrecé 1-3 acciones claras, ordenadas por impacto.\n\n")

	b.WriteString("CATÁLOGO DE DISPOSITIVOS:\n")
	writeDeviceCatalog(&b, devices, true)

	b.WriteString("\nFORMATO DE RESPUESTA:\n")
	b.WriteString("1. Explicación pedagógica breve de por qué el ajuste es adecuado.\n")
	b.WriteString("2. Cómo integrarlo en la actividad descripta.\n")
	b.WriteString("3. Tips prácticos.\n")
	b.WriteString("4. Si recomendás un dispositivo del catálogo, incluí [DEVICE_ID:X] con su ID.\n")
	b.WriteString("5. Al final de tu respuesta, incluí un bloque estructurado con este formato exacto:\n")
	b.WriteString("[ADAPTATION_JSON:{\"title\":\"título corto\",\"type\":\"tipo\",\"strategy\":\"resumen de estrategia\",\"ramp_id\":N,\"device_ids\":[1,2],\"device_names\":[\"nombre1\",\"nombre2\"],\"steps\":[{\"orden\":1,\"texto\":\"primer paso\"},{\"orden\":2,\"texto\":\"segundo paso\"}]}]\n")
	b.WriteString("Los tipos válidos son: actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente.\n")
	b.WriteString("ramp_id = categoría/necesidad del catálogo. steps = el PASO A PASO de la guía (la parte más importante del recurso), claro y accionable.\n")
	b.WriteString("Si la adaptación no usa material físico, usá estrategia_aula con device_ids vacío.\n")

	return b.String()
}

func buildRecommendUserPrompt(student *entities.Student, req RecommendDeviceRequest) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Asignatura: %s\n", req.Subject)
	if req.Objective != "" {
		fmt.Fprintf(&b, "Objetivo de la clase: %s\n", req.Objective)
	}
	if req.Duration != "" {
		fmt.Fprintf(&b, "Duración: %s\n", req.Duration)
	}
	if req.Dynamic != "" {
		fmt.Fprintf(&b, "Dinámica: %s\n", req.Dynamic)
	}
	if req.Materials != "" {
		fmt.Fprintf(&b, "Materiales: %s\n", req.Materials)
	}

	fmt.Fprintf(&b, "\nAlumno: %s\n", student.Name)
	if student.Profile != nil {
		p := student.Profile
		if p.IsTransitory {
			b.WriteString("Condición: transitoria\n")
		} else {
			b.WriteString("Condición: permanente\n")
		}
		if len(p.Difficulties) > 0 {
			fmt.Fprintf(&b, "Dificultades: %s\n", strings.Join(p.Difficulties, ", "))
		}
		if p.FreeDescription != nil && *p.FreeDescription != "" {
			fmt.Fprintf(&b, "Descripción: %s\n", *p.FreeDescription)
		}
	}

	return b.String()
}

// maxKnownStudentsInPrompt acota cuántos alumnos se listan en el prompt. Si el
// docente conoce más, la lista se trunca y se le indica usar find_student_by_name.
const maxKnownStudentsInPrompt = 60

// writeKnownStudents imprime los alumnos que el docente conoce (de toda la org, no
// solo de un aula) con sus dificultades, hasta un tope. Si se supera, nota de que
// use find_student_by_name para encontrar al resto.
func writeKnownStudents(b *strings.Builder, students []entities.Student) {
	if len(students) == 0 {
		return
	}
	b.WriteString("ALUMNOS QUE CONOCÉS:\n")
	shown := students
	truncated := false
	if len(shown) > maxKnownStudentsInPrompt {
		shown = shown[:maxKnownStudentsInPrompt]
		truncated = true
	}
	for i := range shown {
		s := &shown[i]
		fmt.Fprintf(b, "- [ID:%d] %s", s.ID, s.Name)
		if s.Profile != nil && len(s.Profile.Difficulties) > 0 {
			fmt.Fprintf(b, " — Dificultades: %s", strings.Join(s.Profile.Difficulties, ", "))
		}
		b.WriteString("\n")
	}
	if truncated {
		fmt.Fprintf(b, "(y más; si no ves al alumno acá, buscalo con find_student_by_name)\n")
	}
	b.WriteString("\n")
}

// writeStudentContext agrega el contexto rico que arma el Context Assembler
// (PromptContext): docente, alumno foco con su perfil, diagnósticos informados, PPI,
// adaptaciones previas y resúmenes de charlas anteriores. Nil-safe: si no hay
// contexto (sin alumno foco o el assembler falló) no escribe nada y el prompt queda
// igual que antes. Ver alizia-comportamiento-flujo-v1.md §6.
func writeStudentContext(b *strings.Builder, pc *PromptContext) {
	if pc == nil {
		return
	}

	if t := pc.Teacher; t != nil {
		var parts []string
		if t.Specialization != nil && *t.Specialization != "" {
			parts = append(parts, "especialidad "+*t.Specialization)
		}
		if len(t.Subjects) > 0 {
			parts = append(parts, "materias "+strings.Join(t.Subjects, ", "))
		}
		if t.YearsExperience != nil {
			parts = append(parts, fmt.Sprintf("%d años de experiencia", *t.YearsExperience))
		}
		if len(parts) > 0 {
			fmt.Fprintf(b, "DOCENTE: %s.\n\n", strings.Join(parts, "; "))
		}
	}

	s := pc.TargetStudent
	if s == nil {
		writeMissingData(b, pc.MissingData)
		return
	}

	name := s.Name
	if s.PreferredName != nil && *s.PreferredName != "" {
		name = *s.PreferredName
	}
	fmt.Fprintf(b, "ALUMNO FOCO: %s [ID:%d]", name, s.ID)
	if s.GradeLevel != nil && *s.GradeLevel != "" {
		fmt.Fprintf(b, " — %s", *s.GradeLevel)
	}
	b.WriteString("\n")

	if p := s.Profile; p != nil {
		if p.IsTransitory {
			b.WriteString("- Condición: transitoria\n")
		} else {
			b.WriteString("- Condición: permanente\n")
		}
		if p.SupportLevel != nil && *p.SupportLevel != "" {
			fmt.Fprintf(b, "- Nivel de apoyo: %s\n", *p.SupportLevel)
		}
		if len(p.Difficulties) > 0 {
			fmt.Fprintf(b, "- Dificultades observables: %s\n", strings.Join(p.Difficulties, ", "))
		}
		if len(p.Strengths) > 0 {
			fmt.Fprintf(b, "- Fortalezas: %s\n", strings.Join(p.Strengths, ", "))
		}
		if len(p.Interests) > 0 {
			fmt.Fprintf(b, "- Intereses: %s\n", strings.Join(p.Interests, ", "))
		}
		if len(p.EffectiveStrategies) > 0 {
			fmt.Fprintf(b, "- Estrategias que funcionan: %s\n", strings.Join(p.EffectiveStrategies, ", "))
		}
		if len(p.IneffectiveStrategies) > 0 {
			fmt.Fprintf(b, "- Estrategias que NO funcionan: %s\n", strings.Join(p.IneffectiveStrategies, ", "))
		}
		if len(p.Triggers) > 0 {
			fmt.Fprintf(b, "- Disparadores a evitar: %s\n", strings.Join(p.Triggers, ", "))
		}
		if p.FreeDescription != nil && *p.FreeDescription != "" {
			fmt.Fprintf(b, "- Descripción: %s\n", *p.FreeDescription)
		}
		if p.EnvironmentNotes != nil && *p.EnvironmentNotes != "" {
			fmt.Fprintf(b, "- Entorno: %s\n", *p.EnvironmentNotes)
		}
		if p.HasTherapeuticCompanion != nil && *p.HasTherapeuticCompanion {
			b.WriteString("- Tiene acompañante terapéutico (AT) en el aula.\n")
		}
	}

	if len(pc.Diagnoses) > 0 {
		names := make([]string, 0, len(pc.Diagnoses))
		for i := range pc.Diagnoses {
			d := &pc.Diagnoses[i]
			if d.Diagnosis == nil || d.Diagnosis.Name == "" {
				continue
			}
			label := d.Diagnosis.Name
			if d.Severity != nil && *d.Severity != "" {
				label += " (" + *d.Severity + ")"
			}
			names = append(names, label)
		}
		if len(names) > 0 {
			fmt.Fprintf(b, "- Diagnósticos informados por la escuela (son CONTEXTO; no los afirmes como propios, no los repitas si no suma, partí de lo observable): %s\n", strings.Join(names, ", "))
		}
	}

	if ppi := pc.PPI; ppi != nil {
		b.WriteString("PPI (Proyecto Pedagógico Individual):\n")
		if len(ppi.Objectives) > 0 {
			fmt.Fprintf(b, "- Objetivos: %s\n", strings.Join(ppi.Objectives, "; "))
		}
		if ppi.CurricularAdaptations != nil && *ppi.CurricularAdaptations != "" {
			fmt.Fprintf(b, "- Adaptaciones curriculares: %s\n", *ppi.CurricularAdaptations)
		}
		if ppi.FollowUp != nil && *ppi.FollowUp != "" {
			fmt.Fprintf(b, "- Seguimiento: %s\n", *ppi.FollowUp)
		}
	}

	if len(pc.PastAdaptations) > 0 {
		b.WriteString("ADAPTACIONES PREVIAS (no las repitas; construí sobre ellas):\n")
		for i := range pc.PastAdaptations {
			a := &pc.PastAdaptations[i]
			fmt.Fprintf(b, "- %s", a.Title)
			if a.Subject != "" {
				fmt.Fprintf(b, " (%s)", a.Subject)
			}
			if a.Outcome != nil && *a.Outcome != "" {
				fmt.Fprintf(b, " — resultado: %s", *a.Outcome)
			}
			b.WriteString("\n")
		}
	}

	if len(pc.PriorSummaries) > 0 {
		b.WriteString("CHARLAS ANTERIORES (memoria entre clases):\n")
		for i := range pc.PriorSummaries {
			fmt.Fprintf(b, "- %s\n", pc.PriorSummaries[i].Summary)
		}
	}

	writeMissingData(b, pc.MissingData)
	b.WriteString("\n")
}

// writeMissingData enumera, en lenguaje natural, los datos opcionales que faltan
// para que Alizia pueda SUGERIR completarlos (nunca exigirlos ni mostrarlos como "N/A").
func writeMissingData(b *strings.Builder, missing []string) {
	if len(missing) == 0 {
		return
	}
	labels := map[string]string{
		missingTeacherProfile: "el perfil del docente",
		missingStudentProfile: "el perfil del alumno",
		missingPPI:            "el PPI",
		missingDiagnoses:      "los diagnósticos",
	}
	parts := make([]string, 0, len(missing))
	for _, m := range missing {
		if l, ok := labels[m]; ok {
			parts = append(parts, l)
		}
	}
	if len(parts) > 0 {
		fmt.Fprintf(b, "DATOS QUE FALTAN (podés sugerir completarlos, sin exigir): %s.\n", strings.Join(parts, ", "))
	}
}

// buildAssistSystemPrompt arma el prompt de acompañamiento en tiempo real. agentic
// indica si las tools del RAG están disponibles; solo entonces se inyecta FUNDAMENTOS
// (que incluye la cita [CONTENT_ID:X] del corpus).
func buildAssistSystemPrompt(devices []entities.Device, students []entities.Student, pc *PromptContext, agentic bool) string {
	var b strings.Builder

	b.WriteString(aliziaPersona)
	b.WriteString("\nEstás acompañando a un docente DURANTE la clase: sé breve, 1-3 acciones concretas.\n\n")

	b.WriteString("LINEAMIENTOS:\n")
	b.WriteString("- Priorizá adaptar la enseñanza (DUA) por sobre intervenciones individuales.\n")
	b.WriteString("- Liderá con la estrategia pedagógica; el dispositivo es una opción más, no la respuesta.\n")
	b.WriteString("- Si detectás el nombre de un alumno, usá [STUDENT_ID:X]. Si recomendás un dispositivo, usá [DEVICE_ID:X].\n")
	b.WriteString("- Escribí en markdown legible: **negritas** en lo clave, párrafos cortos. Que se lea fácil en pantalla (el docente te lee en plena clase).\n\n")
	b.WriteString("FORMATO DE RESPUESTA CON PASOS:\n")
	b.WriteString("- Cuando proponés acciones concretas (1-3 pasos para implementar en clase), envolvé SOLO esa parte en el bloque [STEPS]...[/STEPS]. Usá lista numerada adentro.\n")
	b.WriteString("- Fuera del bloque: párrafo breve de contexto antes, y si tenés material de fundamento citalo con [CONTENT_ID:X] después.\n")
	b.WriteString("- Si todavía estás preguntando para entender la situación, NO uses [STEPS]: respondé en prosa normal.\n")
	b.WriteString("Ejemplo de turno con pasos:\n")
	b.WriteString("Entiendo, el problema es X. Acá van los pasos para arrancar ahora:\n")
	b.WriteString("[STEPS]\n1. Fragmentá la consigna en dos partes y escribilas en el pizarrón.\n2. Usá el [DEVICE_ID:1] para que pueda seguir el ritmo.\n[/STEPS]\n")
	b.WriteString("Si te parece bien o querés ajustar algo, avisame.\n\n")

	b.WriteString(repreguntaGate)
	b.WriteString("\n")

	if agentic {
		b.WriteString(alumnoFlow)
		b.WriteString("\n")
		b.WriteString(fundamentosRAG)
		b.WriteString("\n")
	}

	writeStudentContext(&b, pc)
	writeKnownStudents(&b, students)

	b.WriteString("DISPOSITIVOS DISPONIBLES:\n")
	writeDeviceCatalog(&b, devices, false)

	b.WriteString("\nGUARDAR COMO RECURSO (bloque estructurado):\n")
	b.WriteString("- Cuando propongas una adaptación concreta, ofrecé guardarla y preguntá si quiere (ej. \"¿Querés que la guarde como recurso?\"). NO incluyas el bloque en ese turno.\n")
	b.WriteString("- Incluí el BLOQUE solo en el turno POSTERIOR, después de que el docente confirme que sí. Nunca en el primer mensaje, ni junto con la pregunta de confirmación, ni en respuestas a consultas o preguntas de aclaración.\n")
	b.WriteString("- Formato exacto, al final del mensaje:\n")
	b.WriteString("[ADAPTATION_JSON:{\"title\":\"título corto\",\"type\":\"tipo\",\"strategy\":\"resumen\",\"ramp_id\":N,\"device_ids\":[1],\"device_names\":[\"nombre\"],\"steps\":[{\"orden\":1,\"texto\":\"primer paso\"}]}]\n")
	b.WriteString("Los tipos válidos son: actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente.\n")
	b.WriteString("ramp_id = categoría/necesidad del catálogo. steps = el PASO A PASO de la guía (lo más importante del recurso), claro y accionable.\n")
	b.WriteString("Si la adaptación no usa material físico, usá estrategia_aula con device_ids vacío.\n")

	return b.String()
}

// buildGuidedAssistPrompt arma el prompt de planificación conversacional. agentic
// indica si las tools del RAG están disponibles; solo entonces se inyecta FUNDAMENTOS.
func buildGuidedAssistPrompt(devices []entities.Device, students []entities.Student, pc *PromptContext, agentic bool) string {
	var b strings.Builder

	b.WriteString(aliziaPersona)
	b.WriteString("\nEl docente quiere planificar una adaptación. Guialo conversacionalmente, sin apurar la propuesta.\n\n")

	b.WriteString("FLUJO GUIADO (una pregunta por vez):\n")
	b.WriteString("1. Para qué alumno es la adaptación (si no lo mencionó).\n")
	b.WriteString("2. Qué materia/actividad están trabajando.\n")
	b.WriteString("3. Qué barrera observable aparece en el aula (concreta, no el diagnóstico).\n")
	b.WriteString("4. Cuando tengas suficiente, generá la propuesta de pasos usando el formato [STEPS] (ver abajo).\n\n")
	b.WriteString("FORMATO DE RESPUESTA CON PASOS:\n")
	b.WriteString("- Una vez que tenés la información necesaria, presentá los pasos concretos SIEMPRE dentro del bloque [STEPS]...[/STEPS]. Usá lista numerada adentro.\n")
	b.WriteString("- Antes del bloque: 1-2 oraciones de cierre que confirmen lo que entendiste. Después: invitá al docente a ajustar si algo no encaja.\n")
	b.WriteString("- Si tenés material de fundamento (del corpus RAG), citalo con [CONTENT_ID:X] fuera del bloque [STEPS], al final.\n")
	b.WriteString("- Mientras estás recopilando información (fases 1-3), NO uses [STEPS]: una sola pregunta en prosa.\n")
	b.WriteString("Ejemplo de turno de propuesta:\n")
	b.WriteString("Entiendo: Lucas tiene dificultades para escribir al dictado por carga cognitiva. Acá van los pasos:\n")
	b.WriteString("[STEPS]\n1. Dictá en fragmentos cortos (4-5 palabras), con pausa entre cada uno.\n2. Escribí la primera palabra en el pizarrón como ancla visual.\n3. Usá el [DEVICE_ID:1] si necesita más tiempo para copiar.\n[/STEPS]\n")
	b.WriteString("Si algo no cierra o querés ajustar, avisame.\n\n")

	b.WriteString(repreguntaGate)
	b.WriteString("\n")

	if agentic {
		b.WriteString(alumnoFlow)
		b.WriteString("\n")
		b.WriteString(fundamentosRAG)
		b.WriteString("\n")
	}

	writeStudentContext(&b, pc)
	writeKnownStudents(&b, students)

	b.WriteString("DISPOSITIVOS DISPONIBLES:\n")
	writeDeviceCatalog(&b, devices, false)

	b.WriteString("\nCuando generes la adaptación final, incluí [STUDENT_ID:X], [DEVICE_ID:X] si aplica, y:\n")
	b.WriteString("[ADAPTATION_JSON:{\"title\":\"título\",\"type\":\"tipo\",\"strategy\":\"resumen\",\"ramp_id\":N,\"device_ids\":[1],\"device_names\":[\"nombre\"],\"steps\":[{\"orden\":1,\"texto\":\"primer paso\"}]}]\n")
	b.WriteString("ramp_id = categoría/necesidad. steps = el PASO A PASO de la guía (lo más importante del recurso). Tipos válidos: actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente. Sin material físico, usá estrategia_aula con device_ids vacío.\n")

	return b.String()
}

// buildSummaryPrompt instruye el resumen de una conversación para guardar memoria
// entre clases (lo consume el cron de summaries). Pide SOLO JSON, sin prosa.
func buildSummaryPrompt() string {
	return `Resumís conversaciones entre Alizia (asistente de inclusión educativa) y un docente, para guardar memoria entre clases.
Devolvé SOLO un JSON válido, sin texto alrededor ni fences, con esta forma exacta:
{"summary":"...","topic_keywords":["...","..."]}
- summary: 2-4 oraciones en español rioplatense, foco pedagógico: qué alumno/barrera/actividad se trabajó, qué se propuso y en qué quedó. Concreto, sin saludos ni relleno.
- topic_keywords: 3-8 palabras o locuciones clave (temas, barreras, dispositivos), en minúscula.
No incluyas marcadores [STUDENT_ID]/[DEVICE_ID] ni IDs crudos.`
}

var deviceIDRegex = regexp.MustCompile(`\[DEVICE_ID:(\d+)\]`)
var studentIDRegex = regexp.MustCompile(`\[STUDENT_ID:(\d+)\]`)
var adaptationJSONRegex = regexp.MustCompile(`\[ADAPTATION_JSON:(\{.+\})\]`)

func extractDeviceID(content string) *int64 {
	matches := deviceIDRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	id, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil
	}
	return &id
}

func extractStudentID(content string) *int64 {
	matches := studentIDRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	id, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil
	}
	return &id
}

func extractAdaptationJSON(content string) *GeneratedAdaptation {
	matches := adaptationJSONRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	var adaptation GeneratedAdaptation
	if err := json.Unmarshal([]byte(matches[1]), &adaptation); err != nil {
		return nil
	}
	return &adaptation
}

var (
	multiSpaceRegex       = regexp.MustCompile(`[ \t]{2,}`)
	spaceBeforePunctRegex = regexp.MustCompile(`[ \t]+([,.;:!?)])`)
)

// stripInternalMarkers quita los marcadores internos ([STUDENT_ID:X], [DEVICE_ID:X],
// [ADAPTATION_JSON:{...}]) del texto del modelo ANTES de mostrarlo al docente o
// persistirlo. Los ids/JSON ya se extrajeron aparte: estos tags son internos del
// backend y nunca deben aparecer en el chat. Limpia los espacios que deja el borrado.
// Lo usa el flujo recommend, cuyo render en el FE no convierte markers en chips.
func stripInternalMarkers(content string) string {
	content = studentIDRegex.ReplaceAllString(content, "")
	content = deviceIDRegex.ReplaceAllString(content, "")
	content = adaptationJSONRegex.ReplaceAllString(content, "")
	content = multiSpaceRegex.ReplaceAllString(content, " ")
	content = spaceBeforePunctRegex.ReplaceAllString(content, "$1")
	return strings.TrimSpace(content)
}

// stripAdaptationBlock quita SOLO el bloque [ADAPTATION_JSON:{...}] (ya extraído a un
// campo estructurado). Lo usa el assist: a diferencia de stripInternalMarkers, deja
// pasar [STUDENT_ID:X]/[DEVICE_ID:X]/[CONTENT_ID:X] porque el FE los renderiza como
// chips (nombre del alumno, material o título del contenido), nunca como id crudo.
func stripAdaptationBlock(content string) string {
	content = adaptationJSONRegex.ReplaceAllString(content, "")
	content = multiSpaceRegex.ReplaceAllString(content, " ")
	content = spaceBeforePunctRegex.ReplaceAllString(content, "$1")
	return strings.TrimSpace(content)
}
