# Conway's Game of Life

## Overview
The **Game of Life** is a cellular automaton devised by mathematician John Horton Conway in 1970. It is a zero-player game, meaning its evolution is determined by its initial state with no further input.

## Rules
The universe of the Game of Life is an infinite two-dimensional grid of square cells (after all, it is a _cellular_ automaton), each of which is in one of two possible states:
- **Alive** (■ or `1`)
- **Dead** (⬛ or `0`)

Every cell interacts with its eight neighbors (horizontal, vertical, _and_ diagonal) following these rules:

1. **Underpopulation**: Any live cell with fewer than 2 live neighbors dies.
2. **Survival**: Any live cell with 2 or 3 live neighbors lives on.
3. **Overpopulation**: Any live cell with more than 3 live neighbors dies.
4. **Reproduction**: Any dead cell with exactly 3 live neighbors becomes alive.

## Implementation Pseudocode
```python
for each generation:
    for each cell:
        count = number of live neighbors
        if cell is alive:
            if count < 2 or count > 3:
                cell dies
            else:
                cell lives
        else:
            if count == 3:
                cell becomes alive
```

## Mathematical Properties
- Turing complete (can perform any computation) and thus can simulate any Turing Machine
- Undecidable problem (no algorithm can predict arbitrary patterns) as a corollary of the Halting problem