package inspection

import (
	"net/http"
)

type RestaritemInspection struct {
	ID   string
	Name string

	Inspection
}

type Button struct {
	Number int
	Active string
	Post   string
}

func (s RestaritemInspection) Buttons(restaritemID int) []Button {
	var target []Button

	switch s.Type {
	case DefaultType:
		//for num := 1; num < 6; num++ {
		//	btn := Button{
		//		Number: num,
		//	}
		//
		//	if num == s.Quality {
		//		btn.Active = "active"
		//	}
		//
		//	// id is restaritemID,
		//	btn.Post = fmt.Sprintf("%d/inspection/%s/%d", restaritemID, s.ID, num)
		//
		//	target = append(target, btn)
		//}
		//
		//return target

	case BoolType:

	default:
		println("unknown type")
	}

	return target

}

func SetInspectionByID(w http.ResponseWriter, r *http.Request) {
	//files := []string{
	//	"./web/template/inspection_view.html",
	//}

	//tmpl, err := template.ParseFiles(files...)
	//if err != nil {
	//	pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")
	//
	//	return
	//}
	//
	//a, ok := RestaritemInspections[inspID]
	//if !ok {
	//	pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "cant find inspection")
	//
	//	return
	//}
	//
	//a.Quality = rating
	//starterInspections[inspID] = a
	//log.Printf("starterInspections: %+v", starterInspections)
	//log.Printf("set rat: %+v", rating)
	//log.Print(rating, inspID, id)
	//
	//if err = tmpl.Execute(w, pkg.JSON{"id": id, "inspections": starterInspections}); err != nil {
	//	pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")
	//
	//	return
	//}
}
