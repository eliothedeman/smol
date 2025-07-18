# Smol Project Plan

## Project Overview

Smol is a lightweight, modular AI system designed to run within a 2GB memory constraint while providing powerful AI capabilities through a network of small neural networks and tools. The system prioritizes CPU execution and uses PyTorch for computation offloading when needed.

## Core Goals

1. **Memory Efficiency**: Operate within 2GB RAM constraint
2. **Modular Architecture**: Design as interconnected units with clear interfaces
3. **Cross-Language Integration**: Seamless Go-Python communication
4. **Generic Compute**: Use PyTorch for computation without hardware-specific dependencies
5. **Self-Contained**: Minimal external dependencies beyond Python 3.12+ and Go 1.24+

## Architecture Overview

### System Components

#### Core Runtime (Go)
- **Main Process**: Manages unit lifecycle and communication
- **Unit Registry**: Dynamic loading and management of units
- **Message Bus**: Inter-unit communication system
- **Python Bridge**: Embedded Python script execution

#### Required Units

**Control Units:**
- InstructionExecutor: Processes and executes commands
- Lifecycle: Manages unit startup/shutdown sequences

**Tool Units:**
- Math: Mathematical operations and calculations
- Storage: Persistent data storage and retrieval
- Registers: Temporary variable storage and management
- CodeExecution: Sandboxed code execution environment
- Memory: System memory management and optimization
- OpenAI: API server for external AI service integration

**Neural Network Units:**
- ImageClassification: CNN-based image classification
- ObjectDetection: Real-time object detection capabilities

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
- [ ] Set up Go project structure
- [ ] Implement core unit interfaces and registry
- [ ] Create basic message passing system
- [ ] Establish Python integration layer

### Phase 2: Core Units (Weeks 3-4)
- [ ] Implement Control units (InstructionExecutor, Lifecycle)
- [ ] Create basic Tool units (Math, Storage, Registers)
- [ ] Add memory management and monitoring
- [ ] Implement unit testing framework

### Phase 3: Advanced Tools (Weeks 5-6)
- [ ] Build CodeExecution sandbox
- [ ] Implement OpenAI API server
- [ ] Add comprehensive error handling
- [ ] Performance optimization for 2GB constraint

### Phase 4: Neural Networks (Weeks 7-8)
- [ ] Integrate PyTorch for neural network support
- [ ] Implement ImageClassification unit
- [ ] Implement ObjectDetection unit
- [ ] Optimize model loading and memory usage

### Phase 5: Integration & Testing (Weeks 9-10)
- [ ] End-to-end system testing
- [ ] Memory profiling and optimization
- [ ] Documentation and examples
- [ ] Performance benchmarking

## Dependencies & Setup

### System Requirements
- **Go**: 1.24 or higher
- **Python**: 3.12 or higher
- **Memory**: 2GB RAM (hard limit)
- **CPU**: x86_64 or ARM64 architecture

### Python Dependencies
```
torch>=2.0.0
torchvision>=0.15.0
numpy>=1.21.0
pillow>=8.0.0
requests>=2.25.0
```

### Go Dependencies
- Standard library only (no external dependencies)
- Embed package for Python scripts
- Subprocess management for Python execution

## Testing Strategy

### Unit Testing
- Each unit tested in isolation
- Mock dependencies for testing
- Memory usage validation per unit

### Integration Testing
- End-to-end workflow testing
- Memory profiling under load
- Cross-language communication testing

### Performance Testing
- Memory usage monitoring
- CPU utilization tracking
- Response time benchmarking

## Risk Mitigation

### Memory Constraints
- Implement memory monitoring in each unit
- Use streaming for large data processing
- Implement graceful degradation

### Python Integration
- Virtual environment auto-creation
- Dependency validation on startup
- Fallback mechanisms for missing packages

### Model Loading
- Lazy loading of neural network models
- Model quantization for memory efficiency
- Cache management for frequently used models

## Success Criteria

1. **Functional**: All units operate correctly within memory constraints
2. **Performant**: Response times under 1 second for common operations
3. **Stable**: No memory leaks over extended operation
4. **Extensible**: Easy to add new units and capabilities
5. **Documented**: Clear documentation for all components and APIs