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
	if a.StudentID == nil {
		return "Situación general"
	}
	return fmt.Sprintf("Alumno #%d", *a.StudentID)
}

func teacherName(a *entities.Adaptation) string {
	if a.Teacher != nil && a.Teacher.Name != "" {
		return a.Teacher.Name
	}
	return ""
}

// studentLine arma "Nombre, Curso" para el encabezado del recurso (ej "Martina Gómez, 3°B").
func studentLine(a *entities.Adaptation) string {
	name := studentName(a)
	if a.Student != nil && a.Student.GradeLevel != nil && strings.TrimSpace(*a.Student.GradeLevel) != "" {
		return fmt.Sprintf("%s, %s", name, strings.TrimSpace(*a.Student.GradeLevel))
	}
	return name
}

// categoryName devuelve la categoría/necesidad del recurso si está cargada.
func categoryName(a *entities.Adaptation) string {
	if a.Ramp != nil {
		return a.Ramp.Name
	}
	return ""
}

func resourceTitle(a *entities.Adaptation) string {
	if strings.TrimSpace(a.Title) != "" {
		return a.Title
	}
	if a.Subject != "" {
		return a.Subject
	}
	return "Recurso"
}

func renderAdaptationMarkdown(a *entities.Adaptation) []byte {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", resourceTitle(a))
	fmt.Fprintf(&b, "**Alumno:** %s  \n", studentLine(a))
	if t := teacherName(a); t != "" {
		fmt.Fprintf(&b, "**Docente:** %s  \n", t)
	}
	if c := categoryName(a); c != "" {
		fmt.Fprintf(&b, "**Categoría:** %s  \n", c)
	}
	fmt.Fprintf(&b, "**Estado:** %s\n\n", statusLabel(a.Status))

	if s := deref(a.ActivityDescription); s != "" {
		fmt.Fprintf(&b, "## Situación trabajada\n\n%s\n\n", s)
	}
	if s := deref(a.AdaptationStrategy); s != "" {
		fmt.Fprintf(&b, "## Adaptación\n\n%s\n\n", s)
	}

	if len(a.Steps) > 0 {
		b.WriteString("## Paso a paso\n\n")
		for i := range a.Steps {
			st := &a.Steps[i]
			box := ""
			if st.Checkbox {
				box = "[ ] "
			}
			fmt.Fprintf(&b, "%d. %s%s\n", i+1, box, strings.TrimSpace(st.Texto))
		}
		b.WriteString("\n")
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
	pdf.MultiCell(contentWidth, 9, tr(resourceTitle(a)), "", "L", false)
	pdf.Ln(2)

	field("Alumno", studentLine(a))
	field("Docente", teacherName(a))
	field("Categoría", categoryName(a))
	field("Estado", statusLabel(a.Status))

	if s := deref(a.ActivityDescription); s != "" {
		heading("Situación trabajada")
		body(s)
	}
	if s := deref(a.AdaptationStrategy); s != "" {
		heading("Adaptación")
		body(s)
	}

	if len(a.Steps) > 0 {
		heading("Paso a paso")
		for i := range a.Steps {
			st := &a.Steps[i]
			prefix := fmt.Sprintf("%d. ", i+1)
			if st.Checkbox {
				prefix = fmt.Sprintf("%d. [ ] ", i+1)
			}
			body(prefix + strings.TrimSpace(st.Texto))
		}
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
