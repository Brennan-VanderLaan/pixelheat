package main

type AIAgent struct {
	Name     string
	Services []*Service
}

var aiAgents = []AIAgent{
	{
		Name: "Agent1",
		Services: []*Service{
			GetService("gpt-4", "gpt-4"),
			GetService("gpt-3.5", "gpt-3.5-turbo"),
		},
	},
	{
		Name: "Agent2",
		Services: []*Service{
			GetService("gpt-4", "gpt-4-32k"),
			GetService("gpt-3.5", "gpt-3.5-turbo-16k"),
		},
	},
	// ... add more agents as needed
}
