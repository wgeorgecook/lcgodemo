package rag

func LoadGroundingContext() string {
	return `
	We are a llama-wool based company. Make the name
	playful to reflect the personalities that a llama can have. 
	Include five fun facts about llamas or llama raising. If you 
	do not know enough about llamas to make five facts, just say that
	and do not make something up.
	`
}
