### **Problem Statement: Digital Evolution Simulator ("Self-Replicating Code Organisms")**  

#### **Objective**  
Create a **closed computational ecosystem** where self-replicating, mutating **goroutines** ("digital organisms") compete for limited system resources (CPU, memory), evolve over time, and exhibit emergent behaviors such as parasitism, cooperation, or optimization.  

---

### **Core Components & Dynamics**  

#### **1. The Organisms**  
- **Representation**:  
  - Each organism is a **goroutine** with a "genome" (a struct defining traits like `replicationRate`, `cpuUsage`, `lifespan`).  
  - Example genome:  
    ```go
    type Genome struct {
        ID          string  // Unique identifier
        ReplicationRate float64 // How often it replicates
        CPUHunger      float64 // CPU cycles consumed per tick
        MutationChance float64 // Chance of mutation in offspring
    }
    ```  
- **Behavior**:  
  - **Replicate**: Fork a new goroutine (child) with a slightly mutated genome.  
  - **Metabolize**: Consume CPU cycles (simulated by performing work).  
  - **Die**: If it exceeds `lifespan` or fails to acquire resources.  

#### **2. The Environment**  
- **Resource Constraints**:  
  - **CPU Pool**: A shared channel (`chan struct{}`) that emits "CPU tokens" (like a semaphore). Organisms must acquire tokens to execute.  
  - **Memory Pressure**: Global counter tracking active organisms; if exceeded, oldest organisms die.  
- **Selection Pressures**:  
  - **Starvation**: Organisms with high `CPUHunger` die if they can’t acquire tokens.  
  - **Predation**: Optional—organisms can "kill" others by stealing their resources.  

#### **3. Evolutionary Mechanics**  
- **Mutation**:  
  - Child organisms inherit parent’s genome with small random changes (e.g., `ReplicationRate += rand.NormFloat64() * 0.1`).  
- **Fitness**:  
  - Survival = Ability to replicate before dying.  
  - Indirect fitness: Organisms leaving the most descendants dominate.  
- **Emergent Scenarios**:  
  - **Parasitism**: Organisms evolve to hijack others’ resources.  
  - **Cooperation**: Organisms might "share" CPU tokens (unlikely but possible).  

---

### **Key Challenges**  

#### **1. Goroutine Management**  
- How to limit goroutines without crashing the host machine?  
  - **Solution**: Use a `sync.Pool` or weighted semaphore.  
- How to track organism states without global locks?  
  - **Solution**: Message-passing via channels (e.g., `chan GenomeUpdate`).  

#### **2. Evolutionary Stability**  
- Avoid "runaway replication" (e.g., organisms evolving infinite `ReplicationRate`).  
  - **Solution**: Penalize high replication with shorter `lifespan`.  

#### **3. Observability**  
- How to visualize evolution in real-time?  
  - **Solution**: Terminal dashboard (BubbleTea) showing:  
    - Population size over time.  
    - Dominant traits (e.g., avg. `CPUHunger`).  

---

### **Hypothetical Outcomes**  
- **Phase 1 (Random Behavior)**: Organisms die quickly due to resource scarcity.  
- **Phase 2 (Adaptation)**: Some evolve optimal `ReplicationRate`/`CPUHunger` ratios.  
- **Phase 3 (Complexity)**: Parasitic strains emerge, triggering "arms races".  

---

### **Non-Goals**  
- No physical-world analogies (e.g., avoid simulating "food" or "predators" literally).  
- No persistence (organisms live only in memory).  

---

### **Why This Is Unique**  
- **Self-Modifying Code**: Organisms’ "DNA" changes runtime behavior.  
- **Concurrency as Evolution**: Goroutines naturally model parallel lifeforms.  
- **Emergent Complexity**: Simple rules → unpredictable outcomes.  

**Next Step**: Define the minimal genome and environment interactions. Want to brainstorm the details?