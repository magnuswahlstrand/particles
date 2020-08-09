package main

//type FloatWidget struct {
//	label string
//
//	f0     float32
//	f1, f2 float32
//
//	target   *generators.Float
//	min, max float32
//	index    int
//}
//
//func (w *FloatWidget) show() {
//	imgui.PushID(w.label)
//	if imgui.BeginComboV("##hidelabel", "", imgui.ComboFlagNoPreview) {
//		if imgui.Selectable("Constant") {
//			w.index = 0
//			w.update()
//		}
//		if imgui.Selectable("Random Between Two Constants") {
//			w.index = 1
//			w.update()
//		}
//
//		imgui.EndCombo()
//	}
//	imgui.SameLine()
//
//	switch w.index {
//	case 0:
//		if imgui.DragFloatV(w.label, &w.f0, 0.1, w.min, w.max, "%0.1f", 1) {
//			*w.target = generators.FloatConstant{float64(w.f0)}
//		}
//	case 1:
//		imgui.BeginGroup()
//		imgui.PushItemWidth(100)
//
//		f1changed := imgui.DragFloatV("##hidelabel", &w.f1, 0.1, w.min, w.max, "%0.1f", 1)
//		imgui.SameLine()
//		f2changed := imgui.DragFloatV(w.label+"##hidelabel", &w.f2, 0.1, w.min, w.max, "%0.1f", 1)
//
//		if f1changed || f2changed {
//			*w.target = generators.FloatRandomBetweenTwoConstants{float64(w.f1), float64(w.f2)}
//		}
//
//		imgui.PopItemWidth()
//		imgui.EndGroup()
//	}
//
//	imgui.PopID()
//}
//
//func (w *FloatWidget) update() {
//	switch w.index {
//	case 0:
//		*w.target = generators.FloatConstant{float64(w.f0)}
//	case 1:
//		*w.target = generators.FloatRandomBetweenTwoConstants{float64(w.f1), float64(w.f2)}
//	}
//}
//
//type ColorWidget struct {
//	label string
//
//	c0     [4]float32
//	c1, c2 [4]float32
//
//	target *generators.Color
//	index  int
//}
//
//func (w *ColorWidget) show() {
//	imgui.PushID(w.label)
//	if imgui.BeginComboV("##hidelabel", "", imgui.ComboFlagNoPreview) {
//		if imgui.Selectable("Constant") {
//			w.index = 0
//			w.update()
//		}
//		if imgui.Selectable("Random Between Two Colors") {
//			w.index = 1
//			w.update()
//		}
//
//		imgui.EndCombo()
//	}
//	imgui.SameLine()
//
//	switch w.index {
//	case 0:
//		if imgui.ColorEdit4("StartColor", &w.c0) {
//			*w.target = generators.ColorConstant{colorFromArray(w.c0)}
//		}
//	case 1:
//		imgui.BeginGroup()
//
//		c1changed := imgui.ColorEdit4("##hidelabel", &w.c1)
//		c2changed := imgui.ColorEdit4(w.label+"##hidelabel", &w.c2)
//
//		if c1changed || c2changed {
//			*w.target = generators.ColorRandomBetweenTwoConstants{colorFromArray(w.c1), colorFromArray(w.c2)}
//		}
//
//		imgui.EndGroup()
//	}
//
//	imgui.PopID()
//}
//
//func (w *ColorWidget) update() {
//	switch w.index {
//	case 0:
//		*w.target = generators.ColorConstant{colorFromArray(w.c0)}
//	case 1:
//		*w.target = generators.ColorRandomBetweenTwoConstants{colorFromArray(w.c1), colorFromArray(w.c2)}
//	}
//}
