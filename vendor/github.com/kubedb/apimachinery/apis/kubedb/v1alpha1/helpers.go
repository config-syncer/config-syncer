package v1alpha1

func (e *BackupScheduleSpec) SetDefaults() {
	if e == nil {
		return
	}
	if e.Resources != nil {
		e.PodTemplate.Spec.Resources = *e.Resources
		e.Resources = nil
	}
}
