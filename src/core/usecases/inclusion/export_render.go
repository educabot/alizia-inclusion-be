package inclusion

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

var statusLabels = map[string]string{
	"en_curso":     "En curso",
	"probado":      "Probado",
	"funciono":     "Funcionó",
	"para_ajustar": "Para ajustar",
}

func statusLabel(status string) string {
	if label, ok := statusLabels[status]; ok {
		return label
	}
	return status
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}

func studentName(a *entities.Adaptation) string {
	if a.Student != nil && a.Student.Name != "" {
		return a.Student.Name
	}
	return fmt.Sprintf("Alumno #%d", a.StudentID)
}

func teacherName(a *entities.Adaptation) string {
	if a.Teacher != nil && a.Teacher.Name != "" {
		return a.Teacher.Name
	}
	return ""
}

func renderAdaptationMarkdown(a *entities.Adaptation) []byte {
	var b strings.Builder

	title := a.Subject
	if title == "" {
		title = "Adaptación"
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "**Alumno:** %s  \n", studentName(a))
	if t := teacherName(a); t != "" {
		fmt.Fprintf(&b, "**Docente:** %s  \n", t)
	}
	fmt.Fprintf(&b, "**Tipo:** %s  \n", a.AdaptationType)
	fmt.Fprintf(&b, "**Estado:** %s\n\n", statusLabel(a.Status))

	if s := deref(a.ActivityDescription); s != "" {
		fmt.Fprintf(&b, "## Actividad\n\n%s\n\n", s)
	}
	if s := deref(a.AdaptationStrategy); s != "" {
		fmt.Fprintf(&b, "## Estrategia\n\n%s\n\n", s)
	}

	devices := adaptationDevices(a)
	if len(devices) > 0 {
		b.WriteString("## Dispositivos sugeridos\n\n")
		for i := range devices {
			d := &devices[i]
			fmt.Fprintf(&b, "### %s\n\n", d.Name)
			if s := deref(d.Rationale); s != "" {
				fmt.Fprintf(&b, "- **Por qué:** %s\n", s)
			}
			if s := deref(d.HowToUse); s != "" {
				fmt.Fprintf(&b, "- **Cómo usarlo:** %s\n", s)
			}
			b.WriteString("\n")
		}
	}

	if s := deref(a.Notes); s != "" {
		fmt.Fprintf(&b, "## Notas para el docente\n\n%s\n\n", s)
	}
	if s := deref(a.Outcome); s != "" {
		fmt.Fprintf(&b, "## Resultado\n\n%s\n\n", s)
	}

	b.WriteString("---\n\n")
	fmt.Fprintf(&b, "_Generado por Alizia · Educabot · %s · Adaptación #%d_\n",
		a.CreatedAt.Format("02/01/2006"), a.ID)

	return []byte(b.String())
}

func adaptationDevices(a *entities.Adaptation) []entities.Device {
	if len(a.Devices) > 0 {
		return a.Devices
	}
	if a.Device != nil {
		return []entities.Device{*a.Device}
	}
	return nil
}

func renderAdaptationPDF(a *entities.Adaptation) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	const contentWidth = 170.0

	heading := func(text string) {
		pdf.Ln(2)
		pdf.SetFont("Arial", "B", 13)
		pdf.MultiCell(contentWidth, 7, tr(text), "", "L", false)
		pdf.Ln(1)
	}
	body := func(text string) {
		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(contentWidth, 6, tr(text), "", "L", false)
	}
	field := func(label, value string) {
		if value == "" {
			return
		}
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 6, tr(label+": "), "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(0, 6, tr(value), "", "L", false)
	}

	pdf.SetFont("Arial", "B", 18)
	title := a.Subject
	if title == "" {
		title = "Adaptación"
	}
	pdf.MultiCell(contentWidth, 9, tr(title), "", "L", false)
	pdf.Ln(2)

	field("Alumno", studentName(a))
	field("Docente", teacherName(a))
	field("Tipo", a.AdaptationType)
	field("Estado", statusLabel(a.Status))

	if s := deref(a.ActivityDescription); s != "" {
		heading("Actividad")
		body(s)
	}
	if s := deref(a.AdaptationStrategy); s != "" {
		heading("Estrategia")
		body(s)
	}

	devices := adaptationDevices(a)
	if len(devices) > 0 {
		heading("Dispositivos sugeridos")
		for i := range devices {
			d := &devices[i]
			pdf.SetFont("Arial", "B", 11)
			pdf.MultiCell(contentWidth, 6, tr(d.Name), "", "L", false)
			if s := deref(d.Rationale); s != "" {
				body("Por qué: " + s)
			}
			if s := deref(d.HowToUse); s != "" {
				body("Cómo usarlo: " + s)
			}
			pdf.Ln(1)
		}
	}

	if s := deref(a.Notes); s != "" {
		heading("Notas para el docente")
		body(s)
	}
	if s := deref(a.Outcome); s != "" {
		heading("Resultado")
		body(s)
	}

	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 9)
	footer := fmt.Sprintf("Generado por Alizia · Educabot · %s · Adaptación #%d",
		footerDate(a.CreatedAt), a.ID)
	pdf.MultiCell(contentWidth, 5, tr(footer), "", "C", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func footerDate(t time.Time) string {
	if t.IsZero() {
		return time.Now().Format("02/01/2006")
	}
	return t.Format("02/01/2006")
}
