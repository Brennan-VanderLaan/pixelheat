// Package main provides a collection of models and services for natural language processing.
package main

// Service represents a natural language processing service.
type Service struct {
	ModelName    string  // The name of the model used for the service.
	Context      int     // The context size of the model.
	InputCost    float64 // The cost of input for the service.
	OutputCost   float64 // The cost of output for the service.
	TrainingCost float64 // The cost of training the model for the service. Only applicable for fine-tuning models.
}

// ModelType represents a type of natural language processing model.
type ModelType struct {
	Name     string    // The name of the model type.
	Services []Service // The services provided by the model type.
}

// models is a collection of natural language processing models.
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

// GetServiceCost returns the input and output cost of a specific natural language processing service.
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

// GetService returns a pointer to a specific natural language processing service.
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
