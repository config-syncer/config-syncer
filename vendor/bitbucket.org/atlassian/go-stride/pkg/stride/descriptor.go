package stride

type Descriptor struct {
	BaseURL   string                  `json:"baseUrl"`
	Key       string                  `json:"key"`
	LifeCycle *LifeCycle              `json:"lifecycle"`
	Modules   map[ModuleType][]Module `json:"modules"`
}
