### **Core Concepts to Learn**  
#### **1. Digital Biology Basics**  
- **Artificial Life**: Study Conway’s *Game of Life*, *Tierra*, and *Avida* (digital evolution simulators).  
- **Genetic Algorithms (GA)**:  
  - How "organisms" (goroutines) encode traits (e.g., CPU usage, replication rate).  
  - *Fitness functions*: What determines "survival" (e.g., resource efficiency).  
- **Mutation/Selection**: Random changes + competition for resources (CPU, memory).  

#### **2. Go-Specific Concurrency Patterns**  
- **Goroutines as "Organisms"**:  
  - Each goroutine is an independent entity with a lifecycle.  
- **Channels for Communication**:  
  - How organisms "interact" (e.g., compete, cooperate, or share resources).  
- **Resource Contention**:  
  - Use `sync.Mutex` or `semaphores` to simulate limited CPU/memory.  

#### **3. Evolutionary Systems Design**  
- **Genotype → Phenotype Mapping**:  
  - Define how "code DNA" (e.g., a struct) translates to behavior.  
- **Reproduction Logic**:  
  - How organisms "replicate" (fork goroutines with mutations).  
- **Death Conditions**:  
  - When should a goroutine die? (e.g., starvation, timeout).  

#### **4. Observability & Metrics**  
- **pprof + Metrics**:  
  - Track CPU/memory per "organism".  
- **Visualization**:  
  - Use **Pixel (2d library)**, **Terminal UI (BubbleTea)** or **Web (WASM + WebGL)** to watch evolution in real-time.  

---

### **Technical Dependencies to Explore**  
- **Go Libraries**:  
  - `runtime` (for goroutine introspection).  
  - `github.com/gizak/termui` (for TUI dashboards).  
- **Math/Stats**:  
  - Probability distributions (for mutation rates).  
  - Evolutionary stable strategies (ESS) for balancing traits.  

---

### **Phased Approach (No Code Yet)**  
#### **Phase 1: Define the "Organism"**  
- What **data** represents an organism? (e.g., struct with `lifespan`, `mutationRate`).  
- What **actions** can it take? (e.g., `reproduce()`, `consumeCPU()`).  

#### **Phase 2: Simulate the Environment**  
- **Resource Pools**:  
  - Global CPU "food" (e.g., a channel emitting work tokens).  
- **Selection Pressure**:  
  - How does the environment kill organisms? (e.g., starvation, predation).  

#### **Phase 3: Evolutionary Mechanics**  
- **Mutation Rules**:  
  - How are child organisms different from parents? (e.g., random `mutationRate` tweak).  
- **Fitness Scoring**:  
  - What makes an organism "successful"? (e.g., replicates most).  

#### **Phase 4: Emergent Behaviors**  
- **Hypothesize**:  
  - Will "parasitic" organisms evolve? (e.g., goroutines that steal resources).  
  - Will cooperation emerge? (e.g., organisms that share CPU).  

---

### **Inspiration & Further Reading**  
- **Books**:  
  - *"The Selfish Gene"* (Dawkins) – Evolutionary biology basics.  
  - *"Complexity: A Guided Tour"* (Mitchell) – Emergent systems.  
- **Papers**:  
  - [*"Avida: A Software Platform for Research in Computational Evolutionary Biology"*](https://doi.org/10.1038/s41596-020-0293-9).  
- **Videos**:  
  - [*"Digital Evolution: Creating Artificial Life"* (Youtube)](https://www.youtube.com/watch?v=oCXzcPNsqGA).  

---

### **Next Steps**  
1. **Sketch a design doc** answering:  
   - What’s the simplest "organism" you can model?  
   - How will you measure "evolution"?  
2. **Prototype subsystems**:  
   - Start with a single goroutine that replicates itself.  
   - Add resource constraints later.  

When you’re ready, we’ll dive into **Go-specific implementation strategies**. Sound good?