# fibonacciseq
illustrates intention space PnR computing 
Explanation:

    Intention Class: Represents an intention with a name and a subset of the PnR set.
    Fbsequence Class: Represents the Fibonacci sequence object and includes methods to receive and reflect intentions.
    DesignChunk Class: Represents a design chunk and includes methods to emit and absorb intentions.
    IntentionRing Class: Manages the process of executing the design chunks and intentions in the specified order.
    PnR Set Initialization: Initializes the pnr object with the necessary fields and instances of Fbsequence and DesignChunk. The designChunks array is populated with the design chunk instances.
    Design Chunk Definitions: Defines the actions for starting the sequence and finding the next Fibonacci number.
    Intention Ring Loop Execution: The IntentionRing class manages the loop to ensure the design chunks and intentions are processed in the specified order until all tasks are completed.

This structure provides a model for how the DC1-I1-O-I2-DC2 sequence operates within the intention ring, maintaining clarity and ensuring each component interacts correctly.
