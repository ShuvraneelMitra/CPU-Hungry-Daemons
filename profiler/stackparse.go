package profiler

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

/*
Credits to DataDog/gostackparse; I read the code once before starting to
write my own. The design of the parser remains similar, though the implementation
is purely mine.
*/

const (
	heading_prefix = "goroutine "
	created_by_prefix = "created by "
)

func firstWord(b []byte) []byte {
	for i, c := range b {
		if (c < 'A' || c > 'Z') &&
		   (c < 'a' || c > 'z') &&
		   (c < '0' || c > '9') {
			return b[:i]
		}
	}
	return b
}


type GoRoutineState int

const (
	RUNNING GoRoutineState = iota
	RUNNABLE

	BLOCK_SLEEP // 2
	BLOCK_CHAN_SEND
	BLOCK_CHAN_RCV
	BLOCK_SELECT
	BLOCK_MUTEX
	BLOCK_SEM
	BLOCK_IO

	WAIT
	SLEEP // 10
	SYSCALL
	GC
	DEAD
)

var stateMap = map[string]GoRoutineState{
	"running": RUNNING,
	"runnable": RUNNABLE,
	"syscall": SYSCALL,
	"waiting": WAIT,
	"dead": DEAD,

	"chan send": BLOCK_CHAN_SEND,
	"chan receive": BLOCK_CHAN_RCV,
	"select": BLOCK_SELECT,
	"sync.Mutex.Lock": BLOCK_MUTEX,
	"sync.RWMutex.RLock": BLOCK_MUTEX,
	"sync.WaitGroup.Wait": WAIT,
	"time.Sleep": SLEEP,
	"IO wait": BLOCK_IO,
	"GC assist wait": GC,
	"GC sweep wait": GC,
	"GC worker (idle)": GC,
	"finalizer wait": WAIT,
}

type Frame struct {
	Func string
	File string
	Line       int
}

type GoRoutine struct {
	Id             int
	State          GoRoutineState
	Waiting        bool
	Stack          []Frame
	CreatedBy      *Frame
}

type Sample struct {
	Timestamp      time.Time
	GoRoutineCount int
	List           []GoRoutine
}

type ParserState int

const (
	StateHeading ParserState = iota
	StateStackFunc
	StateStackFuncAddr
	StateCreatedByFunc
	StateCreatedByAddr
)

func NewSample(t time.Time, nGo int) *Sample {
	return &Sample{
		Timestamp:      t,
		GoRoutineCount: nGo,
		List:           make([]GoRoutine, 0),
	}
}

func Parse(stop <-chan any, dataStream <-chan Metadata) <-chan Sample {

	/*
		The expected behaviour is to block the parsing goroutine until it
		receives a Metadata object from the dataStream created by the sampler.
	*/
	parsedStream := make(chan Sample)

	go func() {
		defer close(parsedStream)
		for {
			select {
			case <-stop:
				return
			case metadata, ok := <-dataStream:
				if ok == false {
					return
				}
				
				sample := NewSample(metadata.Timestamp, metadata.numGoroutines)
				var cur_state ParserState
				next_state := StateHeading

				for line := range bytes.SplitSeq(metadata.stackDump, []byte("\n")) {
					if len(line) == 0 {
						continue
					}

					fmt.Println(string(line))
					cur_state = next_state // Having 2 such states per iteration is something I picked up from my Verilog days

					state_machine:
						switch cur_state {
						case StateHeading:
							if !bytes.HasPrefix(line, []byte(heading_prefix)) {
								next_state = StateStackFunc // retry with a different state
								goto state_machine
							}
							if !parseHeading(sample, line) {
								next_state = StateHeading // look for heading in next line
								continue
							}
							next_state = StateStackFunc
						
						case StateStackFunc:
							if !parseStackFunc(sample, line) {
								next_state = StateCreatedByFunc // retry with a different state
								goto state_machine
							}
							next_state = StateStackFuncAddr

						case StateStackFuncAddr:
							if !parseStackFuncAddr(sample, line) {
								next_state = StateCreatedByAddr // retry with a different state
								goto state_machine
							}
							next_state = StateStackFunc

						case StateCreatedByFunc:
							if !bytes.HasPrefix(line, []byte(created_by_prefix)) {
								next_state = StateHeading // retry with a different state
								goto state_machine
							}

							if !parseCreatedByFunc(sample, line) {
								next_state = StateHeading // retry with a different state
								goto state_machine
							}
							next_state = StateCreatedByAddr

						case StateCreatedByAddr:
							if !parseCreatedByAddr(sample, line) {
								next_state = StateHeading // retry with a different state
								goto state_machine
							}
							next_state = StateHeading
						}
				}
				fmt.Println(sample.GoRoutineCount)
				parsedStream<- *sample
			}
		}
	}()

	return parsedStream
}

func parseHeading(sample *Sample, line []byte) bool {
	//-------------------------------- ID PARSING -----------------------------------
	/*
	We have already determined at the state machine level that the line does 
	begin with "goroutine ..." and so we need not check it once again.
	*/
	
	line = bytes.TrimSpace(line[len(heading_prefix) : ])
	id := bytes.Split(line, []byte(" "))[0]

	newGoRoutine := GoRoutine{
		Stack: make([]Frame, 0),
		CreatedBy: &Frame{},
	}
	idVal, err := strconv.Atoi(string(id))
	if err != nil {
		fmt.Println("[STACK PARSER] Could not parse header pattern:", err)
		return false
	}
	newGoRoutine.Id = idVal


	// ------------------------- GOROUTINE STATE PARSING ------------------------------
	
	idx_open := bytes.LastIndex(line, []byte("["))
	idx_close := bytes.LastIndex(line, []byte("]"))
	stateString := line[idx_open + 1 : idx_close]


	newGoRoutine.State = stateMap[string(firstWord(stateString))]
	if 2 <= newGoRoutine.State && newGoRoutine.State <= 10 {
		newGoRoutine.Waiting = true
	}
	sample.List = append(sample.List, newGoRoutine)
	return true
}

func parseStackFunc(sample *Sample, line []byte) bool {

	idx := bytes.LastIndex(line, []byte("("))
	if idx == -1 {
		return false
	}

	stack_function := line[ : idx]

	latestGoRoutine := &sample.List[len(sample.List) - 1]
	latestGoRoutine.Stack = append(latestGoRoutine.Stack, Frame{
		Func: string(stack_function),
	})
	return true
}

func parseStackFuncAddr(sample *Sample, line []byte) bool {
	/*
	We want to support absolute paths with spaces in intermediate or final
	file/folder names, which means that we will utilise LastIndex multiple
	times in this function.
	*/
	splitfrom := bytes.LastIndex(line, []byte("+"))
	if splitfrom == -1 {
		return false
	}
	line = line[ : splitfrom]

	idx := bytes.LastIndex(line, []byte(":"))
	if idx == -1 {
		return false
	}

	filepath := line[ : idx]
	line_num := line[idx + 1 : ]

	if len(sample.List) == 0 {
		fmt.Println("[STACK PARSER] Could not parse header pattern: No header detected")
		return false
	}
	latestGoRoutine := &sample.List[len(sample.List) - 1]

	if len(latestGoRoutine.Stack) == 0 {
		fmt.Println("[STACK PARSER] Could not parse header pattern: No stack function detected")
		return false
	}
	latestFrame := &latestGoRoutine.Stack[len(latestGoRoutine.Stack) - 1]

	latestFrame.File = string(filepath)
	latestFrame.Line, _ = strconv.Atoi(string(line_num))

	return true
}

func parseCreatedByFunc(sample *Sample, line []byte) bool {
	createdBy := bytes.TrimSpace(line[len(created_by_prefix) : ])
	
	if len(sample.List) == 0 {
		fmt.Println("[STACK PARSER] Could not parse header pattern: No header detected")
		return false
	}
	latestGoRoutine := &sample.List[len(sample.List) - 1]

	latestGoRoutine.CreatedBy.File = string(createdBy)
	return true
}

func parseCreatedByAddr(sample *Sample, line []byte) bool {
	/*
	Again, we want to support absolute paths with spaces in intermediate or final
	file/folder names.
	*/
	splitfrom := bytes.LastIndex(line, []byte("+"))
	if splitfrom == -1 {
		return false
	}
	line = line[ : splitfrom]

	idx := bytes.LastIndex(line, []byte(":"))
	if idx == -1 {
		return false
	}

	filepath := line[ : idx]
	line_num := line[idx + 1 : ]

	if len(sample.List) == 0 {
		fmt.Println("[STACK PARSER] Could not parse header pattern: No header detected")
		return false
	}
	latestGoRoutine := &sample.List[len(sample.List) - 1]

	latestFrame := latestGoRoutine.CreatedBy

	latestFrame.File = string(filepath)
	latestFrame.Line, _ = strconv.Atoi(string(line_num))

	return true
}


/*
DESIGN MUSINGS: Is a state machine the right approach? It is certainly more elegant,
given that in a file with a given structure we always know what to expect next, and
hence can skip having to check multiple "if" statements on strings which seem slightly
contrived or "hardcoded", for lack of a better word.

The gostackparse parser uses a lot more states than I feel are necessary, and I instead
assign one state for each line that is consumed, instead of one state for each semantic
symbol that is consumed.
*/
