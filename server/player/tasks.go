package player

import(
	"github.com/perdokcat/TermoTune/shared"
)

func (p * Player) AddTask(target string, typeTask int){
		p.mu.Lock()
		defer p.mu.Unlock()
		p.Tasks[target] = shared.Task{
			Type : typeTask,
			Error: "",
		}
}


func(p * Player) RemoveTask(target string){
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Tasks, target)
}

func(p * Player) ErrorTask(target string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock() 
	task, ok := p.Tasks[target] 
	if ok {
		task.Error = err.Error()
		p.Tasks[target] = task
	}
}

