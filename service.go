package main

type Service struct {
	ModelName    string
	Context      int
	InputCost    float64
	OutputCost   float64
	TrainingCost float64 // Only for fine-tuning models
}

type ModelType struct {
	Name     string
	Services []Service
}

var models = []ModelType{
	{
		Name: "gpt-4",
		Services: []Service{
			{ModelName: "gpt-4", Context: 8192},
			{ModelName: "gpt-4-0613", Context: 8192},
			{ModelName: "gpt-4-32k", Context: 32768},
			{ModelName: "gpt-4-32k-0613", Context: 32768},
			{ModelName: "gpt-4-0314 (Legacy)", Context: 8192},
			{ModelName: "gpt-4-32k-0314 (Legacy)", Context: 32768},
		},
	},
	{
		Name: "gpt-3.5",
		Services: []Service{
			{ModelName: "gpt-3.5-turbo", Context: 4096},
			{ModelName: "gpt-3.5-turbo-16k", Context: 16384},
			{ModelName: "gpt-3.5-turbo-0613", Context: 4096},
			{ModelName: "gpt-3.5-turbo-16k-0613", Context: 16384},
			{ModelName: "gpt-3.5-turbo-0301 (Legacy)", Context: 4096},
			{ModelName: "text-davinci-003 (Legacy)", Context: 4097},
			{ModelName: "text-davinci-002 (Legacy)", Context: 4097},
			{ModelName: "code-davinci-002 (Legacy)", Context: 8001},
		},
	},
	{
		Name: "Fine-tuning models",
		Services: []Service{
			{ModelName: "babbage-002", TrainingCost: 0.0004, InputCost: 0.0016, OutputCost: 0.0016},
			{ModelName: "davinci-002", TrainingCost: 0.0060, InputCost: 0.0120, OutputCost: 0.0120},
			{ModelName: "GPT-3.5 Turbo", TrainingCost: 0.0080, InputCost: 0.0120, OutputCost: 0.0160},
		},
	},
	{
		Name: "Embedding models",
		Services: []Service{
			{ModelName: "Ada v2", InputCost: 0.0001},
		},
	},
	// ... Add Image models similarly
}

// Example function to get the cost of a specific model and context
func GetServiceCost(modelName, context string) (float64, float64) {
	for _, model := range models {
		if model.Name == modelName {
			for _, service := range model.Services {
				if service.ModelName == context {
					return service.InputCost, service.OutputCost
				}
			}
		}
	}
	return 0, 0 // Return 0 if the model or context is not found
}

func GetService(modelName, context string) *Service {
	for _, model := range models {
		if model.Name == modelName {
			for _, service := range model.Services {
				if service.ModelName == context {
					return &service
				}
			}
		}
	}
	return nil // Return nil if the model or context is not found
}
